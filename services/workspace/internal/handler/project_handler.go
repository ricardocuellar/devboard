package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/ricardocuellar/devboard/services/workspace/internal/domain"
	"github.com/ricardocuellar/devboard/services/workspace/internal/middlewares"
	"github.com/ricardocuellar/devboard/services/workspace/internal/usecase"
	"github.com/ricardocuellar/devboard/services/workspace/internal/validator"
)

// ProjectHandler maneja los endpoints de proyecto
type ProjectHandler struct {
	usecase   *usecase.ProjectUseCase
	validator *validator.Validator
	logger    *slog.Logger
}

// NewProjectHandler constructor
func NewProjectHandler(uc *usecase.ProjectUseCase, v *validator.Validator, logger *slog.Logger) *ProjectHandler {
	return &ProjectHandler{usecase: uc, validator: v, logger: logger}
}

type createProjectRequest struct {
	WorkspaceID string `json:"workspace_id" validate:"required,uuid4"`
	Name        string `json:"name"         validate:"required,min=2,max=100"`
	Description string `json:"description"  validate:"max=500"`
}

type projectResponse struct {
	ID          string    `json:"id"`
	WorkspaceID string    `json:"workspace_id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

func toProjectResponse(p *domain.Project) projectResponse {
	return projectResponse{
		ID:          p.ID,
		WorkspaceID: p.WorkspaceID,
		Name:        p.Name,
		Description: p.Description,
		CreatedBy:   p.CreatedBy,
		CreatedAt:   p.CreatedAt,
	}
}

// Create crea un proyecto
//
// @Summary     Crear proyecto
// @Description Crea un proyecto dentro de un workspace
// @Tags        projects
// @Accept      json
// @Produce     json
// @Param       request body     createProjectRequest true "Datos del proyecto"
// @Success     201     {object} projectResponse
// @Security    BearerAuth
// @Router      /projects [post]
func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createProjectRequest
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

	project, err := h.usecase.CreateProject(r.Context(), req.WorkspaceID, req.Name, req.Description, userID)

	if err != nil {
		RespondError(w, r, h.logger, err)
		return
	}

	if err := RespondJSON(w, http.StatusCreated, toProjectResponse(project)); err != nil {
		h.logger.ErrorContext(r.Context(), "error escribiendo respuesta JSON",
			"error", err, "path", r.URL.Path)
	}

}

// Get obtiene un proyecto por ID
//
// @Summary     Obtener proyecto
// @Description Obtiene un proyecto por su ID
// @Tags        projects
// @Produce     json
// @Param       id  path     string true "ID del proyecto"
// @Success     200 {object} projectResponse
// @Security    BearerAuth
// @Router      /projects/{id} [get]
func (h *ProjectHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID := middlewares.GetUserID(r)
	if userID == "" {
		RespondError(w, r, h.logger, domain.ErrUnauthorized)
		return
	}

	project, err := h.usecase.GetProject(r.Context(), id, userID)

	if err != nil {
		RespondError(w, r, h.logger, err)
		return
	}

	if err := RespondJSON(w, http.StatusOK, toProjectResponse(project)); err != nil {
		h.logger.ErrorContext(r.Context(), "error escribiendo respuesta JSON",
			"error", err, "path", r.URL.Path)
	}
}

// ListByWorkspace lista proyectos de un workspace
//
// @Summary     Listar proyectos del workspace
// @Description Lista todos los proyectos de un workspace
// @Tags        projects
// @Produce     json
// @Param       id  path     string true "ID del workspace"
// @Success     200 {object} PaginatedResponse[projectResponse]
// @Security    BearerAuth
// @Router      /workspaces/{id}/projects [get]
func (h *ProjectHandler) ListByWorkspace(w http.ResponseWriter, r *http.Request) {
	workspaceID := r.PathValue("id")
	userID := middlewares.GetUserID(r)
	if userID == "" {
		RespondError(w, r, h.logger, domain.ErrUnauthorized)
		return
	}

	projects, err := h.usecase.ListWorkspaceProjects(r.Context(), workspaceID, userID)

	if err != nil {
		RespondError(w, r, h.logger, err)
		return
	}

	response := make([]projectResponse, len(projects))
	for i, p := range projects {
		response[i] = toProjectResponse(p)
	}

	if err := RespondJSON(w, http.StatusOK, PaginatedResponse[projectResponse]{
		Data:    response,
		HasMore: false,
	}); err != nil {
		h.logger.ErrorContext(r.Context(), "error escribiendo respuesta JSON",
			"error", err, "path", r.URL.Path)
	}
}
