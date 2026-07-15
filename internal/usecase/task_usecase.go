package usecase

import (
	"context"

	"github.com/ricardocuellar/devboard/internal/domain"
)

// TaskCommandRepository interface CQRS
type TaskCommandRepository interface {
	Create(ctx context.Context, task *domain.Task) error
	Update(ctx context.Context, task *domain.Task) error
}

// TaskQueryRepository interface CQRS
type TaskQueryRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Task, error)
	ListByProject(ctx context.Context, projectID string) ([]*domain.Task, error)
}

// TaskUseCase struct usecase
type TaskUseCase struct {
	commands TaskCommandRepository
	queries  TaskQueryRepository
}

// NewTaskUseCase constructor
func NewTaskUseCase(commands TaskCommandRepository, queries TaskQueryRepository) *TaskUseCase {
	return &TaskUseCase{
		commands: commands,
		queries:  queries,
	}
}
