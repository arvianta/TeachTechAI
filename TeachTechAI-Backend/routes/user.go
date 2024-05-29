package routes

import (
	"teach-tech-ai/controller"
	"teach-tech-ai/middleware"
	"teach-tech-ai/service"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, UserController controller.UserController, jwtService service.JWTService) {
	userRoutes := router.Group("/api/user")
	{
		userRoutes.POST("/register", UserController.RegisterUser)
		userRoutes.GET("", middleware.Authenticate(jwtService), UserController.GetAllUser)
		userRoutes.POST("/login", UserController.LoginUser)
		userRoutes.DELETE("/delete", middleware.Authenticate(jwtService), UserController.DeleteUser)
		userRoutes.PUT("/", middleware.Authenticate(jwtService), UserController.UpdateUser)
		userRoutes.GET("/me", middleware.Authenticate(jwtService), UserController.MeUser)
		userRoutes.POST("/refresh", UserController.RefreshUser)
		userRoutes.POST("/logout", middleware.Authenticate(jwtService), UserController.Logout)
	}
}
