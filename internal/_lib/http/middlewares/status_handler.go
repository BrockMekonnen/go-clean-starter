package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// Dependencies contains the required dependencies for the status handler
type Dependencies struct {
	StartedAt time.Time
}

// StatusResponse defines the JSON response structure
type StatusResponse struct {
	StartedAt time.Time `json:"startedAt"`
	Uptime    float64   `json:"uptime"`
}

// StatusHandler creates a handler function for status checks
func StatusHandler(deps Dependencies) echo.HandlerFunc {
	return func(c echo.Context) error {
		uptime := time.Since(deps.StartedAt).Seconds()

		return c.JSON(http.StatusOK, StatusResponse{
			StartedAt: deps.StartedAt,
			Uptime:    uptime,
		})
	}
}