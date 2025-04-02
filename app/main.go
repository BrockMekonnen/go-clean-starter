package main

import (
	"os"
	"os/signal"
	"syscall"

	app "github.com/BrockMekonnen/go-clean-starter/core"
	di "github.com/BrockMekonnen/go-clean-starter/core/di"
	log "github.com/BrockMekonnen/go-clean-starter/core/lib/logger"
	modules "github.com/BrockMekonnen/go-clean-starter/core/modules"
	"gorm.io/gorm"
)

func main() {
	//* Initialize Dependency Injection Container
	container := di.InitContainer()

	//* Initialize Logger
	logger := log.NewLogger()
	logger.Info("Starting application...")
	err := container.Provide(func() *log.Log { return logger })
	if err != nil {
		logger.Fatal(err)
	}

	//* Load Configurations
	config := app.LoadConfig(logger)
	err = container.Provide(func() *app.AppConfig { return config })
	if err != nil {
		logger.Fatal("Failed to provide logger", err)
	}

	//* Initialize Database
	dbProvider, shutdownDB, err := app.NewDatabase(config.Database, *logger)
	if err != nil {
		logger.Fatal("Failed to initialize database:", err)
	}
	err = container.Provide(func() app.DatabaseProvider { return dbProvider })
	if err != nil {
		logger.Fatal("Failed to provide database provider", err)
	}
	err = container.Provide(func() *gorm.DB { return dbProvider.GetDB() }) 
	if err != nil {
		logger.Fatal("Failed to provide db", err)
	}

	//* Initialize Server
	server, shutdownServer := app.NewServer(*config, container, *logger)
	err = container.Provide(func() *app.ServerRegistry { return server })
	if err != nil {
		logger.Fatal("Failed to provide server", err)
	}

	//* Register internal modules
	modules.RegisterInternalModules()

	//* Start Server
	app.StartServer(server, *logger)

	//* Graceful Shutdown Handling
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	<-shutdownChan
	logger.Info("Shutting down application...")

	//* Perform Cleanup
	shutdownServer()
	shutdownDB()
	logger.Info("Application exited successfully.")
}
