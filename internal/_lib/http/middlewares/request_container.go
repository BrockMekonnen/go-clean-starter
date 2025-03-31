package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
)

// RequestContainer creates a scoped DI container for each request
func RequestContainer(rootContainer *dig.Container, logger *logrus.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Create request-scoped container
			reqContainer := dig.New()
			
			// Clone all providers from root container
			if err := reqContainer.Clone(rootContainer); err != nil {
				logger.WithError(err).Error("Failed to clone container")
				return echo.NewHTTPError(500, "Internal Server Error")
			}

			// Register request-specific values
			err := reqContainer.Provide(func() string {
				return c.Response().Header().Get(echo.HeaderXRequestID)
			}, dig.Name("requestId"))
			
			if err != nil {
				logger.WithError(err).Error("Failed to register request values")
				return echo.NewHTTPError(500, "Internal Server Error")
			}

			// Store container in context
			c.Set("container", reqContainer)

			return next(c)
		}
	}
}