package infrastructure

import (
	"context"
	"errors"

	cErrors "github.com/BrockMekonnen/go-clean-starter/internal/_lib/errors"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/app/query"
	"gorm.io/gorm"
)

type FindUserByIdHandler struct {
	db *gorm.DB
}

// NewFindUserByIdHandler creates a new PostgreSQL implementation
func NewFindUserByIdHandler(db *gorm.DB) query.FindUserById {
	return &FindUserByIdHandler{db: db}
}

func (h *FindUserByIdHandler) Handle(ctx context.Context, id uint) (query.FindUserByIdResult, error) {
	var user UserSchema

	result := h.db.WithContext(ctx).Where("id = ?", id).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return query.FindUserByIdResult{}, cErrors.NotFoundError{}
		}
		return query.FindUserByIdResult{}, result.Error
	}

	userDTO := query.UserDTO{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		Email:     user.Email,
		Roles:     user.Roles,
	}

	return query.FindUserByIdResult{Data: userDTO}, nil
}
