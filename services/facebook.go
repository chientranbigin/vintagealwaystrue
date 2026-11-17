package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"sale-app/models"

	"gorm.io/gorm"
)

type FacebookSyncAction string

const (
	FBActionCreate            FacebookSyncAction = "create"
	FBActionCreateCatalogOnly FacebookSyncAction = "create_catalog_only"
	FBActionSold              FacebookSyncAction = "sold"
	FBActionRevert            FacebookSyncAction = "revert"
	FBActionHold              FacebookSyncAction = "hold"
)

type FacebookSyncJob struct {
	ProductID uint
	Action    FacebookSyncAction
	Retries   int
}

type FacebookAlbum struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Count       int    `json:"count,omitempty"`
	CreatedTime string `json:"created_time,omitempty"`
}

type FacebookService struct {
	PageID           string
	PageToken        string
	MainPageID       string
	MainPageToken    string
	MainPageToken2   string
	BaseURL          string
	Client           *http.Client
	DB               *gorm.DB
	MediaBase        string
	Queue            chan FacebookSyncJob
	ThrottleDelay    time.Duration
	resyncWakeup     chan struct{}
	resyncMu         sync.Mutex
	appMu            sync.Mutex
	appCooldownUntil [3]time.Time
}

var TrouserSizeMap = []struct {
	Size int
	Min  int
	Max  int
}{
	{28, 74, 76},
	{29, 76, 78},
	{30, 78, 80},
	{31, 80, 82},
	{32, 82, 84},
	{33, 84, 86},
	{34, 86, 88},
	{35, 88, 90},
	{36, 90, 92},
	{37, 92, 94},
	{38, 94, 96},
	{39, 96, 98},
	{40, 98, 100},
}

func GetTrouserSize(eoActual int) *int {
	for _, item := range TrouserSizeMap {
		if eoActual >= item.Min && eoActual < item.Max {
			s := item.Size
			return &s
		}
	}
	last := TrouserSizeMap[len(TrouserSizeMap)-1]
	if eoActual == last.Max {
		s := last.Size
		return &s
	}
	return nil
}

func FormatVNDPrice(amount int64) string {
	str := strconv.FormatInt(amount, 10)
	n := len(str)
	if n <= 3 {
		return str + "đ"
	}
	var res []string
	lead := n % 3
	if lead > 0 {
		res = append(res, str[:lead])
	}
	for i := lead; i < n; i += 3 {
		res = append(res, str[i:i+3])
	}
	return strings.Join(res, ".") + "đ"
}

func NewFacebookService(db *gorm.DB, pageID string, pageToken string, mediaBase string, mainPageArgs ...string) *FacebookService {
	if pageID == "" {
		pageID = os.Getenv("FB_PAGE_ID")
	}
	if pageToken == "" {
		pageToken = os.Getenv("FB_PAGE_TOKEN")
	}
	mainPageID := ""
	mainPageToken := ""
	if len(mainPageArgs) > 0 {
		mainPageID = mainPageArgs[0]
	}
	if len(mainPageArgs) > 1 {
		mainPageToken = mainPageArgs[1]
	}
	if mainPageID == "" {
		mainPageID = os.Getenv("FB_MAIN_PAGE_ID")
	}
	if mainPageToken == "" {
		mainPageToken = os.Getenv("FB_MAIN_PAGE_TOKEN")
	}

	mainPageToken2 := os.Getenv("FB_MAIN_PAGE_TOKEN_2")

	if mediaBase == "" {
		mediaBase = os.Getenv("MEDIA_BASE_URL")
		if mediaBase == "" {
			mediaBase = "https://res.cloudinary.com/leq5n17r/image/upload/"
		}
	}

	fb := &FacebookService{
		PageID:         pageID,
		PageToken:      pageToken,
		MainPageID:     mainPageID,
		MainPageToken:  mainPageToken,
		MainPageToken2: mainPageToken2,
		BaseURL:       "https://graph.facebook.com/v19.0/",
		Client:        &http.Client{Timeout: 30 * time.Second},
		DB:            db,
		MediaBase:     mediaBase,
		Queue:         make(chan FacebookSyncJob, 200),
		ThrottleDelay: 1 * time.Second,
		resyncWakeup:  make(chan struct{}, 10),
	}

	fb.StartWorker()
	return fb
}

func (s *FacebookService) mainToken() string {
	if s != nil && s.MainPageToken != "" {
		return s.MainPageToken
	}
	if s != nil {
		return s.PageToken
	}
	return ""
}

func isFBRateLimitError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "(#4)") || strings.Contains(msg, "request limit reached")
}

var backgroundWorkersOnce sync.Once

func (s *FacebookService) StartWorker() {
	go func() {
		for job := range s.Queue {
			s.handleJob(job)
		}
	}()

	backgroundWorkersOnce.Do(func() {
		go s.StartReconciler()
		go s.StartAlbumResyncWorker()
	})
}

func (s *FacebookService) StartReconciler() {
	// Wait 1 minute after server start before initial reconcile
	time.Sleep(1 * time.Minute)
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	s.ReconcileAlbumPhotos()

	for range ticker.C {
		s.ReconcileAlbumPhotos()
	}
}

func (s *FacebookService) ReconcileAlbumPhotos() {
	if s == nil || s.DB == nil || (s.PageToken == "" && s.MainPageToken == "") {
		return
	}

	var staleProducts []models.Product
	if err := s.DB.Where("status IN ('SOLD', 'ON_HOLD') AND fb_album_photo_id IS NOT NULL AND fb_album_photo_id != ''").
		Order("updated_at DESC").
		Limit(30).
		Find(&staleProducts).Error; err != nil {
		log.Printf("ReconcileAlbumPhotos query error: %v\n", err)
		return
	}

	if len(staleProducts) == 0 {
		return
	}

	log.Printf("ReconcileAlbumPhotos: found %d SOLD/ON_HOLD products with album photos to clean up\n", len(staleProducts))
	for _, p := range staleProducts {
		if p.FBAlbumPhotoID == nil || *p.FBAlbumPhotoID == "" {
			continue
		}
		photoID := *p.FBAlbumPhotoID
		pName := fmt.Sprintf("ID:%d", p.ID)
		if p.Name != nil && *p.Name != "" {
			pName = *p.Name
		}

		_, err := s.RemoveAlbumPhotoWithFallback(photoID, p.FBAlbumApp)
		if err != nil {
			if isFBRateLimitError(err) {
				log.Printf("ReconcileAlbumPhotos: Facebook rate limit reached (code 4), pausing until next cycle\n")
				return
			}
			log.Printf("ReconcileAlbumPhotos error deleting photo %s for product %s: %v\n", photoID, pName, err)
		} else {
			s.DB.Model(&models.Product{}).Where("id = ?", p.ID).Update("fb_album_photo_id", nil)
			log.Printf("ReconcileAlbumPhotos: successfully cleaned up album photo %s for product %s\n", photoID, pName)
		}
		delay := 1 * time.Second
		if s.ThrottleDelay > 0 {
			delay = s.ThrottleDelay
		}
		time.Sleep(delay)
	}
}

func (s *FacebookService) Dispatch(productID uint, action FacebookSyncAction) {
	if s == nil || productID == 0 {
		return
	}
	job := FacebookSyncJob{
		ProductID: productID,
		Action:    action,
		Retries:   0,
	}
	select {
	case s.Queue <- job:
	default:
		go func() { s.Queue <- job }()
	}
}

func (s *FacebookService) handleJob(job FacebookSyncJob) {
	if s.PageToken == "" && s.MainPageToken == "" {
		log.Println("SyncFacebookPost: no FB tokens configured, skipping.")
		return
	}

	delay := 1 * time.Second
	if s.ThrottleDelay > 0 {
		delay = s.ThrottleDelay
	}
	time.Sleep(delay)

	err := s.ExecuteSync(job.ProductID, job.Action)
	if err != nil {
		log.Printf("SyncFacebookPost [%s] error for product ID %d: %v\n", job.Action, job.ProductID, err)
		if isFBRateLimitError(err) {
			log.Printf("SyncFacebookPost [%s] rate limited by Facebook. Background Reconciler will retry automatically.\n", job.Action)
			return
		}
		if job.Retries < 3 {
			go func(j FacebookSyncJob) {
				delay := time.Duration((j.Retries+1)*5) * time.Second
				time.Sleep(delay)
				j.Retries++
				select {
				case s.Queue <- j:
				default:
					go func() { s.Queue <- j }()
				}
			}(job)
		}
	}
}

func (s *FacebookService) ExecuteSync(productID uint, action FacebookSyncAction) error {
	if s.DB == nil {
		return fmt.Errorf("database not initialized in FacebookService")
	}

	var product models.Product
	if err := s.DB.Preload("Images").Preload("Sizes").First(&product, productID).Error; err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	pName := ""
	if product.Name != nil {
		pName = *product.Name
	}

	switch action {
	case FBActionCreateCatalogOnly:
		if product.FBPostID != nil && *product.FBPostID != "" {
			_ = s.DeletePost(*product.FBPostID)
			log.Printf("FB post deleted before re-create: %s for %s\n", *product.FBPostID, pName)
		}

		postID, err := s.CreatePost(&product)
		if err != nil {
			log.Printf("CreatePost error on Catalog Page for product %d: %v\n", product.ID, err)
		} else if postID != "" {
			s.DB.Model(&models.Product{}).Where("id = ?", product.ID).Update("fb_post_id", postID)
			product.FBPostID = &postID
			log.Printf("FB post created (Catalog only, skipping Album A): %s for %s\n", postID, pName)
		}

		// Mark category as needing album reorder (user can trigger Re-sync from UI or let 1 AM nightly job handle it)
		if product.Type != nil && *product.Type != "" {
			s.DB.Model(&models.Category{}).Where("code = ? AND is_active = true", *product.Type).Update("needs_album_reorder", true)
		}

	case FBActionCreate:
		if product.FBPostID != nil && *product.FBPostID != "" {
			_ = s.DeletePost(*product.FBPostID)
			log.Printf("FB post deleted before re-create: %s for %s\n", *product.FBPostID, pName)
		}

		postID, err := s.CreatePost(&product)
		if err != nil {
			log.Printf("CreatePost error on Catalog Page for product %d: %v\n", product.ID, err)
		} else if postID != "" {
			s.DB.Model(&models.Product{}).Where("id = ?", product.ID).Update("fb_post_id", postID)
			product.FBPostID = &postID
			log.Printf("FB post created: %s for %s\n", postID, pName)
		}

		// Sync thumbnail photo into Category Album on Main Page (no_story=true)
		if product.Type != nil {
			var cat models.Category
			if err := s.DB.Where("code = ? AND is_active = true", *product.Type).First(&cat).Error; err == nil && cat.FBAlbumID != nil && *cat.FBAlbumID != "" {
				// Delete previous album photo if existing
				if product.FBAlbumPhotoID != nil && *product.FBAlbumPhotoID != "" {
					if _, err := s.RemoveAlbumPhotoWithFallback(*product.FBAlbumPhotoID, product.FBAlbumApp); err != nil {
						log.Printf("FB album photo remove previous error (%s): %v\n", *product.FBAlbumPhotoID, err)
						if isFBRateLimitError(err) {
							return err // Wait for rate limit to clear before uploading new photo, preventing duplicate album photos
						}
					}
				}

				thumb := product.BestQualityThumb()
				if thumb != "" {
					catalogPostID := ""
					if product.FBPostID != nil {
						catalogPostID = *product.FBPostID
					}
					caption := s.BuildAlbumPhotoCaption(&product, catalogPostID)
					photoID, usedApp, _, albumErr := s.UploadAlbumPhotoWithFailover(*cat.FBAlbumID, thumb, caption)
					if albumErr == nil && photoID != "" {
						s.DB.Model(&models.Product{}).Where("id = ?", product.ID).Updates(map[string]interface{}{
							"fb_album_photo_id": photoID,
							"fb_album_app":      usedApp,
						})
						s.DB.Model(&models.Category{}).Where("id = ?", cat.ID).Update("needs_album_reorder", true)
						product.FBAlbumPhotoID = &photoID
						product.FBAlbumApp = usedApp
						log.Printf("FB Album photo uploaded (App %d): %s in album %s for %s\n", usedApp, photoID, *cat.FBAlbumID, pName)
					} else {
						log.Printf("FB Album photo upload error for %s in album %s: %v\n", pName, *cat.FBAlbumID, albumErr)
					}
				}
			}
		}

	case FBActionSold:
		var lastErr error
		// 1. Mark SOLD on Catalog Page B
		if product.FBPostID != nil && *product.FBPostID != "" {
			if err := s.MarkSold(*product.FBPostID, &product); err != nil {
				log.Printf("FB post MarkSold error for %s: %v\n", *product.FBPostID, err)
				lastErr = err
			} else {
				log.Printf("FB post marked SOLD: %s\n", *product.FBPostID)
			}
		}

		// 2. Remove photo from Main Page Album A
		if product.FBAlbumPhotoID != nil && *product.FBAlbumPhotoID != "" {
			if _, err := s.RemoveAlbumPhotoWithFallback(*product.FBAlbumPhotoID, product.FBAlbumApp); err != nil {
				log.Printf("FB album photo remove error (%s): %v\n", *product.FBAlbumPhotoID, err)
				lastErr = err
			} else {
				s.DB.Model(&models.Product{}).Where("id = ?", product.ID).Update("fb_album_photo_id", nil)
				log.Printf("FB album photo removed: %s for %s\n", *product.FBAlbumPhotoID, pName)
			}
		}
		if lastErr != nil {
			return lastErr
		}

	case FBActionHold:
		// Remove photo from Main Page Album A
		if product.FBAlbumPhotoID != nil && *product.FBAlbumPhotoID != "" {
			if _, err := s.RemoveAlbumPhotoWithFallback(*product.FBAlbumPhotoID, product.FBAlbumApp); err != nil {
				log.Printf("FB album photo remove on HOLD error (%s): %v\n", *product.FBAlbumPhotoID, err)
				return err
			} else {
				s.DB.Model(&models.Product{}).Where("id = ?", product.ID).Update("fb_album_photo_id", nil)
				log.Printf("FB album photo removed on HOLD: %s for %s\n", *product.FBAlbumPhotoID, pName)
			}
		}

	case FBActionRevert:
		var lastErr error
		// 1. Revert AVAILABLE on Catalog Page B
		if product.FBPostID != nil && *product.FBPostID != "" {
			if err := s.RevertAvailable(*product.FBPostID, &product); err != nil {
				log.Printf("FB post RevertAvailable error for %s: %v\n", *product.FBPostID, err)
				lastErr = err
			} else {
				log.Printf("FB post reverted AVAILABLE: %s\n", *product.FBPostID)
			}
		}

		// 2. Re-upload photo to Main Page Album A if missing
		if (product.FBAlbumPhotoID == nil || *product.FBAlbumPhotoID == "") && product.Type != nil {
			var cat models.Category
			if err := s.DB.Where("code = ? AND is_active = true", *product.Type).First(&cat).Error; err == nil && cat.FBAlbumID != nil && *cat.FBAlbumID != "" {
				thumb := product.BestQualityThumb()
				if thumb != "" {
					catalogPostID := ""
					if product.FBPostID != nil {
						catalogPostID = *product.FBPostID
					}
					caption := s.BuildAlbumPhotoCaption(&product, catalogPostID)
					photoID, usedApp, _, albumErr := s.UploadAlbumPhotoWithFailover(*cat.FBAlbumID, thumb, caption)
					if albumErr != nil {
						log.Printf("FB Album photo re-upload error for %s: %v\n", pName, albumErr)
						lastErr = albumErr
					} else if photoID != "" {
						s.DB.Model(&models.Product{}).Where("id = ?", product.ID).Updates(map[string]interface{}{
							"fb_album_photo_id": photoID,
							"fb_album_app":      usedApp,
						})
						s.DB.Model(&models.Category{}).Where("id = ?", cat.ID).Update("needs_album_reorder", true)
						log.Printf("FB Album photo re-uploaded (App %d): %s in album %s for %s\n", usedApp, photoID, *cat.FBAlbumID, pName)
					}
				}
			}
		}
		if lastErr != nil {
			return lastErr
		}
	}

	return nil
}

func (s *FacebookService) apiURL() string {
	base := s.BaseURL
	if base == "" {
		base = "https://graph.facebook.com/v19.0/"
	}
	if !strings.HasSuffix(base, "/") {
		base += "/"
	}
	return base
}

func (s *FacebookService) uploadPhoto(imagePath string) (string, error) {
	if strings.HasPrefix(imagePath, "http://") || strings.HasPrefix(imagePath, "https://") {
		// Attempt to upload by URL directly
		endpoint := fmt.Sprintf("%s%s/photos", s.apiURL(), s.PageID)
		data := url.Values{}
		data.Set("url", imagePath)
		data.Set("published", "false")
		data.Set("access_token", s.PageToken)

		resp, err := s.Client.PostForm(endpoint, data)
		if err == nil {
			defer resp.Body.Close()
			respBytes, rErr := io.ReadAll(resp.Body)
			if rErr == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
				var res struct {
					ID string `json:"id"`
				}
				if json.Unmarshal(respBytes, &res) == nil && res.ID != "" {
					return res.ID, nil
				}
			}
		}

		// If URL upload was rejected, download image bytes and upload via multipart source
		dlResp, dlErr := s.Client.Get(imagePath)
		if dlErr != nil {
			return "", dlErr
		}
		defer dlResp.Body.Close()
		return s.uploadPhotoFromReader(dlResp.Body, filepath.Base(imagePath))
	}

	// Local file path handling
	localPath := imagePath
	if !filepath.IsAbs(localPath) {
		if _, err := os.Stat(localPath); os.IsNotExist(err) {
			clean := strings.TrimPrefix(localPath, "/")
			if _, err2 := os.Stat(clean); err2 == nil {
				localPath = clean
			} else if _, err3 := os.Stat("storage/app/public/" + strings.TrimPrefix(clean, "storage/")); err3 == nil {
				localPath = "storage/app/public/" + strings.TrimPrefix(clean, "storage/")
			} else if _, err4 := os.Stat("./storage/" + strings.TrimPrefix(clean, "storage/")); err4 == nil {
				localPath = "./storage/" + strings.TrimPrefix(clean, "storage/")
			}
		}
	}

	if f, err := os.Open(localPath); err == nil {
		defer f.Close()
		return s.uploadPhotoFromReader(f, filepath.Base(localPath))
	}

	// If not found locally and MediaBase is available, try CDN URL
	if s.MediaBase != "" {
		remoteURL := s.MediaBase
		if !strings.HasSuffix(remoteURL, "/") {
			remoteURL += "/"
		}
		remoteURL += strings.TrimPrefix(imagePath, "storage/")
		return s.uploadPhoto(remoteURL)
	}

	return "", fmt.Errorf("image not accessible at path %s", imagePath)
}

func (s *FacebookService) uploadPhotoFromReader(r io.Reader, filename string) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("published", "false")
	_ = writer.WriteField("access_token", s.PageToken)
	part, err := writer.CreateFormFile("source", filename)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(part, r); err != nil {
		return "", err
	}
	_ = writer.Close()

	endpoint := fmt.Sprintf("%s%s/photos", s.apiURL(), s.PageID)
	req, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := s.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("facebook photo upload error (%d): %s", resp.StatusCode, string(respBytes))
	}

	var res struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(respBytes, &res); err != nil {
		return "", err
	}
	return res.ID, nil
}

func (s *FacebookService) CreatePost(product *models.Product) (string, error) {
	// 1. Collect image paths (thumb first, then detail images without duplicates)
	var imagePaths []string
	seen := make(map[string]bool)

	if product.PathThumb != nil && *product.PathThumb != "" {
		p := *product.PathThumb
		imagePaths = append(imagePaths, p)
		seen[p] = true
	}

	for _, img := range product.Images {
		if img.FilePath != "" && !seen[img.FilePath] {
			seen[img.FilePath] = true
			imagePaths = append(imagePaths, img.FilePath)
		}
	}

	if len(imagePaths) == 0 {
		return "", nil
	}

	// 2. Upload each image as unpublished photo
	type attachedMedia struct {
		MediaFBID string `json:"media_fbid"`
	}
	var mediaList []attachedMedia

	for _, path := range imagePaths {
		mediaID, err := s.uploadPhoto(path)
		if err == nil && mediaID != "" {
			mediaList = append(mediaList, attachedMedia{MediaFBID: mediaID})
		}
	}

	if len(mediaList) == 0 {
		return "", fmt.Errorf("no photos uploaded successfully for product %d", product.ID)
	}

	// 3. Create post with all photos
	message := s.BuildMessage(product, "AVAILABLE")
	mediaJSON, err := json.Marshal(mediaList)
	if err != nil {
		return "", err
	}

	endpoint := fmt.Sprintf("%s%s/feed", s.apiURL(), s.PageID)
	data := url.Values{}
	data.Set("message", message)
	data.Set("attached_media", string(mediaJSON))
	data.Set("access_token", s.PageToken)

	resp, err := s.Client.PostForm(endpoint, data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("facebook feed post error (%d): %s", resp.StatusCode, string(respBytes))
	}

	var res struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(respBytes, &res); err != nil {
		return "", err
	}

	return res.ID, nil
}

func (s *FacebookService) DeletePost(fbPostId string) error {
	endpoint := fmt.Sprintf("%s%s?access_token=%s", s.apiURL(), url.PathEscape(fbPostId), url.QueryEscape(s.PageToken))
	req, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// Post may already be deleted/inaccessible — ignore and proceed.
	return nil
}

func (s *FacebookService) MarkSold(fbPostId string, product *models.Product) error {
	message := s.BuildMessage(product, "SOLD")
	return s.updatePostMessage(fbPostId, message)
}

func (s *FacebookService) RevertAvailable(fbPostId string, product *models.Product) error {
	message := s.BuildMessage(product, "AVAILABLE")
	return s.updatePostMessage(fbPostId, message)
}

func (s *FacebookService) updatePostMessage(fbPostId string, message string) error {
	endpoint := fmt.Sprintf("%s%s", s.apiURL(), fbPostId)
	data := url.Values{}
	data.Set("message", message)
	data.Set("access_token", s.PageToken)

	resp, err := s.Client.PostForm(endpoint, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("facebook update post error (%d): %s", resp.StatusCode, string(respBytes))
	}

	return nil
}

func (s *FacebookService) BuildHashtags(product *models.Product, status string) string {
	var tags []string

	pName := ""
	if product.Name != nil {
		pName = *product.Name
	}
	tags = append(tags, "#"+pName)

	pType := ""
	if product.Type != nil {
		pType = *product.Type
	}
	statusTag := pType + "_AVAILABLE"
	if status == "SOLD" {
		statusTag = pType + "_SOLD"
	}
	tags = append(tags, "#"+statusTag)

	if pType == "TROUSERS" {
		var eoSize *models.Size
		for i := range product.Sizes {
			sz := &product.Sizes[i]
			if sz.Name != nil && (*sz.Name == "EO" || *sz.Name == "EO QUẦN") {
				eoSize = sz
				break
			}
		}
		if eoSize != nil && eoSize.Value != nil {
			eoDb := int(*eoSize.Value)
			eoActual := eoDb * 2
			tags = append(tags, fmt.Sprintf("#EO_%d", eoActual))
			tags = append(tags, fmt.Sprintf("#EO_%d_%d", eoActual, eoActual+2))
			trouserSize := GetTrouserSize(eoActual)
			if trouserSize != nil {
				tags = append(tags, fmt.Sprintf("#SIZE%d", *trouserSize))
			}
		}
	} else if pType == "BLAZER" || pType == "SUIT" || pType == "SHIRT" || pType == "GILE" {
		var vaiSize *models.Size
		for i := range product.Sizes {
			sz := &product.Sizes[i]
			if sz.Name != nil && *sz.Name == "VAI" {
				vaiSize = sz
				break
			}
		}
		if vaiSize != nil && vaiSize.Value != nil {
			vai := int(*vaiSize.Value)
			tags = append(tags, fmt.Sprintf("#VAI_%d", vai))
			start := vai
			if vai%2 != 0 {
				start = vai - 1
			}
			tags = append(tags, fmt.Sprintf("#VAI_%d_%d", start, start+1))
		}
	}

	return strings.Join(tags, " ")
}

func (s *FacebookService) BuildMessage(product *models.Product, status string) string {
	pName := ""
	if product.Name != nil {
		pName = *product.Name
	}
	pType := ""
	if product.Type != nil {
		pType = *product.Type
	}

	statusLine := fmt.Sprintf("✅ CÒN HÀNG — %s", pName)
	if status == "SOLD" {
		statusLine = fmt.Sprintf("⛔ ĐÃ BÁN — %s", pName)
	}

	var priceVal int64
	if product.Price != nil {
		priceVal = *product.Price
	}
	priceFormatted := FormatVNDPrice(priceVal)

	var detectedText string
	if s.DB != nil && product.ID > 0 {
		var logRecord models.ProductUploadLog
		if err := s.DB.Where("product_id = ? AND detected_text IS NOT NULL AND detected_text != ''", product.ID).
			Order("id DESC").First(&logRecord).Error; err == nil && logRecord.DetectedText != nil {
			detectedText = *logRecord.DetectedText
		}
	}

	lines := []string{
		statusLine,
		fmt.Sprintf("👔 %s | %s", pType, priceFormatted),
	}

	if detectedText != "" {
		lines = append(lines, "", strings.TrimSpace(detectedText))
	}

	lines = append(lines, "", s.BuildHashtags(product, status))

	return strings.Join(lines, "\n")
}

func (s *FacebookService) BuildAlbumPhotoCaption(product *models.Product, catalogPostID string) string {
	pName := ""
	if product.Name != nil {
		pName = *product.Name
	}
	pType := ""
	if product.Type != nil {
		pType = *product.Type
	}

	var priceVal int64
	if product.Price != nil {
		priceVal = *product.Price
	}
	priceFormatted := FormatVNDPrice(priceVal)

	lines := []string{
		fmt.Sprintf("✅ CÒN HÀNG — %s", pName),
		fmt.Sprintf("👔 %s | Giá: %s", pType, priceFormatted),
	}

	if len(product.Sizes) > 0 {
		var sizeParts []string
		for _, sz := range product.Sizes {
			if sz.Name != nil && sz.Value != nil && *sz.Value > 0 {
				sizeParts = append(sizeParts, fmt.Sprintf("%s: %.0f", *sz.Name, *sz.Value))
			}
		}
		if len(sizeParts) > 0 {
			lines = append(lines, fmt.Sprintf("📐 Số đo: %s", strings.Join(sizeParts, ", ")))
		}
	}

	if catalogPostID != "" {
		cleanPostID := strings.TrimPrefix(catalogPostID, s.PageID+"_")
		catalogURL := fmt.Sprintf("https://www.facebook.com/%s/posts/%s", s.PageID, cleanPostID)
		lines = append(lines, "")
		lines = append(lines, "👉 Xem ảnh chi tiết tại:")
		lines = append(lines, catalogURL)
	}

	lines = append(lines, "")
	lines = append(lines, s.BuildHashtags(product, "AVAILABLE"))

	return strings.Join(lines, "\n")
}

func parseAppUsageMaxPct(header string) int {
	if header == "" {
		return 0
	}
	var usage struct {
		CallCount    int `json:"call_count"`
		TotalCPUTime int `json:"total_cputime"`
		TotalTime    int `json:"total_time"`
	}
	if err := json.Unmarshal([]byte(header), &usage); err != nil {
		return 0
	}
	maxPct := usage.TotalTime
	if usage.CallCount > maxPct {
		maxPct = usage.CallCount
	}
	if usage.TotalCPUTime > maxPct {
		maxPct = usage.TotalCPUTime
	}
	return maxPct
}

func (s *FacebookService) tokenForApp(appIdx int) string {
	if appIdx == 2 && s.MainPageToken2 != "" {
		return s.MainPageToken2
	}
	return s.mainToken()
}

func (s *FacebookService) markAppRateLimited(appIdx int) {
	if appIdx != 1 && appIdx != 2 {
		return
	}
	s.appMu.Lock()
	defer s.appMu.Unlock()
	s.appCooldownUntil[appIdx] = time.Now().Add(resyncPauseDuration)
}

func (s *FacebookService) pickAvailableUploadApp(excludeIdx int) int {
	s.appMu.Lock()
	defer s.appMu.Unlock()
	now := time.Now()
	if excludeIdx != 1 && s.mainToken() != "" && !s.appCooldownUntil[1].After(now) {
		return 1
	}
	if excludeIdx != 2 && s.MainPageToken2 != "" && !s.appCooldownUntil[2].After(now) {
		return 2
	}
	return 0
}

func (s *FacebookService) NextAvailableAppResumeTime() time.Time {
	s.appMu.Lock()
	defer s.appMu.Unlock()
	earliest := s.appCooldownUntil[1]
	if s.MainPageToken2 != "" && (earliest.IsZero() || s.appCooldownUntil[2].Before(earliest)) {
		earliest = s.appCooldownUntil[2]
	}
	if earliest.IsZero() || !earliest.After(time.Now()) {
		return time.Now().Add(resyncPauseDuration)
	}
	return earliest
}

// UploadAlbumPhotoWithFailover uploads an album photo using App 1, and immediately fails over to App 2 if App 1 hits #4 rate limit (or vice-versa).
func (s *FacebookService) UploadAlbumPhotoWithFailover(albumID string, imagePath string, caption string) (string, int, int, error) {
	appIdx := s.pickAvailableUploadApp(0)
	if appIdx == 0 {
		return "", 0, 100, fmt.Errorf("facebook album photo upload error (403): (#4) Application request limit reached (all apps in cooldown)")
	}

	photoID, usagePct, err := s.uploadAlbumPhotoWithToken(albumID, imagePath, caption, s.tokenForApp(appIdx))
	if err != nil && isFBRateLimitError(err) {
		s.markAppRateLimited(appIdx)
		otherIdx := s.pickAvailableUploadApp(appIdx)
		if otherIdx != 0 {
			log.Printf("🔄 App %d chạm rate limit (#4), tự động chuyển sang App %d ngay lập tức!\n", appIdx, otherIdx)
			photoID2, usagePct2, err2 := s.uploadAlbumPhotoWithToken(albumID, imagePath, caption, s.tokenForApp(otherIdx))
			if err2 != nil && isFBRateLimitError(err2) {
				s.markAppRateLimited(otherIdx)
			}
			return photoID2, otherIdx, usagePct2, err2
		}
	}
	return photoID, appIdx, usagePct, err
}

func (s *FacebookService) UploadProductPhotoToAlbum(albumID string, imagePath string, caption string) (string, error) {
	photoID, _, err := s.UploadProductPhotoToAlbumWithUsage(albumID, imagePath, caption)
	return photoID, err
}

func (s *FacebookService) UploadProductPhotoToAlbumWithUsage(albumID string, imagePath string, caption string) (string, int, error) {
	return s.uploadAlbumPhotoWithToken(albumID, imagePath, caption, s.mainToken())
}

func (s *FacebookService) uploadAlbumPhotoWithToken(albumID string, imagePath string, caption string, token string) (string, int, error) {
	if albumID == "" {
		return "", 0, fmt.Errorf("albumID is empty")
	}
	if token == "" {
		return "", 0, fmt.Errorf("facebook main page token is empty")
	}

	endpoint := fmt.Sprintf("%s%s/photos", s.apiURL(), url.PathEscape(albumID))

	// If remote URL
	if strings.HasPrefix(imagePath, "http://") || strings.HasPrefix(imagePath, "https://") {
		data := url.Values{}
		data.Set("url", imagePath)
		data.Set("caption", caption)
		data.Set("no_story", "true")
		data.Set("access_token", token)

		resp, err := s.Client.PostForm(endpoint, data)
		if err == nil {
			usagePct := parseAppUsageMaxPct(resp.Header.Get("x-app-usage"))
			respBytes, rErr := io.ReadAll(resp.Body)
			resp.Body.Close()
			if rErr == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
				var res struct {
					ID string `json:"id"`
				}
				if json.Unmarshal(respBytes, &res) == nil && res.ID != "" {
					return res.ID, usagePct, nil
				}
			}
			// Do not fall back to multipart upload if Facebook rate-limited us
			if resp.StatusCode == 403 || strings.Contains(string(respBytes), "(#4)") || strings.Contains(string(respBytes), "request limit reached") {
				return "", usagePct, fmt.Errorf("facebook album photo upload error (%d): %s", resp.StatusCode, string(respBytes))
			}
		}

		// Fallback: download bytes and upload via multipart source
		dlResp, dlErr := s.Client.Get(imagePath)
		if dlErr != nil {
			return "", 0, dlErr
		}
		defer dlResp.Body.Close()
		if dlResp.StatusCode < 200 || dlResp.StatusCode >= 300 {
			return "", 0, fmt.Errorf("remote image not found or inaccessible (HTTP %d): %s", dlResp.StatusCode, imagePath)
		}
		return s.uploadAlbumPhotoFromReaderWithToken(albumID, dlResp.Body, filepath.Base(imagePath), caption, token)
	}

	// Local file
	localPath := imagePath
	if !filepath.IsAbs(localPath) {
		if _, err := os.Stat(localPath); os.IsNotExist(err) {
			clean := strings.TrimPrefix(localPath, "/")
			if _, err2 := os.Stat(clean); err2 == nil {
				localPath = clean
			} else if _, err3 := os.Stat("storage/app/public/" + strings.TrimPrefix(clean, "storage/")); err3 == nil {
				localPath = "storage/app/public/" + strings.TrimPrefix(clean, "storage/")
			} else if _, err4 := os.Stat("./storage/" + strings.TrimPrefix(clean, "storage/")); err4 == nil {
				localPath = "./storage/" + strings.TrimPrefix(clean, "storage/")
			}
		}
	}

	if f, err := os.Open(localPath); err == nil {
		defer f.Close()
		return s.uploadAlbumPhotoFromReaderWithToken(albumID, f, filepath.Base(localPath), caption, token)
	}

	// If not found locally, try CDN with MediaBase
	if s.MediaBase != "" {
		remoteURL := s.MediaBase
		if !strings.HasSuffix(remoteURL, "/") {
			remoteURL += "/"
		}
		remoteURL += strings.TrimPrefix(imagePath, "storage/")
		return s.uploadAlbumPhotoWithToken(albumID, remoteURL, caption, token)
	}

	return "", 0, fmt.Errorf("image not accessible at path %s", imagePath)
}

func (s *FacebookService) uploadAlbumPhotoFromReaderWithToken(albumID string, r io.Reader, filename string, caption string, token string) (string, int, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("no_story", "true")
	_ = writer.WriteField("caption", caption)
	_ = writer.WriteField("access_token", token)
	part, err := writer.CreateFormFile("source", filename)
	if err != nil {
		return "", 0, err
	}
	if _, err := io.Copy(part, r); err != nil {
		return "", 0, err
	}
	_ = writer.Close()

	endpoint := fmt.Sprintf("%s%s/photos", s.apiURL(), url.PathEscape(albumID))
	req, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		return "", 0, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := s.Client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	usagePct := parseAppUsageMaxPct(resp.Header.Get("x-app-usage"))
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", usagePct, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", usagePct, fmt.Errorf("facebook album photo upload error (%d): %s", resp.StatusCode, string(respBytes))
	}

	var res struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(respBytes, &res); err != nil {
		return "", usagePct, err
	}
	return res.ID, usagePct, nil
}

func (s *FacebookService) RemoveProductPhotoFromAlbum(photoID string) error {
	_, err := s.RemoveAlbumPhotoWithFallback(photoID, 1)
	return err
}

func (s *FacebookService) RemoveProductPhotoFromAlbumWithUsage(photoID string) (int, error) {
	return s.RemoveAlbumPhotoWithFallback(photoID, 1)
}

func (s *FacebookService) RemoveAlbumPhotoWithFallback(photoID string, appIdx int) (int, error) {
	if photoID == "" {
		return 0, nil
	}
	primaryIdx := 1
	if appIdx == 2 && s.MainPageToken2 != "" {
		primaryIdx = 2
	}
	usagePct, err := s.removeAlbumPhotoWithToken(photoID, s.tokenForApp(primaryIdx))
	if err != nil {
		if isFBRateLimitError(err) {
			s.markAppRateLimited(primaryIdx)
		}
		// If Meta says "App can only delete photos created by the same app" (#200), try the other app token
		if strings.Contains(err.Error(), "(#200)") || strings.Contains(err.Error(), "same app") {
			otherIdx := 2
			if primaryIdx == 2 {
				otherIdx = 1
			}
			if otherTok := s.tokenForApp(otherIdx); otherTok != "" && otherTok != s.tokenForApp(primaryIdx) {
				usage2, err2 := s.removeAlbumPhotoWithToken(photoID, otherTok)
				if err2 != nil && isFBRateLimitError(err2) {
					s.markAppRateLimited(otherIdx)
				}
				return usage2, err2
			}
		}
	}
	return usagePct, err
}

func (s *FacebookService) removeAlbumPhotoWithToken(photoID string, token string) (int, error) {
	if photoID == "" {
		return 0, nil
	}
	endpoint := fmt.Sprintf("%s%s?access_token=%s", s.apiURL(), url.PathEscape(photoID), url.QueryEscape(token))
	req, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		return 0, err
	}
	resp, err := s.Client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	usagePct := parseAppUsageMaxPct(resp.Header.Get("x-app-usage"))
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return usagePct, err
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return usagePct, nil
	}

	// Handle case where object is already deleted or not found on Facebook
	var errResp struct {
		Error struct {
			Message      string `json:"message"`
			Type         string `json:"type"`
			Code         int    `json:"code"`
			ErrorSubcode int    `json:"error_subcode"`
		} `json:"error"`
	}
	if json.Unmarshal(respBytes, &errResp) == nil {
		if errResp.Error.Code == 100 || errResp.Error.ErrorSubcode == 33 || strings.Contains(strings.ToLower(errResp.Error.Message), "does not exist") {
			log.Printf("FB album photo %s already deleted or does not exist on Facebook\n", photoID)
			return usagePct, nil
		}
	}

	return usagePct, fmt.Errorf("facebook delete photo error (%d): %s", resp.StatusCode, string(respBytes))
}

func (s *FacebookService) GetPageAlbums(pageID string) ([]FacebookAlbum, error) {
	if pageID == "" {
		pageID = s.MainPageID
		if pageID == "" {
			pageID = s.PageID
		}
	}
	token := s.mainToken()
	endpoint := fmt.Sprintf("%s%s/albums?fields=id,name,count,created_time&limit=100&access_token=%s", s.apiURL(), url.PathEscape(pageID), url.QueryEscape(token))

	resp, err := s.Client.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res struct {
		Data []FacebookAlbum `json:"data"`
	}
	if err := json.Unmarshal(respBytes, &res); err != nil {
		return nil, err
	}
	return res.Data, nil
}

func (s *FacebookService) CreateAlbum(pageID string, name string, description string) (string, error) {
	if pageID == "" {
		pageID = s.MainPageID
	}
	token := s.mainToken()
	endpoint := fmt.Sprintf("%s%s/albums", s.apiURL(), url.PathEscape(pageID))

	data := url.Values{}
	data.Set("name", name)
	if description != "" {
		data.Set("message", description)
	}
	data.Set("access_token", token)

	resp, err := s.Client.PostForm(endpoint, data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("facebook create album error (%d): %s", resp.StatusCode, string(respBytes))
	}

	var res struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(respBytes, &res); err != nil {
		return "", err
	}
	return res.ID, nil
}

type FeedPhotoItem struct {
	Path    string `json:"path"`
	Caption string `json:"caption"`
}

// UploadMainPagePhoto uploads an unpublished photo to the Main Page for feed post attachment with optional caption
func (s *FacebookService) UploadMainPagePhoto(imagePath string, caption string) (string, error) {
	pageID := s.MainPageID
	if pageID == "" {
		pageID = s.PageID
	}
	token := s.mainToken()
	if token == "" {
		return "", fmt.Errorf("facebook main page token is empty")
	}

	endpoint := fmt.Sprintf("%s%s/photos", s.apiURL(), url.PathEscape(pageID))

	// If remote URL
	if strings.HasPrefix(imagePath, "http://") || strings.HasPrefix(imagePath, "https://") {
		data := url.Values{}
		data.Set("url", imagePath)
		data.Set("published", "false")
		if caption != "" {
			data.Set("caption", caption)
		}
		data.Set("access_token", token)

		resp, err := s.Client.PostForm(endpoint, data)
		if err == nil {
			defer resp.Body.Close()
			respBytes, rErr := io.ReadAll(resp.Body)
			if rErr == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
				var res struct {
					ID string `json:"id"`
				}
				if json.Unmarshal(respBytes, &res) == nil && res.ID != "" {
					return res.ID, nil
				}
			}
		}

		// Fallback: download bytes and upload via multipart source
		dlResp, dlErr := s.Client.Get(imagePath)
		if dlErr != nil {
			return "", dlErr
		}
		defer dlResp.Body.Close()
		if dlResp.StatusCode < 200 || dlResp.StatusCode >= 300 {
			return "", fmt.Errorf("remote image inaccessible (HTTP %d): %s", dlResp.StatusCode, imagePath)
		}
		return s.uploadMainPhotoFromReader(pageID, token, dlResp.Body, filepath.Base(imagePath), caption)
	}

	// Local file
	localPath := imagePath
	if !filepath.IsAbs(localPath) {
		if _, err := os.Stat(localPath); os.IsNotExist(err) {
			clean := strings.TrimPrefix(localPath, "/")
			if _, err2 := os.Stat(clean); err2 == nil {
				localPath = clean
			} else if _, err3 := os.Stat("storage/app/public/" + strings.TrimPrefix(clean, "storage/")); err3 == nil {
				localPath = "storage/app/public/" + strings.TrimPrefix(clean, "storage/")
			} else if _, err4 := os.Stat("./storage/" + strings.TrimPrefix(clean, "storage/")); err4 == nil {
				localPath = "./storage/" + strings.TrimPrefix(clean, "storage/")
			}
		}
	}

	if f, err := os.Open(localPath); err == nil {
		defer f.Close()
		return s.uploadMainPhotoFromReader(pageID, token, f, filepath.Base(localPath), caption)
	}

	// If not found locally, try CDN
	if s.MediaBase != "" {
		remoteURL := s.MediaBase
		if !strings.HasSuffix(remoteURL, "/") {
			remoteURL += "/"
		}
		remoteURL += strings.TrimPrefix(imagePath, "storage/")
		return s.UploadMainPagePhoto(remoteURL, caption)
	}

	return "", fmt.Errorf("image not accessible at path %s", imagePath)
}

func (s *FacebookService) uploadMainPhotoFromReader(pageID, token string, r io.Reader, filename string, caption string) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("published", "false")
	if caption != "" {
		_ = writer.WriteField("caption", caption)
	}
	_ = writer.WriteField("access_token", token)
	part, err := writer.CreateFormFile("source", filename)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(part, r); err != nil {
		return "", err
	}
	_ = writer.Close()

	endpoint := fmt.Sprintf("%s%s/photos", s.apiURL(), url.PathEscape(pageID))
	req, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := s.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("facebook main photo upload error (%d): %s", resp.StatusCode, string(respBytes))
	}

	var res struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(respBytes, &res); err != nil {
		return "", err
	}
	return res.ID, nil
}

// CreateMainPageFeedPost creates a multi-photo post on the main page, appending optional endCard as the final photo
func (s *FacebookService) CreateMainPageFeedPost(message string, photoPaths []string, endCard *FeedPhotoItem) (string, error) {
	pageID := s.MainPageID
	if pageID == "" {
		pageID = s.PageID
	}
	token := s.mainToken()
	if token == "" {
		return "", fmt.Errorf("facebook main page token is empty")
	}

	type attachedMedia struct {
		MediaFBID string `json:"media_fbid"`
	}
	var mediaList []attachedMedia

	for i, path := range photoPaths {
		mediaID, err := s.UploadMainPagePhoto(path, "")
		if err != nil {
			log.Printf("Warning: failed to upload photo [%d] (%s): %v\n", i+1, path, err)
			continue
		}
		mediaList = append(mediaList, attachedMedia{MediaFBID: mediaID})
		time.Sleep(150 * time.Millisecond) // Throttling safe delay
	}

	// Upload end-card if provided
	if endCard != nil && endCard.Path != "" {
		endMediaID, err := s.UploadMainPagePhoto(endCard.Path, endCard.Caption)
		if err == nil && endMediaID != "" {
			mediaList = append(mediaList, attachedMedia{MediaFBID: endMediaID})
		} else {
			log.Printf("Warning: failed to upload end-card photo (%s): %v\n", endCard.Path, err)
		}
	}

	if len(mediaList) == 0 {
		return "", fmt.Errorf("no photos uploaded successfully for feed post")
	}

	mediaJSON, err := json.Marshal(mediaList)
	if err != nil {
		return "", err
	}

	endpoint := fmt.Sprintf("%s%s/feed", s.apiURL(), url.PathEscape(pageID))
	data := url.Values{}
	data.Set("message", message)
	data.Set("attached_media", string(mediaJSON))
	data.Set("access_token", token)

	resp, err := s.Client.PostForm(endpoint, data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("facebook feed post error (%d): %s", resp.StatusCode, string(respBytes))
	}

	var res struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(respBytes, &res); err != nil {
		return "", err
	}

	return res.ID, nil
}
