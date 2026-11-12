package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"sale-app/models"
)

func TestSaleV2_Orders_SearchAndFilter(t *testing.T) {
	db, cfg := getTestDB(t)
	tx := db.Begin()
	defer tx.Rollback()

	handler := NewSaleV2Handler(tx, cfg)
	router := setupSaleV2TestRouter(handler)

	// Seed 2 products
	p1Name := fmt.Sprintf("SP_TEST_A_%d", time.Now().UnixNano()%100000)
	p2Name := fmt.Sprintf("SP_TEST_B_%d", time.Now().UnixNano()%100000)
	prod1 := models.Product{Name: &p1Name, Type: strPtr("TROUSERS"), Status: strPtr("AVAILABLE")}
	prod2 := models.Product{Name: &p2Name, Type: strPtr("BLAZER"), Status: strPtr("AVAILABLE")}
	tx.Create(&prod1)
	tx.Create(&prod2)

	// Seed 2 test orders
	phone1 := fmt.Sprintf("0981%06d", time.Now().UnixNano()%1000000)
	name1 := "Nguyen Van TestOne"
	order1 := models.Order{
		Name:        &name1,
		Phone:       &phone1,
		TotalAmount: int64Ptr(250000),
		Status:      strPtr("NEW"),
	}
	tx.Create(&order1)
	tx.Create(&models.OrderProduct{OrderID: &order1.ID, ProductID: &prod1.ID})

	phone2 := fmt.Sprintf("0982%06d", time.Now().UnixNano()%1000000)
	name2 := "Tran Thi TestTwo"
	order2 := models.Order{
		Name:        &name2,
		Phone:       &phone2,
		TotalAmount: int64Ptr(500000),
		Status:      strPtr("COMPLETED"),
	}
	tx.Create(&order2)
	tx.Create(&models.OrderProduct{OrderID: &order2.ID, ProductID: &prod2.ID})

	// 1. Search by Phone
	t.Run("Search by Phone", func(t *testing.T) {
		w := performRequest(router, "GET", "/salev2/api/orders?search="+phone1, nil)
		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d: %s", w.Code, w.Body.String())
		}
		var res struct {
			Data  []map[string]interface{} `json:"data"`
			Total int                       `json:"total"`
		}
		json.Unmarshal(w.Body.Bytes(), &res)
		if len(res.Data) == 0 {
			t.Fatalf("Expected order to be found with phone %s", phone1)
		}
		foundPhone := res.Data[0]["phone"].(string)
		if foundPhone != phone1 {
			t.Errorf("Expected phone %s, got %s", phone1, foundPhone)
		}
	})

	// 2. Search by Customer Name
	t.Run("Search by Name", func(t *testing.T) {
		w := performRequest(router, "GET", "/salev2/api/orders?search=Nguyen+Van+TestOne", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d", w.Code)
		}
		var res struct {
			Data []map[string]interface{} `json:"data"`
		}
		json.Unmarshal(w.Body.Bytes(), &res)
		if len(res.Data) == 0 {
			t.Fatalf("Expected order to be found by name")
		}
	})

	// 3. Search by Product SKU in Order
	t.Run("Search by Product SKU", func(t *testing.T) {
		w := performRequest(router, "GET", "/salev2/api/orders?search="+p1Name, nil)
		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d", w.Code)
		}
		var res struct {
			Data []map[string]interface{} `json:"data"`
		}
		json.Unmarshal(w.Body.Bytes(), &res)
		if len(res.Data) == 0 {
			t.Fatalf("Expected order to be found by product name in order %s", p1Name)
		}
	})

	// 4. Filter by status Tab: NEW vs ALL
	t.Run("Filter by Tab Status NEW", func(t *testing.T) {
		w := performRequest(router, "GET", "/salev2/api/orders?status=NEW&search="+phone2, nil)
		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d", w.Code)
		}
		var res struct {
			Data []map[string]interface{} `json:"data"`
		}
		json.Unmarshal(w.Body.Bytes(), &res)
		// order2 is COMPLETED, so it should NOT appear when status=NEW
		if len(res.Data) != 0 {
			t.Errorf("Expected 0 orders for COMPLETED order under status=NEW, got %d", len(res.Data))
		}
	})

	// 5. Pagination
	t.Run("Pagination Check", func(t *testing.T) {
		w := performRequest(router, "GET", "/salev2/api/orders?page=1&limit=5", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d", w.Code)
		}
		var res struct {
			Data        []map[string]interface{} `json:"data"`
			CurrentPage int                      `json:"current_page"`
			Total       int                      `json:"total"`
		}
		json.Unmarshal(w.Body.Bytes(), &res)
		if res.CurrentPage != 1 {
			t.Errorf("Expected page 1, got %d", res.CurrentPage)
		}
		if len(res.Data) > 5 {
			t.Errorf("Expected max 5 orders per page, got %d", len(res.Data))
		}
	})
}

func TestSaleV2_StoreOrder_CompleteLogic(t *testing.T) {
	db, cfg := getTestDB(t)
	tx := db.Begin()
	defer tx.Rollback()

	handler := NewSaleV2Handler(tx, cfg)
	router := setupSaleV2TestRouter(handler)

	// Seed 2 available products
	p1Price := int64(200000)
	p2Price := int64(150000)
	prod1 := models.Product{Name: strPtr("ORD_P1"), Price: &p1Price, Status: strPtr("AVAILABLE")}
	prod2 := models.Product{Name: strPtr("ORD_P2"), Price: &p2Price, Status: strPtr("AVAILABLE")}
	tx.Create(&prod1)
	tx.Create(&prod2)

	phone := fmt.Sprintf("0977%06d", time.Now().UnixNano()%1000000)
	name := "Khach Hang Mua Don"
	address := "123 Duong Le Loi, Quan 1"

	// Request with duplicated product IDs to verify backend deduplication
	payload := map[string]interface{}{
		"name":              name,
		"phone":             phone,
		"address":           address,
		"note":              "Giao buoi sang",
		"total_amount":      320000, // 200k + 150k - 30k freeship
		"deposit_amount":    50000,
		"additional_amount": 0,
		"freeship_amount":   30000,
		"is_freeship":       true,
		"is_paid":           false,
		"status":            "NEW",
		"products":          []uint{prod1.ID, prod2.ID, prod1.ID}, // prod1 repeated!
	}

	w := performRequest(router, "POST", "/salev2/api/order", payload)
	if w.Code != http.StatusOK {
		t.Fatalf("StoreOrder returned status %d: %s", w.Code, w.Body.String())
	}

	var res struct {
		Success bool `json:"success"`
		Order   struct {
			ID          uint    `json:"id"`
			TotalAmount int64   `json:"total_amount"`
			Freeship    float64 `json:"freeship_amount"`
		} `json:"order"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	createdOrderID := res.Order.ID
	if createdOrderID == 0 {
		t.Fatalf("Expected valid order ID, got 0")
	}

	// 1. Verify DB order record
	var dbOrder models.Order
	if err := tx.First(&dbOrder, createdOrderID).Error; err != nil {
		t.Fatalf("Order not found in DB: %v", err)
	}
	if *dbOrder.Phone != phone || *dbOrder.TotalAmount != 320000 || *dbOrder.FreeshipAmount != 30000 {
		t.Errorf("Order fields mismatch: phone=%v, total=%v, freeship=%v", *dbOrder.Phone, *dbOrder.TotalAmount, *dbOrder.FreeshipAmount)
	}

	// 2. Verify Deduplication in order_products
	var orderProds []models.OrderProduct
	tx.Where("order_id = ?", createdOrderID).Find(&orderProds)
	if len(orderProds) != 2 {
		t.Fatalf("Expected exactly 2 order_products after deduplication, got %d", len(orderProds))
	}

	// 3. Verify Products status set to SOLD
	var checkP1, checkP2 models.Product
	tx.First(&checkP1, prod1.ID)
	tx.First(&checkP2, prod2.ID)
	if *checkP1.Status != "SOLD" {
		t.Errorf("Expected prod1 status=SOLD, got status=%s", *checkP1.Status)
	}
	if *checkP2.Status != "SOLD" {
		t.Errorf("Expected prod2 status=SOLD, got status=%s", *checkP2.Status)
	}

	// 4. Verify Consumer is created
	var consumer models.Consumer
	if err := tx.Where("phone = ?", phone).First(&consumer).Error; err != nil {
		t.Fatalf("Expected consumer to be created with phone %s: %v", phone, err)
	}
	if *consumer.Name != name {
		t.Errorf("Consumer name mismatch: expected %s, got %s", name, *consumer.Name)
	}
}

func TestSaleV2_UpdateOrder_DeduplicationAndProducts(t *testing.T) {
	db, cfg := getTestDB(t)
	tx := db.Begin()
	defer tx.Rollback()

	handler := NewSaleV2Handler(tx, cfg)
	router := setupSaleV2TestRouter(handler)

	// Seed 3 products
	p1 := models.Product{Name: strPtr("UP_P1"), Status: strPtr("SOLD")}
	p2 := models.Product{Name: strPtr("UP_P2"), Status: strPtr("SOLD")}
	p3 := models.Product{Name: strPtr("UP_P3"), Status: strPtr("AVAILABLE")}
	tx.Create(&p1)
	tx.Create(&p2)
	tx.Create(&p3)

	// Seed order with p1 and p2
	orderName := "Update Test Order"
	orderPhone := "0911223344"
	order := models.Order{
		Name:        &orderName,
		Phone:       &orderPhone,
		TotalAmount: int64Ptr(400000),
		Status:      strPtr("NEW"),
	}
	tx.Create(&order)
	tx.Create(&models.OrderProduct{OrderID: &order.ID, ProductID: &p1.ID})
	tx.Create(&models.OrderProduct{OrderID: &order.ID, ProductID: &p2.ID})

	// Update order: remove p2, keep p1, add p3, and duplicate p3 in request
	updatePayload := map[string]interface{}{
		"name":              "Updated Name",
		"phone":             "0911223399",
		"address":           "Updated Address",
		"note":              "Updated Note",
		"total_amount":      350000,
		"freeship_amount":   20000,
		"additional_amount": 10000,
		"status":            "SHIPPED",
		"products":          []uint{p1.ID, p3.ID, p3.ID}, // p3 duplicated! p2 removed!
	}

	w := performRequest(router, "POST", fmt.Sprintf("/salev2/api/order/%d/update", order.ID), updatePayload)
	if w.Code != http.StatusOK {
		t.Fatalf("UpdateOrder failed with status %d: %s", w.Code, w.Body.String())
	}

	// 1. Verify no duplicate order_products
	var orderProds []models.OrderProduct
	tx.Where("order_id = ?", order.ID).Find(&orderProds)
	if len(orderProds) != 2 {
		t.Fatalf("Expected exactly 2 order_products after update deduplication, got %d", len(orderProds))
	}

	// 2. Verify p3 status became SOLD
	var checkP3 models.Product
	tx.First(&checkP3, p3.ID)
	if *checkP3.Status != "SOLD" {
		t.Errorf("Expected newly added p3 to be status=SOLD, got %s", *checkP3.Status)
	}

	// 3. Verify removed product p2 reverted to AVAILABLE
	var checkP2 models.Product
	tx.First(&checkP2, p2.ID)
	if *checkP2.Status != "AVAILABLE" {
		t.Errorf("Expected removed p2 to be released to status=AVAILABLE, got %s", *checkP2.Status)
	}

	// 4. Verify updated amounts
	var updatedOrder models.Order
	tx.First(&updatedOrder, order.ID)
	if *updatedOrder.Status != "SHIPPED" || *updatedOrder.TotalAmount != 350000 || *updatedOrder.FreeshipAmount != 20000 {
		t.Errorf("Updated order attributes mismatch: status=%s, total=%v, freeship=%v",
			*updatedOrder.Status, *updatedOrder.TotalAmount, *updatedOrder.FreeshipAmount)
	}
}

func TestSaleV2_DestroyOrder(t *testing.T) {
	db, cfg := getTestDB(t)
	tx := db.Begin()
	defer tx.Rollback()

	handler := NewSaleV2Handler(tx, cfg)
	router := setupSaleV2TestRouter(handler)

	// Seed product and order
	prod := models.Product{Name: strPtr("DEL_P1"), Status: strPtr("SOLD")}
	tx.Create(&prod)

	order := models.Order{TotalAmount: int64Ptr(100000), Status: strPtr("NEW")}
	tx.Create(&order)
	tx.Create(&models.OrderProduct{OrderID: &order.ID, ProductID: &prod.ID})

	// Call DestroyOrder
	w := performRequest(router, "DELETE", fmt.Sprintf("/salev2/api/order/%d", order.ID), nil)
	if w.Code != http.StatusOK {
		t.Fatalf("DestroyOrder failed with status %d: %s", w.Code, w.Body.String())
	}

	// 1. Verify order is deleted
	var deletedOrder models.Order
	err := tx.First(&deletedOrder, order.ID).Error
	if err == nil {
		t.Errorf("Expected order to be deleted, but it still exists")
	}

	// 2. Verify product is released back to AVAILABLE
	var checkProd models.Product
	tx.First(&checkProd, prod.ID)
	if *checkProd.Status != "AVAILABLE" {
		t.Errorf("Expected product to revert to status=AVAILABLE upon order deletion, got %s", *checkProd.Status)
	}
}

func TestSaleV2_OrderStats(t *testing.T) {
	db, cfg := getTestDB(t)
	tx := db.Begin()
	defer tx.Rollback()

	handler := NewSaleV2Handler(tx, cfg)
	router := setupSaleV2TestRouter(handler)

	// Seed product and order for today
	pPrice := int64(300000)
	prod := models.Product{Name: strPtr("STAT_P1"), Price: &pPrice, Status: strPtr("SOLD")}
	tx.Create(&prod)

	now := time.Now()
	order := models.Order{
		Name:        strPtr("Stat Tester"),
		TotalAmount: int64Ptr(300000),
		FinalAmount: int64Ptr(300000),
		Status:      strPtr("NEW"),
		CreatedAt:   &now,
	}
	tx.Create(&order)
	tx.Create(&models.OrderProduct{OrderID: &order.ID, ProductID: &prod.ID, CreatedAt: &now})

	w := performRequest(router, "GET", "/salev2/api/orders-stats", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var res struct {
		Yesterday struct {
			Count   int64 `json:"count"`
			Revenue int64 `json:"revenue"`
		} `json:"yesterday"`
		Today struct {
			Count   int64 `json:"count"`
			Revenue int64 `json:"revenue"`
		} `json:"today"`
		CurrentMonth struct {
			Count   int64 `json:"count"`
			Revenue int64 `json:"revenue"`
		} `json:"current_month"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("Failed to parse json: %v", err)
	}

	if res.Today.Count < 1 || res.CurrentMonth.Count < 1 {
		t.Errorf("Expected today and month count >= 1, got today=%d, month=%d", res.Today.Count, res.CurrentMonth.Count)
	}
}

func TestSaleV2_FreeshipAndDiscountFormulas_CreateAndUpdate(t *testing.T) {
	db, cfg := getTestDB(t)
	tx := db.Begin()
	defer tx.Rollback()

	handler := NewSaleV2Handler(tx, cfg)
	router := setupSaleV2TestRouter(handler)

	// 1. Order 400k with freeship = true, shipping_fee = 20k -> total = 400k, final = 380k, viettel = 400k
	p1Price := int64(400000)
	prod1 := models.Product{Name: strPtr("FS_P1"), Price: &p1Price, Status: strPtr("AVAILABLE")}
	tx.Create(&prod1)

	phone1 := fmt.Sprintf("0961%06d", time.Now().UnixNano()%1000000)
	payload1 := map[string]interface{}{
		"name":              "Khach Freeship 20k",
		"phone":             phone1,
		"address":           "HCM",
		"total_amount":      400000,
		"additional_amount": 0,
		"freeship_amount":   0,
		"shipping_fee":      20000,
		"deposit_amount":    0,
		"is_freeship":       true,
		"is_paid":           false,
		"status":            "NEW",
		"products":          []uint{prod1.ID},
	}

	w1 := performRequest(router, "POST", "/salev2/api/order", payload1)
	if w1.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d: %s", w1.Code, w1.Body.String())
	}
	var res1 struct {
		Order models.Order `json:"order"`
	}
	json.Unmarshal(w1.Body.Bytes(), &res1)

	if res1.Order.TotalAmount == nil || *res1.Order.TotalAmount != 400000 {
		t.Errorf("Expected TotalAmount=400000, got %v", res1.Order.TotalAmount)
	}
	if res1.Order.FinalAmount == nil || *res1.Order.FinalAmount != 380000 {
		t.Errorf("Expected FinalAmount=380000 (400k - 20k freeship), got %v", res1.Order.FinalAmount)
	}
	if res1.Order.ShippingFee == nil || *res1.Order.ShippingFee != 20000 {
		t.Errorf("Expected ShippingFee=20000, got %v", res1.Order.ShippingFee)
	}

	// Verify GET /salev2/api/orders returns viettel_amount = 400000
	wList := performRequest(router, "GET", "/salev2/api/orders?search="+phone1, nil)
	var listRes struct {
		Data []map[string]interface{} `json:"data"`
	}
	json.Unmarshal(wList.Body.Bytes(), &listRes)
	if len(listRes.Data) != 1 {
		t.Fatalf("Expected 1 order in list, got %d", len(listRes.Data))
	}
	if int64(listRes.Data[0]["viettel_amount"].(float64)) != 400000 {
		t.Errorf("Expected viettel_amount=400000, got %v", listRes.Data[0]["viettel_amount"])
	}
	if int64(listRes.Data[0]["final_amount"].(float64)) != 380000 {
		t.Errorf("Expected final_amount=380000, got %v", listRes.Data[0]["final_amount"])
	}

	// 2. Order with is_freeship = true and omitted shipping_fee -> defaults to 15000
	p2Price := int64(400000)
	prod2 := models.Product{Name: strPtr("FS_P2"), Price: &p2Price, Status: strPtr("AVAILABLE")}
	tx.Create(&prod2)

	phone2 := fmt.Sprintf("0962%06d", time.Now().UnixNano()%1000000)
	payload2 := map[string]interface{}{
		"name":            "Khach Default 15k",
		"phone":           phone2,
		"total_amount":    370000, // 400k - 30k discount
		"freeship_amount": 30000,  // discount
		"deposit_amount":  50000,
		"is_freeship":     true,
		"is_paid":         false,
		"status":          "NEW",
		"products":        []uint{prod2.ID},
	}
	w2 := performRequest(router, "POST", "/salev2/api/order", payload2)
	if w2.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d: %s", w2.Code, w2.Body.String())
	}
	var res2 struct {
		Order models.Order `json:"order"`
	}
	json.Unmarshal(w2.Body.Bytes(), &res2)
	if res2.Order.ShippingFee == nil || *res2.Order.ShippingFee != 15000 {
		t.Errorf("Expected default ShippingFee=15000, got %v", res2.Order.ShippingFee)
	}
	if res2.Order.FinalAmount == nil || *res2.Order.FinalAmount != 355000 {
		t.Errorf("Expected FinalAmount=355000 (370k - 15k), got %v", res2.Order.FinalAmount)
	}

	// 3. Update order2: change shipping_fee to 25000 -> final_amount = 345000
	updatePayload := map[string]interface{}{
		"name":            "Khach Default 15k",
		"phone":           phone2,
		"total_amount":    370000,
		"freeship_amount": 30000,
		"shipping_fee":    25000,
		"deposit_amount":  50000,
		"is_freeship":     true,
		"is_paid":         false,
		"status":          "NEW",
		"products":        []uint{prod2.ID},
	}
	wUp1 := performRequest(router, "POST", fmt.Sprintf("/salev2/api/order/%d/update", res2.Order.ID), updatePayload)
	if wUp1.Code != http.StatusOK {
		t.Fatalf("Expected 200 on update, got %d: %s", wUp1.Code, wUp1.Body.String())
	}
	var checkOrd models.Order
	tx.First(&checkOrd, res2.Order.ID)
	if *checkOrd.ShippingFee != 25000 || *checkOrd.FinalAmount != 345000 {
		t.Errorf("Expected ShippingFee=25000 and FinalAmount=345000, got fee=%v, final=%v", *checkOrd.ShippingFee, *checkOrd.FinalAmount)
	}

	// 4. Update order2: toggle is_freeship = false -> shipping_fee becomes 0, final_amount = 370000
	updatePayload["is_freeship"] = false
	wUp2 := performRequest(router, "POST", fmt.Sprintf("/salev2/api/order/%d/update", res2.Order.ID), updatePayload)
	if wUp2.Code != http.StatusOK {
		t.Fatalf("Expected 200 on update off, got %d", wUp2.Code)
	}
	tx.First(&checkOrd, res2.Order.ID)
	if *checkOrd.ShippingFee != 0 || *checkOrd.FinalAmount != 370000 {
		t.Errorf("Expected ShippingFee=0 and FinalAmount=370000 when freeship off, got fee=%v, final=%v", *checkOrd.ShippingFee, *checkOrd.FinalAmount)
	}
}

func TestSaleV2_MultiDayOrderUpdate_RevenueAndPreservedCreatedAt(t *testing.T) {
	db, cfg := getTestDB(t)
	tx := db.Begin()
	defer tx.Rollback()

	handler := NewSaleV2Handler(tx, cfg)
	router := setupSaleV2TestRouter(handler)

	// Measure baseline stats before seeding
	getStats := func() (yesterdayCount, yesterdayRev, todayCount, todayRev int64) {
		w := performRequest(router, "GET", "/salev2/api/orders-stats", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("OrderStats failed: %d", w.Code)
		}
		var s struct {
			Yesterday struct {
				Count   int64 `json:"count"`
				Revenue int64 `json:"revenue"`
			} `json:"yesterday"`
			Today struct {
				Count   int64 `json:"count"`
				Revenue int64 `json:"revenue"`
			} `json:"today"`
		}
		json.Unmarshal(w.Body.Bytes(), &s)
		return s.Yesterday.Count, s.Yesterday.Revenue, s.Today.Count, s.Today.Revenue
	}

	baseYCount, baseYRev, baseTCount, baseTRev := getStats()

	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		loc = time.FixedZone("GMT+7", 7*3600)
	}
	now := time.Now().In(loc)
	yesterday := now.AddDate(0, 0, -1)
	yesterdayStr := yesterday.Format("2006-01-02")
	todayStr := now.Format("2006-01-02")

	// Seed 2 products
	p1Price := int64(400000)
	p2Price := int64(200000)
	prod1 := models.Product{Name: strPtr("MD_P1"), Type: strPtr("SHIRT"), Price: &p1Price, Status: strPtr("SOLD")}
	prod2 := models.Product{Name: strPtr("MD_P2"), Type: strPtr("PANTS"), Price: &p2Price, Status: strPtr("AVAILABLE")}
	tx.Create(&prod1)
	tx.Create(&prod2)

	// Create order Yesterday with prod1 (400k), discount 10k, freeship 15k -> total 390k, final 375k
	isFS := true
	disc := float64(10000)
	shipFee := float64(15000)
	order := models.Order{
		Name:           strPtr("MultiDay Customer"),
		Phone:          strPtr("0909888777"),
		TotalAmount:    int64Ptr(390000),
		FinalAmount:    int64Ptr(375000),
		FreeshipAmount: &disc,
		ShippingFee:    &shipFee,
		IsFreeship:     isFS,
		Status:         strPtr("ON HOLD"),
		CreatedAt:      &yesterday,
		UpdatedAt:      &yesterday,
	}
	tx.Create(&order)
	tx.Create(&models.OrderProduct{
		OrderID:   &order.ID,
		ProductID: &prod1.ID,
		CreatedAt: &yesterday,
		UpdatedAt: &yesterday,
	})

	// Check stats right after yesterday's order is seeded
	yCount1, yRev1, tCount1, tRev1 := getStats()
	if yCount1-baseYCount != 1 {
		t.Errorf("Expected yesterday count +1, got +%d", yCount1-baseYCount)
	}
	if yRev1-baseYRev != 375000 {
		t.Errorf("Expected yesterday revenue +375000 (400k - 10k discount - 15k freeship), got +%d", yRev1-baseYRev)
	}
	if tCount1-baseTCount != 0 || tRev1-baseTRev != 0 {
		t.Errorf("Expected today stats unchanged before update, got count +%d, rev +%d", tCount1-baseTCount, tRev1-baseTRev)
	}

	// Now Today customer buys prod2 (200k) as well -> update order with [prod1.ID, prod2.ID]
	// Total = 400k + 200k - 10k = 590k, Final = 590k - 15k = 575k
	updatePayload := map[string]interface{}{
		"name":              "MultiDay Customer",
		"phone":             "0909888777",
		"total_amount":      590000,
		"freeship_amount":   10000,
		"shipping_fee":      15000,
		"additional_amount": 0,
		"deposit_amount":    0,
		"is_freeship":       true,
		"is_paid":           false,
		"status":            "NEW",
		"products":          []uint{prod1.ID, prod2.ID},
	}

	wUp := performRequest(router, "POST", fmt.Sprintf("/salev2/api/order/%d/update", order.ID), updatePayload)
	if wUp.Code != http.StatusOK {
		t.Fatalf("UpdateOrder failed: %d %s", wUp.Code, wUp.Body.String())
	}

	// Verify order_products.created_at in DB
	var op1, op2 models.OrderProduct
	if err := tx.Where("order_id = ? AND product_id = ?", order.ID, prod1.ID).First(&op1).Error; err != nil {
		t.Fatalf("Missing order_product for prod1: %v", err)
	}
	if err := tx.Where("order_id = ? AND product_id = ?", order.ID, prod2.ID).First(&op2).Error; err != nil {
		t.Fatalf("Missing order_product for prod2: %v", err)
	}

	op1Date := op1.CreatedAt.In(loc).Format("2006-01-02")
	op2Date := op2.CreatedAt.In(loc).Format("2006-01-02")
	if op1Date != yesterdayStr {
		t.Errorf("Expected prod1 order_product created_at date to remain %s, got %s", yesterdayStr, op1Date)
	}
	if op2Date != todayStr {
		t.Errorf("Expected prod2 order_product created_at date to be %s, got %s", todayStr, op2Date)
	}

	// Verify OrderStats after adding prod2 today:
	// Yesterday revenue stays +375,000; Today revenue becomes +200,000; Sum = 575,000 == FinalAmount!
	yCount2, yRev2, tCount2, tRev2 := getStats()
	if yCount2-baseYCount != 1 || yRev2-baseYRev != 375000 {
		t.Errorf("Expected yesterday stats to stay +1 / +375000 after today's update, got +%d / +%d", yCount2-baseYCount, yRev2-baseYRev)
	}
	if tCount2-baseTCount != 1 || tRev2-baseTRev != 200000 {
		t.Errorf("Expected today stats to be +1 / +200000 after adding prod2 today, got +%d / +%d", tCount2-baseTCount, tRev2-baseTRev)
	}
	if (yRev2-baseYRev)+(tRev2-baseTRev) != 575000 {
		t.Errorf("Expected sum of daily revenues to equal FinalAmount (575000), got %d", (yRev2-baseYRev)+(tRev2-baseTRev))
	}

	// Verify SoldByDateProducts returns prod1 on yesterdayStr and prod2 on todayStr
	wProdY := performRequest(router, "GET", "/salev2/api/sold-by-date/"+yesterdayStr+"/products", nil)
	if wProdY.Code != http.StatusOK {
		t.Fatalf("SoldByDateProducts yesterday failed: %d", wProdY.Code)
	}
	var prodYRes struct {
		Products []models.Product `json:"products"`
	}
	json.Unmarshal(wProdY.Body.Bytes(), &prodYRes)
	foundP1OnY := false
	for _, p := range prodYRes.Products {
		if p.ID == prod1.ID {
			foundP1OnY = true
		}
		if p.ID == prod2.ID {
			t.Errorf("prod2 (added today) should NOT appear in yesterday's SoldByDateProducts")
		}
	}
	if !foundP1OnY {
		t.Errorf("Expected prod1 to appear in yesterday's SoldByDateProducts")
	}

	wProdT := performRequest(router, "GET", "/salev2/api/sold-by-date/"+todayStr+"/products", nil)
	if wProdT.Code != http.StatusOK {
		t.Fatalf("SoldByDateProducts today failed: %d", wProdT.Code)
	}
	var prodTRes struct {
		Products []models.Product `json:"products"`
	}
	json.Unmarshal(wProdT.Body.Bytes(), &prodTRes)
	foundP2OnT := false
	for _, p := range prodTRes.Products {
		if p.ID == prod2.ID {
			foundP2OnT = true
		}
		if p.ID == prod1.ID {
			t.Errorf("prod1 (added yesterday) should NOT appear in today's SoldByDateProducts")
		}
	}
	if !foundP2OnT {
		t.Errorf("Expected prod2 to appear in today's SoldByDateProducts")
	}
}

