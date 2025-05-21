package user

import (
	"github.com/BrockMekonnen/go-clean-starter/core/di"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/app/query"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/app/usecase"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/delivery/handlers"
)

func MakeUserRoutes() {
	apiRouter := di.GetApiRouter()
	authRouter := di.GetAuthRouter()

	//* Get handlers
	createHandler := delivery.CreateUserHandler(di.MustResolve[usecase.CreateUserUsecase]())
	findUsersHandler := delivery.FindUsersHandler(di.MustResolve[query.FindUsers]())
	deleteHandler := delivery.DeleteUserHandler(di.MustResolve[usecase.DeleteUserUsecase]())
	getUserHandler := delivery.GetUserHandler(di.MustResolve[query.FindUserById]())
	getMeHandler := delivery.GetMeHandler(di.MustResolve[query.FindUserById]())
	generateHandler := delivery.GenerateTokenHandler(di.MustResolve[usecase.GenerateTokenUsecase]())

	//* Unauth routes
	apiRouter.HandleFunc("/users", createHandler).Methods("POST")
	apiRouter.HandleFunc("/users/login", generateHandler).Methods("POST")

	//* Auth Routes
	authRouter.HandleFunc("/users", findUsersHandler).Methods("GET")
	authRouter.HandleFunc("/users/me", getMeHandler).Methods("GET")
	authRouter.HandleFunc("/users/{id}", deleteHandler).Methods("DELETE")
	authRouter.HandleFunc("/users/{id}", getUserHandler).Methods("GET")

}
