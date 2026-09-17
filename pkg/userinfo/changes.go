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
	PublishAccountGroupsChanged(ctx context.Context, event *pbuserinfo.AccountGroupsChanged) error
}

// ChangeSubscriber provides process-local fanout of canonical user-info changes.
// A closed subscription means the subscriber must reload authoritative state.
type ChangeSubscriber interface {
	SubscribeUserInfoChanges() chan *pbuserinfo.UserInfoChanged
	UnsubscribeUserInfoChanges(chan *pbuserinfo.UserInfoChanged)
	SubscribeAccountGroupsChanges() chan *pbuserinfo.AccountGroupsChanged
	UnsubscribeAccountGroupsChanges(chan *pbuserinfo.AccountGroupsChanged)
}

type Changes struct {
	logger              *zap.Logger
	js                  *events.JSWrapper
	userInfoBroker      *broker.Broker[*pbuserinfo.UserInfoChanged]
	accountGroupsBroker *broker.Broker[*pbuserinfo.AccountGroupsChanged]
	wg                  sync.WaitGroup
	userInfoJSCons      jetstream.ConsumeContext
	accountGroupsJSCons jetstream.ConsumeContext
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
		logger:         p.Logger.Named("userinfo.changes"),
		js:             p.JS,
		userInfoBroker: broker.NewWithResyncOnSlowSubscriber[*pbuserinfo.UserInfoChanged](32),
		accountGroupsBroker: broker.NewWithResyncOnSlowSubscriber[*pbuserinfo.AccountGroupsChanged](
			32,
		),
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		c.wg.Go(func() {
			c.userInfoBroker.Start(ctx)
		})
		c.wg.Go(func() {
			c.accountGroupsBroker.Start(ctx)
		})

		if err := c.registerUserInfoSubscription(ctxStartup, ctx); err != nil {
			cancel()
			c.wg.Wait()
			return err
		}
		if err := c.registerAccountGroupsSubscription(ctxStartup, ctx); err != nil {
			cancel()
			c.wg.Wait()
			return err
		}

		return nil
	}))

	p.LC.Append(fx.StopHook(func(context.Context) error {
		cancel()
		if c.userInfoJSCons != nil {
			c.userInfoJSCons.Stop()
		}
		if c.accountGroupsJSCons != nil {
			c.accountGroupsJSCons.Stop()
		}
		c.wg.Wait()
		return nil
	}))

	return ChangesResult{Changes: c, Publisher: c, Subscriber: c}
}

func (c *Changes) PublishAccountGroupsChanged(
	ctx context.Context,
	event *pbuserinfo.AccountGroupsChanged,
) error {
	if event == nil || event.GetAccountId() <= 0 {
		return errors.New("account group change requires an account ID")
	}

	if _, err := c.js.PublishProto(
		ctx,
		fmt.Sprintf("%s.%d.account_groups", BaseSubject, event.GetAccountId()),
		event,
	); err != nil {
		getChangesMetrics().publishFailures.Inc()
		return fmt.Errorf("failed to publish account group change: %w", err)
	}

	return nil
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
	return c.userInfoBroker.Subscribe()
}

func (c *Changes) UnsubscribeUserInfoChanges(ch chan *pbuserinfo.UserInfoChanged) {
	c.userInfoBroker.Unsubscribe(ch)
}

func (c *Changes) SubscribeAccountGroupsChanges() chan *pbuserinfo.AccountGroupsChanged {
	return c.accountGroupsBroker.Subscribe()
}

func (c *Changes) UnsubscribeAccountGroupsChanges(ch chan *pbuserinfo.AccountGroupsChanged) {
	c.accountGroupsBroker.Unsubscribe(ch)
}

func (c *Changes) registerUserInfoSubscription(
	ctxStartup context.Context,
	ctx context.Context,
) error {
	consumer, err := CreateOrUpdateUserInfoChangeConsumer(
		ctxStartup,
		c.js,
		jetstream.ConsumerConfig{
			Durable:           instance.ID() + "_ui_changes",
			AckPolicy:         jetstream.AckExplicitPolicy,
			InactiveThreshold: time.Minute,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create user info changes consumer. %w", err)
	}

	c.userInfoJSCons, err = consumer.Consume(c.handleUserInfoMessage,
		c.js.ConsumeErrHandlerWithRestart(ctx, c.logger, c.registerUserInfoSubscription))
	if err != nil {
		return fmt.Errorf("failed to start user info changes consumer. %w", err)
	}

	return nil
}

func (c *Changes) registerAccountGroupsSubscription(
	ctxStartup context.Context,
	ctx context.Context,
) error {
	consumer, err := CreateOrUpdateAccountGroupsChangeConsumer(
		ctxStartup,
		c.js,
		jetstream.ConsumerConfig{
			Durable:           instance.ID() + "_account_groups_changes",
			AckPolicy:         jetstream.AckExplicitPolicy,
			InactiveThreshold: time.Minute,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create account group changes consumer. %w", err)
	}

	c.accountGroupsJSCons, err = consumer.Consume(c.handleAccountGroupsMessage,
		c.js.ConsumeErrHandlerWithRestart(ctx, c.logger, c.registerAccountGroupsSubscription))
	if err != nil {
		return fmt.Errorf("failed to start account group changes consumer. %w", err)
	}

	return nil
}

func (c *Changes) handleUserInfoMessage(msg jetstream.Msg) {
	event := &pbuserinfo.UserInfoChanged{}
	if err := protoutils.UnmarshalPartialJSON(msg.Data(), event); err != nil {
		c.logger.Error(
			"failed to unmarshal user info change",
			zap.Error(err),
			zap.String("subject", msg.Subject()),
		)
		// Malformed events cannot succeed on retry, so terminate them instead
		// of acknowledging and silently losing them.
		if err := msg.Term(); err != nil {
			c.logger.Error("failed to terminate malformed user info change", zap.Error(err))
		}
		return
	}
	if event.GetUserId() <= 0 || event.GetAccountId() <= 0 {
		c.logger.Error(
			"received invalid user info change",
			zap.Int32("userId", event.GetUserId()),
			zap.Int64("accountId", event.GetAccountId()),
			zap.String("subject", msg.Subject()),
		)
		if err := msg.Term(); err != nil {
			c.logger.Error("failed to terminate invalid user info change", zap.Error(err))
		}
		return
	}

	c.userInfoBroker.Publish(event)

	if err := msg.Ack(); err != nil {
		c.logger.Error(
			"failed to ack user info change",
			zap.Error(err),
			zap.String("subject", msg.Subject()),
		)
	}
}

func (c *Changes) handleAccountGroupsMessage(msg jetstream.Msg) {
	event := &pbuserinfo.AccountGroupsChanged{}
	if err := protoutils.UnmarshalPartialJSON(msg.Data(), event); err != nil {
		c.logger.Error(
			"failed to unmarshal account group change",
			zap.Error(err),
			zap.String("subject", msg.Subject()),
		)
		if err := msg.Term(); err != nil {
			c.logger.Error("failed to terminate malformed account group change", zap.Error(err))
		}
		return
	}
	if event.GetAccountId() <= 0 {
		c.logger.Error(
			"received invalid account group change",
			zap.Int64("account_id", event.GetAccountId()),
			zap.String("subject", msg.Subject()),
		)
		if err := msg.Term(); err != nil {
			c.logger.Error("failed to terminate invalid account group change", zap.Error(err))
		}
		return
	}

	c.accountGroupsBroker.Publish(event)
	if err := msg.Ack(); err != nil {
		c.logger.Error(
			"failed to ack account group change",
			zap.Error(err),
			zap.String("subject", msg.Subject()),
		)
	}
}
