package usecase

import (
	"context"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/contracts"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/domain"
)

// DeleteUserContract makes the function signature readable
type DeleteUserUsecase = contracts.ApplicationService[string, contracts.Void]

// NewDeleteUserUsecase implements the contract
func MakeDeleteUserUsecase(userRepo domain.UserRepository) DeleteUserUsecase {
	return func(ctx context.Context, id string) (contracts.Void, error) {
		err := userRepo.DeleteUser(ctx, id)
		if err != nil {
			return contracts.Void{}, err
		}

		return contracts.Void{}, nil
	}
}
