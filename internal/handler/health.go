// Package handler health
package handler

import (
	"encoding/json"
	"net/http"
)

// HealthHandler struct
type HealthHandler struct{}

// NewHealthHandler constructor
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// healthResponse struct
type healthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// HealthCheck verifica que el servidor este funcionando
//
// @Summary	Health check
// @Description Verifica que el servidor está corriendo y respondiendo
// @Tags  system
// @Produce json
// @Success 200 {object} healthResponse
// @Router /health [get]
func (handle *HealthHandler) ServeHTTP(writer http.ResponseWriter, _ *http.Request) {
	resp := healthResponse{
		Status:  "ok",
		Version: "1.0.0",
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK) //opcional

	if err := json.NewEncoder(writer).Encode(resp); err != nil {
		return
	}
}
