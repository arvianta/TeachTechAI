package dto

import "teach-tech-ai/entity"

type MessageRequestDTO struct {
	ConversationID string `json:"conversation_id" form:"conversation_id"`
	Topic 		   string `json:"topic" form:"topic" binding:"required"`
	Request        string `json:"request" form:"request" binding:"required"`
	AIModelName    string `json:"aimodel" form:"aimodel" binding:"required"`
}

type MessageResponseDTO struct {
	ConversationID string `json:"conversation_id"`
	Message_ID     string `json:"message_id" binding:"required"`
	Response       string `json:"response" binding:"required"`
}

type GetMessagesFromConversationDTO struct {
	ConversationID string `json:"conversation_id" form:"conversation_id"`
}

type GetMessagesFromConversationResponseDTO struct {
	Messages []entity.Message `json:"messages" binding:"required"`
}
