package entity

import (
	"time"

	"github.com/google/uuid"
)

type Conversation struct {
    ID        uuid.UUID      `gorm:"type:char(36);primary_key;not null" json:"id"`
    Topic     string         `gorm:"type:varchar(255);not null" json:"topic"`
    StartTime time.Time      `gorm:"type:timestamp;not null" json:"start_time"`
    EndTime   time.Time      `gorm:"type:timestamp;" json:"end_time"`
    UserID    uuid.UUID      `gorm:"type:char(36);not null" json:"user_id"`
    User      User           `gorm:"foreignkey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
    
    Timestamp
}