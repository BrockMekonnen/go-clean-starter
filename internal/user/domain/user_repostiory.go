package domain

import (
	"context"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/contracts"
)

// UserRepository is the interface for interacting with users in the system.
type UserRepository interface {
	contracts.Repository[User]

	Update(ctx context.Context, user *User) (error)

	// GetUserById retrieves a user by their ID
	FindById(ctx context.Context, id string) (*User, error)

	// DeleteUser removes a user by their ID
	DeleteUser(ctx context.Context, id string) error

	// FindByEmail retrieves a user by their email
	FindByEmail(ctx context.Context, email string) (*User, error)

	// FindByPhoneNumber retrieves a user by their phone number
	FindByPhone(ctx context.Context, phone string) (*User, error)
}