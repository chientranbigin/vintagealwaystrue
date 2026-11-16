package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type ContentGenerationRequest struct {
	CategoryCode string `json:"category_code"` // TROUSERS, GILE, BLAZER, JACKET, SHIRT, etc.
	CategoryName string `json:"category_name"`
	AlbumURL     string `json:"album_url"`
	IncludeEvent bool   `json:"include_event"`
	EventName    string `json:"event_name"` // "Giáng Sinh / Noel", "Tết Nguyên Đán", "Mùa Thu Đông", "Tiệc cuối năm", etc.
	CustomPrompt string `json:"custom_prompt,omitempty"`
}

type AIContentService struct {
	OpenAIKey string
	GeminiKey string
	Client    *http.Client
}

func NewAIContentService() *AIContentService {
	return &AIContentService{
		OpenAIKey: os.Getenv("OPENAI_API_KEY"),
		GeminiKey: os.Getenv("GEMINI_API_KEY"),
		Client:    &http.Client{Timeout: 30 * time.Second},
	}
}

// GeneratePostContent produces an engaging vintage Facebook post caption
func (s *AIContentService) GeneratePostContent(ctx context.Context, req ContentGenerationRequest) (string, error) {
	// If OpenAI Key is available, attempt OpenAI generation first
	if s.OpenAIKey != "" {
		if content, err := s.callOpenAI(ctx, req); err == nil && content != "" {
			return content, nil
		}
	}

	// Fallback to high-quality procedural vintage templates
	return s.GenerateProceduralContent(req), nil
}

func (s *AIContentService) callOpenAI(ctx context.Context, req ContentGenerationRequest) (string, error) {
	prompt := fmt.Sprintf(`Bạn là copywriter chuyên nghiệp cho shop thời trang nam Vintage Always True.
Hãy viết một bài đăng Facebook bán hàng cực kỳ hấp dẫn cho danh mục: %s (%s).

Yêu cầu nội dung:
1. Tiêu đề (Hook): Bắt mắt với icon, nêu rõ loại sản phẩm và mức giá đặc trưng (nếu quần âu thì đồng giá 180k, gile 120k, blazer/jacket giá trên từng ảnh).
2. Sự kiện / Mùa lễ: %s
3. Điểm nổi bật: Tôn lên nét đẹp cổ điển (classic vintage), chất liệu vải, độ mới cao, tính đa dụng dễ phối đồ.
4. Thông số: Nhắc nhở người mua đọc kỹ thông số (ngang ngực, vai, eo, dài vạt...) đã ghi trực tiếp trên từng ảnh.
5. Cam kết của shop: Hàng 2hand tuyển chọn, mỗi mẫu 1 chiếc duy nhất, đã giặt ủi sạch sẽ, ship COD toàn quốc hỗ trợ đồng kiểm.
6. Kết bài bắt buộc:
👉 Để xem thêm sp và check sp còn hàng tại:
%s
7. Hashtags phù hợp ở cuối bài.
Không dùng markdown dấu sao quá nhiều, dùng định dạng văn bản Facebook tự nhiên, gạch ngang phân đoạn, icon phù hợp.`,
		req.CategoryCode,
		req.CategoryName,
		func() string {
			if req.IncludeEvent && req.EventName != "" {
				return "BẮT BUỘC lồng ghép khéo léo không khí: " + req.EventName + " (ví dụ: sắm đồ ấm đi chơi Noel, dự tiệc cuối năm, du xuân, lên đồ ấm áp ngày đông...)"
			}
			return "KHÔNG lồng ghép dịp lễ, viết theo phong cách chuẩn mực Vintage hằng ngày (timeless classic)."
		}(),
		req.AlbumURL,
	)

	payload := map[string]interface{}{
		"model": "gpt-4o-mini",
		"messages": []map[string]string{
			{"role": "system", "content": "Bạn là chuyên gia viết content bán hàng thời trang Vintage chất lượng cao, súc tích, tinh tế và có gu."},
			{"role": "user", "content": prompt},
		},
		"temperature": 0.8,
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Authorization", "Bearer "+s.OpenAIKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("openai error (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	var chatResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(bodyBytes, &chatResp); err != nil || len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("invalid openai response")
	}

	return strings.TrimSpace(chatResp.Choices[0].Message.Content), nil
}

// GenerateProceduralContent generates high quality authentic posts based on real page patterns
func (s *AIContentService) GenerateProceduralContent(req ContentGenerationRequest) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := strings.ToUpper(strings.TrimSpace(req.CategoryCode))

	albumLink := req.AlbumURL
	if albumLink == "" {
		albumLink = "https://www.facebook.com/236648362875702"
	}

	// Determine festive/seasonal angle
	eventAngle := ""
	if req.IncludeEvent {
		evt := strings.ToLower(req.EventName)
		if strings.Contains(evt, "giáng sinh") || strings.Contains(evt, "noel") || strings.Contains(evt, "christmas") {
			noelHooks := []string{
				"🎄 MÙA GIÁNG SINH ĐẾN GẦN – THỜI ĐIỂM HOÀN HẢO ĐỂ LÊN ĐỒ ẤM ÁP ĐÓN NOEL! ❄️",
				"🎅 LÊN ĐỒ XUỐNG PHỐ MÙA GIÁNG SINH CÙNG GU VINTAGE CHỈN CHU 🎁",
				"❄️ SE LẠNH MÙA NOEL – SẮM ĐỒ ẤM ĐI TIỆC, CHECK-IN GIÁNG SINH CỰC GU ✨",
			}
			eventAngle = noelHooks[r.Intn(len(noelHooks))]
		} else if strings.Contains(evt, "tết") || strings.Contains(evt, "năm mới") || strings.Contains(evt, "du xuân") {
			tetHooks := []string{
				"🧧 LÊN ĐỒ DU XUÂN ĐÓN TẾT – PHONG CÁCH QUÝ ÔNG LỊCH LÃM, KHÔNG ĐỤNG HÀNG 🌸",
				"✨ DIỆN ĐỒ ĐÓN NĂM MỚI – TỔNG HỢP NHỮNG MẪU VINTAGE ĐỘC BẢN CỰC CHẤT 🧨",
				"🏮 CHUẨN BỊ OUTFIT CHUẨN GU CHO NHỮNG CHUYẾN DU XUÂN ĐẦU NĂM ✨",
			}
			eventAngle = tetHooks[r.Intn(len(tetHooks))]
		} else if strings.Contains(evt, "tiệc") || strings.Contains(evt, "party") || strings.Contains(evt, "cuối năm") {
			partyHooks := []string{
				"🍷 BẬT VIBE LỊCH LÃM CHO NHỮNG BUỔI TIỆC VÀ TỤ HỌP CUỐI NĂM ✨",
				"✨ OUTFIT SANG TRỌNG – ĐI TIỆC CUỐI NĂM CỰC CUỐN VỚI GU VINTAGE ĐẬM CHẤT 🍸",
				"🎩 NÂNG TẦM OUTFIT YEAR-END PARTY VỚI NHỮNG SIÊU PHẨM TUYỂN CHỌN 🥂",
			}
			eventAngle = partyHooks[r.Intn(len(partyHooks))]
		} else {
			// Generic seasonal
			eventAngle = fmt.Sprintf("❄️ %s – SẮM ĐỒ ẤM VÀ NÂNG TẦM OUTFIT CỰC CHẤT ✨", strings.ToUpper(req.EventName))
		}
	}

	switch code {
	case "GILE":
		return s.buildGileContent(eventAngle, albumLink, r)
	case "TROUSERS":
		return s.buildTrousersContent(eventAngle, albumLink, r)
	case "BLAZER":
		return s.buildBlazerContent(eventAngle, albumLink, r)
	case "JACKET", "HARRINGTON_JACKET":
		return s.buildJacketContent(eventAngle, albumLink, r)
	case "SHIRT":
		return s.buildShirtContent(eventAngle, albumLink, r)
	default:
		return s.buildGenericContent(code, req.CategoryName, eventAngle, albumLink, r)
	}
}

func (s *AIContentService) buildGileContent(eventAngle, albumLink string, r *rand.Rand) string {
	header := "✨ VINTAGE WAISTCOAT / GILE COLLECTION – ĐỒNG GIÁ 120K ✨"
	if eventAngle != "" {
		header = eventAngle + "\n" + header
	}

	descs := []string{
		"🧥 Item cực dễ phối trong những ngày se lạnh — vừa tăng độ ấm, vừa nâng tầm outfit lên một bậc!",
		"🧥 Một chiếc gile chuẩn vintage khoác ngoài sơ mi hay áo thun là đủ để tạo điểm nhấn cổ điển, thanh lịch.",
		"🧥 Cả năm mới có dịp thời tiết se lạnh chuẩn bài để anh em tự tin phối gile nhiều lớp cực có gu.",
	}

	materials := []string{
		"✔️ Form dáng classic chuẩn vintage, độ mới rất cao, chất liệu phong phú từ wool dệt đầm tay, tweed nhẹ đến poly đứng form.",
		"✔️ Chất vải tuyển chọn kỹ càng: wool ấm áp, tweed sang trọng và các chất liệu pha giữ form phẳng phiu suốt ngày dài.",
		"✔️ Đa dạng họa tiết từ kẻ caro, xương cá (herringbone) cho đến màu trơn kinh điển, rất dễ mix-match.",
	}

	return fmt.Sprintf(`%s
-----------------------------------
%s
%s
✔️ Phối cùng sơ mi, áo thun dài tay hoặc khoác thêm blazer bên ngoài đều cực kỳ lịch lãm và chuẩn phong thái quý ông.
-----------------------------------
📐 Thông số đo đạc chi tiết đã ghi trực tiếp trên từng ảnh:
➡️ Ngang ngực — Dài trước (tính đến mũi nhọn vạt) — Dài sau
📩 Cần hướng dẫn cách đo hay tư vấn độ fit – cứ inbox shop gửi ảnh chi tiết nhé!
-----------------------------------
📦 Toàn bộ đã được giặt ủi thơm tho, sạch sẽ, sẵn sàng giao – mỗi mẫu chỉ có 1 chiếc duy nhất!
🔎 Cam kết từ Vintage Always True:
• Bao hàng không lỗi (nếu có lỗi nhỏ shop note rõ trên ảnh)
• Ship COD toàn quốc – hỗ trợ đồng kiểm khi nhận hàng 📦
• Không bao vừa hoặc đúng gu vì form mỗi người khác nhau — anh em nên xem kỹ số đo trước khi chốt nha 💬

👉 Để xem thêm sp và check sp còn hàng tại:
%s

#vintagegile #waistcoat #vintagevest #menswearvietnam #vintagealwaystrue #ao2hand`,
		header,
		descs[r.Intn(len(descs))],
		materials[r.Intn(len(materials))],
		albumLink,
	)
}

func (s *AIContentService) buildTrousersContent(eventAngle, albumLink string, r *rand.Rand) string {
	header := "🔥 VINTAGE WOOL TROUSERS COLLECTION 🔥\n💸 Giá trên từng ảnh, các mẫu không để giá là đồng giá 180k"
	if eventAngle != "" {
		header = eventAngle + "\n" + header
	}

	hooks := []string{
		"Khi tiết trời se lạnh, không gì chuẩn mực hơn một chiếc quần Wool Vintage đậm đà vị thu đông:",
		"Một chiếc quần âu vintage chuẩn classic luôn là nền tảng vững chắc cho mọi set đồ thanh lịch:",
		"Mảnh ghép hoàn hảo cho outfit quý ông mùa lạnh – đứng form, đầm tay và hạn chế nhăn tối đa:",
	}

	return fmt.Sprintf(`%s
——————————————
%s
✨ Chất liệu: Wool dệt đầm tay, giữ ấm tốt nhưng vẫn đảm bảo độ mềm mại, bay bổng tự nhiên trong từng bước đi.
✨ Ưu điểm: Lịch sự, đứng dáng và hạn chế nhăn vượt trội, giữ form chỉn chu cả ngày dài làm việc hay đi chơi.
✨ Form dáng: Classic chuẩn mực, cạp cao tôn dáng, phối cùng sơ mi, sweater, gile hay blazer đều \"hết bài\".
——————————————
📐 Thông số ghi rõ trực tiếp trên ảnh:
➡️ Eo — Đáy — Đùi — Dài (dư gấu) — Ống
📩 Cần hướng dẫn cách đo – cứ inbox mình gửi chi tiết tư vấn ngay nhé!
——————————————
⚠️ Hàng 2hand tuyển chọn kỹ lưỡng – mỗi mẫu chỉ có 1 chiếc duy nhất
🧼 Hình chụp thật 100%% – mô tả rõ ràng – giao hàng nhanh, giặt ủi sạch sẽ
📦 Ship COD toàn quốc – hỗ trợ đồng kiểm tận tay

👉 Để xem thêm sp và check sp còn hàng tại:
%s

#vintagetrousers #quanauvintage #wooltrousers #mensclassicwear #vintagealwaystrue #menswearvietnam`,
		header,
		hooks[r.Intn(len(hooks))],
		albumLink,
	)
}

func (s *AIContentService) buildBlazerContent(eventAngle, albumLink string, r *rand.Rand) string {
	header := "✨ VINTAGE BLAZER COLLECTION ✨\n📐 Thông số & 💸 Giá: đã ghi trực tiếp trên từng ảnh"
	if eventAngle != "" {
		header = eventAngle + "\n" + header
	}

	return fmt.Sprintf(`%s
-----------------------------------
🎯 Toàn blazer cực phẩm: hiệu lớn – vải xịn – chi tiết nghệ thuật – độ mới cực cao!
🔹 Chất liệu đa dạng từ wool mịn, tweed đầm ấm đến vải pha mát nhẹ, lót trong xịn sò.
🔹 Form classic dễ phối – từ phong cách smart-casual hàng ngày đến trang trọng, lịch lãm đều chuẩn đẹp.
-----------------------------------
📐 Thông số ghi trên ảnh:
• Vai – Ngực – Dài vạt trước – Dài tay
📸 Xem ảnh từng mẫu và mô tả chi tiết tại page – mỗi chiếc một vibe riêng biệt, không lo đụng hàng.
-----------------------------------
📦 Giao hàng toàn quốc – đóng gói kỹ càng, sạch sẽ, thơm tho
📩 Inbox ngay để xem thêm ảnh góc cận & tư vấn phối đồ phù hợp
📌 Hàng tuyển chọn – mỗi mẫu chỉ có 1 chiếc – chốt sớm kẻo lỡ!

👉 Để xem thêm sp và check sp còn hàng tại:
%s

#vintageblazer #blazernam #vintagestyle #menswearvietnam #vintagealwaystrue #ao2hand #mensclassicwear`,
		header,
		albumLink,
	)
}

func (s *AIContentService) buildJacketContent(eventAngle, albumLink string, r *rand.Rand) string {
	header := "✨ CLASSIC JACKET COLLECTION ✨\n📐 Thông số & 💸 Giá: đã ghi trực tiếp trên từng ảnh"
	if eventAngle != "" {
		header = eventAngle + "\n" + header
	}

	return fmt.Sprintf(`%s
-----------------------------------
🧥 Đợt áo khoác mùa thu – đông tuyển chọn với đầy đủ thể loại:
• Harrington Jacket
• Field Jacket
• Classic Coat & Bomber Jacket
✔️ Dễ phối, tôn dáng và giữ ấm cực tốt trong những ngày trở gió.
✔️ Độ mới rất cao, đường may chỉn chu, khóa kéo và nút áo zin đét.
-----------------------------------
📐 Thông số & 💸 Giá: Đã ghi trực tiếp trên từng ảnh.
⚠️ Lưu ý: Áo khoác có nhiều loại lót (trần bông, lót nỉ, lót dù) nên số đo có thể chênh lệch so với cảm giác mặc thực tế. Anh em nên inbox để shop tư vấn chiều cao cân nặng chuẩn xác nhất nhé!
-----------------------------------
📸 Ảnh chụp thật 100%% góc sáng rõ ràng, chi tiết từng đường nét.
📦 Ship toàn quốc – áo đã giặt ủi sạch sẽ, đóng gói cẩn thận.
👉 Hàng vintage tuyển chọn mỗi chiếc chỉ có 1 cái duy nhất!

👉 Để xem thêm sp và check sp còn hàng tại:
%s

#vintagejacket #harringtonjacket #ao2hand #menswearvietnam #vintagealwaystrue`,
		header,
		albumLink,
	)
}

func (s *AIContentService) buildShirtContent(eventAngle, albumLink string, r *rand.Rand) string {
	header := "👔 VINTAGE SHIRT COLLECTION 👔\n💸 Giá trên từng ảnh, các mẫu không để giá là đồng giá theo album"
	if eventAngle != "" {
		header = eventAngle + "\n" + header
	}

	return fmt.Sprintf(`%s
-----------------------------------
✨ Những mẫu sơ mi vintage với chất vải bền bỉ, đường kim mũi chỉ tinh tế và họa tiết cổ điển đậm chất:
✔️ Chất liệu cotton tự nhiên, đũi mát hoặc oxford dày dặn, thấm hút mồ hôi tốt.
✔️ Form dáng classic thoải mái, diện đơn lẻ hay khoác cùng gile/blazer đều tôn nét nam tính, đĩnh đạc.
-----------------------------------
📐 Thông số đo đạc ghi chi tiết trên từng ảnh:
➡️ Vai — Ngực — Dài áo — Dài tay
📩 Inbox shop để nhận tư vấn phối cùng quần âu hoặc áo khoác chuẩn bài nhé!
-----------------------------------
🧼 Áo đã được giặt là sạch đẹp, thơm tho, sẵn sàng lên đồ ngay khi nhận hàng.
📦 Ship COD toàn quốc – hỗ trợ đồng kiểm tận nơi.

👉 Để xem thêm sp và check sp còn hàng tại:
%s

#vintageshirt #aosomivintage #sominam #menswearvietnam #vintagealwaystrue`,
		header,
		albumLink,
	)
}

func (s *AIContentService) buildGenericContent(code, name, eventAngle, albumLink string, r *rand.Rand) string {
	displayTitle := name
	if displayTitle == "" {
		displayTitle = code
	}
	header := fmt.Sprintf("✨ VINTAGE %s COLLECTION ✨", strings.ToUpper(displayTitle))
	if eventAngle != "" {
		header = eventAngle + "\n" + header
	}

	return fmt.Sprintf(`%s
-----------------------------------
🎯 Những món đồ vintage chuẩn phong cách cổ điển, tuyển chọn kỹ lưỡng từng chiếc:
✔️ Độ mới cao, chất liệu xịn sò và thiết kế độc bản không đụng hàng.
✔️ Dễ phối đồ, mang lại phong thái chỉn chu, có gu cho anh em.
-----------------------------------
📐 Thông số & Giá: Đã ghi trực tiếp trên từng ảnh.
📩 Cần thêm thông tin hay hình ảnh chi tiết – anh em cứ inbox shop hỗ trợ nhanh nhất nhé!
-----------------------------------
📦 Đã giặt ủi sạch sẽ, đóng gói cẩn thận – mỗi mẫu chỉ có 1 chiếc duy nhất!
🧼 Ship COD toàn quốc – hỗ trợ đồng kiểm khi nhận hàng.

👉 Để xem thêm sp và check sp còn hàng tại:
%s

#vintagefashion #menswearvietnam #vintagealwaystrue`,
		header,
		albumLink,
	)
}
