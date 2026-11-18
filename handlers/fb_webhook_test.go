package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"

	"sale-app/config"

	"github.com/gin-gonic/gin"
)

func TestFBWebhook_VerifyAndCompliancePages(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewCopilotHandler(nil, &config.Config{}, nil)

	r := gin.New()
	r.GET("/webhook/facebook", h.VerifyWebhook)
	r.GET("/privacy-policy", h.PrivacyPolicy)
	r.POST("/data-deletion", h.DataDeletion)

	// 1. Valid verify token
	req := httptest.NewRequest("GET", "/webhook/facebook?hub.mode=subscribe&hub.verify_token=vintage-copilot-verify-2026&hub.challenge=123456", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK || w.Body.String() != "123456" {
		t.Fatalf("expected 200 '123456', got %d %q", w.Code, w.Body.String())
	}

	// 2. Invalid verify token
	reqBad := httptest.NewRequest("GET", "/webhook/facebook?hub.mode=subscribe&hub.verify_token=wrong&hub.challenge=123", nil)
	wBad := httptest.NewRecorder()
	r.ServeHTTP(wBad, reqBad)
	if wBad.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", wBad.Code)
	}

	// 3. Privacy Policy
	reqPP := httptest.NewRequest("GET", "/privacy-policy", nil)
	wPP := httptest.NewRecorder()
	r.ServeHTTP(wPP, reqPP)
	if wPP.Code != http.StatusOK {
		t.Fatalf("expected 200 for privacy-policy, got %d", wPP.Code)
	}

	// 4. Data Deletion callback
	reqDD := httptest.NewRequest("POST", "/data-deletion", nil)
	wDD := httptest.NewRecorder()
	r.ServeHTTP(wDD, reqDD)
	if wDD.Code != http.StatusOK {
		t.Fatalf("expected 200 for data-deletion, got %d", wDD.Code)
	}
}

func TestFBWebhook_VerifyMetaSignature(t *testing.T) {
	secret := "test-app-secret"
	body := []byte(`{"object":"page","entry":[]}`)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	sig := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	if !verifyMetaSignature(body, sig, secret) {
		t.Fatalf("expected valid signature to pass")
	}
	if verifyMetaSignature(body, "sha256=deadbeef", secret) {
		t.Fatalf("expected invalid signature to fail")
	}
}
