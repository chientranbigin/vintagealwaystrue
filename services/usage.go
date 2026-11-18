package services

import (
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// APIUsageDaily accumulates paid-API usage per day / provider / feature so the shop can see what
// OCR (Google Vision) and AI (OpenAI) actually cost.
type APIUsageDaily struct {
	ID           uint      `gorm:"primaryKey"`
	Day          string    `gorm:"size:10;uniqueIndex:idx_api_usage_key" json:"day"` // YYYY-MM-DD (UTC+7)
	Provider     string    `gorm:"size:32;uniqueIndex:idx_api_usage_key" json:"provider"`
	Feature      string    `gorm:"size:32;uniqueIndex:idx_api_usage_key" json:"feature"`
	Calls        int64     `json:"calls"`
	InputTokens  int64     `json:"input_tokens"`
	OutputTokens int64     `json:"output_tokens"`
	CostUSD      float64   `json:"cost_usd"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (APIUsageDaily) TableName() string { return "api_usage_daily" }

const (
	ProviderOpenAI = "openai"
	ProviderVision = "google_vision"

	// Google Vision TEXT_DETECTION: first 1000 units / month free, then $1.50 per 1000.
	VisionFreeUnitsPerMonth = 1000
	VisionUSDPerUnit        = 1.5 / 1000
)

// USD per 1M tokens (input, output). Unknown models fall back to gpt-4o-mini pricing.
var openAIPricing = map[string][2]float64{
	"gpt-4o-mini":  {0.15, 0.60},
	"gpt-4o":       {2.50, 10.00},
	"gpt-4.1-mini": {0.40, 1.60},
	"gpt-4.1":      {2.00, 8.00},
}

var usageDB *gorm.DB

var ictZone = time.FixedZone("ICT", 7*3600)

// InitUsageTracking creates the usage table and enables recording.
func InitUsageTracking(db *gorm.DB) {
	if db == nil {
		return
	}
	if err := db.AutoMigrate(&APIUsageDaily{}); err != nil {
		log.Printf("[usage] AutoMigrate api_usage_daily error: %v", err)
		return
	}
	usageDB = db
}

// OpenAITokenUsage mirrors the "usage" object of a chat completion response.
type OpenAITokenUsage struct {
	PromptTokens     int64 `json:"prompt_tokens"`
	CompletionTokens int64 `json:"completion_tokens"`
}

func openAICostUSD(model string, in, out int64) float64 {
	price, ok := openAIPricing[model]
	if !ok {
		// Longest matching prefix wins ("gpt-4o-mini-2024-07-18" must not match "gpt-4o").
		best := ""
		for k, v := range openAIPricing {
			if strings.HasPrefix(model, k) && len(k) > len(best) {
				best, price, ok = k, v, true
			}
		}
	}
	if !ok {
		price = openAIPricing["gpt-4o-mini"]
	}
	return float64(in)/1e6*price[0] + float64(out)/1e6*price[1]
}

// RecordOpenAIUsage adds one chat completion to today's counters (async, never blocks the caller).
func RecordOpenAIUsage(feature, model string, u OpenAITokenUsage) {
	recordUsage(ProviderOpenAI, feature, u.PromptTokens, u.CompletionTokens, openAICostUSD(model, u.PromptTokens, u.CompletionTokens))
}

// RecordVisionCall adds one Vision OCR unit (list price; the monthly free tier is applied in reports).
func RecordVisionCall(feature string) {
	recordUsage(ProviderVision, feature, 0, 0, VisionUSDPerUnit)
}

func recordUsage(provider, feature string, in, out int64, cost float64) {
	db := usageDB
	if db == nil {
		return
	}
	if feature == "" {
		feature = "other"
	}
	row := APIUsageDaily{
		Day:          time.Now().In(ictZone).Format("2006-01-02"),
		Provider:     provider,
		Feature:      feature,
		Calls:        1,
		InputTokens:  in,
		OutputTokens: out,
		CostUSD:      cost,
		UpdatedAt:    time.Now(),
	}
	go func() {
		err := db.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "day"}, {Name: "provider"}, {Name: "feature"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"calls":         gorm.Expr("calls + ?", row.Calls),
				"input_tokens":  gorm.Expr("input_tokens + ?", row.InputTokens),
				"output_tokens": gorm.Expr("output_tokens + ?", row.OutputTokens),
				"cost_usd":      gorm.Expr("cost_usd + ?", row.CostUSD),
				"updated_at":    row.UpdatedAt,
			}),
		}).Create(&row).Error
		if err != nil {
			log.Printf("[usage] record %s/%s error: %v", provider, feature, err)
		}
	}()
}

// UsageReport is returned by the usage API.
type UsageReport struct {
	Days           []APIUsageDaily `json:"days"`
	Month          string          `json:"month"`
	MonthOpenAIUSD float64         `json:"month_openai_usd"`
	MonthOpenAI    int64           `json:"month_openai_calls"`
	MonthVision    int64           `json:"month_vision_calls"`
	MonthVisionUSD float64         `json:"month_vision_usd"` // after the free tier
	MonthTotalUSD  float64         `json:"month_total_usd"`
	ProjectedUSD   float64         `json:"projected_month_usd"`
	TrackingSince  string          `json:"tracking_since"`
}

// BuildUsageReport summarises the last `days` days and the current month.
func BuildUsageReport(db *gorm.DB, days int) (*UsageReport, error) {
	if days <= 0 || days > 120 {
		days = 31
	}
	now := time.Now().In(ictZone)
	from := now.AddDate(0, 0, -days+1).Format("2006-01-02")
	rep := &UsageReport{Month: now.Format("2006-01")}
	if err := db.Where("day >= ?", from).Order("day DESC, provider, feature").Find(&rep.Days).Error; err != nil {
		return nil, err
	}
	var first APIUsageDaily
	if err := db.Order("day ASC").Limit(1).Find(&first).Error; err == nil {
		rep.TrackingSince = first.Day
	}
	var monthRows []APIUsageDaily
	if err := db.Where("day >= ?", now.Format("2006-01")+"-01").Find(&monthRows).Error; err != nil {
		return nil, err
	}
	for _, r := range monthRows {
		switch r.Provider {
		case ProviderOpenAI:
			rep.MonthOpenAIUSD += r.CostUSD
			rep.MonthOpenAI += r.Calls
		case ProviderVision:
			rep.MonthVision += r.Calls
		}
	}
	if billable := rep.MonthVision - VisionFreeUnitsPerMonth; billable > 0 {
		rep.MonthVisionUSD = float64(billable) * VisionUSDPerUnit
	}
	rep.MonthTotalUSD = rep.MonthOpenAIUSD + rep.MonthVisionUSD

	// Simple projection from the average daily rate since tracking started this month.
	startDay := 1
	if rep.TrackingSince != "" && strings.HasPrefix(rep.TrackingSince, rep.Month) {
		if t, err := time.ParseInLocation("2006-01-02", rep.TrackingSince, ictZone); err == nil {
			startDay = t.Day()
		}
	}
	elapsed := now.Day() - startDay + 1
	if elapsed < 1 {
		elapsed = 1
	}
	daysInMonth := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, ictZone).Day()
	remaining := daysInMonth - startDay + 1
	openAIProj := rep.MonthOpenAIUSD / float64(elapsed) * float64(remaining)
	visionProj := float64(rep.MonthVision) / float64(elapsed) * float64(remaining)
	visionProjUSD := 0.0
	if b := visionProj - VisionFreeUnitsPerMonth; b > 0 {
		visionProjUSD = b * VisionUSDPerUnit
	}
	rep.ProjectedUSD = openAIProj + visionProjUSD
	return rep, nil
}
