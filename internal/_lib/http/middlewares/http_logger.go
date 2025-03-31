package middleware

import (
	"math"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

const (
	reqStartTimeKey = "reqStartTime"
	reqIDKey        = "reqID"
)

type LoggerOptions struct {
	IgnorePaths  []string
	CustomProps  func(c echo.Context) logrus.Fields
	GenRequestID func() string
}

func DefaultLoggerOptions() LoggerOptions {
	return LoggerOptions{
		IgnorePaths: []string{"/status", "/favicon.ico"},
		GenRequestID: func() string {
			return "req:" + generateUUID()
		},
	}
}

func HTTPLogger(opts LoggerOptions) echo.MiddlewareFunc {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)

	shouldSkip := func(path string) bool {
		for _, p := range opts.IgnorePaths {
			if p == path {
				return true
			}
		}
		return false
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip logging for ignored paths
			if shouldSkip(c.Path()) {
				return next(c)
			}

			// Generate request ID if not already set
			reqID := c.Response().Header().Get(echo.HeaderXRequestID)
			if reqID == "" {
				reqID = opts.GenRequestID()
				c.Response().Header().Set(echo.HeaderXRequestID, reqID)
			}

			// Store start time and request ID in context
			c.Set(reqStartTimeKey, time.Now())
			c.Set(reqIDKey, reqID)

			// Process request
			err := next(c)

			// Calculate duration
			startTime, _ := c.Get(reqStartTimeKey).(time.Time)
			duration := time.Since(startTime)
			ms := math.Round(float64(duration.Nanoseconds())/1000000) / 1000

			// Prepare base log fields
			fields := logrus.Fields{
				"reqId":   reqID,
				"method":  c.Request().Method,
				"path":    c.Path(),
				"status":  c.Response().Status,
				"latency": ms,
			}

			// Add custom properties if provided
			if opts.CustomProps != nil {
				for k, v := range opts.CustomProps(c) {
					fields[k] = v
				}
			}

			// Determine log level based on status code
			status := c.Response().Status
			entry := logger.WithFields(fields)

			switch {
			case status >= 500 || err != nil:
				if err != nil {
					entry = entry.WithError(err)
				}
				entry.Error("Server error")
			case status >= 400:
				entry.Warn("Client error")
			case status >= 300:
				entry.Trace("Redirection")
			default:
				entry.Info("Request completed")
			}

			return err
		}
	}
}

func generateUUID() string {
	// Implement your preferred UUID generation
	// Example using google/uuid:
	// return uuid.New().String()
	return "temp-uuid" // Replace with actual UUID generation
}