package dto

import (
	"errors"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

const (
	// Failed
	MESSAGE_FAILED_REGISTER_USER           = "failed to create user"
	MESSAGE_FAILED_GET_USER_TOKEN          = "failed to get user token"
	MESSAGE_FAILED_GET_USER                = "failed to get user"
	MESSAGE_FAILED_LOGIN                   = "login failed"
	MESSAGE_FAILED_WRONG_EMAIL_OR_PASSWORD = "wrong email or password"
	MESSAGE_FAILED_UPDATE_USER             = "failed to update user"
	MESSAGE_FAILED_CHANGE_PASSWORD         = "failed to change password"
	MESSAGE_FAILED_RESET_PASSWORD          = "failed to reset password"
	MESSAGE_FAILED_DELETE_USER             = "failed to delete user"
	MESSAGE_FAILED_PROCESSING_REQUEST      = "failed processing request"
	MESSAGE_FAILED_DENIED_ACCESS           = "access denied"
	MESSAGE_FAILED_SEND_OTP_EMAIL          = "failed to send otp verification to email"
	MESSAGE_FAILED_VERIFY_EMAIL            = "failed to verify email"
	MESSAGE_FAILED_REFRESHING_TOKEN        = "failed to refresh token"
	MESSAGE_FAILED_LOGOUT                  = "failed to log out"
	MESSAGE_FAILED_UPLOAD_PROFILE_PICTURE  = "failed to upload picture"
	MESSAGE_FAILED_GET_PROFILE_PICTURE     = "failed to get profile picture"
	MESSAGE_FAILED_DELETE_PROFILE_PICTURE  = "failed to delete profile picture"

	// Success
	MESSAGE_SUCCESS_REGISTER_USER          = "creating user success"
	MESSAGE_SUCCESS_GET_USER               = "getting user success"
	MESSAGE_SUCCESS_LOGIN                  = "login success"
	MESSAGE_SUCCESS_UPDATE_USER            = "updating user success"
	MESSAGE_SUCCESS_CHANGE_PASSWORD        = "change password success"
	MESSAGE_SUCCESS_RESET_PASSWORD         = "reset password success"
	MESSAGE_SUCCESS_DELETE_USER            = "deleting user success"
	MESSAGE_SEND_OTP_EMAIL_SUCCESS         = "sending otp verification to email success"
	MESSAGE_SUCCESS_VERIFY_EMAIL           = "verify email success"
	MESSAGE_SUCCESS_REFRESH_TOKEN          = "refresh token success"
	MESSAGE_SUCCESS_LOGOUT                 = "logout success"
	MESSAGE_SUCCESS_UPLOAD_PROFILE_PICTURE = "upload picture success"
	MESSAGE_SUCCESS_DELETE_PROFILE_PICTURE = "deleting profile picture success"
)

var (
	ErrCreateUser             = errors.New("failed to create user")
	ErrGetAllUser             = errors.New("failed to get all user")
	ErrGetUserById            = errors.New("failed to get user by id")
	ErrGetUserByEmail         = errors.New("failed to get user by email")
	ErrEmailAlreadyExists     = errors.New("email already exist")
	ErrUpdateUser             = errors.New("failed to update user")
	ErrUserNotAdmin           = errors.New("user not admin")
	ErrUserNotFound           = errors.New("user not found")
	ErrEmailNotFound          = errors.New("email not found")
	ErrDeleteUser             = errors.New("failed to delete user")
	ErrPasswordNotMatch       = errors.New("password not match")
	ErrEmailOrPassword        = errors.New("wrong email or password")
	ErrAccountNotVerified     = errors.New("account not verified")
	ErrTokenInvalid           = errors.New("token invalid")
	ErrTokenExpired           = errors.New("token expired")
	ErrAccountAlreadyVerified = errors.New("account already verified")
)

type (
	UserCreateDTO struct {
		ID       uuid.UUID `gorm:"type:char(36);primary_key;not_null" json:"id"`
		Email    string    `json:"email" form:"email" binding:"required"`
		Name     string    `json:"name" form:"name" binding:"required"`
		Password string    `json:"password" form:"password" binding:"required"`
	}

	UserUpdateInfoDTO struct {
		ID uuid.UUID `gorm:"type:char(36);primary_key;not_null" json:"id"`
		// GoogleID     string    `gorm:"type:varchar(255);" json:"google_id"`
		Name         string    `json:"name" form:"name" binding:"required"`
		AsalInstansi string    `json:"asal_instansi" form:"asal_instansi" binding:"required"`
		DateOfBirth  time.Time `json:"date_of_birth" form:"date_of_birth" binding:"required"`
	}

	SendUserOTPByEmail struct {
		Email string `json:"email" form:"email" binding:"required"`
	}

	VerifyUserOTPByEmail struct {
		Email string `json:"email" form:"email" binding:"required"`
		OTP   string `json:"otp" form:"otp" binding:"required"`
	}

	UserUpdateEmailDTO struct {
		Email string `json:"email" form:"email" binding:"required"`
	}

	UserUpdatePhoneDTO struct {
		Phone string `json:"phone" form:"phone" binding:"required"`
	}

	UploadFileDTO struct {
		File *multipart.FileHeader `form:"file" binding:"required"`
	}

	UserChangePassword struct {
		OldPassword string `json:"old_password" form:"old_password" binding:"required"`
		NewPassword string `json:"new_password" form:"new_password" binding:"required"`
	}

	ForgotPassword struct {
		Email string `json:"email" form:"email" binding:"required"`
	}

	UserLoginDTO struct {
		Email    string `json:"email" form:"email" binding:"email"`
		Password string `json:"password" form:"password" binding:"required"`
	}

	UserLoginResponseDTO struct {
		SessionToken string `json:"session_token"`
		RefreshToken string `json:"refresh_token"`
		Role         string `json:"role"`
	}

	UserRefreshDTO struct {
		RefreshToken string `json:"refresh_token" form:"refresh_token" binding:"required"`
	}
)
