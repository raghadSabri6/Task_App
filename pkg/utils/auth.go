package utils

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// UserUUIDKey is the context key for the user UUID
type contextKey string
const UserUUIDKey contextKey = "userUUID"

// GetUserUUIDFromRequest gets the user UUID from the request context
func GetUserUUIDFromRequest(r *http.Request) uuid.UUID {
	// Get user UUID from context
	userUUID, ok := r.Context().Value(UserUUIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil
	}
	
	return userUUID
}

// GetUserUUIDFromContext gets the user UUID from the context
func GetUserUUIDFromContext(ctx context.Context) uuid.UUID {
	// Get user UUID from context
	userUUID, ok := ctx.Value(UserUUIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil
	}
	
	return userUUID
}