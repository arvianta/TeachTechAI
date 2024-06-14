package utils

import (
	"bytes"
	"html/template"
	"os"
	"teach-tech-ai/config"

	"gopkg.in/gomail.v2"
)

func SendMail(toEmail string, subject string, body string) error {
	emailConfig, err := config.NewEmailConfig()
	if err != nil {
		return err
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", emailConfig.AuthEmail)
	mailer.SetHeader("To", toEmail)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body)

	dialer := gomail.NewDialer(
		emailConfig.Host,
		emailConfig.Port,
		emailConfig.AuthEmail,
		emailConfig.AuthPassword,
	)

	err = dialer.DialAndSend(mailer)
	if err != nil {
		return err
	}

	return nil
}

func BuildMail(receiverEmail string, receiverName string, otp string) (map[string]string, error) {
	readHtml, err := os.ReadFile("utils/email-template/otpEmail.html")
	if err != nil {
		return nil, err
	}

	data := struct {
		Name  string
		OTP   string
		Email string
	}{
		Name:  receiverName,
		OTP:   otp,
		Email: receiverEmail,
	}

	tmpl, err := template.New("custom").Parse(string(readHtml))
	if err != nil {
		return nil, err
	}

	var strMail bytes.Buffer
	if err := tmpl.Execute(&strMail, data); err != nil {
		return nil, err
	}

	draftEmail := map[string]string{
		"subject": "Verifikasi OTP TeachTechAI",
		"body":    strMail.String(),
	}

	return draftEmail, nil
}
