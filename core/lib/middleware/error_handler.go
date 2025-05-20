package middleware

import (
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/logger"
	"github.com/BrockMekonnen/go-clean-starter/internal/_shared/delivery"
	"github.com/gorilla/mux"
)

// ErrorConverter defines the interface for converting errors to HTTP responses
type ErrorConverter interface {
	Test(err error) bool
	Convert(err error) (status int, body interface{})
}

// errorConverter implements the ErrorConverter interface
type errorConverter struct {
	test    func(error) bool
	convert func(error) (int, interface{})
}

func (ec *errorConverter) Test(err error) bool {
	return ec.test(err)
}

func (ec *errorConverter) Convert(err error) (int, interface{}) {
	return ec.convert(err)
}

// ErrorConverterFn is a function type for creating error converters
type ErrorConverterFn func(test func(error) bool, convert func(error) (int, interface{})) ErrorConverter

// NewErrorConverter creates a new ErrorConverter
func NewErrorConverter(test func(error) bool, convert func(error) (int, interface{})) ErrorConverter {
	return &errorConverter{
		test:    test,
		convert: convert,
	}
}

// ErrorHandlerOptions configures the error handler middleware
type ErrorHandlerOptions struct {
	Logger logger.Interface
}

// DefaultErrorHandlerOptions returns default options
func DefaultErrorHandlerOptions() ErrorHandlerOptions {
	return ErrorHandlerOptions{
		Logger: logger.NewLogger(), // Using your logger interface
	}
}

// ErrorHandler creates middleware for handling errors
func ErrorHandler(converters []delivery.ErrorConverter, logger logger.Interface) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// Convert panic to error
					var e error
					switch v := err.(type) {
					case error:
						e = v
					default:
						e = fmt.Errorf("%v", v)
					}

					// Find matching converter
					for _, converter := range converters {
						if converter.Test(e) {
							status, body := converter.Convert(e)
							w.Header().Set("Content-Type", "application/json")
							w.WriteHeader(status)
							err := json.NewEncoder(w).Encode(body)
							if err != nil {
								logger.Error("ErrorHandler: json encode failed", err)
							}
							return
						}
					}

					// Default error response
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					err = json.NewEncoder(w).Encode(map[string]interface{}{
						"error":   "InternalServerError",
						"message": e.Error(),
					})
					if err != nil {
						logger.Error("ErrorHandler: json encode failed", err)
					}
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
