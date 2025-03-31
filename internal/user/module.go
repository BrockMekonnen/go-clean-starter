package user

// import (
// 	"github.com/labstack/echo/v4"
// 	"github.com/BrockMekonnen/go-clean-starter/internal/user/app/usecase"
// 	"github.com/BrockMekonnen/go-clean-starter/internal/user/app/query"
// 	"github.com/BrockMekonnen/go-clean-starter/internal/user/delivery"
// 	"github.com/BrockMekonnen/go-clean-starter/internal/user/infrastructure"
// )

// // Module struct encapsulates all dependencies for the user module
// type UserModule struct {
// 	UserHandler *delivery.UserHandler
// }

// // NewUserModule initializes and returns a new instance of the User module
// func NewUserModule(e *echo.Echo) *Module {
// 	// // Init Schema
// 	// infrastructure.InitUserTable()

// 	// // Initialize repository
// 	// userRepo := infrastructure.

// 	// // Initialize use cases
// 	// createUserUsecase := usecase.NewCreateUserUsecase(userRepo)
// 	// deleteUserUsecase := usecase.NewDeleteUserUsecase(userRepo)
// 	// getUserByIdUsecase := usecase.NewGetUserByIdUsecase(userRepo)

// 	// // Initialize handlers
// 	// userHandler := &delivery.UserHandler{
// 	// 	RegisterUserUsecase: registerUserUsecase,
// 	// 	DeleteUserUsecase:   deleteUserUsecase,
// 	// 	GetUserByIdUsecase:  getUserByIdUsecase,
// 	// }

// 	// // Register HTTP routes
// 	// userHandler.RegisterRoutes(e)

// 	// // Return the module instance
// 	// return &Module{
// 	// 	UserHandler: userHandler,
// 	// }
// 	return nil
// }