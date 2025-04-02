package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/logger"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type startTimeKey string
type idKey string

const (
	reqStartTimeKey startTimeKey = "reqStartTime"
	reqIDKey        idKey        = "reqID"
)

type LoggerOptions struct {
	IgnorePaths  []string
	CustomProps  func(*http.Request, *http.Response) map[string]interface{}
	GenRequestID func() string
}

func DefaultLoggerOptions() LoggerOptions {
	return LoggerOptions{
		IgnorePaths: []string{"/status", "/favicon.ico"},
		GenRequestID: func() string {
			return uuid.New().String()
		},
	}
}

func HTTPLoggerMiddleware(log logger.Log, opts LoggerOptions) mux.MiddlewareFunc {
	if opts.GenRequestID == nil {
		opts.GenRequestID = DefaultLoggerOptions().GenRequestID
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip logging for ignored paths
			for _, path := range opts.IgnorePaths {
				if r.URL.Path == path {
					next.ServeHTTP(w, r)
					return
				}
			}

			// Set start time and request ID
			start := time.Now()
			reqID := opts.GenRequestID()
			ctx := context.WithValue(r.Context(), reqStartTimeKey, start)
			ctx = context.WithValue(ctx, reqIDKey, reqID)
			r = r.WithContext(ctx)

			// Create response wrapper to capture status code
			rw := &responseWriter{w, http.StatusOK}

			// Process request
			next.ServeHTTP(rw, r)

			// Calculate duration
			duration := time.Since(start)

			// Create base fields
			fields := map[string]interface{}{
				"reqID":      reqID,
				"method":     r.Method,
				"path":       r.URL.Path,
				"status":     rw.status,
				"durationMs": duration.Milliseconds(),
				"ip":         r.RemoteAddr,
				"userAgent":  r.UserAgent(),
			}

			// Add custom properties if provided
			if opts.CustomProps != nil {
				for k, v := range opts.CustomProps(r, rw.response()) {
					fields[k] = v
				}
			}

			// Create logger with fields
			requestLogger := log.HTTP().WithFields(fields)

			// Log based on status code
			switch {
			case rw.status >= 500:
				requestLogger.Error("server error")
			case rw.status >= 400:
				requestLogger.Warn("client error")
			case rw.status >= 300:
				requestLogger.Info("redirection")
			default:
				requestLogger.Info("request completed")
			}
		})
	}
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) response() *http.Response {
	return &http.Response{
		StatusCode: rw.status,
	}
}

// Helper functions to access request context values
func GetRequestID(r *http.Request) string {
	if id, ok := r.Context().Value(reqIDKey).(string); ok {
		return id
	}
	return ""
}

func GetRequestStartTime(r *http.Request) time.Time {
	if start, ok := r.Context().Value(reqStartTimeKey).(time.Time); ok {
		return start
	}
	return time.Time{}
}
