package events

import (
	"context"
	"strings"

	"github.com/galexrt/fivenet/pkg/config"
	"github.com/nats-io/nats.go/jetstream"
)

const DescriptionPrefix = "FiveNet: "

// Ensures certain NATS config options are applied
type JSWrapper struct {
	jetstream.JetStream

	cfg config.NATS
}

func NewJSWrapper(js jetstream.JetStream, cfg config.NATS) *JSWrapper {
	return &JSWrapper{
		JetStream: js,
		cfg:       cfg,
	}
}

func (j *JSWrapper) CreateOrUpdateStream(ctx context.Context, cfg jetstream.StreamConfig) (jetstream.Stream, error) {
	if cfg.Replicas == 0 || cfg.Replicas > j.cfg.Replicas {
		cfg.Replicas = j.cfg.Replicas
	}

	if !strings.HasPrefix(cfg.Description, DescriptionPrefix) {
		cfg.Description = DescriptionPrefix + cfg.Description
	}

	return j.JetStream.CreateOrUpdateStream(ctx, cfg)
}

func (j *JSWrapper) CreateOrUpdateKeyValue(ctx context.Context, cfg jetstream.KeyValueConfig) (jetstream.KeyValue, error) {
	if cfg.Replicas == 0 || cfg.Replicas > j.cfg.Replicas {
		cfg.Replicas = j.cfg.Replicas
	}

	if !strings.HasPrefix(cfg.Description, DescriptionPrefix) {
		cfg.Description = DescriptionPrefix + cfg.Description
	}

	return j.JetStream.CreateOrUpdateKeyValue(ctx, cfg)
}
