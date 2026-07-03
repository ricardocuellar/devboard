package domain

import (
	"encoding/json"
	"testing"
)

type CreateTaskRequest struct {
	Title  string     `json:"title"`
	Status TaskStatus `json:"status"`
}

func TestCreateTaskRequestWithValidStatus(t *testing.T) {
	body := []byte(`{"title": "Mi tarea", "status": "todo"}`)

	var req CreateTaskRequest

	err := json.Unmarshal(body, &req)
	if err != nil {
		t.Fatalf("no esperaba error, pero recibí: %v", err)
	}

	if req.Status != TaskStatusTodo {
		t.Fatalf("esperaba status todo, recibí: %s", req.Status)
	}
}

func TestCreateTaskRequestWithInvalidStatus(t *testing.T) {
	body := []byte(`{"title": "Mi tarea", "status": "banana"}`)

	var req CreateTaskRequest

	err := json.Unmarshal(body, &req)
	if err == nil {
		t.Fatalf("esperaba error, pero recibí nil")
	}

	expected := `estado de tarea inválido: "banana"`

	if err.Error() != expected {
		t.Fatalf("esperaba error %q, recibí %q", expected, err.Error())
	}

}
