package entity

import (
	"github.com/google/uuid"
)

type AIModel struct {
    ID          uuid.UUID      `gorm:"type:char(36);primary_key;default:uuid_generate_v4();not_null" json:"id"`
    Name        string         `gorm:"type:varchar(100);not null" json:"name"`
    Version     string         `gorm:"type:varchar(30);not null" json:"version"`
    Description string         `gorm:"type:text;not null" json:"description"`
    
    Timestamp
}