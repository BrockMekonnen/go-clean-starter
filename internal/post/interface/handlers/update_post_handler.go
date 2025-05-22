package handlers

import (
	"net/http"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/extension"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/respond"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/validation"
	"github.com/BrockMekonnen/go-clean-starter/internal/post/app/usecase"
)

type UpdatePostBody struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func UpdatePostHandler(
	updatePost usecase.UpdatePostUsecase,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validator := validation.NewValidator(validation.ValidationSchemas{
			Body: &UpdatePostBody{},
		})

		body, err := validator.GetBody(r)
		if err != nil {
			respond.Error(w, err)
		}

		req := body.(*UpdatePostBody)

		authCtx, err := extension.GetAuthContextData(r.Context())
		if err != nil {
			respond.Error(w, err)
			return
		}

		result, err := updatePost(r.Context(), usecase.UpdatePostParams{
			UserId:  authCtx.Credentials.UID,
			Title:   req.Title,
			Content: req.Content,
		})

		if err != nil {
			respond.Error(w, err)
		}

		respond.SuccessWithData(w, http.StatusOK, result)
	}
}
