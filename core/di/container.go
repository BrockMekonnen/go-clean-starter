package di

import (
	"sync"

	"github.com/BrockMekonnen/go-clean-starter/core"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/logger"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/hashids"
	"github.com/gorilla/mux"
	"go.uber.org/dig"
)

var (
	container *dig.Container
	once      sync.Once
)

// Init initializes the global DI container
func InitContainer() *dig.Container {
	once.Do(func() {
		container = dig.New()
	})
	return container
}

// Get returns the global container instance
func Get() *dig.Container {
	if container == nil {
		panic("container not initialized - call Init() first")
	}
	return container
}

// Provide adds a constructor to the container
func Provide(constructor interface{}, opts ...dig.ProvideOption) error {
	return Get().Provide(constructor, opts...)
}

// Invoke calls a function with dependencies from the container
func Invoke(function interface{}, opts ...dig.InvokeOption) error {
	return Get().Invoke(function, opts...)
}

// MustResolve panics if the dependency cannot be resolved
func MustResolve[T any]() T {
	var t T
	if err := Get().Invoke(func(dep T) {
		t = dep
	}); err != nil {
		panic("failed to resolve dependency: " + err.Error())
	}
	return t
}

// helper function to handle dependency injection errors
func ProvideWrapper(name string, provider interface{}, opts ...dig.ProvideOption) error {
	logger := GetLogger()
	if err := Provide(provider, opts...); err != nil {
		logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to register " + name)
		return err
	}
	return nil
}

func GetLogger() *logger.Log {
	var l *logger.Log
	if err := Get().Invoke(func(logger *logger.Log) {
		l = logger
	}); err != nil {
		panic("logger not available in container: " + err.Error())
	}
	return l
}

// GetDatabase now uses core.DatabaseProvider
func GetDatabase() core.DatabaseProvider {
	var db core.DatabaseProvider
	if err := container.Invoke(func(d core.DatabaseProvider) {
		db = d
	}); err != nil {
		panic("database not available in container: " + err.Error())
	}
	return db
}

// Add this to your container.go file
func GetApiRouter() *mux.Router {
	var router *mux.Router
	if err := Get().Invoke(func(s *core.ServerRegistry) {
		router = s.ApiRouter
	}); err != nil {
		panic("api router not available in container: " + err.Error())
	}
	return router
}

// Add this to your container.go file
func GetAuthRouter() *mux.Router {
	var router *mux.Router
	if err := Get().Invoke(func(s *core.ServerRegistry) {
		router = s.AuthRouter
	}); err != nil {
		panic("auth router not available in container: " + err.Error())
	}
	return router
}

// GetHashID retrieves the HashID service from the DI container
func GetHashID() hashids.HashID {
	var h hashids.HashID
	if err := Get().Invoke(func(hid hashids.HashID) {
		h = hid
	}); err != nil {
		panic("HashID service not available in container: " + err.Error())
	}
	return h
}