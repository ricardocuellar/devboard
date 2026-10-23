package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ricardocuellar/devboard/services/notification/internal/domain"
	generated "github.com/ricardocuellar/devboard/services/notification/internal/repository/postgres/generated"
)

type NotificationRepository struct {
	pool    *pgxpool.Pool
	queries *generated.Queries
}

func NewNotificationRepository(
	pool *pgxpool.Pool,
) *NotificationRepository {
	return &NotificationRepository{
		pool:    pool,
		queries: generated.New(pool),
	}
}

func (r *NotificationRepository) CreateOnce(
	ctx context.Context,
	notification *domain.Notification,
) (bool, error) {
	userID, err := parseEventUUID(notification.UserID, "created_by")
	if err != nil {
		return false, err
	}

	taskID, err := parseEventUUID(notification.TaskID, "task_id")
	if err != nil {
		return false, err
	}

	projectID, err := parseEventUUID(notification.ProjectID, "project_id")
	if err != nil {
		return false, err
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return false, fmt.Errorf("iniciando procesamiento: %w", err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	queries := r.queries.WithTx(tx)

	affected, err := queries.RegisterProcessedEvent(
		ctx,
		notification.EventID,
	)
	if err != nil {
		return false, fmt.Errorf("registrando evento procesado: %w", err)
	}

	if affected == 0 {
		return false, nil
	}

	row, err := queries.CreateNotification(
		ctx,
		generated.CreateNotificationParams{
			EventID:   notification.EventID,
			UserID:    userID,
			TaskID:    taskID,
			ProjectID: projectID,
			Message:   notification.Message,
		},
	)
	if err != nil {
		return false, fmt.Errorf("insertando notificación: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return false, fmt.Errorf("confirmando procesamiento: %w", err)
	}

	notification.ID = row.ID.String()
	notification.CreatedAt = row.CreatedAt.Time

	return true, nil
}

func parseEventUUID(value, field string) (pgtype.UUID, error) {
	var id pgtype.UUID

	if err := id.Scan(value); err != nil || !id.Valid {
		return pgtype.UUID{}, fmt.Errorf(
			"%s debe ser un UUID válido: %w",
			field,
			domain.ErrInvalidEvent,
		)
	}

	return id, nil
}
