package handlers

import (
	"sale-app/models"
	"testing"
)

func TestFormatProductSizesTS4717(t *testing.T) {
	eoQuan := 40.0
	day := 36.0
	dui := 33.0
	daiQuan := 106.0
	duLai := 12.0
	ong := 22.0

	sEoQuan := "EO QUẦN"
	sDay := "ĐÁY"
	sDui := "ĐÙI"
	sDaiQuan := "DÀI QUẦN"
	sDuLai := "DƯ LAI"
	sOng := "ỐNG"

	product := models.Product{
		Sizes: []models.Size{
			{Name: &sEoQuan, Value: &eoQuan},
			{Name: &sDay, Value: &day},
			{Name: &sDui, Value: &dui},
			{Name: &sDaiQuan, Value: &daiQuan},
			{Name: &sDuLai, Value: &duLai},
			{Name: &sOng, Value: &ong},
		},
	}

	got := formatProductSizes(&product)
	expected := "E80 - ĐA36 - Đ33 - D106 (+12) - Ô22"

	if got != expected {
		t.Errorf("formatProductSizes() = %q, want %q", got, expected)
	}
}

func TestMappingSizesParityWithOldPHP(t *testing.T) {
	expectedMapping := map[string]string{
		"VAI":        "V",
		"NGỰC":       "N",
		"EO":         "E",
		"DÀI ÁO":     "D",
		"DÀI ÁO SAU": "DS",
		"DÀI TAY":    "DT",
		"EO QUẦN":    "E",
		"ĐÁY":        "ĐA",
		"ĐÙI":        "Đ",
		"DÀI QUẦN":   "D",
		"ỐNG":        "Ô",
		"DƯ LAI":     "L",
	}

	for k, v := range expectedMapping {
		if MappingSizes[k] != v {
			t.Errorf("MappingSizes[%q] = %q, want %q", k, MappingSizes[k], v)
		}
	}
}
