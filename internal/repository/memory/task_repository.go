package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/ricardocuellar/devboard/internal/domain"
)

// TaskRepository struct
type TaskRepository struct {
	mu    sync.RWMutex
	tasks map[string]*domain.Task
}

// NewTaskRepository constructor
func NewTaskRepository() *TaskRepository {
	return &TaskRepository{
		tasks: make(map[string]*domain.Task),
	}
}

// Create método para crear tareas
func (r *TaskRepository) Create(_ context.Context, task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	task.ID = fmt.Sprintf("task_%d", len(r.tasks)+1)
	r.tasks[task.ID] = task
	return nil
}

// Update método para actualizar tarea
func (r *TaskRepository) Update(_ context.Context, task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tasks[task.ID] = task
	return nil
}

// GetByID método para obtener tarea por ID
func (r *TaskRepository) GetByID(_ context.Context, id string) (*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, ok := r.tasks[id]
	if !ok {
		return nil, fmt.Errorf("tarea %s: %w", id, domain.ErrNotFound)
	}

	return task, nil
}

// ListByProject método para listar tareas con base a un ProjectID
func (r *TaskRepository) ListByProject(_ context.Context, projectID string) ([]*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.Task

	for _, t := range r.tasks {
		if t.ProjectID == projectID {
			result = append(result, t)
		}
	}

	return result, nil
}
