package middlewares

import (
	"net/http"
	"task2/helperFunc"
)

func MethodCheck(method string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			helperFunc.RespondJSON(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
			return
		}
		next(w, r)
	}
}
