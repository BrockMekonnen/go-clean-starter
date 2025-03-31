package usecase

import (
	"context"
	"github.com/BrockMekonnen/go-clean-starter/internal/_lib/ddd"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/domain"
	
)

type DeleteUser = ddd.ApplicationService[uint, ddd.Void]

type DeleteUserUsecase struct {
	UserRepo domain.UserRepository // Interface type (not pointer to interface)
}

func NewDeleteUserUsecase(userRepo domain.UserRepository) *DeleteUserUsecase {
	return &DeleteUserUsecase{UserRepo: userRepo}
}

func (uc *DeleteUserUsecase) Execute(ctx context.Context, id uint) (ddd.Void, error) {
	// Call the DeleteUser method on the UserRepo interface
	err := uc.UserRepo.DeleteUser(ctx, id)
	if err != nil {
		return ddd.Void{}, err
	}
	return ddd.Void{}, nil
}