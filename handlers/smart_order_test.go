package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"sale-app/models"
)

func TestSmartOrderCreate_SaleV1(t *testing.T) {
	db, cfg := getTestDB(t)

	// Create a test product with status AVAILABLE
	uniqueCode := "TESTV1" + strconv.FormatInt(time.Now().UnixNano()%100000, 10)
	testProd := models.Product{
		Name:   strPtr("Sơ mi " + uniqueCode),
		Status: strPtr("AVAILABLE"),
		Type:   strPtr("SHIRT"),
		Price:  int64Ptr(150000),
	}
	if err := db.Create(&testProd).Error; err != nil {
		t.Fatalf("Failed to create test product: %v", err)
	}
	defer func() {
		db.Where("id = ?", testProd.ID).Delete(&models.Product{})
	}()

	v1 := NewSaleV1Handler(db, cfg)
	ginRouter := setupSaleV2TestRouter(NewSaleV2Handler(db, cfg))
	ginRouter.POST("/sale/smart-order-create", v1.SmartOrderCreate)

	phone := "0987654321"
	desc := "Anh Hoàng Nam " + phone + "\nSố 456 đường Lê Duẩn, Hà Nội\nLấy áo " + uniqueCode + "\nnote: gọi điện trước khi giao"

	form := url.Values{}
	form.Set("description", desc)

	req, _ := http.NewRequest("POST", "/sale/smart-order-create", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	ginRouter.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Fatalf("Expected 302 redirect, got %d: %s", w.Code, w.Body.String())
	}

	loc := w.Header().Get("Location")
	if !strings.HasPrefix(loc, "/sale/order/") {
		t.Fatalf("Expected redirect to /sale/order/:id, got %s", loc)
	}

	// Verify order was created in DB
	orderIDStr := strings.TrimPrefix(loc, "/sale/order/")
	orderID, _ := strconv.Atoi(orderIDStr)
	var createdOrder models.Order
	if err := db.Preload("Products").First(&createdOrder, orderID).Error; err != nil {
		t.Fatalf("Order was not created: %v", err)
	}
	defer func() {
		db.Where("order_id = ?", createdOrder.ID).Delete(&models.OrderProduct{})
		db.Where("id = ?", createdOrder.ID).Delete(&models.Order{})
		db.Where("phone = ?", phone).Delete(&models.Consumer{})
	}()

	if createdOrder.Phone == nil || *createdOrder.Phone != phone {
		t.Errorf("Expected order phone %s, got %v", phone, createdOrder.Phone)
	}

	if createdOrder.Status == nil || *createdOrder.Status != "NEW" {
		t.Errorf("Expected status NEW, got %v", createdOrder.Status)
	}

	// Verify test product status changed to SOLD
	var updatedProd models.Product
	db.First(&updatedProd, testProd.ID)
	if updatedProd.Status == nil || *updatedProd.Status != "SOLD" {
		t.Errorf("Expected product status SOLD, got %s", strVal(updatedProd.Status))
	}
}

func TestSmartOrderCreateAndParse_SaleV2(t *testing.T) {
	db, cfg := getTestDB(t)

	// Create test product with status AVAILABLE
	uniqueCode := "TESTV2" + strconv.Itoa(int(testProdCounter()))
	testProd := models.Product{
		Name:   strPtr("Quần " + uniqueCode),
		Status: strPtr("AVAILABLE"),
		Type:   strPtr("TROUSERS"),
		Price:  int64Ptr(220000),
	}
	if err := db.Create(&testProd).Error; err != nil {
		t.Fatalf("Failed to create test product: %v", err)
	}
	defer func() {
		db.Where("id = ?", testProd.ID).Delete(&models.Product{})
	}()

	v2 := NewSaleV2Handler(db, cfg)
	r := setupSaleV2TestRouter(v2)

	phone := "0912345678"
	desc := "Chị Linh " + phone + "\nSố 88 đường Nguyễn Trãi, Thanh Xuân, Hà Nội\nLấy quần " + uniqueCode + "\nnote: giao buổi chiều"

	// 1. Test Smart Parse API
	parsePayload, _ := json.Marshal(map[string]string{
		"description": desc,
	})
	parseReq, _ := http.NewRequest("POST", "/salev2/api/order/smart-parse", bytes.NewBuffer(parsePayload))
	parseReq.Header.Set("Content-Type", "application/json")
	parseW := httptest.NewRecorder()
	r.ServeHTTP(parseW, parseReq)

	if parseW.Code != http.StatusOK {
		t.Fatalf("SmartParse expected 200, got %d: %s", parseW.Code, parseW.Body.String())
	}

	var parseResp struct {
		Success         bool             `json:"success"`
		MatchedProducts []models.Product `json:"matched_products"`
		Parsed          struct {
			Name  string `json:"name"`
			Phone string `json:"phone"`
			Note  string `json:"note"`
		} `json:"parsed"`
	}
	if err := json.Unmarshal(parseW.Body.Bytes(), &parseResp); err != nil {
		t.Fatalf("Failed to unmarshal parse response: %v", err)
	}

	if !parseResp.Success {
		t.Errorf("Expected parse success true")
	}
	if len(parseResp.MatchedProducts) != 1 || parseResp.MatchedProducts[0].ID != testProd.ID {
		t.Errorf("Expected 1 matched product with ID %d, got %v", testProd.ID, parseResp.MatchedProducts)
	}
	if parseResp.Parsed.Phone != phone {
		t.Errorf("Expected parsed phone %s, got %s", phone, parseResp.Parsed.Phone)
	}

	// 2. Test Smart Create API
	createPayload, _ := json.Marshal(map[string]string{
		"description": desc,
	})
	createReq, _ := http.NewRequest("POST", "/salev2/api/order/smart-create", bytes.NewBuffer(createPayload))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	r.ServeHTTP(createW, createReq)

	if createW.Code != http.StatusOK {
		t.Fatalf("SmartCreate expected 200, got %d: %s", createW.Code, createW.Body.String())
	}

	var createResp struct {
		Success bool         `json:"success"`
		OrderID uint         `json:"order_id"`
		Order   models.Order `json:"order"`
	}
	if err := json.Unmarshal(createW.Body.Bytes(), &createResp); err != nil {
		t.Fatalf("Failed to unmarshal create response: %v", err)
	}

	if !createResp.Success || createResp.OrderID == 0 {
		t.Fatalf("Expected valid order creation, got %+v", createResp)
	}

	defer func() {
		db.Where("order_id = ?", createResp.OrderID).Delete(&models.OrderProduct{})
		db.Where("id = ?", createResp.OrderID).Delete(&models.Order{})
		db.Where("phone = ?", phone).Delete(&models.Consumer{})
	}()

	// Verify product status is SOLD
	var checkProd models.Product
	db.First(&checkProd, testProd.ID)
	if checkProd.Status == nil || *checkProd.Status != "SOLD" {
		t.Errorf("Expected product status SOLD after smart create, got %s", strVal(checkProd.Status))
	}
}

var counter int64 = 5000

func testProdCounter() int64 {
	counter++
	return counter
}
