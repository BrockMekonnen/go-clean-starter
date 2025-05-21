package core

import (
	"github.com/BrockMekonnen/go-clean-starter/core/lib/logger"
	"github.com/BrockMekonnen/go-clean-starter/internal/auth"
	"github.com/BrockMekonnen/go-clean-starter/internal/user"
)

// InitModules registers all internal modules in the DI container
func MakeInternalModules(logger logger.Log) {

	if err := auth.MakeAuthModule(); err != nil {
		logger.Fatal("Failed to register auth module:", err)
	}

	if err := user.MakeUserModule(); err != nil {
		logger.Fatal("Failed to register user module:", err)
	}

}
