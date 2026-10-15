package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/ricardocuellar/devboard/services/task/internal/domain"
	"github.com/ricardocuellar/devboard/services/task/internal/middlewares"
	"github.com/ricardocuellar/devboard/services/task/internal/usecase"
	"github.com/ricardocuellar/devboard/services/task/internal/validator"
)

// CommentHandler maneja los endpoints de comentarios
type CommentHandler struct {
	usecase   *usecase.CommentUseCase
	validator *validator.Validator
	logger    *slog.Logger
}

// NewCommentHandler constructor
func NewCommentHandler(uc *usecase.CommentUseCase, v *validator.Validator, logger *slog.Logger) *CommentHandler {
	return &CommentHandler{usecase: uc, validator: v, logger: logger}
}

type createCommentRequest struct {
	Content string `json:"content" validate:"required,min=1,max=5000"`
}

type commentResponse struct {
	ID        string    `json:"id"`
	TaskID    string    `json:"task_id"`
	AuthorID  string    `json:"author_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toCommentResponse(c *domain.Comment) commentResponse {
	return commentResponse{
		ID:        c.ID,
		TaskID:    c.TaskID,
		AuthorID:  c.AuthorID,
		Content:   c.Content,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

// Create agrega un comentario a una tarea
//
// @Summary     Crear comentario
// @Description Agrega un comentario a una tarea
// @Failure 401 {object} ProblemDetails
// @Failure 403 {object} ProblemDetails
// @Failure 503 {object} ProblemDetails
// @Failure 504 {object} ProblemDetails
// @Tags        comments
// @Accept      json
// @Produce     json
// @Param       id      path     string               true "ID de la tarea"
// @Param       request body     createCommentRequest true "Contenido"
// @Success     201     {object} commentResponse
// @Security    BearerAuth
// @Router      /tasks/{id}/comments [post]
func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("id")
	authorID := middlewares.GetUserID(r)

	if authorID == "" {
		RespondError(w, r, h.logger, domain.ErrUnauthorized)
		return
	}

	var req createCommentRequest
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

	comment, err := h.usecase.CreateComment(r.Context(), taskID, authorID, req.Content)

	if err != nil {
		RespondError(w, r, h.logger, err)
		return
	}

	if err := RespondJSON(w, http.StatusCreated, toCommentResponse(comment)); err != nil {
		h.logger.ErrorContext(r.Context(), "error escribiendo respuesta JSON",
			"error", err, "path", r.URL.Path)
	}
}

// ListByTask lista los comentarios de una tarea
//
// @Summary     Listar comentarios
// @Description Lista todos los comentarios de una tarea
// @Failure 401 {object} ProblemDetails
// @Failure 403 {object} ProblemDetails
// @Failure 503 {object} ProblemDetails
// @Failure 504 {object} ProblemDetails
// @Tags        comments
// @Produce     json
// @Param       id  path     string true "ID de la tarea"
// @Success     200 {object} PaginatedResponse[commentResponse]
// @Security    BearerAuth
// @Router      /tasks/{id}/comments [get]
func (h *CommentHandler) ListByTask(w http.ResponseWriter, r *http.Request) {
	requesterID := middlewares.GetUserID(r)

	if requesterID == "" {
		RespondError(w, r, h.logger, domain.ErrUnauthorized)
		return
	}

	taskID := r.PathValue("id")

	comments, err := h.usecase.ListTaskComments(r.Context(), taskID, requesterID)

	if err != nil {
		RespondError(w, r, h.logger, err)
		return
	}

	response := make([]commentResponse, len(comments))
	for i, c := range comments {
		response[i] = toCommentResponse(c)
	}

	if err := RespondJSON(w, http.StatusOK, PaginatedResponse[commentResponse]{
		Data:    response,
		HasMore: false,
	}); err != nil {
		h.logger.ErrorContext(r.Context(), "error escribiendo respuesta JSON",
			"error", err, "path", r.URL.Path)
	}
}

// Delete elimina un comentario
//
// @Summary     Eliminar comentario
// @Description Elimina un comentario. Solo el author puede borrarlo.
// @Tags        comments
// @Param       id path string true "ID del comentario"
// @Success     204 "Sin contenido"
// @Failure     403 {object} ProblemDetails
// @Security    BearerAuth
// @Router      /comments/{id} [delete]
func (h *CommentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	commentID := r.PathValue("id")
	requesterID := middlewares.GetUserID(r)

	if requesterID == "" {
		RespondError(w, r, h.logger, domain.ErrUnauthorized)
		return
	}

	if err := h.usecase.DeleteComment(r.Context(), commentID, requesterID); err != nil {
		RespondError(w, r, h.logger, err)
		return

	}

	w.WriteHeader(http.StatusNoContent)
}
