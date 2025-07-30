package services

import (
	"canada-hires/utils"
	"fmt"
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
	port, _ := strconv.Atoi(utils.GetEnv("SMTP_PORT", "587"))

	return &emailService{
		smtpHost:     utils.GetEnv("SMTP_HOST", "localhost"),
		smtpPort:     port,
		smtpUser:     utils.GetEnv("SMTP_USER", ""),
		smtpPassword: utils.GetEnv("SMTP_PASSWORD", ""),
		fromEmail:    utils.GetEnv("FROM_EMAIL", "noreply@jobwatchcanada.com"),
		backendURL:   utils.GetEnv("API_URL", "http://localhost:8000"),
	}
}

func (s *emailService) SendLoginLink(email, token string) error {
	loginURL := fmt.Sprintf("%s/api/auth/verify-login/%s", s.backendURL, token)

	m := gomail.NewMessage()
	m.SetHeader("From", s.fromEmail)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Login to JobWatch Canada")

	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>Login to JobWatch Canada</h2>
			<p>Click the link below to log in to your account:</p>
			<p><a href="%s">Login to JobWatch Canada</a></p>
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

