// Package main
package main

import (
	"log/slog"
	"os"
	"time"

	_ "github.com/ricardocuellar/devboard/docs"
	"github.com/ricardocuellar/devboard/internal/handler"
	"github.com/ricardocuellar/devboard/internal/logger"
	"github.com/ricardocuellar/devboard/internal/middlewares"
	"github.com/ricardocuellar/devboard/internal/server"
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
	env := os.Getenv("GO_ENV")
	slog.Info("entorno cargado", "GO_ENV", env)
	// 2. Logger
	var log *slog.Logger
	if env == "production" {
		log = logger.New(logger.ProductionConfig())
	} else {
		log = logger.New(logger.DefaultConfig())
	}

	// 3. Dependencies compartidas
	// validate := validator.New()
	// notifier := notification.NewLogNotifier(log)
	// _ = notifier

	// 4. Servidor con Functional options
	srv := server.New(":8080",
		server.WithLogger(log),
		server.WithReadTimeout(15*time.Second),
		server.WithWriteTimeout(30*time.Second),
	)

	// 5. Handlers
	healthHandler := handler.NewHealthHandler()

	// 6. Registro de rutas
	srv.RegisterRoutes("GET /docs/", httpSwagger.WrapHandler)
	srv.RegisterRoutes("GET /health", healthHandler)

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
