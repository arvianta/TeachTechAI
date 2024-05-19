package entity

import (
	"github.com/google/uuid"
)

type Role struct {
    ID          uuid.UUID `gorm:"type:char(36);primary_key;not_null" json:"id"`
    Name        string    `gorm:"type:varchar(20);not_null" json:"name"`
    Description string    `gorm:"type:text;not_null" json:"description"`
}