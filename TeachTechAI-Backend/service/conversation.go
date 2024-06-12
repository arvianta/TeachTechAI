package service

import (
	"teach-tech-ai/entity"
	"teach-tech-ai/repository"
	"time"

	"github.com/google/uuid"
)

type ConversationService interface {
	CreateConversation(userID uuid.UUID, topic string) (entity.Conversation, error)
	ValidateUserConversation(userID uuid.UUID, convoID uuid.UUID) (bool, error)
	GetConversationsFromUser(userID uuid.UUID) ([]entity.Conversation, error)
	DeleteConversation(convoID uuid.UUID) error
}

type conversationService struct {
	conversationRepository repository.ConversationRepository
}

func NewConversationService(cr repository.ConversationRepository) ConversationService {
	return &conversationService{
		conversationRepository: cr,
	}
}

func (cs *conversationService) CreateConversation(userID uuid.UUID, topic string) (entity.Conversation, error) {
	conversation := entity.Conversation{
		UserID:    userID,
		StartTime: time.Now(),
		Topic:     topic,
	}

	return cs.conversationRepository.StoreConversation(conversation)
}

func (cs *conversationService) ValidateUserConversation(userID uuid.UUID, convoID uuid.UUID) (bool, error) {
	conversation, err := cs.conversationRepository.GetConversation(convoID)
	if err != nil {
		return false, err
	}

	if conversation.UserID != userID {
		return false, nil
	}

	return true, nil
}

func (cs *conversationService) GetConversationsFromUser(userID uuid.UUID) ([]entity.Conversation, error) {
	return cs.conversationRepository.GetConversationsFromUser(userID)
}

func (cs *conversationService) DeleteConversation(convoID uuid.UUID) error {
	return cs.conversationRepository.DeleteConversation(convoID)
}
