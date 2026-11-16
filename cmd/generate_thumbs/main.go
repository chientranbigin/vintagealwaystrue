package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"sale-app/config"
	"sale-app/models"
	"sale-app/services"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	typeFlag := flag.String("type", "", "Filter by product type (e.g. TROUSERS, BLAZER, or empty for all)")
	statusFlag := flag.String("status", "", "Filter by status (e.g. AVAILABLE, or empty for all)")
	limitFlag := flag.Int("limit", 0, "Limit number of products to process (0 = all)")
	forceFlag := flag.Bool("force", false, "Regenerate thumbnail even if image_thumb_scale already exists")
	workersFlag := flag.Int("workers", 5, "Number of concurrent workers")
	flag.Parse()

	cfg := config.LoadConfig()
	db, err := gorm.Open(mysql.Open(cfg.DatabaseDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("❌ Database connection error: %v", err)
	}

	storage := services.NewStorageService()

	// 1. Query products that need thumbnails
	query := db.Model(&models.Product{}).
		Where("path_thumb IS NOT NULL AND path_thumb != ''")

	if !*forceFlag {
		query = query.Where("image_thumb_scale IS NULL OR image_thumb_scale = ''")
	}

	if *typeFlag != "" {
		query = query.Where("type = ?", strings.ToUpper(*typeFlag))
	}

	if *statusFlag != "" {
		query = query.Where("status = ?", strings.ToUpper(*statusFlag))
	}

	if *limitFlag > 0 {
		query = query.Limit(*limitFlag)
	}

	// Process newest products first
	query = query.Order("id DESC")

	var products []models.Product
	if err := query.Find(&products).Error; err != nil {
		log.Fatalf("❌ Query products error: %v", err)
	}

	total := len(products)
	if total == 0 {
		fmt.Println("ℹ️ Không có sản phẩm nào cần tạo thumbnail scale.")
		return
	}

	fmt.Println("==================================================")
	fmt.Printf("🖼️ BẮT ĐẦU TẠO THUMBNAIL SCALE CHO %d SẢN PHẨM\n", total)
	if *typeFlag != "" {
		fmt.Printf("📦 Danh mục: %s\n", strings.ToUpper(*typeFlag))
	}
	if *statusFlag != "" {
		fmt.Printf("📌 Trạng thái: %s\n", strings.ToUpper(*statusFlag))
	}
	fmt.Printf("⚙️ Số worker đồng thời: %d\n", *workersFlag)
	fmt.Println("==================================================")

	startTime := time.Now()
	var successCount int64
	var failCount int64
	var processedCount int64

	jobs := make(chan models.Product, total)
	for _, p := range products {
		jobs <- p
	}
	close(jobs)

	var wg sync.WaitGroup
	workers := *workersFlag
	if workers < 1 {
		workers = 1
	}

	for w := 1; w <= workers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for p := range jobs {
				pName := ""
				if p.Name != nil {
					pName = *p.Name
				}
				thumb := ""
				if p.PathThumb != nil {
					thumb = *p.PathThumb
				}

				scaledPath, err := storage.CreateScaledThumbnailFromPath(thumb)
				curr := atomic.AddInt64(&processedCount, 1)
				pct := float64(curr) / float64(total) * 100

				if err != nil {
					atomic.AddInt64(&failCount, 1)
					fmt.Printf("[%d/%d (%.1f%%)] ❌ [%s] Thất bại: %v\n", curr, total, pct, pName, err)
					continue
				}

				// Update database
				if err := db.Model(&models.Product{}).Where("id = ?", p.ID).Update("image_thumb_scale", scaledPath).Error; err != nil {
					atomic.AddInt64(&failCount, 1)
					fmt.Printf("[%d/%d (%.1f%%)] ❌ [%s] Lỗi cập nhật DB: %v\n", curr, total, pct, pName, err)
					continue
				}

				atomic.AddInt64(&successCount, 1)
				fmt.Printf("[%d/%d (%.1f%%)] ✅ [%s] Thành công -> %s\n", curr, total, pct, pName, scaledPath)
			}
		}(w)
	}

	wg.Wait()
	duration := time.Since(startTime).Round(time.Second)

	fmt.Println("==================================================")
	fmt.Printf("🎉 HOÀN TẤT TRONG %v:\n", duration)
	fmt.Printf(" - Thành công: %d sản phẩm\n", successCount)
	fmt.Printf(" - Thất bại:   %d sản phẩm\n", failCount)
	fmt.Println("==================================================")
}
