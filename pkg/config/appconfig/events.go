package appconfig

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/pkg/events"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const (
	BaseSubject events.Subject = "appconfig"

	UpdateSubject events.Type = "update"
)

func (c *Config) registerSubscriptions(ctx context.Context, bc context.Context) error {
	cfg := jetstream.StreamConfig{
		Name:        "APPCONFIG",
		Description: "AppConfig update events",
		Retention:   jetstream.InterestPolicy,
		Subjects:    []string{fmt.Sprintf("%s.>", BaseSubject)},
		Discard:     jetstream.DiscardOld,
		MaxAge:      10 * time.Second,
		Storage:     jetstream.MemoryStorage,
	}

	if _, err := c.js.CreateOrUpdateStream(ctx, cfg); err != nil {
		return err
	}

	consumer, err := c.js.CreateConsumer(ctx, cfg.Name, jetstream.ConsumerConfig{
		DeliverPolicy: jetstream.DeliverNewPolicy,
		FilterSubject: fmt.Sprintf("%s.>", BaseSubject),
	})
	if err != nil {
		return err
	}

	if c.jsCons != nil {
		c.jsCons.Stop()
		c.jsCons = nil
	}

	c.jsCons, err = consumer.Consume(c.handleMessageFunc(bc),
		c.js.ConsumeErrHandlerWithRestart(context.Background(), c.logger,
			func(_ context.Context, ctx context.Context) error {
				return c.registerSubscriptions(bc, bc)
			}))
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) handleMessageFunc(ctx context.Context) jetstream.MessageHandler {
	return func(msg jetstream.Msg) {
		remoteCtx, err := events.GetJetstreamMsgContext(msg)
		if err != nil {
			c.logger.Error("failed to get js msg context", zap.Error(err))
		}
		_, span := c.tracer.Start(trace.ContextWithRemoteSpanContext(ctx, remoteCtx), msg.Subject())
		defer span.End()

		if err := msg.Ack(); err != nil {
			c.logger.Error("failed to ack message", zap.Error(err))
		}

		split := strings.Split(msg.Subject(), ".")

		if len(split) < 2 {
			c.logger.Warn("unknown app config subject received", zap.String("subject", msg.Subject()))
			return
		}

		if split[1] == string(UpdateSubject) {
			cfg, err := c.updateConfigFromDB(ctx)
			if err != nil {
				c.logger.Error("failed to update app config from db", zap.Error(err))
				return
			}

			c.broker.Publish(cfg)
		}
	}
}
