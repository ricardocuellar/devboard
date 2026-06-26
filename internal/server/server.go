package server

import (
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
	mux        *http.ServeMux
}

// Crea un nuevo server
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

// Inicia el servidor
func (server *Server) Start() error {
	return server.httpServer.ListenAndServe()
}
