package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OTP struct {
	ID        uuid.UUID `gorm:"type:varchar(255);primary_key;default:not_null" json:"id"`
	Email     string 	`gorm:"type:varchar(255);unique;not_null" json:"email"`
	OTP       string 	`gorm:"type:varchar(255);not_null" json:"otp"`
	CreatedAt time.Time `gorm:"type:timestamp;not_null" json:"created_at"`
	ExpiresAt time.Time `gorm:"type:timestamp;not_null" json:"expires_at"`
}

func (otp *OTP) BeforeCreate(tx *gorm.DB) error {
	otp.ID = uuid.New()
	return nil
}