// Package main
package main

import (
	"os"

	"github.com/ricardocuellar/devboard/internal/handler"
	"github.com/ricardocuellar/devboard/internal/logger"
	"github.com/ricardocuellar/devboard/internal/middlewares"
	"github.com/ricardocuellar/devboard/internal/server"
)

func main() {

	// logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	// 	Level: slog.LevelInfo,
	// }))
	logger := logger.New(logger.DefaultConfig())

	// Aqui arranca el servidor
	srv := server.New(":8080", logger)

	srv.Use(middlewares.Recovery(logger))
	srv.Use(middlewares.Logger(logger))

	healthHander := handler.NewHealthHandler()

	srv.RegisterRoutes("GET /health", healthHander)
	srv.RegisterRoutes("GET /ready", healthHander)

	// Ruta temporal para probar Recovery
	// srv.RegisterRoutes("GET /panic", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	panic("error provocado para probar Recovery")
	// }))

	if err := srv.Start(); err != nil {
		logger.Error("error fatal", "error", err)
		os.Exit(1)
		// log.Fatalf("error al iniciar el servidor %v", err)
	}

}
