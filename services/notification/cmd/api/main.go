package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ricardocuellar/devboard/services/notification/internal/auth"
	"github.com/ricardocuellar/devboard/services/notification/internal/handler"
	"github.com/ricardocuellar/devboard/services/notification/internal/middlewares"
	"github.com/ricardocuellar/devboard/services/notification/internal/repository/postgres"
	"github.com/ricardocuellar/devboard/services/notification/internal/usecase"
)

func main() {
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, nil),
	).With(slog.String("service", "notification-http"))

	if err := run(logger); err != nil {
		logger.Error("API detenida con error", slog.Any("error", err))
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	databaseURL := os.Getenv("NOTIFICATION_DATABASE_URL")
	secret := os.Getenv("JWT_SECRET")
	addr := os.Getenv("NOTIFICATION_HTTP_ADDR")

	if addr == "" {
		addr = ":8084"
	}

	if databaseURL == "" {
		return fmt.Errorf("NOTIFICATION_DATABASE_URL es requerida")
	}

	verifier, err := auth.NewTokenVerifier([]byte(secret))
	if err != nil {
		return err
	}

	poolConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return fmt.Errorf("configurando PostgreSQL: %w", err)
	}
	poolConfig.MaxConns = 5

	startupCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	pool, err := pgxpool.NewWithConfig(startupCtx, poolConfig)
	if err != nil {
		cancel()
		return fmt.Errorf("creando pool: %w", err)
	}
	defer pool.Close()

	err = pool.Ping(startupCtx)
	cancel()

	if err != nil {
		return fmt.Errorf("conectando con PostgreSQL: %w", err)
	}

	repository := postgres.NewNotificationRepository(pool)
	queryUseCase := usecase.NewNotificationQueryUseCase(repository)

	mux := http.NewServeMux()

	mux.HandleFunc(
		"GET /health",
		func(writer http.ResponseWriter, _ *http.Request) {
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write(
				[]byte(`{"status":"ok","service":"notification"}`),
			)
		},
	)

	requireAuth := middlewares.RequireAuth(verifier, logger)

	mux.Handle(
		"GET /api/v1/notifications",
		requireAuth(handler.ListNotifications(queryUseCase, logger)),
	)

	server := &http.Server{
		Addr:              addr,
		Handler:           middlewares.Logger(logger)(mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	logger.Info("iniciando Notification HTTP", slog.String("addr", addr))

	return serve(server, logger)
}
