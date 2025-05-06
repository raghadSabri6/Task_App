package middlewares

import (
	"net/http"
	"task2/helperFunc"

	"github.com/google/uuid"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userUUID := helperFunc.GetUserUUIDFromRequest(r)
		if userUUID == uuid.Nil {
			helperFunc.RespondJSON(w, http.StatusUnauthorized, "User not authenticated", nil)
			return
		}
		next(w, r)
	}
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return Auth(next)
}
