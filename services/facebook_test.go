package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"sale-app/config"
	"sale-app/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	testDBOnce sync.Once
	sharedDB   *gorm.DB
	sharedCfg  *config.Config
)

func getTestDB(t *testing.T) (*gorm.DB, *config.Config) {
	testDBOnce.Do(func() {
		sharedCfg = config.LoadConfig()
		db, err := gorm.Open(mysql.Open(sharedCfg.DatabaseDSN), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			t.Fatalf("failed to connect to test db: %v", err)
		}
		_ = db.AutoMigrate(&models.Category{}, &models.AlbumResyncJob{})
		sharedDB = db
	})
	return sharedDB, sharedCfg
}

func strPtr(s string) *string {
	return &s
}

func floatPtr(f float64) *float64 {
	return &f
}

func int64Ptr(i int64) *int64 {
	return &i
}

func TestGetTrouserSize(t *testing.T) {
	tests := []struct {
		eoActual int
		expected *int
	}{
		{73, nil},
		{74, intPtr(28)},
		{75, intPtr(28)},
		{76, intPtr(29)},
		{77, intPtr(29)},
		{78, intPtr(30)},
		{83, intPtr(32)},
		{84, intPtr(33)},
		{98, intPtr(40)},
		{99, intPtr(40)},
		{100, intPtr(40)}, // inclusive boundary
		{101, nil},
	}

	for _, tt := range tests {
		actual := GetTrouserSize(tt.eoActual)
		if tt.expected == nil {
			if actual != nil {
				t.Errorf("GetTrouserSize(%d) = %v; want nil", tt.eoActual, *actual)
			}
		} else {
			if actual == nil || *actual != *tt.expected {
				t.Errorf("GetTrouserSize(%d) = %v; want %d", tt.eoActual, actual, *tt.expected)
			}
		}
	}
}

func intPtr(i int) *int {
	return &i
}

func TestFormatVNDPrice(t *testing.T) {
	tests := []struct {
		amount   int64
		expected string
	}{
		{0, "0đ"},
		{500, "500đ"},
		{1000, "1.000đ"},
		{180000, "180.000đ"},
		{1250000, "1.250.000đ"},
	}

	for _, tt := range tests {
		actual := FormatVNDPrice(tt.amount)
		if actual != tt.expected {
			t.Errorf("FormatVNDPrice(%d) = %s; want %s", tt.amount, actual, tt.expected)
		}
	}
}

func TestBuildHashtags_Trousers(t *testing.T) {
	fb := &FacebookService{}

	product := &models.Product{
		Name: strPtr("TS1235"),
		Type: strPtr("TROUSERS"),
		Sizes: []models.Size{
			{Name: strPtr("EO"), Value: floatPtr(39)}, // eoActual = 78 -> Size 30
		},
	}

	tagsAvailable := fb.BuildHashtags(product, "AVAILABLE")
	expectedAvailable := "#TS1235 #TROUSERS_AVAILABLE #EO_78 #EO_78_80 #SIZE30"
	if tagsAvailable != expectedAvailable {
		t.Errorf("tagsAvailable = %q; want %q", tagsAvailable, expectedAvailable)
	}

	tagsSold := fb.BuildHashtags(product, "SOLD")
	expectedSold := "#TS1235 #TROUSERS_SOLD #EO_78 #EO_78_80 #SIZE30"
	if tagsSold != expectedSold {
		t.Errorf("tagsSold = %q; want %q", tagsSold, expectedSold)
	}
}

func TestBuildHashtags_BlazerAndShirt(t *testing.T) {
	fb := &FacebookService{}

	// Even vai (44)
	blazer := &models.Product{
		Name: strPtr("BZ456"),
		Type: strPtr("BLAZER"),
		Sizes: []models.Size{
			{Name: strPtr("VAI"), Value: floatPtr(44)},
		},
	}
	tagsBlazer := fb.BuildHashtags(blazer, "AVAILABLE")
	expectedBlazer := "#BZ456 #BLAZER_AVAILABLE #VAI_44 #VAI_44_45"
	if tagsBlazer != expectedBlazer {
		t.Errorf("tagsBlazer = %q; want %q", tagsBlazer, expectedBlazer)
	}

	// Odd vai (43) -> start = 42, range 42_43
	shirt := &models.Product{
		Name: strPtr("SM101"),
		Type: strPtr("SHIRT"),
		Sizes: []models.Size{
			{Name: strPtr("VAI"), Value: floatPtr(43)},
		},
	}
	tagsShirt := fb.BuildHashtags(shirt, "AVAILABLE")
	expectedShirt := "#SM101 #SHIRT_AVAILABLE #VAI_43 #VAI_42_43"
	if tagsShirt != expectedShirt {
		t.Errorf("tagsShirt = %q; want %q", tagsShirt, expectedShirt)
	}
}

func TestBuildMessage(t *testing.T) {
	db, _ := getTestDB(t)
	tx := db.Begin()
	defer tx.Rollback()

	pName := fmt.Sprintf("TS_%d", time.Now().UnixNano()%100000)
	product := models.Product{
		Name:  strPtr(pName),
		Type:  strPtr("TROUSERS"),
		Price: int64Ptr(180000),
		Sizes: []models.Size{
			{Name: strPtr("EO"), Value: floatPtr(40)}, // eoActual = 80 -> Size 31
		},
	}
	tx.Create(&product)

	prodID := int(product.ID)
	log := models.ProductUploadLog{
		ProductID:    &prodID,
		DetectedText: strPtr("EO 80 DAI 100 ONG 20"),
	}
	tx.Create(&log)

	fb := &FacebookService{DB: tx}

	msgAvailable := fb.BuildMessage(&product, "AVAILABLE")
	if !strings.Contains(msgAvailable, fmt.Sprintf("✅ CÒN HÀNG — %s", pName)) {
		t.Errorf("Expected status line with %s, got: %s", pName, msgAvailable)
	}
	if !strings.Contains(msgAvailable, "👔 TROUSERS | 180.000đ") {
		t.Errorf("Expected price and type, got: %s", msgAvailable)
	}
	if !strings.Contains(msgAvailable, "EO 80 DAI 100 ONG 20") {
		t.Errorf("Expected detected text, got: %s", msgAvailable)
	}
	if !strings.Contains(msgAvailable, fmt.Sprintf("#%s #TROUSERS_AVAILABLE #EO_80 #EO_80_82 #SIZE31", pName)) {
		t.Errorf("Expected hashtags, got: %s", msgAvailable)
	}

	msgSold := fb.BuildMessage(&product, "SOLD")
	if !strings.Contains(msgSold, fmt.Sprintf("⛔ ĐÃ BÁN — %s", pName)) {
		t.Errorf("Expected sold status line, got: %s", msgSold)
	}
	if !strings.Contains(msgSold, "#TROUSERS_SOLD") {
		t.Errorf("Expected sold hashtag, got: %s", msgSold)
	}
}

func TestFacebookService_FullSyncLifecycle(t *testing.T) {
	// 1. Mock Facebook Graph API Server
	var uploadedPhotos []string
	var feedMessages []string
	var updatedMessages []string
	var deletedPosts []string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if strings.HasSuffix(path, "/photos") {
			// Photo upload
			uploadedPhotos = append(uploadedPhotos, path)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]string{"id": fmt.Sprintf("media_%d", len(uploadedPhotos))})
			return
		}

		if strings.HasSuffix(path, "/feed") {
			// Feed post create
			_ = r.ParseForm()
			feedMessages = append(feedMessages, r.FormValue("message"))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]string{"id": "fb_post_12345"})
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
			// Update message (mark SOLD / revert)
			_ = r.ParseForm()
			updatedMessages = append(updatedMessages, r.FormValue("message"))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]bool{"success": true})
			return
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	// 2. Setup DB and temporary test images
	db, _ := getTestDB(t)
	tx := db.Begin()
	defer tx.Rollback()

	tmpDir := t.TempDir()
	img1 := filepath.Join(tmpDir, "thumb.jpg")
	_ = os.WriteFile(img1, []byte("fake image data 1"), 0644)
	img2 := filepath.Join(tmpDir, "detail.jpg")
	_ = os.WriteFile(img2, []byte("fake image data 2"), 0644)

	pName := fmt.Sprintf("TS_%d", time.Now().UnixNano()%100000)
	product := models.Product{
		Name:      strPtr(pName),
		Type:      strPtr("TROUSERS"),
		Price:     int64Ptr(200000),
		PathThumb: strPtr(img1),
		Sizes: []models.Size{
			{Name: strPtr("EO"), Value: floatPtr(37)}, // eoActual = 74 -> Size 28
		},
	}
	tx.Create(&product)

	tx.Create(&models.ProductImage{
		ProductID: product.ID,
		FilePath:  img2,
	})

	fb := &FacebookService{
		PageID:    "123456789",
		PageToken: "EAAB_TEST_TOKEN",
		BaseURL:   server.URL + "/",
		Client:    server.Client(),
		DB:        tx,
		Queue:     make(chan FacebookSyncJob, 10),
	}

	// 3. Test Action: CREATE
	err := fb.ExecuteSync(product.ID, FBActionCreate)
	if err != nil {
		t.Fatalf("ExecuteSync(CREATE) failed: %v", err)
	}

	// 2 photos for catalog post + 1 photo for category album = 3 photos total
	if len(uploadedPhotos) < 2 {
		t.Errorf("Expected at least 2 photos uploaded, got %d", len(uploadedPhotos))
	}
	if len(feedMessages) != 1 {
		t.Errorf("Expected 1 feed message, got %d", len(feedMessages))
	}

	// Verify DB updated with fb_post_id
	var updatedProd models.Product
	tx.First(&updatedProd, product.ID)
	if updatedProd.FBPostID == nil || *updatedProd.FBPostID != "fb_post_12345" {
		t.Fatalf("Expected fb_post_id to be 'fb_post_12345', got %v", updatedProd.FBPostID)
	}

	// 4. Test Action: SOLD
	err = fb.ExecuteSync(product.ID, FBActionSold)
	if err != nil {
		t.Fatalf("ExecuteSync(SOLD) failed: %v", err)
	}
	if len(updatedMessages) != 1 || !strings.Contains(updatedMessages[0], fmt.Sprintf("⛔ ĐÃ BÁN — %s", pName)) {
		t.Errorf("Expected sold message update, got %v", updatedMessages)
	}

	// 5. Test Action: REVERT
	err = fb.ExecuteSync(product.ID, FBActionRevert)
	if err != nil {
		t.Fatalf("ExecuteSync(REVERT) failed: %v", err)
	}
	if len(updatedMessages) != 2 || !strings.Contains(updatedMessages[1], fmt.Sprintf("✅ CÒN HÀNG — %s", pName)) {
		t.Errorf("Expected revert message update, got %v", updatedMessages)
	}

	// 6. Test Re-create: Should delete old post before creating new one
	err = fb.ExecuteSync(product.ID, FBActionCreate)
	if err != nil {
		t.Fatalf("Re-creating post failed: %v", err)
	}
	// Old catalog post and previous album photo deleted before re-creating
	if len(deletedPosts) < 1 {
		t.Errorf("Expected delete calls on re-create, got %d", len(deletedPosts))
	}
}

func TestBuildAlbumPhotoCaption(t *testing.T) {
	fb := &FacebookService{
		PageID: "1132727476600909",
	}

	pName := "TS8888"
	pPrice := int64(250000)
	prod := &models.Product{
		Name:  &pName,
		Type:  strPtr("TROUSERS"),
		Price: &pPrice,
		Sizes: []models.Size{
			{Name: strPtr("EO"), Value: floatPtr(40)},
			{Name: strPtr("DÀI"), Value: floatPtr(98)},
		},
	}

	catalogPostID := "post_abc_123"
	caption := fb.BuildAlbumPhotoCaption(prod, catalogPostID)

	if !strings.Contains(caption, "✅ CÒN HÀNG — TS8888") {
		t.Errorf("Caption missing status/name line, got:\n%s", caption)
	}
	if !strings.Contains(caption, "👔 TROUSERS | Giá: 250.000đ") {
		t.Errorf("Caption missing type/price, got:\n%s", caption)
	}
	if !strings.Contains(caption, "EO: 40") || !strings.Contains(caption, "DÀI: 98") {
		t.Errorf("Caption missing sizes, got:\n%s", caption)
	}
	if !strings.Contains(caption, "👉 Xem ảnh chi tiết tại:") {
		t.Errorf("Caption missing link prefix, got:\n%s", caption)
	}
	expectedURL := "https://www.facebook.com/1132727476600909/posts/post_abc_123"
	if !strings.Contains(caption, expectedURL) {
		t.Errorf("Caption missing catalog post URL %s, got:\n%s", expectedURL, caption)
	}
	if !strings.Contains(caption, "#TS8888") || !strings.Contains(caption, "#TROUSERS_AVAILABLE") {
		t.Errorf("Caption missing hashtags, got:\n%s", caption)
	}
}

func TestAlbumPhoto_UploadAndRemoveAPI(t *testing.T) {
	var receivedNoStory string
	var receivedCaption string
	var receivedDeletePhotoID string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Upload photo to album
		if r.Method == "POST" && strings.Contains(r.URL.Path, "/album_999/photos") {
			_ = r.ParseForm()
			receivedNoStory = r.FormValue("no_story")
			receivedCaption = r.FormValue("caption")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"id": "photo_album_111",
			})
			return
		}

		// Delete photo from album
		if r.Method == "DELETE" && strings.Contains(r.URL.Path, "/photo_album_111") {
			parts := strings.Split(r.URL.Path, "/")
			receivedDeletePhotoID = parts[len(parts)-1]
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
			})
			return
		}

		// Delete already deleted photo (test graceful error handling)
		if r.Method == "DELETE" && strings.Contains(r.URL.Path, "/already_deleted_photo") {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]interface{}{
					"message":       "Unsupported delete request. Object with ID 'already_deleted_photo' does not exist.",
					"type":          "GraphMethodException",
					"code":          100,
					"error_subcode": 33,
				},
			})
			return
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	fb := &FacebookService{
		PageID:        "1132727476600909",
		PageToken:     "token_catalog",
		MainPageID:    "236648362875702",
		MainPageToken: "token_main_page",
		BaseURL:       server.URL + "/",
		Client:        server.Client(),
	}

	// 1. Test Upload to Album with no_story=true
	testPhotoURL := "http://example.com/test-thumb.jpg"
	testCaption := "Test Album Caption"
	photoID, err := fb.UploadProductPhotoToAlbum("album_999", testPhotoURL, testCaption)
	if err != nil {
		t.Fatalf("UploadProductPhotoToAlbum failed: %v", err)
	}
	if photoID != "photo_album_111" {
		t.Errorf("Expected photoID 'photo_album_111', got '%s'", photoID)
	}
	if receivedNoStory != "true" {
		t.Errorf("Expected no_story='true', got '%s'", receivedNoStory)
	}
	if receivedCaption != testCaption {
		t.Errorf("Expected caption '%s', got '%s'", testCaption, receivedCaption)
	}

	// 2. Test Remove photo from album
	err = fb.RemoveProductPhotoFromAlbum("photo_album_111")
	if err != nil {
		t.Fatalf("RemoveProductPhotoFromAlbum failed: %v", err)
	}
	if receivedDeletePhotoID != "photo_album_111" {
		t.Errorf("Expected deleted photo ID 'photo_album_111', got '%s'", receivedDeletePhotoID)
	}

	// 3. Test Remove already deleted photo (should not return error)
	err = fb.RemoveProductPhotoFromAlbum("already_deleted_photo")
	if err != nil {
		t.Fatalf("RemoveProductPhotoFromAlbum for already deleted object should succeed gracefully, got: %v", err)
	}
}

func TestExecuteSync_AlbumWorkflow(t *testing.T) {
	db, _ := getTestDB(t)
	tx := db.Begin()
	defer tx.Rollback()

	// Setup Category with fb_album_id
	catCode := fmt.Sprintf("TEST_CAT_%d", time.Now().UnixNano())
	catAlbumID := "album_target_777"
	cat := models.Category{
		Code:      catCode,
		Name:      "Test Category",
		FBAlbumID: &catAlbumID,
		IsActive:  true,
	}
	if err := tx.Create(&cat).Error; err != nil {
		t.Fatalf("Failed to create test category: %v", err)
	}

	// Create test product
	pName := fmt.Sprintf("TS_ALBUM_%d", time.Now().UnixNano()%100000)
	pPrice := int64(300000)
	pThumb := "http://example.com/test-thumb.jpg"
	prod := models.Product{
		Name:      &pName,
		Type:      &catCode,
		Price:     &pPrice,
		PathThumb: &pThumb,
		Status:    strPtr("AVAILABLE"),
	}
	if err := tx.Create(&prod).Error; err != nil {
		t.Fatalf("Failed to create test product: %v", err)
	}

	var albumUploadedWithNoStory bool
	var albumPhotoDeleted bool

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// 1. Photo upload for Catalog Post (feed attached_media)
		if r.Method == "POST" && strings.Contains(r.URL.Path, "/1132727476600909/photos") {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"id": "feed_photo_1"})
			return
		}
		// 2. Feed post creation on Page B
		if r.Method == "POST" && strings.Contains(r.URL.Path, "/1132727476600909/feed") {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"id": "catalog_post_999"})
			return
		}
		// 3. Album photo upload on Page A
		if r.Method == "POST" && strings.Contains(r.URL.Path, "/album_target_777/photos") {
			_ = r.ParseForm()
			if r.FormValue("no_story") == "true" {
				albumUploadedWithNoStory = true
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"id": "album_photo_888"})
			return
		}
		// 4. Update post on Page B (MarkSold)
		if r.Method == "POST" && strings.Contains(r.URL.Path, "/catalog_post_999") {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
			return
		}
		// 5. Delete album photo on Page A
		if r.Method == "DELETE" && strings.Contains(r.URL.Path, "/album_photo_888") {
			albumPhotoDeleted = true
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
			return
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	fb := &FacebookService{
		PageID:        "1132727476600909",
		PageToken:     "token_catalog",
		MainPageID:    "236648362875702",
		MainPageToken: "token_main",
		BaseURL:       server.URL + "/",
		Client:        server.Client(),
		DB:            tx,
	}

	// 1. ExecuteSync(FBActionCreate)
	err := fb.ExecuteSync(prod.ID, FBActionCreate)
	if err != nil {
		t.Fatalf("ExecuteSync(FBActionCreate) failed: %v", err)
	}

	if !albumUploadedWithNoStory {
		t.Errorf("Expected album photo to be uploaded with no_story=true")
	}

	var createdProd models.Product
	tx.First(&createdProd, prod.ID)
	if createdProd.FBPostID == nil || *createdProd.FBPostID != "catalog_post_999" {
		t.Fatalf("Expected fb_post_id 'catalog_post_999', got %v", createdProd.FBPostID)
	}
	if createdProd.FBAlbumPhotoID == nil || *createdProd.FBAlbumPhotoID != "album_photo_888" {
		t.Fatalf("Expected fb_album_photo_id 'album_photo_888', got %v", createdProd.FBAlbumPhotoID)
	}

	// 2. ExecuteSync(FBActionSold)
	err = fb.ExecuteSync(prod.ID, FBActionSold)
	if err != nil {
		t.Fatalf("ExecuteSync(FBActionSold) failed: %v", err)
	}

	if !albumPhotoDeleted {
		t.Errorf("Expected album photo to be deleted via DELETE API")
	}

	var soldProd models.Product
	tx.First(&soldProd, prod.ID)
	if soldProd.FBAlbumPhotoID != nil {
		t.Fatalf("Expected fb_album_photo_id to be NULL after SOLD, got %v", *soldProd.FBAlbumPhotoID)
	}
}

func TestExecuteSync_HoldAndReconcile(t *testing.T) {
	db, _ := getTestDB(t)
	if db == nil {
		t.Skip("Database not available, skipping TestExecuteSync_HoldAndReconcile")
	}

	deletedPhotoIDs := make(map[string]bool)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			parts := strings.Split(r.URL.Path, "/")
			if len(parts) > 1 {
				photoID := parts[len(parts)-1]
				deletedPhotoIDs[photoID] = true
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"success":true}`))
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	tx := db.Begin()
	defer tx.Rollback()

	photoID := "hold_photo_123"
	holdProd := models.Product{
		Name:            strPtr("GL_HOLD_TEST"),
		Status:          strPtr("ON_HOLD"),
		FBAlbumPhotoID:  &photoID,
	}
	tx.Create(&holdProd)

	fb := &FacebookService{
		BaseURL:       server.URL + "/",
		PageToken:     "test_token",
		MainPageToken: "main_token",
		Client:        server.Client(),
		DB:            tx,
		ThrottleDelay: 1 * time.Millisecond,
	}

	// 1. Test ExecuteSync(FBActionHold)
	err := fb.ExecuteSync(holdProd.ID, FBActionHold)
	if err != nil {
		t.Fatalf("ExecuteSync(FBActionHold) failed: %v", err)
	}
	if !deletedPhotoIDs["hold_photo_123"] {
		t.Errorf("Expected hold_photo_123 to be deleted via DELETE API")
	}

	var updated models.Product
	tx.First(&updated, holdProd.ID)
	if updated.FBAlbumPhotoID != nil {
		t.Errorf("Expected fb_album_photo_id to be NULL after HOLD, got %v", *updated.FBAlbumPhotoID)
	}

	// 2. Test ReconcileAlbumPhotos
	stalePhotoID := "stale_photo_456"
	staleProd := models.Product{
		Name:           strPtr("GL_STALE_TEST"),
		Status:         strPtr("SOLD"),
		FBAlbumPhotoID: &stalePhotoID,
	}
	tx.Create(&staleProd)

	fb.ReconcileAlbumPhotos()
	if !deletedPhotoIDs["stale_photo_456"] {
		t.Errorf("Expected stale_photo_456 to be deleted by ReconcileAlbumPhotos")
	}

	var reconciled models.Product
	tx.First(&reconciled, staleProd.ID)
	if reconciled.FBAlbumPhotoID != nil {
		t.Errorf("Expected stale product fb_album_photo_id to be NULL after reconcile, got %v", *reconciled.FBAlbumPhotoID)
	}

	// 3. Test isFBRateLimitError
	rateLimitErr := fmt.Errorf("facebook delete photo error (403): {\"error\":{\"message\":\"(#4) Application request limit reached\"}}")
	if !isFBRateLimitError(rateLimitErr) {
		t.Errorf("Expected isFBRateLimitError to be true for (#4) error")
	}
	otherErr := fmt.Errorf("network timeout")
	if isFBRateLimitError(otherErr) {
		t.Errorf("Expected isFBRateLimitError to be false for timeout error")
	}
}

func TestAlbumResync_RunPauseRun(t *testing.T) {
	db, _ := getTestDB(t)
	if db == nil {
		t.Skip("Database not available")
	}

	uploadCount := 0
	var uploadedCaptions []string
	deletedPhotos := make(map[string]bool)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "GET" && strings.Contains(r.URL.Path, "/photos") {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"data":[{"id":"old_photo_1"}]}`))
			return
		}
		if r.Method == "DELETE" {
			parts := strings.Split(r.URL.Path, "/")
			deletedPhotos[parts[len(parts)-1]] = true
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"success":true}`))
			return
		}
		if r.Method == "POST" && strings.Contains(r.URL.Path, "/photos") {
			_ = r.ParseForm()
			uploadCount++
			uploadedCaptions = append(uploadedCaptions, r.FormValue("caption"))
			// First upload triggers proactive pause via 88% usage header!
			if uploadCount == 1 {
				w.Header().Set("x-app-usage", `{"call_count":10,"total_cputime":0,"total_time":88}`)
			} else {
				w.Header().Set("x-app-usage", `{"call_count":12,"total_cputime":0,"total_time":20}`)
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(fmt.Sprintf(`{"id":"new_photo_%d"}`, uploadCount)))
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	tx := db.Begin()
	defer tx.Rollback()

	catAlbumID := "album_resync_test_123"
	cat := models.Category{
		Code:      "TEST_RESYNC_CAT",
		Name:      "Test Resync Category",
		FBAlbumID: &catAlbumID,
		IsActive:  true,
	}
	tx.Create(&cat)

	t1 := time.Now().Add(-2 * time.Hour)
	t2 := time.Now().Add(-1 * time.Hour)
	pOld := models.Product{
		Name:      strPtr("TEST_OLD"),
		Type:      strPtr("TEST_RESYNC_CAT"),
		Status:    strPtr("AVAILABLE"),
		PathThumb: strPtr("https://example.com/old.jpg"),
		CreatedAt: &t1,
	}
	pNew := models.Product{
		Name:      strPtr("TEST_NEW"),
		Type:      strPtr("TEST_RESYNC_CAT"),
		Status:    strPtr("AVAILABLE"),
		PathThumb: strPtr("https://example.com/new.jpg"),
		CreatedAt: &t2,
	}
	tx.Create(&pOld)
	tx.Create(&pNew)

	fb := &FacebookService{
		BaseURL:       server.URL + "/",
		PageToken:     "test_token",
		MainPageToken: "main_token",
		Client:        server.Client(),
		DB:            tx,
		ThrottleDelay: 1 * time.Millisecond,
	}

	job, err := fb.StartCategoryAlbumResync("TEST_RESYNC_CAT", "MANUAL")
	if err != nil {
		t.Fatalf("StartCategoryAlbumResync failed: %v", err)
	}
	if job.TotalDelete != 1 || job.TotalUpload != 2 {
		t.Fatalf("Expected 1 delete and 2 uploads, got del=%d, up=%d", job.TotalDelete, job.TotalUpload)
	}

	// First run: deletes old_photo_1, uploads newest (TEST_NEW), then proactively pauses due to 88% usage!
	fb.ProcessNextAlbumResyncJob()

	var pausedJob models.AlbumResyncJob
	tx.First(&pausedJob, job.ID)
	if pausedJob.Status != "PAUSED_LIMIT" {
		t.Fatalf("Expected job status PAUSED_LIMIT after 88%% usage, got %s", pausedJob.Status)
	}
	if pausedJob.DeletedCount != 1 || pausedJob.UploadedCount != 1 {
		t.Fatalf("Expected DeletedCount=1, UploadedCount=1, got del=%d, up=%d", pausedJob.DeletedCount, pausedJob.UploadedCount)
	}
	if !deletedPhotos["old_photo_1"] {
		t.Errorf("Expected old_photo_1 to be deleted")
	}

	// Force resume (simulates waiting 45 mins or clicking Resume)
	if err := fb.ForceResumeAlbumResyncJob(job.ID); err != nil {
		t.Fatalf("ForceResumeAlbumResyncJob failed: %v", err)
	}

	// Second run: uploads TEST_OLD and marks COMPLETED
	fb.ProcessNextAlbumResyncJob()

	var completedJob models.AlbumResyncJob
	tx.First(&completedJob, job.ID)
	if completedJob.Status != "COMPLETED" {
		t.Fatalf("Expected job status COMPLETED, got %s (last_error=%s)", completedJob.Status, completedJob.LastError)
	}
	if completedJob.UploadedCount != 2 {
		t.Fatalf("Expected UploadedCount=2, got %d", completedJob.UploadedCount)
	}
	if len(uploadedCaptions) != 2 || !strings.Contains(uploadedCaptions[0], "TEST_NEW") || !strings.Contains(uploadedCaptions[1], "TEST_OLD") {
		t.Errorf("Expected upload order [TEST_NEW, TEST_OLD], got %v", uploadedCaptions)
	}
}
