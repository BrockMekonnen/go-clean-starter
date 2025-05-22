package handlers

import (
	"net/http"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/respond"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/app/usecase"
	"github.com/gorilla/mux"
)

func DeleteUserHandler(
	deleteUser usecase.DeleteUserUsecase,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract and validate ID
		vars := mux.Vars(r)
		userID := vars["id"]

		// Execute use case
		_, err := deleteUser(r.Context(), userID)
		if err != nil {
			respond.Error(w, err)
			return
		}

		respond.SuccessWithData(w, http.StatusNoContent, nil)
	}
}
