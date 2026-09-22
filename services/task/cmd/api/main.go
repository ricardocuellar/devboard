package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/ricardocuellar/devboard/services/task/internal/auth"
	"github.com/ricardocuellar/devboard/services/task/internal/config"
	"github.com/ricardocuellar/devboard/services/task/internal/handler"
	"github.com/ricardocuellar/devboard/services/task/internal/middlewares"
	"github.com/ricardocuellar/devboard/services/task/internal/repository/postgres"
	"github.com/ricardocuellar/devboard/services/task/internal/server"
	"github.com/ricardocuellar/devboard/services/task/internal/usecase"
	"github.com/ricardocuellar/devboard/services/task/internal/validator"
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

	logger.Info("La conexión a la base de datos de task fue establecida")

	tokenVerifier, err := auth.NewTokenVerifier([]byte(cfg.JWTSecret))

	if err != nil {
		return err
	}

	// Repositories
	taskRepository := postgres.NewTaskRepository(pool)
	commentRepository := postgres.NewCommentRepository(pool)

	// Use cases
	taskUseCase := usecase.NewTaskUseCase(taskRepository, taskRepository)

	commentUseCase := usecase.NewCommentUseCase(commentRepository, taskRepository)

	validate := validator.New()

	// handlers
	taskHandler := handler.NewTaskHandler(
		taskUseCase,
		validate,
		logger,
	)

	commentHandler := handler.NewCommentHandler(
		commentUseCase,
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
				`{"status": "ok", "service":"task"}`,
			))
		},
	))

	//Task routes
	srv.RegisterRoutes("POST /api/v1/tasks", middlewares.Chain(http.HandlerFunc(taskHandler.Create), authMiddleware))
	srv.RegisterRoutes("GET /api/v1/tasks/{id}", middlewares.Chain(http.HandlerFunc(taskHandler.Get), authMiddleware))
	srv.RegisterRoutes("PUT /api/v1/tasks/{id}/status", middlewares.Chain(http.HandlerFunc(taskHandler.UpdateStatus), authMiddleware))
	srv.RegisterRoutes("PUT /api/v1/tasks/{id}/assign", middlewares.Chain(http.HandlerFunc(taskHandler.Assign), authMiddleware))
	srv.RegisterRoutes("GET /api/v1/projects/{id}/tasks", middlewares.Chain(http.HandlerFunc(taskHandler.ListByProject), authMiddleware))

	// Comments Routes
	srv.RegisterRoutes("POST /api/v1/tasks/{id}/comments", middlewares.Chain(http.HandlerFunc(commentHandler.Create), authMiddleware))
	srv.RegisterRoutes("GET /api/v1/tasks/{id}/comments", middlewares.Chain(http.HandlerFunc(commentHandler.ListByTask), authMiddleware))
	srv.RegisterRoutes("DELETE /api/v1/comments/{id}", middlewares.Chain(http.HandlerFunc(commentHandler.Delete), authMiddleware))

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

	logger.Info("Se inicio TASK SERVICE",
		slog.String("addr", cfg.HTTPAddr),
	)

	return srv.Start()
}
