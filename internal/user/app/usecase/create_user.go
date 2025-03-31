package usecase

import (
	"context"
	"github.com/BrockMekonnen/go-clean-starter/internal/_lib/ddd"
	SharedDomain "github.com/BrockMekonnen/go-clean-starter/internal/_shared/domain"
	AuthDomain "github.com/BrockMekonnen/go-clean-starter/internal/auth/domain"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/domain"
)

// RegisterUserDTO defines the data transfer object for creating a user
type CreateUserDTO struct {
	FirstName                string
	LastName                 string
	Email                    string
	Phone                    string
	Password                 string
	IsTermAndConditionAgreed bool
}

// CreateUser is the application service type that creates a user
type CreateUser = ddd.ApplicationService[CreateUserDTO, string]

// RegisterUserUsecase implements the CreateUser application service
type CreateUserUsecase struct {
	UserRepo domain.UserRepository
	AuthRepo AuthDomain.AuthRepository
}

// NewRegisterUserUsecase creates a new instance of RegisterUserUsecase
func NewRegisterUserUsecase(userRepo domain.UserRepository, authRepo AuthDomain.AuthRepository) *CreateUserUsecase {
	return &CreateUserUsecase{UserRepo: userRepo, AuthRepo: authRepo}
}

// Execute is the method that processes the user creation
func (uc *CreateUserUsecase) Execute(ctx context.Context, payload CreateUserDTO) (uint, error) {
	// Generate user ID and hash the password
	id, err := uc.UserRepo.GetNextId(ctx)
	if err != nil {
		return 0, SharedDomain.NewBusinessError("Failed to generate user ID.", "USER_ID_GENERATION_FAILED")
	}

	hashedPassword, err := uc.AuthRepo.Hash(payload.Password)
	if err != nil {
		return 0, err
	}

	// Create a new user using the domain logic
	user, err := domain.NewUser(domain.UserProps{
		Id:        id,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Phone:     payload.Phone,
		Email:     payload.Email,
		Password:  hashedPassword,
	})

	// Store the user in the repository
	err = uc.UserRepo.Store(ctx, user)
	if err != nil {
		return 0, err
	}

	// Return the user ID as string
	return user.Id, nil
}
