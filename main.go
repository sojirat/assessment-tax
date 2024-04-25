package main

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/sojirat/assessment-tax/tax"
)

func main() {
	e := echo.New()

	e.POST("tax/calculations", tax.CalculateTaxHandler)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
