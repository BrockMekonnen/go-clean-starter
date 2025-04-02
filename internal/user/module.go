package user

import (
	"github.com/BrockMekonnen/go-clean-starter/core/di"
	authDomain "github.com/BrockMekonnen/go-clean-starter/internal/auth/domain"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/app/query"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/app/usecase"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/delivery"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/domain"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/infrastructure"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

func RegisterUserModule() error {
	// Get dependencies
	logger := di.GetLogger()
	db := di.GetDatabase().GetDB()

	//* Initialize this module tables
	if err := infrastructure.InitUserTable(db); err != nil {
		logger.WithFields(map[string]interface{}{"error": err.Error()}).Error("Failed to initialize user table")
		return err
	}

	//* Register this module repository
	if err := di.ProvideWrapper("UserRepository",
		infrastructure.NewUserRepository, dig.As(new(domain.UserRepository)),
	); err != nil {
		return err
	}

	//* Register this module query
	if err := di.ProvideWrapper("FindUserById",
		func(db *gorm.DB) query.FindUserById {
			return infrastructure.NewFindUserByIdHandler(db)
		},
	); err != nil {
		return err
	}

	//* Register this module use cases
	if err := di.ProvideWrapper("CreateUserUsecase",
		func(repo domain.UserRepository, authRepo authDomain.AuthRepository) usecase.CreateUserUsecase {
			return usecase.NewCreateUserUsecase(usecase.CreateUserDeps{
				UserRepo: repo,
				AuthRepo: authRepo,
			})
		},
	); err != nil {
		return err
	}

	if err := di.ProvideWrapper("DeleteUserUsecase",
		func(repo domain.UserRepository) usecase.DeleteUserUsecase {
			return usecase.NewDeleteUserUsecase(usecase.DeleteUserDeps{
				UserRepo: repo,
			})
		},
	); err != nil {
		return err
	}

	if err := di.ProvideWrapper("GenerateTokenUsecase",
		func(authRepo authDomain.AuthRepository, userRepo domain.UserRepository) usecase.GenerateTokenUsecase {
			return usecase.NewGenerateTokenUsecase(usecase.GenerateTokenDeps{
				AuthRepository: authRepo, UserRepository: userRepo,
			})
		},
	); err != nil {
		return err
	}

	//* Register handlers
	if err := di.ProvideWrapper("CreateUserHandlerDeps",
		func(uc usecase.CreateUserUsecase) delivery.CreateUserHandlerDeps {
			return delivery.CreateUserHandlerDeps{CreateUser: uc}
		},
	); err != nil {
		return err
	}

	if err := di.ProvideWrapper("DeleteUserHandlerDeps",
		func(uc usecase.DeleteUserUsecase) delivery.DeleteUserHandlerDeps {
			return delivery.DeleteUserHandlerDeps{DeleteUser: uc}
		},
	); err != nil {
		return err
	}

	if err := di.ProvideWrapper("GetUserHandlerDeps",
		func(q query.FindUserById) delivery.GetUserHandlerDeps {
			return delivery.GetUserHandlerDeps{FindUserById: q}
		},
	); err != nil {
		return err
	}

	if err := di.ProvideWrapper("GenerateTokenHandlerDeps",
		func(uc usecase.GenerateTokenUsecase) delivery.GenerateTokenHandlerDeps {
			return delivery.GenerateTokenHandlerDeps{GenerateToken: uc}
		},
	); err != nil {
		return err
	}

	//* Register Routes
	RegisterUserRoutes()

	logger.Info("User module initialized successfully.")
	return nil
}
