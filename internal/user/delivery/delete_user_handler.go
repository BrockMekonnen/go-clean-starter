package delivery

import (
	"net/http"

	"github.com/BrockMekonnen/go-clean-starter/internal/user/app/usecase"
	"github.com/labstack/echo/v4"
)

type DeleteUserDependencies struct {
	DeleteUser usecase.DeleteUser
}

type DeleteUserRequest struct {
	Id uint `json:"id"`
}

func DeleteUserHandler(deps DeleteUserDependencies) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req DeleteUserRequest

		// Bind and validate request
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
		}

		_, err := deps.DeleteUser(c.Request().Context(), req.Id)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		// Return response
		return c.JSON(http.StatusCreated, nil)
	}
}
