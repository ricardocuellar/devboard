// Package main
package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	_ "github.com/ricardocuellar/devboard/docs"
	"github.com/ricardocuellar/devboard/internal/handler"
	"github.com/ricardocuellar/devboard/internal/logger"
	"github.com/ricardocuellar/devboard/internal/middlewares"
	pgrepo "github.com/ricardocuellar/devboard/internal/repository/postgres"
	"github.com/ricardocuellar/devboard/internal/server"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Devboard API
// @version 1.0
// @description API REST para gestión de proyectos
// @contact.name Soporte Devboard Ricardo
// @host  localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Escribe el token con el formato: Bearer {token}
func main() {

	// 1. Configuración del entorno
	// 2. Logger
	log := setupLogger()

	dbURL := os.Getenv("DATABASE_URL")

	// Crear pool de conexiones
	pool, err := pgrepo.NewPool(context.Background(), pgrepo.Config{
		URL:             dbURL,
		MaxConns:        25,
		MinConns:        5,
		MaxConnLifetime: 1 * time.Hour,
		MaxConnIdleTime: 30 * time.Minute,
	})

	if err != nil {
		log.Error("no se pudo conectar a la base de datos", slog.Any("error", err))
		os.Exit(1)
	}

	defer pool.Close()

	log.Info("conexión a base de datos establecida")

	// Auth Infrastructure
	// jwtSecret := os.Getenv("JWT_SECRET")
	// refreshTokenDuration := 7 * 24 * time.Hour

	// jwtManager, err := auth.NewJWTManager(auth.JWTConfig{
	// 	SecretKey:            []byte(jwtSecret),
	// 	AccessTokenDuration:  15 * time.Minute,
	// 	RefreshTokenDuration: refreshTokenDuration,
	// })

	// if err != nil {
	// 	log.Error("error configurando JWT", slog.Any("error", err))
	// 	os.Exit(1)
	// }

	// 3. Dependencies compartidas

	// Infrastructure: adapters
	// userRepo := memory.NewUserRepository()
	// taskRepo := memory.NewTaskRepository()

	// 4. Handlers
	healthHandler := handler.NewHealthHandler()

	// notifier := notification.NewLogNotifier(log)
	// _ = notifier

	// 5. Servidor con Functional options
	srv := server.New(":8080",
		server.WithLogger(log),
		server.WithReadTimeout(15*time.Second),
		server.WithWriteTimeout(30*time.Second),
	)

	// 6. Registro de rutas públicas
	srv.RegisterRoutes("GET /docs/", httpSwagger.WrapHandler)
	srv.RegisterRoutes("GET /health", healthHandler)

	// 7. Middleware chain
	srv.UseChain(
		middlewares.Recovery(log),
		middlewares.Logger(log),
		middlewares.SecurityHeaders(middlewares.SecurityConfig{
			AllowedOrigins: []string{
				"http://localhost:3000",
				"http://localhost:5173",
				"https://ricardocuellar.dev",
			},
		}),
	)

	// 8. Arrancar shutdown limpio

	if err := srv.Start(); err != nil {
		log.Error("error fatal", "error", err)
		os.Exit(1)
	}

}

func setupLogger() *slog.Logger {

	if os.Getenv("GO_ENV") == "production" {
		return logger.New(logger.ProductionConfig())
	}

	return logger.New(logger.DefaultConfig())

}
