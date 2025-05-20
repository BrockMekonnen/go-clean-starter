package respond

import (
	"encoding/json"

	"github.com/BrockMekonnen/go-clean-starter/internal/_shared/delivery"
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
	JSON(w, http.StatusOK, data)
}

func SuccessWithData(w http.ResponseWriter, status int, data interface{}) {
	JSON(w, status, map[string]interface{}{
		"data": data,
	})
}

var converters []delivery.ErrorConverter

// RegisterConverters lets you register converters from the middleware
func RegisterConverters(c []delivery.ErrorConverter) {
	converters = c
}

// Error formats error responses, using registered converters if available
func Error(w http.ResponseWriter, err error) {
	for _, converter := range converters {
		if converter.Test(err) {
			status, body := converter.Convert(err)
			JSON(w, status, body)
			return
		}
	}

	// fallback if no converter matched
	JSON(w, http.StatusInternalServerError, map[string]interface{}{
		"error":   "InternalServerError",
		"message": err.Error(),
	})
}
