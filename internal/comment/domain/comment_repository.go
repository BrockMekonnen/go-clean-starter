package domain

import (
	"context"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/contracts"
)

type CommentRepository interface {
	contracts.Repository[Comment]

	FindById(ctx context.Context, id string) (*Comment, error)
}
