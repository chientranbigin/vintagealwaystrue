package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"sale-app/models"
)

const resyncPauseDuration = 45 * time.Minute
const resyncUsageThresholdPct = 85

// FetchAlbumPhotoIDs retrieves all existing photo IDs in a Facebook album.
func (s *FacebookService) FetchAlbumPhotoIDs(albumID string) ([]string, error) {
	if albumID == "" {
		return nil, fmt.Errorf("albumID is empty")
	}
	token := s.mainToken()
	if token == "" {
		return nil, fmt.Errorf("facebook main page token is empty")
	}

	nextURL := fmt.Sprintf("%s%s/photos?fields=id&limit=100&access_token=%s",
		s.apiURL(), url.PathEscape(albumID), url.QueryEscape(token))

	var ids []string
	for nextURL != "" {
		req, err := http.NewRequest("GET", nextURL, nil)
		if err != nil {
			return nil, err
		}
		resp, err := s.Client.Do(req)
		if err != nil {
			return nil, err
		}
		body, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			return nil, readErr
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return nil, fmt.Errorf("facebook list album photos error (%d): %s", resp.StatusCode, string(body))
		}

		var res struct {
			Data []struct {
				ID string `json:"id"`
			} `json:"data"`
			Paging struct {
				Next string `json:"next"`
			} `json:"paging"`
		}
		if err := json.Unmarshal(body, &res); err != nil {
			return nil, err
		}
		for _, item := range res.Data {
			if item.ID != "" {
				ids = append(ids, item.ID)
			}
		}
		nextURL = res.Paging.Next
	}
	return ids, nil
}

// StartCategoryAlbumResync queues a resumable album resync job for a specific category.
// Products are re-uploaded in Newest -> Oldest order so that newest products appear at the top of the Facebook album.
func (s *FacebookService) StartCategoryAlbumResync(categoryCode string, triggerType string) (*models.AlbumResyncJob, error) {
	if s == nil || s.DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	categoryCode = strings.TrimSpace(categoryCode)
	if categoryCode == "" {
		return nil, fmt.Errorf("category_code is required")
	}
	if triggerType == "" {
		triggerType = "MANUAL"
	}

	// Return existing active/paused job if one is already in progress for this category
	var existing models.AlbumResyncJob
	if err := s.DB.Where("category_code = ? AND status IN ('PENDING', 'DELETING', 'UPLOADING', 'PAUSED_LIMIT')", categoryCode).
		Order("id DESC").First(&existing).Error; err == nil {
		return &existing, nil
	}

	var cat models.Category
	if err := s.DB.Where("code = ? AND is_active = true", categoryCode).First(&cat).Error; err != nil {
		return nil, fmt.Errorf("category %s not found or inactive", categoryCode)
	}
	if cat.FBAlbumID == nil || *cat.FBAlbumID == "" {
		return nil, fmt.Errorf("category %s does not have an associated Facebook album", categoryCode)
	}
	albumID := *cat.FBAlbumID

	// 1. Fetch existing photos in the Facebook album to clear them completely (including duplicates/orphans)
	deleteIDs, err := s.FetchAlbumPhotoIDs(albumID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch existing album photos: %w", err)
	}

	// 2. Query AVAILABLE products sorted Newest -> Oldest (so newest are uploaded first and stay at top of album)
	var products []models.Product
	subQuery := "(SELECT product_id, DATE_FORMAT(MAX(created_at), '%Y-%m-%d %H:%i:%s') as latest_upload FROM product_upload_logs WHERE product_id IS NOT NULL GROUP BY product_id) as upload_summary"
	err = s.DB.Model(&models.Product{}).
		Joins("LEFT JOIN "+subQuery+" ON upload_summary.product_id = products.id").
		Preload("Images").
		Where("products.type = ? AND (products.status = 'AVAILABLE' OR products.status = '' OR products.status IS NULL)", cat.Code).
		Where(`(products.path_thumb IS NOT NULL AND products.path_thumb != '') 
			OR (products.image_thumb_scale IS NOT NULL AND products.image_thumb_scale != '') 
			OR EXISTS (SELECT 1 FROM product_images WHERE product_images.product_id = products.id AND product_images.file_path != '')`).
		Select("products.*").
		Order("COALESCE(upload_summary.latest_upload, products.created_at) DESC, products.id DESC").
		Find(&products).Error
	if err != nil {
		return nil, fmt.Errorf("failed to query available products: %w", err)
	}

	var uploadIDs []uint
	for _, p := range products {
		if p.BestQualityThumb() != "" {
			uploadIDs = append(uploadIDs, p.ID)
		}
	}

	delJSON, _ := json.Marshal(deleteIDs)
	upJSON, _ := json.Marshal(uploadIDs)
	now := time.Now()
	phase := "DELETE"
	if len(deleteIDs) == 0 {
		phase = "UPLOAD"
	}

	job := &models.AlbumResyncJob{
		CategoryCode:       cat.Code,
		CategoryName:       cat.Name,
		AlbumID:            albumID,
		TriggerType:        triggerType,
		Status:             "PENDING",
		Phase:              phase,
		TotalDelete:        len(deleteIDs),
		DeletedCount:       0,
		TotalUpload:        len(uploadIDs),
		UploadedCount:      0,
		RemainingDeleteIDs: string(delJSON),
		RemainingUploadIDs: string(upJSON),
		StartedAt:          &now,
	}

	if err := s.DB.Create(job).Error; err != nil {
		return nil, err
	}

	s.DB.Model(&models.Category{}).Where("id = ?", cat.ID).Update("needs_album_reorder", false)

	// Wake up worker immediately
	if s.resyncWakeup != nil {
		select {
		case s.resyncWakeup <- struct{}{}:
		default:
		}
	}

	return job, nil
}

func (s *FacebookService) CancelAlbumResyncJob(jobID uint) error {
	if s == nil || s.DB == nil {
		return fmt.Errorf("database not initialized")
	}
	res := s.DB.Model(&models.AlbumResyncJob{}).
		Where("id = ? AND status IN ('PENDING', 'DELETING', 'UPLOADING', 'PAUSED_LIMIT')", jobID).
		Updates(map[string]interface{}{
			"status":     "CANCELLED",
			"last_error": "Đã hủy bởi người dùng",
		})
	return res.Error
}

func (s *FacebookService) ForceResumeAlbumResyncJob(jobID uint) error {
	if s == nil || s.DB == nil {
		return fmt.Errorf("database not initialized")
	}
	now := time.Now()
	res := s.DB.Model(&models.AlbumResyncJob{}).
		Where("id = ? AND status = 'PAUSED_LIMIT'", jobID).
		Updates(map[string]interface{}{
			"resume_at":  &now,
			"last_error": "Đang tiếp tục chạy theo yêu cầu...",
		})
	if res.Error != nil {
		return res.Error
	}
	if s.resyncWakeup != nil {
		select {
		case s.resyncWakeup <- struct{}{}:
		default:
		}
	}
	return nil
}

func (s *FacebookService) StartAlbumResyncWorker() {
	if s == nil {
		return
	}
	// Wait briefly on startup before resuming any interrupted jobs
	time.Sleep(15 * time.Second)
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	gmt7 := time.FixedZone("GMT+7", 7*3600)
	lastNightlyDate := ""

	for {
		// 1. Check 1:00 AM GMT+7 Nightly Trigger
		nowGMT7 := time.Now().In(gmt7)
		todayStr := nowGMT7.Format("2006-01-02")
		if nowGMT7.Hour() == 1 && nowGMT7.Minute() < 10 && lastNightlyDate != todayStr {
			lastNightlyDate = todayStr
			s.triggerNightlyResyncs()
		}

		// 2. Process any pending or ready-to-resume resync job
		s.ProcessNextAlbumResyncJob()

		select {
		case <-ticker.C:
		case <-s.resyncWakeup:
		}
	}
}

func (s *FacebookService) triggerNightlyResyncs() {
	if s == nil || s.DB == nil || (s.PageToken == "" && s.MainPageToken == "") {
		return
	}
	var cats []models.Category
	if err := s.DB.Where("is_active = true AND fb_album_id IS NOT NULL AND fb_album_id != '' AND needs_album_reorder = true").
		Order("id ASC").Find(&cats).Error; err != nil {
		log.Printf("NightlyResync query error: %v\n", err)
		return
	}
	if len(cats) == 0 {
		return
	}
	log.Printf("NightlyResync (01:00 GMT+7): queueing %d categories with needs_album_reorder=true\n", len(cats))
	for _, cat := range cats {
		if _, err := s.StartCategoryAlbumResync(cat.Code, "NIGHTLY"); err != nil {
			log.Printf("NightlyResync error queueing %s: %v\n", cat.Code, err)
		}
	}
}

func (s *FacebookService) ProcessNextAlbumResyncJob() {
	if s == nil || s.DB == nil || (s.PageToken == "" && s.MainPageToken == "") {
		return
	}
	s.resyncMu.Lock()
	defer s.resyncMu.Unlock()

	var job models.AlbumResyncJob
	now := time.Now()

	// Respect app-level Meta rate-limit cooldown: do not start any other job while cooling down
	var coolingCount int64
	s.DB.Model(&models.AlbumResyncJob{}).
		Where("status = 'PAUSED_LIMIT' AND resume_at > ?", now).
		Count(&coolingCount)
	if coolingCount > 0 {
		return
	}

	err := s.DB.Where("status IN ('PENDING', 'DELETING', 'UPLOADING') OR (status = 'PAUSED_LIMIT' AND resume_at <= ?)", now).
		Order("id ASC").
		First(&job).Error
	if err != nil {
		return
	}

	s.processAlbumResyncJob(&job)
}

func (s *FacebookService) processAlbumResyncJob(job *models.AlbumResyncJob) {
	var deleteIDs []string
	var uploadIDs []uint
	if job.RemainingDeleteIDs != "" {
		_ = json.Unmarshal([]byte(job.RemainingDeleteIDs), &deleteIDs)
	}
	if job.RemainingUploadIDs != "" {
		_ = json.Unmarshal([]byte(job.RemainingUploadIDs), &uploadIDs)
	}

	throttle := 1 * time.Second
	if s.ThrottleDelay > 0 {
		throttle = s.ThrottleDelay
	}

	// Phase 1: DELETE existing photos from album
	if len(deleteIDs) > 0 {
		job.Phase = "DELETE"
		job.Status = "DELETING"
		job.LastError = ""
		job.ResumeAt = nil
		s.saveResyncProgress(job, deleteIDs, uploadIDs)

		for len(deleteIDs) > 0 {
			if s.isResyncJobCancelled(job.ID) {
				return
			}
			photoID := deleteIDs[0]
			usagePct, err := s.RemoveProductPhotoFromAlbumWithUsage(photoID)
			if err != nil && isFBRateLimitError(err) {
				s.pauseAlbumResyncJob(job, deleteIDs, uploadIDs, fmt.Sprintf("Facebook rate limit (#4) khi đang xóa ảnh. Tự động chạy tiếp sau 45 phút."))
				return
			}

			deleteIDs = deleteIDs[1:]
			job.DeletedCount++
			s.DB.Model(&models.Product{}).Where("fb_album_photo_id = ?", photoID).Update("fb_album_photo_id", nil)

			if job.DeletedCount%5 == 0 || len(deleteIDs) == 0 {
				s.saveResyncProgress(job, deleteIDs, uploadIDs)
			}

			if usagePct >= resyncUsageThresholdPct && len(deleteIDs) > 0 {
				s.pauseAlbumResyncJob(job, deleteIDs, uploadIDs, fmt.Sprintf("Tạm nghỉ chủ động: Meta API đạt %d%% giới hạn. Tự động chạy tiếp sau 45 phút.", usagePct))
				return
			}
			time.Sleep(throttle)
		}
	}

	// Phase 2: UPLOAD available products in Newest -> Oldest order
	if len(uploadIDs) > 0 {
		job.Phase = "UPLOAD"
		job.Status = "UPLOADING"
		job.LastError = ""
		job.ResumeAt = nil
		s.saveResyncProgress(job, deleteIDs, uploadIDs)

		for len(uploadIDs) > 0 {
			if s.isResyncJobCancelled(job.ID) {
				return
			}
			prodID := uploadIDs[0]
			var product models.Product
			if err := s.DB.Preload("Sizes").Preload("Images").First(&product, prodID).Error; err != nil {
				uploadIDs = uploadIDs[1:]
				continue
			}
			if product.Status != nil && *product.Status != "" && *product.Status != "AVAILABLE" {
				uploadIDs = uploadIDs[1:]
				continue
			}
			thumb := product.BestQualityThumb()
			if thumb == "" {
				uploadIDs = uploadIDs[1:]
				continue
			}

			catalogPostID := ""
			if product.FBPostID != nil {
				catalogPostID = *product.FBPostID
			}
			caption := s.BuildAlbumPhotoCaption(&product, catalogPostID)
			photoID, usagePct, err := s.UploadProductPhotoToAlbumWithUsage(job.AlbumID, thumb, caption)
			if err != nil {
				if isFBRateLimitError(err) {
					s.pauseAlbumResyncJob(job, deleteIDs, uploadIDs, "Facebook rate limit (#4) khi đang tải ảnh lên. Tự động chạy tiếp sau 45 phút.")
					return
				}
				log.Printf("AlbumResync [%s] error uploading product %d: %v\n", job.CategoryCode, prodID, err)
				job.LastError = err.Error()
				uploadIDs = uploadIDs[1:]
				continue
			}

			if photoID != "" {
				s.DB.Model(&models.Product{}).Where("id = ?", prodID).Update("fb_album_photo_id", photoID)
			}
			uploadIDs = uploadIDs[1:]
			job.UploadedCount++

			if job.UploadedCount%3 == 0 || len(uploadIDs) == 0 {
				s.saveResyncProgress(job, deleteIDs, uploadIDs)
			}

			if usagePct >= resyncUsageThresholdPct && len(uploadIDs) > 0 {
				s.pauseAlbumResyncJob(job, deleteIDs, uploadIDs, fmt.Sprintf("Tạm nghỉ chủ động: Meta API đạt %d%% giới hạn. Tự động chạy tiếp sau 45 phút.", usagePct))
				return
			}
			time.Sleep(throttle * 2)
		}
	}

	now := time.Now()
	job.Status = "COMPLETED"
	job.LastError = ""
	job.ResumeAt = nil
	job.CompletedAt = &now
	s.saveResyncProgress(job, deleteIDs, uploadIDs)
	log.Printf("AlbumResync [%s] COMPLETED: deleted=%d, uploaded=%d\n", job.CategoryCode, job.DeletedCount, job.UploadedCount)
}

func (s *FacebookService) isResyncJobCancelled(jobID uint) bool {
	var status string
	if err := s.DB.Model(&models.AlbumResyncJob{}).Where("id = ?", jobID).Pluck("status", &status).Error; err == nil {
		return status == "CANCELLED"
	}
	return false
}

func (s *FacebookService) saveResyncProgress(job *models.AlbumResyncJob, deleteIDs []string, uploadIDs []uint) {
	delJSON, _ := json.Marshal(deleteIDs)
	upJSON, _ := json.Marshal(uploadIDs)
	job.RemainingDeleteIDs = string(delJSON)
	job.RemainingUploadIDs = string(upJSON)
	s.DB.Model(&models.AlbumResyncJob{}).Where("id = ?", job.ID).Updates(map[string]interface{}{
		"status":               job.Status,
		"phase":                job.Phase,
		"deleted_count":        job.DeletedCount,
		"uploaded_count":       job.UploadedCount,
		"remaining_delete_ids": job.RemainingDeleteIDs,
		"remaining_upload_ids": job.RemainingUploadIDs,
		"resume_at":            job.ResumeAt,
		"last_error":           job.LastError,
		"completed_at":         job.CompletedAt,
	})
}

func (s *FacebookService) pauseAlbumResyncJob(job *models.AlbumResyncJob, deleteIDs []string, uploadIDs []uint, reason string) {
	resumeAt := time.Now().Add(resyncPauseDuration)
	job.Status = "PAUSED_LIMIT"
	job.ResumeAt = &resumeAt
	job.LastError = reason
	s.saveResyncProgress(job, deleteIDs, uploadIDs)
	log.Printf("AlbumResync [%s] PAUSED until %s: %s\n", job.CategoryCode, resumeAt.Format("15:04:05"), reason)
}
