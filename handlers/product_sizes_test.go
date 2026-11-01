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

func TestNormalizeProductSizes(t *testing.T) {
	val44 := 44.0
	val35 := 35.0
	val107 := 107.0

	sVai := "VAI"
	sNguc := "NGỰC"
	sDaiAo := "DÀI ÁO"

	sizes := []models.Size{
		{Name: &sVai, Value: &val44},
		{Name: &sNguc, Value: &val35},
		{Name: &sDaiAo, Value: &val107},
	}

	normalized := normalizeProductSizes("TROUSERS", sizes)
	if *normalized[0].Name != "EO QUẦN" {
		t.Errorf("expected EO QUẦN, got %s", *normalized[0].Name)
	}
	if *normalized[1].Name != "ĐÁY" {
		t.Errorf("expected ĐÁY, got %s", *normalized[1].Name)
	}
	if *normalized[2].Name != "DÀI QUẦN" {
		t.Errorf("expected DÀI QUẦN, got %s", *normalized[2].Name)
	}

	// Should not touch top types
	jacketSizes := []models.Size{
		{Name: &sVai, Value: &val44},
	}
	jacketNorm := normalizeProductSizes("JACKET", jacketSizes)
	if *jacketNorm[0].Name != "VAI" {
		t.Errorf("expected VAI, got %s", *jacketNorm[0].Name)
	}
}
