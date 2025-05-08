package middleware

import (
	"log"
	"net/http"
)

// CorsMiddleware adds CORS headers to allow cross-origin requests
type CorsMiddleware struct {
	logger *log.Logger
}

// NewCorsMiddleware creates a new CORS middleware
func NewCorsMiddleware(logger *log.Logger) *CorsMiddleware {
	return &CorsMiddleware{
		logger: logger,
	}
}

// Middleware returns a middleware function that adds CORS headers
func (m *CorsMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get origin from request
		origin := r.Header.Get("Origin")
		if origin != "" {
			// Set CORS headers
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true") // Important for cookies
			w.Header().Set("Access-Control-Max-Age", "3600")
			
			m.logger.Printf("CORS headers set for origin: %s", origin)
		}

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			m.logger.Printf("Handled OPTIONS preflight request")
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}