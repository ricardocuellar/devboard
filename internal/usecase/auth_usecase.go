package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ricardocuellar/devboard/internal/auth"
	"github.com/ricardocuellar/devboard/internal/domain"
)

// authUserRepository define lo que AuthUsecase necesita del repositorio de usuarios
type authUserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
}

// AuthTokens contiene los tokens resultantes del login
type AuthTokens struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    time.Duration
}

// AuthUseCase orquesta el flujo completo de autenticación
type AuthUseCase struct {
	userRepo             authUserRepository
	jwtManager           *auth.JWTManager
	tokenStore           *auth.TokenStore
	refreshTokenDuration time.Duration
}

// NewAuthUseCase constructor
func NewAuthUseCase(
	userRepo authUserRepository,
	jwtManager *auth.JWTManager,
	tokenStore *auth.TokenStore,
	refreshTokenDuration time.Duration,
) *AuthUseCase {
	return &AuthUseCase{
		userRepo:             userRepo,
		jwtManager:           jwtManager,
		tokenStore:           tokenStore,
		refreshTokenDuration: refreshTokenDuration,
	}
}

// Register crea un nuevo usuario con password hasheado
func (uc *AuthUseCase) Register(ctx context.Context, email, name, password string) (*AuthTokens, error) {
	passwordHash, err := auth.HashPassword(password)

	if err != nil {
		return nil, fmt.Errorf("hasheando password: %w", err)
	}

	user := domain.NewUser(email, name, passwordHash)

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("creando usuario: %w", err)
	}

	// TODO generateTokenPair
	tokens, err := uc.generateTokenPair(ctx, user.ID, user.Email, "member")

	if err != nil {
		return nil, err
	}

	return tokens, nil

}

// generateTokenPair genera un access token y un refresh token para el usuario
func (uc *AuthUseCase) generateTokenPair(ctx context.Context, userID, email, role string) (*AuthTokens, error) {
	accessToken, err := uc.jwtManager.GenerateAccessToken(userID, email, role)

	if err != nil {
		return nil, fmt.Errorf("generando access token: %w", err)
	}

	refreshToken, err := uc.jwtManager.GenerateRefreshToken(userID, email, role)

	if err != nil {
		return nil, fmt.Errorf("generando refresh token: %w", err)
	}

	refreshClaims, err := uc.jwtManager.VerifyToken(refreshToken)

	if err != nil {
		return nil, fmt.Errorf("extrayendo claims del refresh token: %w", err)
	}

	if err := uc.tokenStore.StoreRefreshToken(
		ctx,
		refreshClaims.ID,
		userID,
		uc.refreshTokenDuration,
	); err != nil {
		return nil, fmt.Errorf("guardando refresh token: %w", err)
	}

	return &AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    15 * time.Minute,
	}, nil
}

// Login verifica credenciales y retorna tokens si son válidas
func (uc *AuthUseCase) Login(ctx context.Context, email, password string) (*AuthTokens, error) {
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("credenciales inválidas: %w", domain.ErrUnauthorized)
		}
		return nil, fmt.Errorf("buscando usuario: %w", err)
	}

	if err := auth.VerifyPassword(password, user.PasswordHash); err != nil {
		if errors.Is(err, auth.ErrPasswordMissmatch) {
			return nil, fmt.Errorf("credenciales inválidas: %w", domain.ErrUnauthorized)
		}
		return nil, fmt.Errorf("verificando password: %w", err)
	}

	tokens, err := uc.generateTokenPair(ctx, user.ID, user.Email, "member")
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

// RefreshTokens renueva el par de tokens usando un refresh token válido
func (uc *AuthUseCase) RefreshTokens(ctx context.Context, refreshToken string) (*AuthTokens, error) {
	// Verificar firma
	claims, err := uc.jwtManager.VerifyToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("refresh token inválido: %w", domain.ErrUnauthorized)
	}

	// Verificar que el token existe en Redis y no fue revocado
	_, err = uc.tokenStore.ValidateRefreshToken(ctx, claims.ID)
	if err != nil {

		if errors.Is(err, auth.ErrTokenRevoked) || errors.Is(err, auth.ErrTokenNotFound) {
			return nil, fmt.Errorf("refresh token no válido: %w", domain.ErrUnauthorized)
		}

		return nil, fmt.Errorf("validando refresh token: %w", err)
	}

	// Revocar el token
	if err := uc.tokenStore.RevokedRefreshToken(ctx, claims.ID, uc.refreshTokenDuration); err != nil {
		return nil, fmt.Errorf("revocando token anterior: %w", err)
	}

	// Generar nuevos tokens
	tokens, err := uc.generateTokenPair(ctx, claims.UserID, claims.Email, claims.Role)

	if err != nil {
		return nil, err
	}

	return tokens, nil

}

// Logout revoca el refresh token del usuario
func (uc *AuthUseCase) Logout(ctx context.Context, refreshToken string) error {
	claims, err := uc.jwtManager.VerifyToken(refreshToken)
	if err != nil {
		return nil
	}

	if err := uc.tokenStore.RevokedRefreshToken(ctx, claims.ID, uc.refreshTokenDuration); err != nil {
		return fmt.Errorf("revocando token en logout: %w", err)
	}

	return nil

}
