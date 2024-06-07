package dto

import "teach-tech-ai/entity"

type MessageRequestDTO struct {
	ConversationID string `json:"conversation_id"`
	Request        string `json:"request" binding:"required"`
	AIModelName    string `json:"aimodel" binding:"required"`
}

type MessageResponseDTO struct {
	ConversationID string `json:"conversation_id"`
	Message_ID     string `json:"message_id" binding:"required"`
	Response       string `json:"response" binding:"required"`
}

type GetMessagesFromConversationDTO struct {
	ConversationID string `json:"conversation_id"`
}

type GetMessagesFromConversationResponseDTO struct {
	Messages []entity.Message `json:"messages" binding:"required"`
}
