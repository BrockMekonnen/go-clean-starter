package user

import (
	"github.com/BrockMekonnen/go-clean-starter/core/di"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/app/query"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/app/usecase"
	"github.com/BrockMekonnen/go-clean-starter/internal/user/interface/handlers"
)

func MakeUserRoutes() {
	apiRouter := di.GetApiRouter()
	authRouter := di.GetAuthRouter()

	//* Get handlers
	createUserHandler := handlers.CreateUserHandler(di.MustResolve[usecase.CreateUserUsecase]())
	findUsersHandler := handlers.FindUsersHandler(di.MustResolve[query.FindUsers]())
	deleteHandler := handlers.DeleteUserHandler(di.MustResolve[usecase.DeleteUserUsecase]())
	getUserHandler := handlers.GetUserHandler(di.MustResolve[query.FindUserById]())
	getMeHandler := handlers.GetMeHandler(di.MustResolve[query.FindUserById]())
	generateHandler := handlers.GenerateTokenHandler(di.MustResolve[usecase.GenerateTokenUsecase]())

	//* Unauth routes
	apiRouter.HandleFunc("/users", createUserHandler).Methods("POST")
	apiRouter.HandleFunc("/users/login", generateHandler).Methods("POST")

	//* Auth Routes
	authRouter.HandleFunc("/users", findUsersHandler).Methods("GET")
	authRouter.HandleFunc("/users/me", getMeHandler).Methods("GET")
	authRouter.HandleFunc("/users/{id}", deleteHandler).Methods("DELETE")
	authRouter.HandleFunc("/users/{id}", getUserHandler).Methods("GET")
}
