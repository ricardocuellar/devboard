package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ricardocuellar/devboard/internal/domain"
	generated "github.com/ricardocuellar/devboard/internal/repository/postgres/generated"
)

// TaskRepository implementa usecaseTaskCommandRepository y usecase.TaskQueryRepository
type TaskRepository struct {
	pool    *pgxpool.Pool
	queries *generated.Queries
}

// NewTaskRepository constructor
func NewTaskRepository(pool *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{
		pool:    pool,
		queries: generated.New(pool),
	}
}

// Create guarda una nueva tarea en PostgreSQL
func (r *TaskRepository) Create(ctx context.Context, task *domain.Task) error {
	var projectID pgtype.UUID
	if err := projectID.Scan(task.ProjectID); err != nil {
		return fmt.Errorf("project_id inválido: %w", domain.ErrInvalidInput)
	}

	var createdBy pgtype.UUID
	if err := createdBy.Scan(task.CreatedBy); err != nil {
		return fmt.Errorf("created_by inválido: %w", domain.ErrInvalidInput)
	}

	result, err := r.queries.CreateTask(ctx, generated.CreateTaskParams{
		ProjectID:   projectID,
		Title:       task.Title,
		Description: task.Description,
		Status:      string(task.Status),
		CreatedBy:   createdBy,
	})

	if err != nil {
		return fmt.Errorf("instertando tarea en BD: %w", err)
	}

	task.ID = result.ID.String()
	task.CreatedAt = result.CreatedAt.Time
	task.UpdatedAt = result.UpdatedAt.Time

	return nil

}

// Update actualiza una tarea existente
func (r *TaskRepository) Update(ctx context.Context, task *domain.Task) error {
	var id pgtype.UUID
	if err := id.Scan(task.ID); err != nil {
		return fmt.Errorf("id inválido: %w", domain.ErrInvalidInput)
	}

	var assigneeID pgtype.UUID

	if task.AssigneeID != nil {
		if err := assigneeID.Scan(task.AssigneeID); err != nil {
			return fmt.Errorf("assignee_id inválido: %w", domain.ErrInvalidInput)
		}
	}

	_, err := r.queries.UpdateTask(ctx, generated.UpdateTaskParams{
		ID:          id,
		Title:       task.Title,
		Description: task.Description,
		Status:      string(task.Status),
		AssigneeID:  assigneeID,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("tarea %s: %w", task.ID, domain.ErrNotFound)
		}
		return fmt.Errorf("actualizando tarea %s: %w", task.ID, err)
	}

	return nil

}
