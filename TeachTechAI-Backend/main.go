package main

import (
	"fmt"
	"net/http"
	"os"
	"teach-tech-ai/common"
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
		res := common.BuildErrorResponse("Gagal Terhubung ke Server", err.Error(), common.EmptyObj{})
		(*gin.Context).JSON((&gin.Context{}), http.StatusBadGateway, res)
		return
	}

	var (
		db *gorm.DB = config.SetupDatabaseConnection()

		
		roleRepository repository.RoleRepository = repository.NewRoleRepository(db)
		userRepository repository.UserRepository = repository.NewUserRepository(db)
		
		jwtService 	   service.JWTService 		 = service.NewJWTService(userRepository, roleRepository)
		userService    service.UserService       = service.NewUserService(userRepository, roleRepository)
		roleService    service.RoleService       = service.NewRoleService(roleRepository)
		
		roleController controller.RoleController = controller.NewRoleController(roleService, userService)
		userController controller.UserController = controller.NewUserController(userService, jwtService)
	)

	config.AutoMigrateDatabase(db)
	// config.CloseDatabaseConnection(db)
	fmt.Println("Migration success!")

	server := gin.Default()
	routes.UserRoutes(server, userController, jwtService)
	routes.RoleRoutes(server, roleController)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	server.Run("127.0.0.1:" + port)
}
