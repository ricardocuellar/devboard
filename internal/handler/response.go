package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/ricardocuellar/devboard/internal/domain"
	"github.com/ricardocuellar/devboard/internal/validator"
)

// errorResponse struct
type errorResponse struct {
	Error   string                      `json:"error"`
	Details []validator.ValidationError `json:"details,omitempty"`
}

// RespondError centralizado
func RespondError(writer http.ResponseWriter, request *http.Request, logger *slog.Logger, err error) {
	var status int
	var message string

	switch {
	case errors.Is(err, domain.ErrNotFound):
		status = http.StatusNotFound
		message = err.Error()
	case errors.Is(err, domain.ErrAlreadyExists):
		status = http.StatusConflict
		message = err.Error()
	case errors.Is(err, domain.ErrUnauthorized):
		status = http.StatusUnauthorized
		message = err.Error()
	case errors.Is(err, domain.ErrForbidden):
		status = http.StatusForbidden
		message = err.Error()
	case errors.Is(err, domain.ErrInvalidInput):
		status = http.StatusBadRequest
		message = err.Error()
	default:
		status = http.StatusInternalServerError
		message = "error interno del servidor"
		logger.ErrorContext(request.Context(), "error no clasificado",
			slog.Any("error", err),
			slog.String("path", request.URL.Path))
	}

	respondJSON(writer, status, errorResponse{Error: message})
}

// RespondValidationError responde errores de validación
func RespondValidationError(writer http.ResponseWriter, errs []validator.ValidationError) {
	respondJSON(writer, http.StatusBadRequest, errorResponse{
		Error:   "datos de entrada inválidos",
		Details: errs,
	})
}

// ResponseJSON envia una respuesta exitosa en formato JSON
func ResponseJSON(writer http.ResponseWriter, status int, data any) {
	respondJSON(writer, status, data)
}

func respondJSON(writer http.ResponseWriter, status int, data any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	json.NewEncoder(writer).Encode(data)
}
