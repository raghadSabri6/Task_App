package utils

import (
	"encoding/json"
	"net/http"
)

// Response represents an HTTP response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// RespondJSON sends a JSON response
func RespondJSON(w http.ResponseWriter, status int, message string, data interface{}) {
	// Set content type
	w.Header().Set("Content-Type", "application/json")
	
	// Set status code
	w.WriteHeader(status)
	
	// Create response
	response := Response{
		Success: status >= 200 && status < 300,
		Message: message,
		Data:    data,
	}
	
	// Encode response
	json.NewEncoder(w).Encode(response)
}