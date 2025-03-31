package cqrs

import "context"

type Sort struct {
	Field     string `json:"field"`
	Direction string `json:"direction"` // "asc" or "desc"
}

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type ResultPage struct {
	Current      int  `json:"current"`
	PageSize     int  `json:"page_size"`
	TotalPages   int  `json:"total_pages"`
	TotalElements int  `json:"total_elements"`
	First        bool `json:"first"`
	Last         bool `json:"last"`
}

type Query[F any] struct {
	Filter F `json:"filter"`
}

type PaginatedQuery[F any] struct {
	Query[F]
	Pagination Pagination `json:"pagination"`
}

type SortedQuery[F any] struct {
	Query[F]
	Sort []Sort `json:"sort"`
}

type SortedPaginatedQuery[F any] struct {
	Query[F]
	Sort       []Sort      `json:"sort"`
	Pagination Pagination `json:"pagination"`
}

type QueryResult[T any] struct {
	Data T `json:"data"`
}

type PaginatedQueryResult[T any] struct {
	QueryResult[T]
	Page ResultPage `json:"page"`
}

// QueryHandler with context support
type QueryHandler[P any, R any] interface {
	Handle(ctx context.Context, payload P) (R, error)
}