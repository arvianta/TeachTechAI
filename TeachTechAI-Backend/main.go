package main

import (
	"net/http"
	"teach-tech-ai/common"
	"teach-tech-ai/config"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		res := common.BuildErrorResponse("Gagal Terhubung ke Server", err.Error(), common.EmptyObj{})
		(*gin.Context).JSON((&gin.Context{}), http.StatusBadGateway, res)
		return
	}
	
	var db *gorm.DB = config.SetupDatabaseConnection()
	config.AutoMigrateDatabase(db)
	config.CloseDatabaseConnection(db)
}