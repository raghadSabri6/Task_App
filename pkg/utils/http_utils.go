package utils

import (
	"bytes"
	"io"
	"net/http"
	"strings"
)

// ReadRequestBody reads the request body and returns it as a string
func ReadRequestBody(r *http.Request) (string, error) {
	// Read body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	
	// Restore body
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	
	return string(bodyBytes), nil
}

// GetTokenFromRequest gets the token from the request
func GetTokenFromRequest(r *http.Request) string {
	// Try to get from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		// Check if it's a Bearer token
		if strings.HasPrefix(authHeader, "Bearer ") {
			return strings.TrimPrefix(authHeader, "Bearer ")
		}
		return authHeader
	}
	
	// Try to get from cookie
	cookie, err := r.Cookie("Authorization")
	if err == nil {
		return cookie.Value
	}
	
	return ""
}