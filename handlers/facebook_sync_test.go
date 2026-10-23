package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"sale-app/models"
	"sale-app/services"
)

func TestSaleV2_FacebookSync_UploadAndOrderLifecycle(t *testing.T) {
	db, cfg := getTestDB(t)
	tx := db.Begin()
	defer tx.Rollback()

	// 1. Mock Facebook Graph API
	var mu sync.Mutex
	var uploadedPhotos []string
	var feedPosts []string
	var updatedMessages []string
	var deletedPosts []string

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		path := r.URL.Path

		if r.Method == "POST" && filepath.Base(path) == "photos" {
			uploadedPhotos = append(uploadedPhotos, path)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]string{"id": fmt.Sprintf("media_%d", len(uploadedPhotos))})
			return
		}

		if r.Method == "POST" && filepath.Base(path) == "feed" {
			_ = r.ParseForm()
			feedPosts = append(feedPosts, r.FormValue("message"))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]string{"id": "fb_post_created_999"})
			return
		}

		if r.Method == "DELETE" {
			deletedPosts = append(deletedPosts, path)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]bool{"success": true})
			return
		}

		if r.Method == "POST" {
			_ = r.ParseForm()
			updatedMessages = append(updatedMessages, r.FormValue("message"))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]bool{"success": true})
			return
		}

		http.NotFound(w, r)
	}))
	defer mockServer.Close()

	// 2. Setup handler with FB service pointing to mockServer
	handler := NewSaleV2Handler(tx, cfg)
	handler.FB = &services.FacebookService{
		PageID:    "123456",
		PageToken: "EAAB_TEST_TOKEN",
		BaseURL:   mockServer.URL + "/",
		Client:    mockServer.Client(),
		DB:        tx,
		Queue:     make(chan services.FacebookSyncJob, 20),
	}
	// Run worker synchronously in test
	router := setupSaleV2TestRouter(handler)

	// 3. Test Smart Upload -> Creates FB Post
	pName := fmt.Sprintf("TS%d", time.Now().UnixNano()%100000)
	fVal := 38.0
	prod := models.Product{
		Name:   &pName,
		Type:   strPtr("TROUSERS"),
		Price:  int64Ptr(190000),
		Status: strPtr("AVAILABLE"),
		Sizes: []models.Size{
			{Name: strPtr("EO"), Value: &fVal}, // eoActual = 76 -> size 29
		},
	}
	tx.Create(&prod)

	tmpDir := t.TempDir()
	mainImg := filepath.Join(tmpDir, "main.jpg")
	_ = os.WriteFile(mainImg, []byte("main image data"), 0644)
	prod.PathThumb = &mainImg
	tx.Save(&prod)

	// Mock text detection in DB log or directly trigger ExecuteSync
	err := handler.FB.ExecuteSync(prod.ID, services.FBActionCreate)
	if err != nil {
		t.Fatalf("FB create failed: %v", err)
	}

	mu.Lock()
	if len(feedPosts) == 0 {
		t.Fatalf("Expected feed post to be created on Facebook")
	}
	mu.Unlock()

	var checkProd models.Product
	tx.First(&checkProd, prod.ID)
	if checkProd.FBPostID == nil || *checkProd.FBPostID != "fb_post_created_999" {
		t.Fatalf("Expected fb_post_id to be 'fb_post_created_999', got %v", checkProd.FBPostID)
	}

	// 4. Test Store Order -> Marks FB post as SOLD
	orderPayload := map[string]interface{}{
		"name":         "Customer FB Test",
		"phone":        "0988111222",
		"address":      "123 Nguyen Trai",
		"total_amount": 190000,
		"products":     []uint{prod.ID},
		"status":       "NEW",
	}
	w := performRequest(router, "POST", "/salev2/api/order", orderPayload)
	if w.Code != http.StatusOK {
		t.Fatalf("StoreOrder failed: %s", w.Body.String())
	}

	// Execute sold sync
	err = handler.FB.ExecuteSync(prod.ID, services.FBActionSold)
	if err != nil {
		t.Fatalf("FB sold sync failed: %v", err)
	}

	mu.Lock()
	hasSoldMessage := false
	for _, msg := range updatedMessages {
		if stringContains(msg, "⛔ ĐÃ BÁN") && stringContains(msg, pName) {
			hasSoldMessage = true
			break
		}
	}
	mu.Unlock()

	if !hasSoldMessage {
		t.Errorf("Expected updated message with ⛔ ĐÃ BÁN for %s, got: %v", pName, updatedMessages)
	}

	// 5. Test Destroy Order -> Reverts FB post to AVAILABLE
	var order models.Order
	tx.Where("phone = ?", "0988111222").First(&order)
	if order.ID == 0 {
		t.Fatalf("Order not found")
	}

	wDel := performRequest(router, "DELETE", fmt.Sprintf("/salev2/api/order/%d", order.ID), nil)
	if wDel.Code != http.StatusOK {
		t.Fatalf("DestroyOrder failed: %s", wDel.Body.String())
	}

	// Execute revert sync
	err = handler.FB.ExecuteSync(prod.ID, services.FBActionRevert)
	if err != nil {
		t.Fatalf("FB revert sync failed: %v", err)
	}

	mu.Lock()
	hasRevertMessage := false
	for _, msg := range updatedMessages {
		if stringContains(msg, "✅ CÒN HÀNG") && stringContains(msg, pName) {
			hasRevertMessage = true
			break
		}
	}
	mu.Unlock()

	if !hasRevertMessage {
		t.Errorf("Expected updated message with ✅ CÒN HÀNG for %s, got: %v", pName, updatedMessages)
	}
}

func stringContains(s, substr string) bool {
	return filepath.HasPrefix(s, substr) || len(s) >= len(substr) && (s == substr || len(substr) > 0 && (func() bool {
		for i := 0; i+len(substr) <= len(s); i++ {
			if s[i:i+len(substr)] == substr {
				return true
			}
		}
		return false
	})())
}
