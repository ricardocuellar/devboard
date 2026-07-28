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

// ValidateRefreshToken verifica que un refresh token existe y no fue revocado
func (s *TokenStore) ValidateRefreshToken(ctx context.Context, tokenID string) (string, error) {
	revoked, err := s.redis.Exists(ctx, revokedKey(tokenID)).Result()

	if err != nil {
		return "", fmt.Errorf("verificando revocación: %w", err)
	}

	if revoked > 0 {
		return "", ErrTokenRevoked
	}

	userID, err := s.redis.Get(ctx, refreshKey(tokenID)).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrTokenNotFound
		}

		return "", fmt.Errorf("obteniendo refresh token: %w", err)
	}

	return userID, nil
}

// RevokedRefreshToken marca un token como revocado y elimina el store activo
func (s *TokenStore) RevokedRefreshToken(ctx context.Context, tokenID string, ttl time.Duration) error {
	pipe := s.redis.Pipeline()

	pipe.Set(ctx, revokedKey(tokenID), "1", ttl)
	pipe.Del(ctx, refreshKey(tokenID))

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("revocando refresh token: %w", err)
	}

	return nil
}

// RevokeAllUserTokens revoca todos los refresh tokens de un usuario
func (s *TokenStore) RevokeAllUserTokens(ctx context.Context, userID string, ttl time.Duration) error {
	var cursor uint64
	pattern := "refresh:*"

	for {
		keys, nextCursor, err := s.redis.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return fmt.Errorf("buscando tokens del usuario: %w", err)
		}

		for _, key := range keys {
			storedUserID, err := s.redis.Get(ctx, key).Result()
			if err != nil {
				continue
			}

			if storedUserID == userID {
				tokenID := key[len("refresh:"):]
				_ = s.RevokedRefreshToken(ctx, tokenID, ttl)
			}
		}

		cursor = nextCursor

		if cursor == 0 {
			break
		}
	}

	return nil
}
