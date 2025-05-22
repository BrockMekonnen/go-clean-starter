package infrastructure

import (
	"context"
	"errors"

	customErrors "github.com/BrockMekonnen/go-clean-starter/core/lib/errors"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/hashids"
	"github.com/BrockMekonnen/go-clean-starter/internal/post/domain"

	"gorm.io/gorm"
)

type PostRepository struct {
	db      *gorm.DB
	hashIDs hashids.HashID
}

func MakePostRepository(db *gorm.DB, hashIDs hashids.HashID) domain.PostRepository {
	return &PostRepository{db: db, hashIDs: hashIDs}
}

func (r *PostRepository) GetNextId(ctx context.Context) (string, error) {
	var maxId uint
	err := r.db.WithContext(ctx).Model(&Post{}).Select("COALESCE(MAX(id), 0)").Scan(&maxId).Error
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

func (r *PostRepository) FindById(ctx context.Context, idStr string) (*domain.Post, error) {
	var postSchema Post

	id, err := r.hashIDs.DecodeID(idStr)
	if err != nil {
		return &domain.Post{}, customErrors.NewBadRequestError[any](err.Error(), "", nil)
	}

	result := r.db.WithContext(ctx).First(&postSchema, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, customErrors.NewNotFoundError("", "", result)
		}
		return nil, result.Error
	}

	user, err := ToEntity(postSchema)
	if err != nil {
		return &domain.Post{}, customErrors.NewBadRequestError[any](err.Error(), "", nil)
	}

	return &user, nil
}

func (r *PostRepository) Store(ctx context.Context, entity *domain.Post) error {
	data, err := ToData(*entity)
	if err != nil {
		return customErrors.NewBadRequestError[any](err.Error(), "", nil)
	}

	return r.db.WithContext(ctx).Create(&data).Error
}

func (r *PostRepository) Update(ctx context.Context, post *domain.Post) error {
	data, err := ToData(*post)
	if err != nil {
		return customErrors.NewBadRequestError[any](err.Error(), "", nil)
	}

	result := r.db.WithContext(ctx).Model(&Post{}).
		Where("id = ? AND version = ?", data.ID, data.Version).
		Updates(data)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("optimistic lock failed or user not found")
	}

	return nil
}
