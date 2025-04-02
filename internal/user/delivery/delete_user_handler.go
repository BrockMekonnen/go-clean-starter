package delivery

import (
	"net/http"
	"strconv"

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
		userID, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil || userID == 0 {
			http.Error(w, "Invalid user ID format", http.StatusBadRequest)
			return
		}

		// Execute use case
		if _, err = deps.DeleteUser(r.Context(), uint(userID)); err != nil {
			panic(err)
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
