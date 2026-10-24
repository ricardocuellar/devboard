package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/ricardocuellar/devboard/services/notification/internal/domain"
	"github.com/ricardocuellar/devboard/services/notification/internal/middlewares"
	"github.com/ricardocuellar/devboard/services/notification/internal/usecase"
)

type notificationResponse struct {
	ID        string    `json:"id"`
	TaskID    string    `json:"task_id"`
	ProjectID string    `json:"project_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

func ListNotifications(
	uc *usecase.NotificationQueryUseCase,
	logger *slog.Logger,
) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Query().Has("user_id") {
			writeProblem(
				writer, request, logger,
				http.StatusBadRequest,
				"El destinatario se obtiene del token; no se acepta user_id",
			)
			return
		}

		limit := int64(20)

		if value := request.URL.Query().Get("limit"); value != "" {
			parsed, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				writeProblem(
					writer, request, logger,
					http.StatusBadRequest,
					"limit debe ser un número entero válido",
				)
				return
			}
			limit = parsed
		}

		userID := middlewares.GetUserID(request)

		ctx, cancel := context.WithTimeout(
			request.Context(),
			3*time.Second,
		)
		defer cancel()

		notifications, err := uc.ListMine(ctx, userID, int32(limit))
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrUnauthorized):
				writeProblem(
					writer, request, logger,
					http.StatusUnauthorized,
					"Se requiere un usuario autenticado",
				)

			case errors.Is(err, domain.ErrInvalidInput):
				writeProblem(
					writer, request, logger,
					http.StatusBadRequest,
					err.Error(),
				)

			default:
				logger.ErrorContext(
					ctx,
					"error consultando notificaciones",
					slog.Any("error", err),
				)

				writeProblem(
					writer, request, logger,
					http.StatusInternalServerError,
					"Ha ocurrido un error inesperado",
				)
			}
			return
		}

		data := make([]notificationResponse, 0, len(notifications))

		for _, notification := range notifications {
			data = append(data, notificationResponse{
				ID:        notification.ID,
				TaskID:    notification.TaskID,
				ProjectID: notification.ProjectID,
				Message:   notification.Message,
				CreatedAt: notification.CreatedAt,
			})
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.Header().Set("Cache-Control", "no-store")

		response := struct {
			Data []notificationResponse `json:"data"`
		}{
			Data: data,
		}

		if err := json.NewEncoder(writer).Encode(response); err != nil {
			logger.ErrorContext(
				ctx,
				"error escribiendo notificaciones",
				slog.Any("error", err),
			)
		}
	}
}

func writeProblem(
	writer http.ResponseWriter,
	request *http.Request,
	logger *slog.Logger,
	status int,
	detail string,
) {
	writer.Header().Set("Content-Type", "application/problem+json")

	if status == http.StatusUnauthorized {
		writer.Header().Set("WWW-Authenticate", "Bearer")
	}

	writer.WriteHeader(status)

	problem := struct {
		Type     string `json:"type"`
		Title    string `json:"title"`
		Status   int    `json:"status"`
		Detail   string `json:"detail"`
		Instance string `json:"instance"`
	}{
		Type:     "about:blank",
		Title:    http.StatusText(status),
		Status:   status,
		Detail:   detail,
		Instance: request.URL.Path,
	}

	if err := json.NewEncoder(writer).Encode(problem); err != nil {
		logger.ErrorContext(
			request.Context(),
			"error escribiendo problem details",
			slog.Any("error", err),
		)
	}
}
