package middleware

import (
	"canada-hires/helpers"
	"canada-hires/services"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

// AuthMiddleware creates a middleware that extracts the user from the cookie
// and attaches it to the request context

type AuthMiddleware struct {
	authService services.AuthService
	userService services.UserService
}

func NewAuthMiddleware(authService services.AuthService, userService services.UserService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		userService: userService,
	}
}

func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the cookie
		cookie, err := r.Cookie("session_id")
		if err != nil {
			// No cookie found, continue without user in context
			log.Debug("No auth cookie found")
			next.ServeHTTP(w, r)
			return
		}
		// Validate the token and get the user
		session, err := m.authService.ValidateSession(cookie.Value)

		if err != nil || session == nil {
			next.ServeHTTP(w, r)
			return
		}

		user, err := m.userService.GetUserByID(session.UserID)

		if err != nil || user == nil {
			next.ServeHTTP(w, r)
			return
		}

		// Attach the user to the request context
		ctx := context.WithValue(r.Context(), helpers.ContextKey, user)

		// Call the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

// RequireAuth creates a middleware that requires authentication
// It will return a 401 Unauthorized response if the user is not authenticated
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := helpers.GetUserFromContext(r.Context())
		if user == nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"Unauthorized"}`))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := generateRequestID()
		r.Header.Set("X-Request-ID", requestID)
		w.Header().Set("X-Request-ID", requestID)

		// Add request ID to context for logging
		ctx := context.WithValue(r.Context(), "requestID", requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func generateRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// CSRFStateMiddleware provides CSRF protection for OAuth flows using state parameters
type CSRFStateMiddleware struct {
	cookieName string
	expiry     time.Duration
}

// NewCSRFStateMiddleware creates a new CSRF state middleware
func NewCSRFStateMiddleware(cookieName string, expiry time.Duration) *CSRFStateMiddleware {
	return &CSRFStateMiddleware{
		cookieName: cookieName,
		expiry:     expiry,
	}
}

// GenerateAndSetState generates a CSRF state token and sets it as a cookie
func (c *CSRFStateMiddleware) GenerateAndSetState(w http.ResponseWriter) (string, error) {
	state, err := generateCSRFState()
	if err != nil {
		return "", fmt.Errorf("failed to generate CSRF state: %w", err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     c.cookieName,
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   helpers.IsDev() == false, // Only secure in production
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(c.expiry),
	})

	return state, nil
}

// ValidateState validates the CSRF state parameter against the cookie
func (c *CSRFStateMiddleware) ValidateState(r *http.Request, state string) error {
	if state == "" {
		return fmt.Errorf("missing state parameter")
	}

	cookie, err := r.Cookie(c.cookieName)
	if err != nil {
		return fmt.Errorf("missing CSRF state cookie: %w", err)
	}

	if cookie.Value == "" {
		return fmt.Errorf("empty CSRF state cookie")
	}

	if state != cookie.Value {
		return fmt.Errorf("CSRF state mismatch")
	}

	return nil
}

// ClearState removes the CSRF state cookie
func (c *CSRFStateMiddleware) ClearState(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     c.cookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(-1 * time.Hour),
	})
}

// Middleware provides CSRF protection for specific routes
func (c *CSRFStateMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only validate state for callback endpoints
		if strings.Contains(r.URL.Path, "/callback") {
			state := r.URL.Query().Get("state")
			if err := c.ValidateState(r, state); err != nil {
				log.Error(err)
				return
			}

			// Clear the state cookie after successful validation
			c.ClearState(w)
		}

		next.ServeHTTP(w, r)
	})
}

// generateCSRFState creates a cryptographically secure random state parameter
func generateCSRFState() (string, error) {
	bytes := make([]byte, 32) // 256 bits
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to read random bytes: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
