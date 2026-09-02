package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/ricardocuellar/devboard/services/identity/internal/config"
	"github.com/ricardocuellar/devboard/services/identity/internal/server"
)

func main() {
	if err := run(); err != nil {
		slog.Error("Identity service se ha detenido", slog.Any("error", err))
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	srv := server.New(cfg.HTTPAddr,
		server.WithLogger(logger),
		server.WithReadTimeout(15*time.Second),
		server.WithWriteTimeout(30*time.Second),
	)

	srv.RegisterRoutes("GET /health", http.HandlerFunc(
		func(writer http.ResponseWriter, _ *http.Request) {
			writer.Header().Set(
				"Content-Type", "application/json",
			)

			writer.WriteHeader(
				http.StatusOK,
			)

			_, _ = writer.Write([]byte(
				`{"status": "ok", "service":"identity"}`,
			))
		},
	))

	logger.Info("Se inicio IDENTITY SERVICE",
		slog.String("addr", cfg.HTTPAddr),
	)

	return srv.Start()
}
