package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ricardocuellar/devboard/internal/domain"
	generated "github.com/ricardocuellar/devboard/internal/repository/postgres/generated"
)

// UserRepository implementa usecase.UserRepository usando postgresql
type UserRepository struct {
	pool    *pgxpool.Pool
	queries *generated.Queries
}

// NewUserRepository constructor
func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool:    pool,
		queries: generated.New(pool),
	}
}

// Create guarda un nuevo usuario en la BD de PostgreSQL
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	result, err := r.queries.CreateUser(ctx, generated.CreateUserParams{
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Name:         user.Name,
	})

	if err != nil {
		return fmt.Errorf("insertando usuario en BD: %w", err)
	}

	user.ID = result.ID.String()
	user.CreatedAt = result.CreatedAt.Time //pgtype.Timestamptz
	user.UpdatedAt = result.UpdatedAt.Time

	return nil
}

// GetByID busca usuario por su ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var pgID pgtype.UUID

	if err := pgID.Scan(id); err != nil {
		return nil, fmt.Errorf("id inválido %s: %w", id, domain.ErrNotFound)
	}

	result, err := r.queries.GetUserByID(ctx, pgID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("usuario %s: %w", id, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("consultando usuario %s: %w", id, err)
	}

	return toDomainUser(result), nil
}

// GetByEmail busca un usuario por su email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	result, err := r.queries.GetUserByEmail(ctx, email)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("email %s: %w", email, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("consultando usuario por email: %w", err)
	}

	return toDomainUser(result), nil
}

// toDaminUser convierte el modelo generado por sqlc al modelo de dominio
func toDomainUser(u generated.User) *domain.User {
	return &domain.User{
		ID:           u.ID.String(),
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Name:         u.Name,
		CreatedAt:    u.CreatedAt.Time,
		UpdatedAt:    u.UpdatedAt.Time,
	}
}
