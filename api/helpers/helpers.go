package helpers

import (
	"canada-hires/models"
	"context"
	"net"
	"net/http"
	"os"
	"strings"
)

// Key type for context values
type contextKey string

// UserContextKey is the key used to store and retrieve the user from context
const ContextKey contextKey = "user"

var UserContextKey = contextKey("user")

// GetUserFromContext retrieves the user from the request context
// Returns nil if no user is found or context value is not a valid user
func GetUserFromContext(ctx context.Context) *models.User {
	value := ctx.Value(ContextKey)
	if value == nil {
		return nil
	}

	user, ok := value.(*models.User)
	if !ok {
		return nil
	}

	return user
}

func IsDev() bool {
	if os.Getenv("ENV") == "development" {
		return true
	}

	return false
}

// GetClientIP extracts the real client IP address from request headers
func GetClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		ips := strings.Split(xForwardedFor, ",")
		if len(ips) > 0 {
			ip := strings.TrimSpace(ips[0])
			if IsValidIP(ip) {
				return ip
			}
		}
	}

	// Check X-Real-IP header
	xRealIP := r.Header.Get("X-Real-IP")
	if xRealIP != "" && IsValidIP(xRealIP) {
		return xRealIP
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	if IsValidIP(ip) {
		return ip
	}

	return "127.0.0.1"
}

// IsValidIP validates if a string is a valid IP address
func IsValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}
