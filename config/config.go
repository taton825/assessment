package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvironmentLocal() {
	// Load env file from Local env file
	_, ok := os.LookupEnv("PORT")
	if !ok {
		fmt.Println("Can't Lookup Env file with os lib")
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalln("Error to load .env file")
		}
		_, ok = os.LookupEnv("PORT")
		if !ok {
			log.Fatalln("Error Load ENV: ENVIRONMENT For Local")
		}
	}
}
