package infrastructure

import (
	"context"
	"errors"

	cErrors "github.com/BrockMekonnen/go-clean-starter/core/lib/errors"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/domain"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &UserRepository{db: db}
}

// GetNextId generates and returns the next available ID for a user
func (r *UserRepository) GetNextId(ctx context.Context) (uint, error) {
	var maxId uint
	err := r.db.WithContext(ctx).Model(&User{}).Select("COALESCE(MAX(id), 0)").Scan(&maxId).Error
	if err != nil {
		return 0, err
	}
	return maxId + 1, nil
}

func (r *UserRepository) Store(ctx context.Context, user *domain.User) error {
	userData := ToData(*user)

	return r.db.WithContext(ctx).Create(&userData).Error
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	userData := ToData(*user)

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

func (r *UserRepository) FindById(ctx context.Context, id uint) (*domain.User, error) {
	var userSchema User

	// Use the context in the database query
	result := r.db.WithContext(ctx).First(&userSchema, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, cErrors.NotFoundError{}
		}
		return nil, result.Error
	}

	user := ToEntity(userSchema)
	return &user, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id uint) error {
	var userSchema User

	// Use the context in the database query
	result := r.db.WithContext(ctx).Delete(&userSchema, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return cErrors.NotFoundError{}
		}
		return result.Error
	}

	// Check if the user was actually deleted
	if result.RowsAffected == 0 {
		return cErrors.NotFoundError{}
	}

	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var userSchema User

	// Use the context in the database query
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&userSchema)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, cErrors.NotFoundError{}
		}
		return nil, result.Error
	}

	user := ToEntity(userSchema)
	return &user, nil
}

func (r *UserRepository) FindByPhone(ctx context.Context, phone string) (*domain.User, error) {
	var userSchema User

	// Use the context in the database query
	result := r.db.WithContext(ctx).Where("phone = ?", phone).First(&userSchema)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, cErrors.NotFoundError{}
		}
		return nil, result.Error
	}

	user := ToEntity(userSchema)
	return &user, nil
}
