// Package notification strategy
package notification

import (
	"context"
	"log/slog"
)

// Notifier es la interfaz que define el contrato de notificación
type Notifier interface {
	Notify(ctx context.Context, userID string, message string) error
}

// LogNotifier registra la notificación en los logs
type LogNotifier struct {
	logger *slog.Logger
}

// NewLogNotifier constructor de Lognotifier
func NewLogNotifier(logger *slog.Logger) *LogNotifier {
	return &LogNotifier{logger: logger}
}

// Notify es el método para enviar notificaciones en los logs
func (notifier *LogNotifier) Notify(ctx context.Context, userID string, message string) error {
	notifier.logger.InfoContext(ctx, "notificación enviada",
		slog.String("user_id", userID),
		slog.String("message", message),
	)

	return nil
}

// NoOpNotifier struct que no notifica pero ayuda a probar
type NoOpNotifier struct{}

// Notify método de NoOPNotifier
func (notifier *NoOpNotifier) Notify(_ context.Context, _ string, _ string) error {
	return nil
}
