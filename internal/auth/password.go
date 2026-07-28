package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"golang.org/x/crypto/argon2"
)

// ErrInvalidPasswordHash indica que el has almacenado tiene formato incorrecto
var ErrInvalidPasswordHash = errors.New("formato de hash inválido")

// ErrPasswordMissmatch indica que el passoword no coincide con el hash
var ErrPasswordMissmatch = errors.New("password incorrecto")

// argon2Params define los parámetros del algoritmo
type argon2Params struct {
	memory      uint32 // 64 MB recomendados como mínimo
	iterations  uint32
	parallelims uint8
	saltLength  uint32
	keyLength   uint32
}

// defaultArgon2Params son los parámetros por defecto
var defaultArgon2Params = argon2Params{
	memory:      64 * 1024,
	iterations:  3,
	parallelims: 2,
	saltLength:  16, //bytes
	keyLength:   32, //bytes
}

// HashPassword genera un hash seguro del password usando Argon2id
func HashPassword(password string) (string, error) {
	salt := make([]byte, defaultArgon2Params.saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generando salt: %w", err)
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		defaultArgon2Params.iterations,
		defaultArgon2Params.memory,
		defaultArgon2Params.parallelims,
		defaultArgon2Params.keyLength,
	)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		defaultArgon2Params.memory,
		defaultArgon2Params.iterations,
		defaultArgon2Params.parallelims,
		b64Salt,
		b64Hash,
	)

	return encoded, nil
}
