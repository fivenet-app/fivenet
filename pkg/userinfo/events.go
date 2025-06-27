package userinfo

import (
	"context"
	"fmt"

	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	KVBucketName = "userinfo_poll_ttl"

	// BaseSubject is the base subject for user, job, object and system notification events.
	BaseSubject = "userinfo"

	PollStreamName = "POLL_REQUESTS"
	PollSubject    = "userinfo.poll.request"

	UserInfoStreamName = "USERINFO"
	UserInfoSubject    = "userinfo.*.changes"
)

func registerStreams(ctx context.Context, js *events.JSWrapper) error {
	// Stream for incoming poll requests
	pollStreamCfg := jetstream.StreamConfig{
		Name:        PollStreamName,
		Description: "Stream for userinfo poll requests",
		Subjects:    []string{PollSubject},
		Retention:   jetstream.InterestPolicy,
	}
	if _, err := js.CreateOrUpdateStream(ctx, pollStreamCfg); err != nil {
		return fmt.Errorf("failed to create/update stream %s: %w", PollStreamName, err)
	}

	// Stream for userinfo diffs
	userinfoStreamCfg := jetstream.StreamConfig{
		Name:        UserInfoStreamName,
		Description: "Stream for userinfo changes",
		Subjects:    []string{UserInfoSubject},
		Retention:   jetstream.InterestPolicy,
	}
	if _, err := js.CreateOrUpdateStream(ctx, userinfoStreamCfg); err != nil {
		return fmt.Errorf("failed to create/update stream %s: %w", UserInfoStreamName, err)
	}

	return nil
}
