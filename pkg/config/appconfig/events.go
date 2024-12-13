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

func (c *Config) registerSubscriptions(ctxStartup context.Context, ctxCancel context.Context) error {
	cfg := jetstream.StreamConfig{
		Name:        "APPCONFIG",
		Description: "AppConfig update events",
		Retention:   jetstream.InterestPolicy,
		Subjects:    []string{fmt.Sprintf("%s.>", BaseSubject)},
		Discard:     jetstream.DiscardOld,
		MaxAge:      10 * time.Second,
		Storage:     jetstream.MemoryStorage,
	}

	if _, err := c.js.CreateOrUpdateStream(ctxStartup, cfg); err != nil {
		return err
	}

	consumer, err := c.js.CreateConsumer(ctxStartup, cfg.Name, jetstream.ConsumerConfig{
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

	c.jsCons, err = consumer.Consume(c.handleMessageFunc(ctxCancel),
		c.js.ConsumeErrHandlerWithRestart(ctxCancel, c.logger,
			c.registerSubscriptions,
		))
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) handleMessageFunc(ctx context.Context) jetstream.MessageHandler {
	return func(msg jetstream.Msg) {
		remoteCtx, _ := events.GetJetstreamMsgContext(msg)
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
