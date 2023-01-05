package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

func gracefullyShutdown(e *echo.Echo) {
	go func() {
		if err := e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server not gracefully stopped error:", err.Error())
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM)
	signal.Notify(shutdown, syscall.SIGINT)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Server gracefully stopped")
	}
}

func main() {

	fmt.Println("Please use server.go for main file")

	config.LoadEnvironmentLocal()

	database.InitDB()
	defer database.DB.Close()

	h := expense.NewApplication(database.DB)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(authMiddleware())

	e.POST("/expenses", h.CreateExpenseHandler)
	e.GET("/expenses/:id", h.GetExpenseHandler)
	e.PUT("/expenses/:id", expense.PutExpenseHandler)
	e.GET("/expenses", expense.GetExpensesHandler)

	log.Println("Server started at :", os.Getenv("PORT"))

	gracefullyShutdown(e)
}
