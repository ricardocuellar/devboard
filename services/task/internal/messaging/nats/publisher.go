package nats

import (
	"context"
	"fmt"

	natsgo "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Publisher struct {
	client jetstream.JetStream
	stream string
}

func NewPublisher(
	connection *natsgo.Conn,
	stream string,
) (*Publisher, error) {
	if connection == nil {
		return nil, fmt.Errorf("la conexión NATS es requerida")
	}

	if stream == "" {
		return nil, fmt.Errorf("el nombre del stream es requerido")
	}

	client, err := jetstream.New(connection)
	if err != nil {
		return nil, fmt.Errorf("creando cliente JetStream: %w", err)
	}

	return &Publisher{
		client: client,
		stream: stream,
	}, nil
}

func (publisher *Publisher) Publish(
	ctx context.Context,
	eventID string,
	subject string,
	payload []byte,
) error {
	if eventID == "" || subject == "" || len(payload) == 0 {
		return fmt.Errorf(
			"event_id, subject y payload son requeridos",
		)
	}

	_, err := publisher.client.Publish(
		ctx,
		subject,
		payload,
		jetstream.WithMsgID(eventID),
		jetstream.WithExpectStream(publisher.stream),
	)
	if err != nil {
		return fmt.Errorf(
			"publicando evento %s: %w",
			eventID,
			err,
		)
	}

	return nil
}
