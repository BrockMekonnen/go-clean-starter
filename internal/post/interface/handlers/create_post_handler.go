package handlers

import (
	"net/http"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/extension"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/respond"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/validation"
	"github.com/BrockMekonnen/go-clean-starter/internal/post/app/query"
	"github.com/BrockMekonnen/go-clean-starter/internal/post/app/usecase"
)

type CreatePostBody struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

func CreatePostHandler(
	createPost usecase.CreatePostUsecase,
	findPostById query.FindPostById,
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

		createdId, err := createPost(r.Context(), usecase.CreatePostParams{
			UserId:  authCtx.Credentials.UID,
			Title:   req.Title,
			Content: req.Content,
		})

		if err != nil {
			respond.Error(w, err)
			return
		}

		result, err := findPostById.Execute(r.Context(), createdId)
		if err != nil {
			respond.Error(w, err)
			return
		}

		respond.Success(w, http.StatusCreated, result)
	}
}
