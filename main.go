package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/sojirat/assessment-tax/admin/deductions"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	e := echo.New()

	e.POST("admin/deductions/personal", deductions.UpdatePersonalDeductionHandler)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
