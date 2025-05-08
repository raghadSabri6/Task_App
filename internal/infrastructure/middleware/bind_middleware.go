package middleware

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"task2/internal/infrastructure/validator"
	"task2/pkg/utils"
)

// BindKey is the context key for the bound request
type contextKey string

const BindKey contextKey = "bind"

// BindAndValidate binds and validates the request body
func BindAndValidate(v interface{}) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check content type
			contentType := r.Header.Get("Content-Type")
			if contentType != "application/json" {
				utils.RespondJSON(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json", nil)
				return
			}

			// Read body
			body, err := io.ReadAll(r.Body)
			if err != nil {
				utils.RespondJSON(w, http.StatusBadRequest, "Failed to read request body", nil)
				return
			}
			defer r.Body.Close()

			// Parse JSON
			err = json.Unmarshal(body, v)
			if err != nil {
				utils.RespondJSON(w, http.StatusBadRequest, "Invalid JSON", nil)
				return
			}

			// Validate
			if err := validator.Validate(v); err != nil {
				utils.RespondJSON(w, http.StatusBadRequest, err.Error(), nil)
				return
			}

			// Add to context
			ctx := context.WithValue(r.Context(), BindKey, v)

			// Call next handler
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
