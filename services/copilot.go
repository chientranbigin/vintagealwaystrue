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
	Timestamp   string `json:"timestamp,omitempty"`
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
	// Images are original-resolution photos: main first, then detail shots.
	Images []string `json:"images,omitempty"`
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
	// LastMessageAt is the time of the newest message; the Copilot list is sorted by it.
	LastMessageAt time.Time `json:"last_message_at"`
	// AIStale: messages changed since the AI fields were computed (e.g. reloaded after a restart).
	AIStale bool `json:"ai_stale,omitempty"`
	// AIGenerated: suggestion came from OpenAI (button). false = free rule-based suggestion.
	AIGenerated bool `json:"ai_generated"`
	// AIForKey: last message the AI suggestion was computed for (kept until a new message).
	AIForKey string `json:"ai_for_key,omitempty"`
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
	vision.Feature = "copilot"
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
	sort.SliceStable(res, func(i, j int) bool {
		return res[i].SortTime().After(res[j].SortTime())
	})
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
		LastMessageAt:    time.Now(),
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
		Images:   s.originalImages(p),
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
		return s.fallbackSuggestion(customerName, action, transcript, messages, products)
	}

	pending := s.pendingCustomerCodes(messages)
	productContext, soldCodes := buildProductContext(products, pending)
	systemPrompt := `Bạn là trợ lý bán hàng của shop đồ vintage nam "Vintage Always True", soạn sẵn tin nhắn để chủ shop gửi khách qua Messenger.

DỮ LIỆU SẢN PHẨM được chia 2 nhóm:
- "CÒN HÀNG": chỉ được giới thiệu / chốt các mã trong nhóm này.
- "ĐÃ BÁN": các mã này KHÔNG còn. TUYỆT ĐỐI không nói chúng còn hàng, không báo giá/số đo như hàng còn. Nếu khách hỏi tới thì báo đã bán.

1. "suggested_reply" — trả lời (các) tin MỚI NHẤT của Khách mà Shop chưa trả lời:
   - Nếu có dòng "KHÁCH ĐANG HỎI", chỉ báo còn/hết cho đúng các mã đó; không nhắc lại mã đã trả lời ở lượt trước.
   - Giọng như chủ shop nhắn Messenger: tự nhiên, thân thiện, NGẮN 1-3 câu. Dùng "dạ", "nha", "nhé", "ạ"; được viết tắt quen thuộc (sp, k, dc, r, ib, ship); tối đa 1 emoji hợp ngữ cảnh (😊 🥰 👍 ✨ 🥲).
   - Không mở đầu "Chào anh/chị ..." ở mọi tin; không dùng câu văn mẫu kiểu "Mình có muốn tham khảo thêm không ạ?", "Nếu cần thêm thông tin gì cứ cho em biết", "Em rất sẵn lòng hỗ trợ".
   - Giá viết gọn: 350k. KHÔNG kèm số đo, trừ khi khách hỏi size / số đo / chiều cao cân nặng / có vừa không.
   - Không nhắc lại điều Shop đã trả lời trước đó.
   - Chỉ dùng dữ liệu sản phẩm và "Chính sách shop". KHÔNG bịa phí ship, thời gian giao, đổi trả, khuyến mãi. Thiếu dữ liệu: "needs_human": true và trả lời kiểu "Dạ để em check lại rồi báo mình liền nha".
   - Khách hỏi vừa không: dựa số đo thực tế (đồ vintage: ngực áo thường lớn hơn ngực người 10-16cm; eo quần ghi nửa vòng thì x2).
   - Khách mặc cả / xin giảm / đổi trả / khiếu nại: "needs_human": true + "human_reason", trả lời khéo kiểu đợi em check.
   - Khách đã chốt và gửi đủ tên/SĐT/địa chỉ: xác nhận ngắn mã + giá.
   - Khi liệt kê nhiều mã, không ghi số lượng kiểu "3 chiếc", chỉ liệt kê mã.
   Ví dụ giọng:
   • Khách: "Mẫu này còn k shop" (còn hàng) → "Dạ mẫu này còn nha anh, 350k ạ 😊"
   • Khách gửi 4 ảnh, chỉ TS4764 còn → "Dạ mấy mẫu kia em bán mất r ạ 🥲 còn mỗi TS4764 350k thôi nha anh"
   • Hết hàng → "Dạ mẫu này em bán mất r anh ơi 🥲"
2. "size_advice": nhận xét size ngắn gọn (chỉ khi khách hỏi size/số đo), để trống nếu không.
3. "order_draft": tên, SĐT, địa chỉ, ghi chú khách đã gửi, và mã sản phẩm CÒN HÀNG khách muốn lấy (không bao giờ đưa mã ĐÃ BÁN).

Trả về JSON thuần:
{"suggested_reply":"string","size_advice":"string","needs_human":false,"human_reason":"string","order_draft":{"name":"string","phone":"string","address":"string","note":"string","product_codes":["TS123"]}}`

	policy := strings.TrimSpace(os.Getenv("COPILOT_SHOP_POLICY"))
	if policy == "" {
		policy = "(Chưa cấu hình — không được tự đưa ra phí ship / thời gian giao / đổi trả)"
	}
	userPrompt := fmt.Sprintf(
		"Tên khách Facebook: %s\nAction hiện tại: %s\n\nChính sách shop:\n%s\n\n%s\n\nHội thoại hiện tại:\n%s",
		customerName, action, policy, productContext, transcript,
	)

	out, ok := s.callCopilotAI(ctx, systemPrompt, userPrompt)
	if !ok {
		return s.fallbackSuggestion(customerName, action, transcript, messages, products)
	}
	// Guard: the model must never present a SOLD code as available. Retry once with an
	// explicit correction; if it still slips, flag the reply for a human.
	if bad := soldCodesPresentedAsAvailable(out.SuggestedReply, soldCodes); len(bad) > 0 {
		retryPrompt := userPrompt + "\n\nLƯU Ý: bản nháp trước đã giới thiệu nhầm mã ĐÃ BÁN " + strings.Join(bad, ", ") +
			" như còn hàng. Viết lại, các mã đó phải báo là đã bán hoặc không nhắc tới."
		if out2, ok2 := s.callCopilotAI(ctx, systemPrompt, retryPrompt); ok2 {
			out = out2
		}
		if bad2 := soldCodesPresentedAsAvailable(out.SuggestedReply, soldCodes); len(bad2) > 0 {
			out.NeedsHuman = true
			out.HumanReason = "Gợi ý có nhắc mã đã bán (" + strings.Join(bad2, ", ") + ") — kiểm tra lại trước khi gửi"
		}
	}
	if out.OrderDraft != nil {
		out.OrderDraft.ProductCodes = filterOutCodes(out.OrderDraft.ProductCodes, soldCodes)
		if out.OrderDraft.Name == "" {
			out.OrderDraft.Name = customerName
		}
	}
	return out
}

// callCopilotAI sends one chat completion and decodes the JSON reply.
func (s *CopilotService) callCopilotAI(ctx context.Context, systemPrompt, userPrompt string) (aiCopilotOutput, bool) {
	var out aiCopilotOutput
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
		return out, false
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return out, false
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.OpenAIKey)

	resp, err := s.Client.Do(req)
	if err != nil {
		return out, false
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return out, false
	}

	var openAIResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage OpenAITokenUsage `json:"usage"`
	}
	if err := json.Unmarshal(respBytes, &openAIResp); err != nil || len(openAIResp.Choices) == 0 {
		return out, false
	}
	RecordOpenAIUsage("copilot", "gpt-4o-mini", openAIResp.Usage)

	if err := json.Unmarshal([]byte(openAIResp.Choices[0].Message.Content), &out); err != nil {
		return out, false
	}
	return out, true
}

func (s *CopilotService) fallbackSuggestion(customerName, action, transcript string, messages []CopilotMessage, products []CopilotProductMatch) aiCopilotOutput {
	// Parse order info only from what the customer typed: shop messages and photo OCR (bank QR,
	// product labels, receipts) produced garbage names/addresses.
	parsed := s.SmartOrder.FallbackParse(CustomerTextForOrderParse(messages))
	// The rule parser often takes a chat line ("cho mình lấy quần này nhé") as the name;
	// the Facebook name is the reliable default (editable on the Create Order page).
	parsed.Name = customerName

	// Only answer codes the customer asked about after the shop's last reply.
	pending := s.pendingCustomerCodes(messages)
	isPending := map[string]bool{}
	for _, c := range pending {
		isPending[c] = true
	}
	var availCodes, availParts, soldParts []string
	for _, p := range products {
		avail := strings.EqualFold(p.Status, "AVAILABLE")
		if avail {
			availCodes = append(availCodes, p.Code)
		}
		if len(pending) > 0 && !isPending[strings.ToUpper(p.Code)] {
			continue
		}
		if avail {
			availParts = append(availParts, p.Code+" "+priceK(p.Price))
		} else {
			soldParts = append(soldParts, p.Code)
		}
	}

	lower := strings.ToLower(transcript)
	needsHuman := strings.Contains(lower, "giảm giá") || strings.Contains(lower, "bớt") || strings.Contains(lower, "fix") || strings.Contains(lower, "đổi trả")
	humanReason := ""
	if needsHuman {
		humanReason = "Khách hỏi giảm giá / đổi trả"
	}

	var reply string
	switch {
	case len(availParts) > 0 && len(soldParts) > 0:
		reply = "Dạ " + strings.Join(availParts, ", ") + " còn nha 😊 còn " + strings.Join(soldParts, ", ") + " em bán mất r ạ"
	case len(availParts) > 0:
		reply = "Dạ " + strings.Join(availParts, ", ") + " còn nha anh 😊"
	case len(soldParts) > 0:
		reply = "Dạ " + strings.Join(soldParts, ", ") + " em bán mất r ạ 🥲"
	default:
		reply = "Dạ mình gửi em mã sp hoặc chiều cao cân nặng để em check cho nha 😊"
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

// pendingCustomerCodes returns product codes the customer mentioned (text or photo OCR)
// after the shop's last message, i.e. what still needs an answer.
func (s *CopilotService) pendingCustomerCodes(messages []CopilotMessage) []string {
	start := 0
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Sender == "shop" {
			start = i + 1
			break
		}
	}
	seen := map[string]bool{}
	var codes []string
	for _, m := range messages[start:] {
		if m.Sender != "customer" {
			continue
		}
		for _, c := range append(s.ExtractProductCodesFromText(m.Text), s.ExtractProductCodesFromText(m.OCRText)...) {
			if !seen[c] {
				seen[c] = true
				codes = append(codes, c)
			}
		}
	}
	return codes
}

// buildProductContext renders products for the prompt as two explicit groups so the model
// can't mix up sold and available items, and returns the sold codes for the output guard.
func buildProductContext(products []CopilotProductMatch, pending []string) (string, []string) {
	var avail, sold []string
	var soldCodes []string
	for _, p := range products {
		if strings.EqualFold(p.Status, "AVAILABLE") {
			var sizes []string
			for _, sz := range p.Sizes {
				sizes = append(sizes, fmt.Sprintf("%s %g", sz.Name, sz.Value))
			}
			line := fmt.Sprintf("- %s | %s | %s", p.Code, p.Type, priceK(p.Price))
			if len(sizes) > 0 {
				line += " | số đo: " + strings.Join(sizes, ", ")
			}
			avail = append(avail, line)
		} else {
			sold = append(sold, p.Code)
			soldCodes = append(soldCodes, strings.ToUpper(p.Code))
		}
	}
	var b strings.Builder
	if len(pending) > 0 {
		b.WriteString("KHÁCH ĐANG HỎI (chưa trả lời): " + strings.Join(pending, ", ") + "\n\n")
	}
	b.WriteString("CÒN HÀNG (chỉ được giới thiệu các mã này):\n")
	if len(avail) == 0 {
		b.WriteString("(không có)\n")
	} else {
		b.WriteString(strings.Join(avail, "\n") + "\n")
	}
	b.WriteString("\nĐÃ BÁN (không còn, không được nói là còn hàng): ")
	if len(sold) == 0 {
		b.WriteString("(không có)")
	} else {
		b.WriteString(strings.Join(sold, ", "))
	}
	return b.String(), soldCodes
}

var soldWordsRe = regexp.MustCompile(`(?i)(bán|hết|sold|không còn|k còn|ko còn|chưa còn)`)

// soldCodesPresentedAsAvailable returns sold codes that appear in a sentence of the reply
// without any "sold / out of stock" wording next to them.
func soldCodesPresentedAsAvailable(reply string, soldCodes []string) []string {
	if reply == "" || len(soldCodes) == 0 {
		return nil
	}
	sentences := regexp.MustCompile(`[.!?\n]+`).Split(reply, -1)
	var bad []string
	for _, code := range soldCodes {
		num := strings.TrimLeft(code, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
		codeRe := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(code) + `\b|\b` + regexp.QuoteMeta(num) + `\b`)
		for _, sen := range sentences {
			if codeRe.MatchString(sen) && !soldWordsRe.MatchString(sen) {
				bad = append(bad, code)
				break
			}
		}
	}
	return bad
}

func filterOutCodes(codes, remove []string) []string {
	rm := map[string]bool{}
	for _, c := range remove {
		rm[strings.ToUpper(c)] = true
	}
	var out []string
	for _, c := range codes {
		if !rm[strings.ToUpper(strings.TrimSpace(c))] {
			out = append(out, c)
		}
	}
	return out
}

// priceK formats 350000 as "350k" (and 1250000 as "1tr250").
func priceK(v int64) string {
	if v <= 0 {
		return ""
	}
	if v >= 1000000 {
		tr := v / 1000000
		rest := (v % 1000000) / 1000
		if rest == 0 {
			return fmt.Sprintf("%dtr", tr)
		}
		return fmt.Sprintf("%dtr%03d", tr, rest)
	}
	if v%1000 == 0 {
		return fmt.Sprintf("%dk", v/1000)
	}
	return FormatVNDPrice(v)
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
	if msgType == "image" && imageURL != "" && s.Vision != nil && sender != "shop" {
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

// RebuildAndAnalyzeSession loads the full stored thread (newest copilotDisplayMessages), runs AI and
// product matching on the recent tail, and returns a CopilotSession. Orders no longer cut the history.
func (s *CopilotService) RebuildAndAnalyzeSession(ctx context.Context, psid string, action string) (*CopilotSession, error) {
	return s.RebuildSession(ctx, psid, true)
}

// RebuildSession reloads a thread from the DB and re-checks product status. OpenAI is called only
// when useAI is true (the "Gợi ý AI" button); otherwise the previous AI suggestion is kept if no
// new message arrived since, else a free rule-based suggestion is used.
func (s *CopilotService) RebuildSession(ctx context.Context, psid string, useAI bool) (*CopilotSession, error) {
	action := "suggest"
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
	// Newest messages (whole history up to the display cap), then flip back to chronological order.
	var recs []models.CopilotMessageRecord
	if err := query.Order("msg_timestamp DESC, id DESC").Limit(copilotDisplayMessages).Find(&recs).Error; err != nil {
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
		var tsStr string
		if !r.MsgTimestamp.IsZero() {
			tsStr = r.MsgTimestamp.Format(time.RFC3339)
		}
		msgs = append(msgs, CopilotMessage{
			Sender:    r.Sender,
			Type:      r.Type,
			Text:      r.Text,
			ImageURL:  r.ImageURL,
			OCRText:   r.OCRText,
			Timestamp: tsStr,
		})
	}
	// Products + AI work on the recent part of the chat only, so old purchases don't flood the
	// suggestion; the full history is still returned for display.
	aiMsgs := tailMessages(msgs, copilotAIContextMessages)
	lastMention := map[string]int{}
	for idx, m := range aiMsgs {
		var found []string
		if m.Type == "text" && m.Text != "" {
			found = s.ExtractProductCodesFromText(m.Text)
		}
		if m.Type == "image" && m.OCRText != "" {
			found = append(found, s.ExtractProductCodesFromText(m.OCRText)...)
		}
		addCodes(found)
		for _, c := range found {
			lastMention[c] = idx
		}
	}

	detectedProducts := s.lookupProductsByCodes(allCodes)
	// Most recently mentioned first (list chips + product sheet show the newest codes).
	sortProductsByMention(detectedProducts, lastMention)
	convText := s.buildConversationTranscript(aiMsgs)
	recommendedProducts := s.recommendProductsByMeasurements(convText, detectedProducts)
	allMatched := append(detectedProducts, recommendedProducts...)

	lastKey := ""
	if n := len(msgs); n > 0 {
		lastKey = messageKey(msgs[n-1])
	}
	var aiRes aiCopilotOutput
	aiGenerated := false
	if useAI {
		aiRes = s.generateAISuggestion(ctx, conv.CustomerName, action, aiMsgs, allMatched)
		aiGenerated = strings.TrimSpace(s.OpenAIKey) != ""
	} else if prev, ok := s.GetSession(conv.PSID); ok && prev.AIGenerated && prev.AIForKey == lastKey {
		aiRes = aiCopilotOutput{
			SuggestedReply: prev.SuggestedReply,
			SizeAdvice:     prev.SizeAdvice,
			NeedsHuman:     prev.NeedsHuman,
			HumanReason:    prev.HumanReason,
			OrderDraft:     prev.OrderDraft,
		}
		aiGenerated = true
	} else {
		aiRes = s.fallbackSuggestion(conv.CustomerName, action, convText, aiMsgs, allMatched)
	}

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
		AIGenerated:      aiGenerated,
		AIForKey:         lastKey,
	}
	if n := len(recs); n > 0 {
		sess.LastMessageAt = recs[n-1].MsgTimestamp
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
		if sess, err := s.RebuildSession(ctx, conv.PSID, false); err == nil {
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
	convs, err := s.fetchGraphConversations(ctx, pageID, pageToken, 25, 10, "")
	if err != nil {
		return 0, err
	}
	synced := 0
	for _, conv := range convs {
		psid := s.importGraphConversation(conv, pageID, pageToken)
		if psid == "" {
			continue
		}
		if sess, err := s.RebuildSession(ctx, psid, false); err == nil {
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

// ImportFBHistory pulls the last 100 messages of one customer's thread so a customer first seen
// via webhook starts with full context (not just the message that triggered the webhook).
func (s *CopilotService) ImportFBHistory(ctx context.Context, pageID, psid, pageToken string) error {
	if strings.TrimSpace(psid) == "" {
		return fmt.Errorf("empty psid")
	}
	convs, err := s.fetchGraphConversations(ctx, pageID, pageToken, 100, 1, "&user_id="+url.QueryEscape(psid))
	if err != nil {
		return err
	}
	for _, conv := range convs {
		s.importGraphConversation(conv, pageID, pageToken)
	}
	return nil
}

var fbSystemNoticeRe = regexp.MustCompile(`(?i)^.{1,80}?\s(đã trả lời (về )?(một |1 )?(quảng cáo|bài viết)|replied to (an|a|your) (ad|post))\.?(\s|\(|$)`)

// Business Suite automation lines, e.g. "Đã thêm nhãn tự động: Đã đánh dấu trạng thái đơn đặt
// hàng là Đã đặt hàng." — they arrive as page messages but nobody typed them.
var fbAutomationNoticeRe = regexp.MustCompile(`(?i)^(đã thêm nhãn tự động|đã đánh dấu trạng thái đơn|đã xóa nhãn tự động|đã gỡ nhãn tự động|automatically (added|removed) (the )?label|marked (the )?order status)`)

// IsFBSystemNotice reports Meta-generated banner lines that arrive as page messages.
func IsFBSystemNotice(text string) bool {
	t := strings.TrimSpace(text)
	return t != "" && (fbSystemNoticeRe.MatchString(t) || fbAutomationNoticeRe.MatchString(t))
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

// SortTime is the newest-message time, falling back to UpdatedAt for older stored sessions.
func (sess *CopilotSession) SortTime() time.Time {
	if !sess.LastMessageAt.IsZero() {
		return sess.LastMessageAt
	}
	return sess.UpdatedAt
}

// SendFBText sends a plain text message to a Messenger user (PSID) via the Page Send API.
func (s *CopilotService) SendFBText(ctx context.Context, psid, text, pageToken string) error {
	if psid == "" || pageToken == "" {
		return fmt.Errorf("thiếu PSID hoặc FB_MAIN_PAGE_TOKEN")
	}
	payload, _ := json.Marshal(map[string]interface{}{
		"recipient":      map[string]string{"id": psid},
		"messaging_type": "RESPONSE",
		"message":        map[string]string{"text": text},
	})
	endpoint := "https://graph.facebook.com/v21.0/me/messages?access_token=" + url.QueryEscape(pageToken)
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
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

// originalImages lists a product's original photos (not the scaled thumb): the main photo
// (PathThumb / IsMain image) first, then the detail photos, de-duplicated.
func (s *CopilotService) originalImages(p models.Product) []string {
	var out []string
	seen := map[string]bool{}
	add := func(path string) {
		u := s.formatMediaPath(strings.TrimSpace(path))
		if u == "" || seen[u] {
			return
		}
		seen[u] = true
		out = append(out, u)
	}
	if p.PathThumb != nil {
		add(*p.PathThumb)
	}
	for _, img := range p.Images {
		if img.IsMain {
			add(img.FilePath)
		}
	}
	for _, img := range p.Images {
		if !img.IsMain {
			add(img.FilePath)
		}
	}
	return out
}

// sortProductsByMention orders products by the index of their latest mention (newest first).
func sortProductsByMention(products []CopilotProductMatch, lastMention map[string]int) {
	sort.SliceStable(products, func(i, j int) bool {
		return lastMention[strings.ToUpper(products[i].Code)] > lastMention[strings.ToUpper(products[j].Code)]
	})
}

// productMentionIndex maps each product code to the index of its latest mention in msgs.
func (s *CopilotService) productMentionIndex(msgs []CopilotMessage) map[string]int {
	last := map[string]int{}
	for idx, m := range msgs {
		for _, c := range append(s.ExtractProductCodesFromText(m.Text), s.ExtractProductCodesFromText(m.OCRText)...) {
			last[c] = idx
		}
	}
	return last
}

// ProductImagesForSend returns the original photo URLs of a product: kind "main" = the main
// photo; kind "detail" = main + all detail photos (only for AVAILABLE products).
func (s *CopilotService) ProductImagesForSend(productID uint, kind string) ([]string, string, error) {
	if s.DB == nil {
		return nil, "", fmt.Errorf("database not initialized")
	}
	var p models.Product
	if err := s.DB.Preload("Images").First(&p, productID).Error; err != nil {
		return nil, "", fmt.Errorf("không tìm thấy sản phẩm")
	}
	code := strPtrVal(p.Name)
	all := s.originalImages(p)
	if len(all) == 0 {
		return nil, code, fmt.Errorf("sản phẩm %s chưa có ảnh", code)
	}
	if kind == "main" {
		return all[:1], code, nil
	}
	if !strings.EqualFold(strPtrVal(p.Status), "AVAILABLE") {
		return nil, code, fmt.Errorf("%s đã bán — không gửi ảnh chi tiết", code)
	}
	if len(all) < 2 {
		return nil, code, fmt.Errorf("%s chưa có ảnh chi tiết", code)
	}
	return all, code, nil // main photo first, then the detail photos
}

// SendFBImageURL sends one image to a Messenger user by public URL (Facebook downloads it).
func (s *CopilotService) SendFBImageURL(ctx context.Context, psid, imageURL, pageToken string) error {
	if psid == "" || pageToken == "" {
		return fmt.Errorf("thiếu PSID hoặc FB_MAIN_PAGE_TOKEN")
	}
	payload, _ := json.Marshal(map[string]interface{}{
		"recipient":      map[string]string{"id": psid},
		"messaging_type": "RESPONSE",
		"message": map[string]interface{}{
			"attachment": map[string]interface{}{
				"type":    "image",
				"payload": map[string]interface{}{"url": imageURL, "is_reusable": true},
			},
		},
	})
	endpoint := "https://graph.facebook.com/v21.0/me/messages?access_token=" + url.QueryEscape(pageToken)
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
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

// CustomerTextForOrderParse joins the customer's own text messages (newest 20, no OCR, no shop
// messages) — the input used to extract name / phone / address for an order.
func CustomerTextForOrderParse(messages []CopilotMessage) string {
	var lines []string
	for _, m := range messages {
		if m.Sender == "customer" && m.Type != "image" && strings.TrimSpace(m.Text) != "" {
			lines = append(lines, strings.TrimSpace(m.Text))
		}
	}
	if len(lines) > 20 {
		lines = lines[len(lines)-20:]
	}
	return strings.Join(lines, "\n")
}
