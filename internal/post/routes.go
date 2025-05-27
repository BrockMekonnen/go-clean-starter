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
	createPostHandler := handlers.CreatePostHandler(di.MustResolve[usecase.CreatePostUsecase](), di.MustResolve[query.FindPostById]())
	publishPostHandler := handlers.PublishPostHandler(di.MustResolve[usecase.PublishPostUsecase](), di.MustResolve[query.FindPostById]())
	updatePostHandler := handlers.UpdatePostHandler(di.MustResolve[usecase.UpdatePostUsecase](), di.MustResolve[query.FindPostById]())
	deletePostHandler := handlers.DeletePostHandler(di.MustResolve[usecase.DeletePostUsecase]())
	findPostByIdHandler := handlers.FindPostByIdHandler(di.MustResolve[query.FindPostById]())
	findMyPostsHandler := handlers.FindMyPostsHandler(di.MustResolve[query.FindPosts]())
	findPostsHandler := handlers.FindPostsHandler(di.MustResolve[query.FindPosts]())

	//* Register In Auth Routes
	authRouter.HandleFunc("/posts", createPostHandler).Methods("POST")
	authRouter.HandleFunc("/posts", findPostsHandler).Methods("GET")
	authRouter.HandleFunc("/me/posts", findMyPostsHandler).Methods("GET")
	authRouter.HandleFunc("/posts/{id}", findPostByIdHandler).Methods("GET")
	authRouter.HandleFunc("/posts/{id}", deletePostHandler).Methods("DELETE")
	authRouter.HandleFunc("/posts/{id}", updatePostHandler).Methods("PATCH")
	authRouter.HandleFunc("/posts/{id}/publish", publishPostHandler).Methods("PATCH")
}
