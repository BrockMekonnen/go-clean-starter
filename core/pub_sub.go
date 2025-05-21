package core

import (
	"github.com/BrockMekonnen/go-clean-starter/core/lib/events"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/logger"
)

// InitPubSub initializes the pub/sub system
func InitPubSub(log *logger.Log) *events.EventEmitterPubSub {
	pubsub := events.NewEventEmitterPubSub()

	if err := pubsub.Start(); err != nil {
		log.Fatal("Failed to start pubsub", err)
	}

	log.Info("PubSub Successfully Started.")

	return pubsub
}
