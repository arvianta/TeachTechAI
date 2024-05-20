package config

import (
	"fmt"
	"teach-tech-ai/entity"

	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDatabaseConnection() *gorm.DB {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbTimezone := os.Getenv("DB_TIMEZONE")

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=%v", dbHost, dbUser, dbPass, dbName, dbPort, dbTimezone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return db
} 

func AutoMigrateDatabase(db *gorm.DB) {
	if err := db.AutoMigrate(
		entity.Role{},
		entity.User{},
		entity.Conversation{},
		entity.AIModel{},
		entity.Message{},
	); err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func CloseDatabaseConnection(db *gorm.DB) {
	dbSQL, err := db.DB()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	dbSQL.Close()
}