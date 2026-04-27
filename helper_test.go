package ff4go

import (
	"math"
	"testing"
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
