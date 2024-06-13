package service

import (
	"context"
	"teach-tech-ai/dto"
	"teach-tech-ai/entity"
	"teach-tech-ai/repository"
	"time"

	"github.com/google/uuid"
)

type ConversationService interface {
	CreateConversation(ctx context.Context, userID uuid.UUID, topic string) (entity.Conversation, error)
	ValidateUserConversation(ctx context.Context, userID uuid.UUID, convoID uuid.UUID) (bool, error)
	GetConversationsFromUser(ctx context.Context, userID uuid.UUID) ([]entity.Conversation, error)
	DeleteConversation(ctx context.Context, convoID uuid.UUID) error
}

type conversationService struct {
	conversationRepository repository.ConversationRepository
}

func NewConversationService(cr repository.ConversationRepository) ConversationService {
	return &conversationService{
		conversationRepository: cr,
	}
}

func (cs *conversationService) CreateConversation(ctx context.Context, userID uuid.UUID, topic string) (entity.Conversation, error) {
	conversation := entity.Conversation{
		UserID:    userID,
		StartTime: time.Now(),
		Topic:     topic,
	}

	return cs.conversationRepository.StoreConversation(ctx, conversation)
}

func (cs *conversationService) ValidateUserConversation(ctx context.Context, userID uuid.UUID, convoID uuid.UUID) (bool, error) {
	conversation, err := cs.conversationRepository.GetConversation(ctx, convoID)
	if err != nil {
		return false, err
	}

	if conversation.UserID != userID {
		return false, dto.ErrValidateUserConversation
	}

	return true, nil
}

func (cs *conversationService) GetConversationsFromUser(ctx context.Context, userID uuid.UUID) ([]entity.Conversation, error) {
	return cs.conversationRepository.GetConversationsFromUser(ctx, userID)
}

func (cs *conversationService) DeleteConversation(ctx context.Context, convoID uuid.UUID) error {
	return cs.conversationRepository.DeleteConversation(ctx, convoID)
}
