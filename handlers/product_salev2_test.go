package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"sale-app/models"
)

func TestSaleV2_Products_SearchFilterAndPagination(t *testing.T) {
	db, cfg := getTestDB(t)
	tx := db.Begin()
	defer tx.Rollback()

	handler := NewSaleV2Handler(tx, cfg)
	router := setupSaleV2TestRouter(handler)

	// Seed products
	codeA := fmt.Sprintf("TS_TEST_%d", time.Now().UnixNano()%100000)
	codeB := fmt.Sprintf("BZ_TEST_%d", time.Now().UnixNano()%100000)
	priceA := int64(300000)
	priceB := int64(600000)

	prodA := models.Product{
		Name:   &codeA,
		Type:   strPtr("TROUSERS"),
		Price:  &priceA,
		Status: strPtr("AVAILABLE"),
	}
	prodB := models.Product{
		Name:   &codeB,
		Type:   strPtr("BLAZER"),
		Price:  &priceB,
		Status: strPtr("SOLD"),
	}
	tx.Create(&prodA)
	tx.Create(&prodB)

	// Attach size to prodA: dai_ao = 75.5
	valSize := 75.5
	tx.Create(&models.Size{
		ProductID: &prodA.ID,
		Name:      strPtr("dai_ao"),
		Value:     &valSize,
	})

	// 1. Search by product code
	t.Run("Search by product code", func(t *testing.T) {
		w := performRequest(router, "GET", "/salev2/api/products?search="+codeA, nil)
		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d", w.Code)
		}
		var res struct {
			Data []map[string]interface{} `json:"data"`
		}
		json.Unmarshal(w.Body.Bytes(), &res)
		if len(res.Data) == 0 {
			t.Fatalf("Expected product %s to be found", codeA)
		}
	})

	// 2. Filter by Status (AVAILABLE vs SOLD)
	t.Run("Filter by status AVAILABLE", func(t *testing.T) {
		w := performRequest(router, "GET", "/salev2/api/products?status=AVAILABLE&search="+codeB, nil)
		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d", w.Code)
		}
		var res struct {
			Data []map[string]interface{} `json:"data"`
		}
		json.Unmarshal(w.Body.Bytes(), &res)
		if len(res.Data) != 0 {
			t.Errorf("Expected SOLD product not to appear in status=AVAILABLE")
		}
	})

	// 3. Filter by Category / Type
	t.Run("Filter by Type", func(t *testing.T) {
		w := performRequest(router, "GET", "/salev2/api/products?type=TROUSERS&search="+codeA, nil)
		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d", w.Code)
		}
		var res struct {
			Data []map[string]interface{} `json:"data"`
		}
		json.Unmarshal(w.Body.Bytes(), &res)
		if len(res.Data) == 0 {
			t.Fatalf("Expected product of type TROUSERS to be found")
		}
	})

	// 4. Detailed Size Range Filter
	t.Run("Filter by Size Range", func(t *testing.T) {
		url := fmt.Sprintf("/salev2/api/products?from_search_key=dai_ao&from_search_value=70&to_search_value=80&search=%s", codeA)
		w := performRequest(router, "GET", url, nil)
		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d", w.Code)
		}
		var res struct {
			Data []map[string]interface{} `json:"data"`
		}
		json.Unmarshal(w.Body.Bytes(), &res)
		if len(res.Data) == 0 {
			t.Fatalf("Expected product with dai_ao=75.5 to match filter [70, 80]")
		}
	})

	// 5. Pagination
	t.Run("Pagination check", func(t *testing.T) {
		w := performRequest(router, "GET", "/salev2/api/products?page=1&limit=10", nil)
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
			t.Errorf("Expected current_page=1, got %d", res.CurrentPage)
		}
	})
}

func TestSaleV2_ProductDetail_Store_Update_Destroy(t *testing.T) {
	db, cfg := getTestDB(t)
	tx := db.Begin()
	defer tx.Rollback()

	handler := NewSaleV2Handler(tx, cfg)
	router := setupSaleV2TestRouter(handler)

	// 1. Store Product
	var createdID uint
	t.Run("Store new product with sizes", func(t *testing.T) {
		sizeVal := 82.0
		payload := map[string]interface{}{
			"type":        "BLAZER",
			"price":       450000,
			"description": "Blazer da ca tinh",
			"sizes": []map[string]interface{}{
				{"name": "rong_vai", "value": sizeVal},
			},
		}

		w := performRequest(router, "POST", "/salev2/api/product/store", payload)
		if w.Code != http.StatusOK {
			t.Fatalf("StoreProduct failed with status %d: %s", w.Code, w.Body.String())
		}
		var res struct {
			Success bool `json:"success"`
			Product struct {
				ID    uint   `json:"id"`
				Name  string `json:"name"`
				Price int64  `json:"price"`
			} `json:"product"`
		}
		json.Unmarshal(w.Body.Bytes(), &res)
		if !res.Success || res.Product.ID == 0 {
			t.Fatalf("Expected created product, got %+v", res)
		}
		createdID = res.Product.ID

		// Verify size inserted in DB
		var sz models.Size
		if err := tx.Where("product_id = ? AND name = ?", createdID, "rong_vai").First(&sz).Error; err != nil {
			t.Errorf("Expected size to be saved: %v", err)
		}
	})

	if createdID == 0 {
		t.Fatalf("Cannot continue without created product")
	}

	// 2. Product Detail
	t.Run("Get product detail", func(t *testing.T) {
		w := performRequest(router, "GET", fmt.Sprintf("/salev2/api/product/%d", createdID), nil)
		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d", w.Code)
		}
		var res models.Product
		json.Unmarshal(w.Body.Bytes(), &res)
		if res.ID != createdID {
			t.Errorf("Expected ID %d, got %d", createdID, res.ID)
		}
	})

	// 3. Update Product
	t.Run("Update product details and price", func(t *testing.T) {
		newPrice := int64(520000)
		updatePayload := map[string]interface{}{
			"price":       newPrice,
			"description": "Updated description",
			"status":      "AVAILABLE",
		}
		w := performRequest(router, "POST", fmt.Sprintf("/salev2/api/product/%d/update", createdID), updatePayload)
		if w.Code != http.StatusOK {
			t.Fatalf("UpdateProduct failed: %d, %s", w.Code, w.Body.String())
		}

		var check models.Product
		tx.First(&check, createdID)
		if check.Price == nil || *check.Price != newPrice {
			t.Errorf("Expected price to be updated to %d, got %v", newPrice, check.Price)
		}
	})

	// 4. Destroy Product
	t.Run("Delete product", func(t *testing.T) {
		w := performRequest(router, "DELETE", fmt.Sprintf("/salev2/api/product/%d", createdID), nil)
		if w.Code != http.StatusOK {
			t.Fatalf("DestroyProduct failed: %d", w.Code)
		}

		var check models.Product
		err := tx.First(&check, createdID).Error
		if err == nil {
			t.Errorf("Expected product %d to be deleted", createdID)
		}
	})
}
