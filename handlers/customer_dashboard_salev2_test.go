package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"sale-app/models"
)

func TestSaleV2_Customers_List_Detail_History(t *testing.T) {
	db, cfg := getTestDB(t)
	tx := db.Begin()
	defer tx.Rollback()

	handler := NewSaleV2Handler(tx, cfg)
	router := setupSaleV2TestRouter(handler)

	// Seed customer
	phone := fmt.Sprintf("0933%06d", time.Now().UnixNano()%1000000)
	name := "Tran Dai Nghia"
	address := "99 Vo Thi Sau, Q3"
	consumer := models.Consumer{
		Name:    &name,
		Phone:   &phone,
		Address: &address,
	}
	tx.Create(&consumer)

	// Seed past order for this customer
	order := models.Order{
		Name:        &name,
		Phone:       &phone,
		TotalAmount: int64Ptr(750000),
		Status:      strPtr("COMPLETED"),
	}
	tx.Create(&order)

	// 1. List Consumers
	t.Run("List Consumers with pagination", func(t *testing.T) {
		w := performRequest(router, "GET", "/salev2/api/consumers?page=1&limit=10", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d", w.Code)
		}
		var res struct {
			Data  []map[string]interface{} `json:"data"`
			Total int                       `json:"total"`
		}
		json.Unmarshal(w.Body.Bytes(), &res)
		if len(res.Data) == 0 {
			t.Errorf("Expected consumers in list")
		}
	})

	// 2. Consumer Detail
	t.Run("Get Consumer Detail by ID", func(t *testing.T) {
		w := performRequest(router, "GET", fmt.Sprintf("/salev2/api/consumer/%d", consumer.ID), nil)
		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d", w.Code)
		}
		var res models.Consumer
		json.Unmarshal(w.Body.Bytes(), &res)
		if res.ID != consumer.ID || *res.Phone != phone {
			t.Errorf("Consumer mismatch: expected phone %s, got %v", phone, res.Phone)
		}
	})

	// 3. Customer History by Phone
	t.Run("Get Customer Order History by Phone", func(t *testing.T) {
		w := performRequest(router, "GET", "/salev2/api/customer-history/"+phone, nil)
		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d", w.Code)
		}
		var orders []models.Order
		json.Unmarshal(w.Body.Bytes(), &orders)
		if len(orders) == 0 {
			t.Fatalf("Expected at least 1 historical order for phone %s", phone)
		}
		if orders[0].ID != order.ID {
			t.Errorf("Expected order ID %d, got %d", order.ID, orders[0].ID)
		}
	})
}

func TestSaleV2_DashboardStats(t *testing.T) {
	db, cfg := getTestDB(t)
	tx := db.Begin()
	defer tx.Rollback()

	handler := NewSaleV2Handler(tx, cfg)
	router := setupSaleV2TestRouter(handler)

	// Seed test order for today
	orderAmt := int64(1200000)
	order := models.Order{
		TotalAmount: &orderAmt,
		FinalAmount: &orderAmt,
		Status:      strPtr("COMPLETED"),
	}
	tx.Create(&order)

	w := performRequest(router, "GET", "/salev2/api/dashboard-stats", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("DashboardStats returned status %d: %s", w.Code, w.Body.String())
	}

	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)

	// Verify required dashboard metrics are present
	expectedKeys := []string{"all", "current_month", "prev_month", "monthly_sales", "total_orders", "total_revenue"}
	for _, k := range expectedKeys {
		if _, ok := res[k]; !ok {
			t.Errorf("Expected key %q in dashboard stats response, but missing", k)
		}
	}
}
