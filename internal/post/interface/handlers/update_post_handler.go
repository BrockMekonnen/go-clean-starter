package handlers

import (
	"net/http"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/extension"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/respond"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/validation"
	"github.com/BrockMekonnen/go-clean-starter/internal/post/app/query"
	"github.com/BrockMekonnen/go-clean-starter/internal/post/app/usecase"
	"github.com/gorilla/mux"
)

type updatePostBody struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func UpdatePostHandler(
	updatePost usecase.UpdatePostUsecase,
	findPostById query.FindPostById,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		validator := validation.NewValidator(validation.ValidationSchemas{
			Body: &updatePostBody{},
		})

		body, err := validator.GetBody(r)
		if err != nil {
			respond.Error(w, err)
		}

		req := body.(*updatePostBody)

		authCtx, err := extension.GetAuthContextData(r.Context())
		if err != nil {
			respond.Error(w, err)
			return
		}

		updatedId, err := updatePost(r.Context(), usecase.UpdatePostParams{
			UserId:  authCtx.Credentials.UID,
			Title:   req.Title,
			Content: req.Content,
			ID:      id,
		})

		if err != nil {
			respond.Error(w, err)
			return
		}

		result, err := findPostById.Execute(r.Context(), updatedId)
		if err != nil {
			respond.Error(w, err)
			return
		}

		respond.Success(w, http.StatusOK, result)
	}
}
