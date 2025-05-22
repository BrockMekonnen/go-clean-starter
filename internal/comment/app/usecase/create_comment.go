package usecase

import (
	"context"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/contracts"
	sharedDomain "github.com/BrockMekonnen/go-clean-starter/internal/_shared/domain"
	"github.com/BrockMekonnen/go-clean-starter/internal/comment/domain"
)

type CreateCommentParams struct {
	Body   string
	UserId string
	PostId string
}

type CreateCommentUsecase = contracts.ApplicationService[CreateCommentParams, string]

func MakeCreateCommentUsecase(
	commentRepo domain.CommentRepository,
) CreateCommentUsecase {
	return func(ctx context.Context, params CreateCommentParams) (string, error) {
		id, err := commentRepo.GetNextId(ctx)
		if err != nil {
			return "", sharedDomain.NewBusinessError("Error occurred while generating comment ID", "")
		}

		comment := domain.CreateComment(domain.CommentProps{
			ID:     id,
			Body:   params.Body,
			UserId: params.UserId,
			PostId: params.PostId,
		})

		if err := commentRepo.Store(ctx, &comment); err != nil {
			return "", sharedDomain.NewBusinessError("Error occurred while creating new comment!", "")
		}

		return comment.ID, nil
	}
}
