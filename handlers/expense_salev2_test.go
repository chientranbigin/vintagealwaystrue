package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"sale-app/models"
)

func TestSaleV2_Expense_CRUD(t *testing.T) {
	db, cfg := getTestDB(t)
	tx := db.Begin()
	defer tx.Rollback()

	handler := NewSaleV2Handler(tx, cfg)
	router := setupSaleV2TestRouter(handler)

	// 1. Validation error tests
	t.Run("Create Expense - Missing Description", func(t *testing.T) {
		body := map[string]interface{}{
			"description": "",
			"amount":      100000,
			"entry_date":  "2026-10-29",
		}
		w := performRequest(router, "POST", "/salev2/api/expenses", body)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected 400 for empty description, got %d", w.Code)
		}
	})

	t.Run("Create Expense - Invalid Amount", func(t *testing.T) {
		body := map[string]interface{}{
			"description": "Shipping fee",
			"amount":      0,
			"entry_date":  "2026-10-29",
		}
		w := performRequest(router, "POST", "/salev2/api/expenses", body)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected 400 for zero amount, got %d", w.Code)
		}
	})

	// 2. Successful Create
	var createdID uint
	t.Run("Create Expense - Success", func(t *testing.T) {
		body := map[string]interface{}{
			"description": "Nhập hàng kiện áo khoác",
			"amount":      5500000,
			"entry_date":  "2026-10-15",
		}
		w := performRequest(router, "POST", "/salev2/api/expenses", body)
		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d: %s", w.Code, w.Body.String())
		}

		var resp models.Expense
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}
		if resp.ID == 0 {
			t.Fatalf("Expected non-zero ID")
		}
		if resp.Description != "Nhập hàng kiện áo khoác" {
			t.Fatalf("Expected description 'Nhập hàng kiện áo khoác', got '%s'", resp.Description)
		}
		if resp.Amount != 5500000 {
			t.Fatalf("Expected amount 5500000, got %f", resp.Amount)
		}
		createdID = resp.ID
	})

	// Add another expense for the same month and one for a different month
	exp2 := models.Expense{
		Description: "Tiền bao bì đóng gói",
		Amount:      300000,
		EntryDate:   "2026-10-20",
	}
	tx.Create(&exp2)

	exp3 := models.Expense{
		Description: "Chi phí tháng trước",
		Amount:      1000000,
		EntryDate:   "2026-09-10",
	}
	tx.Create(&exp3)

	// 3. List Expenses with Month Filter
	t.Run("List Expenses - Filter by Month", func(t *testing.T) {
		w := performRequest(router, "GET", "/salev2/api/expenses?month=2026-10-01", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d: %s", w.Code, w.Body.String())
		}

		var res struct {
			Expenses struct {
				Data        []models.Expense `json:"data"`
				Total       int64            `json:"total"`
				PerPage     int              `json:"per_page"`
				CurrentPage int              `json:"current_page"`
			} `json:"expenses"`
			Total float64 `json:"total"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if res.Expenses.Total < 2 {
			t.Fatalf("Expected at least 2 expenses in month 2026-10, got %d", res.Expenses.Total)
		}
		if res.Total < 5800000 {
			t.Fatalf("Expected month total at least 5800000, got %f", res.Total)
		}
	})

	// 4. Delete Expense
	t.Run("Delete Expense - Success", func(t *testing.T) {
		w := performRequest(router, "DELETE", fmt.Sprintf("/salev2/api/expenses/%d", createdID), nil)
		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d: %s", w.Code, w.Body.String())
		}

		// Verify deletion
		var count int64
		tx.Model(&models.Expense{}).Where("id = ?", createdID).Count(&count)
		if count != 0 {
			t.Fatalf("Expense ID %d should have been deleted", createdID)
		}
	})
}
