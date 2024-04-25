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
		allowances  []tax.Allowance
		expected    float64
	}{
		{
			name:        "Basic Test 1",
			totalIncome: 500000.0,
			wht:         0.0,
			allowances: []tax.Allowance{
				{
					AllowanceType: "donation",
					Amount:        100000.0,
				},
			},
			expected: 19000.0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if actual := tax.CalculateTax(tc.totalIncome, tc.wht, tc.allowances); actual != tc.expected {
				t.Errorf("Expected %f, got %f", tc.expected, actual)
			}
		})
	}
}
