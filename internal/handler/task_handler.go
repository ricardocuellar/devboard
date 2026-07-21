package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/ricardocuellar/devboard/internal/domain"
	"github.com/ricardocuellar/devboard/internal/usecase"
	"github.com/ricardocuellar/devboard/internal/validator"
)

// TaskHandler struct
type TaskHandler struct {
	usercase  *usecase.TaskUseCase
	validator *validator.Validator
	logger    *slog.Logger
}

type updateStatusRequest struct {
	Status domain.TaskStatus `json:"status" validate:"required"`
}

type assignTaskRequest struct {
	AssigneeID string `json:"assignee_id" validate:"required,uuid4"`
}

// NewTaskHandler constructor
func NewTaskHandler(uc *usecase.TaskUseCase, v *validator.Validator, logger *slog.Logger) *TaskHandler {
	return &TaskHandler{
		usercase:  uc,
		validator: v,
		logger:    logger,
	}
}

// createTaskRequest estructura
type createTaskRequest struct {
	ProjectID   string `json:"project_id" validate:"required"`
	Title       string `json:"title" validate:"required,min=3,max=200"`
	Description string `json:"description" validate:"max=2000"`
}

// // createTaskResponse estructura
// type createTaskResponse struct {
// 	ProjectID string `json:"project_id"`
// 	Title     string `json:"title"`
// }

// taskResponse estructura de respuesta
type taskResponse struct {
	ID          string            `json:"id"`
	ProjectID   string            `json:"project_id"`
	Title       string            `json:"title"`
	Description string            `json:"description,omitempty"`
	Status      domain.TaskStatus `json:"status"`
	AssigneeID  *string           `json:"assignee_id,omitempty"`
	CreatedBy   string            `json:"created_by"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// toTaskResponse convierte entidad de dominio a DTO de respuesta
func toTaskResponse(t *domain.Task) taskResponse {
	return taskResponse{
		ID:          t.ID,
		ProjectID:   t.ProjectID,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		AssigneeID:  t.AssigneeID,
		CreatedBy:   t.CreatedBy,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

//

// Create método para crear tarea
//
// @Summary	Crear tarea
// @Description Registra una nueva tarea en el sistema
// @Tags  tasks
// @Accept json
// @Produce json
// @Param request body	createTaskRequest	true "Datos de la tarea"
// @Success 201 {object} taskResponse
// @Failure 400 {object} ProblemDetails
// @Router /tasks [post]
func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createTaskRequest

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

	task, err := h.usercase.CreateTask(r.Context(), req.ProjectID, req.Title, req.Description, "system")
	if err != nil {
		RespondError(w, r, h.logger, err)
		return
	}

	if err := RespondJSON(w, http.StatusCreated, toTaskResponse(task)); err != nil {
		h.logger.ErrorContext(
			r.Context(),
			"error al escribir respuesta JSON de tarea",
			"error", err,
			"method", r.Method,
			"path", r.URL.Path,
		)
	}
}

// Get método para obtener tarea por ID
//
// @Summary	Obtener tarea por ID
// @Description Obtiene una tarea por su ID
// @Tags  tasks
// @Produce json
// @Param id path	string	true "ID de la tarea"
// @Success 200 {object} taskResponse
// @Failure 404 {object} ProblemDetails
// @Router /tasks/{id} [get]
func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	task, err := h.usercase.GetTask(r.Context(), id)

	if err != nil {
		RespondError(w, r, h.logger, err)
		return
	}

	if err := RespondJSON(w, http.StatusOK, toTaskResponse(task)); err != nil {
		h.logger.ErrorContext(r.Context(),
			"error al escribir respuesta JSON",
			"error", err,
			"method", r.Method,
			"path", r.URL.Path)
	}
}

// ListByProject método para obtener tareas por medio de un ProjectID
//
// @Summary	Listar tareas del proyecto
// @Description Retorna todas las tareas de un proyecto
// @Tags  tasks
// @Produce json
// @Param id path	string	true "ID del proyecto"
// @Success 200 {object} PaginatedResponse[taskResponse]
// @Failure 404 {object} ProblemDetails
// @Router /projects/{id}/tasks [get]
func (h *TaskHandler) ListByProject(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	tasks, err := h.usercase.ListProjectTasks(r.Context(), projectID)

	if err != nil {
		RespondError(w, r, h.logger, err)
		return
	}

	response := make([]taskResponse, len(tasks))
	for i, t := range tasks {
		response[i] = toTaskResponse(t)
	}

	if err := RespondJSON(w, http.StatusOK, PaginatedResponse[taskResponse]{
		Data:    response,
		HasMore: false,
	}); err != nil {
		h.logger.ErrorContext(r.Context(),
			"error al obtener respuesta JSON para obtener tareas",
			"error", err,
			"method", r.Method,
			"path", r.URL.Path)
	}

}

// Assign asigna una tarea
//
// @Summary	Asignar tarea
// @Description Asigna una tarea a un miembro del equipo
// @Tags  tasks
// @Accept json
// @Produce json
// @Param id 	path	string		true 	"ID de la tarea"
// @Param request body	assignTaskRequest	true "ID del asignado"
// @Success 200 {object} taskResponse
// @Failure 400 {object} ProblemDetails
// @Router /tasks/{id}/assign [put]
func (h *TaskHandler) Assign(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req assignTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, r, h.logger, domain.ErrInvalidInput)
		return
	}

	if errs := h.validator.Validate(req); errs != nil {
		if err := RespondValidationError(w, errs); err != nil {
			h.logger.ErrorContext(r.Context(),
				"error al escribir respuesta de validación",
				"error", err,
				"method", r.Method,
				"path", r.URL.Path)
		}
		return
	}

	task, err := h.usercase.AssignTask(r.Context(), id, req.AssigneeID)

	if err != nil {
		RespondError(w, r, h.logger, err)
		return
	}

	if err := RespondJSON(w, http.StatusOK, toTaskResponse(task)); err != nil {
		h.logger.ErrorContext(r.Context(),
			"error al obtener respuesta JSON al asignar tarea",
			"error", err,
			"method", r.Method,
			"path", r.URL.Path)
	}
}

// UpdateStatus Actualizar estado de tarea
//
// @Summary	Actualizar estado de tarea
// @Description Mueve una tarea a un nuevo estado del tablero
// @Tags  tasks
// @Accept json
// @Produce json
// @Param id 	path	string		true 	"ID de la tarea"
// @Param request body	updateStatusRequest	true "Nuevo estado"
// @Success 200 {object} taskResponse
// @Failure 400 {object} ProblemDetails
// @Router /tasks/{id}/status [put]
func (h *TaskHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req updateStatusRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, r, h.logger, domain.ErrInvalidInput)
		return
	}

	if errs := h.validator.Validate(req); errs != nil {
		if err := RespondValidationError(w, errs); err != nil {
			h.logger.ErrorContext(r.Context(),
				"error al escribir respuesta de validación",
				"error", err,
				"method", r.Method,
				"path", r.URL.Path)
		}
		return
	}

	task, err := h.usercase.UpdateTaskStatus(r.Context(), id, req.Status)

	if err != nil {
		RespondError(w, r, h.logger, err)
		return
	}

	if err := RespondJSON(w, http.StatusOK, toTaskResponse(task)); err != nil {
		h.logger.ErrorContext(r.Context(),
			"error al obtener respuesta JSON",
			"error", err,
			"method", r.Method,
			"path", r.URL.Path)
	}

}
