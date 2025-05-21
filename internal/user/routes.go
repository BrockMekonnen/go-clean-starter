package user

import (
	"github.com/BrockMekonnen/go-clean-starter/core/di"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/delivery/handlers"
)

func RegisterUserRoutes() {
	apiRouter := di.GetApiRouter()
	authRouter := di.GetAuthRouter()

	//* Get handler dependencies
	createHandler := delivery.MakeCreateUserHandler(di.MustResolve[delivery.CreateUserHandlerDeps]())
	findUsersHandler := delivery.MakeFindUsersHandler(di.MustResolve[delivery.FindUsersHandlerDeps]())
	deleteHandler := delivery.MakeDeleteUserHandler(di.MustResolve[delivery.DeleteUserHandlerDeps]())
	getUserHandler := delivery.MakeGetUserHandler(di.MustResolve[delivery.GetUserHandlerDeps]())
	getMeHandler := delivery.MakeGetMeHandler(di.MustResolve[delivery.GetMeHandlerDeps]())
	generateHandler := delivery.MakeGenerateTokenHandler(di.MustResolve[delivery.GenerateTokenHandlerDeps]())

	//* Register routes
	apiRouter.HandleFunc("/users", createHandler).Methods("POST")
	apiRouter.HandleFunc("/users/login", generateHandler).Methods("POST")

	//* Auth Routes
	authRouter.HandleFunc("/users", findUsersHandler).Methods("GET")
	authRouter.HandleFunc("/users/me", getMeHandler).Methods("GET")
	authRouter.HandleFunc("/users/{id}", deleteHandler).Methods("DELETE")
	authRouter.HandleFunc("/users/{id}", getUserHandler).Methods("GET")

}
