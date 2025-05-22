package usecase

import (
	"context"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/contracts"
	sharedDomain "github.com/BrockMekonnen/go-clean-starter/internal/_shared/domain"
	"github.com/BrockMekonnen/go-clean-starter/internal/post/domain"
)

type DeletePostParams struct {
	ID     string
	UserId string
}

type DeletePostUsecase = contracts.ApplicationService[DeletePostParams, contracts.Void]

func MakeDeletePostUsecase(
	postRepo domain.PostRepository,
) DeletePostUsecase {
	return func(ctx context.Context, params DeletePostParams) (contracts.Void, error) {
		post, err := postRepo.FindById(ctx, params.ID)
		if err != nil {
			return contracts.Void{}, err
		}

		if post.UserId != params.UserId {
			return contracts.Void{}, sharedDomain.NewBusinessError("Don't have permission to delete this post!", "")
		}

		post, err = post.MarkAsDeleted()
		if err != nil {
			return contracts.Void{}, err
		}

		err = postRepo.Update(ctx, post)
		if err != nil {
			return contracts.Void{}, err
		}

		return contracts.Void{}, nil
	}
}
