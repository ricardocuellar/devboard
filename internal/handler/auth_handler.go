package handler

import (
	"log/slog"

	"github.com/ricardocuellar/devboard/internal/usecase"
	"github.com/ricardocuellar/devboard/internal/validator"
)

// AuthHandler maneja los endpoints de autenticación
type AuthHandler struct {
	usecase   *usecase.AuthUseCase
	validator *validator.Validator
	logger    *slog.Logger
}

// NewAuthHander constructor
func NewAuthHander(uc *usecase.AuthUseCase, v *validator.Validator, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		usecase:   uc,
		validator: v,
		logger:    logger,
	}
}

// registerRequest define los campos para registro
type registerRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

// loginRequest define los campos para login
type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// refreshRequest define el body para renovar tokens
type refreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// logoutRequest define el body para logout
type logoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// tokenResponse es la respuesta estandar cuando se emiten tokens
type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}
