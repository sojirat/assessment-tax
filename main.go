package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/sojirat/assessment-tax/tax"

	_ "github.com/lib/pq"
)

func main() {
	e := echo.New()
	// e.GET("/", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	// })

	// connectionString := "postgres://aesucrip:QAZ0xl8zt4wZZim5z1W_Qly6bM5FHDSK@rain.db.elephantsql.com/aesucrip"
	connectionString := "postgres://postgres:postgres@localhost:5432/ktaxes?sslmode=disable"

	// fmt.Println(os.Getenv("DATABASE_URL"))

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		e.Logger.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		e.Logger.Fatal("Error pinging the database:", err)
	}

	fmt.Println("Successfully connected to the database!")

	e.POST("tax/calculations", tax.CalculateTaxHandler)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
