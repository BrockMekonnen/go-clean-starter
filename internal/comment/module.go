package comment

import (
	"github.com/BrockMekonnen/go-clean-starter/core/di"
	"github.com/BrockMekonnen/go-clean-starter/internal/comment/app/usecase"
	"github.com/BrockMekonnen/go-clean-starter/internal/comment/infrastructure"
)

func MakeCommentModule() error {
	//* Get dependencies
	logger := di.GetLogger()
	db := di.GetDatabase().GetDB()

	//* Initialize this module tables
	if err := infrastructure.InitCommentTable(db); err != nil {
		logger.WithFields(map[string]interface{}{"error": err.Error()}).Error("Failed to initialize comment table")
		return err
	}

	//* Register this module repository
	if err := di.ProvideWrapper("CommentRepository", infrastructure.MakeCommentRepository); err != nil {
		return err
	}

	//* Register this module app usecases
	if err := di.ProvideWrapper("CreateCommentUsecase", usecase.MakeCreateCommentUsecase); err != nil {
		return err
	}

	MakeCommentRoutes()

	logger.Info("Comment module initialized successfully.")
	return nil
}
