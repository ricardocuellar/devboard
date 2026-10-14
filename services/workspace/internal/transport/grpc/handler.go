package grpc

import (
	"context"
	"errors"
	"log/slog"

	workspacev1 "github.com/ricardocuellar/devboard/contracts/gen/go/devboard/workspace/v1"
	"github.com/ricardocuellar/devboard/services/workspace/internal/domain"
	"github.com/ricardocuellar/devboard/services/workspace/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	workspacev1.UnimplementedWorkspaceServiceServer

	projectAccess usecase.ProjectAccessValidator
	logger        *slog.Logger
}

func NewHandler(
	projectAccess usecase.ProjectAccessValidator,
	logger *slog.Logger,
) *Handler {
	return &Handler{
		projectAccess: projectAccess,

		logger: logger,
	}
}

func (handler *Handler) ValidateProjectAccess(
	ctx context.Context,
	request *workspacev1.ValidateProjectAccessRequest,
) (
	*workspacev1.ValidateProjectAccessResponse,
	error,
) {
	if request == nil {
		return nil,
			status.Error(
				codes.InvalidArgument,
				"request requerida",
			)
	}

	result, err :=
		handler.projectAccess.Validate(
			ctx,
			request.GetProjectId(),
			request.GetUserId(),
		)

	if err != nil {
		return nil,
			handler.mapError(
				ctx,
				err,
			)
	}

	return &workspacev1.ValidateProjectAccessResponse{
		WorkspaceId: result.WorkspaceID,
	}, nil
}

func (handler *Handler) mapError(
	ctx context.Context,
	err error,
) error {

	switch {
	case errors.Is(err, context.Canceled):
		return status.Error(codes.Canceled, "request cancelada")

	case errors.Is(err, context.DeadlineExceeded):
		return status.Error(codes.DeadlineExceeded, "deadline excedido")

	case errors.Is(
		err,
		domain.ErrInvalidInput,
	):
		return status.Error(
			codes.InvalidArgument,
			"project_id y user_id son requeridos",
		)

	case errors.Is(
		err,
		domain.ErrNotFound,
	):
		return status.Error(
			codes.NotFound,
			"project no encontrado",
		)

	case errors.Is(
		err,
		domain.ErrForbidden,
	):
		return status.Error(
			codes.PermissionDenied,
			"el usuario no tiene acceso al project",
		)

	default:
		handler.logger.ErrorContext(
			ctx,
			"workspace grpc request failed",
			slog.Any(
				"error",
				err,
			),
		)

		return status.Error(
			codes.Internal,
			"error interno",
		)
	}
}
