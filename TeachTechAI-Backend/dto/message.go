package dto

type MessageRequest struct {
	ConversationID string `json:"conversation_id"`
	Request        string `json:"request" binding:"required"`
	AIModelName    string `json:"aimodel" binding:"required"`
}

type MessageResponse struct {
	ConversationID string `json:"conversation_id"`
	Message_ID     string `json:"message_id" binding:"required"`
	Response       string `json:"response" binding:"required"`
}
