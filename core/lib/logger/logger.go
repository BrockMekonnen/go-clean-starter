package logger

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

// Interface defines the logger contract
type Interface interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	WithFields(fields map[string]interface{}) Interface
	SetOutput(w io.Writer) // Added for testing
	HTTP() Interface       // Returns a logger configured for HTTP logging
}

// Log is a wrapper that properly maintains log fields
type Log struct {
	entry      *logrus.Entry  // Changed from *logrus.Logger to *logrus.Entry
	httpLogger *logrus.Logger // Separate logger for HTTP
}

// NewLogger creates a new logger instance
func NewLogger() *Log {
	// Create main logger (text format)
	mainLogger := logrus.New()
	mainLogger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02T15:04:05Z07:00",
	})
	mainLogger.SetLevel(logrus.InfoLevel)
	mainLogger.SetOutput(os.Stdout)

	// Create HTTP logger (JSON format)
	httpLogger := logrus.New()
	httpLogger.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true,
		TimestampFormat: "2006-01-02T15:04:05Z07:00",
	})
	httpLogger.SetLevel(logrus.InfoLevel)
	httpLogger.SetOutput(os.Stdout) // Same output as main logger

	return &Log{
		entry:      logrus.NewEntry(mainLogger),
		httpLogger: httpLogger,
	}
}

// HTTP returns a logger instance configured for HTTP logging
func (l *Log) HTTP() Interface {
	return &Log{
		entry: logrus.NewEntry(l.httpLogger),
	}
}

// SetOutput sets the output for both loggers
func (l *Log) SetOutput(w io.Writer) {
	l.entry.Logger.SetOutput(w)
	l.httpLogger.SetOutput(w)
}

// WithFields returns a new logger instance with the specified fields
func (l *Log) WithFields(fields map[string]interface{}) Interface {
	return &Log{
		entry:      l.entry.WithFields(fields),
		httpLogger: l.httpLogger, // Share the same HTTP logger
	}
}

// Implement all Interface methods using the entry
func (l *Log) Debug(args ...interface{}) {
	l.entry.Debug(args...)
}

func (l *Log) Info(args ...interface{}) {
	l.entry.Info(args...)
}

func (l *Log) Warn(args ...interface{}) {
	l.entry.Warn(args...)
}

func (l *Log) Error(args ...interface{}) {
	l.entry.Error(args...)
}

func (l *Log) Fatal(args ...interface{}) {
	l.entry.Fatal(args...)
}

func (l *Log) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args...)
}

func (l *Log) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

func (l *Log) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args...)
}

func (l *Log) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

func (l *Log) Fatalf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}
