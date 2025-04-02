package delivery

import (
	"encoding/json"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/res"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/app/usecase"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type GenerateTokenHandlerDeps struct {
	GenerateToken usecase.GenerateTokenUsecase
}

type GenerateTokenRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func NewGenerateTokenHandler(deps GenerateTokenHandlerDeps) http.HandlerFunc {
	validate := validator.New()

	return func(w http.ResponseWriter, r *http.Request) {
		var req GenerateTokenRequest

		// Bind request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request format", http.StatusBadRequest)
		}

		if err := validate.Struct(req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		token, err := deps.GenerateToken(r.Context(), usecase.GenerateTokenParams{
			Email:    req.Email,
			Password: req.Password,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		respond.Success(w, http.StatusOK, token)
	}
}
