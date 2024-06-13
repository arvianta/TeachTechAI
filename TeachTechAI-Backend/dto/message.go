package dto

import (
	"errors"
	"teach-tech-ai/entity"
)

const (
	// Failed
	MESSAGE_FAILED_CREATING_MESSAGE = "failed to create message"
	MESSAGE_FAILED_GET_MESSAGE      = "failed to fetch message"
	MESSAGE_FAILED_GET_CONVO        = "failed to fetch conversation"
	MESSAGE_FAILED_DELETE_CONVO     = "failed to delete conversation"

	// Success
	MESSAGE_SUCCESS_CREATE_MESSAGE = "message created successfully"
	MESSAGE_SUCCESS_GET_MESSAGE    = "message fetched successfully"
	MESSAGE_SUCCESS_GET_CONVO      = "conversation fetched successfully"
	MESSAGE_SUCCESS_DELETE_CONVO   = "conversation deleted successfuly"
)

var (
	ErrValidateUserConversation = errors.New("user unauthorized to access this conversation")
)

type (
	MessageRequestDTO struct {
		ConversationID string `json:"conversation_id" form:"conversation_id"`
		Topic          string `json:"topic" form:"topic" binding:"required"`
		Request        string `json:"request" form:"request" binding:"required"`
		AIModelName    string `json:"aimodel" form:"aimodel" binding:"required"`
	}

	MessageResponseDTO struct {
		ConversationID string `json:"conversation_id"`
		Message_ID     string `json:"message_id"`
		Response       string `json:"response"`
	}

	GetMessagesFromConversationDTO struct {
		ConversationID string `json:"conversation_id" form:"conversation_id"`
	}

	GetMessagesFromConversationResponseDTO struct {
		Messages []entity.Message `json:"messages"`
	}
)
