package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/taton825/assessment/config"
)

var db *sql.DB

type Expense struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Amount float64  `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}

type Err struct {
	Message string `json:"message"`
}

func createExpenseHandler(c echo.Context) error {
	var e Expense
	err := c.Bind(&e)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	fmt.Println(e)
	row := db.QueryRow("INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4) RETURNING id", e.Title, e.Amount, e.Note, pq.Array(e.Tags))

	err = row.Scan(&e.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't insert data:" + err.Error()})
	}
	return c.JSON(http.StatusCreated, e)
}

func main() {

	fmt.Println("Please use server.go for main file")
	config.LoadEnvironmentLocal()

	databaseUrl, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		log.Fatalln("Error to load DATABASE_URL from .env file")
	}

	var err error
	db, err = sql.Open("postgres", databaseUrl)
	if err != nil {
		log.Fatalln("Connect to Database error", err)
	}
	defer db.Close()

	createTb := `CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	);`
	_, err = db.Exec(createTb)
	if err != nil {
		log.Fatalln("Can't create table", err)
	}

	fmt.Println("create table success!!")

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/expenses", createExpenseHandler)

	log.Println("Server started at :", os.Getenv("PORT"))
	log.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
	log.Println("Server Shutdown!! bye bye")
}
