package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Conversation struct {
    ID        uuid.UUID      `gorm:"type:char(36);primary_key;not null" json:"id"`
    StartTime time.Time      `gorm:"type:timestamp;not null" json:"start_time"`
    EndTime   time.Time      `gorm:"type:timestamp;not null" json:"end_time"`
    UserID    uuid.UUID      `gorm:"type:char(36);not null" json:"user_id"`
    User      User           `gorm:"foreignkey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}