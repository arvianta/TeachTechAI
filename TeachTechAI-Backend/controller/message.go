package controller

import (
	"net/http"
	"teach-tech-ai/common"
	"teach-tech-ai/dto"
	"teach-tech-ai/service"

	"github.com/gin-gonic/gin"
)

type MessageController interface {
	CreateMessage(ctx *gin.Context)
	GetMessagesFromConversation(ctx *gin.Context)
}

type messageController struct {
	jwtService          service.JWTService
	messageService      service.MessageService
	conversationService service.ConversationService
}

func NewMessageController(ms service.MessageService, cs service.ConversationService, jwts service.JWTService) MessageController {
	return &messageController{
		messageService:      ms,
		conversationService: cs,
		jwtService:          jwts,
	}
}

func (mc *messageController) CreateMessage(ctx *gin.Context) {
	var msg dto.MessageRequestDTO
	err := ctx.ShouldBind(&msg)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Membuat Pesan", err.Error(), common.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if msg.ConversationID == "" {
		token := ctx.MustGet("token").(string)
		userID, err := mc.jwtService.GetUserIDByToken(token)
		if err != nil {
			response := common.BuildErrorResponse("Gagal Memproses Request", "Token Tidak Valid", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		conversation, err := mc.conversationService.CreateConversation(userID)
		if err != nil {
			response := common.BuildErrorResponse("Gagal Membuat Pesan", err.Error(), common.EmptyObj{})
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
		msg.ConversationID = conversation.ID.String()
	}

	message, err := mc.messageService.CreateMessage(ctx.Request.Context(), msg)
	if err != nil {
		res := common.BuildErrorResponse("Gagal Membuat Pesan", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := common.BuildResponse(true, "Berhasil Membuat Pesan", message)
	ctx.JSON(http.StatusCreated, res)
}

func (mc *messageController) GetMessagesFromConversation(ctx *gin.Context) {
	var msgs dto.GetMessagesFromConversationDTO
	err := ctx.ShouldBindJSON(&msgs)
	if err != nil {
		response := common.BuildErrorResponse("Failed to get conversation", err.Error(), common.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	messages, err := mc.messageService.GetMessagesFromConversation(ctx.Request.Context(), msgs)
	if err != nil {
		response := common.BuildErrorResponse("Failed to fetch messages", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	res := common.BuildResponse(true, "Messages fetched successfully", messages)
	ctx.JSON(http.StatusOK, res)
}
