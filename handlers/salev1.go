package handlers

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"sale-app/config"
	"sale-app/models"
	"sale-app/services"

	"github.com/flosch/pongo2/v6"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	ProductTypes = []string{
		"TROUSERS",
		"JACKET",
		"SHIRT",
		"BLAZER",
		"TIE",
		"GILE",
		"BELT",
		"POLO SHIRT",
		"HAT",
		"SUIT",
	}

	ProductKeyTypes = map[string]string{
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

	SizeList = []string{
		"VAI",
		"NGỰC",
		"EO",
		"DÀI ÁO",
		"DÀI ÁO SAU",
		"DÀI TAY",
		"EO QUẦN",
		"ĐÁY",
		"ĐÙI",
		"DÀI QUẦN",
		"ỐNG",
		"DƯ LAI",
	}

	MappingSizes = map[string]string{
		"VAI":        "V",
		"NGỰC":       "N",
		"EO":         "E",
		"DÀI ÁO":     "D",
		"DÀI ÁO SAU": "DS",
		"DÀI TAY":    "DT",
		"EO QUẦN":    "E",
		"ĐÁY":        "Đ",
		"ĐÙI":        "ĐU",
		"DÀI QUẦN":   "DQ",
		"ỐNG":        "O",
	}
)

func strPtr(s string) *string {
	return &s
}

func strVal(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func int64Ptr(n int64) *int64 {
	return &n
}

func int64Val(p *int64) int64 {
	if p == nil {
		return 0
	}
	return *p
}

func float64Ptr(f float64) *float64 {
	return &f
}

func float64Val(p *float64) float64 {
	if p == nil {
		return 0
	}
	return *p
}

func init() {
	// Register pongo2 filter: number_format
	_ = pongo2.RegisterFilter("number_format", func(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
		val := in.Integer()
		if in.IsFloat() {
			val = int(in.Float())
		}
		formatted := formatNumberWithCommas(int64(val))
		return pongo2.AsValue(formatted), nil
	})

	// Register pongo2 filter: img_url
	_ = pongo2.RegisterFilter("img_url", func(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
		path := in.String()
		if path == "" {
			return pongo2.AsValue(""), nil
		}
		if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
			return pongo2.AsValue(path), nil
		}
		clean := strings.TrimLeft(path, "/")
		baseURL := os.Getenv("MEDIA_BASE_URL")
		if baseURL == "" {
			baseURL = "https://www.tadimedia.vn/storage/"
		}
		if strings.HasPrefix(clean, "storage/") {
			clean = strings.TrimPrefix(clean, "storage/")
		}
		return pongo2.AsValue(baseURL + clean), nil
	})
}

func formatNumberWithCommas(n int64) string {
	in := strconv.FormatInt(n, 10)
	var out []byte
	isNegative := false
	if strings.HasPrefix(in, "-") {
		isNegative = true
		in = in[1:]
	}
	for i, c := range in {
		if i > 0 && (len(in)-i)%3 == 0 {
			out = append(out, ',')
		}
		out = append(out, byte(c))
	}
	if isNegative {
		return "-" + string(out)
	}
	return string(out)
}

type SaleV1Handler struct {
	DB          *gorm.DB
	Cfg         *config.Config
	Storage     *services.StorageService
	templateSet *pongo2.TemplateSet
}

func NewSaleV1Handler(db *gorm.DB, cfg *config.Config) *SaleV1Handler {
	loader := pongo2.MustNewLocalFileSystemLoader("./templates/sale")
	ts := pongo2.NewSet("sale_templates", loader)

	return &SaleV1Handler{
		DB:          db,
		Cfg:         cfg,
		Storage:     services.NewStorageService(),
		templateSet: ts,
	}
}

func (h *SaleV1Handler) render(c *gin.Context, templateName string, ctx pongo2.Context) {
	tpl, err := h.templateSet.FromFile(templateName)
	if err != nil {
		c.String(http.StatusInternalServerError, "Template error: %v", err)
		return
	}
	out, err := tpl.Execute(ctx)
	if err != nil {
		c.String(http.StatusInternalServerError, "Render error: %v", err)
		return
	}
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, out)
}

type PaginationData struct {
	CurrentPage int
	PrevPage    int
	NextPage    int
	LastPage    int
	PerPage     int
	Total       int64
	From        int64
	To          int64
}

// 1. Home / Dashboard
func (h *SaleV1Handler) Home(c *gin.Context) {
	var totalOrder int64
	var totalFinalPrice float64
	h.DB.Model(&models.Order{}).Count(&totalOrder)
	h.DB.Model(&models.Order{}).Select("COALESCE(SUM(final_amount), 0)").Scan(&totalFinalPrice)

	type TypeStat struct {
		Type   string
		Status string
		Count  int64
	}
	var typeStats []TypeStat
	h.DB.Table("products").
		Select("type, status, count(id) as count").
		Where("path_thumb IS NOT NULL AND path_thumb != ''").
		Group("type, status").
		Scan(&typeStats)

	soldMap := make(map[string]int64)
	availMap := make(map[string]int64)
	var totalSold, totalAvail int64
	for _, s := range typeStats {
		if s.Status == "SOLD" {
			soldMap[s.Type] = s.Count
			totalSold += s.Count
		} else if s.Status == "AVAILABLE" {
			availMap[s.Type] = s.Count
			totalAvail += s.Count
		}
	}

	type ProductRow struct {
		Type      string
		Sold      int64
		Available int64
	}
	var typeProductRows []ProductRow
	for _, t := range ProductTypes {
		typeProductRows = append(typeProductRows, ProductRow{
			Type:      t,
			Sold:      soldMap[t],
			Available: availMap[t],
		})
	}

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	var currentTotalOrder int64
	var currentTotalFinalPrice float64
	h.DB.Model(&models.Order{}).Where("created_at BETWEEN ? AND ?", startOfMonth, endOfMonth).Count(&currentTotalOrder)
	h.DB.Model(&models.Order{}).Where("created_at BETWEEN ? AND ?", startOfMonth, endOfMonth).Select("COALESCE(SUM(final_amount), 0)").Scan(&currentTotalFinalPrice)

	prevMonth := now.AddDate(0, -1, 0)
	startOfPrevMonth := time.Date(prevMonth.Year(), prevMonth.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfPrevMonth := startOfPrevMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	var prevTotalOrder int64
	var prevTotalFinalPrice float64
	h.DB.Model(&models.Order{}).Where("created_at BETWEEN ? AND ?", startOfPrevMonth, endOfPrevMonth).Count(&prevTotalOrder)
	h.DB.Model(&models.Order{}).Where("created_at BETWEEN ? AND ?", startOfPrevMonth, endOfPrevMonth).Select("COALESCE(SUM(final_amount), 0)").Scan(&prevTotalFinalPrice)

	type MonthlyOrder struct {
		Month       string
		Number      int64
		FinalAmount float64
	}
	var orderByMonth []MonthlyOrder
	h.DB.Table("orders").
		Select("DATE_FORMAT(created_at, '%Y-%m') as month, count(id) as number, COALESCE(SUM(final_amount), 0) as final_amount").
		Group("month").
		Order("month DESC").
		Limit(12).
		Scan(&orderByMonth)

	h.render(c, "index.html", pongo2.Context{
		"current_route":             "home",
		"total_order":               totalOrder,
		"total_final_price":         int64(totalFinalPrice),
		"type_product_rows":         typeProductRows,
		"total_sold":                totalSold,
		"total_available":           totalAvail,
		"previous_total_order":      prevTotalOrder,
		"previous_total_final_price": int64(prevTotalFinalPrice),
		"prev_product_rows":         typeProductRows,
		"current_total_order":       currentTotalOrder,
		"current_total_final_price": int64(currentTotalFinalPrice),
		"current_product_rows":      typeProductRows,
		"order_by_month":            orderByMonth,
	})
}

// 2. Products table
func (h *SaleV1Handler) Products(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	perPage := 40
	if c.Query("showAll") == "true" {
		perPage = 500
	}

	search := c.Query("search")
	selectedType := c.Query("type")
	selectedStatus := c.Query("status")
	fromSearchKey := c.Query("from_search_key")
	fromSearchVal := c.Query("from_search_value")
	toSearchVal := c.Query("to_search_value")
	soldFrom := c.Query("sold_from")
	soldTo := c.Query("sold_to")

	query := h.DB.Model(&models.Product{}).Preload("Sizes").Order("status ASC, id DESC")

	if search != "" {
		names := strings.Split(search, ",")
		query = query.Where(func(db *gorm.DB) *gorm.DB {
			sub := db
			for i, name := range names {
				clean := strings.TrimSpace(name)
				if clean != "" {
					if i == 0 {
						sub = sub.Where("name LIKE ?", "%"+clean+"%")
					} else {
						sub = sub.Or("name LIKE ?", "%"+clean+"%")
					}
				}
			}
			return sub
		}(h.DB))
	}

	if selectedType != "" {
		query = query.Where("type = ?", selectedType)
	}
	if selectedStatus != "" {
		query = query.Where("status = ?", selectedStatus)
	}

	if fromSearchKey != "" && (fromSearchVal != "" || toSearchVal != "") {
		query = query.Joins("JOIN sizes ON sizes.product_id = products.id").Where("sizes.name = ?", fromSearchKey)
		if fromSearchVal != "" {
			fVal, _ := strconv.ParseFloat(fromSearchVal, 64)
			query = query.Where("sizes.value >= ?", fVal)
		}
		if toSearchVal != "" {
			tVal, _ := strconv.ParseFloat(toSearchVal, 64)
			query = query.Where("sizes.value <= ?", tVal)
		}
	}

	if soldFrom != "" && soldTo != "" {
		query = query.Joins("JOIN order_products ON products.id = order_products.product_id").
			Joins("JOIN orders ON orders.id = order_products.order_id").
			Where("orders.created_at BETWEEN ? AND ?", soldFrom+" 00:00:00", soldTo+" 23:59:59")
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * perPage
	var products []models.Product
	query.Offset(offset).Limit(perPage).Find(&products)

	type ProductView struct {
		ID              uint
		Type            string
		Name            string
		Price           int64
		PathThumb       string
		Description     string
		Status          string
		SizeFormatted   string
		PriceFormattedK string
	}
	var productViews []ProductView
	for _, p := range products {
		pPrice := int64Val(p.Price)
		sizeStr := formatProductSizes(&p)
		priceK := fmt.Sprintf("%gk", float64(pPrice)/1000)
		productViews = append(productViews, ProductView{
			ID:              p.ID,
			Type:            strVal(p.Type),
			Name:            strVal(p.Name),
			Price:           pPrice,
			PathThumb:       strVal(p.PathThumb),
			Description:     strVal(p.Description),
			Status:          strVal(p.Status),
			SizeFormatted:   sizeStr,
			PriceFormattedK: priceK,
		})
	}

	lastPage := int(math.Ceil(float64(total) / float64(perPage)))
	if lastPage < 1 {
		lastPage = 1
	}
	from := int64(offset + 1)
	if from > total {
		from = total
	}
	to := int64(offset + len(products))

	h.render(c, "products.html", pongo2.Context{
		"current_route":     "products",
		"products":          productViews,
		"product_types":     ProductTypes,
		"size_list":         SizeList,
		"search":            search,
		"selected_type":     selectedType,
		"selected_status":   selectedStatus,
		"from_search_key":   fromSearchKey,
		"from_search_value": fromSearchVal,
		"to_search_value":   toSearchVal,
		"sold_from":         soldFrom,
		"sold_to":           soldTo,
		"pagination": PaginationData{
			CurrentPage: page,
			PrevPage:    page - 1,
			NextPage:    page + 1,
			LastPage:    lastPage,
			PerPage:     perPage,
			Total:       total,
			From:        from,
			To:          to,
		},
	})
}

// 3. Available Products grid
func (h *SaleV1Handler) AvailableProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	perPage := 60
	if c.Query("showAll") == "true" {
		perPage = 500
	}

	search := c.Query("search")
	selectedType := c.Query("type")
	fromSearchKey := c.Query("from_search_key")
	fromSearchVal := c.Query("from_search_value")
	toSearchVal := c.Query("to_search_value")
	fromSearchKey2 := c.Query("from_search_key_2")
	fromSearchVal2 := c.Query("from_search_value_2")
	toSearchVal2 := c.Query("to_search_value_2")

	query := h.DB.Model(&models.Product{}).Preload("Sizes").
		Where("status = 'AVAILABLE' AND path_thumb IS NOT NULL AND path_thumb != ''").
		Order("id DESC")

	if search != "" {
		names := strings.Split(search, ",")
		query = query.Where(func(db *gorm.DB) *gorm.DB {
			sub := db
			for i, name := range names {
				clean := strings.TrimSpace(name)
				if clean != "" {
					if i == 0 {
						sub = sub.Where("name LIKE ?", "%"+clean+"%")
					} else {
						sub = sub.Or("name LIKE ?", "%"+clean+"%")
					}
				}
			}
			return sub
		}(h.DB))
	}

	if selectedType != "" {
		query = query.Where("type = ?", selectedType)
	}

	if fromSearchKey != "" && (fromSearchVal != "" || toSearchVal != "") {
		query = query.Joins("JOIN sizes ON sizes.product_id = products.id").Where("sizes.name = ?", fromSearchKey)
		if fromSearchVal != "" {
			fVal, _ := strconv.ParseFloat(fromSearchVal, 64)
			query = query.Where("sizes.value >= ?", fVal)
		}
		if toSearchVal != "" {
			tVal, _ := strconv.ParseFloat(toSearchVal, 64)
			query = query.Where("sizes.value <= ?", tVal)
		}
	}

	if fromSearchKey2 != "" && (fromSearchVal2 != "" || toSearchVal2 != "") {
		query = query.Joins("JOIN sizes sizes_2 ON sizes_2.product_id = products.id").Where("sizes_2.name = ?", fromSearchKey2)
		if fromSearchVal2 != "" {
			fVal, _ := strconv.ParseFloat(fromSearchVal2, 64)
			query = query.Where("sizes_2.value >= ?", fVal)
		}
		if toSearchVal2 != "" {
			tVal, _ := strconv.ParseFloat(toSearchVal2, 64)
			query = query.Where("sizes_2.value <= ?", tVal)
		}
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * perPage
	var products []models.Product
	query.Offset(offset).Limit(perPage).Find(&products)

	type ProductView struct {
		ID            uint
		Name          string
		Price         int64
		PathThumb     string
		SizeFormatted string
	}
	var productViews []ProductView
	for _, p := range products {
		sizeStr := formatProductSizes(&p)
		productViews = append(productViews, ProductView{
			ID:            p.ID,
			Name:          strVal(p.Name),
			Price:         int64Val(p.Price),
			PathThumb:     strVal(p.PathThumb),
			SizeFormatted: sizeStr,
		})
	}

	lastPage := int(math.Ceil(float64(total) / float64(perPage)))
	if lastPage < 1 {
		lastPage = 1
	}
	from := int64(offset + 1)
	if from > total {
		from = total
	}
	to := int64(offset + len(products))

	h.render(c, "available-products.html", pongo2.Context{
		"current_route":       "available_products",
		"products":            productViews,
		"product_types":       ProductTypes,
		"size_list":           SizeList,
		"search":              search,
		"selected_type":       selectedType,
		"from_search_key":     fromSearchKey,
		"from_search_value":   fromSearchVal,
		"to_search_value":     toSearchVal,
		"from_search_key_2":   fromSearchKey2,
		"from_search_value_2": fromSearchVal2,
		"to_search_value_2":   toSearchVal2,
		"pagination": PaginationData{
			CurrentPage: page,
			PrevPage:    page - 1,
			NextPage:    page + 1,
			LastPage:    lastPage,
			PerPage:     perPage,
			Total:       total,
			From:        from,
			To:          to,
		},
	})
}

// 4. Product Detail & Update
func (h *SaleV1Handler) ProductDetail(c *gin.Context) {
	id := c.Param("id")

	if c.Request.Method == "POST" {
		var p models.Product
		if err := h.DB.First(&p, id).Error; err != nil {
			c.String(http.StatusNotFound, "Product not found")
			return
		}

		p.Name = strPtr(c.PostForm("name"))
		oldType := strVal(p.Type)
		p.Type = strPtr(c.PostForm("type"))
		price, _ := strconv.ParseInt(c.PostForm("price"), 10, 64)
		p.Price = int64Ptr(price)
		p.Description = strPtr(c.PostForm("description"))
		p.Color = strPtr(c.PostForm("color"))
		p.Fabric = strPtr(c.PostForm("fabric"))
		p.DetailFabric = strPtr(c.PostForm("detail_fabric"))
		p.Title = strPtr(c.PostForm("title"))

		if file, err := c.FormFile("banner"); err == nil && file != nil {
			if uploadedPath, err := h.Storage.Upload(file); err == nil && uploadedPath != "" {
				p.PathThumb = strPtr(uploadedPath)
			}
		}

		if oldType != *p.Type {
			var latest models.Product
			h.DB.Where("type = ?", *p.Type).Order("id DESC").First(&latest)
			prefix := ProductKeyTypes[*p.Type]
			re := regexp.MustCompile(`\d+`)
			num := 1
			if latest.Name != nil {
				if m := re.FindString(*latest.Name); m != "" {
					if n, err := strconv.Atoi(m); err == nil {
						num = n + 1
					}
				}
			}
			p.Name = strPtr(fmt.Sprintf("%s%d", prefix, num))
		}

		h.DB.Save(&p)

		// Update sizes
		sizeNames := c.PostFormArray("size_names[]")
		sizeValues := c.PostFormArray("size_values[]")
		h.DB.Where("product_id = ?", p.ID).Delete(&models.Size{})
		for i, name := range sizeNames {
			if i < len(sizeValues) && sizeValues[i] != "" {
				val, err := strconv.ParseFloat(sizeValues[i], 64)
				if err == nil && val > 0 {
					pID := p.ID
					h.DB.Create(&models.Size{
						ProductID: &pID,
						Name:      strPtr(name),
						Value:     float64Ptr(val),
					})
				}
			}
		}

		backRoute := c.PostForm("backRoute")
		if backRoute != "" {
			c.Redirect(http.StatusFound, backRoute)
			return
		}
		c.Redirect(http.StatusFound, "/sale/product/"+id)
		return
	}

	var product models.Product
	if err := h.DB.Preload("Sizes").First(&product, id).Error; err != nil {
		c.String(http.StatusNotFound, "Product not found")
		return
	}

	type SizeInput struct {
		Name  string
		Value string
	}
	var col1, col2 []SizeInput
	sizeMap := make(map[string]float64)
	for _, s := range product.Sizes {
		if s.Name != nil && s.Value != nil {
			sizeMap[*s.Name] = *s.Value
		}
	}

	for i, name := range SizeList {
		valStr := ""
		if v, ok := sizeMap[name]; ok && v > 0 {
			valStr = fmt.Sprintf("%g", v)
		}
		item := SizeInput{Name: name, Value: valStr}
		if i < 6 {
			col1 = append(col1, item)
		} else {
			col2 = append(col2, item)
		}
	}

	h.render(c, "product-detail.html", pongo2.Context{
		"current_route": "products",
		"product": gin.H{
			"ID":           product.ID,
			"Name":         strVal(product.Name),
			"Type":         strVal(product.Type),
			"Price":        int64Val(product.Price),
			"PathThumb":    strVal(product.PathThumb),
			"Description":  strVal(product.Description),
			"Color":        strVal(product.Color),
			"Fabric":       strVal(product.Fabric),
			"DetailFabric": strVal(product.DetailFabric),
			"Title":        strVal(product.Title),
			"Status":       strVal(product.Status),
		},
		"product_types":    ProductTypes,
		"size_inputs_col1": col1,
		"size_inputs_col2": col2,
		"back_route":       c.Request.Referer(),
	})
}

// 5. Product Create
func (h *SaleV1Handler) ProductCreate(c *gin.Context) {
	if c.Request.Method == "POST" {
		pType := c.PostForm("type")
		price, _ := strconv.ParseInt(c.PostForm("price"), 10, 64)
		desc := c.PostForm("description")

		var latest models.Product
		h.DB.Where("type = ?", pType).Order("id DESC").First(&latest)
		prefix := ProductKeyTypes[pType]
		re := regexp.MustCompile(`\d+`)
		num := 1
		if latest.Name != nil {
			if m := re.FindString(*latest.Name); m != "" {
				if n, err := strconv.Atoi(m); err == nil {
					num = n + 1
				}
			}
		}
		skuName := fmt.Sprintf("%s%d", prefix, num)

		p := models.Product{
			Name:        strPtr(skuName),
			Type:        strPtr(pType),
			Price:       int64Ptr(price),
			Description: strPtr(desc),
			Status:      strPtr("AVAILABLE"),
		}

		if file, err := c.FormFile("banner"); err == nil && file != nil {
			if uploadedPath, err := h.Storage.Upload(file); err == nil && uploadedPath != "" {
				p.PathThumb = strPtr(uploadedPath)
			}
		}

		h.DB.Create(&p)

		sizeNames := c.PostFormArray("size_names[]")
		sizeValues := c.PostFormArray("size_values[]")
		for i, name := range sizeNames {
			if i < len(sizeValues) && sizeValues[i] != "" {
				val, err := strconv.ParseFloat(sizeValues[i], 64)
				if err == nil && val > 0 {
					pID := p.ID
					h.DB.Create(&models.Size{
						ProductID: &pID,
						Name:      strPtr(name),
						Value:     float64Ptr(val),
					})
				}
			}
		}

		c.Redirect(http.StatusFound, "/sale/product")
		return
	}

	h.render(c, "product-create.html", pongo2.Context{
		"current_route":    "products",
		"product_types":    ProductTypes,
		"size_inputs_col1": SizeList[:6],
		"size_inputs_col2": SizeList[6:],
	})
}

// 6. Product Delete
func (h *SaleV1Handler) ProductDelete(c *gin.Context) {
	id := c.Param("id")
	h.DB.Where("id = ?", id).Delete(&models.Product{})
	c.Redirect(http.StatusFound, "/sale/product")
}

// 7. Product Upload
func (h *SaleV1Handler) ProductUpload(c *gin.Context) {
	if c.Request.Method == "POST" {
		form, err := c.MultipartForm()
		if err == nil && form != nil {
			files := form.File["banner"]
			pType := c.PostForm("type")
			for _, file := range files {
				if uploadedPath, err := h.Storage.Upload(file); err == nil && uploadedPath != "" {
					cleanName := strings.ToUpper(strings.TrimSuffix(file.Filename, ".jpg"))
					cleanName = strings.TrimSuffix(cleanName, ".JPEG")
					cleanName = strings.TrimSuffix(cleanName, ".PNG")

					var prod models.Product
					if err := h.DB.Where("name = ? AND type = ?", cleanName, pType).First(&prod).Error; err == nil {
						prod.PathThumb = strPtr(uploadedPath)
						h.DB.Save(&prod)
					}
				}
			}
		}
		c.Redirect(http.StatusFound, "/sale/available-product")
		return
	}

	h.render(c, "product-upload.html", pongo2.Context{
		"current_route": "upload",
		"product_types": ProductTypes,
	})
}

// 8. Product Batch Upload
func (h *SaleV1Handler) ProductBatchUpload(c *gin.Context) {
	h.render(c, "product-batch-upload.html", pongo2.Context{
		"current_route": "upload",
		"product_types": ProductTypes,
	})
}

func (h *SaleV1Handler) ProductBatchUploadSingle(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Image required"})
		return
	}
	pType := c.PostForm("type")

	uploadedPath, err := h.Storage.Upload(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	baseName := strings.ToUpper(strings.TrimSuffix(file.Filename, ".jpg"))
	baseName = strings.TrimSuffix(baseName, ".jpeg")
	baseName = strings.TrimSuffix(baseName, ".png")

	var prod models.Product
	matched := false
	if err := h.DB.Where("name = ?", baseName).First(&prod).Error; err == nil {
		prod.PathThumb = strPtr(uploadedPath)
		if pType != "" {
			prod.Type = strPtr(pType)
		}
		h.DB.Save(&prod)
		matched = true
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"path":    uploadedPath,
		"url":     h.Storage.BaseURL + strings.TrimPrefix(uploadedPath, "storage/"),
		"matched": matched,
		"product": baseName,
	})
}

// 9. Orders
func (h *SaleV1Handler) Orders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	perPage := 20

	search := c.Query("search")
	status := c.Query("status")

	query := h.DB.Model(&models.Order{}).Preload("Products").Order("id DESC")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if search != "" {
		sUpper := strings.ToUpper(search)
		query = query.Where("UPPER(name) LIKE ? OR phone LIKE ? OR UPPER(address) LIKE ?", "%"+sUpper+"%", "%"+search+"%", "%"+sUpper+"%")
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * perPage
	var orders []models.Order
	query.Offset(offset).Limit(perPage).Find(&orders)

	var totalOrderNew int64
	var totalOrderNewPrice float64
	h.DB.Model(&models.Order{}).Where("status IN ('NEW', 'ON HOLD')").Count(&totalOrderNew)
	h.DB.Model(&models.Order{}).Where("status IN ('NEW', 'ON HOLD')").Select("COALESCE(SUM(final_amount), 0)").Scan(&totalOrderNewPrice)

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	var currentTotalOrder int64
	var currentTotalFinalPrice float64
	h.DB.Model(&models.Order{}).Where("created_at BETWEEN ? AND ?", startOfMonth, endOfMonth).Count(&currentTotalOrder)
	h.DB.Model(&models.Order{}).Where("created_at BETWEEN ? AND ?", startOfMonth, endOfMonth).Select("COALESCE(SUM(final_amount), 0)").Scan(&currentTotalFinalPrice)

	type ProductItemView struct {
		ID        uint
		Name      string
		Type      string
		Price     int64
		PathThumb string
	}

	type OrderItemView struct {
		ID                 uint
		Name               string
		Phone              string
		Address            string
		Note               string
		Desc               string
		Status             string
		TotalAmount        int64
		FinalAmount        int64
		AdditionalAmount   int64
		FreeshipAmount     int64
		DepositAmount      int64
		RemainingAmount    int64
		ViettelOrderAmount int64
		IsPaid             bool
		Products           []ProductItemView
	}

	var orderViews []OrderItemView
	for _, o := range orders {
		tot := int64Val(o.TotalAmount)
		dep := int64Val(o.DepositAmount)
		rem := tot - dep
		if o.IsPaid {
			rem = 0
		}
		if rem < 0 {
			rem = 0
		}

		var prodList []ProductItemView
		for _, p := range o.Products {
			prodList = append(prodList, ProductItemView{
				ID:        p.ID,
				Name:      strVal(p.Name),
				Type:      strVal(p.Type),
				Price:     int64Val(p.Price),
				PathThumb: strVal(p.PathThumb),
			})
		}

		orderViews = append(orderViews, OrderItemView{
			ID:                 o.ID,
			Name:               strVal(o.Name),
			Phone:              strVal(o.Phone),
			Address:            strVal(o.Address),
			Note:               strVal(o.Note),
			Desc:               strVal(o.Desc),
			Status:             strVal(o.Status),
			TotalAmount:        tot,
			FinalAmount:        int64Val(o.FinalAmount),
			AdditionalAmount:   int64Val(o.AdditionalAmount),
			FreeshipAmount:     int64(float64Val(o.FreeshipAmount)),
			DepositAmount:      dep,
			RemainingAmount:    rem,
			ViettelOrderAmount: rem,
			IsPaid:             o.IsPaid,
			Products:           prodList,
		})
	}

	lastPage := int(math.Ceil(float64(total) / float64(perPage)))
	if lastPage < 1 {
		lastPage = 1
	}
	from := int64(offset + 1)
	if from > total {
		from = total
	}
	to := int64(offset + len(orders))

	h.render(c, "orders.html", pongo2.Context{
		"current_route":            "orders",
		"orders":                   orderViews,
		"search":                   search,
		"selected_status":          status,
		"total_order_new":          totalOrderNew,
		"total_order_new_price":    int64(totalOrderNewPrice),
		"current_total_order":      currentTotalOrder,
		"current_total_final_price": int64(currentTotalFinalPrice),
		"pagination": PaginationData{
			CurrentPage: page,
			PrevPage:    page - 1,
			NextPage:    page + 1,
			LastPage:    lastPage,
			PerPage:     perPage,
			Total:       total,
			From:        from,
			To:          to,
		},
	})
}

// 10. Summary Orders
func (h *SaleV1Handler) SummaryOrders(c *gin.Context) {
	search := c.Query("search")
	status := c.Query("status")

	query := h.DB.Model(&models.Order{}).Preload("Products").Order("id DESC")

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if search != "" {
		sUpper := strings.ToUpper(search)
		query = query.Where("UPPER(name) LIKE ? OR phone LIKE ? OR UPPER(address) LIKE ?", "%"+sUpper+"%", "%"+search+"%", "%"+sUpper+"%")
	}

	var orders []models.Order
	query.Limit(100).Find(&orders)

	var totalOrderNew int64
	var totalOrderNewPrice float64
	h.DB.Model(&models.Order{}).Where("status IN ('NEW', 'ON HOLD')").Count(&totalOrderNew)
	h.DB.Model(&models.Order{}).Where("status IN ('NEW', 'ON HOLD')").Select("COALESCE(SUM(final_amount), 0)").Scan(&totalOrderNewPrice)

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	var currentTotalOrder int64
	var currentTotalFinalPrice float64
	h.DB.Model(&models.Order{}).Where("created_at BETWEEN ? AND ?", startOfMonth, endOfMonth).Count(&currentTotalOrder)
	h.DB.Model(&models.Order{}).Where("created_at BETWEEN ? AND ?", startOfMonth, endOfMonth).Select("COALESCE(SUM(final_amount), 0)").Scan(&currentTotalFinalPrice)

	type ProductItemView struct {
		ID        uint
		Name      string
		PathThumb string
	}

	type SummaryOrderView struct {
		ID                 uint
		ConsumerID         int
		Name               string
		Phone              string
		Address            string
		Note               string
		Desc               string
		Status             string
		CreatedAtFormatted string
		TotalAmount        int64
		FreeshipAmount     int64
		DepositAmount      int64
		RemainingAmount    int64
		ViettelOrderAmount int64
		IsPaid             bool
		Products           []ProductItemView
	}
	var orderViews []SummaryOrderView
	for _, o := range orders {
		tot := int64Val(o.TotalAmount)
		dep := int64Val(o.DepositAmount)
		rem := tot - dep
		if o.IsPaid {
			rem = 0
		}
		if rem < 0 {
			rem = 0
		}

		var prodList []ProductItemView
		for _, p := range o.Products {
			prodList = append(prodList, ProductItemView{
				ID:        p.ID,
				Name:      strVal(p.Name),
				PathThumb: strVal(p.PathThumb),
			})
		}

		cTimeStr := ""
		if o.CreatedAt != nil {
			cTimeStr = o.CreatedAt.Format("02/01/2006 15:04")
		}

		consID := 0
		if o.ConsumerID != nil {
			consID = *o.ConsumerID
		}

		orderViews = append(orderViews, SummaryOrderView{
			ID:                 o.ID,
			ConsumerID:         consID,
			Name:               strVal(o.Name),
			Phone:              strVal(o.Phone),
			Address:            strVal(o.Address),
			Note:               strVal(o.Note),
			Desc:               strVal(o.Desc),
			Status:             strVal(o.Status),
			CreatedAtFormatted: cTimeStr,
			TotalAmount:        tot,
			FreeshipAmount:     int64(float64Val(o.FreeshipAmount)),
			DepositAmount:      dep,
			RemainingAmount:    rem,
			ViettelOrderAmount: rem,
			IsPaid:             o.IsPaid,
			Products:           prodList,
		})
	}

	h.render(c, "orders-summary.html", pongo2.Context{
		"current_route":            "summary",
		"orders":                   orderViews,
		"search":                   search,
		"total_order_new":          totalOrderNew,
		"total_order_new_price":    int64(totalOrderNewPrice),
		"current_total_order":      currentTotalOrder,
		"current_total_final_price": int64(currentTotalFinalPrice),
	})
}

// 11. Order Detail & Update
func (h *SaleV1Handler) OrderDetail(c *gin.Context) {
	id := c.Param("id")

	if c.Request.Method == "POST" {
		var o models.Order
		if err := h.DB.First(&o, id).Error; err != nil {
			c.String(http.StatusNotFound, "Order not found")
			return
		}

		o.Name = strPtr(c.PostForm("name"))
		o.Phone = strPtr(strings.ReplaceAll(c.PostForm("phone"), " ", ""))
		o.Address = strPtr(c.PostForm("address"))
		o.Note = strPtr(c.PostForm("note"))
		o.Desc = strPtr(c.PostForm("desc"))
		o.Status = strPtr(c.PostForm("status"))
		o.IsPaid = c.PostForm("is_paid") == "on"

		additionalAmount, _ := strconv.ParseInt(c.PostForm("additional_amount"), 10, 64)
		freeshipAmount, _ := strconv.ParseFloat(c.PostForm("freeship_amount"), 64)
		depositAmount, _ := strconv.ParseInt(c.PostForm("deposit_amount"), 10, 64)
		o.AdditionalAmount = int64Ptr(additionalAmount)
		o.FreeshipAmount = float64Ptr(freeshipAmount)
		o.DepositAmount = int64Ptr(depositAmount)

		productIDs := c.PostFormArray("products[]")
		var products []models.Product
		if len(productIDs) > 0 {
			h.DB.Where("id IN ?", productIDs).Find(&products)
		}

		var itemsTotal int64
		for _, p := range products {
			itemsTotal += int64Val(p.Price)
		}
		totalAmount := itemsTotal + additionalAmount - int64(freeshipAmount)
		o.TotalAmount = int64Ptr(totalAmount)
		o.FinalAmount = int64Ptr(totalAmount)

		h.DB.Save(&o)

		h.DB.Where("order_id = ?", o.ID).Delete(&models.OrderProduct{})
		for _, p := range products {
			oID := o.ID
			pID := p.ID
			h.DB.Create(&models.OrderProduct{
				OrderID:   &oID,
				ProductID: &pID,
			})
		}
		h.updateProductStatuses()

		backRoute := c.PostForm("backRoute")
		if backRoute != "" {
			c.Redirect(http.StatusFound, backRoute)
			return
		}
		c.Redirect(http.StatusFound, "/sale/order/"+id)
		return
	}

	var order models.Order
	if err := h.DB.Preload("Products").First(&order, id).Error; err != nil {
		c.String(http.StatusNotFound, "Order not found")
		return
	}

	var availableProducts []models.Product
	h.DB.Where("status = 'AVAILABLE'").Order("id DESC").Limit(100).Find(&availableProducts)

	type ProductItemView struct {
		ID    uint
		Name  string
		Type  string
		Price int64
	}

	var orderProds, availProds []ProductItemView
	for _, p := range order.Products {
		orderProds = append(orderProds, ProductItemView{
			ID:    p.ID,
			Name:  strVal(p.Name),
			Type:  strVal(p.Type),
			Price: int64Val(p.Price),
		})
	}
	for _, p := range availableProducts {
		availProds = append(availProds, ProductItemView{
			ID:    p.ID,
			Name:  strVal(p.Name),
			Type:  strVal(p.Type),
			Price: int64Val(p.Price),
		})
	}

	h.render(c, "order-detail.html", pongo2.Context{
		"current_route": "orders",
		"order": gin.H{
			"ID":               order.ID,
			"Name":             strVal(order.Name),
			"Phone":            strVal(order.Phone),
			"Address":          strVal(order.Address),
			"Note":             strVal(order.Note),
			"Desc":             strVal(order.Desc),
			"Status":           strVal(order.Status),
			"AdditionalAmount": int64Val(order.AdditionalAmount),
			"FreeshipAmount":   int64(float64Val(order.FreeshipAmount)),
			"DepositAmount":    int64Val(order.DepositAmount),
			"IsPaid":           order.IsPaid,
			"Products":         orderProds,
		},
		"available_products": availProds,
		"back_route":         c.Request.Referer(),
	})
}

// 12. Order Create
func (h *SaleV1Handler) OrderCreate(c *gin.Context) {
	if c.Request.Method == "POST" {
		name := c.PostForm("name")
		phone := strings.ReplaceAll(c.PostForm("phone"), " ", "")
		address := c.PostForm("address")
		note := c.PostForm("note")
		desc := c.PostForm("desc")
		isPaid := c.PostForm("is_paid") == "on"

		additionalAmount, _ := strconv.ParseInt(c.PostForm("additional_amount"), 10, 64)
		freeshipAmount, _ := strconv.ParseFloat(c.PostForm("freeship_amount"), 64)
		depositAmount, _ := strconv.ParseInt(c.PostForm("deposit_amount"), 10, 64)

		productIDs := c.PostFormArray("products[]")
		var products []models.Product
		if len(productIDs) > 0 {
			h.DB.Where("id IN ?", productIDs).Find(&products)
		}

		var itemsTotal int64
		for _, p := range products {
			itemsTotal += int64Val(p.Price)
		}
		totalAmount := itemsTotal + additionalAmount - int64(freeshipAmount)

		consumerIDStr := c.PostForm("consumer_id")
		consumerID, _ := strconv.Atoi(consumerIDStr)
		if consumerID == 0 && phone != "" {
			var cons models.Consumer
			if err := h.DB.Where("phone = ?", phone).First(&cons).Error; err != nil {
				cons = models.Consumer{Name: strPtr(name), Phone: strPtr(phone), Address: strPtr(address)}
				h.DB.Create(&cons)
			}
			consumerID = int(cons.ID)
		}

		order := models.Order{
			ConsumerID:       &consumerID,
			Name:             strPtr(name),
			Phone:            strPtr(phone),
			Address:          strPtr(address),
			Note:             strPtr(note),
			Desc:             strPtr(desc),
			Status:           strPtr("NEW"),
			TotalAmount:      int64Ptr(totalAmount),
			FinalAmount:      int64Ptr(totalAmount),
			AdditionalAmount: int64Ptr(additionalAmount),
			FreeshipAmount:   float64Ptr(freeshipAmount),
			DepositAmount:    int64Ptr(depositAmount),
			IsPaid:           isPaid,
		}
		h.DB.Create(&order)

		for _, p := range products {
			oID := order.ID
			pID := p.ID
			h.DB.Create(&models.OrderProduct{
				OrderID:   &oID,
				ProductID: &pID,
			})
		}
		h.updateProductStatuses()

		c.Redirect(http.StatusFound, "/sale/order")
		return
	}

	var availableProducts []models.Product
	h.DB.Where("status = 'AVAILABLE'").Order("id DESC").Limit(200).Find(&availableProducts)

	var consumer models.Consumer
	if cID := c.Query("consumer_id"); cID != "" {
		h.DB.First(&consumer, cID)
	}

	type ProductItemView struct {
		ID    uint
		Name  string
		Type  string
		Price int64
	}
	var availProds []ProductItemView
	for _, p := range availableProducts {
		availProds = append(availProds, ProductItemView{
			ID:    p.ID,
			Name:  strVal(p.Name),
			Type:  strVal(p.Type),
			Price: int64Val(p.Price),
		})
	}

	h.render(c, "order-create.html", pongo2.Context{
		"current_route": "orders",
		"products":      availProds,
		"consumer": gin.H{
			"ID":      consumer.ID,
			"Name":    strVal(consumer.Name),
			"Phone":   strVal(consumer.Phone),
			"Address": strVal(consumer.Address),
		},
		"back_route": c.Request.Referer(),
	})
}

// 13. Order Delete
func (h *SaleV1Handler) OrderDelete(c *gin.Context) {
	id := c.Param("id")
	h.DB.Where("id = ?", id).Delete(&models.Order{})
	h.DB.Where("order_id = ?", id).Delete(&models.OrderProduct{})
	h.updateProductStatuses()
	c.Redirect(http.StatusFound, "/sale/order")
}

// 14. Order Complete
func (h *SaleV1Handler) OrderComplete(c *gin.Context) {
	id := c.Param("id")
	h.DB.Model(&models.Order{}).Where("id = ?", id).Update("status", "COMPLETED")
	c.Redirect(http.StatusFound, c.Request.Referer())
}

// 15. Order Ship
func (h *SaleV1Handler) OrderShip(c *gin.Context) {
	id := c.Param("id")
	h.DB.Model(&models.Order{}).Where("id = ?", id).Update("status", "A SHIP NOW")
	c.Redirect(http.StatusFound, c.Request.Referer())
}

// 16. Smart Order Create
func (h *SaleV1Handler) SmartOrderCreate(c *gin.Context) {
	if c.Request.Method == "POST" {
		text := c.PostForm("description")

		phoneRe := regexp.MustCompile(`(0[3|5|7|8|9][0-9]{8})`)
		phone := phoneRe.FindString(text)

		codeRe := regexp.MustCompile(`(?i)\b(TS|HJ|JK|ST|BZ|TI|BE|GL|PS|HA|SU)\d+\b`)
		codes := codeRe.FindAllString(text, -1)

		lines := strings.Split(strings.TrimSpace(text), "\n")
		name := ""
		address := ""
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			if name == "" && !phoneRe.MatchString(line) && !codeRe.MatchString(line) {
				name = line
			} else if strings.Contains(strings.ToLower(line), "địa chỉ") ||
				strings.Contains(strings.ToLower(line), "đường") ||
				strings.Contains(strings.ToLower(line), "quận") ||
				strings.Contains(strings.ToLower(line), "phường") ||
				strings.Contains(strings.ToLower(line), "hà nội") ||
				strings.Contains(strings.ToLower(line), "hcm") ||
				strings.Contains(strings.ToLower(line), "tp") {
				address = line
			}
		}

		searchCodes := strings.Join(codes, ",")
		c.Redirect(http.StatusFound, fmt.Sprintf("/sale/order-create?name=%s&phone=%s&address=%s&note=%s",
			url.QueryEscape(name), url.QueryEscape(phone), url.QueryEscape(address), url.QueryEscape(searchCodes)))
		return
	}

	h.render(c, "smart-order-create.html", pongo2.Context{
		"current_route": "orders",
	})
}

// 17. Consumers
func (h *SaleV1Handler) Consumers(c *gin.Context) {
	type ConsumerView struct {
		ID          uint   `gorm:"column:id"`
		Name        string `gorm:"column:name"`
		Phone       string `gorm:"column:phone"`
		Address     string `gorm:"column:address"`
		OrdersCount int64  `gorm:"column:orders_count"`
	}
	var consumerViews []ConsumerView

	h.DB.Table("consumers").
		Select("consumers.id, COALESCE(consumers.name, '') as name, COALESCE(consumers.phone, '') as phone, COALESCE(consumers.address, '') as address, COUNT(orders.id) as orders_count").
		Joins("LEFT JOIN orders ON orders.consumer_id = consumers.id OR (orders.phone = consumers.phone AND consumers.phone != '')").
		Group("consumers.id, consumers.name, consumers.phone, consumers.address").
		Order("consumers.id DESC").
		Limit(100).
		Scan(&consumerViews)

	h.render(c, "consumers.html", pongo2.Context{
		"current_route": "consumers",
		"consumers":     consumerViews,
	})
}

func (h *SaleV1Handler) updateProductStatuses() {
	h.DB.Exec("UPDATE products SET status = 'AVAILABLE' WHERE status != 'ON_HOLD'")
	h.DB.Exec("UPDATE products SET status = 'SOLD' WHERE id IN (SELECT DISTINCT product_id FROM order_products)")
}

func formatProductSizes(p *models.Product) string {
	var parts []string
	var duLaiVal float64
	hasDuLai := false

	for _, s := range p.Sizes {
		if s.Name != nil && *s.Name == "DƯ LAI" && s.Value != nil {
			duLaiVal = *s.Value
			hasDuLai = true
			break
		}
	}

	for _, s := range p.Sizes {
		if s.Name == nil || s.Value == nil || *s.Name == "DƯ LAI" {
			continue
		}
		sName := *s.Name
		prefix := MappingSizes[sName]
		if prefix == "" {
			prefix = sName
		}
		val := *s.Value
		if sName == "EO QUẦN" {
			val = *s.Value * 2
		}
		part := fmt.Sprintf("%s%g", prefix, val)
		if sName == "DÀI QUẦN" && hasDuLai && duLaiVal > 0 {
			part += fmt.Sprintf(" (+%g)", duLaiVal)
		}
		parts = append(parts, part)
	}
	return strings.Join(parts, " - ")
}
