// Package auth jwt
package auth

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// ErrTokenExpired indica que el token expiró.
var ErrTokenExpired = errors.New("token expirado")

// ErrTokenInvalid indica que el token tiene firma inválida
var ErrTokenInvalid = errors.New("token inválido")

// Claims define el contenido del JWT
type Claims struct {
	UserID               string `json:"user_id"`
	Email                string `json:"email"`
	Role                 string `json:"role"`
	jwt.RegisteredClaims        // campos estándar de JWT
}

type TokenVerifier struct {
	secretKey []byte
}

func NewTokenVerifier(secretKey []byte) (*TokenVerifier, error) {
	if len(secretKey) < 32 {
		return nil, fmt.Errorf("secret key debe tener al menos 32 bytes, tiene %d", len(secretKey))
	}

	return &TokenVerifier{
		secretKey: secretKey,
	}, nil
}

// VerifyToken verifica la firma y extrae los claims de un JWT token
func (verifier *TokenVerifier) VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("algoritmo inesperado: %v", t.Header["alg"])
			}

			return verifier.secretKey, nil
		},
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithIssuer("devboard"),
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}

		return nil, ErrTokenInvalid
	}

	claims, ok := token.Claims.(*Claims)

	if !ok || !token.Valid {
		return nil, ErrTokenInvalid
	}

	return claims, nil
}
