// Package server creates and start server
package server

import (
	"net/http"
	"time"
)

// Server struct
type Server struct {
	httpServer *http.Server
	mux        *http.ServeMux
}

// New creates a new server
func New(addr string) *Server {
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
	}
}

// Start server
func (server *Server) Start() error {
	return server.httpServer.ListenAndServe()
}

// RegisterRoutes method
func (server *Server) RegisterRoutes(pattern string, handler http.Handler) {
	server.mux.Handle(pattern, handler)
}

// ServeHTTP server
func (server *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	server.mux.ServeHTTP(writer, request)
}
