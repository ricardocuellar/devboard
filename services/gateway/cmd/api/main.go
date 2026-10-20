package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/ricardocuellar/devboard/services/gateway/internal/config"
	"github.com/ricardocuellar/devboard/services/gateway/internal/middlewares"
	"github.com/ricardocuellar/devboard/services/gateway/internal/router"
)

func main() {
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, nil),
	)

	if err := run(logger); err != nil {
		logger.Error(
			"Gateway se ha detenido",
			slog.Any("error", err),
		)

		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("cargando configuración: %w", err)
	}

	mux := http.NewServeMux()

	router.Register(mux, cfg, logger)

	mux.HandleFunc(
		"GET /health",
		func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set(
				"Content-Type",
				"application/json",
			)

			if err := json.NewEncoder(writer).Encode(
				map[string]string{
					"status":  "ok",
					"service": "gateway",
				},
			); err != nil {
				logger.ErrorContext(
					request.Context(),
					"error escribiendo health",
					slog.Any("error", err),
				)
			}
		},
	)

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           middlewares.RequestID(logger)(mux),
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	logger.Info(
		"iniciando Gateway",
		slog.String("addr", cfg.HTTPAddr),
		slog.String("identity_url", cfg.IdentityURL.String()),
		slog.String("workspace_url", cfg.WorkspaceURL.String()),
		slog.String("task_url", cfg.TaskURL.String()),
	)

	if err := server.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("ejecutando servidor HTTP: %w", err)
	}

	return nil
}
