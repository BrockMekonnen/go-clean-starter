package domain

import (
	"github.com/BrockMekonnen/go-clean-starter/internal/_shared/domain"
	"time"
)

// User represents a user aggregate in the system
type User struct {
	ID        string
	FirstName string
	LastName  string
	Phone     string
	Email     string
	Password  string
	Roles     []string
	CreatedAt time.Time
	UpdatedAt time.Time
	Version   int
}

// UserProps defines the properties required to create a new User
type UserProps struct {
	ID        string
	FirstName string
	LastName  string
	Phone     string
	Email     string
	Password  string
}

// NewUser creates a new user with invariants
func NewUser(props UserProps) (*User, error) {
	// Invariant checks
	if len(props.FirstName) == 0 {
		return nil, domain.NewBusinessError("first name cannot be empty", "")
	}
	if len(props.LastName) == 0 {
		return nil, domain.NewBusinessError("last name cannot be empty", "")
	}
	if len(props.Phone) == 0 {
		return nil, domain.NewBusinessError("phone cannot be empty", "")
	}
	if len(props.Email) == 0 {
		return nil, domain.NewBusinessError("email cannot be empty", "")
	}
	if len(props.Password) == 0 {
		return nil, domain.NewBusinessError("password cannot be empty", "")
	}

	user := &User{
		ID:        props.ID,
		FirstName: props.FirstName,
		LastName:  props.LastName,
		Phone:     props.Phone,
		Email:     props.Email,
		Password:  props.Password,
		Roles:     []string{"user"}, // Default role
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Version:   0,
	}

	return user, nil
}

// ChangePassword updates the user's password
func (u *User) ChangePassword(newPassword string) *User {
	u.Password = newPassword
	u.UpdatedAt = time.Now() // Update the timestamp
	return u
}
