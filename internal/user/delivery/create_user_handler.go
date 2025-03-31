package delivery

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/app/usecase"
)

type CreateUserDependencies struct {
	CreateUser usecase.CreateUser
}

type CreateUserRequest struct {
	FirstName                string `json:"firstName" validate:"required"`
	LastName                 string `json:"lastName" validate:"required"`
	Phone                    string `json:"phone" validate:"required"`
	Email                    string `json:"email" validate:"required,email"`
	Password                 string `json:"password" validate:"required,min=8"`
	IsTermAndConditionAgreed bool   `json:"isTermAndConditionAgreed" validate:"required"`
}

type CreateUserResponse struct {
	ID string `json:"id"`
}

func CreateUserHandler(deps CreateUserDependencies) echo.HandlerFunc {
	validate := validator.New()

	return func(c echo.Context) error {
		var req CreateUserRequest

		// Bind and validate request
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
		}

		if err := validate.Struct(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		// Call use case
		userID, err := deps.CreateUser(c.Request().Context(), usecase.CreateUserDTO{
			FirstName:                req.FirstName,
			LastName:                 req.LastName,
			Phone:                    req.Phone,
			Email:                    req.Email,
			Password:                 req.Password,
			IsTermAndConditionAgreed: req.IsTermAndConditionAgreed,
		})

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		// Return response
		return c.JSON(http.StatusCreated, CreateUserResponse{
			ID: userID,
		})
	}
}

