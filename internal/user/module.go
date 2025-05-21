package user

import (
	"github.com/BrockMekonnen/go-clean-starter/core/di"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/events"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/app/usecase"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/delivery/listeners"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/infrastructure"
)

func MakeUserModule() error {
	//* Get dependencies
	logger := di.GetLogger()
	db := di.GetDatabase().GetDB()

	//* Initialize this module tables
	if err := infrastructure.InitUserTable(db); err != nil {
		logger.WithFields(map[string]interface{}{"error": err.Error()}).Error("Failed to initialize user table")
		return err
	}

	//* Register this module repository
	if err := di.ProvideWrapper("UserRepository", infrastructure.MakeUserRepository); err != nil {
		return err
	}

	//* Register this app queries
	if err := di.ProvideWrapper("FindUserById", infrastructure.MakeFindUserById); err != nil {
		return err
	}

	if err := di.ProvideWrapper("FindUsers", infrastructure.MakeFindUsers); err != nil {
		return err
	}

	//* Register this module app usecases
	if err := di.ProvideWrapper("CreateUserUsecase", usecase.MakeCreateUserUsecase); err != nil {
		return err
	}

	if err := di.ProvideWrapper("DeleteUserUsecase", usecase.MakeDeleteUserUsecase); err != nil {
		return err
	}

	if err := di.ProvideWrapper("GenerateTokenUsecase", usecase.MakeGenerateTokenUsecase); err != nil {
		return err
	}

	//* Invoke Consumers
	if err := di.Invoke(func(subscriber events.Subscriber) error {
		return listeners.SendOTPEventListener(subscriber)()
	}); err != nil {
		return err
	}

	//* Register Routes
	MakeUserRoutes()

	logger.Info("User module initialized successfully.")
	return nil
}
