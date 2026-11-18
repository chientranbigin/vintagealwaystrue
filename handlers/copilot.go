package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"sale-app/config"
	"sale-app/models"
	"sale-app/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CopilotHandler struct {
	DB      *gorm.DB
	Cfg     *config.Config
	Copilot *services.CopilotService
	FB      *services.FacebookService
}

func NewCopilotHandler(db *gorm.DB, cfg *config.Config, fb *services.FacebookService) *CopilotHandler {
	h := &CopilotHandler{
		DB:      db,
		Cfg:     cfg,
		Copilot: services.NewCopilotService(db, services.NewVisionService(), cfg.MediaBaseURL),
		FB:      fb,
	}
	// Page ids are never a conversation id; drop sessions wrongly stored under them.
	h.Copilot.PurgeThreads(cfg.FBPageID, cfg.FBMainPageID)
	return h
}

// isInvalidThreadID rejects ids that would merge/split conversations wrongly:
// empty ids and the Facebook Page id (asset_id) that old extension builds sent.
func (h *CopilotHandler) isInvalidThreadID(id string) bool {
	id = strings.TrimSpace(id)
	if id == "" {
		return true
	}
	return (h.Cfg.FBPageID != "" && id == h.Cfg.FBPageID) || (h.Cfg.FBMainPageID != "" && id == h.Cfg.FBMainPageID)
}

// Ingest receives conversation snapshots & images from the Chrome Extension.
func (h *CopilotHandler) Ingest(c *gin.Context) {
	expectedSecret := strings.TrimSpace(os.Getenv("COPILOT_SECRET"))
	if expectedSecret != "" {
		provided := strings.TrimSpace(c.GetHeader("X-Copilot-Secret"))
		if provided != expectedSecret {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid X-Copilot-Secret"})
			return
		}
	}

	var req services.CopilotIngestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload: " + err.Error()})
		return
	}

	if h.isInvalidThreadID(req.ThreadID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Thiếu ID hội thoại (selected_item_id), hãy reload extension"})
		return
	}

	sess := h.Copilot.AnalyzeIngest(c.Request.Context(), req)

	if req.Action == "create_order" && sess.OrderDraft != nil && len(sess.OrderDraft.ProductCodes) > 0 {
		order, err := h.createOrderFromDraft(sess.OrderDraft)
		if err == nil && order != nil {
			sess.CreatedOrder = order
		}
	}

	h.Copilot.SaveAndBroadcast(sess)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"session": sess,
	})
}

// ListSessions returns recent copilot sessions for the /salev2/copilot screen.
func (h *CopilotHandler) ListSessions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"sessions": h.Copilot.ListSessions(),
	})
}

// Stream pushes real-time updates to /salev2/copilot via Server-Sent Events (SSE).
func (h *CopilotHandler) Stream(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")

	ch := h.Copilot.Subscribe()
	defer h.Copilot.Unsubscribe(ch)

	ticker := time.NewTicker(25 * time.Second)
	defer ticker.Stop()

	c.Stream(func(w io.Writer) bool {
		select {
		case <-c.Request.Context().Done():
			return false
		case sess, ok := <-ch:
			if !ok {
				return false
			}
			b, err := json.Marshal(sess)
			if err != nil {
				return true
			}
			fmt.Fprintf(w, "event: session\ndata: %s\n\n", string(b))
			return true
		case <-ticker.C:
			fmt.Fprintf(w, "event: ping\ndata: {}\n\n")
			return true
		}
	})
}

// CreateOrderFromSession allows manual 1-click order creation from the Copilot UI.
func (h *CopilotHandler) CreateOrderFromSession(c *gin.Context) {
	var req struct {
		ThreadID     string   `json:"thread_id" binding:"required"`
		Name         string   `json:"name"`
		Phone        string   `json:"phone"`
		Address      string   `json:"address"`
		Note         string   `json:"note"`
		ProductCodes []string `json:"product_codes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	draft := &services.CopilotOrderDraft{
		Name:         strings.TrimSpace(req.Name),
		Phone:        strings.TrimSpace(req.Phone),
		Address:      strings.TrimSpace(req.Address),
		Note:         strings.TrimSpace(req.Note),
		ProductCodes: req.ProductCodes,
	}

	order, err := h.createOrderFromDraft(draft)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if sess, ok := h.Copilot.GetSession(req.ThreadID); ok {
		sess.CreatedOrder = order
		sess.OrderDraft = draft
		sess.UpdatedAt = time.Now()
		h.Copilot.SaveAndBroadcast(sess)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"order":   order,
	})
}

func (h *CopilotHandler) createOrderFromDraft(draft *services.CopilotOrderDraft) (*models.Order, error) {
	if draft == nil || len(draft.ProductCodes) == 0 {
		return nil, fmt.Errorf("không có mã sản phẩm khả dụng để lên đơn")
	}

	var matchedProducts []models.Product
	var cleanProductIDs []uint
	seenIDs := make(map[uint]bool)

	for _, code := range draft.ProductCodes {
		code = strings.TrimSpace(code)
		if code == "" {
			continue
		}
		var prod models.Product
		// Double guard: ONLY pick AVAILABLE products so old sold items from previous weeks are never included
		if err := h.DB.Where("UPPER(name) = ? AND status = 'AVAILABLE'", strings.ToUpper(code)).
			Order("id DESC").
			First(&prod).Error; err != nil {
			if err2 := h.DB.Where("LOWER(name) LIKE ? AND status = 'AVAILABLE'", "%"+strings.ToLower(code)+"%").
				Order("id DESC").
				First(&prod).Error; err2 != nil {
				continue
			}
		}
		if !seenIDs[prod.ID] {
			seenIDs[prod.ID] = true
			matchedProducts = append(matchedProducts, prod)
			cleanProductIDs = append(cleanProductIDs, prod.ID)
		}
	}

	if len(matchedProducts) == 0 {
		return nil, fmt.Errorf("các mã sản phẩm đều đã bán hoặc không tồn tại")
	}

	var itemsTotal int64
	for _, p := range matchedProducts {
		if p.Price != nil {
			itemsTotal += *p.Price
		}
	}

	cleanPhone := strings.ReplaceAll(strings.TrimSpace(draft.Phone), " ", "")
	status := "NEW"
	var order models.Order

	txErr := h.DB.Transaction(func(tx *gorm.DB) error {
		var consumerID *int
		if cleanPhone != "" {
			var cons models.Consumer
			if err := tx.Where("phone = ?", cleanPhone).First(&cons).Error; err != nil {
				cons = models.Consumer{
					Name:    strPtr(draft.Name),
					Phone:   strPtr(cleanPhone),
					Address: strPtr(draft.Address),
				}
				if err := tx.Create(&cons).Error; err == nil {
					cID := int(cons.ID)
					consumerID = &cID
				}
			} else {
				cID := int(cons.ID)
				consumerID = &cID
			}
		}

		order = models.Order{
			ConsumerID:  consumerID,
			Name:        strPtr(draft.Name),
			Phone:       strPtr(cleanPhone),
			Address:     strPtr(draft.Address),
			Note:        strPtr(draft.Note),
			Status:      strPtr(status),
			TotalAmount: int64Ptr(itemsTotal),
			FinalAmount: int64Ptr(itemsTotal),
		}
		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		for _, pID := range cleanProductIDs {
			prodID := pID
			oID := order.ID
			op := models.OrderProduct{
				OrderID:   &oID,
				ProductID: &prodID,
			}
			if err := tx.Create(&op).Error; err != nil {
				return err
			}
			if err := tx.Model(&models.Product{}).Where("id = ?", prodID).Update("status", "SOLD").Error; err != nil {
				return err
			}
		}
		return nil
	})
	if txErr != nil {
		return nil, txErr
	}

	if h.FB != nil {
		for _, pID := range cleanProductIDs {
			h.FB.Dispatch(pID, services.FBActionSold)
		}
	}

	// Reload with Products for Receipt rendering
	h.DB.Preload("Products").First(&order, order.ID)
	return &order, nil
}
