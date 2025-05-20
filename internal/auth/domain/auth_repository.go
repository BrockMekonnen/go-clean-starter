package domain

import "context"

type Credentials struct {
	Uid   string  
	Scope []string
}

type AuthRepository interface {
	Generate(ctx context.Context, payload Credentials) (string, error)
	Decode(ctx context.Context, accessToken string) (*Credentials, error)
	Hash(ctx context.Context, plainText string) (string, error)
	Compare(ctx context.Context, plainText, hashedText string) (bool, error)
}