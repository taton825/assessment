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
	_, err = DB.Exec(createTb)
	if err != nil {
		log.Fatalln("Can't create table", err)
	}

	fmt.Println("create table success!!")
}
