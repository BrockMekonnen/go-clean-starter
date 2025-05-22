package handlers

import (
	"net/http"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/contracts"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/respond"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/validation"
	"github.com/BrockMekonnen/go-clean-starter/internal/post/app/query"
)

func FindPostsHandler(
	findPosts query.FindPosts,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		paginator := validation.NewPaginator()

		// Get pagination data
		page, pageSize, err := paginator.GetPagination(r)
		if err != nil {
			respond.Error(w, err)
			return
		}

		params := query.FindPostsQuery{
			Filter:     contracts.Void{},
			Pagination: contracts.Pagination{Page: page, PageSize: pageSize},
		}

		result, err := findPosts.Execute(r.Context(), params)
		if err != nil {
			respond.Error(w, err)
			return
		}

		respond.Success(w, http.StatusOK, result)
	}

}
