package entity

import (
	"github.com/google/uuid"
)

type Role struct {
	ID          uuid.UUID `gorm:"type:char(36);primary_key;default:uuid_generate_v4();not_null" json:"id"`
	Name        string    `gorm:"type:varchar(20);not_null" json:"name"`
	Description string    `gorm:"type:text;" json:"description"`

	Timestamp
}
