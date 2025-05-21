package usecase

import (
	"context"
	"fmt"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/contracts"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/domain"
)

// DeleteUserContract makes the function signature readable
type DeleteUserUsecase = contracts.ApplicationService[string, contracts.Void]

// DeleteUserDeps declares all required dependencies
type DeleteUserDeps struct {
	UserRepo domain.UserRepository
}

// NewDeleteUserUsecase implements the contract
func MakeDeleteUserUsecase(deps DeleteUserDeps) DeleteUserUsecase {
	return func(ctx context.Context, id string) (contracts.Void, error) {
		err := deps.UserRepo.DeleteUser(ctx, id)
		if err != nil {
			fmt.Printf("DeleteUser error: %v\n", err)
			return contracts.Void{}, err
		}

		return contracts.Void{}, nil
	}
}
