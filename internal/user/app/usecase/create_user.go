package usecase

import (
	"context"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/contracts"
	sharedDomain "github.com/BrockMekonnen/go-clean-starter/internal/_shared/domain"
	authDomain "github.com/BrockMekonnen/go-clean-starter/internal/auth/domain"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/domain"
)

// CreateUserDTO defines the input contract
type CreateUserParams struct {
	FirstName                string
	LastName                 string
	Email                    string
	Phone                    string
	Password                 string
	IsTermAndConditionAgreed bool
}

// CreateUserDeps declares all required dependencies
type CreateUserDeps struct {
	UserRepo domain.UserRepository
	AuthRepo authDomain.AuthRepository
}

// CreateUserContract makes the function signature readable
type CreateUserUsecase = contracts.ApplicationService[CreateUserParams, string]

// NewCreateUserUsecase constructs the use case with explicit dependencies
func NewCreateUserUsecase(deps CreateUserDeps) CreateUserUsecase {
	return func(ctx context.Context, payload CreateUserParams) (string, error) {
		id, err := deps.UserRepo.GetNextId(ctx)
		if err != nil {
			return "", sharedDomain.NewBusinessError("Failed to generate user ID",  "")
		}

		hashedPassword, err := deps.AuthRepo.Hash(ctx, payload.Password)

		if err != nil {
			return "", sharedDomain.NewBusinessError("Failed to hash user password", "")
		}

		user, err := domain.NewUser(domain.UserProps{
			Id:        id,
			FirstName: payload.FirstName,
			LastName:  payload.LastName,
			Phone:     payload.Phone,
			Email:     payload.Email,
			Password:  hashedPassword,
		})

		if err != nil {
			return "", err
		}

		if err := deps.UserRepo.Store(ctx, user); err != nil {
			return "", sharedDomain.NewBusinessError("Error occurred while creating user account!", "")
		}

		return user.Id, nil
	}
}
