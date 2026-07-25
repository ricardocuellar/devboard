package postgres

import (
	"context"
	"fmt"

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
