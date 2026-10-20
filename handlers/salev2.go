package handlers

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"sale-app/config"
	"sale-app/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SaleV2Handler struct {
	DB  *gorm.DB
	Cfg *config.Config
}

func NewSaleV2Handler(db *gorm.DB, cfg *config.Config) *SaleV2Handler {
	return &SaleV2Handler{DB: db, Cfg: cfg}
}

func (h *SaleV2Handler) FormatImagePath(path string) string {
	if path == "" {
		return ""
	}
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}
	clean := strings.TrimLeft(path, "/")
	if strings.HasPrefix(clean, "storage/") {
		clean = strings.TrimPrefix(clean, "storage/")
	}
	if clean == "" {
		return ""
	}
	return h.Cfg.MediaBaseURL + clean
}

func (h *SaleV2Handler) formatProduct(p *models.Product) {
	if p.PathThumb != nil && *p.PathThumb != "" {
		formatted := h.FormatImagePath(*p.PathThumb)
		p.PathThumb = &formatted
	}
	if p.ImageThumbScale != nil && *p.ImageThumbScale != "" {
		formatted := h.FormatImagePath(*p.ImageThumbScale)
		p.ImageThumbScaleURL = &formatted
	}
	for i := range p.Images {
		p.Images[i].FilePath = h.FormatImagePath(p.Images[i].FilePath)
	}
}

// LaravelPagination matches Laravel's JSON pagination structure
type LaravelPagination struct {
	CurrentPage int         `json:"current_page"`
	Data        interface{} `json:"data"`
	From        int         `json:"from"`
	LastPage    int         `json:"last_page"`
	PerPage     int         `json:"per_page"`
	To          int         `json:"to"`
	Total       int64       `json:"total"`
}

func (h *SaleV2Handler) Products(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "60"))
	if limit < 1 {
		limit = 60
	}
	offset := (page - 1) * limit

	query := h.DB.Model(&models.Product{}).Preload("Images").Preload("Sizes")

	// Search filter
	if search := c.Query("search"); search != "" {
		names := strings.Split(search, ",")
		var conds []string
		var args []interface{}
		for _, n := range names {
			n = strings.TrimSpace(n)
			if n != "" {
				conds = append(conds, "name LIKE ?")
				args = append(args, "%"+n+"%")
			}
		}
		if len(conds) > 0 {
			query = query.Where(strings.Join(conds, " OR "), args...)
		}
	}

	// Status filter
	if status := c.Query("status"); status != "" {
		query = query.Where("products.status = ?", status)
		if status == "AVAILABLE" {
			query = query.Where("path_thumb IS NOT NULL AND path_thumb != ''")
		}
	}

	// Type filter
	if pType := c.Query("type"); pType != "" {
		query = query.Where("type = ?", pType)
	}

	var total int64
	query.Count(&total)

	// Sorting
	if c.Query("sort") == "upload" {
		query = query.Order("created_at DESC")
	} else {
		query = query.Order("id DESC")
	}

	var products []models.Product
	query.Limit(limit).Offset(offset).Find(&products)

	for i := range products {
		h.formatProduct(&products[i])
	}

	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	if lastPage == 0 {
		lastPage = 1
	}

	c.JSON(http.StatusOK, LaravelPagination{
		CurrentPage: page,
		Data:        products,
		From:        offset + 1,
		LastPage:    lastPage,
		PerPage:     limit,
		To:          offset + len(products),
		Total:       total,
	})
}

func (h *SaleV2Handler) ProductDetail(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := h.DB.Preload("Images").Preload("Sizes").First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	h.formatProduct(&product)

	// Fetch latest upload time from logs
	var latestUpload *time.Time
	h.DB.Model(&models.ProductUploadLog{}).
		Where("product_id = ?", product.ID).
		Select("MAX(created_at)").
		Scan(&latestUpload)
	if latestUpload != nil {
		formatted := latestUpload.Format("2006-01-02 15:04:05")
		product.LatestUpload = &formatted
	}

	c.JSON(http.StatusOK, product)
}

func (h *SaleV2Handler) Orders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "12"))
	if limit < 1 {
		limit = 12
	}
	offset := (page - 1) * limit

	query := h.DB.Model(&models.Order{}).Preload("Products.Images").Preload("Products.Sizes")

	if status := c.Query("status"); status != "" {
		if status == "NEW" {
			query = query.Where("status != ? AND status != ?", "COMPLETED", "READY FOR SHIPPING").Order("status ASC")
		} else {
			query = query.Where("status = ?", status)
		}
	}

	if search := c.Query("search"); search != "" {
		query = query.Where("name LIKE ? OR phone LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	var total int64
	query.Count(&total)

	var orders []models.Order
	query.Order("id DESC").Limit(limit).Offset(offset).Find(&orders)

	now := time.Now()
	for i := range orders {
		o := &orders[i]
		if o.CreatedAt != nil {
			o.IsOld = now.Sub(*o.CreatedAt).Hours() >= 72 // 3 days
			o.CreatedAtHuman = o.CreatedAt.Format("02/01/2006 15:04")
		}

		if o.IsPaid {
			o.ViettelAmount = 0
		} else {
			tot := int64(0)
			dep := int64(0)
			if o.TotalAmount != nil {
				tot = *o.TotalAmount
			}
			if o.DepositAmount != nil {
				dep = *o.DepositAmount
			}
			o.ViettelAmount = tot - dep
		}

		for j := range o.Products {
			h.formatProduct(&o.Products[j])
		}
	}

	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	if lastPage == 0 {
		lastPage = 1
	}

	c.JSON(http.StatusOK, LaravelPagination{
		CurrentPage: page,
		Data:        orders,
		From:        offset + 1,
		LastPage:    lastPage,
		PerPage:     limit,
		To:          offset + len(orders),
		Total:       total,
	})
}

func (h *SaleV2Handler) OrderDetail(c *gin.Context) {
	id := c.Param("id")
	var order models.Order
	if err := h.DB.Preload("Products.Images").Preload("Products.Sizes").First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	for j := range order.Products {
		h.formatProduct(&order.Products[j])
	}

	c.JSON(http.StatusOK, order)
}

func (h *SaleV2Handler) StoreOrder(c *gin.Context) {
	var req struct {
		Name             string                 `json:"name"`
		Phone            string                 `json:"phone"`
		Address          string                 `json:"address"`
		Note             string                 `json:"note"`
		Desc             string                 `json:"desc"`
		TotalAmount      int64                  `json:"total_amount"`
		DepositAmount    int64                  `json:"deposit_amount"`
		AdditionalAmount int64                  `json:"additional_amount"`
		FreeshipAmount   float64                `json:"freeship_amount"`
		IsFreeship       bool                   `json:"is_freeship"`
		IsPaid           bool                   `json:"is_paid"`
		Status           string                 `json:"status"`
		Products         []uint                 `json:"products"`
		ProductsInfo     []map[string]interface{} `json:"products_info"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order := models.Order{
		Name:             &req.Name,
		Phone:            &req.Phone,
		Address:          &req.Address,
		Note:             &req.Note,
		Desc:             &req.Desc,
		TotalAmount:      &req.TotalAmount,
		DepositAmount:    &req.DepositAmount,
		AdditionalAmount: &req.AdditionalAmount,
		FreeshipAmount:   &req.FreeshipAmount,
		IsFreeship:       req.IsFreeship,
		IsPaid:           req.IsPaid,
		Status:           &req.Status,
	}

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		for _, pID := range req.Products {
			op := models.OrderProduct{OrderID: &order.ID, ProductID: &pID}
			if err := tx.Create(&op).Error; err != nil {
				return err
			}
		}

		// Update product status to SOLD
		soldStatus := "SOLD"
		if err := tx.Model(&models.Product{}).Where("id IN ?", req.Products).Update("status", soldStatus).Error; err != nil {
			return err
		}

		// Price updates if any
		for _, info := range req.ProductsInfo {
			if idVal, ok := info["id"]; ok {
				if priceVal, ok2 := info["price"]; ok2 {
					idFloat, _ := idVal.(float64)
					priceFloat, _ := priceVal.(float64)
					tx.Model(&models.Product{}).Where("id = ?", uint(idFloat)).Update("price", int64(priceFloat))
				}
			}
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "order": order})
}

func (h *SaleV2Handler) CustomerHistory(c *gin.Context) {
	phone := c.Param("phone")
	var orders []models.Order
	h.DB.Where("phone = ?", phone).Preload("Products.Images").Order("id DESC").Find(&orders)
	for i := range orders {
		for j := range orders[i].Products {
			h.formatProduct(&orders[i].Products[j])
		}
	}
	c.JSON(http.StatusOK, orders)
}

func (h *SaleV2Handler) Consumers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	var total int64
	h.DB.Model(&models.Consumer{}).Count(&total)

	var consumers []models.Consumer
	h.DB.Order("id DESC").Limit(limit).Offset(offset).Find(&consumers)

	c.JSON(http.StatusOK, LaravelPagination{
		CurrentPage: page,
		Data:        consumers,
		From:        offset + 1,
		LastPage:    int(math.Ceil(float64(total) / float64(limit))),
		To:          offset + len(consumers),
		Total:       total,
	})
}

func (h *SaleV2Handler) StoreProduct(c *gin.Context) {
	var req struct {
		Type        string `json:"type" form:"type" binding:"required"`
		Price       int64  `json:"price" form:"price"`
		Description string `json:"description" form:"description"`
		Sizes       []struct {
			Name  string   `json:"name"`
			Value *float64 `json:"value"`
		} `json:"sizes"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prefixMap := map[string]string{
		"TROUSERS":          "TS",
		"HARRINGTON_JACKET": "HJ",
		"JACKET":            "JK",
		"SHIRT":             "ST",
		"BLAZER":            "BZ",
		"TIE":               "TI",
		"BELT":              "BE",
		"GILE":              "GL",
		"POLO SHIRT":        "PS",
		"HAT":               "HA",
		"SUIT":              "SU",
	}

	prefix := req.Type
	if p, ok := prefixMap[req.Type]; ok {
		prefix = p
	}

	var last models.Product
	h.DB.Where("type = ?", req.Type).Order("id DESC").First(&last)
	nextID := last.ID + 1
	codeName := fmt.Sprintf("%s%d", prefix, nextID)
	status := "AVAILABLE"

	product := models.Product{
		Type:        &req.Type,
		Name:        &codeName,
		Price:       &req.Price,
		Description: &req.Description,
		Status:      &status,
	}

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&product).Error; err != nil {
			return err
		}

		for _, sz := range req.Sizes {
			if sz.Value != nil {
				s := models.Size{
					ProductID: &product.ID,
					Name:      &sz.Name,
					Value:     sz.Value,
				}
				if err := tx.Create(&s).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.DB.Preload("Sizes").First(&product, product.ID)
	c.JSON(http.StatusOK, gin.H{"success": true, "product": product})
}

func (h *SaleV2Handler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := h.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var req struct {
		Price       *int64  `json:"price" form:"price"`
		Description *string `json:"description" form:"description"`
		Type        *string `json:"type" form:"type"`
		Sizes       []struct {
			Name  string   `json:"name"`
			Value *float64 `json:"value"`
		} `json:"sizes"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if req.Price != nil {
			product.Price = req.Price
		}
		if req.Description != nil {
			product.Description = req.Description
		}
		if req.Type != nil && *req.Type != "" {
			product.Type = req.Type
		}

		if err := tx.Save(&product).Error; err != nil {
			return err
		}

		if len(req.Sizes) > 0 {
			tx.Where("product_id = ?", product.ID).Delete(&models.Size{})
			for _, sz := range req.Sizes {
				if sz.Value != nil {
					s := models.Size{
						ProductID: &product.ID,
						Name:      &sz.Name,
						Value:     sz.Value,
					}
					tx.Create(&s)
				}
			}
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.DB.Preload("Sizes").First(&product, product.ID)
	c.JSON(http.StatusOK, gin.H{"success": true, "product": product})
}

func (h *SaleV2Handler) DestroyProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := h.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		tx.Where("product_id = ?", product.ID).Delete(&models.ProductImage{})
		tx.Where("product_id = ?", product.ID).Delete(&models.Size{})
		return tx.Delete(&product).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *SaleV2Handler) UpdateOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.Order
	if err := h.DB.Preload("Products").First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if val, ok := req["name"].(string); ok {
			order.Name = &val
		}
		if val, ok := req["phone"].(string); ok {
			order.Phone = &val
		}
		if val, ok := req["address"].(string); ok {
			order.Address = &val
		}
		if val, ok := req["note"].(string); ok {
			order.Note = &val
		}
		if val, ok := req["desc"].(string); ok {
			order.Desc = &val
		}
		if val, ok := req["status"].(string); ok {
			order.Status = &val
		}
		if val, ok := req["is_paid"].(bool); ok {
			order.IsPaid = val
		}
		if val, ok := req["is_freeship"].(bool); ok {
			order.IsFreeship = val
		}
		if val, ok := req["total_amount"].(float64); ok {
			tot := int64(val)
			order.TotalAmount = &tot
		}
		if val, ok := req["deposit_amount"].(float64); ok {
			dep := int64(val)
			order.DepositAmount = &dep
		}

		return tx.Save(&order).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "order": order})
}


func (h *SaleV2Handler) ConsumerDetail(c *gin.Context) {
	id := c.Param("id")
	var consumer models.Consumer
	if err := h.DB.First(&consumer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Consumer not found"})
		return
	}
	c.JSON(http.StatusOK, consumer)
}

func (h *SaleV2Handler) DashboardStats(c *gin.Context) {
	var totalProducts int64
	var availableProducts int64
	var soldProducts int64
	var totalOrders int64
	var totalRevenue int64

	h.DB.Model(&models.Product{}).Count(&totalProducts)
	h.DB.Model(&models.Product{}).Where("status = ?", "AVAILABLE").Count(&availableProducts)
	h.DB.Model(&models.Product{}).Where("status = ?", "SOLD").Count(&soldProducts)
	h.DB.Model(&models.Order{}).Count(&totalOrders)
	h.DB.Model(&models.Order{}).Select("COALESCE(SUM(total_amount), 0)").Scan(&totalRevenue)

	type TypeStat struct {
		Type   string `json:"type"`
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	var typeStats []TypeStat
	h.DB.Model(&models.Product{}).Select("type, status, count(*) as count").Group("type, status").Scan(&typeStats)

	c.JSON(http.StatusOK, gin.H{
		"total_products":     totalProducts,
		"available_products": availableProducts,
		"sold_products":      soldProducts,
		"total_orders":       totalOrders,
		"total_revenue":      totalRevenue,
		"type_stats":         typeStats,
	})
}

func (h *SaleV2Handler) UploadSessions(c *gin.Context) {
	type SessionSummary struct {
		SessionID    string    `json:"session_id"`
		TotalFiles   int64     `json:"total_files"`
		SuccessFiles int64     `json:"success_files"`
		FirstUpload  time.Time `json:"first_upload"`
	}

	var sessions []SessionSummary
	h.DB.Model(&models.ProductUploadLog{}).
		Select("session_id, COUNT(*) as total_files, SUM(CASE WHEN status = 'SUCCESS' THEN 1 ELSE 0 END) as success_files, MIN(created_at) as first_upload").
		Group("session_id").
		Order("first_upload DESC").
		Limit(50).
		Scan(&sessions)

	c.JSON(http.StatusOK, gin.H{"sessions": sessions})
}

func (h *SaleV2Handler) UploadLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset := (page - 1) * limit

	var total int64
	h.DB.Model(&models.ProductUploadLog{}).Count(&total)

	var logs []models.ProductUploadLog
	h.DB.Order("id DESC").Limit(limit).Offset(offset).Find(&logs)

	c.JSON(http.StatusOK, LaravelPagination{
		CurrentPage: page,
		Data:        logs,
		From:        offset + 1,
		LastPage:    int(math.Ceil(float64(total) / float64(limit))),
		PerPage:     limit,
		To:          offset + len(logs),
		Total:       total,
	})
}
