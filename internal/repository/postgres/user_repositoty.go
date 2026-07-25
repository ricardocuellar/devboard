package postgres

import (
	"context"
	"fmt"

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
