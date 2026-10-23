package services

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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
