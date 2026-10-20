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
		baseURL = "https://www.tadimedia.vn/storage/"
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

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext == "" {
		ext = ".jpg"
	}
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("%d_%s%s", timestamp, randomString(8), ext)

	// If Cloudinary CloudName is provided, upload to Cloudinary
	if s.CloudName != "" {
		cURL, err := s.uploadToCloudinary(data, filename)
		if err == nil && cURL != "" {
			return cURL, nil
		}
		// If failed, fall back to local/remote storage
		fmt.Printf("Cloudinary upload failed: %v, falling back to local storage\n", err)
	}

	// Save to local storage
	localPath := filepath.Join(s.UploadDir, filename)
	if err := os.WriteFile(localPath, data, 0644); err != nil {
		return "", err
	}

	// Return path format matching Laravel: "storage/sale/filename.jpg"
	return "storage/sale/" + filename, nil
}

func (s *StorageService) uploadToCloudinary(fileData []byte, filename string) (string, error) {
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	paramsToSign := fmt.Sprintf("timestamp=%s%s", ts, s.APISecret)

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
