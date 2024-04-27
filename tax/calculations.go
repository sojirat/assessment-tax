package tax

import (
	"errors"
	"net/http"
	"sort"

	"github.com/labstack/echo/v4"
)

type taxBracket struct {
	threshold float64
	rate      float64
}

var taxBrackets = []taxBracket{
	{threshold: 150000.0, rate: 0.1},
	{threshold: 500000.0, rate: 0.15},
	{threshold: 1000000.0, rate: 0.2},
	{threshold: 2000000.0, rate: 0.35},
}

const (
	baseThreshold     = 150000.0
	personalAllowance = 60000.0
)

type TaxCalculationInput struct {
	TotalIncome float64     `json:"totalIncome"`
	WHT         float64     `json:"wht"`
	Allowances  []Allowance `json:"allowances"`
}

type Allowance struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}

type TaxCalculationResponse struct {
	Tax float64 `json:"tax"`
}

func CalculateTaxHandler(c echo.Context) error {
	var input TaxCalculationInput
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := validateInput(input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tax := CalculateTax(input.TotalIncome)
	response := TaxCalculationResponse{Tax: tax}

	return c.JSON(http.StatusOK, response)
}

func validateInput(input TaxCalculationInput) error {
	if input.TotalIncome < 0 {
		return errors.New("total income cannot be negative")
	}

	return nil
}

func CalculateTax(totalIncome float64) float64 {
	taxableIncome := (totalIncome - personalAllowance)

	if taxableIncome <= 0 {
		return 0
	}

	sort.Slice(taxBrackets, func(i, j int) bool {
		return taxBrackets[i].threshold < taxBrackets[j].threshold
	})

	index := sort.Search(len(taxBrackets), func(i int) bool {
		return taxBrackets[i].threshold > taxableIncome
	})

	var tax float64
	if index > 0 {
		preiousRate := taxBrackets[index-1].rate
		tax = (taxableIncome - baseThreshold) * preiousRate
	}

	if tax < 0 {
		tax = 0
	}

	return tax
}
