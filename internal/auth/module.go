package auth

import (
	"github.com/BrockMekonnen/go-clean-starter/core/di"
	"github.com/BrockMekonnen/go-clean-starter/internal/auth/app/usecase"
	"github.com/BrockMekonnen/go-clean-starter/internal/auth/domain"
	"github.com/BrockMekonnen/go-clean-starter/internal/auth/infrastructure"
	"github.com/BrockMekonnen/go-clean-starter/internal/auth/interface/handlers"
)

func MakeAuthModule() error {
	logger := di.GetLogger()
	authRouter := di.GetAuthRouter()

	//* Register this module repository
	if err := di.ProvideWrapper("AuthRepository", infrastructure.NewAuthRepository); err != nil {
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

	//* Add Verify Token Middleware
	verifyTokenMiddleware := handlers.VerifyTokenMiddleware(di.MustResolve[usecase.VerifyTokenUsecase]())
	authRouter.Use(verifyTokenMiddleware)

	logger.Info("Auth module initialized successfully.")
	return nil
}
