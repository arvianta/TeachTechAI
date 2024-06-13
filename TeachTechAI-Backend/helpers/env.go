package helpers

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func MustGetenv(k string) string {
	if os.Getenv("APP_ENV") != "production" || os.Getenv("APP_ENV") != "staging" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}

	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("%s environment variable not set.", k)
	}
	return v
}

func MustGetenvInt(k string) int {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("%s environment variable not set.", k)
	}
	num, err := strconv.Atoi(v)
	if err != nil {
		log.Fatalf("Failed to convert %s to int: %v", k, err)
	}
	return num
}
