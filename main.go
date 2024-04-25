package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/sojirat/assessment-tax/tax"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	e := echo.New()

	e.POST("tax/calculations", tax.CalculateTaxHandler)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
