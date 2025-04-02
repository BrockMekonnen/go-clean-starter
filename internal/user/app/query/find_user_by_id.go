package query

import (
	"github.com/BrockMekonnen/go-clean-starter/core/lib/cqrs"
)

type UserDTO struct {
	Id              uint  
	FirstName       string  
	LastName        string  
	Phone           string  
	Email           string  
	Roles           any
}


type FindUserByIdResult = cqrs.QueryResult[UserDTO]
type FindUserById = cqrs.QueryHandler[uint, FindUserByIdResult]