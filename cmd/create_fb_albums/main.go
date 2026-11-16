package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"sale-app/config"
	"sale-app/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type AlbumDef struct {
	Code        string
	AlbumName   string
	Description string
}

func main() {
	cfg := config.LoadConfig()

	pageID := cfg.FBMainPageID
	pageToken := cfg.FBMainPageToken

	if pageID == "" || pageToken == "" {
		log.Fatalf("❌ Thiếu FB_MAIN_PAGE_ID hoặc FB_MAIN_PAGE_TOKEN trong file .env")
	}

	fmt.Printf("Connecting to TiDB Cloud MySQL database...\n")
	db, err := gorm.Open(mysql.Open(cfg.DatabaseDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	albumDefs := []AlbumDef{
		{
			Code:        "TROUSERS",
			AlbumName:   "Quần Âu - Trousers",
			Description: "Bộ sưu tập Quần Âu tuyển chọn tại Vintage Always True House",
		},
		{
			Code:        "JACKET",
			AlbumName:   "Áo Khoác - Jacket",
			Description: "Bộ sưu tập Áo Khoác tuyển chọn tại Vintage Always True House",
		},
		{
			Code:        "HARRINGTON_JACKET",
			AlbumName:   "Áo Khoác Harrington",
			Description: "Bộ sưu tập Áo Khoác Harrington Jacket tại Vintage Always True House",
		},
		{
			Code:        "SHIRT",
			AlbumName:   "Áo Sơ Mi - Shirt",
			Description: "Bộ sưu tập Áo Sơ Mi tuyển chọn tại Vintage Always True House",
		},
		{
			Code:        "BLAZER",
			AlbumName:   "Áo Blazer",
			Description: "Bộ sưu tập Áo Blazer cổ điển và hiện đại tại Vintage Always True House",
		},
		{
			Code:        "SUIT",
			AlbumName:   "Bộ Suit",
			Description: "Bộ sưu tập Suits cao cấp tại Vintage Always True House",
		},
		{
			Code:        "GILE",
			AlbumName:   "Áo Gile - Vest",
			Description: "Bộ sưu tập Áo Gile tuyển chọn tại Vintage Always True House",
		},
		{
			Code:        "POLO SHIRT",
			AlbumName:   "Áo Polo",
			Description: "Bộ sưu tập Áo Polo tuyển chọn tại Vintage Always True House",
		},
		{
			Code:        "TIE",
			AlbumName:   "Cà Vạt - Tie",
			Description: "Bộ sưu tập Cà Vạt lụa và vintage tại Vintage Always True House",
		},
		{
			Code:        "BELT",
			AlbumName:   "Thắt Lưng - Belt",
			Description: "Bộ sưu tập Thắt Lưng da cao cấp tại Vintage Always True House",
		},
		{
			Code:        "HAT",
			AlbumName:   "Mũ - Hat",
			Description: "Bộ sưu tập Mũ phong cách vintage tại Vintage Always True House",
		},
	}

	client := &http.Client{Timeout: 30 * time.Second}

	fmt.Println("==================================================")
	fmt.Printf("🚀 Bắt đầu tạo Album trên Page '%s' (ID: %s)\n", "Vintage Always True House", pageID)
	fmt.Println("==================================================")

	createdCount := 0
	for i, def := range albumDefs {
		var cat models.Category
		if err := db.Where("code = ?", def.Code).First(&cat).Error; err != nil {
			fmt.Printf("[%d/%d] ⚠️ Không tìm thấy category %s trong DB\n", i+1, len(albumDefs), def.Code)
			continue
		}

		if cat.FBAlbumID != nil && *cat.FBAlbumID != "" {
			fmt.Printf("[%d/%d] ⏭️ %s (%s) đã có Album ID: %s (Bỏ qua)\n", i+1, len(albumDefs), def.Code, def.AlbumName, *cat.FBAlbumID)
			continue
		}

		fmt.Printf("[%d/%d] Đang tạo Album '%s'... ", i+1, len(albumDefs), def.AlbumName)

		endpoint := fmt.Sprintf("https://graph.facebook.com/v19.0/%s/albums", pageID)
		formData := url.Values{}
		formData.Set("name", def.AlbumName)
		formData.Set("message", def.Description)
		formData.Set("access_token", pageToken)

		resp, err := client.PostForm(endpoint, formData)
		if err != nil {
			fmt.Printf("❌ Lỗi request: %v\n", err)
			continue
		}
		respBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			fmt.Printf("❌ Lỗi Facebook (%d): %s\n", resp.StatusCode, string(respBytes))
			continue
		}

		var res struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(respBytes, &res); err != nil || res.ID == "" {
			fmt.Printf("❌ Không parse được Album ID: %s\n", string(respBytes))
			continue
		}

		// Update database
		db.Model(&models.Category{}).Where("id = ?", cat.ID).Update("fb_album_id", res.ID)
		fmt.Printf("✅ THÀNH CÔNG! Album ID: %s\n", res.ID)
		createdCount++

		// Delay 1.5s to avoid hitting Facebook rate limit
		time.Sleep(1500 * time.Millisecond)
	}

	fmt.Println("==================================================")
	fmt.Printf("🎉 Đã hoàn tất! Tạo thành công %d Album và tự động lưu vào DB.\n", createdCount)
	fmt.Println("==================================================")
}
