package ff4go

import (
	"math"
	"testing"
	"time"
)

func TestHashStringToFloat(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"hello", 0.83},
		{"world", 0.93},
		{"ff4go", 0.87},
	}

	for _, test := range tests {
		result := math.Round(hashStringToFloat(test.input, "testID")*100) / 100
		expected := math.Round(test.expected*100) / 100

		if result != expected {
			t.Errorf("hashStringToFloat(%q) = %f; expected %f", test.input, result, expected)
		}
	}
}

func TestIsExpired(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"2023-01-01T00:00:00Z", true},
		{time.Now().Add(24 * time.Hour).Format(time.RFC3339), false},
		{"invalid-date", false},
		{"", false},
	}

	for _, test := range tests {
		result := isExpired(test.input)

		if result != test.expected {
			t.Errorf("isExpired(%q) = %v; expected %v", test.input, result, test.expected)
		}
	}
}
