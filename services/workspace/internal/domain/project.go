// Package domain project
package domain

import "time"

// Project representa un proyecto dentro de un workspace
type Project struct {
	ID          string
	WorkspaceID string
	Name        string
	Description string
	CreatedBy   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewProject construye un Project válido
func NewProject(workspaceID, name, description, createdBy string) *Project {
	now := time.Now()
	return &Project{
		WorkspaceID: workspaceID,
		Name:        name,
		Description: description,
		CreatedBy:   createdBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
