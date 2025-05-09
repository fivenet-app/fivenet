package croner

import (
	"context"
	"fmt"
	"time"

	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	CronScheduleStreamName = "CRON_SCHEDULE"

	CronScheduleSubject events.Subject = "cron_schedule"
	CronScheduleTopic   events.Topic   = "schedule"
	CronCompleteTopic   events.Topic   = "complete"
)

func registerCronStreams(ctx context.Context, js *events.JSWrapper) error {
	if _, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:      CronScheduleStreamName,
		Storage:   jetstream.MemoryStorage,
		Retention: jetstream.InterestPolicy,
		Subjects:  []string{fmt.Sprintf("%s.>", CronScheduleSubject)},
		Discard:   jetstream.DiscardOld,
		MaxAge:    30 * time.Second,
	}); err != nil {
		return err
	}

	return nil
}
