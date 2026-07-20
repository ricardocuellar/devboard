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
	tasks, err := uc.queries.ListByProject(ctx, projectID)

	if err != nil {
		return nil, fmt.Errorf("listando tareas del proyecto %s: %w", projectID, err)
	}

	return tasks, nil
}

// AssignTask método para asignar tarea
func (uc *TaskUseCase) AssignTask(ctx context.Context, taskID, userID string) (*domain.Task, error) {
	task, err := uc.queries.GetByID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("obteniendo tarea %s: %w", taskID, err)
	}

	if err := task.AssignTo(userID); err != nil {
		return nil, err
	}

	if err := uc.commands.Update(ctx, task); err != nil {
		return nil, fmt.Errorf("actualizando tarea: %w", err)
	}

	return task, nil
}

// GetTask obtiene una tarea por ID
func (uc *TaskUseCase) GetTask(ctx context.Context, id string) (*domain.Task, error) {
	task, err := uc.queries.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("obteniendo tarea %s: %w", id, err)
	}

	return task, nil
}

// UpdateTaskStatus mueve una tarea a un nuevo estado del tablero
func (uc *TaskUseCase) UpdateTaskStatus(ctx context.Context, taskID string, newStatus domain.TaskStatus) (*domain.Task, error) {
	task, err := uc.queries.GetByID(ctx, taskID)

	if err != nil {
		return nil, fmt.Errorf("obteniendo tarea %s: %w", taskID, err)
	}

	if err := task.MoveToStatus(newStatus); err != nil {
		return nil, err
	}

	if err := uc.commands.Update(ctx, task); err != nil {
		return nil, fmt.Errorf("actualizando estado de tarea: %w", err)
	}

	return task, nil
}
