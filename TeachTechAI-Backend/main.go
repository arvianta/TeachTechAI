package main

import (
	"log"
	"os"
	"teach-tech-ai/config"
	"teach-tech-ai/controller"
	"teach-tech-ai/database"
	"teach-tech-ai/repository"
	"teach-tech-ai/routes"
	"teach-tech-ai/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load(".env")
		if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
		}
	}
	
	var (
		db 				*gorm.DB 				   = config.SetupDatabaseConnection()
		roleRepository 	repository.RoleRepository  = repository.NewRoleRepository(db)
		userRepository 	repository.UserRepository  = repository.NewUserRepository(db)
		
		oauthService   	service.OAuthService       = service.NewOAuthService()
		otpService     	service.OTPService         = service.NewOTPService()
		jwtService 	  	service.JWTService 		   = service.NewJWTService(userRepository, roleRepository)
		userService    	service.UserService        = service.NewUserService(userRepository, roleRepository)
		
		oauthController controller.OAuthController = controller.NewOAuthController(oauthService)
		otpController   controller.OTPController   = controller.NewOTPController(otpService)
		userController 	controller.UserController  = controller.NewUserController(userService, jwtService)
	)
	// migrate db
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}

	// seed db
	if err := database.Seeder(db); err != nil {
		log.Fatalf("Error seeding database: %v", err)
	}

	oauthService.InitOAuth()

	server := gin.Default()
	routes.UserRoutes(server, userController, jwtService)
	routes.OAuthRoutes(server, oauthController)
	routes.OTPRoutes(server, otpController)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server.Run("0.0.0.0:" + port)
}
