package grpc

import (
	"fmt"
	"log/slog"
	"net"
	"time"

	workspacev1 "github.com/ricardocuellar/devboard/contracts/gen/go/devboard/workspace/v1"

	googlegrpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	server   *googlegrpc.Server
	listener net.Listener
	logger   *slog.Logger
}

func NewServer(
	addr string,
	handler *Handler,
	logger *slog.Logger,
) (*Server, error) {
	listener, err := net.Listen("tcp", addr)

	if err != nil {
		return nil,
			fmt.Errorf(
				"listen grpc %s: %w",
				addr,
				err,
			)
	}

	grpcServer := googlegrpc.NewServer()

	workspacev1.RegisterWorkspaceServiceServer(
		grpcServer,
		handler,
	)

	reflection.Register(
		grpcServer,
	)

	return &Server{
		server:   grpcServer,
		listener: listener,
		logger:   logger,
	}, nil
}

func (server *Server) Serve() error {
	server.logger.Info(
		"starting workspace grpc server",
		slog.String(
			"addr",
			server.listener.Addr().String(),
		),
	)

	if err := server.server.Serve(server.listener); err != nil {
		return fmt.Errorf(
			"serve grpc: %w",
			err,
		)
	}

	return nil
}

func (server *Server) Shutdown(timeout time.Duration) {
	done := make(chan struct{})

	go func() {
		server.server.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		server.logger.Info(
			"workspace grpc server stopped gracefully",
		)

	case <-time.After(timeout):
		server.logger.Warn(
			"workspace grpc graceful shutdown timed out",
		)

		server.server.Stop()
	}
}
