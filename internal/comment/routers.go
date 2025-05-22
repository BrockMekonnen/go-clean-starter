package comment

import (
	"github.com/BrockMekonnen/go-clean-starter/core/di"
	"github.com/BrockMekonnen/go-clean-starter/internal/comment/app/usecase"
	"github.com/BrockMekonnen/go-clean-starter/internal/comment/interface/handlers"
)

func MakeCommentRoutes() {
	authRouter := di.GetAuthRouter()

	//* Get Handlers
	createCommentHandler := handlers.CreateCommentHandler(di.MustResolve[usecase.CreateCommentUsecase]())

	//* Auth Routes
	authRouter.HandleFunc("/posts/{postId}/comments", createCommentHandler).Methods("POST")
}
