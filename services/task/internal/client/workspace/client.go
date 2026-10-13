package workspace

import (
	"context"
	"errors"
	"fmt"
	"time"

	workspacev1 "github.com/ricardocuellar/devboard/contracts/gen/go/devboard/workspace/v1"
	"github.com/ricardocuellar/devboard/internal/domain"
	"github.com/ricardocuellar/devboard/services/task/internal/usecase"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Client struct {
	connection *grpc.ClientConn
	client     workspacev1.WorkspaceServiceClient
	timeout    time.Duration
}

func NewClient(target string, timeout time.Duration) (*Client, error) {

	if timeout <= 0 {
		return nil, fmt.Errorf("El timeout de GRPC Workspace debe ser positivo")
	}

	connection, err :=
		grpc.NewClient(
			target,
			grpc.WithTransportCredentials(
				insecure.NewCredentials(),
			),
		)

	if err != nil {
		return nil,
			fmt.Errorf(
				"create workspace grpc client: %w",
				err,
			)
	}

	return &Client{
		connection: connection,
		client:     workspacev1.NewWorkspaceServiceClient(connection),
		timeout:    timeout,
	}, nil
}

func (client *Client) ValidateProjectAccess(ctx context.Context, projectID string, userID string) (usecase.ProjectAccess, error) {

	if err := ctx.Err(); err != nil {
		return usecase.ProjectAccess{}, err
	}

	rpcCtx, cancel := context.WithTimeout(ctx, client.timeout)

	defer cancel()

	response, err := client.client.ValidateProjectAccess(ctx, &workspacev1.ValidateProjectAccessRequest{
		ProjectId: projectID,
		UserId:    userID,
	})

	if err != nil {
		return usecase.ProjectAccess{}, mapError(rpcCtx, err)
	}

	return usecase.ProjectAccess{
		WorkpaceID: response.GetWorkspaceId(),
	}, nil
}

func (client *Client) Close() error {
	if err := client.connection.Close(); err != nil {
		return fmt.Errorf("Error al cerrar conexión grpc: %w", err)
	}

	return nil
}

func mapError(ctx context.Context, err error) error {

	if errors.Is(ctx.Err(), context.Canceled) {
		return context.Canceled
	}

	code := status.Code(err)
	switch code {
	case codes.InvalidArgument:
		return domain.ErrInvalidInput

	case codes.NotFound:
		return domain.ErrNotFound

	case codes.PermissionDenied:
		return domain.ErrForbidden

	case codes.Unavailable:
		return fmt.Errorf("Servicio de Workspace no disponible: %w", err)

	case codes.Canceled:
		return context.Canceled

	default:
		return fmt.Errorf("Servicio de Workspace gRPC falló: %w", err)
	}
}
