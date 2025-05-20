package query

import (
	"time"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/contracts"
)

type FindUsersDTO struct {
	Id        string    `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Roles     any       `json:"roles"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type FindUsersQuery = contracts.PaginatedQuery[contracts.Void]
type FindUsersResult = contracts.PaginatedQueryResult[[]FindUsersDTO]

type FindUsers = contracts.QueryHandler[FindUsersQuery, FindUsersResult]
