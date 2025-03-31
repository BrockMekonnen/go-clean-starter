package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
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
}

// NewServer initializes and returns the HTTP server
func NewServer(config ServerConfig, logger *logrus.Logger) (*ServerRegistry, func()) {
	// Create a new router
	rootRouter := mux.NewRouter()

	// Middleware: Request Logging
	rootRouter.Use(loggingMiddleware(logger))

	// Middleware: Graceful shutdown
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt)

	// Middleware: CORS (if enabled)
	if config.Cors {
		corsHandler := cors.Default().Handler
		rootRouter.Use(corsHandler)
	}

	// Health check route
	rootRouter.HandleFunc("/status", statusHandler).Methods("GET")

	// API router
	apiRouter := rootRouter.PathPrefix("/api").Subrouter()
	// Register API endpoints here...

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
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
	}, shutdown
}

// loggingMiddleware logs HTTP requests
func loggingMiddleware(logger *logrus.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			logger.Infof("[%s] %s %s %s", r.Method, r.RequestURI, r.RemoteAddr, time.Since(start))
		})
	}
}

// statusHandler returns a 200 OK response
func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// StartServer starts the HTTP server
func StartServer(server *ServerRegistry, logger *logrus.Logger) {
	go func() {
		logger.Infof("Starting server on %s", server.HttpServer.Addr)
		if err := server.HttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Could not start server: %v", err)
		}
	}()
}
