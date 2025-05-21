package infrastructure

import (
	"context"
	"errors"

	cErrors "github.com/BrockMekonnen/go-clean-starter/core/lib/errors"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/hashids"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/app/query"
	"gorm.io/gorm"
)

type FindUserByIdHandler struct {
	db      *gorm.DB
	hashIDs hashids.HashID
}

// NewFindUserByIdHandler
func MakeFindUserById(db *gorm.DB, hashIDs hashids.HashID) query.FindUserById {
	return &FindUserByIdHandler{db: db, hashIDs: hashIDs}
}

func (h *FindUserByIdHandler) Execute(ctx context.Context, idStr string) (query.FindUserByIdResult, error) {
	var user User

	id, err := h.hashIDs.DecodeID(idStr)
	if err != nil {
		return query.FindUserByIdResult{}, cErrors.NewBadRequestError[any]("Invalid ID format!", "", nil)
	}

	result := h.db.WithContext(ctx).Where("id = ?", id).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return query.FindUserByIdResult{}, cErrors.NewNotFoundError("", "", result)
		}
		return query.FindUserByIdResult{}, result.Error
	}
	hashedId, err := h.hashIDs.EncodeID(uint(user.Id))
	if err != nil {
		return query.FindUserByIdResult{}, err // or wrap the error properly
	}

	userDTO := query.UserDTO{
		Id:        hashedId,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		Email:     user.Email,
		Roles:     user.Roles,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return query.FindUserByIdResult{Data: userDTO}, nil
}
