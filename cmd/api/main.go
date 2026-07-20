// Package main
package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/ricardocuellar/devboard/docs"
	"github.com/ricardocuellar/devboard/internal/handler"
	"github.com/ricardocuellar/devboard/internal/logger"
	"github.com/ricardocuellar/devboard/internal/middlewares"
	"github.com/ricardocuellar/devboard/internal/repository/memory"
	"github.com/ricardocuellar/devboard/internal/server"
	"github.com/ricardocuellar/devboard/internal/usecase"
	"github.com/ricardocuellar/devboard/internal/validator"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Devboard API
// @version 1.0
// @description API REST para gestión de proyectos
// @contact.name Soporte Devboard Ricardo
// @host  localhost:8080
// @BasePath /api/v1
func main() {

	// 1. Configuración del entorno
	// 2. Logger
	log := setupLogger()

	// 3. Dependencies compartidas

	// Infrastructure: adapters
	userRepo := memory.NewUserRepository()
	taskRepo := memory.NewTaskRepository()

	// Use cases: lógica del negocio
	userUC := usecase.NewUserUseCase(userRepo)
	taskUC := usecase.NewTaskUseCase(taskRepo, taskRepo)

	// Validación compartida
	validate := validator.New()

	// 4. Handlers
	healthHandler := handler.NewHealthHandler()
	userHandler := handler.NewUserHandler(userUC, validate, log)
	taskHandler := handler.NewTaskHandler(taskUC, validate, log)

	// notifier := notification.NewLogNotifier(log)
	// _ = notifier

	// 5. Servidor con Functional options
	srv := server.New(":8080",
		server.WithLogger(log),
		server.WithReadTimeout(15*time.Second),
		server.WithWriteTimeout(30*time.Second),
	)

	// 6. Registro de rutas
	srv.RegisterRoutes("GET /docs/", httpSwagger.WrapHandler)
	srv.RegisterRoutes("GET /health", healthHandler)
	srv.RegisterRoutes("POST /api/v1/users", http.HandlerFunc(userHandler.Create))
	srv.RegisterRoutes("GET /api/v1/users/{id}", http.HandlerFunc(userHandler.Get))
	srv.RegisterRoutes("POST /api/v1/tasks", http.HandlerFunc(taskHandler.Create))
	srv.RegisterRoutes("GET /api/v1/tasks/{id}", http.HandlerFunc(taskHandler.Get))
	srv.RegisterRoutes("GET /api/v1/projects/{id}/tasks", http.HandlerFunc(taskHandler.ListByProject))

	// GET /api/v1/tasks/{id}
	// GET /api/v1/projects/{id}/tasks
	// PUT /api/v1/tasks/{id}/assign
	// PUT /api/v1/tasks/{id}/status

	// 7. Middleware chain
	srv.UseChain(
		middlewares.Recovery(log),
		middlewares.Logger(log),
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
