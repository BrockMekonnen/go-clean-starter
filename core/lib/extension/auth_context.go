package extension

import (
	"context"
	"errors"
)

type contextKey string

const (
	AuthContextKey contextKey = "auth"
)

type AuthContext struct {
	Credentials struct {
		UID   string
		Scope []string
	}
	Artifacts struct {
		AccessToken string
	}
}

func GetAuthContextData(ctx context.Context) (AuthContext, error) {
	authCtx, ok := ctx.Value(AuthContextKey).(AuthContext)
	if !ok {
		return AuthContext{}, errors.New("unauthorized: missing auth context")
	}

	// Extract UID
	return authCtx, nil
}
