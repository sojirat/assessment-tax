package tax_test

import (
	"testing"

	"github.com/sojirat/assessment-tax/tax"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected float64
	}{
		{
			name:     "Basic test 1",
			input:    500000.0,
			expected: 29000.0,
		},
		{
			name:     "Basic test 2",
			input:    800000.0,
			expected: 88500.0,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if actual := tax.CalculateTax(tc.input); actual != tc.expected {
				t.Errorf("Expected %f, got %f. Input was %f", tc.expected, actual, tc.input)
			}
		})
	}
}
