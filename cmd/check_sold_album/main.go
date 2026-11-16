package main

import (
	"fmt"
	"log"

	"sale-app/config"
	"sale-app/models"

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

	// 1. Kiểm tra chi tiết TS4093
	var p models.Product
	err = db.Preload("Orders").Where("name LIKE ?", "%TS4093%").First(&p).Error
	if err != nil {
		fmt.Printf("Không tìm thấy sản phẩm TS4093: %v\n", err)
	} else {
		pName := ""
		if p.Name != nil {
			pName = *p.Name
		}
		pType := ""
		if p.Type != nil {
			pType = *p.Type
		}
		albumPhotoID := "NULL/EMPTY"
		if p.FBAlbumPhotoID != nil && *p.FBAlbumPhotoID != "" {
			albumPhotoID = *p.FBAlbumPhotoID
		}
		fbPostID := "NULL/EMPTY"
		if p.FBPostID != nil && *p.FBPostID != "" {
			fbPostID = *p.FBPostID
		}

		fmt.Println("==================================================")
		fmt.Printf("THÔNG TIN SẢN PHẨM %s (ID: %d):\n", pName, p.ID)
		fmt.Printf("  Status: %s\n", p.Status)
		fmt.Printf("  Type: %s\n", pType)
		fmt.Printf("  FBAlbumPhotoID: %s\n", albumPhotoID)
		fmt.Printf("  FBPostID: %s\n", fbPostID)
		fmt.Printf("  Created At: %v\n", p.CreatedAt)
		fmt.Printf("  Updated At: %v\n", p.UpdatedAt)

		// Kiểm tra Order liên quan
		type OrderItem struct {
			OrderID   uint
			ProductID uint
			CreatedAt string
		}
		var orderItems []OrderItem
		db.Table("order_items").Where("product_id = ?", p.ID).Find(&orderItems)
		fmt.Printf("  Số order items liên quan: %d\n", len(orderItems))
		for _, oi := range orderItems {
			var o models.Order
			db.First(&o, oi.OrderID)
			fmt.Printf("    -> Order ID: %d | Status: %s | CreatedAt: %v | UpdatedAt: %v\n", o.ID, o.Status, o.CreatedAt, o.UpdatedAt)
		}
		fmt.Println("==================================================")
	}

	// 2. Tìm tất cả sản phẩm có Status = 'SOLD' nhưng vẫn còn FBAlbumPhotoID
	var soldWithAlbum []models.Product
	db.Where("status = ? AND fb_album_photo_id IS NOT NULL AND fb_album_photo_id != ''", "SOLD").
		Order("updated_at DESC").
		Limit(50).
		Find(&soldWithAlbum)

	fmt.Printf("\nCÁC SẢN PHẨM STATUS = 'SOLD' NHƯNG VẪN CÒN FB_ALBUM_PHOTO_ID (Tổng: %d):\n", len(soldWithAlbum))
	for i, sp := range soldWithAlbum {
		name := ""
		if sp.Name != nil {
			name = *sp.Name
		}
		typ := ""
		if sp.Type != nil {
			typ = *sp.Type
		}
		photoID := ""
		if sp.FBAlbumPhotoID != nil {
			photoID = *sp.FBAlbumPhotoID
		}
		fmt.Printf(" %2d. ID: %-8d | Mã: %-10s | Loại: %-10s | AlbumPhotoID: %-20s | UpdatedAt: %v\n",
			i+1, sp.ID, name, typ, photoID, sp.UpdatedAt)
	}

	// 3. Tìm các sản phẩm được bán gần đây (ví dụ trong 7 ngày qua)
	var recentSold []models.Product
	db.Where("status = ?", "SOLD").
		Order("updated_at DESC").
		Limit(30).
		Find(&recentSold)

	fmt.Printf("\nTOP 30 SẢN PHẨM SOLD GẦN ĐÂY NHẤT:\n")
	for i, sp := range recentSold {
		name := ""
		if sp.Name != nil {
			name = *sp.Name
		}
		typ := ""
		if sp.Type != nil {
			typ = *sp.Type
		}
		photoID := "Đã xóa"
		if sp.FBAlbumPhotoID != nil && *sp.FBAlbumPhotoID != "" {
			photoID = "⚠️ VẪN CÒN: " + *sp.FBAlbumPhotoID
		}
		fmt.Printf(" %2d. ID: %-8d | Mã: %-10s | Loại: %-10s | AlbumPhoto: %-30s | UpdatedAt: %v\n",
			i+1, sp.ID, name, typ, photoID, sp.UpdatedAt)
	}
}
