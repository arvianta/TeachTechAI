package routes

import (
	"teach-tech-ai/controller"
	"teach-tech-ai/middleware"
	"teach-tech-ai/service"

	"github.com/gin-gonic/gin"
)

func ConversationRoutes(router *gin.Engine, ConversationController controller.ConversationController, jwtService service.JWTService) {
	conversationRoutes := router.Group("/api/conversation")
	{
		conversationRoutes.GET("/me", middleware.Authenticate(jwtService), ConversationController.GetConversationsFromUser)
	}
}
