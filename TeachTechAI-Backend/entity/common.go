package entity

import (
	"time"

	"gorm.io/gorm"
)

type Timestamp struct {
	CreatedAt 		time.Time 	`json:"created_at" default:"CURRENT_TIMESTAMP"`
	UpdatedAt 		time.Time 	`json:"updated_at"`
	DeletedAt 		gorm.DeletedAt
}

type OTPData struct {
	PhoneNumber string `json:"phone_number,omitempty" validate:"required"`
}

type VerifyData struct {
	PhoneNumber string `json:"phone_number,omitempty" validate:"required"`
	Code string   `json:"code,omitempty" validate:"required"`
}
