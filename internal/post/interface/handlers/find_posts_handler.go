package handlers

import (
	"net/http"
	"time"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/contracts"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/respond"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/validation"
	"github.com/BrockMekonnen/go-clean-starter/internal/post/app/query"
)

type FindPostsRequestQuery struct {
	Title          string    `json:"title"`
	PublishedStart time.Time `json:"publishedStart" validate:"omitempty"`
	PublishedEnd   time.Time `json:"publishedEnd" validate:"omitempty"`
}

func FindPostsHandler(
	findPosts query.FindPosts,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queryData := &FindPostsRequestQuery{}
		validator := validation.NewValidator(validation.ValidationSchemas{})

		paginator := validation.NewPaginator()

		page, pageSize, err := paginator.GetPagination(r)
		if err != nil {
			respond.Error(w, err)
			return
		}

		if err := validator.BindAndValidateQuery(r, queryData); err != nil {
			respond.Error(w, err)
			return
		}

		var publishedRange []time.Time
		if !queryData.PublishedStart.IsZero() && !queryData.PublishedEnd.IsZero() {
			publishedRange = []time.Time{queryData.PublishedStart, queryData.PublishedEnd}
		}

		params := query.FindPostsQuery{
			Filter: query.FindPostsFilter{
				Title:            queryData.Title,
				PublishedBetween: publishedRange,
				PublishedOnly: true,
			},
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
