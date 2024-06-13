package controller

import (
	"net/http"
	"teach-tech-ai/dto"
	"teach-tech-ai/service"
	"teach-tech-ai/utils"

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
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_PROCESSING_REQUEST, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	conversations, err := cc.conversationService.GetConversationsFromUser(ctx, userID)
	if err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_GET_CONVO, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := utils.BuildSuccessResponse(dto.MESSAGE_SUCCESS_GET_CONVO, conversations)
	ctx.JSON(http.StatusOK, response)
}

func (cc *conversationController) DeleteConversation(ctx *gin.Context) {
	token := ctx.MustGet("token").(string)
	userID, err := cc.jwtService.GetUserIDByToken(token)
	if err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_PROCESSING_REQUEST, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	convoID, err := uuid.Parse(ctx.Param("convoID"))
	if err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_DELETE_CONVO, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if valid, err := cc.conversationService.ValidateUserConversation(ctx, userID, convoID); !valid || err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_DELETE_CONVO, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	err = cc.conversationService.DeleteConversation(ctx, convoID)
	if err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_DELETE_CONVO, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := utils.BuildSuccessResponse(dto.MESSAGE_SUCCESS_DELETE_CONVO, utils.EmptyObj{})
	ctx.JSON(http.StatusOK, response)
}
