package handlers

import (
	"net/http"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/respond"
	"github.com/BrockMekonnen/go-clean-starter/internal/post/app/query"
	"github.com/gorilla/mux"
)

func FindPostByIdHandler(
	findPostById query.FindPostById,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		result, err := findPostById.Execute(r.Context(), id)
		if err != nil {
			respond.Error(w, err)
			return
		}

		respond.Success(w, http.StatusOK, result)
	}
}
