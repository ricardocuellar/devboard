package nats

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	natsgo "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"

	taskv1 "github.com/ricardocuellar/devboard/contracts/events/task/v1"
	"github.com/ricardocuellar/devboard/services/notification/internal/domain"
)

type TaskCreatedHandler interface {
	Handle(
		ctx context.Context,
		event taskv1.TaskCreated,
	) (bool, error)
}

func Run(
	ctx context.Context,
	connection *natsgo.Conn,
	handler TaskCreatedHandler,
	logger *slog.Logger,
) error {
	client, err := jetstream.New(connection)
	if err != nil {
		return fmt.Errorf("creando cliente JetStream: %w", err)
	}

	setupCtx, cancel := context.WithTimeout(ctx, 5*time.Second)

	consumer, err := client.CreateOrUpdateConsumer(
		setupCtx,
		"DEVBOARD_TASKS",
		jetstream.ConsumerConfig{
			Durable:       "notification-task-created-v1",
			FilterSubject: taskv1.TaskCreatedSubject,
			DeliverPolicy: jetstream.DeliverAllPolicy,
			AckPolicy:     jetstream.AckExplicitPolicy,
			AckWait:       15 * time.Second,
			MaxAckPending: 1,
			MaxDeliver:    -1,
		},
	)
	cancel()

	if err != nil {
		return fmt.Errorf("configurando consumidor: %w", err)
	}

	messages, err := consumer.Messages(
		jetstream.PullMaxMessages(1),
	)
	if err != nil {
		return fmt.Errorf("iniciando recepción: %w", err)
	}
	defer messages.Stop()

	logger.Info("consumidor preparado")

	for {
		message, err := messages.Next(jetstream.NextContext(ctx))
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("recibiendo mensaje: %w", err)
		}

		if ctx.Err() != nil {
			return nil
		}

		processMessage(ctx, message, handler, logger)
	}
}

func processMessage(
	ctx context.Context,
	message jetstream.Msg,
	handler TaskCreatedHandler,
	logger *slog.Logger,
) {
	var event taskv1.TaskCreated

	if err := json.Unmarshal(message.Data(), &event); err != nil {
		retryMessage(
			ctx,
			message,
			logger,
			fmt.Errorf("%w: JSON inválido: %v", domain.ErrInvalidEvent, err),
			event.EventID,
		)
		return
	}

	processCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	created, err := handler.Handle(processCtx, event)
	cancel()

	if ctx.Err() != nil {
		return
	}

	if err != nil {
		retryMessage(ctx, message, logger, err, event.EventID)
		return
	}

	ackCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	err = message.DoubleAck(ackCtx)
	cancel()

	if err != nil {
		logger.WarnContext(
			ctx,
			"procesamiento guardado; confirmación no completada",
			slog.String("event_id", event.EventID),
			slog.Any("error", err),
		)
		return
	}

	logger.InfoContext(
		ctx,
		"TaskCreated procesado y confirmado",
		slog.String("event_id", event.EventID),
		slog.String("request_id", event.CorrelationID),
		slog.Bool("notification_created", created),
	)
}

func retryMessage(
	ctx context.Context,
	message jetstream.Msg,
	logger *slog.Logger,
	processErr error,
	eventID string,
) {
	if ctx.Err() != nil {
		return
	}

	logger.WarnContext(
		ctx,
		"evento no procesado; solicitando reentrega",
		slog.String("event_id", eventID),
		slog.Bool("invalid_event", errors.Is(processErr, domain.ErrInvalidEvent)),
		slog.Any("error", processErr),
	)

	if err := message.NakWithDelay(5 * time.Second); err != nil {
		logger.ErrorContext(
			ctx,
			"no se pudo solicitar la reentrega",
			slog.Any("error", err),
		)
	}
}
