package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/taton825/assessment/config"
	"github.com/taton825/assessment/database"
	"github.com/taton825/assessment/expense"
)

func authMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Header.Get("Authorization") == os.Getenv("AUTH_TOKEN") {
				return next(c)
			} else {
				return &echo.HTTPError{
					Code:    http.StatusUnauthorized,
					Message: "Unauthorized",
				}
			}
		}
	}
}
func main() {

	fmt.Println("Please use server.go for main file")

	config.LoadEnvironmentLocal()

	database.InitDB()
	defer database.DB.Close()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(authMiddleware())

	e.POST("/expenses", expense.CreateExpenseHandler)

	log.Println("Server started at :", os.Getenv("PORT"))
	log.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
	log.Println("Server Shutdown!! bye bye")
}
