// Package auth argon2id
package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

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

// VerifyPassword compara un password en texto plano con su hash almacenado
func VerifyPassword(password, encodedHash string) error {
	params, salt, hash, err := parseHash(encodedHash)
	if err != nil {
		return fmt.Errorf("parseando hash: %w", err)
	}

	otherHash := argon2.IDKey(
		[]byte(password),
		salt,
		params.iterations,
		params.memory,
		params.parallelims,
		params.keyLength,
	)

	if subtle.ConstantTimeCompare(hash, otherHash) != 1 {
		return ErrPasswordMissmatch
	}

	return nil
}

// parseHash extrae los parámetros y valores del hash en formato Argon2.
func parseHash(encodedHash string) (*argon2Params, []byte, []byte, error) {
	//$argon2id$v=19$m=65536,t=3,p=4$c2FsdA$hashGenerado
	parts := strings.Split(encodedHash, "$")
	// Formato esperado: ["", "argon2id", "v=19", "m=65536,t=3,p=2", "salt", "hash"]
	if len(parts) != 6 {
		return nil, nil, nil, ErrInvalidPasswordHash
	}

	if parts[1] != "argon2id" {
		return nil, nil, nil, ErrInvalidPasswordHash
	}

	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return nil, nil, nil, ErrInvalidPasswordHash
	}

	if version != argon2.Version {
		return nil, nil, nil, ErrInvalidPasswordHash
	}

	params := &argon2Params{}
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d",
		&params.memory, &params.iterations, &params.parallelims); err != nil {
		return nil, nil, nil, ErrInvalidPasswordHash
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, ErrInvalidPasswordHash
	}
	params.saltLength = uint32(len(salt))

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])

	if err != nil {
		return nil, nil, nil, ErrInvalidPasswordHash
	}

	params.keyLength = uint32(len(hash))

	return params, salt, hash, nil
}
