package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ricardocuellar/devboard/internal/validator"
)

// errorResponse struct
type errorResponse struct {
	Error   string                      `json:"error"`
	Details []validator.ValidationError `json:"details,omitempty"`
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
