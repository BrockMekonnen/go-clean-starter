package handlers

import (
	"net/http"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/extension"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/respond"
	"github.com/BrockMekonnen/go-clean-starter/internal/post/app/usecase"
	"github.com/gorilla/mux"
)

func DeletePostHandler(
	deletePost usecase.DeletePostUsecase,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		postId := vars["id"]

		authCtx, err := extension.GetAuthContextData(r.Context())
		if err != nil {
			respond.Error(w, err)
			return
		}

		_, err = deletePost(r.Context(), usecase.DeletePostParams{
			UserId: authCtx.Credentials.UID,
			ID:     postId,
		})

		if err != nil {
			respond.Error(w, err)
			return
		}

		respond.Success(w, http.StatusNoContent, nil)
	}
}
