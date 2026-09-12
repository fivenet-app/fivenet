package userinfo

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/events"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/broker"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/instance"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/protoutils"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var ChangesModule = fx.Module(
	"userinfo.changes",
	fx.Provide(NewChanges),
)

// ChangePublisher publishes canonical user-info changes for all processes.
type ChangePublisher interface {
	PublishUserInfoChanged(ctx context.Context, event *pbuserinfo.UserInfoChanged) error
}

// ChangeSubscriber provides process-local fanout of canonical user-info changes.
// A closed subscription means the subscriber must reload authoritative state.
type ChangeSubscriber interface {
	SubscribeUserInfoChanges() chan *pbuserinfo.UserInfoChanged
	UnsubscribeUserInfoChanges(chan *pbuserinfo.UserInfoChanged)
}

type Changes struct {
	logger *zap.Logger
	js     *events.JSWrapper
	broker *broker.Broker[*pbuserinfo.UserInfoChanged]
	wg     sync.WaitGroup
	jsCons jetstream.ConsumeContext
}

type ChangesParams struct {
	fx.In

	LC     fx.Lifecycle
	Logger *zap.Logger
	JS     *events.JSWrapper
}

type ChangesResult struct {
	fx.Out

	Changes    *Changes
	Publisher  ChangePublisher
	Subscriber ChangeSubscriber
}

func NewChanges(p ChangesParams) ChangesResult {
	ctx, cancel := context.WithCancel(context.Background())
	c := &Changes{
		logger: p.Logger.Named("userinfo.changes"),
		js:     p.JS,
		broker: broker.NewWithResyncOnSlowSubscriber[*pbuserinfo.UserInfoChanged](32),
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		if err := registerUserInfoStream(ctxStartup, c.js); err != nil {
			return fmt.Errorf("failed to register user info streams. %w", err)
		}

		c.wg.Go(func() {
			c.broker.Start(ctx)
		})

		if err := c.registerSubscription(ctxStartup, ctx); err != nil {
			cancel()
			c.wg.Wait()
			return err
		}

		return nil
	}))

	p.LC.Append(fx.StopHook(func(context.Context) error {
		cancel()
		if c.jsCons != nil {
			c.jsCons.Stop()
		}
		c.wg.Wait()
		return nil
	}))

	return ChangesResult{Changes: c, Publisher: c, Subscriber: c}
}

func (c *Changes) PublishUserInfoChanged(
	ctx context.Context,
	event *pbuserinfo.UserInfoChanged,
) error {
	if event == nil || event.GetAccountId() <= 0 || event.GetUserId() <= 0 {
		return errors.New("user info change requires account and user IDs")
	}

	if _, err := c.js.PublishProto(
		ctx,
		fmt.Sprintf("%s.%d.changes", BaseSubject, event.GetAccountId()),
		event,
	); err != nil {
		getChangesMetrics().publishFailures.Inc()
		return fmt.Errorf("failed to publish user info change. %w", err)
	}

	return nil
}

func (c *Changes) SubscribeUserInfoChanges() chan *pbuserinfo.UserInfoChanged {
	return c.broker.Subscribe()
}

func (c *Changes) UnsubscribeUserInfoChanges(ch chan *pbuserinfo.UserInfoChanged) {
	c.broker.Unsubscribe(ch)
}

func (c *Changes) registerSubscription(ctxStartup context.Context, ctx context.Context) error {
	consumer, err := c.js.CreateOrUpdateConsumer(
		ctxStartup,
		UserInfoStreamName,
		jetstream.ConsumerConfig{
			Durable:           instance.ID() + "_ui_changes",
			AckPolicy:         jetstream.AckExplicitPolicy,
			FilterSubject:     UserInfoSubject,
			InactiveThreshold: time.Minute,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create user info changes consumer. %w", err)
	}

	c.jsCons, err = consumer.Consume(c.handleMessage,
		c.js.ConsumeErrHandlerWithRestart(ctx, c.logger, c.registerSubscription))
	if err != nil {
		return fmt.Errorf("failed to start user info changes consumer. %w", err)
	}

	return nil
}

func (c *Changes) handleMessage(msg jetstream.Msg) {
	defer func() {
		if err := msg.Ack(); err != nil {
			c.logger.Error(
				"failed to ack user info change",
				zap.Error(err),
				zap.String("subject", msg.Subject()),
			)
		}
	}()

	event := &pbuserinfo.UserInfoChanged{}
	if err := protoutils.UnmarshalPartialJSON(msg.Data(), event); err != nil {
		c.logger.Error(
			"failed to unmarshal user info change",
			zap.Error(err),
			zap.String("subject", msg.Subject()),
		)
		return
	}

	c.broker.Publish(event)
}
