package usecase

import (
	"context"
	"fmt"

	"github.com/ricardocuellar/devboard/internal/domain"
)

// WorkspaceRepository define lo que WorkspaceUseCase necesita del repositorio
type WorkspaceRepository interface {
	CreateWorkspace(ctx context.Context, ws *domain.Workspace) error
	GetWorkspaceByID(ctx context.Context, id string) (*domain.Workspace, error)
	GetWorkspaceBySlug(ctx context.Context, slug string) (*domain.Workspace, error)
	ListWorkspacesByUser(ctx context.Context, userID string) ([]*domain.Workspace, error)
	AddMember(ctx context.Context, member *domain.WorkspaceMember) error
	GetMember(ctx context.Context, workspaceID, userID string) (*domain.WorkspaceMember, error)
	ListMembers(ctx context.Context, workspaceID string) ([]*domain.WorkspaceMember, error)
	RemoveMember(ctx context.Context, workspaceID, userID string) error
}

// WorkspaceUseCase orquesta las operaciones de workspace
type WorkspaceUseCase struct {
	repo WorkspaceRepository
}

// NewWorkspaceUseCase constructor
func NewWorkspaceUseCase(repo WorkspaceRepository) *WorkspaceUseCase {
	return &WorkspaceUseCase{repo: repo}
}

// CreateWorkspace crea un workspace y agrega al creador como admin
func (uc *WorkspaceUseCase) CreateWorkspace(ctx context.Context, name, description, createdBy string) (*domain.Workspace, error) {
	slug := domain.GenerateSlug(name)
	existing, err := uc.repo.GetWorkspaceBySlug(ctx, slug)

	if err == nil && existing != nil {
		return nil, fmt.Errorf("workspace con slug %q: %w", slug, domain.ErrAlreadyExists)
	}

	ws := domain.NewWorkspace(name, description, createdBy)

	if err := uc.repo.CreateWorkspace(ctx, ws); err != nil {
		return nil, fmt.Errorf("creando workspace: %w", err)
	}

	member := domain.NewWorkspaceMember(ws.ID, createdBy, domain.WorkspaceRoleAdmin)

	if err := uc.repo.AddMember(ctx, member); err != nil {
		return nil, fmt.Errorf("agregando creador como admin: %w", err)
	}

	return ws, nil
}

// GetWorkspace obtiene un workspace verificando que el usuario sea miembro
func (uc *WorkspaceUseCase) GetWorkspace(ctx context.Context, workspaceID, requesterID string) (*domain.Workspace, error) {
	if err := uc.requireMember(ctx, workspaceID, requesterID); err != nil {
		return nil, err
	}

	ws, err := uc.repo.GetWorkspaceByID(ctx, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("obteniendo workspace: %w", err)
	}

	return ws, nil
}

// ListUserWorkspaces lista los workspaces del usuario autenticado
func (uc *WorkspaceUseCase) ListUserWorkspaces(ctx context.Context, userID string) ([]*domain.Workspace, error) {
	workspaces, err := uc.repo.ListWorkspacesByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("listando workspaces del usuario: %w", err)
	}

	return workspaces, nil
}

// AddMember agrega un miembro al workspace.
func (uc *WorkspaceUseCase) AddMember(ctx context.Context, workspaceID, requesterID, newUserID string, role domain.WorkspaceRole) (*domain.WorkspaceMember, error) {
	// TODO
}

// ListMembers lista los miembros de un workspace.
func (uc *WorkspaceUseCase) ListMembers(ctx context.Context, workspaceID, requesterID string) ([]*domain.WorkspaceMember, error) {
	// TODO
}

// RemoveMember elimina un miembro del workspace.
func (uc *WorkspaceUseCase) RemoveMember(ctx context.Context, workspaceID, requesterID, targetUserID string) error {
	// TODO
}

// requireMember verifica que el usuario sea miembro del workspace.
func (uc *WorkspaceUseCase) requireMember(ctx context.Context, workspaceID, userID string) error {
	_, err := uc.repo.GetMember(ctx, workspaceID, userID)
	if err != nil {
		return fmt.Errorf("acceso al workspace denegado: %w", domain.ErrForbidden)
	}

	return nil
}

// requireAdmin verifica que el usuario sea admin del workspace.
func (uc *WorkspaceUseCase) requireAdmin(ctx context.Context, workspaceID, userID string) error {
	// TODO
}
