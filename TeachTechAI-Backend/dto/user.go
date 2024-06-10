package dto

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type UserCreateDTO struct {
	ID             uuid.UUID `gorm:"type:char(36);primary_key;not_null" json:"id"`
	Email          string    `json:"email" form:"email" binding:"required"`
	Name           string    `json:"name" form:"name" binding:"required"`
	Password       string    `json:"password" form:"password" binding:"required"`
}

type UserUpdateInfoDTO struct {
	ID             uuid.UUID `gorm:"type:char(36);primary_key;not_null" json:"id"`
	GoogleID       string    `gorm:"type:varchar(255);" json:"google_id"`
	Name           string    `json:"name" binding:"required"`
	AsalInstansi   string    `json:"asal_instansi" binding:"required"`
	DateOfBirth    time.Time `json:"date_of_birth" binding:"required"`
}

type UserUpdateEmailDTO struct {
	Email string `json:"email" form:"email" binding:"required"`
}

type UserUpdatePhoneDTO struct {
	Phone string `json:"phone" form:"phone" binding:"required"`
}

type UploadFileDTO struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type UserLoginDTO struct {
	Email    string `json:"email" form:"email" binding:"email"`
	Password string `json:"password" form:"password" binding:"required"`
}

type UserRefreshDTO struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}