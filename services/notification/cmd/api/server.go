package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func serve(server *http.Server, logger *slog.Logger) error {
	signalCtx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return fmt.Errorf("ejecutando servidor HTTP: %w", err)

	case <-signalCtx.Done():
		logger.Info("iniciando cierre del Notification")
	}

	stop()

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error(
			"se agotó o falló el cierre ordenado; forzando cierre",
			slog.Any("error", err),
		)

		_ = server.Close()

		return fmt.Errorf("cerrando servidor HTTP: %w", err)
	}

	err := <-serverErrors
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("deteniendo servidor HTTP: %w", err)
	}

	logger.Info("Notification detenido correctamente")

	return nil
}
