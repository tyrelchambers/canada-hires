package services

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/log"
	"gopkg.in/gomail.v2"
)

type EmailService interface {
	SendLoginLink(email, token string) error
}

type emailService struct {
	smtpHost     string
	smtpPort     int
	smtpUser     string
	smtpPassword string
	fromEmail    string
	backendURL   string
}

func NewEmailService() EmailService {
	port, _ := strconv.Atoi(getEnv("SMTP_PORT", "587"))

	return &emailService{
		smtpHost:     getEnv("SMTP_HOST", "localhost"),
		smtpPort:     port,
		smtpUser:     getEnv("SMTP_USER", ""),
		smtpPassword: getEnv("SMTP_PASSWORD", ""),
		fromEmail:    getEnv("FROM_EMAIL", "noreply@canada-hires.com"),
		backendURL:   getEnv("API_URL", "http://localhost:8000"),
	}
}

func (s *emailService) SendLoginLink(email, token string) error {
	loginURL := fmt.Sprintf("%s/api/auth/verify-login/%s", s.backendURL, token)

	m := gomail.NewMessage()
	m.SetHeader("From", s.fromEmail)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Login to Canada Hires")

	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>Login to Canada Hires</h2>
			<p>Click the link below to log in to your account:</p>
			<p><a href="%s">Login to Canada Hires</a></p>
			<p>This link will expire in 15 minutes.</p>
			<p>If you didn't request this login, please ignore this email.</p>
		</body>
		</html>
	`, loginURL)

	log.Info(loginURL)

	m.SetBody("text/html", body)

	d := gomail.NewDialer(s.smtpHost, s.smtpPort, s.smtpUser, s.smtpPassword)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
