package usecase

import (
	"context"

	"github.com/BrockMekonnen/go-clean-starter/core/di"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/ddd"
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
type CreateUserUsecase = ddd.ApplicationService[CreateUserParams, uint]

// NewCreateUserUsecase constructs the use case with explicit dependencies
func NewCreateUserUsecase(deps CreateUserDeps) CreateUserUsecase {
	return func(ctx context.Context, payload CreateUserParams) (uint, error) {
		id, err := deps.UserRepo.GetNextId(ctx)
		if err != nil {
			return 0, sharedDomain.NewBusinessError(
				"Failed to generate user ID", 
				"USER_ID_GENERATION_FAILED",
			)
		}

		hashedPassword, err := deps.AuthRepo.Hash(ctx, payload.Password)

		di.GetLogger().Info(`hashedPassword: `)

		if err != nil {
			return 0, sharedDomain.NewBusinessError("Failed to hash user password", "PASSWORD_HASH_FAILED")
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
			return 0, err
		}

		if err := deps.UserRepo.Store(ctx, user); err != nil {
			return 0, err
		}

		return user.Id, nil
	}
}

// import (
// 	"context"
// 	"github.com/BrockMekonnen/go-clean-starter/core/lib/ddd"
// 	SharedDomain "github.com/BrockMekonnen/go-clean-starter/internal/_shared/domain"
// 	// AuthDomain "github.com/BrockMekonnen/go-clean-starter/internal/auth/domain"
// 	"github.com/BrockMekonnen/go-clean-starter/internal/user/domain"
// )

// // RegisterUserDTO defines the data transfer object for creating a user
// type CreateUserDTO struct {
// 	FirstName                string
// 	LastName                 string
// 	Email                    string
// 	Phone                    string
// 	Password                 string
// 	IsTermAndConditionAgreed bool
// }

// // * CreateUser is the application service type that creates a user
// type CreateUserUsecase = ddd.ApplicationService[CreateUserDTO, string]

// // * RegisterUserUsecase implements the CreateUser application service
// type CreateUserDeps struct {
// 	UserRepo domain.UserRepository
// 	// AuthRepo AuthDomain.AuthRepository
// }

// // NewRegisterUserUsecase creates a new instance of RegisterUserUsecase
// func NewRegisterUserUsecase(userRepo domain.UserRepository) *CreateUserDeps {
// 	return &CreateUserDeps{UserRepo: userRepo}
// }

// // Execute is the method that processes the user creation
// func (uc *CreateUserDeps) Execute(ctx context.Context, payload CreateUserDTO) (uint, error) {
// 	//* Generate user ID and hash the password
// 	id, err := uc.UserRepo.GetNextId(ctx)
// 	if err != nil {
// 		return 0, SharedDomain.NewBusinessError("Failed to generate user ID.", "USER_ID_GENERATION_FAILED")
// 	}

// 	// hashedPassword, err := uc.AuthRepo.Hash(payload.Password)
// 	// if err != nil {
// 	// 	return 0, err
// 	// }

// 	//* Create a new user using the domain logic
// 	user, err := domain.NewUser(domain.UserProps{
// 		Id:        id,
// 		FirstName: payload.FirstName,
// 		LastName:  payload.LastName,
// 		Phone:     payload.Phone,
// 		Email:     payload.Email,
// 		Password:  payload.Password,
// 	})

// 	//* Store the user in the repository
// 	err = uc.UserRepo.Store(ctx, user)
// 	if err != nil {
// 		return 0, err
// 	}

// 	//* Return the user ID as string
// 	return user.Id, nil
// }
