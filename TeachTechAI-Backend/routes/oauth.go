package routes

import (
	"teach-tech-ai/controller"

	"github.com/gin-gonic/gin"
)

func OAuthRoutes(router *gin.Engine, OAuthController controller.OAuthController) {
	oauthRoutes := router.Group("/api/oauth")
	{
		oauthRoutes.GET("/auth/callback", OAuthController.GetAuthCallbackFunction)
		oauthRoutes.GET("/logout", OAuthController.Logout)
		oauthRoutes.GET("/auth", OAuthController.Authenticate)
	}
}