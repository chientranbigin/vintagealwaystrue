package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"sale-app/config"
	"sale-app/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type FBAlbum struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func main() {
	cfg := config.LoadConfig()
	pageID := cfg.FBMainPageID
	pageToken := cfg.FBMainPageToken

	if pageID == "" || pageToken == "" {
		log.Fatalf("❌ Thiếu FB_MAIN_PAGE_ID hoặc FB_MAIN_PAGE_TOKEN trong .env")
	}

	db, err := gorm.Open(mysql.Open(cfg.DatabaseDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	endpoint := fmt.Sprintf("https://graph.facebook.com/v19.0/%s/albums?fields=id,name,count&limit=100&access_token=%s", pageID, url.QueryEscape(pageToken))

	resp, err := client.Get(endpoint)
	if err != nil {
		log.Fatalf("Lỗi gọi API: %v", err)
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	var res struct {
		Data []FBAlbum `json:"data"`
	}
	if err := json.Unmarshal(respBytes, &res); err != nil {
		log.Fatalf("Lỗi parse JSON: %v", err)
	}

	fmt.Println("==================================================")
	fmt.Printf("📋 TÌM THẤY %d ALBUM TRÊN PAGE '%s':\n", len(res.Data), "Vintage Always True House")
	fmt.Println("==================================================")
	for i, alb := range res.Data {
		fmt.Printf(" %2d. ID: %-20s | Tên: %-30s (%d ảnh)\n", i+1, alb.ID, alb.Name, alb.Count)
	}
	fmt.Println("==================================================")

	// Keyword mapping table
	keywordsMap := map[string][]string{
		"TROUSERS":          {"quần", "trouser", "trousers", "quan au", "quần âu"},
		"JACKET":            {"áo khoác", "jacket", "khoác", "ao khoac"},
		"HARRINGTON_JACKET": {"harrington", "harrington jacket"},
		"SHIRT":             {"sơ mi", "shirt", "so mi", "áo sơ mi"},
		"BLAZER":            {"blazer", "áo blazer", "ao blazer"},
		"SUIT":              {"suit", "suits", "bộ suit", "bo suit"},
		"GILE":              {"gile", "vest", "áo gile"},
		"POLO SHIRT":        {"polo", "polo shirt", "áo polo"},
		"TIE":               {"cà vạt", "ca vat", "tie", "ties", "caravat"},
		"BELT":              {"thắt lưng", "that lung", "belt", "dây nịt"},
		"HAT":               {"mũ", "hat", "nón", "non"},
	}

	var categories []models.Category
	db.Find(&categories)

	mappedCount := 0
	for _, cat := range categories {
		kwList, exists := keywordsMap[cat.Code]
		if !exists {
			continue
		}

		// Look for matching album
		var matchedAlbum *FBAlbum
		for _, alb := range res.Data {
			lowerName := strings.ToLower(alb.Name)
			// Skip default system albums
			if lowerName == "ảnh" || lowerName == "tải lên từ di động" || lowerName == "ảnh đại diện" || lowerName == "ảnh bìa" {
				continue
			}

			for _, kw := range kwList {
				if strings.Contains(lowerName, strings.ToLower(kw)) {
					matchedAlbum = &alb
					break
				}
			}
			if matchedAlbum != nil {
				break
			}
		}

		if matchedAlbum != nil {
			db.Model(&models.Category{}).Where("id = ?", cat.ID).Update("fb_album_id", matchedAlbum.ID)
			fmt.Printf("🎯 Tự động map: Danh mục [%-10s - %-18s] ➡️  Album: '%s' (ID: %s)\n", cat.Code, cat.Name, matchedAlbum.Name, matchedAlbum.ID)
			mappedCount++
		}
	}

	fmt.Println("==================================================")
	if mappedCount > 0 {
		fmt.Printf("✅ Đã cập nhật thành công %d danh mục với Album ID tương ứng!\n", mappedCount)
	} else {
		fmt.Println("ℹ️  Chưa tìm thấy Album danh mục phù hợp. Hãy tạo các Album trên Facebook rồi chạy lại lệnh này để hệ thống tự map.")
	}
	fmt.Println("==================================================")
}
