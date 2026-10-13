package grpc

import (
	"fmt"
	"log/slog"
	"net"
	"time"

	workspacev1 "github.com/ricardocuellar/devboard/contracts/gen/go/devboard/workspace/v1"

	googlegrpc "google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	server       *googlegrpc.Server
	healthServer *health.Server
	listener     net.Listener
	logger       *slog.Logger
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

	healthServer := health.NewServer()

	healthgrpc.RegisterHealthServer(grpcServer, healthServer)

	healthServer.SetServingStatus("", healthgrpc.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus(workspacev1.WorkspaceService_ServiceDesc.ServiceName, healthgrpc.HealthCheckResponse_SERVING)

	reflection.Register(
		grpcServer,
	)

	return &Server{
		server:       grpcServer,
		healthServer: healthServer,
		listener:     listener,
		logger:       logger,
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

	server.healthServer.Shutdown()

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
