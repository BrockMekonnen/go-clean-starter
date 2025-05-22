package usecase

import (
	"context"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/contracts"
	sharedDomain "github.com/BrockMekonnen/go-clean-starter/internal/_shared/domain"
	"github.com/BrockMekonnen/go-clean-starter/internal/post/domain"
)

type CreatePostParams struct {
	Title   string
	Content string
	UserId  string
}

type CreatePostUsecase = contracts.ApplicationService[CreatePostParams, string]

func MakeCreatePostUsecase(
	postRepo domain.PostRepository,
) CreatePostUsecase {
	return func(ctx context.Context, params CreatePostParams) (string, error) {
		id, err := postRepo.GetNextId(ctx)
		if err != nil {
			return "", sharedDomain.NewBusinessError("Error occurred while generating post ID", "")
		}

		post, err := domain.NewPost(domain.PostProps{
			ID:      id,
			Title:   params.Title,
			Content: params.Content,
			UserId:  params.UserId,
		})

		if err != nil {
			return "", err
		}

		if err := postRepo.Store(ctx, post); err != nil {
			return "", sharedDomain.NewBusinessError("Error occurred while creating new post", "")
		}

		return post.ID, nil
	}
}
