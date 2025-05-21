package usecase

import (
	"context"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/contracts"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/events"
	sharedDomain "github.com/BrockMekonnen/go-clean-starter/internal/_shared/domain"
	authDomain "github.com/BrockMekonnen/go-clean-starter/internal/auth/domain"
	userEvents "github.com/BrockMekonnen/go-clean-starter/internal/user/app/events"
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
func MakeCreateUserUsecase(deps CreateUserDeps, publisher events.Publisher) CreateUserUsecase {
	return events.EventProvider(func(deps CreateUserDeps, enqueue events.EnqueueFunc) CreateUserUsecase {
		return func(ctx context.Context, payload CreateUserParams) (string, error) {
			if existingUser, _ := deps.UserRepo.FindByEmail(ctx, payload.Email); existingUser != nil {
				return "", sharedDomain.NewBusinessError("Email is already in use", "")
			}
			if existingUser, _ := deps.UserRepo.FindByPhone(ctx, payload.Phone); existingUser != nil {
				return "", sharedDomain.NewBusinessError("Phone number is already in use", "")
			}

			id, err := deps.UserRepo.GetNextId(ctx)
			if err != nil {
				return "", sharedDomain.NewBusinessError("Failed to generate user ID", "")
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
				return "", sharedDomain.NewBusinessError("Error occurred while creating user account", "")
			}

			enqueue(userEvents.NewSendOTPEvent(user.Email, userEvents.Verification))

			return user.Id, nil
		}
	})(deps, publisher)
}
