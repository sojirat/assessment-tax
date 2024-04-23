package tax_test

import (
	"testing"

	"github.com/sojirat/assessment-tax/tax"
)

type Allowance struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}

func TestCalculateTax(t *testing.T) {
	tests := []struct {
		name        string
		totalIncome float64
		wht         float64
		allowances  []tax.Allowance
		expectedTax float64
	}{
		{
			name:        "No income",
			totalIncome: 0,
			wht:         0,
			allowances:  []tax.Allowance{{AllowanceType: "personal", Amount: 0}},
			expectedTax: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actualTax := tax.CalculateTax(tc.totalIncome, tc.wht, tc.allowances)
			if actualTax != tc.expectedTax {
				t.Errorf("Expected tax: %f, got: %f", tc.expectedTax, actualTax)
			}
		})
	}
}
