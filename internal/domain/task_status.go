package domain

import (
	"encoding/json"
	"fmt"
)

type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusDone       TaskStatus = "done"
)

func (status TaskStatus) IsValid() bool {
	switch status {
	case TaskStatusTodo, TaskStatusInProgress, TaskStatusDone:
		return true
	default:
		return false
	}
}

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
