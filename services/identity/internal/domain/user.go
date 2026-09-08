package domain

import "time"

// User representa al usuario del sistema
type User struct {
	ID           string
	Email        string
	PasswordHash string
	Name         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewUser construye un User válido
func NewUser(email, name, passwordHash string) *User {
	now := time.Now()

	return &User{
		Email:        email,
		Name:         name,
		PasswordHash: passwordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// IsValid verifica las invariantes de un usuario (email y name con string vacío)
func (u *User) IsValid() bool {
	return u.Email != "" && u.Name != ""
}
