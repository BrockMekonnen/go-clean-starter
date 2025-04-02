package middleware

import (
	"encoding/json"
	"fmt"

	"net/http"
	"time"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/errors"
)

// StatusResponse defines the JSON response structure
type StatusResponse struct {
	StartedAt time.Time `json:"startedAt"`
	Uptime    float64   `json:"uptime"`
}

// StatusHandler creates a handler function for status checks
func StatusHandler(startedAt time.Time) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		errorType := r.URL.Query().Get("error")

		switch errorType {
		case "bad-request":
			// This will be caught by the BadRequestError converter
			panic(errors.NewBadRequestError[any]("", "", nil))

		case "not-found":
			// This will be caught by the NotFoundError converter
			panic(errors.NewNotFoundError(
				"Resource not found",
				"RESOURCE_NOT_FOUND",
				map[string]interface{}{"resource": "user"},
			))

		case "validation":
			// This will be caught by the ValidationError converter
			panic(errors.NewValidationError("body", "invalid email", nil))

		}

		// Normal response
		response := StatusResponse{
			StartedAt: startedAt,
			Uptime:    time.Since(startedAt).Seconds(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			fmt.Print("status_handler_error: json encode failed")
		}

	}
}
