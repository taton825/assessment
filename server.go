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

func getExpenseHandler(c echo.Context) error {
	rowid := c.Param("id")

	stmt, err := database.DB.Prepare("SELECT id, title, amount, note, tags FROM expenses WHERE id = $1")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, expense.Err{Message: "can't create prepare statement for get one expense: " + err.Error()})
	}

	row := stmt.QueryRow(rowid)
	// log.Println("row:" + row.Err().Error())
	var e expense.Expense
	err = row.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, expense.Err{Message: "can't scan row into variable: " + err.Error()})
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
	e.GET("/expenses/:id", getExpenseHandler)

	log.Println("Server started at :", os.Getenv("PORT"))
	log.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
	log.Println("Server Shutdown!! bye bye")
}
