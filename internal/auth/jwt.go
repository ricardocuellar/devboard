// Package auth jwt
package auth

import (
	"crypto/rand"
	"errors"
	"fmt"
	"time"

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

// JWTConfig contiene la configuración del manager de tokens
type JWTConfig struct {
	SecretKey            []byte
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

// JWTManager maneja la creación y verificación de tokens JWT
type JWTManager struct {
	cfg JWTConfig
}

// NewJWTManager crea un JWT Manager
func NewJWTManager(cfg JWTConfig) (*JWTManager, error) {
	// OWASP -> 32 bytes HS256
	if len(cfg.SecretKey) < 32 {
		return nil, fmt.Errorf("secret key debe tener al menos 32 bytes, tiene %d", len(cfg.SecretKey))
	}

	return &JWTManager{cfg: cfg}, nil
}

// GenerateAccessToken crea un access token de corta duración
func (m *JWTManager) GenerateAccessToken(userID, email, role string) (string, error) {
	return m.generateToken(userID, email, role, m.cfg.AccessTokenDuration)
}

// GenerateRefreshToken crea un refresh token de larga duración
func (m *JWTManager) GenerateRefreshToken(userID, email, role string) (string, error) {
	return m.generateToken(userID, email, role, m.cfg.RefreshTokenDuration)
}

// VerifyToken verifica la firma y extrae los claims de un JWT token
func (m *JWTManager) VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("algoritmo inesperado: %v", t.Header["alg"])
			}

			return m.cfg.SecretKey, nil
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

// generateToken crea un token JWT firmado con los claims dados
func (m *JWTManager) generateToken(userID, email, role string, duration time.Duration) (string, error) {
	now := time.Now()

	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        generateTokenID(),
			Issuer:    "devboard",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(m.cfg.SecretKey)

	if err != nil {
		return "", fmt.Errorf("firmando token: %w", err)
	}

	return signed, nil
}

// generateTokenID genera un ID único para el token usando crypto/rand
func generateTokenID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
