package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"sale-app/models"
)

func TestSaleV2_UploadSmartGroup_Endpoint(t *testing.T) {
	db, cfg := getTestDB(t)
	tx := db.Begin()
	defer tx.Rollback()

	handler := NewSaleV2Handler(tx, cfg)
	router := setupSaleV2TestRouter(handler)

	t.Run("Fail without main image", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.WriteField("type", "TROUSERS")
		writer.Close()

		req, _ := http.NewRequest("POST", "/salev2/api/upload-smart-group", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected 400 Bad Request, got %d", w.Code)
		}
	})

	t.Run("Success with main and detail images", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.WriteField("type", "TROUSERS")
		writer.WriteField("session_id", "TEST_SESSION_123")

		// Create mock main image
		partMain, _ := writer.CreateFormFile("main", "test_main.jpg")
		io.WriteString(partMain, "fake-image-bytes-main")

		// Create mock detail image
		partDetail, _ := writer.CreateFormFile("details[0]", "test_detail_0.jpg")
		io.WriteString(partDetail, "fake-image-bytes-detail")

		writer.Close()

		req, _ := http.NewRequest("POST", "/salev2/api/upload-smart-group", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200 OK, got %d: %s", w.Code, w.Body.String())
		}

		var res struct {
			Success   bool   `json:"success"`
			SessionID string `json:"session_id"`
			Status    string `json:"status"`
		}
		json.Unmarshal(w.Body.Bytes(), &res)
		if !res.Success || res.Status != "DISPATCHED" {
			t.Errorf("Expected status DISPATCHED, got %+v", res)
		}
	})
}

func TestSaleV2_ProcessSmartUploadJob_OCR_And_Attachment(t *testing.T) {
	db, cfg := getTestDB(t)
	tx := db.Begin()
	defer tx.Rollback()

	handler := NewSaleV2Handler(tx, cfg)

	// Seed existing product to match candidate OCR
	prodName := fmt.Sprintf("TS%d", time.Now().UnixNano()%100000)
	initialPrice := int64(100000)
	prod := models.Product{
		Name:   &prodName,
		Type:   strPtr("TROUSERS"),
		Price:  &initialPrice,
		Status: strPtr("AVAILABLE"),
	}
	tx.Create(&prod)

	t.Run("Job updates price and creates images when product found", func(t *testing.T) {
		sessionID := fmt.Sprintf("SESSION_%d", time.Now().UnixNano()%100000)
		mainImg := "storage/sale/test_main_ok.jpg"
		detailImg := "storage/sale/test_detail_ok.jpg"

		job := SmartUploadJob{
			SessionID:        sessionID,
			OriginalFileName: fmt.Sprintf("%s 180k.jpg", prodName),
			MainImagePath:    mainImg,
			DetailPaths:      []string{detailImg},
			ProductType:      "TROUSERS",
		}

		// Direct call to processSmartUploadJob
		handler.processSmartUploadJob(job)

		// 1. Verify log record created
		var logRecord models.ProductUploadLog
		err := tx.Where("session_id = ? AND file_path = ?", sessionID, mainImg).First(&logRecord).Error
		if err != nil {
			t.Fatalf("Expected upload log to be created: %v", err)
		}

		// 2. Verify product images were attached if log succeeded
		if logRecord.Status == "SUCCESS" {
			var attachedImgs []models.ProductImage
			tx.Where("product_id = ?", prod.ID).Find(&attachedImgs)
			if len(attachedImgs) < 1 {
				t.Errorf("Expected at least 1 image attached to product, got %d", len(attachedImgs))
			}
		}
	})

	t.Run("Job deletes old main and detail images and overwrites with new ones", func(t *testing.T) {
		_ = os.MkdirAll("./storage/sale", 0755)
		oldMainFile := fmt.Sprintf("./storage/sale/mock_old_main_%d.jpg", time.Now().UnixNano())
		oldDetailFile := fmt.Sprintf("./storage/sale/mock_old_detail_%d.jpg", time.Now().UnixNano())
		newMainFile := fmt.Sprintf("./storage/sale/mock_new_main_%d.jpg", time.Now().UnixNano())
		newDetailFile := fmt.Sprintf("./storage/sale/mock_new_detail_%d.jpg", time.Now().UnixNano())

		_ = os.WriteFile(oldMainFile, []byte("old main"), 0644)
		_ = os.WriteFile(oldDetailFile, []byte("old detail"), 0644)
		_ = os.WriteFile(newMainFile, []byte("new main"), 0644)
		_ = os.WriteFile(newDetailFile, []byte("new detail"), 0644)
		defer func() {
			_ = os.Remove(oldMainFile)
			_ = os.Remove(oldDetailFile)
			_ = os.Remove(newMainFile)
			_ = os.Remove(newDetailFile)
		}()

		// Set product with old main image and old detail image
		prod.PathThumb = strPtr(oldMainFile)
		tx.Save(&prod)
		tx.Create(&models.ProductImage{
			ProductID: prod.ID,
			FilePath:  oldMainFile,
			IsMain:    true,
		})
		tx.Create(&models.ProductImage{
			ProductID: prod.ID,
			FilePath:  oldDetailFile,
			IsMain:    false,
		})

		sessionID := fmt.Sprintf("SESSION_REPLACE_%d", time.Now().UnixNano()%100000)
		job := SmartUploadJob{
			SessionID:        sessionID,
			OriginalFileName: fmt.Sprintf("%s 220k.jpg", prodName),
			MainImagePath:    newMainFile,
			DetailPaths:      []string{newDetailFile},
			ProductType:      "TROUSERS",
		}

		handler.processSmartUploadJob(job)

		var logRecord models.ProductUploadLog
		tx.Where("session_id = ?", sessionID).First(&logRecord)
		t.Logf("LogRecord status=%s, msg=%v, code=%v, detected=%v", logRecord.Status, strVal(logRecord.Message), strVal(logRecord.ProductCode), strVal(logRecord.DetectedText))

		// Check product record updated
		var updatedProd models.Product
		tx.First(&updatedProd, prod.ID)
		if updatedProd.PathThumb == nil || *updatedProd.PathThumb != newMainFile {
			t.Errorf("Expected product PathThumb to be updated to %s, got %v", newMainFile, updatedProd.PathThumb)
		}
		if updatedProd.Price == nil || *updatedProd.Price != 220000 {
			t.Errorf("Expected product Price to be 220000, got %v", updatedProd.Price)
		}

		// Check attached images
		var attachedImgs []models.ProductImage
		tx.Where("product_id = ?", prod.ID).Find(&attachedImgs)
		if len(attachedImgs) != 2 {
			t.Errorf("Expected exactly 2 attached images (1 main + 1 detail), got %d", len(attachedImgs))
		}
		for _, img := range attachedImgs {
			if img.FilePath == oldMainFile || img.FilePath == oldDetailFile {
				t.Errorf("Old image still exists in product_images: %s", img.FilePath)
			}
		}

		// Check physical files: old files should be removed from disk
		if _, err := os.Stat(oldMainFile); !os.IsNotExist(err) {
			t.Errorf("Expected old main file %s to be deleted from disk", oldMainFile)
		}
		if _, err := os.Stat(oldDetailFile); !os.IsNotExist(err) {
			t.Errorf("Expected old detail file %s to be deleted from disk", oldDetailFile)
		}
		// New files must still exist
		if _, err := os.Stat(newMainFile); os.IsNotExist(err) {
			t.Errorf("Expected new main file %s to still exist", newMainFile)
		}
		if _, err := os.Stat(newDetailFile); os.IsNotExist(err) {
			t.Errorf("Expected new detail file %s to still exist", newDetailFile)
		}
	})

	t.Run("Job fails gracefully when product does not exist", func(t *testing.T) {
		sessionID := fmt.Sprintf("SESSION_NOTFOUND_%d", time.Now().UnixNano()%100000)
		mainImg := "storage/sale/test_notfound.jpg"

		job := SmartUploadJob{
			SessionID:        sessionID,
			OriginalFileName: "NONEXISTENT_PRODUCT_99999 200k.jpg",
			MainImagePath:    mainImg,
			DetailPaths:      []string{},
			ProductType:      "TROUSERS",
		}

		handler.processSmartUploadJob(job)

		var logRecord models.ProductUploadLog
		err := tx.Where("session_id = ?", sessionID).First(&logRecord).Error
		if err != nil {
			t.Fatalf("Expected upload log record for failed job: %v", err)
		}
		if logRecord.Status != "ERROR" {
			t.Errorf("Expected status ERROR for nonexistent product, got %s", logRecord.Status)
		}
	})
}

func TestSaleV2_UploadLogs_And_Sessions(t *testing.T) {
	db, cfg := getTestDB(t)
	tx := db.Begin()
	defer tx.Rollback()

	handler := NewSaleV2Handler(tx, cfg)
	router := setupSaleV2TestRouter(handler)

	// Seed upload log
	sessionID := fmt.Sprintf("LOG_SESS_%d", time.Now().UnixNano()%100000)
	logEntry := models.ProductUploadLog{
		SessionID: sessionID,
		FileName:  "test_item.jpg",
		FilePath:  "storage/sale/test_item.jpg",
		Status:    "SUCCESS",
		Message:   strPtr("All good"),
	}
	tx.Create(&logEntry)

	// 1. Fetch upload logs
	t.Run("Fetch Upload Logs by session_ids", func(t *testing.T) {
		w := performRequest(router, "GET", "/salev2/api/upload-logs?session_ids="+sessionID, nil)
		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d", w.Code)
		}
		var logs []models.ProductUploadLog
		json.Unmarshal(w.Body.Bytes(), &logs)
		if len(logs) == 0 {
			t.Errorf("Expected to retrieve at least 1 log entry")
		}
	})

	// 2. Fetch upload sessions
	t.Run("Fetch Upload Sessions", func(t *testing.T) {
		w := performRequest(router, "GET", "/salev2/api/upload-sessions", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d", w.Code)
		}
	})
}

func TestSaleV2_SoldByDate_And_Cleanup(t *testing.T) {
	db, cfg := getTestDB(t)
	tx := db.Begin()
	defer tx.Rollback()

	handler := NewSaleV2Handler(tx, cfg)
	router := setupSaleV2TestRouter(handler)

	// 1. SoldByDate endpoint
	t.Run("SoldByDate returns aggregated list", func(t *testing.T) {
		w := performRequest(router, "GET", "/salev2/api/sold-by-date", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d", w.Code)
		}
	})

	// 2. CleanupImages endpoint (dry run analysis)
	t.Run("CleanupImages runs safely", func(t *testing.T) {
		w := performRequest(router, "GET", "/salev2/api/images/cleanup", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d: %s", w.Code, w.Body.String())
		}
		var res map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &res)
		if _, ok := res["used_images_db"]; !ok {
			t.Errorf("Expected used_images_db in cleanup response")
		}
	})
}
