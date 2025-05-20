package usecase

import (
	"context"
	"fmt"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/contracts"
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
type GenerateTokenUsecase = contracts.ApplicationService[GenerateTokenParams, string]

func NewGenerateTokenUsecase(deps GenerateTokenDeps) GenerateTokenUsecase {
	return func(ctx context.Context, payload GenerateTokenParams) (string, error) {
		user, err := deps.UserRepository.FindByEmail(ctx, payload.Email)

		if err != nil || user == nil {
			return "", sharedDomain.NewBusinessError("Incorrect Email.", "")
		}

		fmt.Println("password: " + user.Password)
		isMatch, err := deps.AuthRepository.Compare(ctx, payload.Password, user.Password)
		fmt.Printf("after is match %t \n", isMatch)
		fmt.Printf("%v", err)
		if err != nil || !isMatch {
			return "", sharedDomain.NewBusinessError("Incorrect Password.", "")
		}

		token, err := deps.AuthRepository.Generate(ctx, domain.Credentials{
			Uid: user.Id, Scope: user.Roles})

		if err != nil {
			return "", sharedDomain.NewBusinessError("Decryption Failed", "")
		}

		return token, nil
	}
}
