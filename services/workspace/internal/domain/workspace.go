// Package domain workspace
package domain

import (
	"strings"
	"time"
	"unicode"
)

// WorkspaceRole define los roles posibles dentro de un workspace.
type WorkspaceRole string

const (
	// WorkspaceRoleAdmin puede gestionar miembros y proyectos
	WorkspaceRoleAdmin WorkspaceRole = "admin"
	// WorkspaceRoleMember puede ver y crear tareas
	WorkspaceRoleMember WorkspaceRole = "member"
	// WorkspaceRoleViewer solo puede ver
	WorkspaceRoleViewer WorkspaceRole = "viewer"
)

// IsValid verifica que el rol sea uno de los permitidos
func (r WorkspaceRole) IsValid() bool {
	switch r {
	case WorkspaceRoleAdmin, WorkspaceRoleMember, WorkspaceRoleViewer:
		return true
	default:
		return false
	}
}

// Workspace representa un espacio de trabajo que agrupa usuarios y proyectos
type Workspace struct {
	ID          string
	Name        string
	Slug        string
	Description string
	CreatedBy   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// WorkspaceMember representa la membresía de un usuario en un workspace
type WorkspaceMember struct {
	ID          string
	WorkspaceID string
	UserID      string
	Role        WorkspaceRole
	JoinedAt    time.Time
}

// NewWorkspace construye un Workspace con slug generado automáticamente
func NewWorkspace(name, description, createdBy string) *Workspace {
	now := time.Now()
	return &Workspace{
		Name:        name,
		Slug:        GenerateSlug(name),
		Description: description,
		CreatedBy:   createdBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// GenerateSlug convierte un nombre a un slug URL-friendly. "Mi Equipo Dev" -> "mi-equipo-dev"
func GenerateSlug(name string) string {
	slug := strings.ToLower(name)
	var result strings.Builder

	for _, r := range slug {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result.WriteRune(r)
		} else {
			result.WriteRune('-')
		}
	}

	cleaned := strings.Join(strings.FieldsFunc(result.String(), func(r rune) bool {
		return r == '-'
	}), "-")

	return cleaned

}

// NewWorkspaceMember construye una membresía con rol por defecto
func NewWorkspaceMember(workspaceID, userID string, role WorkspaceRole) *WorkspaceMember {
	return &WorkspaceMember{
		WorkspaceID: workspaceID,
		UserID:      userID,
		Role:        role,
		JoinedAt:    time.Now(),
	}
}
