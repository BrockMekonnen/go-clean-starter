package query

import (
	"time"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/contracts"
)

type UserDTO struct {
	Id        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type FindPostsDTO struct {
	Id       string     `json:"id"`
	Title    string     `json:"title"`
	Content  string     `json:"content"`
	Author   UserDTO    `json:"user"`
	State    string     `json:"state"`
	PostedAt *time.Time `json:"postedAt"`
	CreatedAt *time.Time `json:"createdAt"`
}

type FindPostsFilter struct {
	UserId           string
	Title            string
	PublishedBetween []time.Time
	PublishedOnly    bool
}

type FindPostsQuery = contracts.PaginatedQuery[FindPostsFilter]
type FindPostsResult = contracts.PaginatedQueryResult[[]FindPostsDTO]

type FindPosts = contracts.QueryHandler[FindPostsQuery, FindPostsResult]
