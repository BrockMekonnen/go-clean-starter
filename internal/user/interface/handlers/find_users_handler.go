package handlers

import (
	"net/http"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/contracts"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/respond"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/validation"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/app/query"
)

func FindUsersHandler(
	findUsers query.FindUsers,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Setup paginator with defaults
		paginator := validation.NewPaginator()

		// Get pagination data
		page, pageSize, err := paginator.GetPagination(r)
		if err != nil {
			respond.Error(w, err)
			return
		}

		// Build the query
		params := query.FindUsersQuery{
			Filter:     contracts.Void{},
			Pagination: contracts.Pagination{Page: page, PageSize: pageSize},
		}

		// Call the use case handler
		result, err := findUsers.Execute(r.Context(), params)
		if err != nil {
			respond.Error(w, err)
			return
		}

		// Return the result
		respond.Success(w, http.StatusOK, result)
	}
}
