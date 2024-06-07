package repository

import (
	"teach-tech-ai/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ConversationRepository interface {
	StoreConversation(conversation entity.Conversation) (entity.Conversation, error)
}

type conversationConnection struct {
	connection *gorm.DB
}

func NewconversationRepository(db *gorm.DB) ConversationRepository {
	return &conversationConnection{
		connection: db,
	}
}

func (db *conversationConnection) StoreConversation(conversation entity.Conversation) (entity.Conversation, error) {
	conversation.ID = uuid.New()
	uc := db.connection.Create(&conversation)
	if uc.Error != nil {
		return entity.Conversation{}, uc.Error
	}
	return conversation, nil
}