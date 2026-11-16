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
	typeFlag := flag.String("type", "ALL", "Product type code (e.g. ALL, TROUSERS, BLAZER, or comma-separated: SHIRT,GILE)")
	limitFlag := flag.Int("limit", 1000, "Max products to backfill per category")
	delayFlag := flag.Int("delay", 4, "Delay in seconds between uploads to strictly prevent Facebook rate limits")
	dryRun := flag.Bool("dry-run", false, "Preview products and captions without uploading to Facebook")
	force := flag.Bool("force", false, "Re-upload even if product already has fb_album_photo_id")
	reverseOrder := flag.Bool("order-asc", true, "Upload older items first so the newest uploaded item stays on TOP of Facebook Album")
	flag.Parse()

	cfg := config.LoadConfig()

	db, err := gorm.Open(mysql.Open(cfg.DatabaseDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("❌ Lỗi kết nối database: %v", err)
	}

	fbService := services.NewFacebookService(db, cfg.FBPageID, cfg.FBPageToken, cfg.MediaBaseURL, cfg.FBMainPageID, cfg.FBMainPageToken)
	if fbService.MainPageToken == "" {
		log.Fatalf("❌ Chưa cấu hình FB_MAIN_PAGE_TOKEN!")
	}

	// Xác định danh sách Category cần backfill
	var categories []models.Category
	reqType := strings.ToUpper(strings.TrimSpace(*typeFlag))

	if reqType == "ALL" {
		if err := db.Where("is_active = true AND fb_album_id IS NOT NULL AND fb_album_id != ''").Order("id ASC").Find(&categories).Error; err != nil {
			log.Fatalf("❌ Lỗi query categories: %v", err)
		}
	} else if strings.Contains(reqType, ",") {
		codes := strings.Split(reqType, ",")
		for i := range codes {
			codes[i] = strings.TrimSpace(codes[i])
		}
		if err := db.Where("code IN ? AND is_active = true", codes).Find(&categories).Error; err != nil {
			log.Fatalf("❌ Lỗi query categories: %v", err)
		}
	} else {
		var cat models.Category
		if err := db.Where("code = ? AND is_active = true", reqType).First(&cat).Error; err != nil {
			log.Fatalf("❌ Không tìm thấy danh mục '%s' trong bảng categories", reqType)
		}
		categories = append(categories, cat)
	}

	if len(categories) == 0 {
		log.Fatalf("❌ Không có danh mục nào hợp lệ được chọn.")
	}

	fmt.Println("==================================================")
	fmt.Printf("📦 TỔNG SỐ DANH MỤC CẦN XỬ LÝ: %d\n", len(categories))
	for i, c := range categories {
		album := "Chưa có album"
		if c.FBAlbumID != nil && *c.FBAlbumID != "" {
			album = *c.FBAlbumID
		}
		fmt.Printf("  %d. [%-10s] %-25s | Album ID: %s\n", i+1, c.Code, c.Name, album)
	}
	fmt.Println("==================================================")

	totalSuccess := 0
	totalFail := 0

	for catIdx, cat := range categories {
		if cat.FBAlbumID == nil || *cat.FBAlbumID == "" {
			fmt.Printf("\n⚠️ Bỏ qua [%s]: Chưa cấu hình FBAlbumID\n", cat.Code)
			continue
		}

		fmt.Printf("\n==================================================\n")
		fmt.Printf("[%d/%d] 📂 ĐANG XỬ LÝ: %s (%s)\n", catIdx+1, len(categories), cat.Code, cat.Name)
		fmt.Printf("🎯 Facebook Album ID: %s\n", *cat.FBAlbumID)
		fmt.Println("==================================================")

		subQuery := "(SELECT product_id, DATE_FORMAT(MAX(created_at), '%Y-%m-%d %H:%i:%s') as latest_upload FROM product_upload_logs WHERE product_id IS NOT NULL GROUP BY product_id) as upload_summary"

		query := db.Model(&models.Product{}).
			Joins("LEFT JOIN " + subQuery + " ON upload_summary.product_id = products.id").
			Preload("Images").
			Preload("Sizes").
			Where("products.status = ?", "AVAILABLE").
			Where("products.type = ?", cat.Code).
			Where("products.path_thumb IS NOT NULL AND products.path_thumb != ''")

		if !*force {
			query = query.Where("products.fb_album_photo_id IS NULL OR products.fb_album_photo_id = ''")
		}

		query = query.Select("products.*, upload_summary.latest_upload").
			Order("COALESCE(upload_summary.latest_upload, products.created_at) DESC")

		var products []models.Product
		if err := query.Limit(*limitFlag).Find(&products).Error; err != nil {
			log.Printf("❌ Lỗi query sản phẩm cho %s: %v", cat.Code, err)
			continue
		}

		if len(products) == 0 {
			fmt.Printf("ℹ️ Danh mục %s không có sản phẩm AVAILABLE nào cần backfill.\n", cat.Code)
			continue
		}

		fmt.Printf("Tìm thấy %d sản phẩm AVAILABLE cần backfill cho %s:\n", len(products), cat.Code)
		for i, p := range products {
			pName := ""
			if p.Name != nil {
				pName = *p.Name
			}
			pPrice := int64(0)
			if p.Price != nil {
				pPrice = *p.Price
			}
			fbPost := "Chưa có post"
			if p.FBPostID != nil && *p.FBPostID != "" {
				fbPost = *p.FBPostID
			}
			fmt.Printf(" %2d. ID: %-8d | Mã: %-10s | Giá: %s | Post: %s\n",
				i+1, p.ID, pName, services.FormatVNDPrice(pPrice), fbPost)
		}

		// Đảo ngược thứ tự: cũ upload trước, mới upload sau -> mới nhất nằm trên cùng của Album Facebook
		itemsToUpload := make([]models.Product, len(products))
		if *reverseOrder && len(products) > 1 {
			fmt.Println("🔄 Thứ tự upload: Đẩy sản phẩm cũ trước, sản phẩm mới nhất sau cùng để hiển thị trên đầu Album.")
			for i := range products {
				itemsToUpload[i] = products[len(products)-1-i]
			}
		} else {
			copy(itemsToUpload, products)
		}

		if *dryRun {
			fmt.Printf("\n[DRY RUN] Đã kiểm tra %d sản phẩm của danh mục %s. Bỏ qua upload thật.\n", len(itemsToUpload), cat.Code)
			continue
		}

		// Upload thật
		delayDuration := time.Duration(*delayFlag) * time.Second
		fmt.Printf("\n🚀 Bắt đầu upload %d sản phẩm vào Album Facebook '%s' (giãn cách %ds/ảnh)...\n", len(itemsToUpload), *cat.FBAlbumID, *delayFlag)
		catSuccess := 0
		catFail := 0

		for i, p := range itemsToUpload {
			pName := ""
			if p.Name != nil {
				pName = *p.Name
			}
			thumb := p.BestQualityThumb()
			if thumb == "" {
				fmt.Printf("[%d/%d] ⚠️ Bỏ qua %s: Không có ảnh thumbnail.\n", i+1, len(itemsToUpload), pName)
				continue
			}

			postID := ""
			if p.FBPostID != nil {
				postID = *p.FBPostID
			}

			caption := fbService.BuildAlbumPhotoCaption(&p, postID)

			pct := float64(i+1) / float64(len(itemsToUpload)) * 100
			fmt.Printf("[%d/%d (%.1f%%)] Đang upload %s vào Album... ", i+1, len(itemsToUpload), pct, pName)

			maxRetries := 2
			uploaded := false
			for attempt := 0; attempt <= maxRetries; attempt++ {
				photoID, err := fbService.UploadProductPhotoToAlbum(*cat.FBAlbumID, thumb, caption)
				if err != nil {
					errStr := strings.ToLower(err.Error())
					if strings.Contains(errStr, "rate limit") || strings.Contains(errStr, "(#4)") || strings.Contains(errStr, "(#17)") || strings.Contains(errStr, "calls reached") {
						fmt.Printf("\n⚠️ Phát hiện Rate Limit từ Facebook! Tạm nghỉ 60s để hồi phục... (Thử lại lần %d)\n", attempt+1)
						time.Sleep(60 * time.Second)
						continue
					}
					fmt.Printf("❌ Thất bại: %v\n", err)
					break
				}

				db.Model(&models.Product{}).Where("id = ?", p.ID).Update("fb_album_photo_id", photoID)
				fmt.Printf("✅ THÀNH CÔNG! Photo ID: %s\n", photoID)
				catSuccess++
				uploaded = true
				break
			}

			if !uploaded {
				catFail++
			}

			if i < len(itemsToUpload)-1 {
				time.Sleep(delayDuration)
			}
		}

		fmt.Printf("📊 Kết quả danh mục %s: %d thành công, %d thất bại.\n", cat.Code, catSuccess, catFail)
		totalSuccess += catSuccess
		totalFail += catFail
	}

	fmt.Println("\n==================================================")
	fmt.Printf("🎉 TỔNG KẾT TOÀN BỘ: %d thành công, %d thất bại.\n", totalSuccess, totalFail)
	fmt.Println("==================================================")
}
