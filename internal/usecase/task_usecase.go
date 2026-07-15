package usecase

import (
	"context"
	"fmt"

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

// CreateTask método para crear tarea
func (uc *TaskUseCase) CreateTask(ctx context.Context, projectID, title, createdBy string) (*domain.Task, error) {
	task := domain.NewTask(projectID, title, createdBy)

	if err := uc.commands.Create(ctx, task); err != nil {
		return nil, fmt.Errorf("creando tarea: %w", err)
	}

	return task, nil
}

// ListProjectTasks método para crear tarea
func (uc *TaskUseCase) ListProjectTasks(ctx context.Context, projectID string) ([]*domain.Task, error) {
	return uc.queries.ListByProject(ctx, projectID)
}
