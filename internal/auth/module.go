package auth

import (
	"github.com/BrockMekonnen/go-clean-starter/core/di"
	"github.com/BrockMekonnen/go-clean-starter/internal/auth/app/usecase"
	"github.com/BrockMekonnen/go-clean-starter/internal/auth/delivery"
	"github.com/BrockMekonnen/go-clean-starter/internal/auth/domain"
	"github.com/BrockMekonnen/go-clean-starter/internal/auth/infrastructure"

	"go.uber.org/dig"
)

func RegisterAuthModule() error {
	logger := di.GetLogger()
	authRouter := di.GetAuthRouter()

	//* Register this module repository
	if err := di.ProvideWrapper("AuthRepository",
		infrastructure.NewAuthRepository, dig.As(new(domain.AuthRepository)),
	); err != nil {
		return err
	}

	//* Register this module use cases
	if err := di.ProvideWrapper("VerifyTokenUsecase",
		func(authRepo domain.AuthRepository) usecase.VerifyTokenUsecase {
			return usecase.NewVerifyTokenUsecase(usecase.VerifyTokenDeps{AuthRepository: authRepo})
		},
	); err != nil {
		return err
	}

	//* Register handlers
	if err := di.ProvideWrapper("VerifyTokenHandlerDeps",
		func(uc usecase.VerifyTokenUsecase) delivery.VerifyTokenHandlerDeps {
			return delivery.VerifyTokenHandlerDeps{VerifyToken: uc}
		},
	); err != nil {
		return err
	}

	//* Add Verify Token Middleware
	verifyHandler := delivery.NewVerifyTokenHandler(di.MustResolve[delivery.VerifyTokenHandlerDeps]())
	authRouter.Use(verifyHandler)

	logger.Info("Auth module initialized successfully.")
	return nil
}
