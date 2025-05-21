package delivery

import (
	"net/http"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/respond"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/app/query"
	"github.com/gorilla/mux"
)

type GetUserHandlerDeps struct {
	FindUserById query.FindUserById
}

type GetUserRequest struct {
	Id uint `json:"id"`
}

func MakeGetUserHandler(deps GetUserHandlerDeps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get ID from URL path
		vars := mux.Vars(r)
		id := vars["id"]

		// Call query to find user
		result, err := deps.FindUserById.Handle(r.Context(), id)
		if err != nil {
			respond.Error(w, err)
			return
		}

		// Return user data
		respond.Success(w, http.StatusOK, result)
	}
}
