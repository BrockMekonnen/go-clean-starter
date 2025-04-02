package delivery

import (
	"net/http"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/res"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/app/query"
	"github.com/gorilla/mux"
	"strconv"
)

type GetUserHandlerDeps struct {
	FindUserById query.FindUserById
}

type GetUserRequest struct {
	Id uint `json:"id"`
}

func GetUserHandler(deps GetUserHandlerDeps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get ID from URL path
		vars := mux.Vars(r)
		idStr := vars["id"]

		// Convert ID to uint
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// Call query to find user
		result, err := deps.FindUserById.Handle(r.Context(), uint(id))
		if err != nil {
			// Handle specific error types
			// if errors.Is(err, query.ErrUserNotFound) {
			//     http.Error(w, "User not found", http.StatusNotFound)
			//     return
			// }
			http.Error(w, "Failed to find user", http.StatusInternalServerError)
			return
		}

		// Return user data
		respond.Success(w, http.StatusOK, result)
	}
}
