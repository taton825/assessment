package main

import (
	"fmt"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/taton825/assessment/config"
	"github.com/taton825/assessment/database"
	"github.com/taton825/assessment/expense"
)

func main() {

	fmt.Println("Please use server.go for main file")

	config.LoadEnvironmentLocal()

	database.InitDB()
	defer database.DB.Close()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/expenses", expense.CreateExpenseHandler)

	log.Println("Server started at :", os.Getenv("PORT"))
	log.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
	log.Println("Server Shutdown!! bye bye")
}
