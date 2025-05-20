package infrastructure

import (
	"context"
	"errors"
	"time"

	"github.com/BrockMekonnen/go-clean-starter/core"
	"github.com/BrockMekonnen/go-clean-starter/internal/auth/domain"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthRepositoryImpl implements AuthRepository
type AuthRepository struct {
	config *core.AppConfig
}

// NewAuthRepository creates a new AuthRepository instance
func NewAuthRepository(config *core.AppConfig) domain.AuthRepository {
	return &AuthRepository{config}
}

// Generate generates a JWT token
func (r *AuthRepository) Generate(ctx context.Context, payload domain.Credentials) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":   payload.Uid,
		"scope": payload.Scope,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})

	JWT_SECRET_KEY := r.config.Encryption.JWTKey
	signedToken, err := token.SignedString([]byte(JWT_SECRET_KEY))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// Decode decodes a JWT token
func (r *AuthRepository) Decode(ctx context.Context, accessToken string) (*domain.Credentials, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		JWT_SECRET_KEY := r.config.Encryption.JWTKey
		return []byte(JWT_SECRET_KEY), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &domain.Credentials{
			Uid:   claims["uid"].(string),
			Scope: convertToStringSlice(claims["scope"]),
		}, nil
	}

	return nil, errors.New("invalid token")
}

// Hash hashes a plaintext password using bcrypt
func (r *AuthRepository) Hash(ctx context.Context, plainText string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// Compare compares a plaintext password with a hashed password
func (r *AuthRepository) Compare(ctx context.Context, plainText string, hashedText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedText), []byte(plainText))
	if err != nil {
		return false, err
	}
	return true, nil
}

// convertToStringSlice safely converts an interface{} to []string
func convertToStringSlice(input interface{}) []string {
	if input == nil {
		return []string{}
	}
	if strSlice, ok := input.([]interface{}); ok {
		result := make([]string, len(strSlice))
		for i, v := range strSlice {
			if str, ok := v.(string); ok {
				result[i] = str
			}
		}
		return result
	}
	return []string{}
}
