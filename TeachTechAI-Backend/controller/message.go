package controller

import (
	"net/http"
	"teach-tech-ai/common"
	"teach-tech-ai/dto"
	"teach-tech-ai/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	token := ctx.MustGet("token").(string)
	userID, err := mc.jwtService.GetUserIDByToken(token)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Memproses Request", "Token Tidak Valid", nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	if msg.ConversationID == "" {
		conversation, err := mc.conversationService.CreateConversation(userID, msg.Topic)
		if err != nil {
			response := common.BuildErrorResponse("Gagal Membuat Pesan", err.Error(), common.EmptyObj{})
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
		msg.ConversationID = conversation.ID.String()
	}

	convoID, err := uuid.Parse(msg.ConversationID)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Membuat Pesan", "Invalid conversation ID", common.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if valid, err := mc.conversationService.ValidateUserConversation(userID, convoID); !valid || err != nil{
		response := common.BuildErrorResponse("Gagal Membuat Pesan", "Anda Tidak Memiliki Akses", common.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
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
	conversationID := ctx.Param("id")
	if conversationID == "" {
		response := common.BuildErrorResponse("Failed to get conversation", "Invalid conversation ID", common.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	var convoID dto.GetMessagesFromConversationDTO
	convoID.ConversationID = conversationID

	messages, err := mc.messageService.GetMessagesFromConversation(ctx.Request.Context(), convoID)
	if err != nil {
		response := common.BuildErrorResponse("Failed to fetch messages", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	res := common.BuildResponse(true, "Messages fetched successfully", messages)
	ctx.JSON(http.StatusOK, res)
}