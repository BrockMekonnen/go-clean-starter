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
		validator := validation.NewValidator(validation.ValidationSchemas{
			Body:   &createCommentBody{},
			Params: &createCommentParams{},
		})

		body, err := validator.GetBody(r)
		if err != nil {
			respond.Error(w, err)
			return
		}
		req := body.(*createCommentBody)

		params, err := validator.GetParams(r)
		if err != nil {
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
			PostId: params["postId"],
			Body:   req.Body,
		})

		if err != nil {
			respond.Error(w, err)
			return
		}

		respond.SuccessWithData(w, http.StatusCreated, result)
	}
}
