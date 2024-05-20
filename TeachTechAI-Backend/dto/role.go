package dto

import (
	"github.com/google/uuid"
)

type RoleCreateDto struct {
	ID          uuid.UUID `gorm:"type:char(36);primary_key;not_null" json:"id"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
}
