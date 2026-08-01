package domain

import (
	"errors"
	"time"
)

// ErrCommentEmpty indica que el contenido del comentario está vacío
var ErrCommentEmpty = errors.New("el comentario no puede estar vacío")

// Comment representa un comentario en una tarea
type Comment struct {
	ID        string
	TaskID    string
	AuthorID  string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewComment construye un Comment válido.
func NewComment(taskID, authorID, content string) (*Comment, error) {
	if content == "" {
		return nil, ErrCommentEmpty
	}

	now := time.Now()
	return &Comment{
		TaskID:    taskID,
		AuthorID:  authorID,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
