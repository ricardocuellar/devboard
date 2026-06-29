// Package server creates and start server
package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Server struct
type Server struct {
	httpServer *http.Server
	mux        *http.ServeMux
	logger     *slog.Logger
}

// New creates a new server
func New(addr string, logger *slog.Logger) *Server {
	mux := http.NewServeMux() //Enrutador HTTP

	httpServer := &http.Server{
		Addr:    addr,
		Handler: mux,

		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &Server{
		httpServer: httpServer,
		mux:        mux,
		logger:     logger,
	}
}

// Start server
func (server *Server) Start() error {

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	defer signal.Stop(quit)

	serveErr := make(chan error, 1)
	go func() {
		server.logger.Info("servidor iniciado", slog.String("addr", server.httpServer.Addr))
		if err := server.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serveErr <- err
		}
	}()

	select {
	case err := <-serveErr:
		return err
	case sig := <-quit:
		server.logger.Info("señal recibida, iniciando shutdown",
			slog.String("signal", sig.String()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.httpServer.Shutdown(ctx); err != nil {
		server.logger.Error("error en shutdown", slog.Any("error", err))
		return err
	}

	server.logger.Info("Servidor detenido correctamente")

	return nil
}

// RegisterRoutes method
func (server *Server) RegisterRoutes(pattern string, handler http.Handler) {
	server.mux.Handle(pattern, handler)
}

// ServeHTTP server
func (server *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	server.mux.ServeHTTP(writer, request)
}

// Use function for middlewares
func (server *Server) Use(middleware func(http.Handler) http.Handler) {
	server.httpServer.Handler = middleware(server.httpServer.Handler)
}
