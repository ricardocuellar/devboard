package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	taskv1 "github.com/ricardocuellar/devboard/contracts/events/task/v1"
	"github.com/ricardocuellar/devboard/services/task/internal/domain"
	generated "github.com/ricardocuellar/devboard/services/task/internal/repository/postgres/generated"
	"github.com/ricardocuellar/devboard/services/task/internal/requestid"
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
func (r *TaskRepository) CreateWithEvent(ctx context.Context, task *domain.Task) error {
	var projectID pgtype.UUID
	if err := projectID.Scan(task.ProjectID); err != nil {
		return fmt.Errorf(
			"project_id inválido: %w",
			domain.ErrInvalidInput,
		)
	}

	var createdBy pgtype.UUID
	if err := createdBy.Scan(task.CreatedBy); err != nil {
		return fmt.Errorf("created_by inválido: %w", domain.ErrInvalidInput)
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("iniciando transacción de tarea: %w", err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	queries := r.queries.WithTx(tx)

	row, err := queries.CreateTask(
		ctx,
		generated.CreateTaskParams{
			ProjectID:   projectID,
			Title:       task.Title,
			Description: task.Description,
			Status:      string(task.Status),
			CreatedBy:   createdBy,
		},
	)
	if err != nil {
		return fmt.Errorf("insertando tarea: %w", err)
	}

	event := taskv1.NewTaskCreated(
		taskv1.TaskCreatedData{
			TaskID:    row.ID.String(),
			ProjectID: row.ProjectID.String(),
			CreatedBy: row.CreatedBy.String(),
			Title:     row.Title,
		},
		row.CreatedAt.Time,
		requestid.FromContext(ctx),
	)

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("serializando TaskCreated: %w", err)
	}

	if err := queries.InsertOutboxEvent(
		ctx,
		generated.InsertOutboxEventParams{
			EventID: event.EventID,
			Subject: taskv1.TaskCreatedSubject,
			Payload: payload,
		},
	); err != nil {
		return fmt.Errorf("guardando TaskCreated en Outbox: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("confirmando tarea y evento: %w", err)
	}

	task.ID = row.ID.String()
	task.CreatedAt = row.CreatedAt.Time
	task.UpdatedAt = row.UpdatedAt.Time

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
		if err := assigneeID.Scan(*task.AssigneeID); err != nil {
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

// GetByID obtiene una tarea por ID.
func (r *TaskRepository) GetByID(ctx context.Context, id string) (*domain.Task, error) {
	var pgID pgtype.UUID

	if err := pgID.Scan(id); err != nil {
		return nil, fmt.Errorf("id inválido %s: %w", id, domain.ErrNotFound)
	}

	result, err := r.queries.GetTaskByID(ctx, pgID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("tarea %s: %w", id, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("consultando tarea %s: %w", id, err)
	}

	return toDomainTask(result), nil
}

// ListByProject lista las tareas de un proyecto.
func (r *TaskRepository) ListByProject(ctx context.Context, projectID string) ([]*domain.Task, error) {
	var pgProjectID pgtype.UUID

	if err := pgProjectID.Scan(projectID); err != nil {
		return nil, fmt.Errorf("project_id inválido: %w", domain.ErrNotFound)
	}

	results, err := r.queries.ListTasksByProject(ctx, pgProjectID)

	if err != nil {
		return nil, fmt.Errorf("listanto tareas del proyecto %s: %w", projectID, err)
	}

	tasks := make([]*domain.Task, len(results))

	for i, result := range results {
		tasks[i] = toDomainTask(result)
	}

	return tasks, nil
}

// toDomainTask convierte el modelo de sqlc a modelo de dominio
func toDomainTask(t generated.Task) *domain.Task {
	task := &domain.Task{
		ID:          t.ID.String(),
		ProjectID:   t.ProjectID.String(),
		Title:       t.Title,
		Description: t.Description,
		Status:      domain.TaskStatus(t.Status),
		CreatedBy:   t.CreatedBy.String(),
		CreatedAt:   t.CreatedAt.Time,
		UpdatedAt:   t.UpdatedAt.Time,
	}

	if t.AssigneeID.Valid {
		id := t.AssigneeID.String()
		task.AssigneeID = &id
	}

	return task
}
