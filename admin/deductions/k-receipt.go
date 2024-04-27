package deductions

import (
	"net/http"

	"github.com/labstack/echo/v4"
	sql "github.com/sojirat/assessment-tax/sql"
)

type UpdateKReceiptDeductionInput struct {
	Amount float64 `json:"amount" validate:"required"`
}

type UpdateKReceiptDeductionResponse struct {
	KReceiptDeduction float64 `json:"kReceipt"`
}

func UpdateKReceiptDeductionHandler(c echo.Context) error {
	var input UpdateKReceiptDeductionInput
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if input.Amount < 0 || input.Amount > 100000 {
		return echo.NewHTTPError(http.StatusBadRequest, "Amount must be between 0 and 100,000")
	}

	if err := UpdateKReceiptDeduction(input.Amount); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := UpdateKReceiptDeductionResponse{
		KReceiptDeduction: input.Amount,
	}

	return c.JSON(http.StatusOK, response)
}

func UpdateKReceiptDeduction(amount float64) error {
	db := sql.Connect()
	defer db.Close()

	stmt, err := db.Prepare(`UPDATE setting SET k_receipt = $1 WHERE id = 1;`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(amount)
	return err
}
