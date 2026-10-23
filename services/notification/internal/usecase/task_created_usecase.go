package usecase

import (
	"context"
	"fmt"
	"strings"

	taskv1 "github.com/ricardocuellar/devboard/contracts/events/task/v1"
	"github.com/ricardocuellar/devboard/services/notification/internal/domain"
)

type NotificationRepository interface {
	CreateOnce(
		ctx context.Context,
		notification *domain.Notification,
	) (bool, error)
}

type TaskCreatedUseCase struct {
	repository NotificationRepository
}

func NewTaskCreatedUseCase(
	repository NotificationRepository,
) *TaskCreatedUseCase {
	return &TaskCreatedUseCase{
		repository: repository,
	}
}

func (uc *TaskCreatedUseCase) Handle(
	ctx context.Context,
	event taskv1.TaskCreated,
) (bool, error) {
	if event.EventID == "" ||
		len(event.EventID) > 128 ||
		event.Type != taskv1.TaskCreatedType ||
		event.Version != taskv1.TaskCreatedVersion ||
		event.Source != "task" ||
		event.OccurredAt.IsZero() ||
		strings.TrimSpace(event.Data.Title) == "" {
		return false, fmt.Errorf(
			"contrato TaskCreated no válido: %w",
			domain.ErrInvalidEvent,
		)
	}

	notification := &domain.Notification{
		EventID:   event.EventID,
		UserID:    event.Data.CreatedBy,
		TaskID:    event.Data.TaskID,
		ProjectID: event.Data.ProjectID,
		Message: fmt.Sprintf(
			"Se creó la tarea: %s",
			event.Data.Title,
		),
	}

	created, err := uc.repository.CreateOnce(ctx, notification)
	if err != nil {
		return false, fmt.Errorf(
			"procesando TaskCreated: %w",
			err,
		)
	}

	return created, nil
}
