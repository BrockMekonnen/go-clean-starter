package usecase

import (
	"context"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/contracts"
	sharedDomain "github.com/BrockMekonnen/go-clean-starter/internal/_shared/domain"
	"github.com/BrockMekonnen/go-clean-starter/internal/auth/domain"
	userDomain "github.com/BrockMekonnen/go-clean-starter/internal/user/domain"
)

// LoginParams defines the input structure for the service
type GenerateTokenParams struct {
	Email    string
	Password string
}

// GenerateTokenContract makes the function signature readable
type GenerateTokenUsecase = contracts.ApplicationService[GenerateTokenParams, string]

func MakeGenerateTokenUsecase(
	authRepo domain.AuthRepository,
	userRepo userDomain.UserRepository,
) GenerateTokenUsecase {
	return func(ctx context.Context, payload GenerateTokenParams) (string, error) {
		user, err := userRepo.FindByEmail(ctx, payload.Email)

		if err != nil || user == nil {
			return "", sharedDomain.NewBusinessError("Invalid email or password.", "")
		}

		isMatch, err := authRepo.Compare(ctx, payload.Password, user.Password)
		if err != nil || !isMatch {
			return "", sharedDomain.NewBusinessError("Invalid email or password.", "")
		}

		token, err := authRepo.Generate(ctx, domain.Credentials{
			Uid: user.ID, Scope: user.Roles})

		if err != nil {
			return "", sharedDomain.NewBusinessError("Failed to generate authentication token.", "")
		}

		return token, nil
	}
}
