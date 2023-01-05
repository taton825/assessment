package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

var DB *sql.DB

func InitDB() {
	databaseUrl, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		log.Fatalln("Error to load DATABASE_URL from .env file")
	}

	var err error
	DB, err = sql.Open("postgres", databaseUrl)
	if err != nil {
		log.Fatalln("Connect to Database error", err)
	}

	createTb := `CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	);`
	CreateTable(DB, createTb)
}

func InitITDB() {
	var err error
	DB, err = sql.Open("postgres", "postgresql://root:root@db/expenses-db?sslmode=disable")
	if err != nil {
		log.Fatalln("Connect to Database error", err)
	}

	createTb := `CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	);`
	CreateTable(DB, createTb)
}

func CreateTable(db *sql.DB, query string) (sql.Result, error) {
	result, err := db.Exec(query)
	fmt.Println(result)
	if err != nil {
		log.Fatalln("Can't create table", err)
		return nil, err
	}
	log.Println("create table success!!")
	return result, nil
}
