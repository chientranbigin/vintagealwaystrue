package handlers

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"sale-app/services"

	"github.com/gin-gonic/gin"
)

func (h *CopilotHandler) getFBVerifyToken() string {
	tok := strings.TrimSpace(os.Getenv("FB_VERIFY_TOKEN"))
	if tok == "" {
		return "vintage-copilot-verify-2026"
	}
	return tok
}

func (h *CopilotHandler) getMainPageCredentials() (string, string) {
	pageID := strings.TrimSpace(h.Cfg.FBMainPageID)
	token := strings.TrimSpace(h.Cfg.FBMainPageToken)
	if pageID == "" {
		pageID = strings.TrimSpace(h.Cfg.FBPageID)
	}
	if token == "" {
		token = strings.TrimSpace(h.Cfg.FBPageToken)
	}
	return pageID, token
}

// VerifyWebhook handles Facebook Messenger Webhook verification (GET /webhook/facebook).
func (h *CopilotHandler) VerifyWebhook(c *gin.Context) {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	if mode == "subscribe" && token == h.getFBVerifyToken() {
		c.String(http.StatusOK, challenge)
		return
	}
	c.String(http.StatusForbidden, "Verification failed")
}

// ReceiveWebhook handles incoming Facebook Messenger messages & echoes (POST /webhook/facebook).
func (h *CopilotHandler) ReceiveWebhook(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	// Verify X-Hub-Signature-256 if FB_APP_SECRET is configured
	if appSecret := strings.TrimSpace(os.Getenv("FB_APP_SECRET")); appSecret != "" {
		sigHeader := c.GetHeader("X-Hub-Signature-256")
		if !verifyMetaSignature(body, sigHeader, appSecret) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
			return
		}
	}

	var payload fbWebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if payload.Object != "page" {
		c.Status(http.StatusOK)
		return
	}

	// Acknowledge immediately so Meta never retries due to timeout, then process async
	go h.processFBWebhookPayload(payload)
	c.String(http.StatusOK, "EVENT_RECEIVED")
}

func verifyMetaSignature(body []byte, sigHeader string, appSecret string) bool {
	if !strings.HasPrefix(sigHeader, "sha256=") {
		return false
	}
	expectedHex := strings.TrimPrefix(sigHeader, "sha256=")
	mac := hmac.New(sha256.New, []byte(appSecret))
	mac.Write(body)
	actualHex := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expectedHex), []byte(actualHex))
}

type fbWebhookPayload struct {
	Object string `json:"object"`
	Entry  []struct {
		ID        string `json:"id"`
		Time      int64  `json:"time"`
		Messaging []struct {
			Sender struct {
				ID string `json:"id"`
			} `json:"sender"`
			Recipient struct {
				ID string `json:"id"`
			} `json:"recipient"`
			Timestamp int64 `json:"timestamp"`
			Message   *struct {
				MID         string `json:"mid"`
				IsEcho      bool   `json:"is_echo"`
				Text        string `json:"text"`
				StickerID   int64  `json:"sticker_id"`
				Attachments []struct {
					Type    string `json:"type"`
					Payload struct {
						URL       string `json:"url"`
						StickerID int64  `json:"sticker_id"`
					} `json:"payload"`
				} `json:"attachments"`
			} `json:"message"`
		} `json:"messaging"`
	} `json:"entry"`
}

func (h *CopilotHandler) processFBWebhookPayload(payload fbWebhookPayload) {
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	pageID, pageToken := h.getMainPageCredentials()
	historyLoaded := map[string]bool{}
	triggerNorm := strings.ToLower(strings.TrimSpace(os.Getenv("COPILOT_ORDER_TRIGGER")))
	if triggerNorm == "" {
		triggerNorm = "ok ạ, để em lên đơn ngay"
	}

	for _, entry := range payload.Entry {
		for _, event := range entry.Messaging {
			if event.Message == nil {
				continue
			}

			msg := event.Message
			senderRole := "customer"
			psid := event.Sender.ID
			if msg.IsEcho {
				senderRole = "shop"
				psid = event.Recipient.ID
			}
			if psid == "" {
				continue
			}

			// First time we see this customer: import their recent thread from Graph API so the
			// AI has the earlier context. The current message is part of it (deduped by mid).
			importedHistory := false
			if !historyLoaded[psid] && !h.Copilot.HasFBConversation(psid) {
				historyLoaded[psid] = true
				if err := h.Copilot.ImportFBHistory(ctx, pageID, psid, pageToken); err != nil {
					log.Printf("[copilot] import history for %s failed: %v", psid, err)
				} else {
					importedHistory = true
				}
			}

			msgTime := time.Now()
			if event.Timestamp > 0 {
				msgTime = time.UnixMilli(event.Timestamp)
			}

			insertedAny := false
			shouldCreateOrder := false

			text := strings.TrimSpace(msg.Text)
			if services.IsFBSystemNotice(text) {
				text = ""
			}
			// Recognize Facebook thumb-up Like stickers
			if text == "" && isLikeSticker(msg.StickerID) {
				text = "👍"
			}

			if text != "" {
				_, inserted, _ := h.Copilot.UpsertFBMessage(
					psid, "", "", msg.MID, senderRole, "text", text, "", msgTime, pageToken,
				)
				if inserted {
					insertedAny = true
					if senderRole == "shop" && strings.Contains(strings.ToLower(text), triggerNorm) {
						shouldCreateOrder = true
					}
				}
			}

			for idx, att := range msg.Attachments {
				if att.Type == "image" && att.Payload.URL != "" {
					attID := fmt.Sprintf("%s_att_%d", msg.MID, idx)
					if isLikeSticker(att.Payload.StickerID) {
						_, inserted, _ := h.Copilot.UpsertFBMessage(
							psid, "", "", attID, senderRole, "text", "👍", "", msgTime, pageToken,
						)
						if inserted {
							insertedAny = true
						}
						continue
					}
					_, inserted, _ := h.Copilot.UpsertFBMessage(
						psid, "", "", attID, senderRole, "image", "", att.Payload.URL, msgTime, pageToken,
					)
					if inserted {
						insertedAny = true
					}
				}
			}

			// The imported history may already contain this message (same mid), so still rebuild.
			if !insertedAny && !importedHistory {
				continue
			}

			action := "suggest"
			if shouldCreateOrder {
				action = "create_order"
			}

			sess, err := h.Copilot.RebuildAndAnalyzeSession(ctx, psid, action)
			if err != nil || sess == nil {
				continue
			}

			if shouldCreateOrder && sess.OrderDraft != nil && len(sess.OrderDraft.ProductCodes) > 0 {
				if order, err := h.createOrderFromDraft(sess.OrderDraft); err == nil && order != nil {
					sess.CreatedOrder = order
					h.Copilot.MarkConversationOrdered(psid, order.ID)
				}
			}

			h.Copilot.SaveAndBroadcast(sess)
		}
	}
}

func isLikeSticker(stickerID int64) bool {
	return services.IsLikeStickerID(stickerID)
}

// SyncFBConversations pulls recent Messenger threads from Facebook Graph API.
func (h *CopilotHandler) SyncFBConversations(c *gin.Context) {
	pageID, pageToken := h.getMainPageCredentials()
	count, err := h.Copilot.SyncRecentFBConversations(c.Request.Context(), pageID, pageToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"synced":   count,
		"sessions": h.Copilot.ListSessions(),
	})
}

// PrivacyPolicy serves the public Privacy Policy page required for Meta App Review.
func (h *CopilotHandler) PrivacyPolicy(c *gin.Context) {
	html := `<!DOCTYPE html>
<html lang="vi">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Chính sách bảo mật - Vintage Always True</title>
  <style>
    body { font-family: system-ui, -apple-system, sans-serif; max-width: 780px; margin: 40px auto; padding: 0 20px; line-height: 1.7; color: #1e293b; }
    h1 { color: #0f172a; }
    h2 { margin-top: 24px; color: #334155; font-size: 18px; }
    .box { background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 10px; padding: 16px; margin: 16px 0; }
  </style>
</head>
<body>
  <h1>Chính Sách Bảo Mật (Privacy Policy)</h1>
  <div class="box">
    <strong>Đơn vị vận hành:</strong> Hộ kinh doanh Vintage Always True<br/>
    <strong>Website:</strong> https://vintagealwaystrue.id.vn<br/>
    <strong>Email liên hệ:</strong> support@vintagealwaystrue.id.vn
  </div>
  <h2>1. Dữ liệu chúng tôi tiếp nhận</h2>
  <p>Khi khách hàng chủ động nhắn tin tới Trang Facebook chính thức của Vintage Always True, hệ thống nội bộ chỉ tiếp nhận: tên hiển thị, ảnh đại diện công khai, PSID cuộc hội thoại, nội dung tin nhắn và hình ảnh sản phẩm mà khách hàng gửi.</p>
  <h2>2. Mục đích sử dụng</h2>
  <p>Dữ liệu chỉ được sử dụng cho mục đích chăm sóc khách hàng nội bộ (Agent Assist): kiểm tra tình trạng tồn kho sản phẩm, tư vấn kích thước phù hợp và tạo đơn hàng theo yêu cầu của khách hàng. Chúng tôi tuyệt đối không gửi tin nhắn quảng cáo hàng loạt và không bán hay chia sẻ dữ liệu cho bên thứ ba.</p>
  <h2>3. Lưu trữ và xóa dữ liệu</h2>
  <p>Dữ liệu hội thoại được lưu trữ bảo mật trên hệ thống quản trị nội bộ của cửa hàng. Khách hàng có thể yêu cầu xóa toàn bộ dữ liệu bất cứ lúc nào theo hướng dẫn tại <a href="/data-deletion">/data-deletion</a>.</p>
</body>
</html>`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// DataDeletion serves the Data Deletion instructions (GET) and Meta callback (POST).
func (h *CopilotHandler) DataDeletion(c *gin.Context) {
	if c.Request.Method == http.MethodPost {
		code := fmt.Sprintf("del_%d", time.Now().Unix())
		c.JSON(http.StatusOK, gin.H{
			"url":               "https://vintagealwaystrue.id.vn/data-deletion",
			"confirmation_code": code,
		})
		return
	}

	html := `<!DOCTYPE html>
<html lang="vi">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Hướng dẫn xóa dữ liệu - Vintage Always True</title>
  <style>
    body { font-family: system-ui, -apple-system, sans-serif; max-width: 780px; margin: 40px auto; padding: 0 20px; line-height: 1.7; color: #1e293b; }
    h1 { color: #0f172a; }
  </style>
</head>
<body>
  <h1>Hướng Dẫn Xóa Dữ Liệu Người Dùng (User Data Deletion)</h1>
  <p>Nếu bạn đã từng nhắn tin với Trang Facebook của <strong>Vintage Always True</strong> và muốn xóa toàn bộ dữ liệu hội thoại lưu trữ trên hệ thống quản trị nội bộ của chúng tôi, vui lòng thực hiện một trong hai cách sau:</p>
  <ol>
    <li>Nhắn tin trực tiếp tới Fanpage với nội dung: <strong>"Yêu cầu xóa dữ liệu hội thoại"</strong>.</li>
    <li>Gửi email tới <strong>support@vintagealwaystrue.id.vn</strong> kèm tên tài khoản Facebook của bạn.</li>
  </ol>
  <p>Yêu cầu của bạn sẽ được xử lý và xóa vĩnh viễn khỏi cơ sở dữ liệu trong vòng 24–48 giờ.</p>
</body>
</html>`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// TermsOfService serves the public Terms of Service page required by Meta App Settings.
func (h *CopilotHandler) TermsOfService(c *gin.Context) {
	html := `<!DOCTYPE html>
<html lang="vi">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Điều khoản dịch vụ - Vintage Always True</title>
  <style>
    body { font-family: system-ui, -apple-system, sans-serif; max-width: 780px; margin: 40px auto; padding: 0 20px; line-height: 1.7; color: #1e293b; }
    h1 { color: #0f172a; }
    h2 { margin-top: 24px; color: #334155; font-size: 18px; }
    .box { background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 10px; padding: 16px; margin: 16px 0; }
  </style>
</head>
<body>
  <h1>Điều Khoản Dịch Vụ (Terms of Service)</h1>
  <div class="box">
    <strong>Đơn vị vận hành:</strong> Hộ kinh doanh Vintage Always True<br/>
    <strong>Website:</strong> https://vintagealwaystrue.id.vn<br/>
    <strong>Email liên hệ:</strong> support@vintagealwaystrue.id.vn
  </div>
  <h2>1. Phạm vi áp dụng</h2>
  <p>Hệ thống quản trị nội bộ và công cụ hỗ trợ tư vấn khách hàng của Vintage Always True được sử dụng để tiếp nhận thông tin sản phẩm, kiểm tra tồn kho, tư vấn kích thước và xử lý đơn hàng khi khách hàng chủ động liên hệ qua Trang Facebook chính thức.</p>
  <h2>2. Quyền và trách nhiệm</h2>
  <p>Chúng tôi cam kết sử dụng thông tin liên hệ và đơn hàng đúng mục đích giao nhận hàng hóa, bảo mật thông tin khách hàng và tuân thủ đầy đủ chính sách nền tảng của Meta.</p>
  <h2>3. Liên hệ</h2>
  <p>Mọi thắc mắc về đơn hàng hoặc điều khoản dịch vụ, vui lòng liên hệ trực tiếp qua Fanpage Vintage Always True hoặc email support@vintagealwaystrue.id.vn.</p>
</body>
</html>`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// SendReceipt sends the rendered receipt PNG to the customer's Messenger thread via the Send API.
func (h *CopilotHandler) SendReceipt(c *gin.Context) {
	var req struct {
		ThreadID    string `json:"thread_id" binding:"required"`
		ImageBase64 string `json:"image_base64" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if h.isInvalidThreadID(req.ThreadID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hội thoại không hợp lệ"})
		return
	}
	raw := req.ImageBase64
	if i := strings.Index(raw, ","); i >= 0 && strings.HasPrefix(raw, "data:") {
		raw = raw[i+1:]
	}
	img, err := base64.StdEncoding.DecodeString(raw)
	if err != nil || len(img) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ảnh receipt không hợp lệ"})
		return
	}
	_, pageToken := h.getMainPageCredentials()
	if err := h.Copilot.SendFBImage(c.Request.Context(), req.ThreadID, img, "receipt.png", pageToken); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// SendText sends a (suggested / edited) reply to the customer's Messenger thread.
func (h *CopilotHandler) SendText(c *gin.Context) {
	var req struct {
		ThreadID string `json:"thread_id" binding:"required"`
		Text     string `json:"text" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	text := strings.TrimSpace(req.Text)
	if text == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nội dung trống"})
		return
	}
	if len([]rune(text)) > 2000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tin nhắn quá dài (tối đa 2000 ký tự)"})
		return
	}
	if h.isInvalidThreadID(req.ThreadID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hội thoại không hợp lệ"})
		return
	}
	_, pageToken := h.getMainPageCredentials()
	if err := h.Copilot.SendFBText(c.Request.Context(), req.ThreadID, text, pageToken); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
