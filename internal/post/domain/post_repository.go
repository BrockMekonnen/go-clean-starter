package domain

import (
	"context"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/contracts"
)

type PostRepository interface {
	contracts.Repository[Post]

	FindById(ctx context.Context, id string) (*Post, error)

	Update(ctx context.Context, article *Post) error
}
