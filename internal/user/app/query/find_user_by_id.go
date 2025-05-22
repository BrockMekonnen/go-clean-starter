package query

import (
	"time"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/contracts"
)

type FindUserByIdDTO struct {
	Id        string    `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Roles     any       `json:"roles"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type FindUserByIdResult = contracts.QueryResult[FindUserByIdDTO]

type FindUserById = contracts.QueryHandler[string, FindUserByIdResult]
