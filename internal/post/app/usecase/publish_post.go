package usecase

import (
	"context"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/contracts"
	sharedDomain "github.com/BrockMekonnen/go-clean-starter/internal/_shared/domain"
	"github.com/BrockMekonnen/go-clean-starter/internal/post/domain"
)

type PublishPostParams struct {
	ID     string
	UserId string
}

type PublishPostUsecase = contracts.ApplicationService[PublishPostParams, contracts.Void]

func MakePublishPostUsecase(
	postRepo domain.PostRepository,
) PublishPostUsecase {
	return func(ctx context.Context, params PublishPostParams) (contracts.Void, error) {
		post, err := postRepo.FindById(ctx, params.ID)
		if err != nil {
			return contracts.Void{}, err
		}

		if post.UserId != params.UserId {
			return contracts.Void{}, sharedDomain.NewBusinessError("Don't have permission to update this post!", "")
		}

		if post.IsPublished() {
			return contracts.Void{}, sharedDomain.NewBusinessError("Can't Republish a published Post!", "")
		}

		post, err = post.PublishPost()
		if err != nil {
			return contracts.Void{}, err
		}

		if err = postRepo.Update(ctx, post); err != nil {
			return contracts.Void{}, err
		}

		return contracts.Void{}, nil
	}
}
