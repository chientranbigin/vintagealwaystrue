package main

import (
	"fmt"
	"log"
	"time"

	"sale-app/config"
	"sale-app/models"
	"sale-app/services"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	cfg := config.LoadConfig()
	db, err := gorm.Open(mysql.Open(cfg.DatabaseDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("Lỗi kết nối DB: %v", err)
	}

	fbService := services.NewFacebookService(db, cfg.FBPageID, cfg.FBPageToken, cfg.MediaBaseURL, cfg.FBMainPageID, cfg.FBMainPageToken)
	if fbService.MainPageToken == "" {
		log.Fatalf("Thiếu FB_MAIN_PAGE_TOKEN")
	}

	// Lấy tất cả sản phẩm SOLD nhưng vẫn còn fb_album_photo_id
	var soldProducts []models.Product
	db.Where("status = ? AND fb_album_photo_id IS NOT NULL AND fb_album_photo_id != ''", "SOLD").
		Order("updated_at DESC").
		Find(&soldProducts)

	fmt.Println("==================================================")
	fmt.Printf("🔍 TÌM THẤY %d SẢN PHẨM ĐÃ BÁN (SOLD) NHƯNG VẪN CÒN ẢNH TRONG ALBUM:\n", len(soldProducts))
	fmt.Println("==================================================")

	successCount := 0
	failCount := 0

	for i, p := range soldProducts {
		pName := ""
		if p.Name != nil {
			pName = *p.Name
		}
		photoID := *p.FBAlbumPhotoID

		fmt.Printf("[%d/%d] Đang xóa ảnh %s (%s) khỏi Album Facebook... ", i+1, len(soldProducts), pName, photoID)

		err := fbService.RemoveProductPhotoFromAlbum(photoID)
		if err != nil {
			fmt.Printf("❌ Lỗi: %v\n", err)
			failCount++
		} else {
			// Cập nhật DB
			db.Model(&models.Product{}).Where("id = ?", p.ID).Update("fb_album_photo_id", nil)
			fmt.Printf("✅ ĐÃ XÓA THÀNH CÔNG!\n")
			successCount++
		}

		time.Sleep(1 * time.Second)
	}

	fmt.Println("==================================================")
	fmt.Printf("🎉 HOÀN TẤT DỌN DẸP: %d thành công, %d thất bại.\n", successCount, failCount)
	fmt.Println("==================================================")
}
