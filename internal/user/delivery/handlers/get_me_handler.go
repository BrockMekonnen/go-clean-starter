package delivery

import (
	"errors"
	"net/http"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/extension"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/respond"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/app/query"
)

type GetMeHandlerDeps struct {
	FindUserById query.FindUserById
}

func MakeGetMeHandler(deps GetMeHandlerDeps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve AuthContext from context
		authCtx, ok := r.Context().Value(extension.AuthContextKey).(extension.AuthContext)
		if !ok {
			respond.Error(w, errors.New("unauthorized: missing auth context"))
			return
		}

		// Extract UID
		userID := authCtx.Credentials.UID
		
		// Call query to find user
		result, err := deps.FindUserById.Handle(r.Context(), userID)
		if err != nil {
			respond.Error(w, err)
			return
		}

		// Return user data
		respond.Success(w, http.StatusOK, result)
	}
}
