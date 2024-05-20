package dto

import (
	"time"

	"github.com/google/uuid"
)

type UserCreateDto struct {
	ID             uuid.UUID `gorm:"type:char(36);primary_key;not_null" json:"id"`
	GoogleID       string    `gorm:"type:varchar(255);" json:"google_id"`
	Email          string    `json:"email" binding:"required"`
	Name           string    `json:"name" binding:"required"`
	Phone          string    `json:"phone" binding:"required"`
	Password       string    `json:"password" binding:"required"`
	ProfilePicture string    `json:"profile_picture" binding:"required"`
	AsalInstansi   string    `json:"asal_instansi" binding:"required"`
	DateOfBirth    time.Time `json:"date_of_birth" binding:"required"`
}

// type UserUpdateDto struct {
// 	ID       uuid.UUID `gorm:"primary_key" json:"id" form:"id"`
// 	Name     string    `json:"name" form:"name"`
// 	Email    string    `json:"email" form:"email"`
// 	NoTelp   string    `json:"no_telp" form:"no_telp"`
// 	Password string    `json:"password" form:"password"`
// }

// type UserLoginDTO struct {
// 	Email    string `json:"email" form:"email" binding:"email"`
// 	Password string `json:"password" form:"password" binding:"required"`
// }
