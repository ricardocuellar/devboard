package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// ErrTokenRevoked indica que el token fue revocado
var ErrTokenRevoked = errors.New("token revocado")

// ErrTokenNotFound indica que el refresh token no existe (en redis)
var ErrTokenNotFound = errors.New("refresh token no encontrado")

// TokenStore maneja el almacenamiento de refresh tokens en Redis
type TokenStore struct {
	redis *redis.Client
}

// NewTokenStore crea un TokenStore con el cliente Redis
func NewTokenStore(redisClient *redis.Client) *TokenStore {
	return &TokenStore{redis: redisClient}
}

// refreshKey construye clave Redis para refresh token
func refreshKey(tokenID string) string {
	return "refresh:" + tokenID
}

// revokedKey construye la clave Redis para un token revocado
func revokedKey(tokenID string) string {
	return "revoked:" + tokenID
}

// StoreRefreshToken guarda un refresh token en Redis
func (s *TokenStore) StoreRefreshToken(ctx context.Context, tokenID, userID string, ttl time.Duration) error {
	err := s.redis.Set(
		ctx,
		refreshKey(tokenID),
		userID,
		ttl,
	).Err()

	if err != nil {
		return fmt.Errorf("guardando refresh token en Redis: %w", err)
	}

	return nil
}
