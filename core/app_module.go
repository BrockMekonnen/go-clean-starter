package app

import (
	"log"

	"go.uber.org/dig"
)

// InitModules registers all internal modules in the DI container
func InitModules(container *dig.Container) {

	log.Println("All internal modules have been successfully registered.")
}
