package domain

import (
	"errors"
	"time"
)

// User represents a user aggregate in the system
type User struct {
	Id        uint
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

// NewUser creates a new user with invariants
func NewUser(props UserProps) (*User, error) {
	// Invariant checks
	if len(props.FirstName) == 0 {
		return nil, errors.New("first name cannot be empty")
	}
	if len(props.LastName) == 0 {
		return nil, errors.New("last name cannot be empty")
	}
	if len(props.Phone) == 0 {
		return nil, errors.New("phone cannot be empty")
	}
	if len(props.Email) == 0 {
		return nil, errors.New("email cannot be empty")
	}
	if len(props.Password) == 0 {
		return nil, errors.New("password cannot be empty")
	}

	user := &User{
		Id:        props.Id,
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

// UserProps defines the properties required to create a new User
type UserProps struct {
	Id        uint
	FirstName string
	LastName  string
	Phone     string
	Email     string
	Password  string
}
