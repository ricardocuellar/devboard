package memory

import (
	"context"
	"errors"
	"testing"

	"github.com/ricardocuellar/devboard/internal/domain"
)

func TestUserRepository_Create(t *testing.T) {
	repo := NewUserRepository()
	ctx := context.Background()

	user := &domain.User{
		Name:  "Ana",
		Email: "ana@devboard.dev",
	}

	err := repo.Create(ctx, user)

	if err != nil {
		t.Fatalf("Create() devolvió un error inesperado: %v", err)
	}

	if user.ID == "" {
		t.Fatalf("Create() no estableció un ID al usuario")
	}

	if user.ID != "user_1" {
		t.Fatalf("Create() estableció un ID %q; se esperaba %q", user.ID, "user_1")
	}
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	repo := NewUserRepository()
	ctx := context.Background()

	_, err := repo.GetByID(ctx, "user_999")

	if err == nil {
		t.Fatal("GetByID() debía devolver un error")
	}

	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("GetByID() devolvió %v; se esperaba ErrNotFound", err)
	}
}

// Test _GetByID
func TestUserRepository_GetByID(t *testing.T) {
	repo := NewUserRepository()
	ctx := context.Background()

	user := &domain.User{
		Name:  "Ana",
		Email: "ana@devboard.dev",
	}

	if err := repo.Create(ctx, user); err != nil {
		t.Fatalf("no se pudo preparar el usuario: %v", err)
	}

	found, err := repo.GetByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("GetByID() devolvió un error inesperado: %v", err)
	}

	if found.ID != user.ID {
		t.Errorf(
			"GetByID() devolvió ID %q; se esperaba %q",
			found.ID,
			user.ID,
		)
	}

	if found.Email != user.Email {
		t.Errorf(
			"GetByID() devolvió email %q; se esperaba %q",
			found.Email,
			user.Email,
		)
	}
}

// Test _GetByEmail
func TestUserRepository_GetByEmail(t *testing.T) {
	repo := NewUserRepository()
	ctx := context.Background()

	user := &domain.User{
		Name:  "Ana",
		Email: "ana@devboard.dev",
	}

	if err := repo.Create(ctx, user); err != nil {
		t.Fatalf("no se pudo preparar el usuario: %v", err)
	}

	found, err := repo.GetByEmail(ctx, "ana@devboard.dev")
	if err != nil {
		t.Fatalf(
			"GetByEmail() devolvió un error inesperado: %v",
			err,
		)
	}

	if found.ID != user.ID {
		t.Errorf(
			"GetByEmail() devolvió ID %q; se esperaba %q",
			found.ID,
			user.ID,
		)
	}
}
