package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/ricardocuellar/devboard/services/task/internal/domain"
	"github.com/ricardocuellar/devboard/services/task/internal/validator"
)

// errorResponse struct
type errorResponse struct {
	Error   string                      `json:"error"`
	Details []validator.ValidationError `json:"details,omitempty"`
}

// ProblemDetails struct generalizado usando rfc 7807
type ProblemDetails struct {
	Type     string                      `json:"type"`
	Title    string                      `json:"title"`
	Status   int                         `json:"status"`
	Detail   string                      `json:"detail"`
	Instance string                      `json:"instance"`
	Errors   []validator.ValidationError `json:"errors,omitempty"`
}

const problemBaseURL = "https://ricardocuellar.dev/errors"

// RespondError centralizado
func RespondError(writer http.ResponseWriter, request *http.Request, logger *slog.Logger, err error) {
	var problem ProblemDetails
	problem.Instance = request.URL.Path

	switch {
	case errors.Is(err, domain.ErrDependencyUnavailable):
		problem.Type = problemBaseURL + "/not-found"
		problem.Title = "Recurso no disponible"
		problem.Status = http.StatusServiceUnavailable
		problem.Detail = err.Error()
	case errors.Is(err, domain.ErrDependencyTimeout):
		problem.Type = problemBaseURL + "/not-found"
		problem.Title = "Recurso tiempo de respuesta expirado"
		problem.Status = http.StatusGatewayTimeout
		problem.Detail = err.Error()
	case errors.Is(err, domain.ErrNotFound):
		problem.Type = problemBaseURL + "/not-found"
		problem.Title = "Recurso no encontrado"
		problem.Status = http.StatusNotFound
		problem.Detail = err.Error()
	case errors.Is(err, domain.ErrAlreadyExists):
		problem.Type = problemBaseURL + "/conflict"
		problem.Title = "El recurso ya existe"
		problem.Status = http.StatusConflict
		problem.Detail = err.Error()
	case errors.Is(err, domain.ErrUnauthorized):
		problem.Type = problemBaseURL + "/unauthorized"
		problem.Title = "No autorizado"
		problem.Status = http.StatusUnauthorized
		problem.Detail = err.Error()
	case errors.Is(err, domain.ErrForbidden):
		problem.Type = problemBaseURL + "/forbidden"
		problem.Title = "Recurso prohíbido"
		problem.Status = http.StatusForbidden
		problem.Detail = err.Error()
	case errors.Is(err, domain.ErrInvalidInput):
		problem.Type = problemBaseURL + "/bad-request"
		problem.Title = "Solicitud no válida"
		problem.Status = http.StatusBadRequest
		problem.Detail = err.Error()
	default:
		problem.Type = problemBaseURL + "/internal"
		problem.Title = "Error interno del servidor"
		problem.Status = http.StatusInternalServerError
		problem.Detail = "ha ocurrido un error inesperado"
		logger.ErrorContext(request.Context(), "error no clasificado",
			slog.Any("error", err),
			slog.String("path", problem.Instance))
	}

	writer.Header().Set("Content-Type", "application/problem+json")
	writer.WriteHeader(problem.Status)

	if encodeErr := json.NewEncoder(writer).Encode(problem); encodeErr != nil {
		logger.Error("error al escribir problem details",
			slog.String("error", encodeErr.Error()))
	}
}

// RespondValidationError responde errores de validación
func RespondValidationError(writer http.ResponseWriter, errs []validator.ValidationError) error {
	return respondJSON(writer, http.StatusBadRequest, errorResponse{
		Error:   "datos de entrada inválidos",
		Details: errs,
	})
}

// RespondJSON envia una respuesta exitosa en formato JSON
func RespondJSON(writer http.ResponseWriter, status int, data any) error {
	return respondJSON(writer, status, data)
}

func respondJSON(writer http.ResponseWriter, status int, data any) error {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	return json.NewEncoder(writer).Encode(data)
}
