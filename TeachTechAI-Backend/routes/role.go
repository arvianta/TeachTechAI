package routes

import (
	"teach-tech-ai/controller"
	// "teach-tech-ai/middleware"
	// "teach-tech-ai/service"

	"github.com/gin-gonic/gin"
)

func RoleRoutes(router *gin.Engine, RoleController controller.RoleController) {
	roleRoutes := router.Group("/api/role")
	{
		roleRoutes.POST("/create", RoleController.CreateRole)
		// roleRoutes.GET("", middleware.Authenticate(jwtService), UserController.GetAllUser)
		// roleRoutes.POST("/login", UserController.LoginUser)
		// roleRoutes.DELETE("/", middleware.Authenticate(jwtService), UserController.DeleteUser)
		// roleRoutes.PUT("/", middleware.Authenticate(jwtService), UserController.UpdateUser)
		// roleRoutes.GET("/me", middleware.Authenticate(jwtService), UserController.MeUser)
	}
}
