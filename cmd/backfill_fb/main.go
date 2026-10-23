package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"sale-app/config"
	"sale-app/models"
	"sale-app/services"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	sessionFlag := flag.String("session", "", "Filter by specific upload session ID (default: all or latest)")
	latestSessionFlag := flag.Bool("latest-session", false, "Only backfill products from the latest upload session")
	typeFlag := flag.String("type", "", "Filter by product type (e.g. TROUSERS, BLAZER)")
	productFlag := flag.String("product", "", "Specific product name/code (e.g. TS1235)")
	dryRun := flag.Bool("dry-run", false, "Preview what would be posted without syncing")
	delaySec := flag.Int("delay", 3, "Delay in seconds between posts to prevent FB rate limit")
	force := flag.Bool("force", false, "Re-create FB post even if product already has fb_post_id")
	flag.Parse()

	cfg := config.LoadConfig()
	db, err := gorm.Open(mysql.Open(cfg.DatabaseDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 1. Inspect recent upload sessions
	var recentLogs []models.ProductUploadLog
	db.Select("session_id, MAX(created_at) as created_at, COUNT(*) as id").
		Group("session_id").
		Order("created_at DESC").
		Limit(5).
		Find(&recentLogs)

	fmt.Println("==================================================")
	fmt.Println("📋 5 Upload Sessions gần đây nhất:")
	for _, rl := range recentLogs {
		createdStr := ""
		if rl.CreatedAt != nil {
			createdStr = rl.CreatedAt.Format("2006-01-02 15:04:05")
		}
		fmt.Printf(" - Session ID: %s (%s, %d ảnh/logs)\n", rl.SessionID, createdStr, rl.ID)
	}
	fmt.Println("==================================================")

	targetSession := *sessionFlag
	if *latestSessionFlag && len(recentLogs) > 0 {
		targetSession = recentLogs[0].SessionID
		fmt.Printf("🎯 Chọn phiên mới nhất: %s\n", targetSession)
	}

	// 2. Build Query
	var productIDs []uint
	if targetSession != "" {
		db.Table("product_upload_logs").
			Where("session_id = ? AND product_id IS NOT NULL", targetSession).
			Distinct("product_id").
			Pluck("product_id", &productIDs)
		fmt.Printf("Tìm thấy %d sản phẩm trong session '%s'\n", len(productIDs), targetSession)
	}

	query := db.Model(&models.Product{}).
		Preload("Images").
		Preload("Sizes").
		Where("path_thumb IS NOT NULL AND path_thumb != ''")

	if targetSession != "" {
		if len(productIDs) == 0 {
			fmt.Println("Không tìm thấy sản phẩm nào trong session chỉ định.")
			return
		}
		query = query.Where("id IN ?", productIDs)
	} else if *productFlag != "" {
		query = query.Where("name = ?", *productFlag)
	} else {
		// Default: AVAILABLE products with no fb_post_id
		if !*force {
			query = query.Where("fb_post_id IS NULL OR fb_post_id = ''")
		}
		query = query.Where("status = 'AVAILABLE'")
	}

	if *typeFlag != "" {
		query = query.Where("type = ?", strings.ToUpper(*typeFlag))
	}

	if !*force && targetSession != "" {
		query = query.Where("fb_post_id IS NULL OR fb_post_id = ''")
	}

	var products []models.Product
	if err := query.Order("id ASC").Find(&products).Error; err != nil {
		log.Fatalf("Query products error: %v", err)
	}

	if len(products) == 0 {
		fmt.Println("✅ Không có sản phẩm nào cần backfill (tất cả đã có fb_post_id hoặc không khớp bộ lọc).")
		return
	}

	fmt.Printf("🔍 Tìm thấy %d sản phẩm cần đồng bộ lên Facebook:\n", len(products))
	for i, p := range products {
		pName := ""
		if p.Name != nil {
			pName = *p.Name
		}
		pType := ""
		if p.Type != nil {
			pType = *p.Type
		}
		pStatus := ""
		if p.Status != nil {
			pStatus = *p.Status
		}
		fbID := "Chưa có"
		if p.FBPostID != nil && *p.FBPostID != "" {
			fbID = *p.FBPostID
		}
		fmt.Printf(" %2d. ID: %-5d | Tên: %-10s | Loại: %-10s | Status: %-10s | FB Post: %s\n",
			i+1, p.ID, pName, pType, pStatus, fbID)
	}

	if *dryRun {
		fmt.Println("\n[DRY RUN] Đã xem trước danh sách, không thực hiện gọi Facebook API.")
		return
	}

	// 3. Initialize FacebookService
	fbService := services.NewFacebookService(db, cfg.FBPageID, cfg.FBPageToken, cfg.MediaBaseURL)
	if fbService.PageToken == "" {
		log.Fatalf("❌ Lỗi: FB_PAGE_TOKEN chưa được cấu hình. Vui lòng kiểm tra biến môi trường hoặc file .env")
	}

	fmt.Printf("\n🚀 Bắt đầu đồng bộ %d sản phẩm lên Facebook (giãn cách %ds/post)...\n", len(products), *delaySec)
	successCount := 0
	failCount := 0

	for i, p := range products {
		pName := ""
		if p.Name != nil {
			pName = *p.Name
		}
		fmt.Printf("[%d/%d] Đang đồng bộ %s (ID %d)... ", i+1, len(products), pName, p.ID)

		err := fbService.ExecuteSync(p.ID, services.FBActionCreate)
		if err != nil {
			fmt.Printf("❌ Thất bại: %v\n", err)
			failCount++
		} else {
			var updated models.Product
			db.Select("fb_post_id").First(&updated, p.ID)
			newPostID := ""
			if updated.FBPostID != nil {
				newPostID = *updated.FBPostID
			}
			fmt.Printf("✅ Thành công! Post ID: %s\n", newPostID)
			successCount++
		}

		if i < len(products)-1 && *delaySec > 0 {
			time.Sleep(time.Duration(*delaySec) * time.Second)
		}
	}

	fmt.Println("==================================================")
	fmt.Printf("🎉 Hoàn tất backfill: %d thành công, %d thất bại.\n", successCount, failCount)
	fmt.Println("==================================================")
}
