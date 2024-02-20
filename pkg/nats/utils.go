package natsutils

import (
	"context"
	"errors"

	"github.com/nats-io/nats.go"
)

const (
	Description = "FiveNet"
)

func CreateKeyValue(js nats.JetStreamContext, bucket string, config *nats.KeyValueConfig) (nats.KeyValue, error) {
	kv, err := js.KeyValue(bucket)
	if err != nil {
		if !errors.Is(err, nats.ErrBucketNotFound) {
			return nil, err
		}

		kv, err = js.CreateKeyValue(config)
		if err != nil {
			return nil, err
		}
	}

	return kv, err
}

func CreateOrUpdateStream(ctx context.Context, js nats.JetStreamContext, config *nats.StreamConfig) (*nats.StreamInfo, error) {
	sub, err := js.UpdateStream(config, nats.Context(ctx))
	if err != nil {
		if !errors.Is(err, nats.ErrStreamNotFound) {
			return nil, err
		}

		if _, err := js.AddStream(config, nats.Context(ctx)); err != nil {
			return nil, err
		}
	}

	return sub, nil
}
