package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func MustGetenv(k string) string {

	if os.Getenv("APP_ENV") == "production" {
		v := os.Getenv(k)
		if v == "" {
			log.Fatalf("%s environment variable not set.", k)
		}
		return v
	}

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("%s environment variable not set.", k)
	}
	return v
}
