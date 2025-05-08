package middleware

import (
	"net/http"
	"task2/pkg/utils"
)

// MethodCheck checks if the request method is allowed
func MethodCheck(methods ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if method is allowed
			allowed := false
			for _, method := range methods {
				if r.Method == method {
					allowed = true
					break
				}
			}
			
			if !allowed {
				utils.RespondJSON(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
				return
			}
			
			// Call next handler
			next.ServeHTTP(w, r)
		})
	}
}