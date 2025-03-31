package app

import (
	"github.com/sirupsen/logrus"
	"os"
)

type LoggerRegistry struct {
	log *logrus.Logger
}

// Logger is the global instance of logrus.Logger
var Logger = logrus.New()

func InitLogger() {
	// Set log format (JSON format like Pino)
	Logger.SetFormatter(&logrus.JSONFormatter{})

	// Set log level (default: Info)
	Logger.SetLevel(logrus.InfoLevel)

	// Set output to stdouta
	Logger.SetOutput(os.Stdout)
}
