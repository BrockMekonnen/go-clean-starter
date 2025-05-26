package handlers

import (
	"net/http"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/extension"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/respond"
	"github.com/BrockMekonnen/go-clean-starter/internal/post/app/usecase"
	"github.com/gorilla/mux"
)

func PublishPostHandler(
	publishPost usecase.PublishPostUsecase,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		postId := vars["id"]

		authCtx, err := extension.GetAuthContextData(r.Context())
		if err != nil {
			respond.Error(w, err)
			return
		}

		_, err = publishPost(r.Context(), usecase.PublishPostParams{
			ID:     postId,
			UserId: authCtx.Credentials.UID,
		})

		if err != nil {
			respond.Error(w, err)
		}

		respond.SuccessWithData(w, http.StatusOK, postId)
	}
}
