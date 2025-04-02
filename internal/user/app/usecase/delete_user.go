package usecase

import (
	"context"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/ddd"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/domain"
)

// DeleteUserContract makes the function signature readable
type DeleteUserUsecase = ddd.ApplicationService[uint, ddd.Void]

// DeleteUserDeps declares all required dependencies
type DeleteUserDeps struct {
	UserRepo domain.UserRepository
}

// NewDeleteUserUsecase implements the contract
func NewDeleteUserUsecase(deps DeleteUserDeps) DeleteUserUsecase {
	return func(ctx context.Context, id uint) (ddd.Void, error) {
		err := deps.UserRepo.DeleteUser(ctx, id)
		if err != nil {
			return ddd.Void{}, err
		}
		return ddd.Void{}, nil
	}
}