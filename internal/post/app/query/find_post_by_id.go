package query

import (
	"time"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/contracts"
)

type FindPostCommentDTO struct {
	Id        string `json:"id"`
	Body      string `json:"body"`
	CreatedAt string `json:"createdAt"`
}

type FindPostByIdDTO struct {
	Id        string               `json:"id"`
	Title     string               `json:"title"`
	Content   string               `json:"content"`
	User      UserDTO              `json:"user"`
	State     string               `json:"state"`
	PostedAt  *time.Time           `json:"postedAt"`
	CreatedAt time.Time            `json:"createdAt"`
	UpdatedAt time.Time            `json:"updatedAt"`
	Comments  []FindPostCommentDTO `json:"comments"`
}

type FindPostByIdResult = contracts.QueryResult[FindPostByIdDTO]

type FindPostById = contracts.QueryHandler[string, FindPostByIdResult]
