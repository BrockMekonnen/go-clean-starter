package core

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/logger"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/middleware"
	"github.com/BrockMekonnen/go-clean-starter/internal/_shared/delivery"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.uber.org/dig"
)

// ServerConfig holds the configuration for the server
type ServerConfig struct {
	Host string
	Port int
	Cors bool
}

// ServerRegistry holds the server components
type ServerRegistry struct {
	HttpServer *http.Server
	RootRouter *mux.Router
	ApiRouter  *mux.Router
	AuthRouter *mux.Router
}

// NewServer initializes and returns the HTTP server
func NewServer(config AppConfig, container *dig.Container, logger logger.Log) (*ServerRegistry, func()) {
	// Create a new router
	rootRouter := mux.NewRouter()

	// Middleware: Request Logging
	logOpts := middleware.DefaultLoggerOptions()
	rootRouter.Use(middleware.HTTPLoggerMiddleware(logger, logOpts))
	// rootRouter.Use(loggingMiddleware(logger))
	rootRouter.Use(middleware.RequestContainerMiddleware(container, logger))

	// Middleware: Graceful shutdown
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt)

	// Middleware: CORS (if enabled)
	if config.HTTP.Cors {
		corsHandler := cors.Default().Handler
		rootRouter.Use(corsHandler)
	}

	// Status check route
	rootRouter.HandleFunc("/status", middleware.StatusHandler(config.StartedAt)).Methods("GET")

	// API router
	apiRouter := rootRouter.PathPrefix("/api").Subrouter()
	authRouter := rootRouter.PathPrefix("/api").Subrouter()

	opts := delivery.DefaultErrorConverters
	rootRouter.Use(middleware.ErrorHandler(opts, &logger))

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.HTTP.Host, config.HTTP.Port),
		Handler:      rootRouter,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown function
	shutdown := func() {
		logger.Info("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatalf("Server forced to shutdown: %v", err)
		}
		close(shutdownChan)
		logger.Info("Server gracefully stopped")
	}

	return &ServerRegistry{
		HttpServer: server,
		RootRouter: rootRouter,
		ApiRouter:  apiRouter,
		AuthRouter: authRouter,
	}, shutdown
}

// StartServer starts the HTTP server
func StartServer(server *ServerRegistry, logger logger.Log) {
	go func() {
		logger.Infof("Starting server on %s", server.HttpServer.Addr)
		if err := server.HttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Could not start server: %v", err)
		}
	}()
}
