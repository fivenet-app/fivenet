package appconfig

import (
	"context"
	"fmt"
	"strings"

	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

const (
	BaseSubject events.Subject = "appconfig"

	UpdateSubject events.Type = "update"
)

func (c *Config) registerSubscriptions(ctxCancel context.Context) error {
	if c.ncSub != nil {
		c.ncSub.Unsubscribe()
		c.ncSub = nil
	}

	ncSub, err := c.nc.Subscribe(fmt.Sprintf("%s.>", BaseSubject), c.handleMessageFunc(ctxCancel))
	if err != nil {
		return fmt.Errorf("failed to subscribe to events. %w", err)
	}
	c.ncSub = ncSub

	return nil
}

func (c *Config) handleMessageFunc(ctx context.Context) nats.MsgHandler {
	return func(msg *nats.Msg) {
		if err := msg.Ack(); err != nil {
			c.logger.Error("failed to ack message", zap.Error(err))
		}

		split := strings.Split(msg.Subject, ".")

		if len(split) < 2 {
			c.logger.Warn("unknown app config subject received", zap.String("subject", msg.Subject))
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
