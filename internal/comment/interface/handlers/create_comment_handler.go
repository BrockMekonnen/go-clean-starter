package handlers

import (
	"net/http"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/extension"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/respond"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/validation"
	"github.com/BrockMekonnen/go-clean-starter/internal/comment/app/usecase"
)

type createCommentParams struct {
	PostId string `json:"postId" validate:"required"`
}

type createCommentBody struct {
	Body string `json:"body" validate:"required"`
}

func CreateCommentHandler(
	createComment usecase.CreateCommentUsecase,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validator := validation.NewValidator(validation.ValidationSchemas{})

		body := &createCommentBody{}
		if err := validator.BindAndValidateBody(r, body); err != nil {
			respond.Error(w, err)
			return
		}

		params := &createCommentParams{}
		if err := validator.BindAndValidateParams(r, params); err != nil {
			respond.Error(w, err)
			return
		}

		authCtx, err := extension.GetAuthContextData(r.Context())
		if err != nil {
			respond.Error(w, err)
			return
		}

		result, err := createComment(r.Context(), usecase.CreateCommentParams{
			UserId: authCtx.Credentials.UID,
			PostId: params.PostId,
			Body:   body.Body,
		})

		if err != nil {
			respond.Error(w, err)
			return
		}

		respond.SuccessWithData(w, http.StatusCreated, result)
	}
}