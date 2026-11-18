package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"sale-app/models"

	"gorm.io/gorm"
)

// ProductPrefixes matches ProductKeyTypes in handlers/salev1.go
var copilotProductPrefixes = []string{
	"TS", "HJ", "JK", "ST", "BZ", "TI", "BE", "GL", "PS", "HA", "SU",
}

type CopilotMessage struct {
	Sender      string `json:"sender"` // "customer" or "shop"
	Type        string `json:"type"`   // "text" or "image"
	Text        string `json:"text"`
	ImageBase64 string `json:"image_base64,omitempty"`
	ImageURL    string `json:"image_url,omitempty"`
	OCRText     string `json:"ocr_text,omitempty"`
}

type CopilotIngestRequest struct {
	ThreadID     string           `json:"thread_id"`
	CustomerName string           `json:"customer_name"`
	Action       string           `json:"action"` // "suggest" or "create_order"
	CapturedAt   string           `json:"captured_at"`
	Messages     []CopilotMessage `json:"messages"`
}

type CopilotSizeInfo struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

type CopilotProductMatch struct {
	ID       uint              `json:"id"`
	Code     string            `json:"code"`
	Type     string            `json:"type"`
	Price    int64             `json:"price"`
	Status   string            `json:"status"`
	ThumbURL string            `json:"thumb_url"`
	Sizes    []CopilotSizeInfo `json:"sizes"`
}

type CopilotOrderDraft struct {
	Name         string   `json:"name"`
	Phone        string   `json:"phone"`
	Address      string   `json:"address"`
	Note         string   `json:"note"`
	ProductCodes []string `json:"product_codes"`
}

type CopilotSession struct {
	ThreadID         string                `json:"thread_id"`
	CustomerName     string                `json:"customer_name"`
	AvatarURL        string                `json:"avatar_url,omitempty"`
	Action           string                `json:"action"`
	Messages         []CopilotMessage      `json:"messages"`
	DetectedProducts []CopilotProductMatch `json:"detected_products"`
	SuggestedReply   string                `json:"suggested_reply"`
	SizeAdvice       string                `json:"size_advice"`
	NeedsHuman       bool                  `json:"needs_human"`
	HumanReason      string                `json:"human_reason"`
	OrderDraft       *CopilotOrderDraft    `json:"order_draft,omitempty"`
	CreatedOrder     *models.Order         `json:"created_order,omitempty"`
	UpdatedAt        time.Time             `json:"updated_at"`
}

type CopilotService struct {
	DB         *gorm.DB
	Vision     *VisionService
	SmartOrder *SmartOrderService
	MediaBase  string
	OpenAIKey  string
	Client     *http.Client

	mu          sync.RWMutex
	sessions    map[string]*CopilotSession
	sessionList []string // ordered newest first
	subscribers map[chan *CopilotSession]struct{}
}

func NewCopilotService(db *gorm.DB, vision *VisionService, mediaBase string) *CopilotService {
	if vision == nil {
		vision = NewVisionService()
	}
	svc := &CopilotService{
		DB:          db,
		Vision:      vision,
		SmartOrder:  NewSmartOrderService(),
		MediaBase:   strings.TrimRight(mediaBase, "/"),
		OpenAIKey:   os.Getenv("OPENAI_API_KEY"),
		Client:      &http.Client{Timeout: 45 * time.Second},
		sessions:    make(map[string]*CopilotSession),
		subscribers: make(map[chan *CopilotSession]struct{}),
	}
	svc.initStore()
	return svc
}

func (s *CopilotService) Subscribe() chan *CopilotSession {
	ch := make(chan *CopilotSession, 16)
	s.mu.Lock()
	s.subscribers[ch] = struct{}{}
	s.mu.Unlock()
	return ch
}

func (s *CopilotService) Unsubscribe(ch chan *CopilotSession) {
	s.mu.Lock()
	if _, ok := s.subscribers[ch]; ok {
		delete(s.subscribers, ch)
		close(ch)
	}
	s.mu.Unlock()
}

func (s *CopilotService) broadcast(sess *CopilotSession) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for ch := range s.subscribers {
		select {
		case ch <- sess:
		default:
			// Drop if slow subscriber
		}
	}
}

func (s *CopilotService) SaveAndBroadcast(sess *CopilotSession) {
	s.mu.Lock()
	s.sessions[sess.ThreadID] = sess
	// Move thread_id to front
	newOrder := make([]string, 0, len(s.sessionList)+1)
	newOrder = append(newOrder, sess.ThreadID)
	for _, id := range s.sessionList {
		if id != sess.ThreadID {
			newOrder = append(newOrder, id)
		}
	}
	if len(newOrder) > copilotMemSessions {
		for _, oldID := range newOrder[copilotMemSessions:] {
			delete(s.sessions, oldID)
		}
		newOrder = newOrder[:copilotMemSessions]
	}
	s.sessionList = newOrder
	s.mu.Unlock()

	s.broadcast(sess)
	go s.persistSession(sess)
}

func (s *CopilotService) ListSessions() []*CopilotSession {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]*CopilotSession, 0, len(s.sessionList))
	for _, id := range s.sessionList {
		if sess, ok := s.sessions[id]; ok {
			res = append(res, sess)
		}
	}
	return res
}

func (s *CopilotService) GetSession(threadID string) (*CopilotSession, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sess, ok := s.sessions[threadID]
	return sess, ok
}

// ExtractProductCodesFromText extracts product codes using the same prefix convention as Smart Upload.
func (s *CopilotService) ExtractProductCodesFromText(text string) []string {
	if strings.TrimSpace(text) == "" {
		return nil
	}
	pattern := fmt.Sprintf(`(?i)\b(%s)\s*[-_]?\s*(\d+)\b`, strings.Join(copilotProductPrefixes, "|"))
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(text, -1)

	seen := make(map[string]bool)
	var codes []string
	for _, m := range matches {
		if len(m) >= 3 {
			code := strings.ToUpper(m[1] + m[2])
			if !seen[code] {
				seen[code] = true
				codes = append(codes, code)
			}
		}
	}
	return codes
}

// AnalyzeIngest runs OCR on incoming images, matches DB products, and generates AI reply / order draft.
func (s *CopilotService) AnalyzeIngest(ctx context.Context, req CopilotIngestRequest) *CopilotSession {
	threadID := strings.TrimSpace(req.ThreadID)
	if threadID == "" {
		threadID = req.CustomerName
	}
	if threadID == "" {
		threadID = fmt.Sprintf("thread-%d", time.Now().Unix())
	}

	customerName := strings.TrimSpace(req.CustomerName)
	if customerName == "" {
		customerName = "Khách Facebook"
	}

	// Preserve created_order if same session already created an order recently
	var prevCreatedOrder *models.Order
	var history []CopilotMessage
	if prev, ok := s.loadSession(threadID); ok {
		history = prev.Messages
		if prev.CreatedOrder != nil && req.Action != "create_order" {
			prevCreatedOrder = prev.CreatedOrder
		}
	}

	// 1. Run OCR on images (matching Smart Upload's Google Vision pipeline) & collect all codes
	var allCodes []string
	seenCode := make(map[string]bool)
	addCodes := func(codes []string) {
		for _, c := range codes {
			if !seenCode[c] {
				seenCode[c] = true
				allCodes = append(allCodes, c)
			}
		}
	}

	for i := range req.Messages {
		msg := &req.Messages[i]
		if msg.Type == "text" && msg.Text != "" {
			addCodes(s.ExtractProductCodesFromText(msg.Text))
		}
		if msg.Type == "image" && msg.ImageBase64 != "" && s.Vision != nil {
			if ocrText, err := s.Vision.ExtractTextFromBase64(msg.ImageBase64); err == nil {
				msg.OCRText = strings.TrimSpace(ocrText)
				addCodes(s.ExtractProductCodesFromText(ocrText))
			}
		}
	}

	// Merge the snapshot into the stored full history (extension only sends the latest ~30 messages).
	merged := mergeCopilotMessages(history, req.Messages)

	// AI context: "suggest" sees the longer history; "create_order" keeps the extension's
	// slice (messages since the last order trigger) so old orders don't leak into the new one.
	aiMessages := tailMessages(merged, copilotAIContextMessages)
	if req.Action == "create_order" {
		aiMessages = req.Messages
	}

	// 2. Lookup matched products in DB
	detectedProducts := s.lookupProductsByCodes(allCodes)

	// 3. Check if customer mentioned height/weight and recommend available items if no product detected yet
	convText := s.buildConversationTranscript(aiMessages)
	recommendedProducts := s.recommendProductsByMeasurements(convText, detectedProducts)
	allMatched := append(detectedProducts, recommendedProducts...)

	// 4. Generate AI Suggestion & Order Draft
	aiRes := s.generateAISuggestion(ctx, customerName, req.Action, aiMessages, allMatched)

	sess := &CopilotSession{
		ThreadID:         threadID,
		CustomerName:     customerName,
		Action:           req.Action,
		Messages:         merged,
		DetectedProducts: allMatched,
		SuggestedReply:   aiRes.SuggestedReply,
		SizeAdvice:       aiRes.SizeAdvice,
		NeedsHuman:       aiRes.NeedsHuman,
		HumanReason:      aiRes.HumanReason,
		OrderDraft:       aiRes.OrderDraft,
		CreatedOrder:     prevCreatedOrder,
		UpdatedAt:        time.Now(),
	}

	return sess
}

func (s *CopilotService) lookupProductsByCodes(codes []string) []CopilotProductMatch {
	if len(codes) == 0 || s.DB == nil {
		return nil
	}
	var out []CopilotProductMatch
	seenID := make(map[uint]bool)

	for _, code := range codes {
		var prod models.Product
		if err := s.DB.Where("UPPER(name) = ?", strings.ToUpper(code)).
			Preload("Images").
			Preload("Sizes").
			Order("id DESC").
			First(&prod).Error; err != nil {
			// Fallback LIKE search if exact match not found
			if err2 := s.DB.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(code)+"%").
				Preload("Images").
				Preload("Sizes").
				Order("id DESC").
				First(&prod).Error; err2 != nil {
				continue
			}
		}
		if seenID[prod.ID] {
			continue
		}
		seenID[prod.ID] = true
		out = append(out, s.toProductMatch(prod))
	}
	return out
}

func (s *CopilotService) toProductMatch(p models.Product) CopilotProductMatch {
	// Prefer the scaled thumb (ImageThumbScale), then PathThumb, then gallery images.
	thumb := p.DisplayThumb()
	if thumb == "" {
		thumb = p.BestQualityThumb()
	}
	thumb = s.formatMediaPath(thumb)
	var sizes []CopilotSizeInfo
	for _, sz := range p.Sizes {
		if sz.Name != nil && sz.Value != nil {
			sizes = append(sizes, CopilotSizeInfo{
				Name:  *sz.Name,
				Value: *sz.Value,
			})
		}
	}
	return CopilotProductMatch{
		ID:       p.ID,
		Code:     strPtrVal(p.Name),
		Type:     strPtrVal(p.Type),
		Price:    int64PtrVal(p.Price),
		Status:   strPtrVal(p.Status),
		ThumbURL: thumb,
		Sizes:    sizes,
	}
}

func (s *CopilotService) recommendProductsByMeasurements(convText string, existing []CopilotProductMatch) []CopilotProductMatch {
	if len(existing) > 0 || s.DB == nil {
		return nil
	}
	heightCm, weightKg := extractHeightWeight(convText)
	if heightCm == 0 && weightKg == 0 {
		return nil
	}
	// Estimate rough chest & waist in cm for vintage sizing suggestions
	estEo := estimateWaistCm(heightCm, weightKg)
	if estEo <= 0 {
		return nil
	}

	var prods []models.Product
	if err := s.DB.Where("status = 'AVAILABLE'").
		Preload("Images").
		Preload("Sizes").
		Order("id DESC").
		Limit(60).
		Find(&prods).Error; err != nil {
		return nil
	}

	type scored struct {
		prod models.Product
		diff float64
	}
	var candidates []scored
	for _, p := range prods {
		for _, sz := range p.Sizes {
			if sz.Name == nil || sz.Value == nil {
				continue
			}
			n := strings.ToUpper(*sz.Name)
			if strings.Contains(n, "EO") {
				diff := absFloat(*sz.Value - estEo)
				if diff <= 4.0 {
					candidates = append(candidates, scored{prod: p, diff: diff})
				}
				break
			}
		}
	}
	sort.Slice(candidates, func(i, j int) bool { return candidates[i].diff < candidates[j].diff })
	limit := 3
	if len(candidates) < limit {
		limit = len(candidates)
	}
	var out []CopilotProductMatch
	for i := 0; i < limit; i++ {
		out = append(out, s.toProductMatch(candidates[i].prod))
	}
	return out
}

type aiCopilotOutput struct {
	SuggestedReply string             `json:"suggested_reply"`
	SizeAdvice     string             `json:"size_advice"`
	NeedsHuman     bool               `json:"needs_human"`
	HumanReason    string             `json:"human_reason"`
	OrderDraft     *CopilotOrderDraft `json:"order_draft,omitempty"`
}

func (s *CopilotService) generateAISuggestion(
	ctx context.Context,
	customerName string,
	action string,
	messages []CopilotMessage,
	products []CopilotProductMatch,
) aiCopilotOutput {
	transcript := s.buildConversationTranscript(messages)

	// Fast fallback if OpenAI key is missing
	if strings.TrimSpace(s.OpenAIKey) == "" {
		return s.fallbackSuggestion(customerName, action, transcript, products)
	}

	prodJSON, _ := json.Marshal(products)
	systemPrompt := `Bạn là trợ lý bán hàng cho shop thời trang Vintage Always True.
Nhiệm vụ: Đọc hội thoại giữa Khách và Shop cùng danh sách sản phẩm khớp từ Database (kèm giá, trạng thái AVAILABLE/SOLD, và số đo chi tiết VAI/NGỰC/EO/DÀI...) để:
1. Soạn "suggested_reply": câu trả lời tự nhiên, lịch sự bằng tiếng Việt (xưng "em" - gọi khách "mình/anh/chị") để chủ shop copy dán qua Facebook.
   - Chỉ trả lời (các) tin nhắn MỚI NHẤT của Khách mà Shop chưa trả lời. Ngắn gọn 1-4 câu, KHÔNG nhắc lại thông tin Shop/AI đã trả lời ở các lượt trước.
   - Chỉ dùng thông tin có trong dữ liệu sản phẩm và "Chính sách shop" bên dưới. TUYỆT ĐỐI KHÔNG tự bịa phí ship, thời gian giao, chính sách đổi trả, khuyến mãi. Nếu khách hỏi điều không có dữ liệu: đặt "needs_human": true và trả lời kiểu "Dạ để em kiểm tra và báo lại mình ngay ạ".
   - Nếu khách đã chốt mua và gửi đủ tên/SĐT/địa chỉ: xác nhận lại ngắn gọn mã + giá + thông tin nhận hàng.
   - Nếu khách hỏi còn hàng: báo rõ mã nào còn (kèm giá + số đo nổi bật), mã nào đã bán.
   - Nếu khách hỏi có vừa chiều cao/cân nặng không: phân tích dựa trên số đo thực tế trong DB (lưu ý đồ vintage thường form rộng/ôm tùy loại, ngực áo thường lớn hơn ngực cơ thể 10-16cm, eo quần x2 nếu số đo ghi nửa vòng hoặc đúng cm nếu ghi cả vòng).
   - Nếu khách mặc cả / xin giảm giá / đổi trả / khiếu nại: đặt "needs_human": true, ghi rõ "human_reason", và gợi ý câu trả lời khéo đợi nhân viên kiểm tra.
2. Ghi ngắn gọn nhận xét vào "size_advice" (nếu có hỏi size/số đo).
3. Trích xuất "order_draft": lấy tên khách, SĐT, địa chỉ, ghi chú và danh sách mã sản phẩm đang AVAILABLE được nhắc tới trong hội thoại (TUYỆT ĐỐI KHÔNG đưa mã đã SOLD vào order_draft).

Trả về JSON thuần theo schema:
{
  "suggested_reply": "string",
  "size_advice": "string",
  "needs_human": false,
  "human_reason": "string",
  "order_draft": {
    "name": "string",
    "phone": "string",
    "address": "string",
    "note": "string",
    "product_codes": ["TS123"]
  }
}`

	policy := strings.TrimSpace(os.Getenv("COPILOT_SHOP_POLICY"))
	if policy == "" {
		policy = "(Chưa cấu hình — không được tự đưa ra phí ship / thời gian giao / đổi trả)"
	}
	userPrompt := fmt.Sprintf(
		"Tên khách Facebook: %s\nAction hiện tại: %s\n\nChính sách shop:\n%s\n\nSản phẩm tìm thấy trong DB:\n%s\n\nHội thoại hiện tại:\n%s",
		customerName, action, policy, string(prodJSON), transcript,
	)

	payload := map[string]interface{}{
		"model": "gpt-4o-mini",
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userPrompt},
		},
		"response_format": map[string]string{"type": "json_object"},
		"temperature":     0.3,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return s.fallbackSuggestion(customerName, action, transcript, products)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return s.fallbackSuggestion(customerName, action, transcript, products)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.OpenAIKey)

	resp, err := s.Client.Do(req)
	if err != nil {
		return s.fallbackSuggestion(customerName, action, transcript, products)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return s.fallbackSuggestion(customerName, action, transcript, products)
	}

	var openAIResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBytes, &openAIResp); err != nil || len(openAIResp.Choices) == 0 {
		return s.fallbackSuggestion(customerName, action, transcript, products)
	}

	var out aiCopilotOutput
	if err := json.Unmarshal([]byte(openAIResp.Choices[0].Message.Content), &out); err != nil {
		return s.fallbackSuggestion(customerName, action, transcript, products)
	}
	if out.OrderDraft != nil && out.OrderDraft.Name == "" {
		out.OrderDraft.Name = customerName
	}
	return out
}

func (s *CopilotService) fallbackSuggestion(customerName, action, transcript string, products []CopilotProductMatch) aiCopilotOutput {
	parsed := s.SmartOrder.FallbackParse(transcript)
	if parsed.Name == "" {
		parsed.Name = customerName
	}

	var availCodes []string
	var replyParts []string
	for _, p := range products {
		if strings.EqualFold(p.Status, "AVAILABLE") {
			availCodes = append(availCodes, p.Code)
			replyParts = append(replyParts, fmt.Sprintf("Mã %s còn hàng, giá %sđ ạ", p.Code, FormatVNDPrice(p.Price)))
		} else {
			replyParts = append(replyParts, fmt.Sprintf("Mã %s bên em vừa bán mất rồi ạ", p.Code))
		}
	}

	lower := strings.ToLower(transcript)
	needsHuman := strings.Contains(lower, "giảm giá") || strings.Contains(lower, "bớt") || strings.Contains(lower, "fix") || strings.Contains(lower, "đổi trả")
	humanReason := ""
	if needsHuman {
		humanReason = "Khách hỏi giảm giá / đổi trả"
	}

	reply := strings.Join(replyParts, ". ")
	if reply == "" {
		reply = "Dạ mình gửi giúp em mã sản phẩm hoặc chiều cao cân nặng để em kiểm tra cho mình nha ạ."
	}

	return aiCopilotOutput{
		SuggestedReply: reply,
		NeedsHuman:     needsHuman,
		HumanReason:    humanReason,
		OrderDraft: &CopilotOrderDraft{
			Name:         parsed.Name,
			Phone:        parsed.Phone,
			Address:      parsed.Address,
			Note:         parsed.Note,
			ProductCodes: availCodes,
		},
	}
}

func (s *CopilotService) buildConversationTranscript(messages []CopilotMessage) string {
	var lines []string
	for _, m := range messages {
		role := "Khách"
		if m.Sender == "shop" {
			role = "Shop"
		}
		if m.Type == "image" {
			if m.OCRText != "" {
				lines = append(lines, fmt.Sprintf("%s: [Gửi ảnh - OCR đọc được: %s]", role, m.OCRText))
			} else {
				lines = append(lines, fmt.Sprintf("%s: [Gửi ảnh]", role))
			}
		} else if m.Text != "" {
			lines = append(lines, fmt.Sprintf("%s: %s", role, m.Text))
		}
	}
	return strings.Join(lines, "\n")
}

func extractHeightWeight(text string) (int, int) {
	var h, w int
	reH1 := regexp.MustCompile(`(?i)\b1\s*m\s*(\d{2})\b`)
	if m := reH1.FindStringSubmatch(text); len(m) >= 2 {
		if v, err := strconv.Atoi(m[1]); err == nil {
			h = 100 + v
		}
	} else {
		reH2 := regexp.MustCompile(`(?i)\b(1[4-9]\d)\s*cm\b`)
		if m := reH2.FindStringSubmatch(text); len(m) >= 2 {
			h, _ = strconv.Atoi(m[1])
		}
	}
	reW := regexp.MustCompile(`(?i)\b(\d{2,3})\s*(kg|kí|ký|cân)\b`)
	if m := reW.FindStringSubmatch(text); len(m) >= 2 {
		w, _ = strconv.Atoi(m[1])
	}
	return h, w
}

func estimateWaistCm(heightCm, weightKg int) float64 {
	if weightKg <= 0 {
		return 0
	}
	// Rough waist heuristic for Asian men/unisex sizing
	base := float64(weightKg) * 1.25
	if heightCm > 175 {
		base -= 2.0
	}
	return base
}

func absFloat(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}

func strPtrVal(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func int64PtrVal(p *int64) int64 {
	if p == nil {
		return 0
	}
	return *p
}

// formatMediaPath mirrors SaleV2Handler.FormatImagePath so Copilot thumbnails resolve
// exactly like the rest of SaleV2 (ProductVault, OrderCreate...).
func (s *CopilotService) formatMediaPath(path string) string {
	if path == "" {
		return ""
	}
	if strings.Contains(path, "tadimedia.vn") {
		path = strings.Replace(path, "www.tadimedia.vn", "vintagealwaystrue.id.vn", 1)
		path = strings.Replace(path, "tadimedia.vn", "vintagealwaystrue.id.vn", 1)
	}
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}
	clean := strings.TrimPrefix(strings.TrimLeft(path, "/"), "storage/")
	if clean == "" || s.MediaBase == "" {
		return clean
	}
	return s.MediaBase + "/" + clean
}

// FetchFBUserProfile queries Facebook Graph API for the real customer name and profile picture.
func (s *CopilotService) FetchFBUserProfile(psid, pageToken string) (string, string) {
	if psid == "" || pageToken == "" {
		return "", ""
	}
	endpoint := fmt.Sprintf("https://graph.facebook.com/v21.0/%s?fields=name,profile_pic&access_token=%s", url.PathEscape(psid), url.QueryEscape(pageToken))
	resp, err := s.Client.Get(endpoint)
	if err != nil {
		return "", ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", ""
	}
	var out struct {
		Name       string `json:"name"`
		ProfilePic string `json:"profile_pic"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", ""
	}
	return strings.TrimSpace(out.Name), strings.TrimSpace(out.ProfilePic)
}

// UpsertFBMessage persists an incoming webhook/Graph API message to DB (deduplicating by fb_message_id)
// and runs Vision OCR if it's an image.
func (s *CopilotService) UpsertFBMessage(
	psid, customerName, avatarURL, fbMsgID, sender, msgType, text, imageURL string,
	msgTime time.Time,
	pageToken string,
) (*models.CopilotConversation, bool, error) {
	if s.DB == nil {
		return nil, false, fmt.Errorf("database not initialized")
	}
	if strings.TrimSpace(psid) == "" {
		return nil, false, fmt.Errorf("empty psid")
	}
	if msgTime.IsZero() {
		msgTime = time.Now()
	}

	var conv models.CopilotConversation
	err := s.DB.Where(&models.CopilotConversation{PSID: psid}).First(&conv).Error
	if err != nil {
		if customerName == "" && pageToken != "" {
			name, pic := s.FetchFBUserProfile(psid, pageToken)
			if name != "" {
				customerName = name
			}
			if avatarURL == "" && pic != "" {
				avatarURL = pic
			}
		}
		if customerName == "" {
			customerName = "Khách Facebook"
		}
		conv = models.CopilotConversation{
			PSID:         psid,
			CustomerName: customerName,
			AvatarURL:    avatarURL,
		}
		if err := s.DB.Create(&conv).Error; err != nil {
			return nil, false, err
		}
	} else {
		updates := map[string]interface{}{}
		if (conv.CustomerName == "" || conv.CustomerName == "Khách Facebook") && pageToken != "" && customerName == "" {
			name, pic := s.FetchFBUserProfile(psid, pageToken)
			if name != "" {
				customerName = name
			}
			if pic != "" {
				avatarURL = pic
			}
		}
		if customerName != "" && customerName != "Khách Facebook" && conv.CustomerName != customerName {
			updates["customer_name"] = customerName
			conv.CustomerName = customerName
		}
		if avatarURL != "" && conv.AvatarURL != avatarURL {
			updates["avatar_url"] = avatarURL
			conv.AvatarURL = avatarURL
		}
		if len(updates) > 0 {
			s.DB.Model(&conv).Updates(updates)
		}
	}

	if fbMsgID != "" {
		var existing models.CopilotMessageRecord
		if err := s.DB.Where("fb_message_id = ?", fbMsgID).First(&existing).Error; err == nil {
			// Already processed; fix records stored with an older classification.
			if existing.Type != msgType || existing.Text != text {
				s.DB.Model(&existing).Updates(map[string]interface{}{"type": msgType, "text": text, "image_url": imageURL})
			}
			return &conv, false, nil
		}
	}

	var ocrText string
	if msgType == "image" && imageURL != "" && s.Vision != nil {
		if txt, err := s.Vision.ExtractTextFromImage(imageURL); err == nil {
			ocrText = strings.TrimSpace(txt)
		}
	}

	rec := models.CopilotMessageRecord{
		ConversationID: conv.ID,
		FBMessageID:    fbMsgID,
		Sender:         sender,
		Type:           msgType,
		Text:           text,
		ImageURL:       imageURL,
		OCRText:        ocrText,
		MsgTimestamp:   msgTime,
	}
	if err := s.DB.Create(&rec).Error; err != nil {
		return &conv, false, err
	}
	s.DB.Model(&conv).Update("updated_at", time.Now())
	return &conv, true, nil
}

// RebuildAndAnalyzeSession loads messages after LastOrderAt from DB, runs AI & DB matching, and returns a CopilotSession.
func (s *CopilotService) RebuildAndAnalyzeSession(ctx context.Context, psid string, action string) (*CopilotSession, error) {
	if s.DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	if strings.TrimSpace(psid) == "" {
		return nil, fmt.Errorf("empty psid")
	}
	var conv models.CopilotConversation
	if err := s.DB.Where(&models.CopilotConversation{PSID: psid}).First(&conv).Error; err != nil {
		return nil, err
	}

	query := s.DB.Where("conversation_id = ?", conv.ID)
	if conv.LastOrderAt != nil {
		query = query.Where("msg_timestamp > ?", *conv.LastOrderAt)
	}
	// Newest 30 messages, then flip back to chronological order.
	var recs []models.CopilotMessageRecord
	if err := query.Order("msg_timestamp DESC, id DESC").Limit(30).Find(&recs).Error; err != nil {
		return nil, err
	}
	for i, j := 0, len(recs)-1; i < j; i, j = i+1, j-1 {
		recs[i], recs[j] = recs[j], recs[i]
	}

	var msgs []CopilotMessage
	var allCodes []string
	seenCode := make(map[string]bool)
	addCodes := func(codes []string) {
		for _, c := range codes {
			if !seenCode[c] {
				seenCode[c] = true
				allCodes = append(allCodes, c)
			}
		}
	}

	for _, r := range recs {
		if r.Type == "text" && IsFBSystemNotice(r.Text) {
			continue
		}
		msgs = append(msgs, CopilotMessage{
			Sender:   r.Sender,
			Type:     r.Type,
			Text:     r.Text,
			ImageURL: r.ImageURL,
			OCRText:  r.OCRText,
		})
		if r.Type == "text" && r.Text != "" {
			addCodes(s.ExtractProductCodesFromText(r.Text))
		}
		if r.Type == "image" && r.OCRText != "" {
			addCodes(s.ExtractProductCodesFromText(r.OCRText))
		}
	}

	detectedProducts := s.lookupProductsByCodes(allCodes)
	convText := s.buildConversationTranscript(msgs)
	recommendedProducts := s.recommendProductsByMeasurements(convText, detectedProducts)
	allMatched := append(detectedProducts, recommendedProducts...)

	aiRes := s.generateAISuggestion(ctx, conv.CustomerName, action, msgs, allMatched)

	var createdOrder *models.Order
	if conv.LastOrderID != nil {
		var ord models.Order
		if err := s.DB.Preload("Products").First(&ord, *conv.LastOrderID).Error; err == nil {
			s.FormatOrderMedia(&ord)
			createdOrder = &ord
		}
	}

	sess := &CopilotSession{
		ThreadID:         conv.PSID,
		CustomerName:     conv.CustomerName,
		AvatarURL:        conv.AvatarURL,
		Action:           action,
		Messages:         msgs,
		DetectedProducts: allMatched,
		SuggestedReply:   aiRes.SuggestedReply,
		SizeAdvice:       aiRes.SizeAdvice,
		NeedsHuman:       aiRes.NeedsHuman,
		HumanReason:      aiRes.HumanReason,
		OrderDraft:       aiRes.OrderDraft,
		CreatedOrder:     createdOrder,
		UpdatedAt:        conv.UpdatedAt,
	}
	return sess, nil
}

func (s *CopilotService) MarkConversationOrdered(psid string, orderID uint) {
	if s.DB == nil || strings.TrimSpace(psid) == "" {
		return
	}
	now := time.Now()
	s.DB.Model(&models.CopilotConversation{}).
		Where(&models.CopilotConversation{PSID: psid}).
		Updates(map[string]interface{}{
			"last_order_at": now,
			"last_order_id": orderID,
			"updated_at":    now,
		})
}

// LoadRecentSessionsFromDB hydrates in-memory sessions from DB after server restart.
func (s *CopilotService) LoadRecentSessionsFromDB(ctx context.Context) []*CopilotSession {
	if s.DB == nil {
		return nil
	}
	var convs []models.CopilotConversation
	if err := s.DB.Order("updated_at DESC").Limit(15).Find(&convs).Error; err != nil {
		return nil
	}
	var out []*CopilotSession
	for _, conv := range convs {
		if sess, err := s.RebuildAndAnalyzeSession(ctx, conv.PSID, "suggest"); err == nil {
			out = append(out, sess)
			s.mu.Lock()
			s.sessions[sess.ThreadID] = sess
			s.sessionList = append(s.sessionList, sess.ThreadID)
			s.mu.Unlock()
		}
	}
	return out
}

type graphConversation struct {
	ID           string `json:"id"`
	Participants struct {
		Data []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"data"`
	} `json:"participants"`
	Messages struct {
		Data []struct {
			ID          string `json:"id"`
			Message     string `json:"message"`
			CreatedTime string `json:"created_time"`
			From        struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"from"`
			Attachments struct {
				Data []struct {
					ID        string `json:"id"`
					Name      string `json:"name"`
					MimeType  string `json:"mime_type"`
					ImageData struct {
						URL string `json:"url"`
					} `json:"image_data"`
				} `json:"data"`
			} `json:"attachments"`
		} `json:"data"`
	} `json:"messages"`
}

// fetchGraphConversations calls /{page}/conversations. extraQuery is appended as-is
// (e.g. "&user_id=<PSID>" to fetch one customer's thread).
func (s *CopilotService) fetchGraphConversations(ctx context.Context, pageID, pageToken string, msgLimit, convLimit int, extraQuery string) ([]graphConversation, error) {
	if pageID == "" || pageToken == "" {
		return nil, fmt.Errorf("thiếu FB_MAIN_PAGE_ID hoặc FB_MAIN_PAGE_TOKEN")
	}
	fields := fmt.Sprintf("participants,updated_time,messages.limit(%d){id,message,from,attachments,created_time}", msgLimit)
	endpoint := fmt.Sprintf(
		"https://graph.facebook.com/v21.0/%s/conversations?platform=messenger&fields=%s&limit=%d%s&access_token=%s",
		url.PathEscape(pageID),
		url.QueryEscape(fields),
		convLimit,
		extraQuery,
		url.QueryEscape(pageToken),
	)
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Graph API error (%d): %s", resp.StatusCode, string(bodyBytes))
	}
	var parsed struct {
		Data []graphConversation `json:"data"`
	}
	if err := json.Unmarshal(bodyBytes, &parsed); err != nil {
		return nil, err
	}
	return parsed.Data, nil
}

// importGraphConversation stores one Graph conversation's messages (deduplicated by message id)
// and returns the customer's PSID.
func (s *CopilotService) importGraphConversation(conv graphConversation, pageID, pageToken string) string {
	var customerPSID, customerName string
	for _, p := range conv.Participants.Data {
		if p.ID != pageID {
			customerPSID = p.ID
			customerName = p.Name
			break
		}
	}
	if customerPSID == "" {
		return ""
	}

	// Messages from Graph API are newest-first; iterate backwards
	msgs := conv.Messages.Data
	for i := len(msgs) - 1; i >= 0; i-- {
		m := msgs[i]
		sender := "customer"
		if m.From.ID == pageID {
			sender = "shop"
		}
		ts, _ := time.Parse("2006-01-02T15:04:05-0700", m.CreatedTime)
		if ts.IsZero() {
			ts = time.Now()
		}
		if IsFBSystemNotice(m.Message) {
			// "<Tên> đã trả lời một quảng cáo." is a Meta banner, not a real shop message.
			s.DeleteFBMessage(m.ID)
		} else if m.Message != "" {
			_, _, _ = s.UpsertFBMessage(customerPSID, customerName, "", m.ID, sender, "text", m.Message, "", ts, pageToken)
		}
		for idx, att := range m.Attachments.Data {
			attID := fmt.Sprintf("%s_att_%d", m.ID, idx)
			if IsLikeStickerAttachment(att.ID, att.Name) {
				_, _, _ = s.UpsertFBMessage(customerPSID, customerName, "", attID, sender, "text", "👍", "", ts, pageToken)
				continue
			}
			if att.ImageData.URL != "" {
				_, _, _ = s.UpsertFBMessage(customerPSID, customerName, "", attID, sender, "image", "", att.ImageData.URL, ts, pageToken)
			}
		}
	}
	return customerPSID
}

// SyncRecentFBConversations fetches the latest Messenger threads via Graph API and stores them in DB.
func (s *CopilotService) SyncRecentFBConversations(ctx context.Context, pageID, pageToken string) (int, error) {
	convs, err := s.fetchGraphConversations(ctx, pageID, pageToken, 15, 10, "")
	if err != nil {
		return 0, err
	}
	synced := 0
	for _, conv := range convs {
		psid := s.importGraphConversation(conv, pageID, pageToken)
		if psid == "" {
			continue
		}
		if sess, err := s.RebuildAndAnalyzeSession(ctx, psid, "suggest"); err == nil {
			s.SaveAndBroadcast(sess)
			synced++
		}
	}
	return synced, nil
}

// HasFBConversation reports whether this PSID already has a stored conversation.
func (s *CopilotService) HasFBConversation(psid string) bool {
	if s.DB == nil || strings.TrimSpace(psid) == "" {
		return false
	}
	var n int64
	s.DB.Model(&models.CopilotConversation{}).Where(&models.CopilotConversation{PSID: psid}).Count(&n)
	return n > 0
}

// ImportFBHistory pulls the last 30 messages of one customer's thread so a customer first seen
// via webhook starts with full context (not just the message that triggered the webhook).
func (s *CopilotService) ImportFBHistory(ctx context.Context, pageID, psid, pageToken string) error {
	if strings.TrimSpace(psid) == "" {
		return fmt.Errorf("empty psid")
	}
	convs, err := s.fetchGraphConversations(ctx, pageID, pageToken, 30, 1, "&user_id="+url.QueryEscape(psid))
	if err != nil {
		return err
	}
	for _, conv := range convs {
		s.importGraphConversation(conv, pageID, pageToken)
	}
	return nil
}

var fbSystemNoticeRe = regexp.MustCompile(`(?i)^.{1,80}?\s(đã trả lời (về )?(một |1 )?(quảng cáo|bài viết)|replied to (an|a|your) (ad|post))\.?(\s|\(|$)`)

// IsFBSystemNotice reports Meta-generated banner lines that arrive as page messages.
func IsFBSystemNotice(text string) bool {
	t := strings.TrimSpace(text)
	return t != "" && fbSystemNoticeRe.MatchString(t)
}

// Facebook "Like" (thumbs up) sticker ids in small / medium / large.
var fbLikeStickerIDs = map[string]bool{"369239263222822": true, "369239343222814": true, "369239383222810": true}

func IsLikeStickerID(id int64) bool {
	return fbLikeStickerIDs[strconv.FormatInt(id, 10)]
}

// IsLikeStickerAttachment matches Graph API attachments like {"id":"sticker_369239263222822"}.
func IsLikeStickerAttachment(id, name string) bool {
	for _, v := range []string{id, name} {
		v = strings.TrimPrefix(strings.TrimPrefix(v, "sticker_"), "sticker-")
		if fbLikeStickerIDs[v] {
			return true
		}
	}
	return false
}

func (s *CopilotService) DeleteFBMessage(fbMsgID string) {
	if s.DB == nil || fbMsgID == "" {
		return
	}
	s.DB.Where("fb_message_id = ?", fbMsgID).Delete(&models.CopilotMessageRecord{})
}

// FormatOrderMedia resolves product image paths inside an order to absolute URLs (scale image
// first, like SaleV2's formatProduct), so the Copilot receipt can render them.
func (s *CopilotService) FormatOrderMedia(order *models.Order) {
	if order == nil {
		return
	}
	for i := range order.Products {
		p := &order.Products[i]
		if p.PathThumb != nil && *p.PathThumb != "" {
			v := s.formatMediaPath(*p.PathThumb)
			p.PathThumb = &v
		}
		if p.ImageThumbScale != nil && *p.ImageThumbScale != "" {
			v := s.formatMediaPath(*p.ImageThumbScale)
			p.ImageThumbScale = &v
			p.ImageThumbScaleURL = &v
		}
		if d := s.formatMediaPath(p.DisplayThumb()); d != "" {
			p.DisplayThumbURL = &d
		}
	}
}

// SendFBImage uploads an image and sends it to a Messenger user (PSID) via the Page Send API.
func (s *CopilotService) SendFBImage(ctx context.Context, psid string, img []byte, filename, pageToken string) error {
	if psid == "" || pageToken == "" {
		return fmt.Errorf("thiếu PSID hoặc FB_MAIN_PAGE_TOKEN")
	}
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	_ = w.WriteField("recipient", fmt.Sprintf(`{"id":%q}`, psid))
	_ = w.WriteField("messaging_type", "RESPONSE")
	_ = w.WriteField("message", `{"attachment":{"type":"image","payload":{"is_reusable":false}}}`)
	hdr := make(textproto.MIMEHeader)
	hdr.Set("Content-Disposition", fmt.Sprintf(`form-data; name="filedata"; filename=%q`, filename))
	hdr.Set("Content-Type", "image/png")
	part, err := w.CreatePart(hdr)
	if err != nil {
		return err
	}
	if _, err := part.Write(img); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	endpoint := "https://graph.facebook.com/v21.0/me/messages?access_token=" + url.QueryEscape(pageToken)
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		var fbErr struct {
			Error struct {
				Message string `json:"message"`
				Code    int    `json:"code"`
			} `json:"error"`
		}
		if json.Unmarshal(body, &fbErr) == nil && fbErr.Error.Message != "" {
			return fmt.Errorf("Facebook (%d): %s", fbErr.Error.Code, fbErr.Error.Message)
		}
		return fmt.Errorf("Facebook HTTP %d", resp.StatusCode)
	}
	return nil
}
