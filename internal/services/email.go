package services

import (
	"fmt"

	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/config"
	"gopkg.in/gomail.v2"
)

func SendEmail(email, username, verificationCode string) error {
	var credentials = map[string]interface{}{
		"from":     config.Config.FROM_EMAIL,
		"username": config.Config.SMTP_USERNAME,
		"password": config.Config.SMTP_PASSWORD,
		"host":     config.Config.SMTP_HOST,
		"port":     config.Config.SMTP_PORT,
	}

	for key, value := range credentials {
		if value == "" {
			return fmt.Errorf("email credential %v is not set", key)
		}
		fmt.Println(value)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", credentials["from"].(string))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "🔐 Your Verification Code for CodePulse")
	m.SetBody("text/html", fmt.Sprintf("Hello %s! This is your verification code for CodePulse: <b>%s</b>", username, verificationCode))
	dialer := gomail.NewDialer(credentials["host"].(string), credentials["port"].(int), credentials["username"].(string), credentials["password"].(string))

	if err := dialer.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
