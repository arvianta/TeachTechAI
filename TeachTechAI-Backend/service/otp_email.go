package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"teach-tech-ai/entity"
	"teach-tech-ai/repository"
	"teach-tech-ai/utils"
	"time"

	"gorm.io/gorm"
)

type OTPEmailService interface {
	SendOTPByEmail(ctx context.Context, email string, name string) (string, error)
	VerifyOTPByEmail(ctx context.Context, email, otp string) error
}

type otpEmailService struct {
	otpRepository repository.OTPEmailRepository
}

func NewOTPEmailService(or repository.OTPEmailRepository) OTPEmailService {
	return &otpEmailService{
		otpRepository: or,
	}
}

func GenerateOTP() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func (oes *otpEmailService) SendOTPByEmail(ctx context.Context, email string, name string) (string, error) {
	existingOTP, err := oes.otpRepository.GetValidOTPByEmail(ctx, email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}

	now, err := utils.GetCurrentTime()
	if err != nil {
		return "", err
	}

	if existingOTP != nil {
		minutes, seconds, err := utils.CalculateRemainingCooldown(existingOTP.CreatedAt)
		if err != nil {
			return "", err
		}
		if minutes > 0 || seconds > 0 {
			return "", fmt.Errorf("silakan tunggu %d menit dan %d detik sebelum meminta kode otp lagi", minutes, seconds)
		}
	}

	// Generate OTP
	randomOTP := GenerateOTP()
	expiresAt, err := utils.GetExpiryTime()
	if err != nil {
		return "", err
	}

	// Store OTP in the database
	if existingOTP != nil {
		// Replace the previous OTP record with the new one
		existingOTP.CreatedAt = now
		existingOTP.ExpiresAt = expiresAt
		existingOTP.OTP = randomOTP
		err = oes.otpRepository.UpdateOTP(ctx, *existingOTP)
		if err != nil {
			return "", err
		}
	} else {
		// Create a new OTP record
		otp := entity.OTP{
			Email:     email,
			OTP:       randomOTP,
			CreatedAt: now,
			ExpiresAt: expiresAt,
		}
		err = oes.otpRepository.CreateOTP(ctx, otp)
		if err != nil {
			return "", err
		}
	}

	// Send OTP via email
	draftEmail, err := utils.BuildMail(email, name, randomOTP)
	if err != nil {
		return "", err
	}

	err = utils.SendMail(email, draftEmail["subject"], draftEmail["body"])
	if err != nil {
		return "", err
	}

	return randomOTP, nil
}

func (oes *otpEmailService) VerifyOTPByEmail(ctx context.Context, email, otp string) error {
	storedOTP, err := oes.otpRepository.GetOTPByEmail(ctx, email)
	if err != nil {
		return err
	}

	if storedOTP.OTP != otp {
		return errors.New("invalid OTP")
	}

	if storedOTP.ExpiresAt.Before(time.Now()) {
		return errors.New("OTP has expired")
	}

	// Delete OTP from database after successful verification
	err = oes.otpRepository.DeleteOTP(ctx, email)
	if err != nil {
		return err
	}

	return nil
}
