package tax

import (
	"encoding/csv"
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type CSVFileResponse struct {
	Taxes []CSVFileDetailResponse `json:"taxes"`
}

type CSVFileDetailResponse struct {
	TotalIncome float64 `json:"totalIncome"`
	Tax         float64 `json:"tax"`
}

func ReadCSVFileHandler(c echo.Context) error {

	data, err := ReadCSVFile(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var responseDetail []CSVFileDetailResponse
	for _, d := range data {
		if err := validateInput(d); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		result := CalculateTax(d.TotalIncome, d.WHT, d.Allowances)

		responseDetail = append(responseDetail, CSVFileDetailResponse{
			TotalIncome: d.TotalIncome,
			Tax:         result.Tax,
		})
	}

	response := CSVFileResponse{
		Taxes: responseDetail,
	}

	return c.JSON(http.StatusOK, response)
}

func ReadCSVFile(c echo.Context) ([]TaxCalculationInput, error) {
	form, err := c.FormFile("taxFile")
	if err != nil {
		return nil, err
	}

	if form.Filename != "taxes.csv" {
		return nil, errors.New("filename must be taxes.csv")
	}

	f, err := form.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)

	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	records = records[1:]
	for _, record := range records {
		for _, value := range record {
			if value == "" {
				return nil, errors.New("value must not be empty")
			}

			_, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, errors.New("value must be number only")
			}
		}
	}

	result := make([]TaxCalculationInput, 0)

	for _, record := range records {
		totalIncome, _ := strconv.ParseFloat(record[0], 64)
		wht, _ := strconv.ParseFloat(record[1], 64)
		allowances, _ := strconv.ParseFloat(record[2], 64)

		result = append(result, TaxCalculationInput{
			TotalIncome: totalIncome,
			WHT:         wht,
			Allowances:  []Allowance{{AllowanceType: "donation", Amount: allowances}},
		})
	}

	return result, nil
}
