package main

import (
	"os"
	"os/signal"
	"syscall"

	app "github.com/BrockMekonnen/go-clean-starter/core"
)

func main() {
	// Initialize Logger
	app.InitLogger()
	logger := app.Logger
	logger.Info("Starting application...")

	// Load Configurations
	config := app.LoadConfig()

	// Initialize Dependency Injection Container
	container := app.InitContainer()

	// Initialize Server
	server, shutdownServer := app.NewServer(config.HTTP, logger)

	// Initialize Database
	db, shutdownDB, err := app.NewDatabase(config.Database)
	if err != nil {
		logger.Fatal("Failed to initialize database:", err)
	}

	// Provide Dependencies to Container
	container.Provide(func() *app.AppConfig { return config })
	container.Provide(func() *app.ServerRegistry { return server })
	container.Provide(func() *app.DatabaseRegistry { return db })
	container.Provide(func() *app.LoggerRegistry { return &app.LoggerRegistry{} })

	// Start Server
	app.StartServer(server, logger)

	// Graceful Shutdown Handling
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	<-shutdownChan
	logger.Info("Shutting down application...")

	// Perform Cleanup
	shutdownServer()
	shutdownDB()
	logger.Info("Application exited successfully.")
}
