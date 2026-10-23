package domain

import (
	"errors"
	"time"
)

var ErrInvalidEvent = errors.New("evento inválido")

type Notification struct {
	ID        string
	EventID   string
	UserID    string
	TaskID    string
	ProjectID string
	Message   string
	CreatedAt time.Time
}
