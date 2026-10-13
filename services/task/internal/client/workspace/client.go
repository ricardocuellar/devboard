package workspace

import (
	"fmt"

	workspacev1 "github.com/ricardocuellar/devboard/contracts/gen/go/devboard/workspace/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	connection *grpc.ClientConn
	client     workspacev1.WorkspaceServiceClient
}

func NewClient(target string) (*Client, error) {
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
	}, nil
}
