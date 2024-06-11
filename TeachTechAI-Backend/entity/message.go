package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
    ID             uuid.UUID      `gorm:"type:char(36);primary_key;not null" json:"id"`
    Request        string         `gorm:"type:text;not null" json:"request"`
    Response       string         `gorm:"type:text;not null" json:"response"`
    Datetime       time.Time      `gorm:"type:timestamp;not null" json:"datetime"`
    ConversationID uuid.UUID      `gorm:"type:char(36);not null" json:"conversation_id"`
    Conversation   Conversation   `gorm:"foreignkey:ConversationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
    AIModelID      uuid.UUID      `gorm:"type:char(36);not null" json:"aimodel_id"`
    NumOfTokens    int            `gorm:"type:int;not null" json:"num_of_tokens"`
    AIModel        AIModel        `gorm:"foreignkey:AIModelID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
    DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (msg *Message) BeforeCreate(tx *gorm.DB) error {
	msg.ID = uuid.New()
	return nil
}