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

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	cfg := config.LoadConfig()
	albumID := "122252355908280390" // Trousers Album ID
	pageToken := cfg.FBMainPageToken

	if pageToken == "" {
		log.Fatalf("Thiếu token FB_MAIN_PAGE_TOKEN")
	}

	db, err := gorm.Open(mysql.Open(cfg.DatabaseDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("Lỗi DB: %v", err)
	}

	client := &http.Client{Timeout: 30 * time.Second}

	// 1. Lấy danh sách ảnh hiện có trong Album
	endpoint := fmt.Sprintf("https://graph.facebook.com/v19.0/%s/photos?fields=id&limit=200&access_token=%s", albumID, url.QueryEscape(pageToken))
	resp, err := client.Get(endpoint)
	if err != nil {
		log.Fatalf("Lỗi fetch ảnh: %v", err)
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	var res struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(respBytes, &res); err != nil {
		log.Fatalf("Lỗi parse JSON: %v", err)
	}

	fmt.Printf("Tìm thấy %d ảnh trong Album Facebook '%s'. Đang xóa dọn dẹp...\n", len(res.Data), albumID)
	for i, item := range res.Data {
		delURL := fmt.Sprintf("https://graph.facebook.com/v19.0/%s?access_token=%s", item.ID, url.QueryEscape(pageToken))
		req, _ := http.NewRequest("DELETE", delURL, nil)
		delResp, dErr := client.Do(req)
		if dErr == nil {
			delResp.Body.Close()
			fmt.Printf(" [%d/%d] Đã xóa photo %s\n", i+1, len(res.Data), item.ID)
		}
		time.Sleep(300 * time.Millisecond)
	}

	// 2. Reset fb_album_photo_id trong database
	db.Exec("UPDATE products SET fb_album_photo_id = NULL WHERE type = 'TROUSERS'")
	fmt.Println("✅ Đã reset sạch toàn bộ ảnh cũ trong Album và Database!")
}
