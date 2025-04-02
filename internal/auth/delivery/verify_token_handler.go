package delivery

import (
	"context"
	"net/http"
	"strings"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/errors"
	"github.com/BrockMekonnen/go-clean-starter/internal/auth/app/usecase"
	"github.com/gorilla/mux"
)

type VerifyTokenHandlerDeps struct {
	VerifyToken usecase.VerifyTokenUsecase
}

type contextKey string

const (
	authContextKey contextKey = "auth"
)

type AuthContext struct {
	IsAuthenticated bool
	IsAuthorized    bool
	IsInjected      bool
	Credentials     struct {
		UID   uint
		Scope []string
	}
	Artifacts struct {
		AccessToken string
	}
}

func NewVerifyTokenHandler(deps VerifyTokenHandlerDeps) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				panic(errors.NewBadRequestError("Missing or wrong Authorization request header", "", w))
			}

			// Extract token (fixed: there was a missing space after "Bearer")
			accessToken := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

			// Verify token
			credentials, err := deps.VerifyToken(r.Context(), accessToken)
			if err != nil {
				panic(errors.NewBadRequestError("Invalid token", err.Error(), w))
			}

			// Set auth context
			ctx := context.WithValue(r.Context(), authContextKey, AuthContext{
				IsAuthenticated: true,
				IsAuthorized:    true,
				IsInjected:      true,
				Credentials: struct {
					UID   uint
					Scope []string
				}{
					UID:   credentials.Uid,
					Scope: credentials.Scope,
				},
				Artifacts: struct {
					AccessToken string
				}{
					AccessToken: accessToken,
				},
			})

			// Call next handler with updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
