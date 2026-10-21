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
