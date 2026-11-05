package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"sale-app/models"
)

func TestSaleV2_Products_SortByUpload(t *testing.T) {
	db, cfg := getTestDB(t)
	tx := db.Begin()
	defer tx.Rollback()

	handler := NewSaleV2Handler(tx, cfg)
	router := setupSaleV2TestRouter(handler)

	// Create 2 products
	name1 := fmt.Sprintf("UPLOAD_TEST_A_%d", time.Now().UnixNano()%100000)
	name2 := fmt.Sprintf("UPLOAD_TEST_B_%d", time.Now().UnixNano()%100000)
	price := int64(100000)

	prodA := models.Product{
		Name:      &name1,
		Price:     &price,
		Status:    strPtr("AVAILABLE"),
		PathThumb: strPtr("https://res.cloudinary.com/test/prodA.jpg"),
	}
	prodB := models.Product{
		Name:      &name2,
		Price:     &price,
		Status:    strPtr("AVAILABLE"),
		PathThumb: strPtr("https://res.cloudinary.com/test/prodB.jpg"),
	}
	tx.Create(&prodA)
	// Give a slight delay or create prodB later so prodB.ID > prodA.ID
	tx.Create(&prodB)

	// Upload log for prodA is NEWER (created_at in future or now)
	uploadTime := time.Now().Add(1 * time.Hour)
	prodAID := int(prodA.ID)
	logA := models.ProductUploadLog{
		SessionID: "sess_test_123",
		FileName:  "test.jpg",
		FilePath:  "test/path.jpg",
		ProductID: &prodAID,
		Status:    "SUCCESS",
		CreatedAt: &uploadTime,
	}
	tx.Create(&logA)

	t.Run("Default sort (ID DESC) returns ProdB before ProdA", func(t *testing.T) {
		w := performRequest(router, "GET", fmt.Sprintf("/salev2/api/products?search=%s,%s", name1, name2), nil)
		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d: %s", w.Code, w.Body.String())
		}
		var res struct {
			Data []models.Product `json:"data"`
		}
		json.Unmarshal(w.Body.Bytes(), &res)
		if len(res.Data) < 2 {
			t.Fatalf("Expected at least 2 products, got %d", len(res.Data))
		}
		if res.Data[0].ID != prodB.ID {
			t.Errorf("Expected first product to be B (%d), got %d", prodB.ID, res.Data[0].ID)
		}
	})

	t.Run("Sort by upload returns ProdA before ProdB and has latest_upload", func(t *testing.T) {
		w := performRequest(router, "GET", fmt.Sprintf("/salev2/api/products?sort=upload&search=%s,%s", name1, name2), nil)
		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d: %s", w.Code, w.Body.String())
		}
		var res struct {
			Data []models.Product `json:"data"`
		}
		json.Unmarshal(w.Body.Bytes(), &res)
		if len(res.Data) < 2 {
			t.Fatalf("Expected at least 2 products, got %d", len(res.Data))
		}
		if res.Data[0].ID != prodA.ID {
			t.Errorf("Expected first product to be A (%d) because it was uploaded later, got %d", prodA.ID, res.Data[0].ID)
		}
		if res.Data[0].LatestUpload == nil {
			t.Errorf("Expected ProdA to have latest_upload not nil")
		}
	})
}
