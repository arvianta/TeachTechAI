package service

import (
	"teach-tech-ai/entity"
	"teach-tech-ai/repository"
	"time"

	"github.com/google/uuid"
)

type ConversationService interface {
	CreateConversation(userID uuid.UUID) (entity.Conversation, error)
}

type conversationService struct {
	conversationRepository 	repository.ConversationRepository
}

func NewConversationService(cr repository.ConversationRepository) ConversationService {
	return &conversationService{
		conversationRepository: cr,
	}
}

func (ms *conversationService) CreateConversation(userID uuid.UUID) (entity.Conversation, error) {
	conversation := entity.Conversation{
		UserID: userID,
		StartTime: time.Now(),
	}
	
	return ms.conversationRepository.StoreConversation(conversation)
}