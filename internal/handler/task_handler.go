package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

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
	ProjectID string `json:"project_id" validate:"required"`
	Title     string `json:"title" validate:"required,min=3,max=200"`
}

// createTaskResponse estructura
type createTaskResponse struct {
	ProjectID string `json:"project_id"`
	Title     string `json:"title"`
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
// @Success 201 {object} createTaskResponse
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

	task, err := h.usercase.CreateTask(r.Context(), req.ProjectID, req.Title, "system")
	if err != nil {
		RespondError(w, r, h.logger, err)
		return
	}

	if err := RespondJSON(w, http.StatusCreated, createTaskResponse{
		ProjectID: task.ProjectID,
		Title:     task.Title,
	}); err != nil {
		h.logger.ErrorContext(
			r.Context(),
			"error al escribir respuesta JSON de tarea",
			"error", err,
			"method", r.Method,
			"path", r.URL.Path,
		)
	}
}

// Get

// ListByProject

// Assign

// UpdateStatus
