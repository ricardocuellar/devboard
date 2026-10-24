package router

import (
	"log/slog"
	"net/http"

	"github.com/ricardocuellar/devboard/services/gateway/internal/config"
	"github.com/ricardocuellar/devboard/services/gateway/internal/proxy"
)

func Register(
	mux *http.ServeMux,
	cfg config.Config,
	logger *slog.Logger,
) {
	identityProxy := proxy.New(cfg.IdentityURL, cfg.ResponseHeaderTimeout, logger)
	workspaceProxy := proxy.New(cfg.WorkspaceURL, cfg.ResponseHeaderTimeout, logger)
	taskProxy := proxy.New(cfg.TaskURL, cfg.ResponseHeaderTimeout, logger)
	notificationProxy := proxy.New(cfg.NotificationURL, cfg.ResponseHeaderTimeout, logger)

	// Identity.
	mux.Handle("/api/v1/auth/", identityProxy)
	mux.Handle("/api/v1/users", identityProxy)
	mux.Handle("/api/v1/users/", identityProxy)

	// Workspace.
	mux.Handle("/api/v1/workspaces", workspaceProxy)
	mux.Handle("/api/v1/workspaces/", workspaceProxy)
	mux.Handle("/api/v1/projects", workspaceProxy)
	mux.Handle("/api/v1/projects/", workspaceProxy)

	// Las tareas de un proyecto pertenecen a Task.
	mux.Handle("/api/v1/projects/{id}/tasks", taskProxy)

	// Task y comentarios.
	mux.Handle("/api/v1/tasks", taskProxy)
	mux.Handle("/api/v1/tasks/", taskProxy)
	mux.Handle("/api/v1/comments", taskProxy)
	mux.Handle("/api/v1/comments/", taskProxy)

	// Notification
	mux.Handle("/api/v1/notifications", notificationProxy)
}
