package tax_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/sojirat/assessment-tax/tax"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		name        string
		totalIncome float64
		wht         float64
		allowances  []tax.Allowance
		expected    tax.TaxCalculationResponse
	}{
		{
			name:        "Valid tax calculation",
			totalIncome: 500000.0,
			wht:         0.0,
			allowances: []tax.Allowance{
				{
					AllowanceType: "donation",
					Amount:        100000.0,
				},
			},
			expected: tax.TaxCalculationResponse{
				Tax: 19000.0,
				TaxLevel: []tax.LevelCalculationResponse{
					{
						Level: "0-150,000",
						Tax:   0.0,
					},
					{
						Level: "150,001-500,000",
						Tax:   19000.0,
					},
					{
						Level: "500,001-1,000,000",
						Tax:   0.0,
					},
					{
						Level: "1,000,001-2,000,000",
						Tax:   0.0,
					},
					{
						Level: "2,000,001 ขึ้นไป",
						Tax:   0.0,
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := tax.CalculateTax(tc.totalIncome, tc.wht, tc.allowances)

			sort.Slice(actual.TaxLevel, func(i, j int) bool {
				return actual.TaxLevel[i].Level < actual.TaxLevel[j].Level
			})

			sort.Slice(tc.expected.TaxLevel, func(i, j int) bool {
				return tc.expected.TaxLevel[i].Level < tc.expected.TaxLevel[j].Level
			})

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Expected %v, got %v", tc.expected, actual)
			}
		})
	}

}
