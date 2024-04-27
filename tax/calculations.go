package tax

import (
	"errors"
	"net/http"
	"sort"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sojirat/assessment-tax/sql"
)

type taxBracket struct {
	threshold float64
	rate      float64
	level     string
}

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
	Tax       float64                    `json:"tax"`
	TaxRefund float64                    `json:"taxRefund,omitempty"`
	TaxLevel  []LevelCalculationResponse `json:"taxLevel"`
}

type LevelCalculationResponse struct {
	Level string  `json:"level"`
	Tax   float64 `json:"tax"`
}

type SelectSettingResponse struct {
	Personal float64 `json:"personal"`
	KReceipt float64 `json:"k_receipt"`
}

var taxBrackets = []taxBracket{
	{threshold: 0.0, rate: 0.0, level: "0-150,000"},
	{threshold: 150000.0, rate: 0.1, level: "150,001-500,000"},
	{threshold: 500000.0, rate: 0.15, level: "500,001-1,000,000"},
	{threshold: 1000000.0, rate: 0.2, level: "1,000,001-2,000,000"},
	{threshold: 2000000.0, rate: 0.35, level: "2,000,001 ขึ้นไป"},
}

var (
	maxKReceiptAllowance = 50000.0
	personalAllowance    = 60000.0
)

const (
	baseThreshold = 150000.0

	minPersonalAllowance = 10000.0
	maxPersonalAllowance = 100000.0
	maxDonationAllowance = 100000.0
)

func CalculateTaxHandler(c echo.Context) error {
	var input TaxCalculationInput
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := ValidateInput(input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	response := CalculateTax(input.TotalIncome, input.WHT, input.Allowances)

	return c.JSON(http.StatusOK, response)
}

func ValidateInput(input TaxCalculationInput) error {
	setting, _ := SelectSetting()
	maxKReceiptAllowance = setting.KReceipt
	personalAllowance = setting.Personal

	if input.TotalIncome < 0 {
		return errors.New("total income cannot be negative")
	}

	if input.WHT < 0 || input.WHT > input.TotalIncome {
		return errors.New("invalid withholding tax")
	}

	var validAllowanceTypes = map[string]bool{
		"personal":  true,
		"donation":  true,
		"k-receipt": true,
	}

	totalPersonalAllowance := 0.0
	totalDonationAllowance := 0.0
	totalKReceiptAllowance := 0.0

	for _, allowance := range input.Allowances {
		if allowance.Amount < 0 {
			return errors.New("allowance amount cannot be negative")
		}

		if !validAllowanceTypes[allowance.AllowanceType] {
			return errors.New("invalid allowance type")
		}

		switch allowance.AllowanceType {
		case "personal":
			totalPersonalAllowance += allowance.Amount
		case "donation":
			totalDonationAllowance += allowance.Amount
		case "k-receipt":
			totalKReceiptAllowance += allowance.Amount
		}
	}

	if totalPersonalAllowance != 0 && (totalPersonalAllowance < minPersonalAllowance) {
		return errors.New("personal allowance must be at least " + strconv.FormatFloat(minPersonalAllowance, 'f', 0, 64) + " baht")
	}

	if totalPersonalAllowance > maxPersonalAllowance {
		return errors.New("personal allowance exceeds maximum " + strconv.FormatFloat(maxPersonalAllowance, 'f', 0, 64) + " baht")
	}

	if totalDonationAllowance > maxDonationAllowance {
		return errors.New("donation allowance exceeds maximum " + strconv.FormatFloat(maxDonationAllowance, 'f', 0, 64) + " baht")
	}

	if totalKReceiptAllowance > maxKReceiptAllowance {
		return errors.New("k-receipt allowance exceeds maximum " + strconv.FormatFloat(maxKReceiptAllowance, 'f', 0, 64) + " baht")
	}

	if totalKReceiptAllowance < 0 {
		return errors.New("k-receipt allowance cannot be negative")
	}

	return nil
}

func CalculateTax(totalIncome, wht float64, allowances []Allowance) TaxCalculationResponse {
	totalAllowance := 0.0
	for _, allowance := range allowances {
		totalAllowance += allowance.Amount
	}
	taxableIncome := (totalIncome - personalAllowance) - totalAllowance

	if taxableIncome <= 0 {
		return TaxCalculationResponse{}
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

	taxRate := tax

	var taxRefund float64 = 0
	if wht > tax {
		taxRefund = wht - tax
		tax = 0
	} else {
		tax -= wht
	}

	var taxBracketsInterface []LevelCalculationResponse
	for i, v := range taxBrackets {
		var isTax = 0.0
		if index-1 == i {
			isTax = taxRate
		}

		taxBracketsInterface = append(taxBracketsInterface, LevelCalculationResponse{
			Level: v.level,
			Tax:   isTax,
		})
	}

	var response TaxCalculationResponse
	response.Tax = tax
	response.TaxRefund = taxRefund
	response.TaxLevel = taxBracketsInterface

	return response
}

func SelectSetting() (SelectSettingResponse, error) {
	db := sql.Connect()
	defer db.Close()

	stmt, err := db.Prepare(`SELECT personal, k_receipt FROM setting WHERE id = 1;`)
	if err != nil {
		return SelectSettingResponse{}, err
	}
	defer stmt.Close()

	var response SelectSettingResponse
	err = stmt.QueryRow().Scan(&response.Personal, &response.KReceipt)

	return response, err
}
