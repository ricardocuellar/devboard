package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/ricardocuellar/devboard/internal/domain"
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

func toTokenResponse(tokens *usecase.AuthTokens) tokenResponse {
	return tokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(tokens.ExpiresIn.Seconds()),
	}
}

// Register endpoint
//
// @Summary     Registrar usuario
// @Description Crea un nuevo usuario y retorna tokens de acceso
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body     registerRequest true "Datos de registro"
// @Success     201     {object} tokenResponse
// @Failure     400     {object} ProblemDetails
// @Failure     409     {object} ProblemDetails
// @Router      /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, r, h.logger, domain.ErrInvalidInput)
		return
	}

	if errs := h.validator.Validate(req); errs != nil {
		if err := RespondValidationError(w, errs); err != nil {
			h.logger.ErrorContext(
				r.Context(),
				"error al escribir respuesta de validación",
				"error", err,
				"method", r.Method,
				"path", r.URL.Path,
			)
		}
		return
	}

	tokens, err := h.usecase.Register(r.Context(), req.Email, req.Name, req.Password)

	if err != nil {
		RespondError(w, r, h.logger, err)
		return
	}

	if err := RespondJSON(w, http.StatusCreated, toTokenResponse(tokens)); err != nil {
		h.logger.ErrorContext(
			r.Context(),
			"error al escribir respuesta JSON",
			"error", err,
			"method", r.Method,
			"path", r.URL.Path,
		)
	}

}

// Login endpoint
//
// @Summary     Iniciar sesión
// @Description Verifica credenciales y retorna tokens de acceso
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body     loginRequest true "Credenciales"
// @Success     200     {object} tokenResponse
// @Failure     401     {object} ProblemDetails
// @Router      /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

}

// Refresh endpoint
//
// @Summary     Renovar tokens
// @Description Renueva el par de tokens usando un refresh token válido
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body     refreshRequest true "Refresh token"
// @Success     200     {object} tokenResponse
// @Failure     401     {object} ProblemDetails
// @Router      /auth/refresh [post]
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {

}

// Logout endpoint
//
// @Summary     Cerrar sesión
// @Description Revoca el refresh token del usuario
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body refreshRequest true "Refresh token a revocar"
// @Success     204     "Sin contenido"
// @Router      /auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {

}

// Respuestas con validaciones
/*
if err := RespondValidationError(w, errs); err != nil {
	h.logger.ErrorContext(
		r.Context(),
		"error al escribir respuesta de validación",
		"error", err,
		"method", r.Method,
		"path", r.URL.Path,
	)
}
return
*/

/*
if err := RespondJSON(w, http.StatusCreated, toTaskResponse(task)); err != nil {
	h.logger.ErrorContext(
		r.Context(),
		"error al escribir respuesta JSON",
		"error", err,
		"method", r.Method,
		"path", r.URL.Path,
	)
}
*/
