package workspace

import (
	"context"
	"log/slog"
	"time"

	"github.com/ricardocuellar/devboard/services/task/internal/requestid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func callerServiceUnaryInterceptor(serviceName string) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req any,
		reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {

		md := metadata.MD{}

		if existing, ok := metadata.FromOutgoingContext(ctx); ok {
			md = existing.Copy()
		}

		md.Set("x-caller-service", serviceName)

		if id := requestid.FromContext(ctx); id != "" {
			md.Set("x-request-id", id)
		}

		ctx = metadata.NewOutgoingContext(ctx, md)

		return invoker(
			ctx,
			method,
			req,
			reply,
			cc,
			opts...,
		)
	}
}

func loggingUnaryClientInterceptor(logger *slog.Logger) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req any,
		reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		start := time.Now()
		err := invoker(
			ctx,
			method,
			req,
			reply,
			cc,
			opts...,
		)

		duration := time.Since(start)

		code := status.Code(err)

		logger.InfoContext(
			ctx,
			"grpc client call completed",
			slog.String("request_id", requestid.FromContext(ctx)),
			slog.String("method", method),
			slog.String("status", code.String()),
			slog.Duration("duration", duration),
		)

		return err
	}
}
