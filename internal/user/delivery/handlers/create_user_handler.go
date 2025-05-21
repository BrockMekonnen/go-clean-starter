package delivery

import (
	"net/http"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/respond"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/validation"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/app/usecase"
)

type CreateUserRequest struct {
	FirstName                string `json:"firstName" validate:"required"`
	LastName                 string `json:"lastName" validate:"required"`
	Phone                    string `json:"phone" validate:"required"`
	Email                    string `json:"email" validate:"required,email"`
	Password                 string `json:"password" validate:"required,min=8"`
	IsTermAndConditionAgreed bool   `json:"isTermAndConditionAgreed" validate:"required"`
}

// NewCreateUserHandler creates the handler with explicit dependencies
func CreateUserHandler(
	createUser usecase.CreateUserUsecase,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set up validator with request body schema
		validator := validation.NewValidator(validation.ValidationSchemas{
			Body: &CreateUserRequest{},
		})

		// Get and validate request body
		body, err := validator.GetBody(r)
		if err != nil {
			respond.Error(w, err)
			return
		}

		req := body.(*CreateUserRequest)

		// Execute use case
		userID, err := createUser(r.Context(), usecase.CreateUserParams{
			FirstName:                req.FirstName,
			LastName:                 req.LastName,
			Phone:                    req.Phone,
			Email:                    req.Email,
			Password:                 req.Password,
			IsTermAndConditionAgreed: req.IsTermAndConditionAgreed,
		})

		if err != nil {
			respond.Error(w, err)
			return
		}

		respond.SuccessWithData(w, http.StatusCreated, userID)
	}
}
