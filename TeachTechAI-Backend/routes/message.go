package routes

import (
	"teach-tech-ai/controller"
	"teach-tech-ai/middleware"
	"teach-tech-ai/service"

	"github.com/gin-gonic/gin"
)

func MessageRoutes(router *gin.Engine, MessageController controller.MessageController, jwtService service.JWTService) {
	messageRoutes := router.Group("/api/message")
	{
		messageRoutes.POST("/prompt", middleware.Authenticate(jwtService), MessageController.CreateMessage)
		// messageRoutes.POST("/prompt-stream", middleware.Authenticate(jwtService), MessageController.CreateMessageStream)
		messageRoutes.GET("/conversation/:id", middleware.Authenticate(jwtService), MessageController.GetMessagesFromConversation)
	}
}
