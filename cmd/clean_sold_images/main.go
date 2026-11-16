package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"sale-app/config"
	"sale-app/models"
	"sale-app/services"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	daysFlag := flag.Int("days", 60, "Only clean images of products whose order was completed or created more than N days ago")
	dryRun := flag.Bool("dry-run", false, "Preview which products and images will be cleaned without deleting")
	forceAllSold := flag.Bool("all-sold", false, "Clean images for ALL products with status = SOLD regardless of order date")
	flag.Parse()

	cfg := config.LoadConfig()
	db, err := gorm.Open(mysql.Open(cfg.DatabaseDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("❌ Database connection error: %v", err)
	}

	storage := services.NewStorageService()

	// 1. Query SOLD products that still have images
	querySQL := `
		SELECT p.id, p.name, p.type, p.status, p.path_thumb, p.image_thumb_scale, p.fb_album_photo_id,
		       COALESCE(MAX(o.created_at), p.updated_at) as last_order_date,
		       COALESCE(MAX(o.status), 'NO_ORDER') as order_status
		FROM products p
		LEFT JOIN order_products op ON p.id = op.product_id
		LEFT JOIN orders o ON op.order_id = o.id
		WHERE p.status = 'SOLD'
		  AND ((p.path_thumb IS NOT NULL AND p.path_thumb != '') 
		       OR (p.image_thumb_scale IS NOT NULL AND p.image_thumb_scale != ''))
	`

	if !*forceAllSold {
		querySQL += fmt.Sprintf(" GROUP BY p.id HAVING (order_status = 'COMPLETED' OR last_order_date < NOW() - INTERVAL %d DAY)", *daysFlag)
	} else {
		querySQL += " GROUP BY p.id"
	}

	type SoldItem struct {
		ID              uint
		Name            *string
		Type            *string
		Status          *string
		PathThumb       *string
		ImageThumbScale *string
		FBAlbumPhotoID  *string
		LastOrderDate   *string
		OrderStatus     string
	}

	var items []SoldItem
	if err := db.Raw(querySQL).Scan(&items).Error; err != nil {
		log.Fatalf("❌ Lỗi query sản phẩm đã bán: %v", err)
	}

	fmt.Println("==================================================")
	fmt.Printf("🧹 DỌN DẸP ẢNH SẢN PHẨM ĐÃ BÁN (SOLD)\n")
	fmt.Printf("📌 Điều kiện: Đã bán > %d ngày hoặc đơn hàng COMPLETED\n", *daysFlag)
	fmt.Printf("🔍 Tìm thấy: %d sản phẩm có ảnh cần clear\n", len(items))
	if *dryRun {
		fmt.Println("⚠️ CHẾ ĐỘ DRY-RUN: Chỉ xem trước, không xóa bất kỳ file hoặc dữ liệu nào.")
	}
	fmt.Println("==================================================")

	if len(items) == 0 {
		fmt.Println("ℹ️ Không có ảnh sản phẩm đã bán nào cần dọn dẹp.")
		return
	}

	var cloudinaryIDs []string
	var localFilesToDelete []string
	var productIDsToUpdate []uint

	for i, item := range items {
		name := ""
		if item.Name != nil {
			name = *item.Name
		}
		thumb := ""
		if item.PathThumb != nil {
			thumb = *item.PathThumb
		}
		scale := ""
		if item.ImageThumbScale != nil {
			scale = *item.ImageThumbScale
		}

		fmt.Printf("[%d/%d] ID: %-8d | Mã: %-10s | Đơn: %-10s\n", i+1, len(items), item.ID, name, item.OrderStatus)
		if thumb != "" {
			fmt.Printf("   - Ảnh gốc:  %s\n", thumb)
			if pid := services.ExtractPublicID(thumb); pid != "" {
				cloudinaryIDs = append(cloudinaryIDs, pid)
			}
			localFilesToDelete = append(localFilesToDelete, thumb)
		}
		if scale != "" {
			fmt.Printf("   - Ảnh scale: %s\n", scale)
			if pid := services.ExtractPublicID(scale); pid != "" {
				cloudinaryIDs = append(cloudinaryIDs, pid)
			}
			localFilesToDelete = append(localFilesToDelete, scale)
		}

		productIDsToUpdate = append(productIDsToUpdate, item.ID)
	}

	if *dryRun {
		fmt.Println("\n[DRY RUN] Kết thúc xem trước.")
		return
	}

	// 2. Xóa Cloudinary assets (cả ảnh gốc và scale)
	if len(cloudinaryIDs) > 0 && storage.CloudName != "" {
		fmt.Printf("\n☁️ Đang xóa %d file trên Cloudinary (cả ảnh gốc + scale)...\n", len(cloudinaryIDs))
		deletedCld, err := storage.DeleteFromCloudinary(cloudinaryIDs)
		if err != nil {
			fmt.Printf("⚠️ Lỗi Cloudinary: %v\n", err)
		} else {
			fmt.Printf("✅ Đã xóa thành công %d asset trên Cloudinary.\n", deletedCld)
		}
	}

	// 3. Xóa file local (nếu có)
	localDeleted := 0
	for _, rawPath := range localFilesToDelete {
		trimmed := strings.TrimSpace(rawPath)
		cleanLoc := strings.TrimLeft(trimmed, "/")
		if strings.HasPrefix(cleanLoc, "storage/") {
			cleanLoc = strings.TrimPrefix(cleanLoc, "storage/")
		}
		candidates := []string{
			trimmed,
			filepath.Join("./storage/sale", filepath.Base(cleanLoc)),
			filepath.Join("./storage", cleanLoc),
		}
		for _, cand := range candidates {
			if info, err := os.Stat(cand); err == nil && !info.IsDir() {
				if err := os.Remove(cand); err == nil {
					localDeleted++
					break
				}
			}
		}
	}
	if localDeleted > 0 {
		fmt.Printf("📁 Đã xóa %d file trên ổ đĩa local.\n", localDeleted)
	}

	// 4. Cập nhật Database: Clear cả path_thumb VÀ image_thumb_scale thành NULL
	batchSize := 500
	totalUpdated := int64(0)
	for i := 0; i < len(productIDsToUpdate); i += batchSize {
		end := i + batchSize
		if end > len(productIDsToUpdate) {
			end = len(productIDsToUpdate)
		}
		batch := productIDsToUpdate[i:end]
		res := db.Model(&models.Product{}).
			Where("id IN ?", batch).
			Updates(map[string]interface{}{
				"path_thumb":        nil,
				"image_thumb_scale": nil,
			})
		totalUpdated += res.RowsAffected
	}

	fmt.Println("==================================================")
	fmt.Printf("🎉 ĐÃ HOÀN TẤT DỌN DẸP:\n")
	fmt.Printf(" - Đã reset path_thumb & image_thumb_scale của %d sản phẩm về NULL\n", totalUpdated)
	fmt.Printf(" - Đã xóa ảnh gốc và ảnh scale tương ứng trên Cloudinary & Storage\n")
	fmt.Println("==================================================")
}
