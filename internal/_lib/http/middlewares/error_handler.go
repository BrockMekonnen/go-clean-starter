package middleware

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Exception represents a base error interface
type Exception interface {
	error
	StatusCode() int
	Body() interface{}
}

// ErrorConverter defines a function that tests and converts errors
type ErrorConverter func(err error) (status int, body interface{}, ok bool)

// ErrorHandlerConfig contains configuration for the error handler
type ErrorHandlerConfig struct {
	Logger *log.Logger
}

// DefaultErrorHandlerConfig provides default configuration
var DefaultErrorHandlerConfig = ErrorHandlerConfig{
	Logger: log.Default(),
}

// ErrorConverter creates a new error converter
func NewErrorConverter[E Exception](
	test func(err error) bool,
	converter func(err E) (int, interface{}),
) ErrorConverter {
	return func(err error) (int, interface{}, bool) {
		if !test(err) {
			return 0, nil, false
		}
		// Type assert the error to our specific type
		if e, ok := err.(E); ok {
			status, body := converter(e)
			return status, body, true
		}
		return 0, nil, false
	}
}

// ErrorHandler creates an Echo middleware error handler
func ErrorHandler(converters []ErrorConverter, config ...ErrorHandlerConfig) echo.HTTPErrorHandler {
	cfg := DefaultErrorHandlerConfig
	if len(config) > 0 {
		cfg = config[0]
	}

	return func(err error, c echo.Context) {
		// Log the error stack
		cfg.Logger.Printf("Error: %v\nStack: %+v\n", err, err)

		// Try all converters
		for _, converter := range converters {
			if status, body, ok := converter(err); ok {
				_ = c.JSON(status, body)
				return
			}
		}

		// Default error response
		if httpErr, ok := err.(*echo.HTTPError); ok {
			_ = c.JSON(httpErr.Code, map[string]interface{}{"error": httpErr.Message})
			return
		}

		_ = c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}
}