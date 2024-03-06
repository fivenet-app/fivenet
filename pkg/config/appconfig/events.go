package appconfig

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/galexrt/fivenet/pkg/events"
	natsutils "github.com/galexrt/fivenet/pkg/nats"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

const (
	BaseSubject events.Subject = "appconfig"

	UpdateSubject events.Type = "update"
)

func (c *Config) registerEvents(ctx context.Context) error {
	cfg := &nats.StreamConfig{
		Name:        "APPCONFIG",
		Description: natsutils.Description,
		Retention:   nats.InterestPolicy,
		Subjects:    []string{fmt.Sprintf("%s.>", BaseSubject)},
		Discard:     nats.DiscardOld,
		MaxAge:      10 * time.Second,
		Storage:     nats.MemoryStorage,
	}

	if _, err := natsutils.CreateOrUpdateStream(ctx, c.js, cfg); err != nil {
		return err
	}

	sub, err := c.js.Subscribe(fmt.Sprintf("%s.>", BaseSubject), c.handleMessage, nats.DeliverNew())
	if err != nil {
		return err
	}
	c.jsSub = sub

	return nil
}

func (c *Config) handleMessage(msg *nats.Msg) {
	split := strings.Split(msg.Subject, ".")

	if len(split) < 2 {
		c.logger.Warn("unknown app config subject received", zap.String("subject", msg.Subject))
		return
	}

	if split[1] == string(UpdateSubject) {
		if err := c.updateConfigFromDB(c.ctx); err != nil {
			c.logger.Error("failed to update app config from db", zap.Error(err))
		}
	}
}
