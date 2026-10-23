package handlers

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"sale-app/config"
	"sale-app/models"
	"sale-app/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SmartUploadJob struct {
	SessionID        string
	OriginalFileName string
	MainImagePath    string
	DetailPaths      []string
	ProductType      string
}

type SaleV2Handler struct {
	DB          *gorm.DB
	Cfg         *config.Config
	Storage     *services.StorageService
	Vision      *services.VisionService
	FB          *services.FacebookService
	UploadQueue chan SmartUploadJob
}

func NewSaleV2Handler(db *gorm.DB, cfg *config.Config) *SaleV2Handler {
	fbPageID := ""
	fbPageToken := ""
	mediaBase := ""
	if cfg != nil {
		fbPageID = cfg.FBPageID
		fbPageToken = cfg.FBPageToken
		mediaBase = cfg.MediaBaseURL
	}
	h := &SaleV2Handler{
		DB:          db,
		Cfg:         cfg,
		Storage:     services.NewStorageService(),
		Vision:      services.NewVisionService(),
		FB:          services.NewFacebookService(db, fbPageID, fbPageToken, mediaBase),
		UploadQueue: make(chan SmartUploadJob, 100),
	}
	h.StartWorkers(3)
	return h
}

func (h *SaleV2Handler) StartWorkers(count int) {
	for i := 0; i < count; i++ {
		go func(workerID int) {
			for job := range h.UploadQueue {
				h.processSmartUploadJob(job)
			}
		}(i + 1)
	}
}

func (h *SaleV2Handler) FormatImagePath(path string) string {
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
	if limit > 2000 {
		limit = 2000
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
		// Only filter out empty thumbnails for standard Product Vault grid (limit <= 60 and no search)
		if status == "AVAILABLE" && limit <= 60 && c.Query("search") == "" && c.Query("include_empty_thumb") != "true" {
			query = query.Where("path_thumb IS NOT NULL AND path_thumb != ''")
		}
	}

	// Type filter
	if pType := c.Query("type"); pType != "" {
		query = query.Where("type = ?", pType)
	}

	// Size filters (EXISTS subquery handles ranges and avoid duplicates/join collision)
	applySizeFilter := func(key, fromVal, toVal string) {
		key = strings.TrimSpace(key)
		fromVal = strings.TrimSpace(fromVal)
		toVal = strings.TrimSpace(toVal)
		if key == "" {
			return
		}
		subQuery := "SELECT 1 FROM sizes WHERE sizes.product_id = products.id AND sizes.name = ?"
		args := []interface{}{key}
		if fromVal != "" {
			if fVal, err := strconv.ParseFloat(fromVal, 64); err == nil {
				subQuery += " AND sizes.value >= ?"
				args = append(args, fVal)
			}
		}
		if toVal != "" {
			if tVal, err := strconv.ParseFloat(toVal, 64); err == nil {
				subQuery += " AND sizes.value <= ?"
				args = append(args, tVal)
			}
		}
		query = query.Where("EXISTS ("+subQuery+")", args...)
	}

	applySizeFilter(c.Query("from_search_key"), c.Query("from_search_value"), c.Query("to_search_value"))
	applySizeFilter(c.Query("from_search_key_2"), c.Query("from_search_value_2"), c.Query("to_search_value_2"))
	applySizeFilter(c.Query("from_search_key_3"), c.Query("from_search_value_3"), c.Query("to_search_value_3"))

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
	status := c.Query("status")
	isNewTab := status == "NEW"

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	limitStr := c.Query("limit")
	limit := 10
	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}

	// For NEW tab, show all orders without pagination limit
	noLimit := isNewTab || (limitStr != "" && limit <= 0)

	query := h.DB.Model(&models.Order{})

	if status != "" {
		if status == "NEW" {
			query = query.Where("orders.status != ? AND orders.status != ?", "COMPLETED", "READY FOR SHIPPING").Order("orders.status ASC")
		} else {
			query = query.Where("orders.status = ?", status)
		}
	}

	if search := c.Query("search"); search != "" {
		sUpper := "%" + strings.ToUpper(search) + "%"
		query = query.Joins("LEFT JOIN order_products ON order_products.order_id = orders.id").
			Joins("LEFT JOIN products ON products.id = order_products.product_id").
			Where("UPPER(orders.name) LIKE ? OR orders.phone LIKE ? OR UPPER(orders.address) LIKE ? OR UPPER(products.name) LIKE ?", sUpper, "%"+search+"%", sUpper, sUpper)
	}

	var total int64
	query.Distinct("orders.id").Count(&total)

	var orders []models.Order
	dbQuery := query.Distinct().Select("orders.*").Preload("Products.Images").Preload("Products.Sizes").Order("orders.id DESC")
	offset := 0
	if !noLimit {
		offset = (page - 1) * limit
		dbQuery = dbQuery.Limit(limit).Offset(offset)
	}
	dbQuery.Find(&orders)

	now := time.Now()
	for i := range orders {
		o := &orders[i]
		if o.CreatedAt != nil {
			o.IsOld = now.Sub(*o.CreatedAt).Hours() >= 72 // 3 days
			o.CreatedAtHuman = FormatTimeAgo(o.CreatedAt)
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

		var uniqueProducts []models.Product
		seen := make(map[uint]bool)
		for _, p := range o.Products {
			if !seen[p.ID] {
				seen[p.ID] = true
				uniqueProducts = append(uniqueProducts, p)
			}
		}
		o.Products = uniqueProducts

		for j := range o.Products {
			h.formatProduct(&o.Products[j])
		}
	}

	perPage := limit
	lastPage := 1
	if noLimit {
		perPage = int(total)
		if perPage == 0 {
			perPage = 1
		}
	} else if limit > 0 {
		lastPage = int(math.Ceil(float64(total) / float64(limit)))
		if lastPage == 0 {
			lastPage = 1
		}
	}

	c.JSON(http.StatusOK, LaravelPagination{
		CurrentPage: page,
		Data:        orders,
		From:        offset + 1,
		LastPage:    lastPage,
		PerPage:     perPage,
		To:          offset + len(orders),
		Total:       total,
	})
}

// OrderStats returns statistics for Yesterday, Today, and Current Month (parity with V1)
func (h *SaleV2Handler) OrderStats(c *gin.Context) {
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil || loc == nil {
		loc = time.FixedZone("UTC+7", 7*3600)
	}
	now := time.Now().In(loc)
	todayStr := now.Format("2006-01-02")
	yesterdayStr := now.AddDate(0, 0, -1).Format("2006-01-02")

	type StatResult struct {
		Count int64   `json:"count"`
		Total float64 `json:"revenue"`
	}

	soldByDay := func(dateStr string) StatResult {
		var res StatResult
		h.DB.Table("order_products").
			Select("COUNT(products.id) as count, COALESCE(SUM(products.price), 0) as total").
			Joins("JOIN orders ON orders.id = order_products.order_id").
			Joins("JOIN products ON products.id = order_products.product_id").
			Where("DATE(orders.created_at) = ?", dateStr).
			Scan(&res)
		return res
	}

	todayStat := soldByDay(todayStr)
	yesterdayStat := soldByDay(yesterdayStr)

	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Second)

	var currentTotalOrder int64
	var currentTotalFinalPrice float64
	h.DB.Table("orders").
		Where("created_at BETWEEN ? AND ?", startOfMonth, endOfMonth).
		Count(&currentTotalOrder)
	h.DB.Table("orders").
		Where("created_at BETWEEN ? AND ?", startOfMonth, endOfMonth).
		Select("COALESCE(SUM(final_amount), 0)").
		Scan(&currentTotalFinalPrice)

	c.JSON(http.StatusOK, gin.H{
		"yesterday": gin.H{
			"count":   yesterdayStat.Count,
			"revenue": int64(yesterdayStat.Total),
		},
		"today": gin.H{
			"count":   todayStat.Count,
			"revenue": int64(todayStat.Total),
		},
		"current_month": gin.H{
			"count":   currentTotalOrder,
			"revenue": int64(currentTotalFinalPrice),
		},
	})
}

func (h *SaleV2Handler) OrderDetail(c *gin.Context) {
	id := c.Param("id")
	var order models.Order
	if err := h.DB.Preload("Products.Images").Preload("Products.Sizes").First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	var uniqueProducts []models.Product
	seen := make(map[uint]bool)
	for _, p := range order.Products {
		if !seen[p.ID] {
			seen[p.ID] = true
			uniqueProducts = append(uniqueProducts, p)
		}
	}
	order.Products = uniqueProducts

	for j := range order.Products {
		h.formatProduct(&order.Products[j])
	}

	if order.CreatedAt != nil {
		order.CreatedAtHuman = FormatTimeAgo(order.CreatedAt)
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
		// Auto create/link Consumer
		cleanedPhone := strings.ReplaceAll(req.Phone, " ", "")
		if cleanedPhone != "" {
			var cons models.Consumer
			if err := tx.Where("phone = ?", cleanedPhone).First(&cons).Error; err != nil {
				cons = models.Consumer{
					Name:    &req.Name,
					Phone:   &cleanedPhone,
					Address: &req.Address,
				}
				if err := tx.Create(&cons).Error; err == nil {
					cID := int(cons.ID)
					order.ConsumerID = &cID
				}
			} else {
				cID := int(cons.ID)
				order.ConsumerID = &cID
				updates := map[string]interface{}{}
				if (cons.Address == nil || *cons.Address == "") && req.Address != "" {
					updates["address"] = req.Address
				}
				if (cons.Name == nil || *cons.Name == "") && req.Name != "" {
					updates["name"] = req.Name
				}
				if len(updates) > 0 {
					tx.Model(&cons).Updates(updates)
				}
			}
		}

		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		var cleanProducts []uint
		seenP := make(map[uint]bool)
		for _, pID := range req.Products {
			if pID > 0 && !seenP[pID] {
				seenP[pID] = true
				cleanProducts = append(cleanProducts, pID)
			}
		}

		for _, pID := range cleanProducts {
			op := models.OrderProduct{OrderID: &order.ID, ProductID: &pID}
			if err := tx.Create(&op).Error; err != nil {
				return err
			}
		}

		// Update product status to SOLD
		soldStatus := "SOLD"
		if len(cleanProducts) > 0 {
			if err := tx.Model(&models.Product{}).Where("id IN ?", cleanProducts).Update("status", soldStatus).Error; err != nil {
				return err
			}
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

	h.updateProductStatuses()

	if h.FB != nil {
		for _, pID := range req.Products {
			if pID > 0 {
				h.FB.Dispatch(pID, services.FBActionSold)
			}
		}
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

	var latest models.Product
	h.DB.Where("type = ?", req.Type).Order("id DESC").First(&latest)
	re := regexp.MustCompile(`\d+`)
	num := 1
	if latest.Name != nil {
		if m := re.FindString(*latest.Name); m != "" {
			if n, err := strconv.Atoi(m); err == nil {
				num = n + 1
			}
		}
	}
	codeName := fmt.Sprintf("%s%d", prefix, num)
	status := "AVAILABLE"

	product := models.Product{
		Type:        &req.Type,
		Name:        &codeName,
		Price:       &req.Price,
		Description: &req.Description,
		Status:      &status,
	}

	if file, err := c.FormFile("banner"); err == nil && file != nil {
		if uploadedPath, err := h.Storage.Upload(file); err == nil && uploadedPath != "" {
			product.PathThumb = &uploadedPath
		}
	}

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&product).Error; err != nil {
			return err
		}

		savedSizes := 0
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
				savedSizes++
			}
		}

		// Fallback for multipart form data if req.Sizes wasn't bound
		if savedSizes == 0 {
			for i := 0; i < 30; i++ {
				sName := c.PostForm(fmt.Sprintf("sizes[%d][name]", i))
				sValStr := c.PostForm(fmt.Sprintf("sizes[%d][value]", i))
				if sName != "" && sValStr != "" {
					val, err := strconv.ParseFloat(sValStr, 64)
					if err == nil {
						s := models.Size{
							ProductID: &product.ID,
							Name:      &sName,
							Value:     &val,
						}
						if err := tx.Create(&s).Error; err != nil {
							return err
						}
					}
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
	// Note: DO NOT preload Products here! GORM tx.Save(&order) would duplicate many2many join rows
	if err := h.DB.First(&order, id).Error; err != nil {
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
			cleaned := strings.ReplaceAll(val, " ", "")
			order.Phone = &cleaned
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
		if val, ok := req["additional_amount"].(float64); ok {
			add := int64(val)
			order.AdditionalAmount = &add
		}
		if val, ok := req["freeship_amount"].(float64); ok {
			order.FreeshipAmount = &val
		}
		if val, ok := req["final_amount"].(float64); ok {
			fin := int64(val)
			order.FinalAmount = &fin
		} else if order.TotalAmount != nil {
			order.FinalAmount = order.TotalAmount
		}

		// Save order fields omitting associations to prevent join table duplicates
		if err := tx.Omit("Products").Save(&order).Error; err != nil {
			return err
		}

		// If products are provided in payload, sync associations cleanly
		if prodsRaw, ok := req["products"]; ok {
			var cleanIDs []uint
			seen := make(map[uint]bool)

			if prodList, okList := prodsRaw.([]interface{}); okList {
				for _, item := range prodList {
					var pid uint
					switch v := item.(type) {
					case float64:
						pid = uint(v)
					case int:
						pid = uint(v)
					case int64:
						pid = uint(v)
					}
					if pid > 0 && !seen[pid] {
						seen[pid] = true
						cleanIDs = append(cleanIDs, pid)
					}
				}
			}

			// Clear previous join table rows for this order
			if err := tx.Where("order_id = ?", order.ID).Delete(&models.OrderProduct{}).Error; err != nil {
				return err
			}

			// Insert deduplicated association rows
			for _, pid := range cleanIDs {
				op := models.OrderProduct{
					OrderID:   &order.ID,
					ProductID: &pid,
				}
				if err := tx.Create(&op).Error; err != nil {
					return err
				}
			}
		}

		// Update product prices if provided in products_info
		if prodsInfoRaw, ok := req["products_info"].([]interface{}); ok {
			for _, item := range prodsInfoRaw {
				if infoMap, okMap := item.(map[string]interface{}); okMap {
					if idVal, okID := infoMap["id"]; okID {
						if priceVal, okPrice := infoMap["price"]; okPrice {
							var pid uint
							switch v := idVal.(type) {
							case float64:
								pid = uint(v)
							case int:
								pid = uint(v)
							case int64:
								pid = uint(v)
							}
							var price int64
							switch v := priceVal.(type) {
							case float64:
								price = int64(v)
							case int:
								price = int64(v)
							case int64:
								price = v
							}
							if pid > 0 {
								tx.Model(&models.Product{}).Where("id = ?", pid).Update("price", price)
							}
						}
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update product statuses across DB so unlinked products become AVAILABLE and linked products become SOLD
	h.updateProductStatuses()

	// Preload products to return full fresh order
	h.DB.Preload("Products.Images").Preload("Products.Sizes").First(&order, id)
	var uniqueProducts []models.Product
	seen := make(map[uint]bool)
	for _, p := range order.Products {
		if !seen[p.ID] {
			seen[p.ID] = true
			uniqueProducts = append(uniqueProducts, p)
		}
	}
	order.Products = uniqueProducts

	for j := range order.Products {
		h.formatProduct(&order.Products[j])
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "order": order})
}

func (h *SaleV2Handler) DestroyOrder(c *gin.Context) {
	id := c.Param("id")
	h.DB.Where("id = ?", id).Delete(&models.Order{})
	h.DB.Where("order_id = ?", id).Delete(&models.OrderProduct{})
	h.updateProductStatuses()
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *SaleV2Handler) updateProductStatuses() {
	var oldSoldIDs []uint
	h.DB.Model(&models.Product{}).Where("status = 'SOLD'").Pluck("id", &oldSoldIDs)

	h.DB.Exec("UPDATE products SET status = 'AVAILABLE' WHERE status != 'ON_HOLD'")
	h.DB.Exec("UPDATE products SET status = 'SOLD' WHERE id IN (SELECT DISTINCT product_id FROM order_products)")

	var newSoldIDs []uint
	h.DB.Model(&models.OrderProduct{}).Distinct("product_id").Pluck("product_id", &newSoldIDs)

	if h.FB != nil {
		oldMap := make(map[uint]bool)
		for _, id := range oldSoldIDs {
			oldMap[id] = true
		}
		newMap := make(map[uint]bool)
		for _, id := range newSoldIDs {
			newMap[id] = true
		}

		for _, id := range newSoldIDs {
			if !oldMap[id] {
				h.FB.Dispatch(id, services.FBActionSold)
			}
		}
		for _, id := range oldSoldIDs {
			if !newMap[id] {
				h.FB.Dispatch(id, services.FBActionRevert)
			}
		}
	}
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
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil || loc == nil {
		loc = time.FixedZone("Asia/Ho_Chi_Minh", 7*3600)
	}

	now := time.Now().In(loc)
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	prevMonthDate := startOfMonth.AddDate(0, -1, 0)
	startOfPrevMonth := time.Date(prevMonthDate.Year(), prevMonthDate.Month(), 1, 0, 0, 0, 0, loc)
	endOfPrevMonth := startOfPrevMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	// All orders
	var allOrdersCount int64
	var allOrdersRevenue int64
	h.DB.Model(&models.Order{}).Count(&allOrdersCount)
	h.DB.Model(&models.Order{}).Select("COALESCE(SUM(final_amount), 0)").Scan(&allOrdersRevenue)

	type TypeStat struct {
		Type   string `json:"type"`
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	var allTypeStats []TypeStat
	h.DB.Model(&models.Product{}).Select("type, status, count(*) as count").Group("type, status").Scan(&allTypeStats)

	// New/Hold orders
	var newOrdersCount int64
	var newOrdersRevenue int64
	h.DB.Model(&models.Order{}).Where("status IN ('NEW', 'ON HOLD')").Count(&newOrdersCount)
	h.DB.Model(&models.Order{}).Where("status IN ('NEW', 'ON HOLD')").Select("COALESCE(SUM(final_amount), 0)").Scan(&newOrdersRevenue)

	// Current Month
	var curMonthOrdersCount int64
	var curMonthOrdersRevenue int64
	h.DB.Model(&models.Order{}).Where("created_at BETWEEN ? AND ?", startOfMonth, endOfMonth).Count(&curMonthOrdersCount)
	h.DB.Model(&models.Order{}).Where("created_at BETWEEN ? AND ?", startOfMonth, endOfMonth).Select("COALESCE(SUM(final_amount), 0)").Scan(&curMonthOrdersRevenue)

	type SoldTypeStat struct {
		Type  string `json:"type"`
		Count int64  `json:"count"`
	}
	var currentTypeStats []SoldTypeStat
	h.DB.Table("products").
		Joins("JOIN order_products ON products.id = order_products.product_id").
		Joins("JOIN orders ON orders.id = order_products.order_id").
		Where("orders.created_at BETWEEN ? AND ? AND products.status = 'SOLD'", startOfMonth, endOfMonth).
		Select("products.type, count(*) as count").
		Group("products.type").
		Scan(&currentTypeStats)

	// Prev Month
	var prevMonthOrdersCount int64
	var prevMonthOrdersRevenue int64
	h.DB.Model(&models.Order{}).Where("created_at BETWEEN ? AND ?", startOfPrevMonth, endOfPrevMonth).Count(&prevMonthOrdersCount)
	h.DB.Model(&models.Order{}).Where("created_at BETWEEN ? AND ?", startOfPrevMonth, endOfPrevMonth).Select("COALESCE(SUM(final_amount), 0)").Scan(&prevMonthOrdersRevenue)

	var prevTypeStats []SoldTypeStat
	h.DB.Table("products").
		Joins("JOIN order_products ON products.id = order_products.product_id").
		Joins("JOIN orders ON orders.id = order_products.order_id").
		Where("orders.created_at BETWEEN ? AND ? AND products.status = 'SOLD'", startOfPrevMonth, endOfPrevMonth).
		Select("products.type, count(*) as count").
		Group("products.type").
		Scan(&prevTypeStats)

	// Order counts by status
	type OrderCount struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	var orderCounts []OrderCount
	h.DB.Model(&models.Order{}).Select("status, count(*) as count").Group("status").Scan(&orderCounts)

	// Monthly Sales (last 12 months)
	sub12Months := now.AddDate(0, -12, 0)
	type MonthlySale struct {
		Month       string `json:"month"`
		Number      int64  `json:"number"`
		FinalAmount int64  `json:"final_amount"`
	}
	var monthlySales []MonthlySale
	h.DB.Model(&models.Order{}).
		Select("DATE_FORMAT(created_at, '%Y-%m') as month, count(*) as number, COALESCE(sum(final_amount), 0) as final_amount").
		Where("created_at >= ?", sub12Months).
		Group("month").
		Order("month ASC").
		Scan(&monthlySales)

	c.JSON(http.StatusOK, gin.H{
		"new_orders_count":     newOrdersCount,
		"new_orders_revenue":   newOrdersRevenue,
		"month_orders_count":   curMonthOrdersCount,
		"month_orders_revenue": curMonthOrdersRevenue,
		"all": gin.H{
			"total_orders":  allOrdersCount,
			"total_revenue": allOrdersRevenue,
			"type_stats":    allTypeStats,
		},
		"current_month": gin.H{
			"total_orders":  curMonthOrdersCount,
			"total_revenue": curMonthOrdersRevenue,
			"type_stats":    currentTypeStats,
		},
		"prev_month": gin.H{
			"total_orders":  prevMonthOrdersCount,
			"total_revenue": prevMonthOrdersRevenue,
			"type_stats":    prevTypeStats,
		},
		"order_counts":  orderCounts,
		"monthly_sales": monthlySales,

		// Flat keys for backwards compatibility:
		"total_orders":  allOrdersCount,
		"total_revenue": allOrdersRevenue,
		"type_stats":    allTypeStats,
	})
}

func (h *SaleV2Handler) UploadSessions(c *gin.Context) {
	type RawSession struct {
		SessionID string    `json:"session_id"`
		CreatedAt time.Time `json:"created_at"`
	}

	var rawSessions []RawSession
	h.DB.Model(&models.ProductUploadLog{}).
		Select("session_id, MIN(created_at) as created_at").
		Group("session_id").
		Order("created_at DESC").
		Scan(&rawSessions)

	if len(rawSessions) == 0 {
		c.JSON(http.StatusOK, gin.H{"sessions": []interface{}{}})
		return
	}

	var sessionIDs []string
	for _, s := range rawSessions {
		sessionIDs = append(sessionIDs, s.SessionID)
	}

	type ProductStat struct {
		SessionID string `json:"session_id"`
		Type      string `json:"type"`
		Status    string `json:"status"`
		Count     int    `json:"count"`
	}

	var stats []ProductStat
	h.DB.Table("product_upload_logs").
		Joins("JOIN products ON products.id = product_upload_logs.product_id").
		Where("product_upload_logs.session_id IN ?", sessionIDs).
		Select("product_upload_logs.session_id, products.type, products.status, count(*) as count").
		Group("product_upload_logs.session_id, products.type, products.status").
		Scan(&stats)

	statsMap := make(map[string][]ProductStat)
	for _, st := range stats {
		statsMap[st.SessionID] = append(statsMap[st.SessionID], st)
	}

	type SessionOut struct {
		SessionID     string         `json:"session_id"`
		CreatedAt     string         `json:"created_at"`
		TotalProducts int            `json:"total_products"`
		SoldProducts  int            `json:"sold_products"`
		TypeBreakdown map[string]int `json:"type_breakdown"`
	}

	var sessions []SessionOut
	for _, s := range rawSessions {
		total := 0
		sold := 0
		breakdown := make(map[string]int)

		for _, st := range statsMap[s.SessionID] {
			total += st.Count
			if st.Status == "SOLD" {
				sold += st.Count
			}
			if st.Type != "" {
				breakdown[st.Type] += st.Count
			}
		}

		sessions = append(sessions, SessionOut{
			SessionID:     s.SessionID,
			CreatedAt:     s.CreatedAt.Format("2006-01-02 15:04:05"),
			TotalProducts: total,
			SoldProducts:  sold,
			TypeBreakdown: breakdown,
		})
	}

	c.JSON(http.StatusOK, gin.H{"sessions": sessions})
}

func (h *SaleV2Handler) UploadSessionProducts(c *gin.Context) {
	sessionID := c.Param("id")
	if sessionID == "" {
		sessionID = c.Param("sessionId")
	}

	var logs []models.ProductUploadLog
	h.DB.Table("product_upload_logs").
		Where("product_upload_logs.session_id = ?", sessionID).
		Joins("LEFT JOIN products ON products.id = product_upload_logs.product_id").
		Preload("Product.Images").
		Preload("Product.Sizes").
		Order("CASE WHEN products.status = 'SOLD' THEN 0 ELSE 1 END, product_upload_logs.id ASC").
		Find(&logs)

	for i := range logs {
		if logs[i].Product != nil {
			h.formatProduct(logs[i].Product)
		}
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs})
}

func (h *SaleV2Handler) ProductsByLatestUpload(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "60"))
	if limit < 1 {
		limit = 60
	}
	offset := (page - 1) * limit

	subQuery := "(SELECT product_id, MAX(created_at) as latest_upload FROM product_upload_logs WHERE product_id IS NOT NULL GROUP BY product_id) as upload_summary"

	query := h.DB.Model(&models.Product{}).
		Joins("LEFT JOIN " + subQuery + " ON upload_summary.product_id = products.id").
		Where("products.path_thumb IS NOT NULL AND products.path_thumb != ''")

	if status := c.Query("status"); status != "" {
		query = query.Where("products.status = ?", status)
	}
	if pType := c.Query("type"); pType != "" {
		query = query.Where("products.type = ?", pType)
	}

	var total int64
	query.Count(&total)

	var products []models.Product
	query.Select("products.*, upload_summary.latest_upload").
		Preload("Images").
		Preload("Sizes").
		Order("COALESCE(upload_summary.latest_upload, products.created_at) DESC").
		Limit(limit).
		Offset(offset).
		Find(&products)

	for i := range products {
		h.formatProduct(&products[i])
	}

	c.JSON(http.StatusOK, LaravelPagination{
		CurrentPage: page,
		Data:        products,
		From:        offset + 1,
		LastPage:    int(math.Ceil(float64(total) / float64(limit))),
		PerPage:     limit,
		To:          offset + len(products),
		Total:       total,
	})
}

func (h *SaleV2Handler) SoldByDate(c *gin.Context) {
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil || loc == nil {
		loc = time.FixedZone("Asia/Ho_Chi_Minh", 7*3600)
	}

	type ActualRow struct {
		Date         string `json:"date"`
		Total        int    `json:"total"`
		TotalRevenue int64  `json:"total_revenue"`
		Types        string `json:"types"`
	}

	var actualRows []ActualRow
	h.DB.Table("order_products").
		Joins("JOIN products ON products.id = order_products.product_id").
		Joins("JOIN orders ON orders.id = order_products.order_id").
		Select("DATE_FORMAT(orders.created_at, '%Y-%m-%d') as date, COUNT(*) as total, COALESCE(SUM(products.price), 0) as total_revenue, GROUP_CONCAT(products.type ORDER BY products.type SEPARATOR ',') as types").
		Group("DATE_FORMAT(orders.created_at, '%Y-%m-%d')").
		Scan(&actualRows)

	if len(actualRows) == 0 {
		c.JSON(http.StatusOK, gin.H{"dates": []interface{}{}})
		return
	}

	actualMap := make(map[string]ActualRow)
	var earliestStr string
	for _, r := range actualRows {
		actualMap[r.Date] = r
		if earliestStr == "" || r.Date < earliestStr {
			earliestStr = r.Date
		}
	}

	earliest, err := time.ParseInLocation("2006-01-02", earliestStr, loc)
	if err != nil {
		earliest = time.Now().In(loc).AddDate(0, -1, 0)
	}

	today := time.Now().In(loc)
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, loc)

	type DateOut struct {
		Date          string         `json:"date"`
		Total         int            `json:"total"`
		TotalRevenue  int64          `json:"total_revenue"`
		TypeBreakdown map[string]int `json:"type_breakdown"`
		IsEmpty       bool           `json:"is_empty"`
	}

	var allDates []DateOut
	for d := today; !d.Before(earliest); d = d.AddDate(0, 0, -1) {
		dateStr := d.Format("2006-01-02")
		if row, ok := actualMap[dateStr]; ok {
			breakdown := make(map[string]int)
			if row.Types != "" {
				for _, t := range strings.Split(row.Types, ",") {
					t = strings.TrimSpace(t)
					if t != "" {
						breakdown[t]++
					}
				}
			}
			allDates = append(allDates, DateOut{
				Date:          dateStr,
				Total:         row.Total,
				TotalRevenue:  row.TotalRevenue,
				TypeBreakdown: breakdown,
				IsEmpty:       false,
			})
		} else {
			allDates = append(allDates, DateOut{
				Date:          dateStr,
				Total:         0,
				TotalRevenue:  0,
				TypeBreakdown: make(map[string]int),
				IsEmpty:       true,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{"dates": allDates})
}

func (h *SaleV2Handler) SoldByDateProducts(c *gin.Context) {
	date := c.Param("date")
	var products []models.Product

	h.DB.Table("products").
		Joins("JOIN order_products ON order_products.product_id = products.id").
		Joins("JOIN orders ON orders.id = order_products.order_id").
		Where("DATE_FORMAT(orders.created_at, '%Y-%m-%d') = ?", date).
		Select("products.*").
		Preload("Images").
		Preload("Sizes").
		Order("orders.created_at DESC").
		Find(&products)

	for i := range products {
		h.formatProduct(&products[i])
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

func (h *SaleV2Handler) UploadLogs(c *gin.Context) {
	sessionIdsParam := c.Query("session_ids")
	if sessionIdsParam != "" {
		sIDs := strings.Split(sessionIdsParam, ",")
		var logs []models.ProductUploadLog
		h.DB.Where("session_id IN ?", sIDs).Find(&logs)
		c.JSON(http.StatusOK, logs)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset := (page - 1) * limit

	var total int64
	h.DB.Model(&models.ProductUploadLog{}).Count(&total)

	var logs []models.ProductUploadLog
	h.DB.Preload("Product.Images").Preload("Product.Sizes").Order("id DESC").Limit(limit).Offset(offset).Find(&logs)

	for i := range logs {
		if logs[i].Product != nil {
			h.formatProduct(logs[i].Product)
		}
	}

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

// UploadSmartGroup accepts a main image and optional detail images and dispatches a background job
func (h *SaleV2Handler) UploadSmartGroup(c *gin.Context) {
	sessionID := c.DefaultPostForm("session_id", fmt.Sprintf("SMART_%d", time.Now().UnixNano()))
	pType := c.DefaultPostForm("type", "TROUSERS")

	// 1. Store Main Image
	mainFile, err := c.FormFile("main")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Main image is required"})
		return
	}

	mainStoredPath, err := h.Storage.Upload(mainFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": fmt.Sprintf("Failed to store main image: %v", err)})
		return
	}

	// 2. Store Detail Images
	var detailPaths []string
	form, err := c.MultipartForm()
	if err == nil && form != nil {
		for key, files := range form.File {
			if strings.HasPrefix(key, "details") {
				for _, f := range files {
					if path, err := h.Storage.Upload(f); err == nil && path != "" {
						detailPaths = append(detailPaths, path)
					}
				}
			}
		}
	}

	// 3. Dispatch Job into in-process Queue Channel
	job := SmartUploadJob{
		SessionID:        sessionID,
		OriginalFileName: mainFile.Filename,
		MainImagePath:    mainStoredPath,
		DetailPaths:      detailPaths,
		ProductType:      pType,
	}

	select {
	case h.UploadQueue <- job:
	default:
		go func() { h.UploadQueue <- job }()
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"session_id": sessionID,
		"status":     "DISPATCHED",
	})
}

// processSmartUploadJob processes OCR recognition and product image attachment in background
func (h *SaleV2Handler) processSmartUploadJob(job SmartUploadJob) {
	// 1. Init Log (Step 1 start)
	logRecord := models.ProductUploadLog{
		SessionID: job.SessionID,
		FileName:  job.OriginalFileName,
		FilePath:  job.MainImagePath,
		Status:    "PROCESSING",
		Message:   strPtr("Step 1: Recognition & Product Search..."),
	}
	h.DB.Create(&logRecord)

	// 2. STEP 1: RECOGNITION & FIND PRODUCT
	detectedText, err := h.Vision.ExtractTextFromImage(job.MainImagePath)
	if err != nil {
		errMsg := fmt.Sprintf("Không thể đọc được chữ từ ảnh (Google Vision Error): %v", err)
		h.DB.Model(&models.ProductUploadLog{}).
			Where("session_id = ? AND file_path = ?", job.SessionID, job.MainImagePath).
			Updates(map[string]interface{}{
				"status":  "ERROR",
				"message": errMsg,
			})
		return
	}

	// A. Extract Code based on category prefix
	prefix := job.ProductType
	if p, ok := ProductKeyTypes[job.ProductType]; ok {
		prefix = p
	}

	reCode := regexp.MustCompile(fmt.Sprintf(`(?i)%s\s*(\d+)`, regexp.QuoteMeta(prefix)))
	matches := reCode.FindStringSubmatch(detectedText)
	if len(matches) < 2 {
		errMsg := fmt.Sprintf("Không detect được mã sản phẩm (%s...) từ ảnh hoặc tên file.", prefix)
		h.DB.Model(&models.ProductUploadLog{}).
			Where("session_id = ? AND file_path = ?", job.SessionID, job.MainImagePath).
			Updates(map[string]interface{}{
				"status":        "ERROR",
				"detected_text": detectedText,
				"message":       errMsg,
			})
		return
	}
	candidateName := strings.ToUpper(prefix + matches[1])

	// B. Extract Price
	var price int64
	rePrice := regexp.MustCompile(`(?i)(\d+)\s*k\b`)
	pMatches := rePrice.FindStringSubmatch(detectedText)
	if len(pMatches) >= 2 {
		if val, err := strconv.ParseInt(pMatches[1], 10, 64); err == nil {
			price = val * 1000
		}
	}

	// Fallback to default category price if price is 0
	if price == 0 {
		categoryMap := map[string]int64{
			"TROUSERS": 180000,
			"BLAZER":   280000,
			"GILE":     120000,
		}
		if defPrice, ok := categoryMap[job.ProductType]; ok {
			price = defPrice
		}
	}

	// C. Search Existing Product
	var existing models.Product
	if err := h.DB.Where("type = ? AND name = ?", job.ProductType, candidateName).First(&existing).Error; err != nil {
		errMsg := fmt.Sprintf("Sản phẩm mã (%s) không tồn tại. Vui lòng tạo sp trước.", candidateName)
		h.DB.Model(&models.ProductUploadLog{}).
			Where("session_id = ? AND file_path = ?", job.SessionID, job.MainImagePath).
			Updates(map[string]interface{}{
				"status":        "ERROR",
				"product_code":  candidateName,
				"detected_text": detectedText,
				"message":       errMsg,
			})
		return
	}

	// Update product record: Price & thumbnail
	if price > 0 {
		existing.Price = &price
	}
	if existing.PathThumb == nil || *existing.PathThumb == "" {
		existing.PathThumb = &job.MainImagePath
	}
	h.DB.Save(&existing)

	// Step 2: Attaching images
	prodIDInt := int(existing.ID)
	priceFloat := float64(price)
	h.DB.Model(&models.ProductUploadLog{}).
		Where("session_id = ? AND file_path = ?", job.SessionID, job.MainImagePath).
		Updates(map[string]interface{}{
			"product_code":  candidateName,
			"product_id":    prodIDInt,
			"amount":        priceFloat,
			"detected_text": detectedText,
			"message":       "Step 2: Attaching images...",
		})

	// Attach Main Image as a detail
	h.DB.Create(&models.ProductImage{
		ProductID: existing.ID,
		FilePath:  job.MainImagePath,
		IsMain:    true,
	})

	// Attach Detail Images
	for _, dPath := range job.DetailPaths {
		if dPath != "" {
			h.DB.Create(&models.ProductImage{
				ProductID: existing.ID,
				FilePath:  dPath,
				IsMain:    false,
			})
		}
	}

	// Log Success
	h.DB.Model(&models.ProductUploadLog{}).
		Where("session_id = ? AND file_path = ?", job.SessionID, job.MainImagePath).
		Updates(map[string]interface{}{
			"status":  "SUCCESS",
			"message": "Xử lý thành công (2 bước)",
		})

	// Sync to Facebook page
	if h.FB != nil {
		h.FB.Dispatch(existing.ID, services.FBActionCreate)
	}
}

// CleanupImages scans storage for orphaned images and removes them safely
func (h *SaleV2Handler) CleanupImages(c *gin.Context) {
	// 1. Query whitelist of in-use images:
	// - Products with status AVAILABLE
	// - Products with status SOLD where order is not COMPLETED or created within the last 60 days
	var activeThumbs []string
	h.DB.Raw(`
		SELECT DISTINCT p.path_thumb 
		FROM products p 
		LEFT JOIN order_products op ON p.id = op.product_id 
		LEFT JOIN orders o ON op.order_id = o.id 
		WHERE p.path_thumb IS NOT NULL AND p.path_thumb != ''
		  AND (p.status = 'AVAILABLE' OR (p.status = 'SOLD' AND (o.status != 'COMPLETED' OR o.created_at >= NOW() - INTERVAL 60 DAY)))
	`).Scan(&activeThumbs)

	var activeScaleThumbs []string
	h.DB.Raw(`
		SELECT DISTINCT p.image_thumb_scale 
		FROM products p 
		LEFT JOIN order_products op ON p.id = op.product_id 
		LEFT JOIN orders o ON op.order_id = o.id 
		WHERE p.image_thumb_scale IS NOT NULL AND p.image_thumb_scale != ''
		  AND (p.status = 'AVAILABLE' OR (p.status = 'SOLD' AND (o.status != 'COMPLETED' OR o.created_at >= NOW() - INTERVAL 60 DAY)))
	`).Scan(&activeScaleThumbs)

	var activeDetailImages []string
	h.DB.Raw(`
		SELECT DISTINCT pi.file_path 
		FROM product_images pi 
		JOIN products p ON pi.product_id = p.id 
		LEFT JOIN order_products op ON p.id = op.product_id 
		LEFT JOIN orders o ON op.order_id = o.id 
		WHERE pi.file_path IS NOT NULL AND pi.file_path != ''
		  AND (p.status = 'AVAILABLE' OR (p.status = 'SOLD' AND (o.status != 'COMPLETED' OR o.created_at >= NOW() - INTERVAL 60 DAY)))
	`).Scan(&activeDetailImages)

	usedFiles := make(map[string]bool)
	cleanAndAdd := func(p string) {
		p = strings.TrimSpace(p)
		if p == "" {
			return
		}
		base := filepath.Base(p)
		if base != "" && base != "." {
			usedFiles[base] = true
		}
	}

	for _, p := range activeThumbs {
		cleanAndAdd(p)
	}
	for _, p := range activeScaleThumbs {
		cleanAndAdd(p)
	}
	for _, p := range activeDetailImages {
		cleanAndAdd(p)
	}

	// 2. Scan local storage directory ./storage/sale
	storageDir := "./storage/sale"
	entries, err := os.ReadDir(storageDir)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success":        false,
			"error":          fmt.Sprintf("Directory %s not found or unreadable: %v", storageDir, err),
			"used_images_db": len(usedFiles),
		})
		return
	}

	deletedCount := 0
	skippedRecent := 0
	var reclaimedBytes int64 = 0
	safeThreshold := 24 * time.Hour

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		filename := entry.Name()
		filePath := filepath.Join(storageDir, filename)

		info, err := entry.Info()
		if err != nil {
			continue
		}

		// Skip files modified in last 24h for safety
		if time.Since(info.ModTime()) < safeThreshold {
			skippedRecent++
			continue
		}

		if !usedFiles[filename] {
			size := info.Size()
			if err := os.Remove(filePath); err == nil {
				deletedCount++
				reclaimedBytes += size
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success":         true,
		"scanned_files":   len(entries),
		"deleted_files":   deletedCount,
		"skipped_recent":  skippedRecent,
		"used_images_db":  len(usedFiles),
		"reclaimed_mb":    fmt.Sprintf("%.2f MB", float64(reclaimedBytes)/(1024*1024)),
		"reclaimed_bytes": reclaimedBytes,
		"message":         "Image cleanup executed successfully",
	})
}


