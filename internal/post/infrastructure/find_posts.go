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

	postsData := f.db.WithContext(ctx).
		Model(&Post{}).
		Where("deleted = ?", false).
		Preload("User")

	// Filter by UserId
	if params.Filter.UserId != "" {
		userId, err := f.hashIDs.DecodeID(params.Filter.UserId)
		if err != nil {
			return query.FindPostsResult{}, err
		}
		postsData = postsData.Where("user_id = ?", userId)
	}

	// Filter by Title (case-insensitive partial match)
	if params.Filter.Title != "" {
		postsData = postsData.Where("LOWER(title) LIKE ?", "%"+params.Filter.Title+"%")
	}

	if params.Filter.PublishedOnly {
		postsData = postsData.Where("state = ?", "PUBLISHED")
	}

	// Filter by PublishedBetween (expects exactly 2 values: start and end)
	if len(params.Filter.PublishedBetween) == 2 {
		start := params.Filter.PublishedBetween[0]
		end := params.Filter.PublishedBetween[1]
		postsData = postsData.Where("posted_at BETWEEN ? AND ?", start, end)
	}

	// Sort by created_at DESC
	postsData = postsData.Order("created_at DESC")

	// Count total
	var totalElements int64
	if err := postsData.Count(&totalElements).Error; err != nil {
		return query.FindPostsResult{}, err
	}

	// Pagination
	page := params.Pagination.Page
	pageSize := params.Pagination.PageSize
	offset := (page - 1) * pageSize

	if err := postsData.Offset(offset).Limit(pageSize).Find(&posts).Error; err != nil {
		return query.FindPostsResult{}, err
	}

	// Transform to DTO
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
			CreatedAt: &post.CreatedAt,
		}
	}

	// Return result
	totalPages := int((totalElements + int64(pageSize) - 1) / int64(pageSize))

	return query.FindPostsResult{
		Data: convertedPosts,
		Page: contracts.ResultPage{
			Current:       page,
			PageSize:      pageSize,
			TotalPages:    totalPages,
			TotalElements: int(totalElements),
			First:         page == 1,
			Last:          page == totalPages,
		},
	}, nil
}
