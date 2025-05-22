package infrastructure

import (
	"context"
	"errors"

	customErrors "github.com/BrockMekonnen/go-clean-starter/core/lib/errors"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/hashids"
	"github.com/BrockMekonnen/go-clean-starter/internal/comment/domain"
	"gorm.io/gorm"
)

type CommentRepository struct {
	db      *gorm.DB
	hashIDs hashids.HashID
}

func MakeCommentRepository(db *gorm.DB, hashIDs hashids.HashID) domain.CommentRepository {
	return &CommentRepository{db: db, hashIDs: hashIDs}
}

// GetNextId implements domain.CommentRepository.
func (r *CommentRepository) GetNextId(ctx context.Context) (string, error) {
	var maxId uint
	err := r.db.WithContext(ctx).Model(&Comment{}).Select("COALESCE(MAX(id), 0)").Scan(&maxId).Error
	if err != nil {
		return "", err
	}
	nextId := maxId + 1
	hashedId, err := r.hashIDs.EncodeID(nextId)
	if err != nil {
		return "", customErrors.NewBadRequestError[any](err.Error(), "", nil)
	}
	return hashedId, nil
}

// FindById implements domain.CommentRepository.
func (r *CommentRepository) FindById(ctx context.Context, idStr string) (*domain.Comment, error) {
	var commentSchema Comment

	id, err := r.hashIDs.DecodeID(idStr)
	if err != nil {
		return &domain.Comment{}, customErrors.NewBadRequestError[any](err.Error(), "", nil)
	}

	result := r.db.WithContext(ctx).First(&commentSchema, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, customErrors.NewNotFoundError("", "", result)
		}
		return nil, result.Error
	}

	user, err := ToEntity(commentSchema)
	if err != nil {
		return &domain.Comment{}, customErrors.NewBadRequestError[any](err.Error(), "", nil)
	}

	return &user, nil
}

// Store implements domain.CommentRepository.
func (r *CommentRepository) Store(ctx context.Context, entity *domain.Comment) error {
	data, err := ToData(*entity)
	if err != nil {
		return customErrors.NewBadRequestError[any](err.Error(), "", nil)
	}

	return r.db.WithContext(ctx).Create(&data).Error
}
