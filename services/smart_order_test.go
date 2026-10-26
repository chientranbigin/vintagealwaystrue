package services

import (
	"context"
	"testing"
)

func TestSmartOrderService_FallbackParse(t *testing.T) {
	svc := NewSmartOrderService()

	sampleText := `Anh Tuấn 0912345678
Số 123 đường Cầu Giấy, Hà Nội
Lấy quần ts1169, áo st567
note: giao giờ hành chính, gọi trước khi giao`

	res, err := svc.ParseDescription(context.Background(), sampleText)
	if err != nil {
		t.Fatalf("ParseDescription failed: %v", err)
	}

	if res.Phone != "0912345678" {
		t.Errorf("Expected phone 0912345678, got %s", res.Phone)
	}

	if res.Name == "" || res.Name != "Anh Tuấn" {
		t.Errorf("Expected name 'Anh Tuấn', got '%s'", res.Name)
	}

	if res.Note != "giao giờ hành chính, gọi trước khi giao" {
		t.Errorf("Expected note 'giao giờ hành chính, gọi trước khi giao', got '%s'", res.Note)
	}

	expectedCodes := map[string]bool{"ts1169": true, "st567": true}
	for _, c := range res.Products {
		if !expectedCodes[c] {
			t.Errorf("Unexpected product code %s", c)
		}
	}
	if len(res.Products) != 2 {
		t.Errorf("Expected 2 products, got %d: %v", len(res.Products), res.Products)
	}

	if res.Address == "" {
		t.Errorf("Expected non-empty address")
	}
}
