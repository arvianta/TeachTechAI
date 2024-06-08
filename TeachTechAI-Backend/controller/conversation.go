package controller

import (
	"net/http"
	"teach-tech-ai/common"
	"teach-tech-ai/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ConversationController interface {
	GetConversationsFromUser(ctx *gin.Context)
	DeleteConversation(ctx *gin.Context)
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

func (cc *conversationController) DeleteConversation(ctx *gin.Context) {
	token := ctx.MustGet("token").(string)
	userID, err := cc.jwtService.GetUserIDByToken(token)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Memproses Request", "Token Tidak Valid", nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	convoID, err := uuid.Parse(ctx.Param("convoID"))
	if err != nil {
		response := common.BuildErrorResponse("Gagal Memproses Request", "ID Tidak Valid", nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if valid, err := cc.conversationService.ValidateUserConversation(userID, convoID); !valid || err != nil{
		response := common.BuildErrorResponse("Gagal Membuat Pesan", "Anda Tidak Memiliki Akses", common.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	err = cc.conversationService.DeleteConversation(convoID)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Menghapus Data", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := common.BuildResponse(true, "Berhasil Menghapus Conversation", nil)
	ctx.JSON(http.StatusOK, response)
}