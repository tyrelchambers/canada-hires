package controllers

import (
	"canada-hires/dto"
	"canada-hires/services"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type AuthController interface {
	SendLoginLink(w http.ResponseWriter, r *http.Request)
	VerifyLogin(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

type authController struct {
	authService services.AuthService
	userService services.UserService
}

func NewAuthController(authService services.AuthService, userService services.UserService) AuthController {
	return &authController{authService: authService, userService: userService}
}

func (c *authController) SendLoginLink(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	clientIP := getClientIP(r)

	if err := c.authService.SendLoginLink(req.Email, clientIP); err != nil {
		log.Error("Failed to send login link", "error", err, "email", req.Email)
		http.Error(w, "Failed to send login link", http.StatusInternalServerError)
		return
	}

	response := dto.LoginResponse{
		Message: "Login link sent to your email",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (c *authController) VerifyLogin(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	clientIP := getClientIP(r)
	userAgent := r.Header.Get("User-Agent")

	session, _, err := c.authService.VerifyLoginToken(token, clientIP, userAgent)
	if err != nil {
		log.Error("Failed to verify login token", "error", err, "token", token)
		// Redirect to a failure page or show an error
		http.Redirect(w, r, "/auth/login?error=verification_failed", http.StatusSeeOther)
		return
	}

	// Set session ID as HTTP-only cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		HttpOnly: true,
		Secure:   true, // Set to true in production
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Expires:  session.ExpiresAt,
	})

	// Redirect to the frontend home page
	http.Redirect(w, r, os.Getenv("FRONTEND_URL"), http.StatusSeeOther)
}

func (c *authController) Logout(w http.ResponseWriter, r *http.Request) {
	// Get session from cookie to delete it from database
	cookie, err := r.Cookie("session_id")
	if err == nil {
		if sessionID, parseErr := uuid.Parse(cookie.Value); parseErr == nil {
			// Delete session from database
			if deleteErr := c.authService.DeleteSession(sessionID.String()); deleteErr != nil {
				log.Error("Failed to delete session", "error", deleteErr, "session_id", sessionID)
			}
		}
	}

	// Clear the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   -1, // Delete cookie
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logged out successfully",
	})
}

func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		ips := strings.Split(xForwardedFor, ",")
		if len(ips) > 0 {
			ip := strings.TrimSpace(ips[0])
			if isValidIP(ip) {
				return ip
			}
		}
	}

	// Check X-Real-IP header
	xRealIP := r.Header.Get("X-Real-IP")
	if xRealIP != "" && isValidIP(xRealIP) {
		return xRealIP
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	if isValidIP(ip) {
		return ip
	}

	return "127.0.0.1"
}

func isValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}
