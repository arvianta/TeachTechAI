package repository

import (
	"context"
	"teach-tech-ai/entity"
	"time"

	"gorm.io/gorm"
)

type OTPEmailRepository interface {
	CreateOTP(ctx context.Context, otp entity.OTP) error
	GetOTPByEmail(ctx context.Context, email string) (entity.OTP, error)
	UpdateOTP(ctx context.Context, otp entity.OTP) error
	DeleteOTP(ctx context.Context, email string) error
	GetValidOTPByEmail(ctx context.Context, email string) (*entity.OTP, error)
}

type otpEmailRepository struct {
	connection *gorm.DB
}

func NewOTPEmailRepository(db *gorm.DB) OTPEmailRepository {
	return &otpEmailRepository{
		connection: db,
	}
}

func (db *otpEmailRepository) CreateOTP(ctx context.Context, otp entity.OTP) error {
	err := db.connection.WithContext(ctx).Create(&otp).Error
	return err
}

func (db *otpEmailRepository) GetOTPByEmail(ctx context.Context, email string) (entity.OTP, error) {
	var otp entity.OTP
	err := db.connection.WithContext(ctx).Where("email = ?", email).First(&otp).Error
	return otp, err
}

func (db *otpEmailRepository) UpdateOTP(ctx context.Context, otp entity.OTP) error {
	result := db.connection.WithContext(ctx).Model(&entity.OTP{}).Where("id = ?", otp.ID).Updates(map[string]interface{}{
		"otp":        otp.OTP,
		"created_at": otp.CreatedAt,
		"expires_at": otp.ExpiresAt,
	})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (db *otpEmailRepository) DeleteOTP(ctx context.Context, email string) error {
	var otp entity.OTP
	err := db.connection.WithContext(ctx).Where("email = ?", email).Delete(&otp).Error
	return err
}

func (db *otpEmailRepository) GetValidOTPByEmail(ctx context.Context, email string) (*entity.OTP, error) {
	var otp entity.OTP
	err := db.connection.WithContext(ctx).Where("email = ? AND expires_at > ?", email, time.Now()).First(&otp).Error
	if err != nil {
		return nil, err
	}
	return &otp, nil
}
