package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// AuthService implements the authentication service
type AuthService struct {
	jwtSecret []byte
}

// NewAuthService creates a new authentication service
func NewAuthService(jwtSecret string) *AuthService {
	return &AuthService{
		jwtSecret: []byte(jwtSecret),
	}
}

// GenerateToken generates a JWT token for a user
func (s *AuthService) GenerateToken(userUUID uuid.UUID) (string, error) {
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userUUID.String(),
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
	})
	
	// Sign token
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}
	
	return tokenString, nil
}

// ValidateToken validates a JWT token
func (s *AuthService) ValidateToken(tokenString string) (uuid.UUID, error) {
	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		
		return s.jwtSecret, nil
	})
	
	if err != nil {
		return uuid.Nil, err
	}
	
	// Validate claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Get user UUID
		userUUIDStr, ok := claims["sub"].(string)
		if !ok {
			return uuid.Nil, errors.New("invalid token claims")
		}
		
		// Parse UUID
		userUUID, err := uuid.Parse(userUUIDStr)
		if err != nil {
			return uuid.Nil, err
		}
		
		return userUUID, nil
	}
	
	return uuid.Nil, errors.New("invalid token")
}

// HashPassword hashes a password
func (s *AuthService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	
	return string(hashedPassword), nil
}

// ValidatePassword validates a password
func (s *AuthService) ValidatePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}