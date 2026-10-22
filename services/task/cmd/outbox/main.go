package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	natsgo "github.com/nats-io/nats.go"

	natsadapter "github.com/ricardocuellar/devboard/services/task/internal/messaging/nats"
	"github.com/ricardocuellar/devboard/services/task/internal/outbox"
	"github.com/ricardocuellar/devboard/services/task/internal/repository/postgres"
	generated "github.com/ricardocuellar/devboard/services/task/internal/repository/postgres/generated"
)

func main() {
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, nil),
	).With(slog.String("service", "task-outbox"))

	if err := run(logger); err != nil {
		logger.Error("publicador detenido con error", slog.Any("error", err))
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	databaseURL := os.Getenv("TASK_DATABASE_URL")
	natsURL := os.Getenv("NATS_URL")

	if databaseURL == "" || natsURL == "" {
		return fmt.Errorf("TASK_DATABASE_URL y NATS_URL son requeridas")
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	startupCtx, cancel := context.WithTimeout(ctx, 10*time.Second)

	pool, err := postgres.NewPool(
		startupCtx,
		postgres.Config{
			URL:             databaseURL,
			MaxConns:        3,
			MinConns:        1,
			MaxConnLifetime: time.Hour,
			MaxConnIdleTime: 5 * time.Minute,
		},
	)
	cancel()

	if err != nil {
		return err
	}
	defer pool.Close()

	connection, err := natsgo.Connect(
		natsURL,
		natsgo.Name("task-outbox"),
		natsgo.Timeout(3*time.Second),
		natsgo.MaxReconnects(-1),
		natsgo.ReconnectWait(time.Second),
		natsgo.ReconnectBufSize(-1),
	)
	if err != nil {
		return fmt.Errorf("conectando con NATS: %w", err)
	}
	defer connection.Close()

	publisher, err := natsadapter.NewPublisher(
		connection,
		"DEVBOARD_TASKS",
	)
	if err != nil {
		return err
	}

	logger.Info("publicador iniciado")

	outbox.Run(
		ctx,
		generated.New(pool),
		publisher,
		logger,
	)

	logger.Info("publicador detenido")

	return nil
}
