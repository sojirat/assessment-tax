package deductions

import (
	"net/http"

	sql "github.com/sojirat/assessment-tax/sql"

	"github.com/labstack/echo/v4"
)

type UpdatePersonalDeductionInput struct {
	Amount float64 `json:"amount" validate:"required"`
}

type UpdatePersonalDeductionResponse struct {
	PersonalDeduction float64 `json:"personalDeduction"`
}

func UpdatePersonalDeductionHandler(c echo.Context) error {
	var input UpdatePersonalDeductionInput
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if input.Amount < 10000 || input.Amount > 100000 {
		return echo.NewHTTPError(http.StatusBadRequest, "Amount must be between 10,000 and 100,000")
	}

	if err := UpdatePersonalDeduction(input.Amount); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := UpdatePersonalDeductionResponse{
		PersonalDeduction: input.Amount,
	}

	return c.JSON(http.StatusOK, response)
}

func UpdatePersonalDeduction(amount float64) error {
	db := sql.Connect()
	defer db.Close()

	stmt, err := db.Prepare(`UPDATE setting SET personal = $1 WHERE id = 1;`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(amount)
	return err
}
