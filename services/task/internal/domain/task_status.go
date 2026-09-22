// Package domain task status
package domain

import (
	"encoding/json"
	"fmt"
)

// TaskStatus representa el estado de una tarea
type TaskStatus string

const (
	// TaskStatusTodo representa una tarea no iniciada
	TaskStatusTodo TaskStatus = "todo"
	// TaskStatusInProgress representa una tarea en proceso
	TaskStatusInProgress TaskStatus = "in_progress"
	// TaskStatusDone representa una tarea finalizada
	TaskStatusDone TaskStatus = "done"
)

// IsValid verifica que el status sea uno de los permitidos
func (status TaskStatus) IsValid() bool {
	switch status {
	case TaskStatusTodo, TaskStatusInProgress, TaskStatusDone:
		return true
	default:
		return false
	}
}

// UnmarshalJSON valida el estatus al deserializar
func (status *TaskStatus) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	statusTask := TaskStatus(raw)
	if !statusTask.IsValid() {
		return fmt.Errorf("estado de tarea inválido: %q", raw)
	}

	*status = statusTask
	return nil

}
