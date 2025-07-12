package utils

import (
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func SendEmail(to, subject, body string) error {
	godotenv.Load()
	smtpHost := "smtp.gmail.com"
	smtpPort := 587

	sender := os.Getenv("EMAIL_SENDER")
	password := os.Getenv("EMAIL_PASSWORD")

	msg := gomail.NewMessage()
	msg.SetHeader("From", sender)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	dialer := gomail.NewDialer(smtpHost, smtpPort, sender, password)
	return dialer.DialAndSend(msg)
}
