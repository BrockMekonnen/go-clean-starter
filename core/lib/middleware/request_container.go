package middleware

import (
	"context"
	"net/http"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/extension"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/logger"
	"github.com/gorilla/mux"
	"go.uber.org/dig"
)

func RequestContainerMiddleware(rootContainer *dig.Container, logger logger.Log) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create named request scope (name is required)
			reqScope := rootContainer.Scope("http-request")

			// Register request-specific values
			err := reqScope.Provide(func() string {
				if id := r.Header.Get("X-Request-ID"); id != "" {
					return id
				}
				return generateRequestID()
			}, dig.Name("requestId"))

			if err != nil {
				logger.Error("failed to register request ID")
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Register HTTP objects in a separate function for clarity
			if err := registerHTTPObjects(reqScope, r, w, logger); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Store scope in context
			ctx := context.WithValue(r.Context(), extension.ContainerContextKey, reqScope)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func registerHTTPObjects(scope *dig.Scope, r *http.Request, w http.ResponseWriter, logger logger.Log) error {
	// Register *http.Request
	if err := scope.Provide(func() *http.Request { return r }); err != nil {
		logger.Error("failed to register http request")
		return err
	}

	// Register http.ResponseWriter
	if err := scope.Provide(func() http.ResponseWriter { return w }); err != nil {
		logger.Error("failed to register http response writer")
		return err
	}

	return nil
}

// GetContainerFromRequest retrieves the request scope
func GetContainerFromRequest(r *http.Request) *dig.Scope {
	if container, ok := r.Context().Value(extension.ContainerContextKey).(*dig.Scope); ok {
		return container
	}
	return nil
}

func generateRequestID() string {
	// Implement your preferred request ID generation
	// Example: return uuid.New().String()
	return ""
}
