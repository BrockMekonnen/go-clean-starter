package infrastructure

import (
	"encoding/json"
	"fmt"

	"github.com/BrockMekonnen/go-clean-starter/core/di"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/contracts"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/domain"
	"github.com/jackc/pgtype"
)

type UserMapper struct{}

var _ contracts.DataMapper[domain.User, User] = (*UserMapper)(nil)

func (m *UserMapper) ToEntity(schema User) (domain.User, error) {
	return ToEntity(schema)
}

func (m *UserMapper) ToData(user domain.User) (User, error) {
	return ToData(user)
}

func ToEntity(schema User) (domain.User, error) {
	hashids := di.GetHashID()
	hashedId, err := hashids.EncodeID(schema.Id)
	if err != nil {
		return domain.User{}, err
	}

	var roles []string
	if schema.Roles.Status == pgtype.Present {
		err := json.Unmarshal(schema.Roles.Bytes, &roles)
		if err != nil {
			return domain.User{}, err
		}
	}

	return domain.User{
		ID:        hashedId,
		FirstName: schema.FirstName,
		LastName:  schema.LastName,
		Phone:     schema.Phone,
		Email:     schema.Email,
		Password:  schema.Password,
		Roles:     roles,
		CreatedAt: schema.CreatedAt,
		UpdatedAt: schema.UpdatedAt,
		Version:   schema.Version,
	}, nil
}

func ToData(user domain.User) (User, error) {
	hashids := di.GetHashID()
	id, err := hashids.DecodeID(user.ID)

	if err != nil {
		return User{}, err
	}

	jsonb := pgtype.JSONB{}
	err = jsonb.Set(user.Roles)
	if err != nil {
		fmt.Println("Error:UserMapper:ToEntity: ", err)
		return User{}, err
	}

	return User{
		Id:        id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		Email:     user.Email,
		Password:  user.Password,
		Roles:     jsonb,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Version:   user.Version,
	}, nil
}
