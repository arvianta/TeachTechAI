package repository

import (
	"teach-tech-ai/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ConversationRepository interface {
	StoreConversation(conversation entity.Conversation) (entity.Conversation, error)
	GetConversation(convoID uuid.UUID) (entity.Conversation, error)
	GetConversationsFromUser(user uuid.UUID) ([]entity.Conversation, error)
	DeleteConversation(convoID uuid.UUID) error
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

func (db *conversationConnection) GetConversation(convoID uuid.UUID) (entity.Conversation, error) {
	var conversation entity.Conversation
	err := db.connection.Where("id = ?", convoID).First(&conversation).Error
	if err != nil {
		return entity.Conversation{}, err
	}
	return conversation, nil
}

func (db *conversationConnection) GetConversationsFromUser(user uuid.UUID) ([]entity.Conversation, error) {
	var conversations []entity.Conversation
	err := db.connection.Where("user_id = ?", user).Order("start_time ASC").Find(&conversations).Error
	if err != nil {
		return []entity.Conversation{}, err
	}
	return conversations, nil
}

func (db *conversationConnection) DeleteConversation(convoID uuid.UUID) error {
	uc := db.connection.Where("id = ?", convoID).Delete(&entity.Conversation{}, &convoID)
	if uc.Error != nil {
		return uc.Error
	}
	return nil
}