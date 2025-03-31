package usecase

// import (
// 	"github.com/BrockMekonnen/go-clean-starter/internal/auth/domain"
// 	userDomain "github.com/BrockMekonnen/go-clean-starter/internal/user/domain"
// 	sharedDomain "github.com/BrockMekonnen/go-clean-starter/internal/_shared/domain"
// )

// // Dependencies struct for passing dependencies
// type Dependencies struct {
// 	AuthRepository domain.AuthRepository
// 	UserRepository userDomain.UserRepository
// }

// // LoginParams defines the input structure for the service
// type LoginParams struct {
// 	Email    string
// 	Password string
// }

// // TokenAndUser defines the structure for the result
// type TokenAndUser struct {
// 	Token string `json:"token"`
// }

// // GenerateToken is the interface for the service
// type GenerateToken interface {
// 	Execute(params LoginParams) (*TokenAndUser, error)
// }

// // GenerateTokenService implements the GenerateToken interface
// type GenerateTokenService struct {
// 	AuthRepo    domain.AuthRepository
// 	UserRepo    userDomain.UserRepository
// }

// // NewGenerateTokenService creates a new instance of GenerateTokenService
// func NewGenerateTokenService(deps Dependencies) GenerateToken {
// 	return &GenerateTokenService{
// 		AuthRepo: deps.AuthRepository,
// 		UserRepo: deps.UserRepository,
// 	}
// }

// // Execute validates user credentials and generates a token
// func (s *GenerateTokenService) Execute(params LoginParams) (*TokenAndUser, error) {
// 	user, err := s.UserRepo.FindByEmail(params.Email)
// 	if err != nil || user == nil {
// 		return nil, sharedDomain.NewBusinessError("Incorrect Email.", "USER_NOT_FOUND")
// 	}

// 	isMatch, err := s.AuthRepo.Compare(params.Password, user.Password)
// 	if err != nil || !isMatch {
// 		return nil, sharedDomain.NewBusinessError("Incorrect Password.", "INVALID_PASSWORD")
// 	}

// 	token, err := s.AuthRepo.Generate(domain.Credentials{
// 		Uid:   user.Id,
// 		Scope: user.Roles,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &TokenAndUser{Token: token}, nil
// }