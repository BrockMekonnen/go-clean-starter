package delivery

import (
	"net/http"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/respond"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/validation"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/app/usecase"
)

type GenerateTokenHandlerDeps struct {
	GenerateToken usecase.GenerateTokenUsecase
}

type GenerateTokenRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func MakeGenerateTokenHandler(deps GenerateTokenHandlerDeps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set up validator with request body schema
		validator := validation.NewValidator(validation.ValidationSchemas{
			Body: &GenerateTokenRequest{},
		})

		// Get and validate request body
		body, err := validator.GetBody(r)
		if err != nil {
			respond.Error(w, err)
			return
		}

		req := body.(*GenerateTokenRequest)

		token, err := deps.GenerateToken(r.Context(), usecase.GenerateTokenParams{
			Email:    req.Email,
			Password: req.Password,
		})

		if err != nil {
			respond.Error(w, err)
			return
		}

		respond.SuccessWithData(w, http.StatusOK, map[string]interface{}{
			"token": token,
		})
	}
}
