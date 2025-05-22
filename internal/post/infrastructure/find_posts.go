package infrastructure

import (
	"context"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/contracts"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/hashids"
	"github.com/BrockMekonnen/go-clean-starter/internal/post/app/query"
	"gorm.io/gorm"
)

type FindPostsDeps struct {
	db      *gorm.DB
	hashIDs hashids.HashID
}

func MakeFindPosts(db *gorm.DB, hashIDs hashids.HashID) query.FindPosts {
	return &FindPostsDeps{db: db, hashIDs: hashIDs}
}

// Execute implements contracts.QueryHandler.
func (f *FindPostsDeps) Execute(ctx context.Context, params query.FindPostsQuery) (query.FindPostsResult, error) {
	var posts []Post

	postsData := f.db.WithContext(ctx).Model(&Post{}).Preload("User")

	// Count total posts
	var totalElements int64
	if err := postsData.Count(&totalElements).Error; err != nil {
		return query.FindPostsResult{}, err
	}

	page := params.Pagination.Page
	pageSize := params.Pagination.PageSize
	offset := (page - 1) * pageSize

	if err := postsData.Offset(offset).Limit(pageSize).Find(&posts).Error; err != nil {
		return query.FindPostsResult{}, err
	}

	convertedPosts := make([]query.FindPostsDTO, len(posts))
	for i, post := range posts {
		hashedId, _ := f.hashIDs.EncodeID(post.ID)
		hashedUserId, _ := f.hashIDs.EncodeID(post.UserId)
		convertedPosts[i] = query.FindPostsDTO{
			Id:      hashedId,
			Title:   post.Title,
			Content: post.Content,
			Author: query.UserDTO{
				Id:        hashedUserId,
				FirstName: post.User.FirstName,
				LastName:  post.User.LastName,
			},
			State:    post.State,
			PostedAt: post.PublishedAt,
		}
	}

	totalPages := int((totalElements + int64(pageSize) - 1) / int64(pageSize)) // ceil

	result := query.FindPostsResult{
		Data: convertedPosts,
		Page: contracts.ResultPage{
			Current:       page,
			PageSize:      pageSize,
			TotalPages:    totalPages,
			TotalElements: int(totalElements),
			First:         page == 1,
			Last:          page == totalPages,
		},
	}

	return result, nil
}
