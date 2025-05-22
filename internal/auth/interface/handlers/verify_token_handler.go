package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/errors"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/extension"
	"github.com/BrockMekonnen/go-clean-starter/internal/auth/app/usecase"
	"github.com/gorilla/mux"
)

func VerifyTokenMiddleware(verifyToken usecase.VerifyTokenUsecase) mux.MiddlewareFunc {
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
			credentials, err := verifyToken(r.Context(), accessToken)
			if err != nil {
				panic(errors.NewBadRequestError("Invalid token", err.Error(), w))
			}

			// Set auth context
			ctx := context.WithValue(r.Context(), extension.AuthContextKey, extension.AuthContext{
				Credentials: struct {
					UID   string
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
