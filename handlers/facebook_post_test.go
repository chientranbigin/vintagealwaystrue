package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFacebookPostInitAndContent(t *testing.T) {
	db, cfg := getTestDB(t)
	v2 := NewSaleV2Handler(db, cfg)

	router := setupSaleV2TestRouter(v2)
	router.GET("/salev2/api/facebook/post-init", v2.FacebookPostInit)
	router.GET("/salev2/api/facebook/category-products", v2.FacebookCategoryProducts)
	router.POST("/salev2/api/facebook/generate-content", v2.FacebookGenerateContent)
	router.POST("/salev2/api/facebook/publish-main-post", v2.FacebookPublishMainPost)

	// 1. Test FacebookPostInit
	reqInit, _ := http.NewRequest("GET", "/salev2/api/facebook/post-init", nil)
	wInit := httptest.NewRecorder()
	router.ServeHTTP(wInit, reqInit)

	if wInit.Code != http.StatusOK {
		t.Fatalf("FacebookPostInit failed with code %d: %s", wInit.Code, wInit.Body.String())
	}

	var initResp struct {
		Categories []struct {
			Code           string `json:"code"`
			Name           string `json:"name"`
			AvailableCount int64  `json:"available_count"`
		} `json:"categories"`
		SuggestedEvent string `json:"suggested_event"`
		EndCardURL     string `json:"end_card_url"`
	}
	if err := json.Unmarshal(wInit.Body.Bytes(), &initResp); err != nil {
		t.Fatalf("Failed to parse FacebookPostInit response: %v", err)
	}

	if len(initResp.Categories) == 0 {
		t.Errorf("Expected categories to not be empty")
	}
	if initResp.SuggestedEvent == "" {
		t.Errorf("Expected suggested event to not be empty")
	}
	if initResp.EndCardURL == "" {
		t.Errorf("Expected end card URL to not be empty")
	}

	// 2. Test FacebookGenerateContent with Holiday ON
	genBodyHoliday, _ := json.Marshal(map[string]interface{}{
		"category_code": "GILE",
		"category_name": "Áo Gile",
		"album_url":     "https://www.facebook.com/media/set/?set=a.122252359550280390",
		"include_event": true,
		"event_name":    "Giáng Sinh / Noel",
	})
	reqGen1, _ := http.NewRequest("POST", "/salev2/api/facebook/generate-content", bytes.NewBuffer(genBodyHoliday))
	reqGen1.Header.Set("Content-Type", "application/json")
	wGen1 := httptest.NewRecorder()
	router.ServeHTTP(wGen1, reqGen1)

	if wGen1.Code != http.StatusOK {
		t.Fatalf("FacebookGenerateContent (Holiday ON) failed with code %d: %s", wGen1.Code, wGen1.Body.String())
	}

	var genResp1 struct {
		Content string `json:"content"`
	}
	if err := json.Unmarshal(wGen1.Body.Bytes(), &genResp1); err != nil {
		t.Fatalf("Failed to parse generate content response: %v", err)
	}
	if genResp1.Content == "" {
		t.Errorf("Expected generated content to not be empty")
	}

	// 3. Test FacebookGenerateContent with Holiday OFF (Standard)
	genBodyStd, _ := json.Marshal(map[string]interface{}{
		"category_code": "TROUSERS",
		"category_name": "Quần Âu",
		"album_url":     "https://www.facebook.com/media/set/?set=a.122252355908280390",
		"include_event": false,
	})
	reqGen2, _ := http.NewRequest("POST", "/salev2/api/facebook/generate-content", bytes.NewBuffer(genBodyStd))
	reqGen2.Header.Set("Content-Type", "application/json")
	wGen2 := httptest.NewRecorder()
	router.ServeHTTP(wGen2, reqGen2)

	if wGen2.Code != http.StatusOK {
		t.Fatalf("FacebookGenerateContent (Holiday OFF) failed with code %d: %s", wGen2.Code, wGen2.Body.String())
	}

	// 4. Test validation of publish endpoint (rejects empty product_ids)
	pubEmpty, _ := json.Marshal(map[string]interface{}{
		"message":     "Test message",
		"product_ids": []uint{},
	})
	reqPub, _ := http.NewRequest("POST", "/salev2/api/facebook/publish-main-post", bytes.NewBuffer(pubEmpty))
	reqPub.Header.Set("Content-Type", "application/json")
	wPub := httptest.NewRecorder()
	router.ServeHTTP(wPub, reqPub)

	if wPub.Code != http.StatusBadRequest {
		t.Errorf("Expected Bad Request for empty products, got %d", wPub.Code)
	}
}
