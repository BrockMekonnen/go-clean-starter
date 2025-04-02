package usecase

import (
	"context"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/ddd"
	sharedDomain "github.com/BrockMekonnen/go-clean-starter/internal/_shared/domain"
	"github.com/BrockMekonnen/go-clean-starter/internal/auth/domain"
)

type VerifyTokenDeps struct {
	AuthRepository domain.AuthRepository
}

type VerifyResponse struct {
	Uid   uint
	Scope []string
}

type VerifyTokenUsecase = ddd.ApplicationService[string, VerifyResponse]

func NewVerifyTokenUsecase(deps VerifyTokenDeps) VerifyTokenUsecase {
	return func(ctx context.Context, token string) (VerifyResponse, error) {
		decoded, err := deps.AuthRepository.Decode(ctx, token)

		if err != nil || decoded == nil {
			return VerifyResponse{}, sharedDomain.NewBusinessError("Authentication Error Invalid Token", "Unauthorized")
		}

		return VerifyResponse{
			Uid:   decoded.Uid,
			Scope: decoded.Scope,
		}, err
	}
}
