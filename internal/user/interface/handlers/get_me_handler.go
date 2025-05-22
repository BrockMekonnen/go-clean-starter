package handlers

import (
	"net/http"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/extension"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/respond"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/app/query"
)

func GetMeHandler(
	findUserById query.FindUserById,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve AuthContext from context
		authCtx, err := extension.GetAuthContextData(r.Context())
		if err != nil {
			respond.Error(w, err)
			return
		}

		// Extract UID
		userID := authCtx.Credentials.UID

		// Call query to find user
		result, err := findUserById.Execute(r.Context(), userID)
		if err != nil {
			respond.Error(w, err)
			return
		}

		// Return user data
		respond.Success(w, http.StatusOK, result)
	}
}
