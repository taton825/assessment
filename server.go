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

func getExpensesHandler(c echo.Context) error {
	stmt, err := database.DB.Prepare("SELECT id, title, amount, note, tags FROM expenses")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, expense.Err{Message: "can't create prepare statement for get one expense: " + err.Error()})
	}

	rows, err := stmt.Query()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, expense.Err{Message: "can't query all expenses:" + err.Error()})
	}

	var expenses []expense.Expense

	for rows.Next() {
		var e expense.Expense
		err := rows.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, expense.Err{Message: "can't scan row into variable:" + err.Error()})
		}
		expenses = append(expenses, e)
	}
	log.Println("query all expenses success")

	return c.JSON(http.StatusOK, expenses)
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
	e.PUT("/expenses/:id", expense.PutExpenseHandler)
	e.GET("/expenses", getExpensesHandler)

	log.Println("Server started at :", os.Getenv("PORT"))
	log.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
	log.Println("Server Shutdown!! bye bye")
}
