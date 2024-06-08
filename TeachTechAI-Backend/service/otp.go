package service

import (
	"errors"
	"teach-tech-ai/utils"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/verify/v2"
)

type OTPService interface {
	TwilioSendOTP(phoneNumber string) (string, error)
	TwilioVerifyOTP(phoneNumber string, code string) error
}

type otpService struct {
	twilioAccountSID string
	twilioAuthToken  string
	twilioServiceSID string
}

func NewOTPService() OTPService {
	return &otpService{
		twilioAccountSID: getTwilioAccountSID(),
		twilioAuthToken:  getTwilioAuthToken(),
		twilioServiceSID: getTwilioServiceSID(),
	}
}


func getTwilioAccountSID() string {
	twilioAccountSID := utils.MustGetenv("TWILIO_ACCOUNT_SID")
	return twilioAccountSID
}

func getTwilioAuthToken() string {
	twilioAuthToken := utils.MustGetenv("TWILIO_AUTH_TOKEN")
	return twilioAuthToken
}

func getTwilioServiceSID() string {
	twilioServiceSID := utils.MustGetenv("TWILIO_SERVICE_SID")
	return twilioServiceSID
}

var client *twilio.RestClient = twilio.NewRestClientWithParams(twilio.ClientParams{
	Username: getTwilioAccountSID(),
	Password: getTwilioAuthToken(),
})


func (o *otpService) TwilioSendOTP(phoneNumber string) (string, error) {
	params := &twilioApi.CreateVerificationParams{}
	params.SetTo(phoneNumber)
	params.SetChannel("sms")

	resp, err := client.VerifyV2.CreateVerification(o.twilioServiceSID, params)
	if err != nil {
		return "", err
	}

	return *resp.Sid, nil
}

func (o *otpService) TwilioVerifyOTP(phoneNumber string, code string) error {
	params := &twilioApi.CreateVerificationCheckParams{}
	params.SetTo(phoneNumber)
	params.SetCode(code)

	resp, err := client.VerifyV2.CreateVerificationCheck(o.twilioServiceSID, params)
	if err != nil {
		return err
	}

	// BREAKING CHANGE IN THE VERIFY API
	// https://www.twilio.com/docs/verify/quickstarts/verify-totp-change-in-api-response-when-authpayload-is-incorrect
	if *resp.Status != "approved" {
		return errors.New("not a valid code")
	}

	return nil
}