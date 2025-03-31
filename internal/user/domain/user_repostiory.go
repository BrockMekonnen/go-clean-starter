package domain

import (
	"context"

	"github.com/BrockMekonnen/go-clean-starter/internal/_lib/ddd"
)

// UserRepository is the interface for interacting with users in the system.
type UserRepository interface {
	ddd.Repository[User]

	// GetNextId generates and returns the next available ID for a user
	GetNextId(ctx context.Context) (uint, error)

	// GetUserById retrieves a user by their ID
	GetUserById(ctx context.Context, id uint) (*User, error)

	// DeleteUser removes a user by their ID
	DeleteUser(ctx context.Context, id uint) error

	// FindByEmail retrieves a user by their email
	FindByEmail(ctx context.Context, email string) (*User, error)

	// FindByPhoneNumber retrieves a user by their phone number
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (*User, error)
}