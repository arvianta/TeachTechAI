package main

import (
	"fmt"
	"log"
	"os"
	"teach-tech-ai/config"
	"teach-tech-ai/controller"
	"teach-tech-ai/repository"
	"teach-tech-ai/routes"
	"teach-tech-ai/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	var db *gorm.DB
	if os.Getenv("ENV") == "production" {
		db = config.ConnectWithConnector()
	} else {
		db = config.SetupDatabaseConnection()
	}

	defer config.CloseDatabaseConnection(db)

	config.AutoMigrateDatabase(db)
	fmt.Println("Migration success!")

	var (
		roleRepository 	repository.RoleRepository  = repository.NewRoleRepository(db)
		userRepository 	repository.UserRepository  = repository.NewUserRepository(db)
		
		oauthService   	service.OAuthService       = service.NewOAuthService()
		jwtService 	  	service.JWTService 		   = service.NewJWTService(userRepository, roleRepository)
		userService    	service.UserService        = service.NewUserService(userRepository, roleRepository)
		roleService    	service.RoleService        = service.NewRoleService(roleRepository)
		
		oauthController controller.OAuthController = controller.NewOAuthController(oauthService)
		roleController 	controller.RoleController  = controller.NewRoleController(roleService, userService)
		userController 	controller.UserController  = controller.NewUserController(userService, jwtService)
	)

	oauthService.InitOAuth()

	server := gin.Default()
	routes.UserRoutes(server, userController, jwtService)
	routes.RoleRoutes(server, roleController)
	routes.OAuthRoutes(server, oauthController)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server.Run("0.0.0.0:" + port)
}
