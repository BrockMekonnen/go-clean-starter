package app

import (
	"sync"

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
