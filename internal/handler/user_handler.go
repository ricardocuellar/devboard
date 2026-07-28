package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/ricardocuellar/devboard/internal/auth"
	"github.com/ricardocuellar/devboard/internal/domain"
	"github.com/ricardocuellar/devboard/internal/usecase"
	"github.com/ricardocuellar/devboard/internal/validator"
)

// UserHandler traduce request HTTP a llamadas del UserUseCase
type UserHandler struct {
	usecase   *usecase.UserUseCase
	validator *validator.Validator
	logger    *slog.Logger
}

// NewUserHandler constructor
func NewUserHandler(uc *usecase.UserUseCase, v *validator.Validator, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		usecase:   uc,
		validator: v,
		logger:    logger,
	}
}

type createUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Password string `json:"password" validate:"required,min=8"`
}

type userResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Create user endpoint
//
// @Summary	Crear usuario
// @Description Registra un nuevo usuario en el sistema
// @Tags  users
// @Accept json
// @Produce json
// @Param request body	createUserRequest	true "Datos del usuario"
// @Success 201 {object} userResponse
// @Failure 400 {object} ProblemDetails
// @Router /users [post]
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, r, h.logger, domain.ErrInvalidInput)
		return
	}

	if errs := h.validator.Validate(req); errs != nil {
		if err := RespondValidationError(w, errs); err != nil {
			h.logger.ErrorContext(r.Context(), "error al escribir respuesta de validación", "error", err, "method", r.Method, "path", r.URL.Path)
		}
		return
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		RespondError(w, r, h.logger, domain.ErrInternal)
		return
	}

	user, err := h.usecase.CreateUser(r.Context(), req.Email, req.Name, passwordHash)

	if err != nil {
		RespondError(w, r, h.logger, err)
		return
	}

	if err := RespondJSON(w, http.StatusCreated, userResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}); err != nil {
		h.logger.ErrorContext(r.Context(), "error al escribir respuesta JSON", "error", err, "method", r.Method, "path", r.URL.Path)
	}
}

// Get user endpoint
//
// @Summary	Obtener usuario
// @Description Obtiene un usuario por su ID
// @Tags  users
// @Produce json
// @Param id path	string	true "ID del usuario"
// @Success 200 {object} userResponse
// @Failure 404 {object} ProblemDetails
// @Router /users/{id} [get]
func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	user, err := h.usecase.GetUser(r.Context(), id)

	if err != nil {
		RespondError(w, r, h.logger, err)
		return
	}

	if err := RespondJSON(w, http.StatusOK, userResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}); err != nil {
		h.logger.ErrorContext(
			r.Context(),
			"error al escribir respuesta JSON al obtener usuarios",
			"error",
			err,
			"method",
			r.Method,
			"path", r.URL.Path)

	}
}
