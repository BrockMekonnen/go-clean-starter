package middleware

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type ShutdownMiddleware struct {
	isShuttingDown *atomic.Bool
	server         *http.Server
	forceTimeout   time.Duration
	logger         *logrus.Logger
}

// NewGracefulShutdown creates a new graceful shutdown middleware with Logrus
func NewGracefulShutdown(server *http.Server, forceTimeout time.Duration, logger *logrus.Logger) *ShutdownMiddleware {
	return &ShutdownMiddleware{
		isShuttingDown: &atomic.Bool{},
		server:        server,
		forceTimeout:  forceTimeout,
		logger:        logger,
	}
}

// Handler returns the middleware function that rejects new requests during shutdown
func (s *ShutdownMiddleware) Handler() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if s.isShuttingDown.Load() {
				c.Response().Header().Set("Connection", "close")
				return c.String(http.StatusServiceUnavailable, "Server is in the process of restarting")
			}
			return next(c)
		}
	}
}

// Hook initiates the graceful shutdown process
func (s *ShutdownMiddleware) Hook() error {
	// Skip if not in production or server isn't running
	if os.Getenv("APP_ENV") != "production" || s.server == nil {
		return nil
	}

	s.isShuttingDown.Store(true)
	s.logger.Warn("Shutting down server")

	// Create context with force timeout
	ctx, cancel := context.WithTimeout(context.Background(), s.forceTimeout)
	defer cancel()

	// Start force shutdown timer
	forceTimer := time.AfterFunc(s.forceTimeout, func() {
		s.logger.Error("Could not close connections in time, forcefully shutting down")
		cancel()
	})
	defer forceTimer.Stop()

	// Shutdown the server
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}

	s.logger.Info("Closed out remaining connections")
	return nil
}

// Listen listens for OS signals and triggers shutdown
func (s *ShutdownMiddleware) Listen() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err := s.Hook(); err != nil {
		s.logger.WithError(err).Error("Error during shutdown")
	}
}