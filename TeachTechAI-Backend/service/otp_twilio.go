package service

import (
	"errors"
	"teach-tech-ai/helpers"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/verify/v2"
)

type OTPTwilioService interface {
	TwilioSendOTP(phoneNumber string) (string, error)
	TwilioVerifyOTP(phoneNumber string, code string) error
}

type otpTwilioService struct {
	twilioAccountSID string
	twilioAuthToken  string
	twilioServiceSID string
}

func NewOTPTwilioService() OTPTwilioService {
	return &otpTwilioService{
		twilioAccountSID: getTwilioAccountSID(),
		twilioAuthToken:  getTwilioAuthToken(),
		twilioServiceSID: getTwilioServiceSID(),
	}
}

func getTwilioAccountSID() string {
	twilioAccountSID := helpers.MustGetenv("TWILIO_ACCOUNT_SID")
	return twilioAccountSID
}

func getTwilioAuthToken() string {
	twilioAuthToken := helpers.MustGetenv("TWILIO_AUTH_TOKEN")
	return twilioAuthToken
}

func getTwilioServiceSID() string {
	twilioServiceSID := helpers.MustGetenv("TWILIO_SERVICE_SID")
	return twilioServiceSID
}

var client *twilio.RestClient = twilio.NewRestClientWithParams(twilio.ClientParams{
	Username: getTwilioAccountSID(),
	Password: getTwilioAuthToken(),
})

func (o *otpTwilioService) TwilioSendOTP(phoneNumber string) (string, error) {
	params := &twilioApi.CreateVerificationParams{}
	params.SetTo(phoneNumber)
	params.SetChannel("sms")

	resp, err := client.VerifyV2.CreateVerification(o.twilioServiceSID, params)
	if err != nil {
		return "", err
	}

	return *resp.Sid, nil
}

func (o *otpTwilioService) TwilioVerifyOTP(phoneNumber string, code string) error {
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
