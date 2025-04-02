package user

import (
	"github.com/BrockMekonnen/go-clean-starter/core/di"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/delivery"
)

func RegisterUserRoutes() {
	apiRouter := di.GetApiRouter()
	authRouter := di.GetAuthRouter()

	//* Get handler dependencies
	createHandler := delivery.NewCreateUserHandler(di.MustResolve[delivery.CreateUserHandlerDeps]())
	deleteHandler := delivery.NewDeleteUserHandler(di.MustResolve[delivery.DeleteUserHandlerDeps]())
	getUserHandler := delivery.GetUserHandler(di.MustResolve[delivery.GetUserHandlerDeps]())
	generateHandler := delivery.NewGenerateTokenHandler(di.MustResolve[delivery.GenerateTokenHandlerDeps]())

	//* Register routes
	apiRouter.HandleFunc("/users/login", generateHandler).Methods("POST")
	authRouter.HandleFunc("/users/{id}", deleteHandler).Methods("DELETE")
	authRouter.HandleFunc("/users/{id}", getUserHandler).Methods("GET")
	apiRouter.HandleFunc("/users", createHandler).Methods("POST")
}
