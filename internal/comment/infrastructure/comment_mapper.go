package infrastructure

import (
	"github.com/BrockMekonnen/go-clean-starter/core/di"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/contracts"
	"github.com/BrockMekonnen/go-clean-starter/internal/comment/domain"
)

type CommentMapper struct{}

var _ contracts.DataMapper[domain.Comment, Comment] = (*CommentMapper)(nil)

func (c *CommentMapper) ToData(entity domain.Comment) (Comment, error) {
	return ToData(entity)
}

func ToData(entity domain.Comment) (Comment, error) {
	hashids := di.GetHashID()
	id, err := hashids.DecodeID(entity.ID)
	if err != nil {
		return Comment{}, err
	}

	userId, err := hashids.DecodeID(entity.UserId)
	if err != nil {
		return Comment{}, err
	}

	postId, err := hashids.DecodeID(entity.PostId)
	if err != nil {
		return Comment{}, err
	}
	return Comment{
		ID:        id,
		UserId:    userId,
		PostId:    postId,
		Body:      entity.Body,
		Status:    string(entity.Status),
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		Version:   entity.Version,
		Deleted:   entity.Status == "DELETED",
	}, nil
}

// ToEntity implements contracts.DataMapper.
func (c *CommentMapper) ToEntity(data Comment) (domain.Comment, error) {
	return ToEntity(data)
}

func ToEntity(data Comment) (domain.Comment, error) {
	hashids := di.GetHashID()
	hashedId, err := hashids.EncodeID(data.ID)
	if err != nil {
		return domain.Comment{}, err
	}

	hashedUserId, err := hashids.EncodeID(data.UserId)
	if err != nil {
		return domain.Comment{}, err
	}

	hashedPostId, err := hashids.EncodeID(data.PostId)
	if err != nil {
		return domain.Comment{}, err
	}

	return domain.Comment{
		ID:        hashedId,
		Body:      data.Body,
		UserId:    hashedUserId,
		PostId:    hashedPostId,
		Status:    domain.Status(data.Status),
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		Version:   data.Version,
	}, nil
}
