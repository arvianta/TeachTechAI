package controller

import (
	"net/http"
	"teach-tech-ai/dto"
	"teach-tech-ai/service"
	"teach-tech-ai/utils"

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
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	token := ctx.MustGet("token").(string)
	userID, err := mc.jwtService.GetUserIDByToken(token)
	if err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_PROCESSING_REQUEST, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	if msg.ConversationID == "" {
		conversation, err := mc.conversationService.CreateConversation(ctx.Request.Context(), userID, msg.Topic)
		if err != nil {
			response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_CREATING_MESSAGE, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
		msg.ConversationID = conversation.ID.String()
	}

	convoID, err := uuid.Parse(msg.ConversationID)
	if err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_CREATING_MESSAGE, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if valid, err := mc.conversationService.ValidateUserConversation(ctx.Request.Context(), userID, convoID); !valid || err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_CREATING_MESSAGE, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	message, err := mc.messageService.CreateMessage(ctx.Request.Context(), msg)
	if err != nil {
		res := utils.BuildErrorResponse(dto.MESSAGE_FAILED_CREATING_MESSAGE, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildSuccessResponse(dto.MESSAGE_SUCCESS_CREATE_MESSAGE, message)
	ctx.JSON(http.StatusCreated, res)
}

// func (mc *messageController) CreateMessageStream(ctx *gin.Context) {
// 	var msg dto.MessageRequestDTO
// 	err := ctx.ShouldBind(&msg)
// 	if err != nil {
// 		response := utils.BuildErrorResponse("Gagal Membuat Pesan", err.Error(), utils.EmptyObj{})
// 		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
// 		return
// 	}
// 	token := ctx.MustGet("token").(string)
// 	userID, err := mc.jwtService.GetUserIDByToken(token)
// 	if err != nil {
// 		response := utils.BuildErrorResponse("Gagal Memproses Request", "Token Tidak Valid", nil)
// 		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
// 		return
// 	}

// 	if msg.ConversationID == "" {
// 		conversation, err := mc.conversationService.CreateConversation(userID, msg.Topic)
// 		if err != nil {
// 			response := utils.BuildErrorResponse("Gagal Membuat Pesan", err.Error(), utils.EmptyObj{})
// 			ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
// 			return
// 		}
// 		msg.ConversationID = conversation.ID.String()
// 	}

// 	convoID, err := uuid.Parse(msg.ConversationID)
// 	if err != nil {
// 		response := utils.BuildErrorResponse("Gagal Membuat Pesan", "Invalid conversation ID", utils.EmptyObj{})
// 		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
// 		return
// 	}

// 	if valid, err := mc.conversationService.ValidateUserConversation(userID, convoID); !valid || err != nil {
// 		response := utils.BuildErrorResponse("Gagal Membuat Pesan", "Anda Tidak Memiliki Akses", utils.EmptyObj{})
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
// 		response := utils.BuildErrorResponse("Gagal Membuat Pesan", err.Error(), utils.EmptyObj{})
// 		ctx.JSON(http.StatusInternalServerError, response)
// 		return
// 	}

// 	// Stream response to client
// 	flusher, ok := ctx.Writer.(http.Flusher)
// 	if !ok {
// 		response := utils.BuildErrorResponse("Streaming not supported", "Failed to stream response", utils.EmptyObj{})
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
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_GET_MESSAGE, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	token := ctx.MustGet("token").(string)
	userID, err := mc.jwtService.GetUserIDByToken(token)
	if err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_GET_MESSAGE, err.Error(), nil)
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}

	if valid, err := mc.conversationService.ValidateUserConversation(ctx.Request.Context(), userID, convoID); !valid || err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_GET_MESSAGE, err.Error(), nil)
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}

	var convoDTO dto.GetMessagesFromConversationDTO
	convoDTO.ConversationID = conversationID

	messagesDTO, err := mc.messageService.GetMessagesFromConversation(ctx.Request.Context(), convoDTO)
	if err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_GET_MESSAGE, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	res := utils.BuildSuccessResponse(dto.MESSAGE_SUCCESS_GET_MESSAGE, messagesDTO)
	ctx.JSON(http.StatusOK, res)
}
