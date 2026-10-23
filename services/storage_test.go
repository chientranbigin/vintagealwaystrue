package services

import (
	"testing"
)

func TestExtractPublicID(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "https://res.cloudinary.com/leq5n17r/image/upload/v1792547220/sale/KRxM8jmLHY0g0X1q92OjWD2X2H1ee4ACgq1g8MFc.jpg",
			expected: "sale/KRxM8jmLHY0g0X1q92OjWD2X2H1ee4ACgq1g8MFc",
		},
		{
			input:    "https://res.cloudinary.com/leq5n17r/image/upload/sale/BuXwUd02ZDFO7InXEBDcGWA5a76I9N3CUVo7B4Iq.jpg",
			expected: "sale/BuXwUd02ZDFO7InXEBDcGWA5a76I9N3CUVo7B4Iq",
		},
		{
			input:    "storage/sale/X13R1eAp7Fh0TARDsf2gpadC4ybog85DQSbo8DEn.png",
			expected: "sale/X13R1eAp7Fh0TARDsf2gpadC4ybog85DQSbo8DEn",
		},
		{
			input:    "sale/p5USe7OyDBnLyBJtwyBk7D74RMdWNzeRFj2BEgMi_scaled.jpg",
			expected: "sale/p5USe7OyDBnLyBJtwyBk7D74RMdWNzeRFj2BEgMi_scaled",
		},
		{
			input:    "Kkedlt9W3MEpK50j4HrvsAz7qAPvHpwTESsMX0MM.jpg",
			expected: "sale/Kkedlt9W3MEpK50j4HrvsAz7qAPvHpwTESsMX0MM",
		},
		{
			input:    "",
			expected: "",
		},
	}

	for _, tc := range tests {
		actual := ExtractPublicID(tc.input)
		if actual != tc.expected {
			t.Errorf("ExtractPublicID(%q) = %q, expected %q", tc.input, actual, tc.expected)
		}
	}
}
