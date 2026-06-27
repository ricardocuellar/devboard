// Package main
package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/ricardocuellar/devboard/internal/handler"
	"github.com/ricardocuellar/devboard/internal/middlewares"
	"github.com/ricardocuellar/devboard/internal/server"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Aqui arranca el servidor
	srv := server.New(":8080")

	srv.Use(middlewares.Logger(logger))

	healthHander := handler.NewHealthHandler()

	srv.RegisterRoutes("GET /health", healthHander)
	srv.RegisterRoutes("GET /ready", healthHander)

	// log.Println("Servidor iniciando en :8080")
	logger.Info("servidor iniciando", slog.String("addr", ":8080"))

	if err := srv.Start(); err != nil {
		log.Fatalf("error al iniciar el servidor %v", err)
	}

}
