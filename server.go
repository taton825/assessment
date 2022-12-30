package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lib/pq"
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

func putExpenseHandler(c echo.Context) error {
	rowid := c.Param("id")

	var e expense.Expense
	err := c.Bind(&e)
	if err != nil {
		return c.JSON(http.StatusBadRequest, expense.Err{Message: err.Error()})
	}

	row := database.DB.QueryRow("UPDATE expenses SET title = $2, amount = $3, note = $4, tags = $5 WHERE id = $1 RETURNING id", rowid, e.Title, e.Amount, e.Note, pq.Array(e.Tags))

	err = row.Scan(&e.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, expense.Err{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, e)
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
	e.GET("/expenses/:id", expense.GetExpenseHandler)
	e.PUT("/expenses/:id", putExpenseHandler)

	log.Println("Server started at :", os.Getenv("PORT"))
	log.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
	log.Println("Server Shutdown!! bye bye")
}
