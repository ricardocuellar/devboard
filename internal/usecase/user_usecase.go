package usecase

import (
	"context"

	"github.com/ricardocuellar/devboard/internal/domain"
)

// UserRepository define lo que el Use Case necesita
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

// UserUseCase orquesta las operaciones de negocio sobre usuario
type UserUseCase struct {
	repo UserRepository
}

// NewUserUseCase constructor
func NewUserUseCase(repo UserRepository) *UserUseCase {
	return &UserUseCase{repo: repo}
}
