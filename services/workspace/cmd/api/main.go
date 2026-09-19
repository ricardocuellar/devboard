package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/ricardocuellar/devboard/services/workspace/internal/auth"
	"github.com/ricardocuellar/devboard/services/workspace/internal/config"
	"github.com/ricardocuellar/devboard/services/workspace/internal/handler"
	"github.com/ricardocuellar/devboard/services/workspace/internal/middlewares"
	"github.com/ricardocuellar/devboard/services/workspace/internal/repository/postgres"
	"github.com/ricardocuellar/devboard/services/workspace/internal/server"
	"github.com/ricardocuellar/devboard/services/workspace/internal/usecase"
	"github.com/ricardocuellar/devboard/services/workspace/internal/validator"
)

func main() {
	if err := run(); err != nil {
		slog.Error("Workspace service se ha detenido", slog.Any("error", err))
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	ctx := context.Background()

	pool, err := postgres.NewPool(
		ctx,
		postgres.Config{
			URL:             cfg.DatabaseURL,
			MaxConns:        10,
			MinConns:        2,
			MaxConnLifetime: time.Hour,
			MaxConnIdleTime: 30 * time.Minute,
		},
	)

	if err != nil {
		return err
	}

	defer pool.Close()

	logger.Info("La conexión a la base de datos de workspace fue establecida")

	tokenVerifier, err := auth.NewTokenVerifier([]byte(cfg.JWTSecret))

	if err != nil {
		return err
	}

	// Repositories
	workspaceRepository := postgres.NewWorkspaceRepository(pool)
	projectRepository := postgres.NewProjectRepository(pool)

	// Use cases
	workspaceUseCase := usecase.NewWorkspaceUseCase(workspaceRepository)

	projectUseCase := usecase.NewProjectUseCase(projectRepository, workspaceRepository)

	validate := validator.New()

	// handlers
	workspaceHandler := handler.NewWorkspaceHandler(
		workspaceUseCase,
		validate,
		logger,
	)

	projectHandler := handler.NewProjectHandler(
		projectUseCase,
		validate,
		logger,
	)

	// Server
	srv := server.New(cfg.HTTPAddr,
		server.WithLogger(logger),
		server.WithReadTimeout(15*time.Second),
		server.WithWriteTimeout(30*time.Second),
	)

	authMiddleware := middlewares.RequireAuth(tokenVerifier, logger)

	// Rutas públicas
	srv.RegisterRoutes("GET /health", http.HandlerFunc(
		func(writer http.ResponseWriter, _ *http.Request) {
			writer.Header().Set(
				"Content-Type", "application/json",
			)

			writer.WriteHeader(
				http.StatusOK,
			)

			_, _ = writer.Write([]byte(
				`{"status": "ok", "service":"identity"}`,
			))
		},
	))

	srv.RegisterRoutes("POST /api/v1/auth/register", http.HandlerFunc(authHandler.Register))
	srv.RegisterRoutes("POST /api/v1/auth/login", http.HandlerFunc(authHandler.Login))
	srv.RegisterRoutes("POST /api/v1/auth/refresh", http.HandlerFunc(authHandler.Refresh))
	srv.RegisterRoutes("POST /api/v1/auth/logout", http.HandlerFunc(authHandler.Logout))
	srv.RegisterRoutes("POST /api/v1/users", http.HandlerFunc(userHandler.Create))
	srv.RegisterRoutes("GET /api/v1/users/{id}", http.HandlerFunc(userHandler.Get))

	// HTTP middlewares
	srv.Use(
		middlewares.SecurityHeaders(
			middlewares.SecurityConfig{
				AllowedOrigins: []string{
					"http://localhost:3000",
					"http://localhost:5173",
					"https://ricardocuellar.dev",
				},
			},
		),
	)

	srv.Use(
		middlewares.Logger(logger),
	)

	srv.Use(
		middlewares.Recovery(logger),
	)

	logger.Info("Se inicio IDENTITY SERVICE",
		slog.String("addr", cfg.HTTPAddr),
	)

	return srv.Start()
}
