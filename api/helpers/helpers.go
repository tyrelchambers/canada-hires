package helpers

import (
	"canada-hires/models"
	"context"
	"os"
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
