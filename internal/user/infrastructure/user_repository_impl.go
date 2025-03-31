package infrastructure

import (
	"context"
	"errors"

	cErrors "github.com/BrockMekonnen/go-clean-starter/internal/_lib/errors"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/domain"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetNextId generates and returns the next available ID for a user
func (r *UserRepository) GetNextId(ctx context.Context) (uint, error) {
	var maxId uint
	err := r.db.WithContext(ctx).Model(&UserSchema{}).Select("COALESCE(MAX(id), 0)").Scan(&maxId).Error
	if err != nil {
		return 0, err
	}
	return maxId + 1, nil
}

func (r *UserRepository) Store(ctx context.Context, user *domain.User) error {
	userData := ToData(*user)

	// Check if user exists
	var existingUser UserSchema
	result := r.db.WithContext(ctx).First(&existingUser, userData.Id)

	if result.Error == nil {
		// Update existing user with optimistic locking
		updateResult := r.db.WithContext(ctx).Model(&UserSchema{}).
			Where("id = ? AND version = ?", userData.Id, userData.Version).
			Updates(map[string]interface{}{
				"first_name": userData.FirstName,
				"last_name":  userData.LastName,
				"phone":      userData.Phone,
				"email":      userData.Email,
				"password":   userData.Password,
				"roles":      userData.Roles,
				"version":    userData.Version + 1,
			})

		if updateResult.RowsAffected == 0 {
			return errors.New("optimistic lock failed - user was modified by another transaction")
		}
		return nil
	} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Create new user
		return r.db.WithContext(ctx).Create(&userData).Error
	} else {
		return result.Error
	}
}

func (r *UserRepository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	var userSchema UserSchema

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

// DeleteUser removes a user by their ID.
func (r *UserRepository) DeleteUser(ctx context.Context, id uint) error {
	var userSchema UserSchema

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

// FindByEmail retrieves a user by their email.
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var userSchema UserSchema

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

// FindByPhone retrieves a user by their phone number.
func (r *UserRepository) FindByPhone(ctx context.Context, phone string) (*domain.User, error) {
	var userSchema UserSchema

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
