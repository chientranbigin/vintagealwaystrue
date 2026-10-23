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
	"time"

	"sale-app/models"

	"gorm.io/gorm"
)

type FacebookSyncAction string

const (
	FBActionCreate FacebookSyncAction = "create"
	FBActionSold   FacebookSyncAction = "sold"
	FBActionRevert FacebookSyncAction = "revert"
)

type FacebookSyncJob struct {
	ProductID uint
	Action    FacebookSyncAction
	Retries   int
}

type FacebookService struct {
	PageID    string
	PageToken string
	BaseURL   string
	Client    *http.Client
	DB        *gorm.DB
	MediaBase string
	Queue     chan FacebookSyncJob
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

func NewFacebookService(db *gorm.DB, pageID string, pageToken string, mediaBase string) *FacebookService {
	if pageID == "" {
		pageID = os.Getenv("FB_PAGE_ID")
	}
	if pageToken == "" {
		pageToken = os.Getenv("FB_PAGE_TOKEN")
	}
	if mediaBase == "" {
		mediaBase = os.Getenv("MEDIA_BASE_URL")
		if mediaBase == "" {
			mediaBase = "https://res.cloudinary.com/leq5n17r/image/upload/"
		}
	}

	fb := &FacebookService{
		PageID:    pageID,
		PageToken: pageToken,
		BaseURL:   "https://graph.facebook.com/v19.0/",
		Client:    &http.Client{Timeout: 30 * time.Second},
		DB:        db,
		MediaBase: mediaBase,
		Queue:     make(chan FacebookSyncJob, 200),
	}

	fb.StartWorker()
	return fb
}

func (s *FacebookService) StartWorker() {
	go func() {
		for job := range s.Queue {
			s.handleJob(job)
		}
	}()
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
	if s.PageToken == "" {
		log.Println("SyncFacebookPost: no FB token configured, skipping.")
		return
	}

	err := s.ExecuteSync(job.ProductID, job.Action)
	if err != nil {
		log.Printf("SyncFacebookPost [%s] error for product ID %d: %v\n", job.Action, job.ProductID, err)
		if job.Retries < 3 {
			go func(j FacebookSyncJob) {
				time.Sleep(5 * time.Second)
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
	case FBActionCreate:
		if product.FBPostID != nil && *product.FBPostID != "" {
			_ = s.DeletePost(*product.FBPostID)
			log.Printf("FB post deleted before re-create: %s for %s\n", *product.FBPostID, pName)
		}

		postID, err := s.CreatePost(&product)
		if err != nil {
			return err
		}
		if postID != "" {
			s.DB.Model(&models.Product{}).Where("id = ?", product.ID).Update("fb_post_id", postID)
			log.Printf("FB post created: %s for %s\n", postID, pName)
		}

	case FBActionSold:
		if product.FBPostID == nil || *product.FBPostID == "" {
			return nil
		}
		if err := s.MarkSold(*product.FBPostID, &product); err != nil {
			return err
		}
		log.Printf("FB post marked SOLD: %s\n", *product.FBPostID)

	case FBActionRevert:
		if product.FBPostID == nil || *product.FBPostID == "" {
			return nil
		}
		if err := s.RevertAvailable(*product.FBPostID, &product); err != nil {
			return err
		}
		log.Printf("FB post reverted AVAILABLE: %s\n", *product.FBPostID)
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
