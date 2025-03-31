package infrastructure

import (
	"github.com/BrockMekonnen/go-clean-starter/internal/_lib/ddd"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/domain"
)

// private type to satisfy ddd.DataMapper interface check
type userMapper struct{}

// Verify interface compliance at compile time
var _ ddd.DataMapper[domain.User, UserSchema] = (*userMapper)(nil)

// ToEntity converts UserSchema to domain.User (package-level function)
func ToEntity(schema UserSchema) domain.User {
	return domain.User{
		Id:        schema.Id,
		FirstName: schema.FirstName,
		LastName:  schema.LastName,
		Phone:     schema.Phone,
		Email:     schema.Email,
		Password:  schema.Password,
		Roles:     schema.Roles,
		CreatedAt: schema.CreatedAt,
		UpdatedAt: schema.UpdatedAt,
		Version:   schema.Version,
	}
}

// ToData converts domain.User to UserSchema (package-level function)
func ToData(user domain.User) UserSchema {
	return UserSchema{
		Id:        uint(user.Id),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		Email:     user.Email,
		Password:  user.Password,
		Roles:     user.Roles,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Version:   user.Version,
	}
}

// Interface implementation methods
func (m *userMapper) ToEntity(schema UserSchema) domain.User {
	return ToEntity(schema)
}

func (m *userMapper) ToData(user domain.User) UserSchema {
	return ToData(user)
}