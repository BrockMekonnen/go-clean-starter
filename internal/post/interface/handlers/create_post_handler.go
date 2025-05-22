package handlers

import (
	"net/http"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/extension"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/respond"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/validation"
	"github.com/BrockMekonnen/go-clean-starter/internal/post/app/usecase"
)

type CreatePostBody struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

func CreatePostHandler(
	createPost usecase.CreatePostUsecase,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validator := validation.NewValidator(validation.ValidationSchemas{
			Body: &CreatePostBody{},
		})

		body, err := validator.GetBody(r)
		if err != nil {
			respond.Error(w, err)
			return
		}
		req := body.(*CreatePostBody)

		authCtx, err := extension.GetAuthContextData(r.Context())
		if err != nil {
			respond.Error(w, err)
			return
		}

		result, err := createPost(r.Context(), usecase.CreatePostParams{
			UserId:  authCtx.Credentials.UID,
			Title:   req.Title,
			Content: req.Content,
		})

		if err != nil {
			respond.Error(w, err)
			return
		}

		respond.SuccessWithData(w, http.StatusCreated, result)
	}
}
