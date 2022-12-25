package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	// Load env file from Local env file
	fmt.Println("Please use server.go for main file")
	port, ok := os.LookupEnv("PORT")
	if !ok {
		fmt.Println("Can't Lookup Env file with os lib")
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalln("Error to load .env file")
		}
		port, ok = os.LookupEnv("PORT")
		if !ok {
			log.Fatalln("Error Load ENV: ENVIRONMENT For Local")
		}
	}
	fmt.Println("start at port:", port)

	databaseUrl, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		log.Fatalln("Error to load DATABASE_URL from .env file")
	}

	db, err := sql.Open("postgres", databaseUrl)
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
}
