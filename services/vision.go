package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type VisionService struct {
	APIKey string
	Client *http.Client
}

func NewVisionService() *VisionService {
	apiKey := os.Getenv("GOOGLE_VISION_API_KEY")
	if apiKey == "" {
		apiKey = "AIzaSyBCQwssnHoUNbueD6DGggilZou-Sd0NjU0"
	}

	return &VisionService{
		APIKey: apiKey,
		Client: &http.Client{Timeout: 60 * time.Second},
	}
}

// ExtractTextFromImage extracts all detected text from an image path or URL
func (v *VisionService) ExtractTextFromImage(imagePathOrURL string) (string, error) {
	var imgBytes []byte
	var err error

	if strings.HasPrefix(imagePathOrURL, "http://") || strings.HasPrefix(imagePathOrURL, "https://") {
		// Download image from URL
		resp, err := v.Client.Get(imagePathOrURL)
		if err != nil {
			return "", fmt.Errorf("failed to fetch image from URL: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("failed to download image, status code: %d", resp.StatusCode)
		}

		imgBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read image body: %w", err)
		}
	} else {
		// Read from local storage
		localPath := imagePathOrURL
		if !filepath.IsAbs(localPath) {
			// Try checking directly or in ./storage/sale
			if _, err := os.Stat(localPath); os.IsNotExist(err) {
				clean := strings.TrimPrefix(localPath, "storage/sale/")
				clean = strings.TrimPrefix(clean, "storage/")
				clean = strings.TrimPrefix(clean, "/")
				localPath = filepath.Join("./storage/sale", clean)
			}
		}

		imgBytes, err = os.ReadFile(localPath)
		if err != nil {
			return "", fmt.Errorf("failed to read local image file (%s): %w", localPath, err)
		}
	}

	if len(imgBytes) == 0 {
		return "", fmt.Errorf("empty image data")
	}

	b64Data := base64.StdEncoding.EncodeToString(imgBytes)

	reqPayload := map[string]interface{}{
		"requests": []map[string]interface{}{
			{
				"image": map[string]string{
					"content": b64Data,
				},
				"features": []map[string]string{
					{"type": "TEXT_DETECTION"},
				},
			},
		},
	}

	reqBody, err := json.Marshal(reqPayload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal vision request: %w", err)
	}

	endpoint := fmt.Sprintf("https://vision.googleapis.com/v1/images:annotate?key=%s", v.APIKey)
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := v.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("vision api request failed: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read vision response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("vision api error (%d): %s", resp.StatusCode, string(respBytes))
	}

	var visionResp struct {
		Responses []struct {
			FullTextAnnotation struct {
				Text string `json:"text"`
			} `json:"fullTextAnnotation"`
			Error *struct {
				Message string `json:"message"`
			} `json:"error"`
		} `json:"responses"`
	}

	if err := json.Unmarshal(respBytes, &visionResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal vision response: %w", err)
	}

	if len(visionResp.Responses) == 0 {
		return "", nil
	}

	firstResp := visionResp.Responses[0]
	if firstResp.Error != nil && firstResp.Error.Message != "" {
		return "", fmt.Errorf("vision error: %s", firstResp.Error.Message)
	}

	return firstResp.FullTextAnnotation.Text, nil
}
