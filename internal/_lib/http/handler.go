package http

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

// AsyncHandler defines the async handler function type
type AsyncHandler func(c echo.Context) error

// ControllerHandler defines a factory function that creates handlers with dependencies
type ControllerHandler func(dependencies interface{}) AsyncHandler

// HandlerMiddleware creates echo middleware that injects dependencies
func HandlerMiddleware(container *dig.Container, handlerFactory ControllerHandler) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get the request-scoped container from context
			reqContainer, ok := c.Get("container").(*dig.Container)
			if !ok {
				return echo.NewHTTPError(500, "Can't find the request container! Have you registered the container middleware?")
			}

			// Resolve dependencies and create the handler
			var handler AsyncHandler
			err := reqContainer.Invoke(func(deps interface{}) {
				handler = handlerFactory(deps)
			})
			if err != nil {
				return err
			}

			// Execute the handler
			return handler(c)
		}
	}
}