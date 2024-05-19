package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AIModel struct {
    ID          uuid.UUID      `gorm:"type:char(36);primary_key;not null" json:"id"`
    Name        string         `gorm:"type:varchar(100);not null" json:"name"`
    Version     string         `gorm:"type:varchar(30);not null" json:"version"`
    Description string         `gorm:"type:text;not null" json:"description"`
    CreatedAt   time.Time      `gorm:"type:timestamp;not null" json:"created_at"`
    UpdatedAt   time.Time      `gorm:"type:timestamp;not null" json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}