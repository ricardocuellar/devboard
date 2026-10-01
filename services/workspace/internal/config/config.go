package config

import (
	"fmt"
	"os"
)

type Config struct {
	HTTPAddr    string
	GRPCAddr    string
	DatabaseURL string
	JWTSecret   string
}

func Load() (Config, error) {

	cfg := Config{
		HTTPAddr:    envOrDefault("WORKSPACE_HTTP_ADDR", ":8082"),
		GRPCAddr:    envOrDefault("WORKSPACE_GRPC_ADDR", ":9092"),
		DatabaseURL: envOrDefault("WORKSPACE_DATABASE_URL", ""),
		JWTSecret:   envOrDefault("JWT_SECRET", "dev-secret-key-change-in-production-min-32-chars$$"),
	}

	// httpAddr := os.Getenv("WORKSPACE_HTTP_ADDR")
	// if httpAddr == "" {
	// 	httpAddr = ":8081"
	// }

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("WORKSPACE_DATABASE_URL es requerida")
	}

	if cfg.JWTSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET es requerida")
	}

	if len(cfg.JWTSecret) < 32 {
		return Config{}, fmt.Errorf("JWT_SECRET debe tener al menos 32 bytes")
	}

	if cfg.HTTPAddr[0] != ':' {
		return Config{}, fmt.Errorf("WORKSPACE_HTTP_ADDR debe usar el formato de puertos :port")
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
