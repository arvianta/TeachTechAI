package repository

import (
	"teach-tech-ai/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageRepository interface {
	StoreMessage(message entity.Message) (entity.Message, error)
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