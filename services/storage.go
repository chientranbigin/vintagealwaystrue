package services

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"
)

type StorageService struct {
	CloudName  string
	APIKey     string
	APISecret  string
	UploadDir  string
	BaseURL    string
}

func NewStorageService() *StorageService {
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	if cloudName == "" {
		cloudName = "leq5n17r"
	}
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	if apiKey == "" {
		apiKey = "859827999643589"
	}
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")
	if apiSecret == "" {
		apiSecret = "2fLjdB1nZnBDenXAKKzh4tXAdug"
	}

	uploadDir := "./storage/sale"
	_ = os.MkdirAll(uploadDir, 0755)

	baseURL := os.Getenv("MEDIA_BASE_URL")
	if baseURL == "" {
		baseURL = "https://res.cloudinary.com/leq5n17r/image/upload/"
	}

	return &StorageService{
		CloudName: cloudName,
		APIKey:    apiKey,
		APISecret: apiSecret,
		UploadDir: uploadDir,
		BaseURL:   baseURL,
	}
}

// Upload processes a file header and returns the storage path / URL
func (s *StorageService) Upload(fileHeader *multipart.FileHeader) (string, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		return "", err
	}

	return s.UploadBytes(data, fileHeader.Filename)
}

// UploadBytes uploads raw image bytes to Cloudinary with local storage fallback
func (s *StorageService) UploadBytes(data []byte, filename string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		ext = ".jpg"
	}
	timestamp := time.Now().UnixNano()
	storedFilename := fmt.Sprintf("%d_%s%s", timestamp, randomString(8), ext)

	// If Cloudinary CloudName is provided, upload to Cloudinary
	if s.CloudName != "" {
		cURL, err := s.uploadToCloudinary(data, storedFilename)
		if err == nil && cURL != "" {
			return cURL, nil
		}
		// If failed, fall back to local/remote storage
		fmt.Printf("Cloudinary upload failed: %v, falling back to local storage\n", err)
	}

	// Save to local storage
	localPath := filepath.Join(s.UploadDir, storedFilename)
	if err := os.WriteFile(localPath, data, 0644); err != nil {
		return "", err
	}

	// Return path format matching Laravel: "storage/sale/filename.jpg"
	return "storage/sale/" + storedFilename, nil
}

// CreateScaledThumbnail resizes image to maxWidth 400px and saves as _scaled.jpg
func (s *StorageService) CreateScaledThumbnail(sourceData []byte, sourceFilenameOrPath string) (string, error) {
	defer debug.FreeOSMemory()
	if len(sourceData) == 0 {
		return "", fmt.Errorf("sourceData is empty")
	}

	src, _, err := image.Decode(bytes.NewReader(sourceData))
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	maxWidth := 400
	newWidth := width
	newHeight := height

	if width > maxWidth {
		newWidth = maxWidth
		newHeight = int(float64(height) * (float64(maxWidth) / float64(width)))
		if newHeight < 1 {
			newHeight = 1
		}
	}

	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	draw.BiLinear.Scale(dst, dst.Rect, src, bounds, draw.Over, nil)

	buf := &bytes.Buffer{}
	if err := jpeg.Encode(buf, dst, &jpeg.Options{Quality: 75}); err != nil {
		return "", fmt.Errorf("failed to encode scaled jpeg: %w", err)
	}

	scaledBytes := buf.Bytes()

	// Extract base name to name the scaled file
	base := filepath.Base(sourceFilenameOrPath)
	ext := filepath.Ext(base)
	baseWithoutExt := strings.TrimSuffix(base, ext)
	baseWithoutExt = strings.TrimSuffix(baseWithoutExt, "_scaled")
	scaledFilename := baseWithoutExt + "_scaled.jpg"

	// Upload to Cloudinary with public_id "sale/xxx_scaled"
	if s.CloudName != "" {
		_, err := s.uploadToCloudinary(scaledBytes, scaledFilename)
		if err == nil {
			return "sale/" + scaledFilename, nil
		}
	}

	// Local fallback
	localPath := filepath.Join(s.UploadDir, scaledFilename)
	_ = os.WriteFile(localPath, scaledBytes, 0644)

	return "sale/" + scaledFilename, nil
}

// CreateScaledThumbnailFromPath downloads or reads image and creates scaled thumbnail
func (s *StorageService) CreateScaledThumbnailFromPath(imagePath string) (string, error) {
	if imagePath == "" {
		return "", fmt.Errorf("imagePath is empty")
	}

	var data []byte
	var err error

	// If remote URL
	if strings.HasPrefix(imagePath, "http://") || strings.HasPrefix(imagePath, "https://") {
		client := &http.Client{Timeout: 30 * time.Second}
		resp, gErr := client.Get(imagePath)
		if gErr != nil {
			return "", gErr
		}
		defer resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return "", fmt.Errorf("remote image returned HTTP %d", resp.StatusCode)
		}
		data, err = io.ReadAll(resp.Body)
	} else {
		// Local file path
		cleanPath := imagePath
		if !filepath.IsAbs(cleanPath) {
			if _, sErr := os.Stat(cleanPath); os.IsNotExist(sErr) {
				clean := strings.TrimPrefix(cleanPath, "/")
				if _, sErr2 := os.Stat(clean); sErr2 == nil {
					cleanPath = clean
				} else if _, sErr3 := os.Stat("storage/sale/" + filepath.Base(clean)); sErr3 == nil {
					cleanPath = "storage/sale/" + filepath.Base(clean)
				}
			}
		}
		data, err = os.ReadFile(cleanPath)
		if err != nil && s.BaseURL != "" {
			// Try download from CDN
			cdnURL := s.BaseURL
			if !strings.HasSuffix(cdnURL, "/") {
				cdnURL += "/"
			}
			cdnURL += strings.TrimPrefix(imagePath, "storage/")
			return s.CreateScaledThumbnailFromPath(cdnURL)
		}
	}

	if err != nil || len(data) == 0 {
		return "", fmt.Errorf("failed to read image data from %s: %v", imagePath, err)
	}

	return s.CreateScaledThumbnail(data, imagePath)
}

func (s *StorageService) uploadToCloudinary(fileData []byte, filename string) (string, error) {
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	nameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))
	publicID := "sale/" + nameWithoutExt
	paramsToSign := fmt.Sprintf("public_id=%s&timestamp=%s%s", publicID, ts, s.APISecret)

	hasher := sha1.New()
	hasher.Write([]byte(paramsToSign))
	signature := hex.EncodeToString(hasher.Sum(nil))

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}
	if _, err := part.Write(fileData); err != nil {
		return "", err
	}

	_ = writer.WriteField("api_key", s.APIKey)
	_ = writer.WriteField("public_id", publicID)
	_ = writer.WriteField("timestamp", ts)
	_ = writer.WriteField("signature", signature)
	_ = writer.Close()

	endpoint := fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/image/upload", s.CloudName)
	req, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 45 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("cloudinary error: %s", string(respBytes))
	}

	var res struct {
		SecureURL string `json:"secure_url"`
		URL       string `json:"url"`
	}
	if err := json.Unmarshal(respBytes, &res); err != nil {
		return "", err
	}

	if res.SecureURL != "" {
		return res.SecureURL, nil
	}
	return res.URL, nil
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

// ExtractPublicID extracts the Cloudinary public_id (e.g. "sale/abc") from any image path or URL
func ExtractPublicID(pathOrURL string) string {
	pathOrURL = strings.TrimSpace(pathOrURL)
	if pathOrURL == "" {
		return ""
	}
	// If it's a URL, parse its Path
	if strings.HasPrefix(pathOrURL, "http://") || strings.HasPrefix(pathOrURL, "https://") {
		if u, err := url.Parse(pathOrURL); err == nil && u.Path != "" {
			pathOrURL = u.Path
		}
	}
	// Extract filename after /sale/ or sale/
	var filename string
	if idx := strings.Index(pathOrURL, "/sale/"); idx != -1 {
		filename = pathOrURL[idx+6:]
	} else if strings.HasPrefix(pathOrURL, "sale/") {
		filename = strings.TrimPrefix(pathOrURL, "sale/")
	} else {
		filename = filepath.Base(pathOrURL)
	}
	ext := filepath.Ext(filename)
	baseWithoutExt := strings.TrimSuffix(filename, ext)
	baseWithoutExt = strings.Trim(baseWithoutExt, "/")
	if baseWithoutExt == "" || baseWithoutExt == "." {
		return ""
	}
	return "sale/" + baseWithoutExt
}

// DeleteFromCloudinary deletes multiple assets from Cloudinary in batches of up to 100
func (s *StorageService) DeleteFromCloudinary(publicIDs []string) (int, error) {
	if s.CloudName == "" || s.APIKey == "" || s.APISecret == "" || len(publicIDs) == 0 {
		return 0, nil
	}

	totalDeleted := 0
	batchSize := 100
	client := &http.Client{Timeout: 60 * time.Second}

	for i := 0; i < len(publicIDs); i += batchSize {
		end := i + batchSize
		if end > len(publicIDs) {
			end = len(publicIDs)
		}
		chunk := publicIDs[i:end]

		endpoint := fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/resources/image/upload", s.CloudName)
		req, err := http.NewRequest("DELETE", endpoint, nil)
		if err != nil {
			return totalDeleted, err
		}
		req.SetBasicAuth(s.APIKey, s.APISecret)

		q := req.URL.Query()
		for _, pid := range chunk {
			q.Add("public_ids[]", pid)
		}
		req.URL.RawQuery = q.Encode()

		resp, err := client.Do(req)
		if err != nil {
			return totalDeleted, fmt.Errorf("cloudinary delete request failed: %w", err)
		}

		respBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return totalDeleted, fmt.Errorf("cloudinary delete error (status %d): %s", resp.StatusCode, string(respBytes))
		}

		var res struct {
			Deleted map[string]string `json:"deleted"`
		}
		if err := json.Unmarshal(respBytes, &res); err == nil {
			for _, status := range res.Deleted {
				if status == "deleted" {
					totalDeleted++
				}
			}
		}
	}

	return totalDeleted, nil
}
