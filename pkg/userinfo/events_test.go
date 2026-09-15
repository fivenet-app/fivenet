package userinfo

import (
	"context"
	"fmt"
	"testing"
	"time"

	accounts "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs"
	pbtimestamp "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/internal/tests/nats"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/broker"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type labelEnricher struct{}

func TestCreateOrUpdateChangeConsumerUsesCanonicalStreamAndSubject(t *testing.T) {
	t.Parallel()

	_, js, shutdown, err := nats.NewInProcessNATSServer()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, shutdown()) })

	consumer, err := CreateOrUpdateChangeConsumer(
		t.Context(),
		js,
		jetstream.ConsumerConfig{
			Durable:       "test_userinfo_change_consumer",
			DeliverPolicy: jetstream.DeliverNewPolicy,
			AckPolicy:     jetstream.AckExplicitPolicy,
			FilterSubject: "wrong.subject",
		},
	)
	require.NoError(t, err)

	info, err := consumer.Info(t.Context())
	require.NoError(t, err)
	assert.Equal(t, UserInfoSubject, info.Config.FilterSubject)

	stream, err := js.Stream(t.Context(), UserInfoStreamName)
	require.NoError(t, err)
	streamInfo, err := stream.Info(t.Context())
	require.NoError(t, err)
	assert.Equal(t, jetstream.InterestPolicy, streamInfo.Config.Retention)
}

func (labelEnricher) EnrichJobInfo(user common.IJobInfo) {
	user.SetJobLabel(user.GetJob())
	user.SetJobGradeLabel(fmt.Sprintf("Rank %d", user.GetJobGrade()))
}

func (labelEnricher) EnrichJobInfoNoFallback(common.IJobInfo) {}

func (labelEnricher) EnrichJobName(common.IJobName) {}

func (labelEnricher) GetJobByName(string) *jobs.Job           { return nil }
func (labelEnricher) GetHighestJobGrade(string) (int32, bool) { return 0, false }

func (labelEnricher) GetJobGrade(string, int32) (*jobs.Job, *jobs.JobGrade) {
	return nil, nil
}

func TestBuildUserInfoChangedEvent(t *testing.T) {
	t.Parallel()

	evt := BuildUserInfoChangedEvent(
		42,
		77,
		nil,
		"police",
		3,
		labelEnricher{},
	)

	require.NotNil(t, evt)
	require.NotNil(t, evt.GetChangedAt())

	assert.Equal(t, int64(42), evt.GetAccountId())
	assert.Equal(t, int32(77), evt.GetUserId())
	assert.Equal(t, "police", evt.GetNewJob())
	assert.Equal(t, "police", evt.GetNewJobLabel())
	assert.Equal(t, int32(3), evt.GetNewJobGrade())
	assert.Equal(t, "Rank 3", evt.GetNewJobGradeLabel())
}

func TestBuildUserInfoChangedEventUsesProvidedTimestamp(t *testing.T) {
	t.Parallel()

	changedAt := pbtimestamp.New(time.Unix(123, 456))
	evt := BuildUserInfoChangedEvent(1, 2, changedAt, "ems", 4, nil)

	require.NotNil(t, evt)
	assert.Same(t, changedAt, evt.GetChangedAt())
	assert.Equal(t, "ems", evt.GetNewJob())
	assert.Equal(t, int32(4), evt.GetNewJobGrade())
}

func TestBuildAccountGroupsChangedEvent(t *testing.T) {
	t.Parallel()

	groups := &accounts.AccountGroups{Groups: []string{"supporter", "donator"}}
	evt := BuildAccountGroupsChangedEvent(42, nil, groups, true, false)

	require.NotNil(t, evt)
	require.NotNil(t, evt.GetChangedAt())
	require.NotNil(t, evt.GetNewGroups())

	assert.Equal(t, int64(42), evt.GetAccountId())
	assert.Equal(t, []string{"supporter", "donator"}, evt.GetNewGroups().GetGroups())
	assert.True(t, evt.GetCanBeSuperuser())
	assert.False(t, evt.GetCanBeConfigAdmin())

	groups.Groups[0] = "modified"
	assert.Equal(t, []string{"supporter", "donator"}, evt.GetNewGroups().GetGroups())
}

func TestChangesSubscriberClosesWhenItFallsBehind(t *testing.T) {
	t.Parallel()

	changes := &Changes{
		broker: broker.NewWithResyncOnSlowSubscriber[*pbuserinfo.UserInfoChanged](32),
	}
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	go changes.broker.Start(ctx)

	sub := changes.SubscribeUserInfoChanges()
	for i := range int32(34) {
		changes.broker.Publish(&pbuserinfo.UserInfoChanged{UserId: i + 1})
	}

	for range 32 {
		<-sub
	}
	_, ok := <-sub
	assert.False(t, ok)
}

func TestChangesPublishesAndConsumesUserInfoChangedEvent(t *testing.T) {
	t.Parallel()

	_, js, shutdown, err := nats.NewInProcessNATSServer()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, shutdown()) })

	changes := &Changes{
		logger: zap.NewNop(),
		js:     js,
		broker: broker.NewWithResyncOnSlowSubscriber[*pbuserinfo.UserInfoChanged](10),
	}
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	go changes.broker.Start(ctx)

	sub := changes.SubscribeUserInfoChanges()
	consumer, err := CreateOrUpdateChangeConsumer(ctx, js, jetstream.ConsumerConfig{
		Durable:       "changes_test_consumer",
		AckPolicy:     jetstream.AckExplicitPolicy,
		DeliverPolicy: jetstream.DeliverNewPolicy,
	})
	require.NoError(t, err)
	consumeCtx, err := consumer.Consume(changes.handleMessage)
	require.NoError(t, err)
	defer consumeCtx.Stop()

	want := &pbuserinfo.UserInfoChanged{AccountId: 42, UserId: 7}
	require.NoError(t, changes.PublishUserInfoChanged(ctx, want))

	select {
	case got := <-sub:
		require.NotNil(t, got)
		assert.Equal(t, want.AccountId, got.AccountId)
		assert.Equal(t, want.UserId, got.UserId)
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for consumed user-info event")
	}
}

func TestChangesTerminatesInvalidConsumedEvent(t *testing.T) {
	t.Parallel()

	_, js, shutdown, err := nats.NewInProcessNATSServer()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, shutdown()) })

	changes := &Changes{
		logger: zap.NewNop(),
		broker: broker.NewWithResyncOnSlowSubscriber[*pbuserinfo.UserInfoChanged](10),
	}
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	go changes.broker.Start(ctx)

	sub := changes.SubscribeUserInfoChanges()
	consumer, err := CreateOrUpdateChangeConsumer(ctx, js, jetstream.ConsumerConfig{
		Durable:       "invalid_changes_test_consumer",
		AckPolicy:     jetstream.AckExplicitPolicy,
		DeliverPolicy: jetstream.DeliverNewPolicy,
		AckWait:       50 * time.Millisecond,
	})
	require.NoError(t, err)
	consumeCtx, err := consumer.Consume(changes.handleMessage)
	require.NoError(t, err)
	defer consumeCtx.Stop()

	_, err = js.Publish(
		ctx,
		fmt.Sprintf("%s.%d.changes", BaseSubject, 42),
		[]byte(`{"accountId":"0","userId":7}`),
	)
	require.NoError(t, err)

	assert.Never(t, func() bool {
		select {
		case <-sub:
			return true
		default:
			return false
		}
	}, 200*time.Millisecond, 10*time.Millisecond)

	require.Eventually(t, func() bool {
		info, err := consumer.Info(ctx)
		return err == nil && info.NumAckPending == 0
	}, time.Second, 10*time.Millisecond)
}

func TestPublishUserInfoChangedRejectsInvalidEvents(t *testing.T) {
	t.Parallel()

	changes := &Changes{}
	require.Error(t, changes.PublishUserInfoChanged(t.Context(), nil))
	require.Error(
		t,
		changes.PublishUserInfoChanged(
			t.Context(),
			&pbuserinfo.UserInfoChanged{AccountId: 0, UserId: 1},
		),
	)
	require.Error(
		t,
		changes.PublishUserInfoChanged(
			t.Context(),
			&pbuserinfo.UserInfoChanged{AccountId: 1, UserId: 0},
		),
	)
}

func TestChangeConsumerRedeliversUnacknowledgedEventAfterRestart(t *testing.T) {
	t.Parallel()

	_, js, shutdown, err := nats.NewInProcessNATSServer()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, shutdown()) })

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	consumer, err := CreateOrUpdateChangeConsumer(ctx, js, jetstream.ConsumerConfig{
		Durable:       "handoff_test_consumer",
		AckPolicy:     jetstream.AckExplicitPolicy,
		DeliverPolicy: jetstream.DeliverNewPolicy,
		AckWait:       50 * time.Millisecond,
	})
	require.NoError(t, err)

	firstDelivery := make(chan struct{}, 1)
	consumeCtx, err := consumer.Consume(func(jetstream.Msg) {
		firstDelivery <- struct{}{}
	})
	require.NoError(t, err)

	_, err = js.Publish(
		ctx,
		fmt.Sprintf("%s.%d.changes", BaseSubject, 42),
		[]byte(`{"accountId":"1","userId":42}`),
	)
	require.NoError(t, err)
	select {
	case <-firstDelivery:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for first delivery")
	}
	consumeCtx.Stop()

	consumer, err = CreateOrUpdateChangeConsumer(ctx, js, jetstream.ConsumerConfig{
		Durable:       "handoff_test_consumer",
		AckPolicy:     jetstream.AckExplicitPolicy,
		DeliverPolicy: jetstream.DeliverNewPolicy,
		AckWait:       50 * time.Millisecond,
	})
	require.NoError(t, err)

	redelivered := make(chan struct{}, 1)
	consumeCtx, err = consumer.Consume(func(msg jetstream.Msg) {
		require.NoError(t, msg.Ack())
		redelivered <- struct{}{}
	})
	require.NoError(t, err)
	defer consumeCtx.Stop()

	select {
	case <-redelivered:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for redelivery after consumer restart")
	}
}
