package dto

type GenerateOTPRequest struct {
	PhoneNumber string `json:"phone_number,omitempty" binding:"required"`
}

type VerifyOTPRequest struct {
	PhoneNumber string `json:"phone_number,omitempty" validate:"required"`
	Code        string `json:"code,omitempty" validate:"required"`
}