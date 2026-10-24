package usecase

import (
	"context"
	"fmt"

	"github.com/ricardocuellar/devboard/services/notification/internal/domain"
)

type NotificationQueryRepository interface {
	ListByUser(
		ctx context.Context,
		userID string,
		limit int32,
	) ([]*domain.Notification, error)
}

type NotificationQueryUseCase struct {
	repository NotificationQueryRepository
}

func NewNotificationQueryUseCase(
	repository NotificationQueryRepository,
) *NotificationQueryUseCase {
	return &NotificationQueryUseCase{
		repository: repository,
	}
}

func (uc *NotificationQueryUseCase) ListMine(
	ctx context.Context,
	requesterID string,
	limit int32,
) ([]*domain.Notification, error) {
	if requesterID == "" {
		return nil, domain.ErrUnauthorized
	}

	if limit < 1 || limit > 100 {
		return nil, fmt.Errorf(
			"limit debe estar entre 1 y 100: %w",
			domain.ErrInvalidInput,
		)
	}

	notifications, err := uc.repository.ListByUser(
		ctx,
		requesterID,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"consultando notificaciones del usuario: %w",
			err,
		)
	}

	return notifications, nil
}
