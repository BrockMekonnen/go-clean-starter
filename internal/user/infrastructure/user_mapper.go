package infrastructure

import (
	"encoding/json"
	"fmt"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/ddd"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/domain"
	"github.com/jackc/pgtype"
)

type userMapper struct{}

var _ ddd.DataMapper[domain.User, User] = (*userMapper)(nil)

func ToEntity(schema User) domain.User {
	var roles []string
	if schema.Roles.Status == pgtype.Present {
		err := json.Unmarshal(schema.Roles.Bytes, &roles)
		if err != nil {
			fmt.Println("Error:UserMapper:ToEntity: ", err)
		}
	}

	return domain.User{
		Id:        schema.Id,
		FirstName: schema.FirstName,
		LastName:  schema.LastName,
		Phone:     schema.Phone,
		Email:     schema.Email,
		Password:  schema.Password,
		Roles:     roles,
		CreatedAt: schema.CreatedAt,
		UpdatedAt: schema.UpdatedAt,
		Version:   schema.Version,
	}
}

func ToData(user domain.User) User {
	jsonb := pgtype.JSONB{}
	err := jsonb.Set(user.Roles)
	if err != nil {
		fmt.Println("Error:UserMapper:ToEntity: ", err)
	}

	return User{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		Email:     user.Email,
		Password:  user.Password,
		Roles:     jsonb,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Version:   user.Version,
	}
}

// * Interface implementation with error handling
func (m *userMapper) ToEntity(schema User) domain.User {
	return ToEntity(schema)
}

func (m *userMapper) ToData(user domain.User) User {
	return ToData(user)
}
