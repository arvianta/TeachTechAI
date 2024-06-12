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
	// CreateMessageStream(ctx *gin.Context)
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

	if valid, err := mc.conversationService.ValidateUserConversation(userID, convoID); !valid || err != nil {
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

// func (mc *messageController) CreateMessageStream(ctx *gin.Context) {
// 	var msg dto.MessageRequestDTO
// 	err := ctx.ShouldBind(&msg)
// 	if err != nil {
// 		response := common.BuildErrorResponse("Gagal Membuat Pesan", err.Error(), common.EmptyObj{})
// 		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
// 		return
// 	}
// 	token := ctx.MustGet("token").(string)
// 	userID, err := mc.jwtService.GetUserIDByToken(token)
// 	if err != nil {
// 		response := common.BuildErrorResponse("Gagal Memproses Request", "Token Tidak Valid", nil)
// 		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
// 		return
// 	}

// 	if msg.ConversationID == "" {
// 		conversation, err := mc.conversationService.CreateConversation(userID, msg.Topic)
// 		if err != nil {
// 			response := common.BuildErrorResponse("Gagal Membuat Pesan", err.Error(), common.EmptyObj{})
// 			ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
// 			return
// 		}
// 		msg.ConversationID = conversation.ID.String()
// 	}

// 	convoID, err := uuid.Parse(msg.ConversationID)
// 	if err != nil {
// 		response := common.BuildErrorResponse("Gagal Membuat Pesan", "Invalid conversation ID", common.EmptyObj{})
// 		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
// 		return
// 	}

// 	if valid, err := mc.conversationService.ValidateUserConversation(userID, convoID); !valid || err != nil {
// 		response := common.BuildErrorResponse("Gagal Membuat Pesan", "Anda Tidak Memiliki Akses", common.EmptyObj{})
// 		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
// 		return
// 	}

// 	// Set headers for streaming
// 	ctx.Writer.Header().Set("Content-Type", "application/json")
// 	ctx.Writer.Header().Set("Transfer-Encoding", "chunked")

// 	// Create a context with a cancel function
// 	streamCtx, cancel := context.WithCancel(ctx.Request.Context())
// 	defer cancel()

// 	// Call service layer to create message and stream response
// 	responseChan, err := mc.messageService.CreateMessageStream(streamCtx, msg)
// 	if err != nil {
// 		response := common.BuildErrorResponse("Gagal Membuat Pesan", err.Error(), common.EmptyObj{})
// 		ctx.JSON(http.StatusInternalServerError, response)
// 		return
// 	}

// 	// Stream response to client
// 	flusher, ok := ctx.Writer.(http.Flusher)
// 	if !ok {
// 		response := common.BuildErrorResponse("Streaming not supported", "Failed to stream response", common.EmptyObj{})
// 		ctx.JSON(http.StatusInternalServerError, response)
// 		return
// 	}

// 	for content := range responseChan {
// 		fmt.Fprintf(ctx.Writer, "%s", content)
// 		flusher.Flush()
// 	}
// }

func (mc *messageController) GetMessagesFromConversation(ctx *gin.Context) {
	conversationID := ctx.Param("id")

	convoID, err := uuid.Parse(conversationID)
	if err != nil {
		response := common.BuildErrorResponse("Invalid conversation ID", "Invalid ID format", common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	token := ctx.MustGet("token").(string)
	userID, err := mc.jwtService.GetUserIDByToken(token)
	if err != nil {
		response := common.BuildErrorResponse("Failed to fetch messages", "Invalid token", nil)
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}

	if valid, err := mc.conversationService.ValidateUserConversation(userID, convoID); !valid || err != nil {
		response := common.BuildErrorResponse("Failed to fetch messages", "You do not have access to this conversation", common.EmptyObj{})
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}

	var convoDTO dto.GetMessagesFromConversationDTO
	convoDTO.ConversationID = conversationID

	messagesDTO, err := mc.messageService.GetMessagesFromConversation(ctx.Request.Context(), convoDTO)
	if err != nil {
		response := common.BuildErrorResponse("Failed to fetch messages", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	// IF WAMT RETURN ERROR WHEN MESSAGES IN A CONVO IS NULL
	// if len(messagesDTO.Messages) == 0 {
	// 	response := common.BuildErrorResponse("No messages found", "The conversation does not contain any messages", common.EmptyObj{})
	// 	ctx.JSON(http.StatusNotFound, response)
	// 	return
	// }

	res := common.BuildResponse(true, "Messages fetched successfully", messagesDTO)
	ctx.JSON(http.StatusOK, res)
}
