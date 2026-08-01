package usecase

import (
	"context"
	"fmt"

	"github.com/ricardocuellar/devboard/internal/domain"
)

// ProjectRepository define lo que ProjectUseCase necesita
type ProjectRepository interface {
	Create(ctx context.Context, project *domain.Project) error
	GetByID(ctx context.Context, id string) (*domain.Project, error)
	ListByWorkspace(ctx context.Context, workspaceID string) ([]*domain.Project, error)
}

// ProjectUseCase orquesta las operaciones de proyecto
type ProjectUseCase struct {
	projectRepo   ProjectRepository
	workspaceRepo WorkspaceRepository
}

// NewProjectUseCase constructor
func NewProjectUseCase(projectRepo ProjectRepository, workspaceRepo WorkspaceRepository) *ProjectUseCase {
	return &ProjectUseCase{
		projectRepo:   projectRepo,
		workspaceRepo: workspaceRepo,
	}
}

// CreateProject crea un proyecto dentro de un workspace.
func (uc *ProjectUseCase) CreateProject(ctx context.Context, workspaceID, name, description, createdBy string) (*domain.Project, error) {
	_, err := uc.workspaceRepo.GetMember(ctx, workspaceID, createdBy)
	if err != nil {
		return nil, fmt.Errorf("acceso al workspace denegado: %w", domain.ErrForbidden)
	}

	project := domain.NewProject(workspaceID, name, description, createdBy)

	if err := uc.projectRepo.Create(ctx, project); err != nil {
		return nil, fmt.Errorf("creando proyecto: %w", err)
	}

	return project, nil
}

// GetProject obtiene un proyecto verificando membresía en el workspace
func (uc *ProjectUseCase) GetProject(ctx context.Context, projectID, requesterID string) (*domain.Project, error) {
	// TODO
}

// ListWorkspaceProjects lista los proyectos de un workspace
func (uc *ProjectUseCase) ListWorkspaceProjects(ctx context.Context, workspaceID, requesterID string) ([]*domain.Project, error) {
	// TODO
}
