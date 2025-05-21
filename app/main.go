package main

import (
	"os"
	"os/signal"
	"syscall"

	core "github.com/BrockMekonnen/go-clean-starter/core"
	di "github.com/BrockMekonnen/go-clean-starter/core/di"
	events "github.com/BrockMekonnen/go-clean-starter/core/lib/events"
	hashids "github.com/BrockMekonnen/go-clean-starter/core/lib/hashids"
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
	config := core.LoadConfig(logger)
	err = container.Provide(func() *core.AppConfig { return config })
	if err != nil {
		logger.Fatal("Failed to provide logger", err)
	}

	//* Initialize Database
	dbProvider, shutdownDB, err := core.NewDatabase(config.Database, *logger)
	if err != nil {
		logger.Fatal("Failed to initialize database:", err)
	}
	err = container.Provide(func() core.DatabaseProvider { return dbProvider })
	if err != nil {
		logger.Fatal("Failed to provide database provider", err)
	}
	err = container.Provide(func() *gorm.DB { return dbProvider.GetDB() })
	if err != nil {
		logger.Fatal("Failed to provide db", err)
	}
	//* Initialize Id Hasher
	err = container.Provide(func() (hashids.HashID, error) {
		return hashids.NewHashIDService(config.Encryption)
	})
	if err != nil {
		logger.Fatal("Failed to provide HashID", err)
	}

	//* Initialize Pub/Sub
	pubsub := core.InitPubSub(logger)
	// Register *EventEmitterPubSub directly (optional, but helpful)
	err = container.Provide(func() *events.EventEmitterPubSub { return pubsub })
	if err != nil {
		logger.Fatal("Failed to provide *EventEmitterPubSub", err)
	}
	err = container.Provide(func() events.Subscriber { return pubsub })
	if err != nil {
		logger.Fatal("Failed to provide Subscriber", err)
	}

	err = container.Provide(func() events.Publisher { return pubsub })
	if err != nil {
		logger.Fatal("Failed to provide Publisher", err)
	}

	//* Initialize Server
	server, shutdownServer := core.NewServer(*config, container, *logger)
	err = container.Provide(func() *core.ServerRegistry { return server })
	if err != nil {
		logger.Fatal("Failed to provide server", err)
	}

	//* Register internal modules
	modules.RegisterInternalModules()

	//* Start Server
	core.StartServer(server, *logger)

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
