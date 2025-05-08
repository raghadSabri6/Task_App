package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"
	"task2/internal/infrastructure/auth"
	"task2/pkg/utils"
)

// AuthMiddleware is a middleware for authentication
type AuthMiddleware struct {
	authService *auth.AuthService
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(authService *auth.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// Middleware authenticates the request
func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to get token from different sources
		var tokenString string
		
		// First try Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString = parts[1]
				log.Printf("Auth: Found token in Authorization header")
			}
		}
		
		// If no token in header, try cookie
		if tokenString == "" {
			cookie, err := r.Cookie("Authorization")
			if err == nil && cookie.Value != "" {
				tokenString = cookie.Value
				log.Printf("Auth: Found token in cookie")
			}
		}
		
		// If still no token, return error
		if tokenString == "" {
			log.Printf("Auth: No token found in request")
			utils.RespondJSON(w, http.StatusUnauthorized, "Authentication required", nil)
			return
		}
		
		// Validate token
		userUUID, err := m.authService.ValidateToken(tokenString)
		if err != nil {
			log.Printf("Auth: Invalid token: %v", err)
			utils.RespondJSON(w, http.StatusUnauthorized, "Invalid token", nil)
			return
		}
		
		log.Printf("Auth: Token validated successfully for user %s", userUUID)
		
		// Add user UUID to context
		ctx := context.WithValue(r.Context(), utils.UserUUIDKey, userUUID)
		
		// Call next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}