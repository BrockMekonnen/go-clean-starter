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

// CreateUserContract makes the function signature readable
type CreateUserUsecase = contracts.ApplicationService[CreateUserParams, string]

// NewCreateUserUsecase constructs the use case with explicit dependencies
func MakeCreateUserUsecase(
	userRepo domain.UserRepository,
	authRepo authDomain.AuthRepository,
	publisher events.Publisher,
) CreateUserUsecase {
	return events.EventProvider(func(enqueue events.EnqueueFunc) CreateUserUsecase {
		return func(ctx context.Context, payload CreateUserParams) (string, error) {
			if existingUser, _ := userRepo.FindByEmail(ctx, payload.Email); existingUser != nil {
				return "", sharedDomain.NewBusinessError("Email is already in use", "")
			}
			if existingUser, _ := userRepo.FindByPhone(ctx, payload.Phone); existingUser != nil {
				return "", sharedDomain.NewBusinessError("Phone number is already in use", "")
			}

			id, err := userRepo.GetNextId(ctx)
			if err != nil {
				return "", sharedDomain.NewBusinessError("Failed to generate user ID", "")
			}

			hashedPassword, err := authRepo.Hash(ctx, payload.Password)
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

			if err := userRepo.Store(ctx, user); err != nil {
				return "", sharedDomain.NewBusinessError("Error occurred while creating user account", "")
			}

			enqueue(userEvents.NewSendOTPEvent(user.Email, userEvents.Verification))

			return user.Id, nil
		}
	})(publisher)
}
