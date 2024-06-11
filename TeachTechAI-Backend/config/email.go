package config

import "teach-tech-ai/helpers"

type EmailConfig struct {
	Host         string `mapstructure:"SMTP_HOST"`
	Port         int    `mapstructure:"SMTP_PORT"`
	SenderName   string `mapstructure:"SMTP_SENDER_NAME"`
	AuthEmail    string `mapstructure:"SMTP_AUTH_EMAIL"`
	AuthPassword string `mapstructure:"SMTP_AUTH_PASSWORD"`
}

func NewEmailConfig() (*EmailConfig, error) {
	var config EmailConfig
	config.Host = helpers.MustGetenv("SMTP_HOST")
	config.Port = helpers.MustGetenvInt("SMTP_PORT")
	config.SenderName = helpers.MustGetenv("SMTP_SENDER_NAME")
	config.AuthEmail = helpers.MustGetenv("SMTP_AUTH_EMAIL")
	config.AuthPassword = helpers.MustGetenv("SMTP_AUTH_PASSWORD")

	return &config, nil
}