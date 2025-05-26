package post

import (
	"github.com/BrockMekonnen/go-clean-starter/core/di"
	"github.com/BrockMekonnen/go-clean-starter/internal/post/app/query"
	"github.com/BrockMekonnen/go-clean-starter/internal/post/app/usecase"
	"github.com/BrockMekonnen/go-clean-starter/internal/post/interface/handlers"
)

func MakePostRoutes() {
	authRouter := di.GetAuthRouter()

	//* Get Handlers
	createPostHandler := handlers.CreatePostHandler(di.MustResolve[usecase.CreatePostUsecase]())
	deletePostHandler := handlers.DeletePostHandler(di.MustResolve[usecase.DeletePostUsecase]())
	publishPostHandler := handlers.PublishPostHandler(di.MustResolve[usecase.PublishPostUsecase]())
	updatePostHandler := handlers.UpdatePostHandler(di.MustResolve[usecase.UpdatePostUsecase]())
	findPostByIdHandler := handlers.FindPostByIdHandler(di.MustResolve[query.FindPostById]())
	findPostsHandler := handlers.FindPostsHandler(di.MustResolve[query.FindPosts]())
	findMyPostsHandler := handlers.FindMyPostsHandler(di.MustResolve[query.FindPosts]())

	//* Register In Auth Routes
	authRouter.HandleFunc("/posts", createPostHandler).Methods("POST")
	authRouter.HandleFunc("/posts", findPostsHandler).Methods("GET")
	authRouter.HandleFunc("/posts/me", findMyPostsHandler).Methods("GET")
	authRouter.HandleFunc("/posts/{id}", findPostByIdHandler).Methods("GET")
	authRouter.HandleFunc("/posts/{id}", deletePostHandler).Methods("DELETE")
	authRouter.HandleFunc("/posts/{id}", updatePostHandler).Methods("PATCH")
	authRouter.HandleFunc("/posts/publish/{id}", publishPostHandler).Methods("PATCH")
}
