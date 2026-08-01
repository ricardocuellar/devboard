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

// CommentRepository implementa las interfaces de comment para los use cases
type CommentRepository struct {
	pool    *pgxpool.Pool
	queries *generated.Queries
}

// NewCommentRepository constructor
func NewCommentRepository(pool *pgxpool.Pool) *CommentRepository {
	return &CommentRepository{
		pool:    pool,
		queries: generated.New(pool),
	}
}

// Create guarda un comentario en PostgreSQL
func (r *CommentRepository) Create(ctx context.Context, comment *domain.Comment) error {
	var taskID, authorID pgtype.UUID

	if err := taskID.Scan(comment.TaskID); err != nil {
		return fmt.Errorf("task_id inválido: %w", domain.ErrInvalidInput)
	}
	if err := authorID.Scan(comment.AuthorID); err != nil {
		return fmt.Errorf("author_id inválido: %w", domain.ErrInvalidInput)
	}

	result, err := r.queries.CreateComment(ctx, generated.CreateCommentParams{
		TaskID:   taskID,
		AuthorID: authorID,
		Content:  comment.Content,
	})

	if err != nil {
		return fmt.Errorf("insertando comentario en BD: %w", err)
	}

	comment.ID = result.ID.String()
	comment.CreatedAt = result.CreatedAt.Time
	comment.UpdatedAt = result.UpdatedAt.Time

	return nil
}

// GetByID obtiene un comentario por ID
func (r *CommentRepository) GetByID(ctx context.Context, id string) (*domain.Comment, error) {
	var pgID pgtype.UUID
	if err := pgID.Scan(id); err != nil {
		return nil, fmt.Errorf("id inválido %s: %w", id, domain.ErrNotFound)
	}

	result, err := r.queries.GetCommentByID(ctx, pgID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("comentario %s : %w", id, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("consultando comentario %s: %w", id, err)
	}

	return toDomainComment(result), nil
}

// ListByTask lista los comentarios de una tarea en orden cronológico
func (r *CommentRepository) ListByTask(ctx context.Context, taskID string) ([]*domain.Comment, error) {
	var pgTaskID pgtype.UUID
	if err := pgTaskID.Scan(taskID); err != nil {
		return nil, fmt.Errorf("task_id inválido: %w", domain.ErrInvalidInput)
	}

	results, err := r.queries.ListCommentsByTask(ctx, pgTaskID)
	if err != nil {
		return nil, fmt.Errorf("listando comentarios de la tarea: %w", err)
	}

	comments := make([]*domain.Comment, len(results))

	for i, result := range results {
		comments[i] = toDomainComment(result)
	}

	return comments, nil

}

// Delete elimina un comentario
func (r *CommentRepository) Delete(ctx context.Context, id string) error {
	var pgID pgtype.UUID
	if err := pgID.Scan(id); err != nil {
		return fmt.Errorf("id inválido: %w", domain.ErrInvalidInput)
	}

	return r.queries.DeleteComment(ctx, pgID)
}

// toDomainComment convierte el modelo sqlc al dominio
func toDomainComment(c generated.Comment) *domain.Comment {
	return &domain.Comment{
		ID:        c.ID.String(),
		TaskID:    c.TaskID.String(),
		AuthorID:  c.AuthorID.String(),
		Content:   c.Content,
		CreatedAt: c.CreatedAt.Time,
		UpdatedAt: c.UpdatedAt.Time,
	}
}
