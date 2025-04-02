package usecase

import (
	"context"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/ddd"
	sharedDomain "github.com/BrockMekonnen/go-clean-starter/internal/_shared/domain"
	"github.com/BrockMekonnen/go-clean-starter/internal/auth/domain"
	userDomain "github.com/BrockMekonnen/go-clean-starter/internal/user/domain"
)

// Dependencies struct for passing dependencies
type GenerateTokenDeps struct {
	AuthRepository domain.AuthRepository
	UserRepository userDomain.UserRepository
}

// LoginParams defines the input structure for the service
type GenerateTokenParams struct {
	Email    string
	Password string
}

// GenerateTokenContract makes the function signature readable
type GenerateTokenUsecase = ddd.ApplicationService[GenerateTokenParams, string]

func NewGenerateTokenUsecase(deps GenerateTokenDeps) GenerateTokenUsecase {
	return func(ctx context.Context, payload GenerateTokenParams) (string, error) {
		user, err := deps.UserRepository.FindByEmail(ctx, payload.Email)

		if err != nil || user == nil {
			return "", sharedDomain.NewBusinessError("Incorrect Email.", "USER_NOT_FOUND")
		}

		isMatch, err := deps.AuthRepository.Compare(ctx, payload.Password, user.Password)
		if err != nil || !isMatch {
			return "", sharedDomain.NewBusinessError("Incorrect Password.", "INVALID_PASSWORD")
		}

		token, err := deps.AuthRepository.Generate(ctx, domain.Credentials{
			Uid: user.Id, Scope: user.Roles})

		if err != nil {
			return "", sharedDomain.NewBusinessError("Decryption Failed", "INVALID_TOKEN")
		}

		return token, nil
	}
}
