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
	SendOTPByEmail(ctx context.Context, email string) (string, error)
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

const (
	otpCooldown = 2 * time.Minute
	otpExpiry   = 5 * time.Minute
)

func (oes *otpEmailService) SendOTPByEmail(ctx context.Context, email string) (string, error) {
	existingOTP, err := oes.otpRepository.GetValidOTPByEmail(ctx, email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}

	if existingOTP != nil {
		timeSinceCreation := time.Since(existingOTP.CreatedAt)
		remainingCooldown := otpCooldown - timeSinceCreation
		if remainingCooldown > 0 {
			remainingCooldownSeconds := int(remainingCooldown.Seconds())
			minutes := remainingCooldownSeconds / 60
			seconds := remainingCooldownSeconds % 60
			return "", fmt.Errorf("please wait for %d minute(s) and %d second(s) before requesting another OTP", minutes, seconds)
		}
	}

	// Generate OTP
	randomOTP := GenerateOTP()
	expiresAt := time.Now().Add(otpExpiry)

	otp := entity.OTP{
		Email:     email,
		OTP:       randomOTP,
		ExpiresAt: expiresAt,
	}

	// Store OTP in the database
	if existingOTP != nil {
		// Replace the previous OTP record with the new one
		otp.ID = existingOTP.ID
		otp.CreatedAt = time.Now()
		err = oes.otpRepository.UpdateOTP(ctx, otp)
		if err != nil {
			return "", err
		}
	} else {
		// Create a new OTP record
		err = oes.otpRepository.CreateOTP(ctx, otp)
		if err != nil {
			return "", err
		}
	}

	// Send OTP via email
	subject := "TeachTechAI OTP Verification"
	body := fmt.Sprintf("Your OTP for verification: %s", randomOTP)

	err = utils.SendMail(email, subject, body)
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
