package respond

import (
	"encoding/json"
	"net/http"
)

// JSON sends a standardized JSON response
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data) // Error handling omitted for brevity
}

// Success formats successful responses
func Success(w http.ResponseWriter, status int, data interface{}) {
	JSON(w, http.StatusOK, map[string]interface{}{
		"data": data,
	})
}

// Error formats error responses
func Error(w http.ResponseWriter, status int, message string, details interface{}) {
	JSON(w, status, map[string]interface{}{
		"error":   message,
		"details": details,
	})
}
