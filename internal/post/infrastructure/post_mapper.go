package infrastructure

import (
	"github.com/BrockMekonnen/go-clean-starter/core/di"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/contracts"
	"github.com/BrockMekonnen/go-clean-starter/internal/post/domain"
)

type PostMapper struct{}

var _ contracts.DataMapper[domain.Post, Post] = (*PostMapper)(nil)

func (m *PostMapper) ToEntity(data Post) (domain.Post, error) {
	return ToEntity(data)
}

func (m *PostMapper) ToData(entity domain.Post) (Post, error) {
	return ToData(entity)
}

func ToEntity(data Post) (domain.Post, error) {
	hashids := di.GetHashID()

	hashedId, err := hashids.EncodeID(data.ID)
	if err != nil {
		return domain.Post{}, err
	}

	hashedUserId, err := hashids.EncodeID(data.UserId)
	if err != nil {
		return domain.Post{}, err
	}

	return domain.Post{
		ID:          hashedId,
		Title:       data.Title,
		Content:     data.Content,
		UserId:      hashedUserId,
		State:       domain.PostState(data.State),
		PublishedAt: data.PublishedAt,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
		Version:     data.Version,
	}, nil
}

func ToData(entity domain.Post) (Post, error) {
	hashids := di.GetHashID()

	id, err := hashids.DecodeID(entity.ID)
	if err != nil {
		return Post{}, err
	}

	userId, err := hashids.DecodeID(entity.UserId)
	if err != nil {
		return Post{}, err
	}

	return Post{
		ID:          id,
		Title:       entity.Title,
		Content:     entity.Content,
		UserId:      userId,
		State:       string(entity.State),
		PublishedAt: entity.PublishedAt,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		Version:     entity.Version,
		Deleted:     entity.State == "DELETED",
	}, nil
}
