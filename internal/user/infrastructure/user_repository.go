package infrastructure

import (
	"context"
	"errors"

	customErrors "github.com/BrockMekonnen/go-clean-starter/core/lib/errors"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/hashids"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/domain"
	"gorm.io/gorm"
)

type UserRepository struct {
	db      *gorm.DB
	hashIDs hashids.HashID
}

func MakeUserRepository(db *gorm.DB, hashIDs hashids.HashID) domain.UserRepository {
	return &UserRepository{db: db, hashIDs: hashIDs}
}

// GetNextId generates and returns the next available ID for a user
func (r *UserRepository) GetNextId(ctx context.Context) (string, error) {
	var maxId uint
	err := r.db.WithContext(ctx).Model(&User{}).Select("COALESCE(MAX(id), 0)").Scan(&maxId).Error
	if err != nil {
		return "", err
	}
	nextId := maxId + 1
	hashedId, err := r.hashIDs.EncodeID(nextId)
	if err != nil {
		return "", customErrors.NewBadRequestError[any](err.Error(), "", nil)
	}
	return hashedId, nil
}

func (r *UserRepository) Store(ctx context.Context, user *domain.User) error {
	userData, err := ToData(*user)

	if err != nil {
		return customErrors.NewBadRequestError[any](err.Error(), "", nil)
	}

	return r.db.WithContext(ctx).Create(&userData).Error
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	userData, err := ToData(*user)
	if err != nil {
		return customErrors.NewBadRequestError[any](err.Error(), "", nil)
	}

	result := r.db.WithContext(ctx).Model(&User{}).
		Where("id = ? AND version = ?", userData.Id, userData.Version).
		Updates(userData)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("optimistic lock failed or user not found")
	}
	return nil
}

func (r *UserRepository) FindById(ctx context.Context, idStr string) (*domain.User, error) {
	var userSchema User

	id, err := r.hashIDs.DecodeID(idStr)
	if err != nil {
		return &domain.User{}, customErrors.NewBadRequestError[any](err.Error(), "", nil)
	}

	result := r.db.WithContext(ctx).First(&userSchema, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, customErrors.NewNotFoundError("", "", result)
		}
		return nil, result.Error
	}

	user, err := ToEntity(userSchema)
	if err != nil {
		return &domain.User{}, customErrors.NewBadRequestError[any](err.Error(), "", nil)
	}

	return &user, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, idStr string) error {
	var userSchema User

	id, err := r.hashIDs.DecodeID(idStr)
	if err != nil {
		return customErrors.NewBadRequestError[any](err.Error(), "", nil)
	}

	// Use the context in the database query
	result := r.db.WithContext(ctx).Delete(&userSchema, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return customErrors.NewNotFoundError("", "", result)
		}
		return result.Error
	}

	// Check if the user was actually deleted
	if result.RowsAffected == 0 {
		return customErrors.NewNotFoundError[any]("", "", nil)
	}

	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var userSchema User

	// Use the context in the database query
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&userSchema)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, customErrors.NewNotFoundError("", "", result)
		}
		return nil, result.Error
	}

	user, err := ToEntity(userSchema)
	if err != nil {
		return &domain.User{}, customErrors.NewBadRequestError[any](err.Error(), "", nil)
	}

	return &user, nil
}

func (r *UserRepository) FindByPhone(ctx context.Context, phone string) (*domain.User, error) {
	var userSchema User

	// Use the context in the database query
	result := r.db.WithContext(ctx).Where("phone = ?", phone).First(&userSchema)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, customErrors.NewNotFoundError("", "", result)
		}
		return nil, result.Error
	}

	user, err := ToEntity(userSchema)
	if err != nil {
		return &domain.User{}, customErrors.NewBadRequestError[any](err.Error(), "", nil)
	}

	return &user, nil
}
