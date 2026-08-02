package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/ricardocuellar/devboard/internal/domain"
	"github.com/ricardocuellar/devboard/internal/middlewares"
	"github.com/ricardocuellar/devboard/internal/usecase"
	"github.com/ricardocuellar/devboard/internal/validator"
)

// WorkspaceHandler maneja los endpoints de workspace
type WorkspaceHandler struct {
	usecase   *usecase.WorkspaceUseCase
	validator *validator.Validator
	logger    *slog.Logger
}

// NewWorkspaceHandler constructor
func NewWorkspaceHandler(uc *usecase.WorkspaceUseCase, v *validator.Validator, logger *slog.Logger) *WorkspaceHandler {
	return &WorkspaceHandler{usecase: uc, validator: v, logger: logger}
}

// DTOs de request
type createWorkspaceRequest struct {
	Name        string `json:"name"        validate:"required,min=2,max=100"`
	Description string `json:"description" validate:"max=500"`
}

type addMemberRequest struct {
	UserID string `json:"user_id" validate:"required,uuid4"`
	Role   string `json:"role"    validate:"required,oneof=admin member viewer"`
}

// DTOs de respuesta
type workspaceResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description,omitempty"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

type memberResponse struct {
	ID          string    `json:"id"`
	WorkspaceID string    `json:"workspace_id"`
	UserID      string    `json:"user_id"`
	Role        string    `json:"role"`
	JoinedAt    time.Time `json:"joined_at"`
}

// Funciones para Response
func toWorkspaceResponse(ws *domain.Workspace) workspaceResponse {
	return workspaceResponse{
		ID:          ws.ID,
		Name:        ws.Name,
		Slug:        ws.Slug,
		Description: ws.Description,
		CreatedBy:   ws.CreatedBy,
		CreatedAt:   ws.CreatedAt,
	}
}

func toMemberResponse(m *domain.WorkspaceMember) memberResponse {
	return memberResponse{
		ID:          m.ID,
		WorkspaceID: m.WorkspaceID,
		UserID:      m.UserID,
		Role:        string(m.Role),
		JoinedAt:    m.JoinedAt,
	}
}

// Create crea un workspace
//
// @Summary     Crear workspace
// @Description Crea un nuevo workspace. El creador queda como admin automáticamente.
// @Tags        workspaces
// @Accept      json
// @Produce     json
// @Param       request body     createWorkspaceRequest true "Datos del workspace"
// @Success     201     {object} workspaceResponse
// @Failure     400     {object} ProblemDetails
// @Failure     409     {object} ProblemDetails
// @Security    BearerAuth
// @Router      /workspaces [post]
func (h *WorkspaceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createWorkspaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, r, h.logger, domain.ErrInvalidInput)
		return
	}

	if errs := h.validator.Validate(req); errs != nil {
		if err := RespondValidationError(w, errs); err != nil {
			h.logger.ErrorContext(r.Context(), "error escribiendo respuesta de validación",
				"error", err, "path", r.URL.Path)
		}
		return
	}

	userID := middlewares.GetUserID(r)

	if userID == "" {
		RespondError(w, r, h.logger, domain.ErrUnauthorized)
		return
	}

	ws, err := h.usecase.CreateWorkspace(r.Context(), req.Name, req.Description, userID)

	if err != nil {
		RespondError(w, r, h.logger, err)
		return
	}

	if err := RespondJSON(w, http.StatusCreated, toWorkspaceResponse(ws)); err != nil {
		h.logger.ErrorContext(r.Context(), "error escribiendo respuesta JSON",
			"error", err, "path", r.URL.Path)
	}

}

// List lista los workspaces del usuario autenticado
//
// @Summary     Listar mis workspaces
// @Description Lista todos los workspaces donde el usuario autenticado es miembro
// @Tags        workspaces
// @Produce     json
// @Success     200 {object} PaginatedResponse[workspaceResponse]
// @Security    BearerAuth
// @Router      /workspaces [get]
func (h *WorkspaceHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)
	if userID == "" {
		RespondError(w, r, h.logger, domain.ErrUnauthorized)
		return
	}

	workspaces, err := h.usecase.ListUserWorkspaces(r.Context(), userID)

	if err != nil {
		RespondError(w, r, h.logger, err)
		return
	}

	response := make([]workspaceResponse, len(workspaces))
	for i, ws := range workspaces {
		response[i] = toWorkspaceResponse(ws)
	}

	if err := RespondJSON(w, http.StatusOK, PaginatedResponse[workspaceResponse]{
		Data:    response,
		HasMore: false,
	}); err != nil {
		h.logger.ErrorContext(r.Context(), "error escribiendo respuesta JSON",
			"error", err, "path", r.URL.Path)
	}

}

// Get obtiene un workspace por ID
//
// @Summary     Obtener workspace
// @Description Obtiene un workspace. Solo accesible para miembros.
// @Tags        workspaces
// @Produce     json
// @Param       id  path     string true "ID del workspace"
// @Success     200 {object} workspaceResponse
// @Failure     403 {object} ProblemDetails
// @Failure     404 {object} ProblemDetails
// @Security    BearerAuth
// @Router      /workspaces/{id} [get]
func (h *WorkspaceHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID := middlewares.GetUserID(r)
	if userID == "" {
		RespondError(w, r, h.logger, domain.ErrUnauthorized)
		return
	}

	ws, err := h.usecase.GetWorkspace(r.Context(), id, userID)

	if err != nil {
		RespondError(w, r, h.logger, err)
		return
	}

	if err := RespondJSON(w, http.StatusOK, toWorkspaceResponse(ws)); err != nil {
		h.logger.ErrorContext(r.Context(), "error escribiendo respuesta JSON",
			"error", err, "path", r.URL.Path)
	}
}

// AddMember agrega un miembro al workspace
//
// @Summary     Agregar miembro
// @Description Agrega un usuario al workspace. Solo admins.
// @Tags        workspaces
// @Accept      json
// @Produce     json
// @Param       id      path     string           true "ID del workspace"
// @Param       request body     addMemberRequest true "Datos del miembro"
// @Success     201     {object} memberResponse
// @Failure     403     {object} ProblemDetails
// @Security    BearerAuth
// @Router      /workspaces/{id}/members [post]
func (h *WorkspaceHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	// TODO : Get user ID

	var req addMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, r, h.logger, domain.ErrInvalidInput)
		return
	}

	if errs := h.validator.Validate(req); errs != nil {
		if err := RespondValidationError(w, errs); err != nil {
			h.logger.ErrorContext(r.Context(), "error escribiendo respuesta de validación",
				"error", err, "path", r.URL.Path)
		}
		return
	}

	// TODO
}

// ListMembers lista los miembros de un workspace
//
// @Summary     Listar miembros
// @Description Lista todos los miembros del workspace
// @Tags        workspaces
// @Produce     json
// @Param       id  path     string true "ID del workspace"
// @Success     200 {object} PaginatedResponse[memberResponse]
// @Security    BearerAuth
// @Router      /workspaces/{id}/members [get]
func (h *WorkspaceHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// RemoveMember elimina un miembro del workspace
//
// @Summary     Eliminar miembro
// @Description Elimina un usuario del workspace. Solo admins.
// @Tags        workspaces
// @Param       id     path string true "ID del workspace"
// @Param       userID path string true "ID del usuario a eliminar"
// @Success     204    "Sin contenido"
// @Failure     403    {object} ProblemDetails
// @Security    BearerAuth
// @Router      /workspaces/{id}/members/{userID} [delete]
func (h *WorkspaceHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	// TODO
}

/*
if err := RespondJSON(w, http.StatusCreated, toWorkspaceResponse(ws)); err != nil {
		h.logger.ErrorContext(r.Context(), "error escribiendo respuesta JSON",
			"error", err, "path", r.URL.Path)
	}

*/
