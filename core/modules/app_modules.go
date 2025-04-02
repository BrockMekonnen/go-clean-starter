package modules

import (
	"github.com/BrockMekonnen/go-clean-starter/core/di"
	"github.com/BrockMekonnen/go-clean-starter/internal/auth"
	"github.com/BrockMekonnen/go-clean-starter/internal/user"
)

// InitModules registers all internal modules in the DI container
func RegisterInternalModules() {
	logger := di.GetLogger()

	if err := auth.RegisterAuthModule(); err != nil {
		logger.Fatal("Failed to register auth module:", err)
	}

	if err := user.RegisterUserModule(); err != nil {
		logger.Fatal("Failed to register user module:", err)
	}

	di.GetLogger().Info("All internal modules have been successfully registered.")
}
