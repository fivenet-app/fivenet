package natsutils

import (
	"context"

	"github.com/nats-io/nats.go/jetstream"
)

const (
	Description = "FiveNet"
)

func CreateKeyValue(ctx context.Context, js jetstream.JetStream, config jetstream.KeyValueConfig) (jetstream.KeyValue, error) {
	kv, err := js.CreateOrUpdateKeyValue(ctx, config)
	if err != nil {
		return nil, err
	}

	return kv, err
}

func CreateOrUpdateStream(ctx context.Context, js jetstream.JetStream, config jetstream.StreamConfig) (jetstream.Stream, error) {
	stream, err := js.CreateOrUpdateStream(ctx, config)
	if err != nil {
		return nil, err
	}

	return stream, nil
}
