package dto

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type UserCreateDTO struct {
	ID       uuid.UUID `gorm:"type:char(36);primary_key;not_null" json:"id"`
	Email    string    `json:"email" form:"email" binding:"required"`
	Name     string    `json:"name" form:"name" binding:"required"`
	Password string    `json:"password" form:"password" binding:"required"`
}

type UserUpdateInfoDTO struct {
	ID           uuid.UUID `gorm:"type:char(36);primary_key;not_null" json:"id"`
	GoogleID     string    `gorm:"type:varchar(255);" json:"google_id"`
	Name         string    `json:"name" binding:"required"`
	AsalInstansi string    `json:"asal_instansi" binding:"required"`
	DateOfBirth  time.Time `json:"date_of_birth" binding:"required"`
}

type SendUserOTPByEmail struct {
	Email string `json:"email" form:"email" binding:"required"`
}

type VerifyUserOTPByEmail struct {
	Email string `json:"email" form:"email" binding:"required"`
	OTP   string `json:"otp" form:"otp" binding:"required"`
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

type UserChangePassword struct {
	OldPassword string `json:"old_password" form:"old_password" binding:"required"`
	NewPassword string `json:"new_password" form:"new_password" binding:"required"`
}

type ForgotPassword struct {
	Email string `json:"email" form:"email" binding:"required"`
}

type UserLoginDTO struct {
	Email    string `json:"email" form:"email" binding:"email"`
	Password string `json:"password" form:"password" binding:"required"`
}

type UserLoginResponseDTO struct {
	SessionToken string `json:"session_token"`
	RefreshToken string `json:"refresh_token"`
	Role         string `json:"role"`
}

type UserRefreshDTO struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
