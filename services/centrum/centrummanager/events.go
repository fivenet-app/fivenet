package centrummanager

import (
	"context"
	"fmt"

	eventscentrum "github.com/fivenet-app/fivenet/v2025/services/centrum/events"
	"github.com/nats-io/nats.go/jetstream"
)

func (s *Manager) registerStream(ctx context.Context) (jetstream.StreamConfig, error) {
	streamCfg, err := eventscentrum.RegisterStream(ctx, s.js)
	if err != nil {
		return streamCfg, fmt.Errorf("failed to register stream. %w", err)
	}

	return streamCfg, nil
}
