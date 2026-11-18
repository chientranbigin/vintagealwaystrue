package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

// SmartOrderParsedResult represents the extracted order details
type SmartOrderParsedResult struct {
	Name     string   `json:"name"`
	Phone    string   `json:"phone"`
	Address  string   `json:"address"`
	Products []string `json:"products"`
	Note     string   `json:"note"`
}

type SmartOrderService struct {
	OpenAIKey  string
	Model      string
	HTTPClient *http.Client
}

func NewSmartOrderService() *SmartOrderService {
	key := os.Getenv("OPENAI_API_KEY")
	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = "gpt-4o-mini"
	}

	return &SmartOrderService{
		OpenAIKey: key,
		Model:     model,
		HTTPClient: &http.Client{
			Timeout: 25 * time.Second,
		},
	}
}

// ParseDescription parses order message text into structured order data
func (s *SmartOrderService) ParseDescription(ctx context.Context, text string) (*SmartOrderParsedResult, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return &SmartOrderParsedResult{}, nil
	}

	// 1. If OpenAI API Key is configured, attempt OpenAI API call first
	if s.OpenAIKey != "" {
		res, err := s.callOpenAI(ctx, text)
		if err == nil && res != nil {
			return res, nil
		}
		// Log error and fall through to heuristic parser
		fmt.Printf("[SmartOrder] OpenAI parsing failed or timed out: %v. Falling back to local heuristic parser.\n", err)
	}

	// 2. Fallback heuristic regex parser
	return s.FallbackParse(text), nil
}

type openAIChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIResponseFormat struct {
	Type string `json:"type"`
}

type openAIChatRequest struct {
	Model          string                `json:"model"`
	Messages       []openAIChatMessage   `json:"messages"`
	Temperature    float64               `json:"temperature"`
	ResponseFormat *openAIResponseFormat `json:"response_format,omitempty"`
}

type openAIChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
	Usage OpenAITokenUsage `json:"usage"`
}

func (s *SmartOrderService) callOpenAI(ctx context.Context, text string) (*SmartOrderParsedResult, error) {
	prompt := fmt.Sprintf(`Phân tích tin nhắn đặt hàng/giao hàng của khách (có thể viết tự do trên 1 dòng hoặc nhiều dòng, đảo thứ tự, viết tắt địa danh như KTT, HBT, HN, HCM, Q1, P., KĐT..., dùng dấu -, |, /, dấu chấm trong SĐT) và xuất ra JSON với đúng 5 trường sau:

- "name": tên người nhận (chỉ lấy tên, không dính SĐT hay địa chỉ; nếu không có thì để "")
- "phone": số điện thoại Việt Nam chỉ gồm chữ số bắt đầu bằng 0 (ví dụ: 0945435858; bỏ dấu chấm/cách, đổi +84/84 đầu thành 0; nếu không có thì để "")
- "address": địa chỉ giao hàng đầy đủ (giữ nguyên thông tin địa chỉ kể cả phần trong ngoặc hoặc chữ viết tắt, bỏ chữ "Địa chỉ:" ở đầu nếu có, KHÔNG kèm tên người nhận hay SĐT)
- "products": mảng chuỗi các mã sản phẩm (ví dụ: ["ts1169", "ts1158"]). Lưu ý: KHÔNG nhầm số nhà, lô, tòa, phòng, quận (như "Nhà E", "E12", "P302", "Q1", "Lô A2") với mã sản phẩm. Nếu không có thì để [].
- "note": nội dung ghi chú giao hàng (như nội dung sau "note:", "ghi chú:", "lưu ý:", hoặc yêu cầu giao giờ hành chính/gọi trước...). Nếu nhiều ý thì ngăn cách bằng dấu phẩy. Nếu không có thì để "".

Ví dụ:
Input: "Địa chỉ: Nhà E, khu tập thể Nguyễn Công Trứ, Hai Bà Trưng, Hà Nội - HUY - 0945435858"
Output: {"name":"HUY","phone":"0945435858","address":"Nhà E, khu tập thể Nguyễn Công Trứ, Hai Bà Trưng, Hà Nội","products":[],"note":""}

Input: "ts1169, st220 - Lan Anh - 0912.345.678 | P302 KTT Thành Công, Ba Đình, HN (note: giao sau 5h)"
Output: {"name":"Lan Anh","phone":"0912345678","address":"P302 KTT Thành Công, Ba Đình, HN","products":["ts1169","st220"],"note":"giao sau 5h"}

Yêu cầu: Trả về duy nhất 1 JSON object hợp lệ.

Đoạn text cần phân tích:
"""%s"""`, text)

	reqBody := openAIChatRequest{
		Model: s.Model,
		Messages: []openAIChatMessage{
			{Role: "system", Content: "Bạn là trợ lý trích xuất dữ liệu đơn hàng tiếng Việt cực kỳ chính xác. Luôn trả về JSON object hợp lệ."},
			{Role: "user", Content: prompt},
		},
		Temperature:    0.1,
		ResponseFormat: &openAIResponseFormat{Type: "json_object"},
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+s.OpenAIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("openai API status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var chatResp openAIChatResponse
	if err := json.Unmarshal(bodyBytes, &chatResp); err != nil {
		return nil, err
	}

	if len(chatResp.Choices) == 0 {
		return nil, fmt.Errorf("empty choices from openai")
	}
	RecordOpenAIUsage("smart_order", s.Model, chatResp.Usage)

	content := strings.TrimSpace(chatResp.Choices[0].Message.Content)
	// Strip markdown code fences if present (e.g. ```json ... ```)
	if strings.HasPrefix(content, "```") {
		lines := strings.Split(content, "\n")
		var cleanLines []string
		for _, l := range lines {
			trimmed := strings.TrimSpace(l)
			if strings.HasPrefix(trimmed, "```") {
				continue
			}
			cleanLines = append(cleanLines, l)
		}
		content = strings.TrimSpace(strings.Join(cleanLines, "\n"))
	}

	var result SmartOrderParsedResult
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal openai json: %w (raw: %s)", err, content)
	}

	// Clean and normalize phone
	cleanPhoneRe := regexp.MustCompile(`\D`)
	result.Phone = cleanPhoneRe.ReplaceAllString(result.Phone, "")
	if strings.HasPrefix(result.Phone, "84") && len(result.Phone) == 11 {
		result.Phone = "0" + result.Phone[2:]
	}
	result.Name = strings.TrimSpace(result.Name)
	result.Address = strings.TrimSpace(result.Address)
	result.Note = strings.TrimSpace(result.Note)
	if result.Products == nil {
		result.Products = []string{}
	}

	return &result, nil
}

// FallbackParse extracts order information using heuristics & regex when AI is unavailable
func (s *SmartOrderService) FallbackParse(text string) *SmartOrderParsedResult {
	result := &SmartOrderParsedResult{
		Products: []string{},
	}

	// 1. Extract Note
	// Matches 'note:', 'ghi chú:', 'lưu ý:' etc.
	noteRe := regexp.MustCompile(`(?i)(?:note|ghi\s*chú|ghi\s*chu|lưu\s*ý|luu\s*y)\s*[:：]\s*([^\n\r]+)`)
	if match := noteRe.FindStringSubmatch(text); len(match) > 1 {
		result.Note = strings.TrimSpace(match[1])
	}

	// Remove note lines for subsequent parsing
	textWithoutNote := noteRe.ReplaceAllString(text, "")

	// 2. Extract Phone
	phoneRe := regexp.MustCompile(`(?:\+84|0)[35789]\d{8}|\b0\d{9}\b`)
	rawPhone := phoneRe.FindString(textWithoutNote)
	if rawPhone != "" {
		digitsRe := regexp.MustCompile(`\D`)
		cleanPhone := digitsRe.ReplaceAllString(rawPhone, "")
		if strings.HasPrefix(cleanPhone, "84") && len(cleanPhone) == 11 {
			cleanPhone = "0" + cleanPhone[2:]
		}
		result.Phone = cleanPhone
	}

	// 3. Extract Product Codes (e.g. ts1169, TS1234, ST567, hj12, BE999, etc.)
	codeRe := regexp.MustCompile(`(?i)\b([a-zA-Z]{1,6}\d{1,7})\b`)
	matches := codeRe.FindAllString(textWithoutNote, -1)
	seenCodes := make(map[string]bool)
	for _, m := range matches {
		upper := strings.ToUpper(m)
		// Filter out obvious false positives like years (e.g. 2024), street numbers, or words
		if !seenCodes[upper] {
			seenCodes[upper] = true
			result.Products = append(result.Products, strings.ToLower(m))
		}
	}

	// 4. Lines breakdown for Name and Address
	lines := strings.Split(textWithoutNote, "\n")
	var remainingLines []string
	addressKeywords := []string{
		"địa chỉ", "dia chi", "đường", "duong", "quận", "quan", "huyện", "huyen",
		"phường", "phuong", "xã", "xa", "thành phố", "thanh pho", "tp.", "tp ",
		"hà nội", "ha noi", "hcm", "hồ chí minh", "đà nẵng", "da nang", "tỉnh", "tinh",
		"số nhà", "so nha", "ngõ", "ngo", "ngách", "ngach", "thôn", "thon", "xóm", "xom",
		"ấp", "ap", "tổ", "to ", "khu phố", "chung cư", "tòa", "phòng",
	}

	namePrefixes := []string{"anh ", "chị ", "chi ", "bác ", "bac ", "cô ", "co ", "chú ", "chu ", "em "}

	var detectedAddressParts []string
	var candidateNames []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		// Skip lines that only contain the phone number
		if rawPhone != "" && strings.Contains(trimmed, rawPhone) && len(trimmed) <= len(rawPhone)+4 {
			continue
		}

		lowerLine := strings.ToLower(trimmed)

		// Check if line looks like an address
		isAddress := false
		for _, kw := range addressKeywords {
			if strings.Contains(lowerLine, kw) {
				isAddress = true
				break
			}
		}

		if isAddress {
			// Strip leading "Địa chỉ:" if present
			addrPrefixRe := regexp.MustCompile(`(?i)^(?:địa\s*chỉ|dia\s*chi)\s*[:：]\s*`)
			cleanAddr := addrPrefixRe.ReplaceAllString(trimmed, "")
			detectedAddressParts = append(detectedAddressParts, cleanAddr)
			continue
		}

		// Check if line contains name prefix or could be name
		hasNamePrefix := false
		for _, pfx := range namePrefixes {
			if strings.HasPrefix(lowerLine, pfx) {
				hasNamePrefix = true
				break
			}
		}

		if hasNamePrefix {
			candidateNames = append(candidateNames, trimmed)
			continue
		}

		// If line doesn't match address and isn't just product codes
		cleanedWithoutCodes := codeRe.ReplaceAllString(trimmed, "")
		if rawPhone != "" {
			cleanedWithoutCodes = strings.ReplaceAll(cleanedWithoutCodes, rawPhone, "")
		}
		cleanedWithoutCodes = strings.TrimSpace(cleanedWithoutCodes)

		if cleanedWithoutCodes != "" {
			remainingLines = append(remainingLines, cleanedWithoutCodes)
		}
	}

	// Assign Name
	if len(candidateNames) > 0 {
		result.Name = candidateNames[0]
		// If name has phone in it, strip it
		if rawPhone != "" {
			result.Name = strings.TrimSpace(strings.ReplaceAll(result.Name, rawPhone, ""))
		}
	} else if len(remainingLines) > 0 {
		// First remaining non-address line is typically the customer name
		result.Name = remainingLines[0]
		if rawPhone != "" {
			result.Name = strings.TrimSpace(strings.ReplaceAll(result.Name, rawPhone, ""))
		}
		remainingLines = remainingLines[1:]
	}

	// Assign Address
	if len(detectedAddressParts) > 0 {
		result.Address = strings.Join(detectedAddressParts, ", ")
	} else if len(remainingLines) > 0 {
		result.Address = strings.Join(remainingLines, ", ")
	}

	return result
}
