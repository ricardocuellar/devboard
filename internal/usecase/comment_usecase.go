package usecase

import (
	"context"
	"fmt"

	"github.com/ricardocuellar/devboard/internal/domain"
)

// CommentRepository define lo que CommentUseCase necesita
type CommentRepository interface {
	Create(ctx context.Context, comment *domain.Comment) error
	GetByID(ctx context.Context, id string) (*domain.Comment, error)
	ListByTask(ctx context.Context, taskID string) ([]*domain.Comment, error)
	Delete(ctx context.Context, id string) error
}

// CommentUseCase orquesta las operaciones de comentarios
type CommentUseCase struct {
	commentRepo CommentRepository
	taskQuery   TaskQueryRepository
}

// NewCommentUseCase constructor
func NewCommentUseCase(commentRepo CommentRepository, taskQuery TaskQueryRepository) *CommentUseCase {
	return &CommentUseCase{
		commentRepo: commentRepo,
		taskQuery:   taskQuery,
	}
}

// CreateComment agrega un comentario a una tarea.
func (uc *CommentUseCase) CreateComment(ctx context.Context, taskID, authorID, content string) (*domain.Comment, error) {
	_, err := uc.taskQuery.GetByID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("tarea no encontrada: %w", err)
	}

	comment, err := domain.NewComment(taskID, authorID, content)

	if err != nil {
		return nil, fmt.Errorf("creando comentario: %w", err)
	}

	if err := uc.commentRepo.Create(ctx, comment); err != nil {
		return nil, fmt.Errorf("guardando comentario: %w", err)
	}

	return comment, nil
}

// ListTaskComments lista los comentarios de una tarea
func (uc *CommentUseCase) ListTaskComments(ctx context.Context, taskID string) ([]*domain.Comment, error) {
	comments, err := uc.commentRepo.ListByTask(ctx, taskID)

	if err != nil {
		return nil, fmt.Errorf("listando comentarios: %w", err)
	}

	return comments, nil
}

// DeleteComment elimina un comentario.
func (uc *CommentUseCase) DeleteComment(ctx context.Context, commentID, requesterID string) error {
	comment, err := uc.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return fmt.Errorf("obteniendo comentario: %w", err)
	}

	if comment.AuthorID != requesterID {
		return fmt.Errorf("solo el author puede borrar su comentario: %w", domain.ErrForbidden)
	}

	return uc.commentRepo.Delete(ctx, commentID)
}
