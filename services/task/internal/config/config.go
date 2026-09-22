package config

import (
	"fmt"
	"os"
)

type Config struct {
	HTTPAddr    string
	DatabaseURL string
	JWTSecret   string
}

func Load() (Config, error) {

	cfg := Config{
		HTTPAddr:    envOrDefault("TASK_HTTP_ADDR", ":8083"),
		DatabaseURL: envOrDefault("TASK_DATABASE_URL", ""),
		JWTSecret:   envOrDefault("JWT_SECRET", "dev-secret-key-change-in-production-min-32-chars$$"),
	}

	// httpAddr := os.Getenv("TASK_HTTP_ADDR")
	// if httpAddr == "" {
	// 	httpAddr = ":8081"
	// }

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("TASK_DATABASE_URL es requerida")
	}

	if cfg.JWTSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET es requerida")
	}

	if len(cfg.JWTSecret) < 32 {
		return Config{}, fmt.Errorf("JWT_SECRET debe tener al menos 32 bytes")
	}

	if cfg.HTTPAddr[0] != ':' {
		return Config{}, fmt.Errorf("TASK_HTTP_ADDR debe usar el formato de puertos :port")
	}

	return cfg, nil

}

func envOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}
