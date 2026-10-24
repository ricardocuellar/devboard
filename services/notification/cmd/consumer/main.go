package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	natsgo "github.com/nats-io/nats.go"

	natsadapter "github.com/ricardocuellar/devboard/services/notification/internal/messaging/nats"
	"github.com/ricardocuellar/devboard/services/notification/internal/repository/postgres"
	"github.com/ricardocuellar/devboard/services/notification/internal/usecase"
)

func main() {
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, nil),
	).With(slog.String("service", "notification-consumer"))

	if err := run(logger); err != nil {
		logger.Error("consumidor detenido con error", slog.Any("error", err))
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	databaseURL := os.Getenv("NOTIFICATION_DATABASE_URL")
	natsURL := os.Getenv("NATS_URL")

	if databaseURL == "" || natsURL == "" {
		return fmt.Errorf("NOTIFICATION_DATABASE_URL y NATS_URL son requeridas")
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	poolConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return fmt.Errorf("configurando PostgreSQL: %w", err)
	}
	poolConfig.MaxConns = 3

	startupCtx, cancel := context.WithTimeout(ctx, 10*time.Second)

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

	connection, err := natsgo.Connect(
		natsURL,
		natsgo.Name("notification-consumer"),
		natsgo.Timeout(3*time.Second),
		natsgo.MaxReconnects(-1),
		natsgo.ReconnectWait(time.Second),
	)
	if err != nil {
		return fmt.Errorf("conectando con NATS: %w", err)
	}
	defer connection.Close()

	repository := postgres.NewNotificationRepository(pool)
	handler := usecase.NewTaskCreatedUseCase(repository)

	err = natsadapter.Run(ctx, connection, handler, logger)
	if err != nil {
		return err
	}

	logger.Info("consumidor detenido")
	return nil
}
