package usecase

import (
	"context"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/contracts"
	sharedDomain "github.com/BrockMekonnen/go-clean-starter/internal/_shared/domain"
	"github.com/BrockMekonnen/go-clean-starter/internal/post/domain"
)

type UpdatePostParams struct {
	ID      string
	Title   string
	Content string
	UserId  string
}

type UpdatePostUsecase = contracts.ApplicationService[UpdatePostParams, string]

func MakeUpdatePostUsecase(
	postRepo domain.PostRepository,
) UpdatePostUsecase {
	return func(ctx context.Context, params UpdatePostParams) (string, error) {
		post, err := postRepo.FindById(ctx, params.ID)
		if err != nil {
			return "", err
		}

		if post.UserId != params.UserId {
			return "", sharedDomain.NewBusinessError("Don't have permission to edit this post!", "")
		}

		if params.Content != "" {
			if post, err = post.ChangeContent(params.Content); err != nil {
				return "", err
			}
		}

		if params.Title != "" {
			if post, err = post.ChangeTitle(params.Title); err != nil {
				return "", err
			}
		}

		if err = postRepo.Update(ctx, post); err != nil {
			return "", err
		}

		return post.ID, nil
	}
}
