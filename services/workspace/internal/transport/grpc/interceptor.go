package grpc

import (
	"context"
	"log/slog"
	"time"

	"github.com/ricardocuellar/devboard/services/workspace/internal/requestid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func loggingUnaryServerInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (
		any,
		error,
	) {

		md, _ := metadata.FromIncomingContext(ctx)

		var id string

		if values := md.Get("x-request-id"); len(values) > 0 {
			id = values[0]
		}

		ctx = requestid.WithContext(ctx, id)

		start := time.Now()

		callerService := callerServiceFromContext(ctx)

		response, err := handler(ctx, req)

		duration := time.Since(start)

		code := status.Code(err)

		logger.InfoContext(ctx, "grpc server call completed",
			slog.String("request_id", requestid.FromContext(ctx)),
			slog.String("method", info.FullMethod),
			slog.String(
				"caller_service",
				callerService,
			),
			slog.String(
				"status",
				code.String(),
			),
			slog.Duration(
				"duration",
				duration,
			),
		)

		return response, err
	}
}

func callerServiceFromContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return "unknown"
	}

	values := md.Get("x-caller-service")

	if len(values) == 0 {
		return "unknown"
	}

	return values[0]
}
