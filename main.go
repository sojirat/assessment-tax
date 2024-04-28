package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/sojirat/assessment-tax/admin/deductions"
	"github.com/sojirat/assessment-tax/middleware"
	"github.com/sojirat/assessment-tax/tax"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	e := echo.New()
	e.Logger.SetLevel(log.INFO)

	e.POST("tax/calculations", tax.CalculateTaxHandler)
	e.POST("tax/calculations/upload-csv", tax.ReadCSVFileHandler)

	e.PATCH("/admin/deductions/personal", deductions.UpdatePersonalDeductionHandler, middleware.AuthMiddleware)
	e.PATCH("/admin/deductions/k-receipt", deductions.UpdateKReceiptDeductionHandler, middleware.AuthMiddleware)

	go func() {
		if err := e.Start(":" + os.Getenv("PORT")); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
