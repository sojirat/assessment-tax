package tax_test

import (
	"testing"

	"github.com/sojirat/assessment-tax/tax"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		name        string
		totalIncome float64
		wht         float64
		expected    float64
	}{
		{
			name:        "Basic test 1",
			totalIncome: 500000.0,
			wht:         25000.0,
			expected:    4000.0,
		},
		{
			name:        "Basic test 2",
			totalIncome: 800000.0,
			wht:         25000.0,
			expected:    63500.0,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if actual := tax.CalculateTax(tc.totalIncome, tc.wht); actual != tc.expected {
				t.Errorf("Expected %f, got %f", tc.expected, actual)
			}
		})
	}
}
