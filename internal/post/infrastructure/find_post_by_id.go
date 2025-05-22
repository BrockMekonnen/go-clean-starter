package infrastructure

import (
	"context"
	"time"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/hashids"
	"github.com/BrockMekonnen/go-clean-starter/internal/post/app/query"
	"gorm.io/gorm"
)

type FindPostByIdDeps struct {
	db      *gorm.DB
	hashIDs hashids.HashID
}

func MakeFindPostById(db *gorm.DB, hashIDs hashids.HashID) query.FindPostById {
	return &FindPostByIdDeps{db: db, hashIDs: hashIDs}
}

func (f *FindPostByIdDeps) Execute(ctx context.Context, postIdStr string) (query.FindPostByIdResult, error) {
	postId, err := f.hashIDs.DecodeID(postIdStr)
	if err != nil {
		return query.FindPostByIdResult{}, err
	}

	var post Post
	err = f.db.WithContext(ctx).
		Preload("User").
		Preload("Comments").
		First(&post, "id = ?", postId).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return query.FindPostByIdResult{}, nil
		}
		return query.FindPostByIdResult{}, err
	}

	comments := make([]query.FindPostCommentDTO, len(post.Comments))
	for i, c := range post.Comments {
		hashedCommentID, _ := f.hashIDs.EncodeID(c.ID)
		comments[i] = query.FindPostCommentDTO{
			Id:        hashedCommentID,
			Body:      c.Body,
			CreatedAt: c.CreatedAt.Format(time.RFC3339),
		}
	}

	hashedUserID, _ := f.hashIDs.EncodeID(post.UserId)
	hashedPostID, _ := f.hashIDs.EncodeID(post.ID)

	result := query.FindPostByIdDTO{
		Id:      hashedPostID,
		Title:   post.Title,
		Content: post.Content,
		User: query.UserDTO{
			Id:        hashedUserID,
			FirstName: post.User.FirstName,
			LastName:  post.User.LastName,
		},
		State:     post.State,
		PostedAt:  post.PublishedAt,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
		Comments:  comments,
	}

	return query.FindPostByIdResult{Data: result}, nil
}
