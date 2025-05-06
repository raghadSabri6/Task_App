package helperFunc

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateJWTToken(userID uuid.UUID, expiryDuration time.Duration) (string, error) {
	secret := os.Getenv("SECRET")
	if secret == "" {
		log.Fatal("SECRET key not set in environment variables")
	}

	now := time.Now().Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID.String(),
		"iat": now,
		"nbf": now,
		"exp": now + int64(expiryDuration.Seconds()),
	})

	return token.SignedString([]byte(secret))
}

func SetJWTTokenCookie(w http.ResponseWriter, token string, expiryDuration time.Duration) {
	http.SetCookie(w, &http.Cookie{
		Name:     "Authorization",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(expiryDuration),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
}

func GetTokenExpiration() time.Duration {
	expiryStr := os.Getenv("TOKEN_EXPIRATION")
	expirySeconds, err := strconv.Atoi(expiryStr)
	if err != nil || expirySeconds <= 0 {
		expirySeconds = 3600 * 24 * 30 // 30 days
	}
	return time.Duration(expirySeconds) * time.Second
}

func GetUserUUIDFromRequest(r *http.Request) uuid.UUID {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			userID := validateToken(tokenStr)
			if userID != uuid.Nil {
				return userID
			}
		}
	}

	cookie, err := r.Cookie("Authorization")
	if err != nil {
		return uuid.Nil
	}

	tokenStr := strings.TrimSpace(cookie.Value)
	if tokenStr == "" {
		return uuid.Nil
	}

	return validateToken(tokenStr)
}

func validateToken(tokenStr string) uuid.UUID {
	secret := os.Getenv("SECRET")
	if secret == "" {
		return uuid.Nil
	}

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return uuid.Nil
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if sub, ok := claims["sub"].(string); ok {
			userID, err := uuid.Parse(sub)
			if err == nil {
				return userID
			}
		}
	}

	return uuid.Nil
}
