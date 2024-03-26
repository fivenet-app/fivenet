package appconfig

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/galexrt/fivenet/pkg/events"
	"github.com/galexrt/fivenet/pkg/nats"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

const (
	BaseSubject events.Subject = "appconfig"

	UpdateSubject events.Type = "update"
)

func (c *Config) registerSubscriptions(ctx context.Context) error {
	cfg := jetstream.StreamConfig{
		Name:        "APPCONFIG",
		Description: "AppConfig update events",
		Retention:   jetstream.InterestPolicy,
		Subjects:    []string{fmt.Sprintf("%s.>", BaseSubject)},
		Discard:     jetstream.DiscardOld,
		MaxAge:      10 * time.Second,
		Storage:     jetstream.MemoryStorage,
		Replicas:    2,
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

	cons, err := consumer.Consume(c.handleMessageFunc(ctx), nats.ConsumeErrHandler(c.logger))
	if err != nil {
		return err
	}
	c.jsCons = cons

	return nil
}

func (c *Config) handleMessageFunc(ctx context.Context) jetstream.MessageHandler {
	return func(msg jetstream.Msg) {
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
