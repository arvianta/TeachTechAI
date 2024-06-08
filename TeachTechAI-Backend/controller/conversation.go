package controller

import (
	"net/http"
	"teach-tech-ai/common"
	"teach-tech-ai/service"

	"github.com/gin-gonic/gin"
)

type ConversationController interface {
	GetConversationsFromUser(ctx *gin.Context)
}

type conversationController struct {
	jwtService          service.JWTService
	conversationService service.ConversationService
}

func NewConversationController(cs service.ConversationService, jwts service.JWTService) ConversationController {
	return &conversationController{
		conversationService: cs,
		jwtService:          jwts,
	}
}

func (cc *conversationController) GetConversationsFromUser(ctx *gin.Context) {
	token := ctx.MustGet("token").(string)
	userID, err := cc.jwtService.GetUserIDByToken(token)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Memproses Request", "Token Tidak Valid", nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	conversations, err := cc.conversationService.GetConversationsFromUser(userID)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Mengambil Data", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := common.BuildResponse(true, "OK", conversations)
	ctx.JSON(http.StatusOK, response)
}