package perms

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/galexrt/fivenet/pkg/events"
	natsutils "github.com/galexrt/fivenet/pkg/nats"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

const (
	BaseSubject events.Subject = "perms"

	RolePermUpdateSubject events.Type = "roleperm.update"
	RoleAttrUpdateSubject events.Type = "roleattr.update"
	JobAttrUpdateSubject  events.Type = "jobattr.update"
)

type RolePermUpdateEvent struct {
	RoleID uint64
}

type RoleAttrUpdateEvent struct {
	RoleID uint64
}

type JobAttrUpdateEvent struct {
	Job string
}

func (p *Perms) registerSubscriptions(ctx context.Context, c context.Context) error {
	cfg := jetstream.StreamConfig{
		Name:        "PERMS",
		Description: natsutils.Description,
		Retention:   jetstream.InterestPolicy,
		Subjects:    []string{fmt.Sprintf("%s.>", BaseSubject)},
		Discard:     jetstream.DiscardOld,
		MaxAge:      15 * time.Second,
		Storage:     jetstream.MemoryStorage,
	}

	if _, err := p.js.CreateOrUpdateStream(ctx, cfg); err != nil {
		return err
	}

	consumer, err := p.js.CreateConsumer(ctx, cfg.Name, jetstream.ConsumerConfig{
		DeliverPolicy: jetstream.DeliverNewPolicy,
		FilterSubject: fmt.Sprintf("%s.>", BaseSubject),
	})
	if err != nil {
		return err
	}

	cons, err := consumer.Consume(p.handleMessageFunc(c))
	if err != nil {
		return err
	}
	p.jsCons = cons

	return nil
}

func (p *Perms) handleMessageFunc(ctx context.Context) jetstream.MessageHandler {
	return func(msg jetstream.Msg) {
		p.logger.Debug("received event message", zap.String("subject", msg.Subject()))

		switch events.Type(strings.TrimPrefix(msg.Subject(), string(BaseSubject)+".")) {
		case RolePermUpdateSubject:
			event := &RolePermUpdateEvent{}
			if err := json.Unmarshal(msg.Data(), event); err != nil {
				p.logger.Error("failed to unmarshal message event data", zap.Error(err))
				return
			}

			if err := p.loadRolePermissions(ctx, event.RoleID); err != nil {
				p.logger.Error("failed to update role permissions", zap.Error(err))
				return
			}

		case RoleAttrUpdateSubject:
			event := &RoleAttrUpdateEvent{}
			if err := json.Unmarshal(msg.Data(), event); err != nil {
				p.logger.Error("failed to unmarshal message event data", zap.Error(err))
				return
			}

			if err := p.loadRoleAttributes(ctx, event.RoleID); err != nil {
				p.logger.Error("failed to update role permissions", zap.Error(err))
				return
			}

		case JobAttrUpdateSubject:
			event := &JobAttrUpdateEvent{}
			if err := json.Unmarshal(msg.Data(), event); err != nil {
				p.logger.Error("failed to unmarshal message event data", zap.Error(err))
				return
			}

			if err := p.loadJobAttrs(ctx, event.Job); err != nil {
				p.logger.Error("failed to update job attributes", zap.Error(err))
				return
			}

		default:
			p.logger.Error("unknown type of perms events message received")
		}
	}
}

func (p *Perms) publishMessage(ctx context.Context, subj events.Type, data any) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if _, err := p.js.Publish(ctx, fmt.Sprintf("%s.%s", BaseSubject, subj), out); err != nil {
		return err
	}

	return nil
}
