package services

import (
	"canada-hires/models"
	"canada-hires/repos"
	"canada-hires/utils"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	SendLoginLink(email, ipAddress string) error
	VerifyLoginToken(token, ipAddress, userAgent string) (*models.Session, *models.User, error)
	CreateSession(userID string, ipAddress, userAgent string) (*models.Session, error)
	ValidateSession(sessionID string) (*models.Session, error)
	DeleteSession(sessionID string) error
	DeleteUserSessions(userID string) error
	GenerateJWT(userID string) (string, error)
	ValidateJWT(tokenString string) (*jwt.Token, error)
}

type authService struct {
	userRepo        repos.UserRepository
	tokenRepo       repos.LoginTokenRepository
	sessionRepo     repos.SessionRepository
	emailService    EmailService
	jwtSecret       []byte
	tokenDuration   time.Duration
	sessionDuration time.Duration
}

func NewAuthService(userRepo repos.UserRepository, tokenRepo repos.LoginTokenRepository, sessionRepo repos.SessionRepository, emailService EmailService) AuthService {
	return &authService{
		userRepo:        userRepo,
		tokenRepo:       tokenRepo,
		sessionRepo:     sessionRepo,
		emailService:    emailService,
		jwtSecret:       []byte(utils.GetEnv("JWT_SECRET", "your-secret-key")),
		tokenDuration:   15 * time.Minute,
		sessionDuration: 30 * 24 * time.Hour, // 30 days
	}
}

func (s *authService) SendLoginLink(email, ipAddress string) error {
	// Check if email is in approved list
	if !s.isEmailApproved(email) {
		return fmt.Errorf("email not approved for registration")
	}

	// Get or create user
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		if err == sql.ErrNoRows {
			// Create new user
			user = &models.User{
				Email: email,
			}
			if err := s.userRepo.Create(user); err != nil {
				return fmt.Errorf("failed to create user: %w", err)
			}
		} else {
			return fmt.Errorf("failed to get user: %w", err)
		}
	}

	// Generate secure token
	token, err := s.generateSecureToken()
	if err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}

	// Create login token
	loginToken := &models.LoginToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(s.tokenDuration),
		IPAddress: &ipAddress,
	}

	if err := s.tokenRepo.Create(loginToken); err != nil {
		return fmt.Errorf("failed to create login token: %w", err)
	}

	// Send email
	if err := s.emailService.SendLoginLink(email, token); err != nil {
		return fmt.Errorf("failed to send login email: %w", err)
	}

	// Add IP to user's IP list
	if err := s.userRepo.AddIPAddress(user.ID, ipAddress); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: failed to add IP address: %v", err)
	}

	return nil
}

func (s *authService) VerifyLoginToken(token, ipAddress, userAgent string) (*models.Session, *models.User, error) {
	// Get token from database
	loginToken, err := s.tokenRepo.GetByToken(token)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, fmt.Errorf("invalid or expired token")
		}
		return nil, nil, fmt.Errorf("failed to get token: %w", err)
	}

	// Check if token is expired
	if loginToken.ExpiresAt.Before(time.Now()) {
		return nil, nil, fmt.Errorf("token has expired")
	}

	// Check if token is already used
	if loginToken.UsedAt != nil {
		return nil, nil, fmt.Errorf("token has already been used")
	}

	// Get user
	user, err := s.userRepo.GetByID(loginToken.UserID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Mark token as used
	if err := s.tokenRepo.MarkAsUsed(loginToken.ID); err != nil {
		return nil, nil, fmt.Errorf("failed to mark token as used: %w", err)
	}

	// Update user's last active time
	if err := s.userRepo.UpdateLastActive(user.ID); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: failed to update last active: %v", err)
	}

	// Add IP to user's IP list
	if err := s.userRepo.AddIPAddress(user.ID, ipAddress); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: failed to add IP address: %v", err)
	}

	// Create session
	session, err := s.CreateSession(user.ID, ipAddress, userAgent)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, user, nil
}

func (s *authService) CreateSession(userID string, ipAddress, userAgent string) (*models.Session, error) {
	session := &models.Session{
		UserID:    userID,
		ExpiresAt: time.Now().UTC().Add(s.sessionDuration),
		IPAddress: &ipAddress,
		UserAgent: &userAgent,
	}

	if err := s.sessionRepo.Create(session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

func (s *authService) ValidateSession(sessionID string) (*models.Session, error) {
	session, err := s.sessionRepo.GetByID(sessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	if session.IsExpired() {
		// Delete expired session
		s.sessionRepo.DeleteByID(sessionID)
		return nil, fmt.Errorf("session has expired")
	}

	// Update session last used
	if err := s.sessionRepo.UpdateLastUsed(sessionID); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: failed to update session last used: %v", err)
	}

	return session, nil
}

func (s *authService) DeleteSession(sessionID string) error {
	return s.sessionRepo.DeleteByID(sessionID)
}

func (s *authService) DeleteUserSessions(userID string) error {
	return s.sessionRepo.DeleteByUserID(userID)
}

func (s *authService) GenerateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *authService) ValidateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})
}

func (s *authService) generateSecureToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func getClientIP(r interface{}) string {
	// This is a placeholder - in a real implementation you'd extract IP from http.Request
	// For now, return a default IP
	return "127.0.0.1"
}

func isValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

// isEmailApproved checks if the email is in the APPROVED_EMAILS environment variable
func (s *authService) isEmailApproved(email string) bool {
	approvedEmails := utils.GetEnv("APPROVED_EMAILS", "")
	if approvedEmails == "" {
		// If no approved emails are set, allow all emails (backward compatibility)
		return true
	}

	// Split by comma and check each email
	emails := strings.Split(approvedEmails, ",")
	for _, approvedEmail := range emails {
		if strings.TrimSpace(approvedEmail) == strings.TrimSpace(email) {
			return true
		}
	}

	return false
}

