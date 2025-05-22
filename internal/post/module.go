package post

import (
	"github.com/BrockMekonnen/go-clean-starter/core/di"
	"github.com/BrockMekonnen/go-clean-starter/internal/post/app/usecase"
	"github.com/BrockMekonnen/go-clean-starter/internal/post/infrastructure"
)

func MakePostModule() error {
	//* Get dependencies
	logger := di.GetLogger()
	db := di.GetDatabase().GetDB()

	if err := infrastructure.InitPostTable(db); err != nil {
		logger.WithFields(map[string]interface{}{"error": err.Error()}).Error("Failed to initialize post table")
		return err
	}

	//* Register this module repository
	if err := di.ProvideWrapper("UserRepository", infrastructure.MakePostRepository); err != nil {
		return err
	}

	//* Register this app queries
	if err := di.ProvideWrapper("FindPostById", infrastructure.MakeFindPostById); err != nil {
		return err
	}
	if err := di.ProvideWrapper("FindPosts", infrastructure.MakeFindPosts); err != nil {
		return err
	}

	//* Register this module app usecases
	if err := di.ProvideWrapper("CreatePostUsecase", usecase.MakeCreatePostUsecase); err != nil {
		return err
	}
	if err := di.ProvideWrapper("DeletePostUsecase", usecase.MakeDeletePostUsecase); err != nil {
		return err
	}
	if err := di.ProvideWrapper("PublishPostUsecase", usecase.MakePublishPostUsecase); err != nil {
		return err
	}
	if err := di.ProvideWrapper("UpdatePostUsecase", usecase.MakeUpdatePostUsecase); err != nil {
		return err
	}

	MakePostRoutes()

	logger.Info("Post module initialized successfully.")
	return nil
}
