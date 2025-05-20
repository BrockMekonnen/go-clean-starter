package delivery

import (
	"net/http"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/respond"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/app/usecase"
	"github.com/gorilla/mux"
)

type DeleteUserHandlerDeps struct {
	DeleteUser usecase.DeleteUserUsecase
}

func NewDeleteUserHandler(deps DeleteUserHandlerDeps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract and validate ID
		vars := mux.Vars(r)
		userID := vars["id"]

		// Execute use case
		_, err := deps.DeleteUser(r.Context(), userID)
		if err != nil {
			respond.Error(w, err)
			return
		}

		// w.WriteHeader(http.StatusNoContent)
		respond.SuccessWithData(w, http.StatusNoContent, nil)
	}
}
