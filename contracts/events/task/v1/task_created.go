package taskv1

import (
	"crypto/rand"
	"time"
)

const (
	TaskCreatedType    = "task.created"
	TaskCreatedVersion = 1
	TaskCreatedSubject = "devboard.task.created.v1"
)

type TaskCreated struct {
	EventID       string          `json:"event_id"`
	Type          string          `json:"type"`
	Version       int             `json:"version"`
	OccurredAt    time.Time       `json:"occurred_at"`
	Source        string          `json:"source"`
	CorrelationID string          `json:"correlation_id,omitempty"`
	Data          TaskCreatedData `json:"data"`
}

type TaskCreatedData struct {
	TaskID    string `json:"task_id"`
	ProjectID string `json:"project_id"`
	CreatedBy string `json:"created_by"`
	Title     string `json:"title"`
}

func NewTaskCreated(
	data TaskCreatedData,
	occurredAt time.Time,
	correlationID string,
) TaskCreated {
	return TaskCreated{
		EventID:       "evt_" + rand.Text(),
		Type:          TaskCreatedType,
		Version:       TaskCreatedVersion,
		OccurredAt:    occurredAt.UTC(),
		Source:        "task",
		CorrelationID: correlationID,
		Data:          data,
	}
}
