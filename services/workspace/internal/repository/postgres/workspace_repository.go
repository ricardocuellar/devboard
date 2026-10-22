package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ricardocuellar/devboard/services/workspace/internal/domain"
	generated "github.com/ricardocuellar/devboard/services/workspace/internal/repository/postgres/generated"
)

// WorkspaceRepository implementa las interfaces de workspace para los use cases
type WorkspaceRepository struct {
	pool    *pgxpool.Pool
	queries *generated.Queries
}

// NewWorkspaceRepository constructor
func NewWorkspaceRepository(pool *pgxpool.Pool) *WorkspaceRepository {
	return &WorkspaceRepository{
		pool:    pool,
		queries: generated.New(pool),
	}
}

// CreateWorkspace guarda un workspace en PostgreSQL
func (r *WorkspaceRepository) CreateWorkspaceWithMember(ctx context.Context, ws *domain.Workspace, member *domain.WorkspaceMember) error {
	var createdBy, userID pgtype.UUID
	if err := createdBy.Scan(ws.CreatedBy); err != nil {
		return fmt.Errorf("created_by inválido: %w", domain.ErrInvalidInput)
	}

	if err := userID.Scan(member.UserID); err != nil {
		return fmt.Errorf("user_id inválido: %w", domain.ErrInvalidInput)
	}

	tx, err := r.pool.Begin(ctx)

	if err != nil {
		return fmt.Errorf("iniciando transacción: %w", err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	queries := r.queries.WithTx(tx)

	workspaceRow, err := queries.CreateWorkspace(ctx, generated.CreateWorkspaceParams{
		Name:        ws.Name,
		Slug:        ws.Slug,
		Description: ws.Description,
		CreatedBy:   createdBy,
	})

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "workspaces_slug_key" {
			return fmt.Errorf("workspace con slug %q: %w", ws.Slug, domain.ErrAlreadyExists)
		}

		return fmt.Errorf("insertando workspace en BD: %w", err)
	}

	memberRow, err := queries.AddWorkspaceMember(ctx, generated.AddWorkspaceMemberParams{
		WorkspaceID: workspaceRow.ID,
		UserID:      userID,
		Role:        string(member.Role),
	})

	if err != nil {
		return fmt.Errorf("insertando miembro inicial: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("confimando creación del workspace: %w", err)
	}

	ws.ID = workspaceRow.ID.String()
	ws.CreatedAt = workspaceRow.CreatedAt.Time
	ws.UpdatedAt = workspaceRow.UpdatedAt.Time

	member.ID = memberRow.ID.String()
	member.WorkspaceID = workspaceRow.ID.String()
	member.JoinedAt = memberRow.JoinedAt.Time

	return nil
}

// GetWorkspaceByID obtiene un workspace por ID
func (r *WorkspaceRepository) GetWorkspaceByID(ctx context.Context, id string) (*domain.Workspace, error) {
	var pgID pgtype.UUID
	if err := pgID.Scan(id); err != nil {
		return nil, fmt.Errorf("id inválido %s: %w", id, domain.ErrNotFound)
	}

	result, err := r.queries.GetWorkspaceByID(ctx, pgID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("workspace %s: %w", id, domain.ErrNotFound)
		}

		return nil, fmt.Errorf("consultando workspace %s: %w", id, err)
	}

	return toDomainWorkspace(result), nil

}

// GetWorkspaceBySlug obtiene un workspace por slug
func (r *WorkspaceRepository) GetWorkspaceBySlug(ctx context.Context, slug string) (*domain.Workspace, error) {
	result, err := r.queries.GetWorkspaceBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("workspace slug %s: %w", slug, domain.ErrNotFound)
		}

		return nil, fmt.Errorf("consultando workspace por slug: %w", err)
	}

	return toDomainWorkspace(result), nil
}

// ListWorkspacesByUser lista los workspaces donde el usuario es miembro
func (r *WorkspaceRepository) ListWorkspacesByUser(ctx context.Context, userID string) ([]*domain.Workspace, error) {
	var pgUserID pgtype.UUID
	if err := pgUserID.Scan(userID); err != nil {
		return nil, fmt.Errorf("user_id inválido: %w", domain.ErrInvalidInput)
	}

	results, err := r.queries.ListWorkspacesByUser(ctx, pgUserID)

	if err != nil {
		return nil, fmt.Errorf("listando workspaces del usuario: %w", err)
	}

	workspaces := make([]*domain.Workspace, len(results))
	for i, result := range results {
		workspaces[i] = toDomainWorkspace(result)
	}

	return workspaces, nil
}

// AddMember agrega un usuario a un workspace con un rol
func (r *WorkspaceRepository) AddMember(ctx context.Context, member *domain.WorkspaceMember) error {
	var wsID, userID pgtype.UUID

	if err := wsID.Scan(member.WorkspaceID); err != nil {
		return fmt.Errorf("workspace_id inválido: %w", domain.ErrInvalidInput)
	}
	if err := userID.Scan(member.UserID); err != nil {
		return fmt.Errorf("user_id inválido: %w", domain.ErrInvalidInput)
	}

	result, err := r.queries.AddWorkspaceMember(ctx, generated.AddWorkspaceMemberParams{
		WorkspaceID: wsID,
		UserID:      userID,
		Role:        string(member.Role),
	})

	if err != nil {
		return fmt.Errorf("agregando miembro al workspace: %w", err)
	}

	member.ID = result.ID.String()
	member.JoinedAt = result.JoinedAt.Time

	return nil
}

// GetMember obtiene la membresía de un usuario en un workspace
func (r *WorkspaceRepository) GetMember(ctx context.Context, workspaceID, userID string) (*domain.WorkspaceMember, error) {
	var pgWsID, pgUserID pgtype.UUID

	if err := pgWsID.Scan(workspaceID); err != nil {
		return nil, fmt.Errorf("workspace_id inválido: %w", domain.ErrInvalidInput)
	}
	if err := pgUserID.Scan(userID); err != nil {
		return nil, fmt.Errorf("user_id inválido: %w", domain.ErrInvalidInput)
	}

	result, err := r.queries.GetWorkspaceMember(ctx, generated.GetWorkspaceMemberParams{
		WorkspaceID: pgWsID,
		UserID:      pgUserID,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("membresía no encontrada: %w", domain.ErrNotFound)
		}

		return nil, fmt.Errorf("consultando membresía: %w", err)
	}

	return toDomainMember(result), nil

}

// ListMembers lista todos los miembros de un workspace
func (r *WorkspaceRepository) ListMembers(ctx context.Context, workspaceID string) ([]*domain.WorkspaceMember, error) {
	var pgWsID pgtype.UUID
	if err := pgWsID.Scan(workspaceID); err != nil {
		return nil, fmt.Errorf("workspace_id inválido: %w", domain.ErrInvalidInput)
	}

	results, err := r.queries.ListWorkspaceMembers(ctx, pgWsID)

	if err != nil {
		return nil, fmt.Errorf("listando miembros del workspace: %w", err)
	}

	members := make([]*domain.WorkspaceMember, len(results))
	for i, result := range results {
		members[i] = toDomainMember(result)
	}

	return members, nil
}

// RemoveMember elimina un miembro del workspace
func (r *WorkspaceRepository) RemoveMember(ctx context.Context, workspaceID, userID string) error {
	var pgWsID, pgUserID pgtype.UUID

	if err := pgWsID.Scan(workspaceID); err != nil {
		return fmt.Errorf("workspace_id inválido: %w", domain.ErrInvalidInput)
	}
	if err := pgUserID.Scan(userID); err != nil {
		return fmt.Errorf("user_id inválido: %w", domain.ErrInvalidInput)
	}

	return r.queries.RemoveWorkspaceMember(ctx, generated.RemoveWorkspaceMemberParams{
		WorkspaceID: pgWsID,
		UserID:      pgUserID,
	})
}

// toDomainWorkspace convierte el modelo generado por sqlc al modelo de dominio
func toDomainWorkspace(w generated.Workspace) *domain.Workspace {
	return &domain.Workspace{
		ID:          w.ID.String(),
		Name:        w.Name,
		Slug:        w.Slug,
		Description: w.Description,
		CreatedBy:   w.CreatedBy.String(),
		CreatedAt:   w.CreatedAt.Time,
		UpdatedAt:   w.UpdatedAt.Time,
	}
}

// toDomainMember convierte el modelo de membresía de sqlc al dominio
func toDomainMember(m generated.WorkspaceMember) *domain.WorkspaceMember {
	return &domain.WorkspaceMember{
		ID:          m.ID.String(),
		WorkspaceID: m.WorkspaceID.String(),
		UserID:      m.UserID.String(),
		Role:        domain.WorkspaceRole(m.Role),
		JoinedAt:    m.JoinedAt.Time,
	}
}
