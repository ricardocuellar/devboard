package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ricardocuellar/devboard/services/workspace/internal/domain"
	generated "github.com/ricardocuellar/devboard/services/workspace/internal/repository/postgres/generated"
)

// ProjectRepository implementa las interfaces de project para los use cases
type ProjectRepository struct {
	pool    *pgxpool.Pool
	queries *generated.Queries
}

// NewProjectRepository constructor
func NewProjectRepository(pool *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{
		pool:    pool,
		queries: generated.New(pool),
	}
}

// Create guarda un proyecto en PostgreSQL
func (r *ProjectRepository) Create(ctx context.Context, project *domain.Project) error {
	var wsID, createdBy pgtype.UUID

	if err := wsID.Scan(project.WorkspaceID); err != nil {
		return fmt.Errorf("workspace_id inválido: %w", domain.ErrInvalidInput)
	}
	if err := createdBy.Scan(project.CreatedBy); err != nil {
		return fmt.Errorf("created_by inválido: %w", domain.ErrInvalidInput)
	}

	result, err := r.queries.CreateProject(ctx, generated.CreateProjectParams{
		WorkspaceID: wsID,
		Name:        project.Name,
		Description: project.Description,
		CreatedBy:   createdBy,
	})

	if err != nil {
		return fmt.Errorf("insertando proyecto en BD: %w", err)
	}

	project.ID = result.ID.String()
	project.CreatedAt = result.CreatedAt.Time
	project.UpdatedAt = result.UpdatedAt.Time

	return nil
}

// GetByID obtiene un proyecto por ID
func (r *ProjectRepository) GetByID(ctx context.Context, id string) (*domain.Project, error) {
	var pgID pgtype.UUID
	if err := pgID.Scan(id); err != nil {
		return nil, fmt.Errorf("id inválido %s: %w", id, domain.ErrNotFound)
	}

	result, err := r.queries.GetProjectByID(ctx, pgID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("proyecto %s : %w", id, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("consultando proyecto %s: %w", id, err)
	}

	return toDomainProject(result), nil
}

// ListByWorkspace lista los proyectos de un workspace
func (r *ProjectRepository) ListByWorkspace(ctx context.Context, workspaceID string) ([]*domain.Project, error) {
	var pgWsID pgtype.UUID
	if err := pgWsID.Scan(workspaceID); err != nil {
		return nil, fmt.Errorf("workspace_id inválido: %w", domain.ErrInvalidInput)
	}

	results, err := r.queries.ListProjectsByWorkspace(ctx, pgWsID)
	if err != nil {
		return nil, fmt.Errorf("listando proyectos del workspace: %w", err)
	}

	projects := make([]*domain.Project, len(results))

	for i, result := range results {
		projects[i] = toDomainProject(result)
	}

	return projects, nil

}

// toDomainProject convierte el modelo sqlc al dominio
func toDomainProject(p generated.Project) *domain.Project {
	return &domain.Project{
		ID:          p.ID.String(),
		WorkspaceID: p.WorkspaceID.String(),
		Name:        p.Name,
		Description: p.Description,
		CreatedBy:   p.CreatedBy.String(),
		CreatedAt:   p.CreatedAt.Time,
		UpdatedAt:   p.UpdatedAt.Time,
	}
}
