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

	if len(uploadedPhotos) != 2 {
		t.Errorf("Expected 2 photos uploaded, got %d", len(uploadedPhotos))
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
	if len(deletedPosts) != 1 {
		t.Errorf("Expected 1 delete post call, got %d", len(deletedPosts))
	}
}
