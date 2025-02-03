package discord

import (
	"context"
	"fmt"
	"time"

	"github.com/fivenet-app/fivenet/pkg/events"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	StreamName = "DISCORD_BOT"

	BaseSubject events.Subject = "discord_bot"

	TriggerSyncType events.Type = "trigger_sync"
)

func registerStreams(ctx context.Context, js *events.JSWrapper) error {
	cfg := jetstream.StreamConfig{
		Name:        StreamName,
		Description: "Discord Bot Events",
		Retention:   jetstream.InterestPolicy,
		Subjects:    []string{fmt.Sprintf("%s.>", BaseSubject)},
		Discard:     jetstream.DiscardOld,
		MaxAge:      2 * time.Minute,
		Storage:     jetstream.MemoryStorage,
	}
	if _, err := js.CreateOrUpdateStream(ctx, cfg); err != nil {
		return err
	}

	return nil
}
