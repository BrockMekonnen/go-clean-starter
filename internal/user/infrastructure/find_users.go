package infrastructure

import (
	"context"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/contracts"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/hashids"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/app/query"
	"gorm.io/gorm"
)

type FindUsersDeps struct {
	db      *gorm.DB
	hashIDs hashids.HashID
}

// NewFindUsersHandler
func MakeFindUsers(db *gorm.DB, hashIDs hashids.HashID) query.FindUsers {
	return &FindUsersDeps{db: db, hashIDs: hashIDs}
}

func (h *FindUsersDeps) Handle(ctx context.Context, params query.FindUsersQuery) (query.FindUsersResult, error) {
	var users []User

	usersData := h.db.WithContext(ctx).Model(&User{})

	// Count total users
	var totalElements int64
	if err := usersData.Count(&totalElements).Error; err != nil {
		return query.FindUsersResult{}, err
	}

	// Apply pagination
	page := params.Pagination.Page
	pageSize := params.Pagination.PageSize
	offset := (page - 1) * pageSize

	if err := usersData.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return query.FindUsersResult{}, err
	}

	// Convert to DTOs
	convertedUsers := make([]query.FindUsersDTO, len(users))
	for i, user := range users {
		hashedId, _ := h.hashIDs.EncodeID(user.Id) // ignoring error for now
		convertedUsers[i] = query.FindUsersDTO{
			Id:        hashedId,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Phone:     user.Phone,
			Email:     user.Email,
			Roles:     user.Roles,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}

	// Calculate total pages
	totalPages := int((totalElements + int64(pageSize) - 1) / int64(pageSize)) // ceil

	// Prepare the result
	result := query.FindUsersResult{
		Data: convertedUsers,
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
