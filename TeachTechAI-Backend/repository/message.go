package repository

import (
	"context"
	"teach-tech-ai/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageRepository interface {
	StoreMessage(message entity.Message) (entity.Message, error)
	GetMessagesFromConversation(ctx context.Context, conversationID uuid.UUID) ([]entity.Message, error)
}

type messageConnection struct {
	connection *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageConnection{
		connection: db,
	}
}

func (db *messageConnection) StoreMessage(message entity.Message) (entity.Message, error) {
	message.ID = uuid.New()
	uc := db.connection.Create(&message)
	if uc.Error != nil {
		return entity.Message{}, uc.Error
	}
	return message, nil
}

func (db *messageConnection) GetMessagesFromConversation(ctx context.Context, conversationID uuid.UUID) ([]entity.Message, error) {
	var messages []entity.Message
	result := db.connection.Where("conversation_id = ?", conversationID).Order("datetime ASC").Find(&messages)
	if result.Error != nil {
		return nil, result.Error
	}
	return messages, nil
}
