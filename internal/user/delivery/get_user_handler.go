package delivery

import (
	"net/http"

	"github.com/BrockMekonnen/go-clean-starter/internal/user/app/query"
	"github.com/labstack/echo/v4"
)

type Dependencies struct {
	FindUserById query.FindUserById
}

type GetUserRequest struct {
	Id uint `json:"id"`
}

func GetUserHandler(deps Dependencies) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Bind and validate request
		var req GetUserRequest
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
		}

		if err := c.Validate(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "User ID is required")
		}

		// Call query to find user
		user, err := deps.FindUserById.Handle(c.Request().Context(), req.Id)
		if err != nil {
			// Handle specific error types
			// if err == query.ErrUserNotFound {
			// 	return echo.NewHTTPError(http.StatusNotFound, "User not found")
			// }
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to find user")
		}

		// Return user data
		return c.JSON(http.StatusOK, user)
	}
}
