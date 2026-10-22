package outbox

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	generated "github.com/ricardocuellar/devboard/services/task/internal/repository/postgres/generated"
)

const (
	batchSize       = 20
	pollInterval    = time.Second
	publishTimeout  = 3 * time.Second
	databaseTimeout = 3 * time.Second
	retryDelay      = 5 * time.Second
)

type Publisher interface {
	Publish(
		ctx context.Context,
		eventID string,
		subject string,
		payload []byte,
	) error
}

func Run(
	ctx context.Context,
	queries *generated.Queries,
	publisher Publisher,
	logger *slog.Logger,
) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		if ctx.Err() != nil {
			return
		}

		processBatch(ctx, queries, publisher, logger)

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func processBatch(
	ctx context.Context,
	queries *generated.Queries,
	publisher Publisher,
	logger *slog.Logger,
) {
	queryCtx, cancel := context.WithTimeout(ctx, databaseTimeout)
	events, err := queries.ListPendingOutboxEvents(queryCtx, batchSize)
	cancel()

	if err != nil {
		if ctx.Err() == nil {
			logger.ErrorContext(
				ctx,
				"error consultando Outbox",
				slog.Any("error", err),
			)
		}
		return
	}

	for _, event := range events {
		if ctx.Err() != nil {
			return
		}

		publishCtx, cancel := context.WithTimeout(ctx, publishTimeout)
		publishErr := publisher.Publish(
			publishCtx,
			event.EventID,
			event.Subject,
			event.Payload,
		)
		cancel()

		if ctx.Err() != nil {
			return
		}

		resultCtx, cancel := context.WithTimeout(ctx, databaseTimeout)

		var affected int64

		if publishErr != nil {
			affected, err = queries.RecordOutboxEventFailure(
				resultCtx,
				generated.RecordOutboxEventFailureParams{
					EventID:   event.EventID,
					LastError: publishErr.Error(),
					NextAttemptAt: pgtype.Timestamptz{
						Time:  time.Now().UTC().Add(retryDelay),
						Valid: true,
					},
				},
			)
		} else {
			affected, err = queries.MarkOutboxEventPublished(
				resultCtx,
				event.EventID,
			)
		}

		cancel()

		if ctx.Err() != nil {
			return
		}

		if err != nil {
			logger.ErrorContext(
				ctx,
				"error registrando resultado de publicación",
				slog.String("event_id", event.EventID),
				slog.Any("error", err),
				slog.Any("publish_error", publishErr),
			)
			continue
		}

		if affected == 0 {
			logger.WarnContext(
				ctx,
				"evento no actualizado; puede estar ya publicado",
				slog.String("event_id", event.EventID),
			)
			continue
		}

		if publishErr != nil {
			logger.WarnContext(
				ctx,
				"publicación fallida; evento pendiente de reintento",
				slog.String("event_id", event.EventID),
				slog.Any("error", publishErr),
			)
			continue
		}

		logger.InfoContext(
			ctx,
			"evento publicado y registrado",
			slog.String("event_id", event.EventID),
			slog.String("subject", event.Subject),
		)
	}
}
