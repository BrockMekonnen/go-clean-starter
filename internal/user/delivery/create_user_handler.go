package delivery

import (
	"encoding/json"
	"net/http"

	respond "github.com/BrockMekonnen/go-clean-starter/core/lib/res"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/app/usecase"
	"github.com/go-playground/validator/v10"
)

type CreateUserHandlerDeps struct {
	CreateUser usecase.CreateUserUsecase
}

type CreateUserRequest struct {
	FirstName                string `json:"firstName" validate:"required"`
	LastName                 string `json:"lastName" validate:"required"`
	Phone                    string `json:"phone" validate:"required"`
	Email                    string `json:"email" validate:"required,email"`
	Password                 string `json:"password" validate:"required,min=8"`
	IsTermAndConditionAgreed bool   `json:"isTermAndConditionAgreed" validate:"required"`
}

// NewCreateUserHandler creates the handler with explicit dependencies
func NewCreateUserHandler(deps CreateUserHandlerDeps) http.HandlerFunc {
	validate := validator.New()

	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateUserRequest

		// Bind request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}

		// Validate request
		if err := validate.Struct(req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Execute use case
		userID, err := deps.CreateUser(r.Context(), usecase.CreateUserParams{
			FirstName:                req.FirstName,
			LastName:                 req.LastName,
			Phone:                    req.Phone,
			Email:                    req.Email,
			Password:                 req.Password,
			IsTermAndConditionAgreed: req.IsTermAndConditionAgreed,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		respond.Success(w, http.StatusCreated, userID)
	}
}
