package croner

import (
	"context"
	"fmt"
	"time"

	"github.com/fivenet-app/fivenet/v2026/pkg/events"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	CronScheduleStreamName = "CRON_SCHEDULE"
	// CronExecutorConsumerName is shared by all executor instances so JetStream
	// distributes each schedule event to exactly one executor.
	CronExecutorConsumerName = "cron_executor"
	// CronSchedulerConsumerName is shared by all scheduler instances so each
	// completion event is reconciled exactly once.
	CronSchedulerConsumerName = "cron_scheduler"
	// CronScheduleMaxAge must exceed the maximum handler timeout and leave room
	// for consumer recovery and redelivery.
	CronScheduleMaxAge = 2 * time.Hour

	CronScheduleSubject events.Subject = "cron_schedule"
	CronScheduleTopic   events.Topic   = "schedule"
	CronCompleteTopic   events.Topic   = "complete"
)

func registerCronStreams(ctx context.Context, js *events.JSWrapper) error {
	if _, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:        CronScheduleStreamName,
		Description: "Cron schedule events stream",
		Storage:     jetstream.FileStorage,
		Retention:   jetstream.InterestPolicy,
		Subjects:    []string{fmt.Sprintf("%s.>", CronScheduleSubject)},
		Discard:     jetstream.DiscardOld,
		MaxAge:      CronScheduleMaxAge,
	}); err != nil {
		return err
	}

	return nil
}
