package perms

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/pkg/events"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

const (
	BaseSubject events.Subject = "perms"

	RoleCreatedSubject    events.Type = "roleperm.create"
	RolePermUpdateSubject events.Type = "roleperm.update"
	RoleDeletedSubject    events.Type = "roleperm.delete"
	RoleAttrUpdateSubject events.Type = "roleattr.update"
	JobAttrUpdateSubject  events.Type = "jobattr.update"
)

type RoleIDEvent struct {
	RoleID uint64
	Job    string
	Grade  int32
}

type JobAttrUpdateEvent struct {
	Job string
}

func (p *Perms) registerSubscriptions(ctxStartup context.Context, ctxCancel context.Context) error {
	cfg := jetstream.StreamConfig{
		Name:        "PERMS",
		Description: "Perms system events",
		Retention:   jetstream.InterestPolicy,
		Subjects:    []string{fmt.Sprintf("%s.>", BaseSubject)},
		Discard:     jetstream.DiscardOld,
		MaxAge:      15 * time.Second,
		Storage:     jetstream.MemoryStorage,
	}

	if _, err := p.js.CreateOrUpdateStream(ctxStartup, cfg); err != nil {
		return err
	}

	consumer, err := p.js.CreateConsumer(ctxStartup, cfg.Name, jetstream.ConsumerConfig{
		DeliverPolicy: jetstream.DeliverNewPolicy,
		FilterSubject: fmt.Sprintf("%s.>", BaseSubject),
	})
	if err != nil {
		return err
	}

	if p.jsCons != nil {
		p.jsCons.Stop()
		p.jsCons = nil
	}

	p.jsCons, err = consumer.Consume(p.handleMessageFunc(ctxCancel),
		p.js.ConsumeErrHandlerWithRestart(ctxCancel, p.logger, p.registerSubscriptions))
	if err != nil {
		return err
	}

	return nil
}

func (p *Perms) handleMessageFunc(ctx context.Context) jetstream.MessageHandler {
	return func(msg jetstream.Msg) {
		if err := msg.Ack(); err != nil {
			p.logger.Error("failed to ack message", zap.Error(err))
		}

		p.logger.Debug("received event message", zap.String("subject", msg.Subject()))

		switch events.Type(strings.TrimPrefix(msg.Subject(), string(BaseSubject)+".")) {
		case RoleCreatedSubject:
			fallthrough
		case RolePermUpdateSubject:
			event := &RoleIDEvent{}
			if err := json.Unmarshal(msg.Data(), event); err != nil {
				p.logger.Error("failed to unmarshal message event data", zap.Error(err))
				return
			}

			if err := p.loadRoles(ctx, event.RoleID); err != nil {
				p.logger.Error("failed to load role for role data load", zap.Error(err))
				return
			}

			if err := p.loadRolePermissions(ctx, event.RoleID); err != nil {
				p.logger.Error("failed to load updated role permissions", zap.Error(err))
				return
			}

		case RoleAttrUpdateSubject:
			event := &RoleIDEvent{}
			if err := json.Unmarshal(msg.Data(), event); err != nil {
				p.logger.Error("failed to unmarshal message event data", zap.Error(err))
				return
			}

			if err := p.loadRoles(ctx, event.RoleID); err != nil {
				p.logger.Error("failed to load role for role data load", zap.Error(err))
				return
			}

			if err := p.loadRoleAttributes(ctx, event.RoleID); err != nil {
				p.logger.Error("failed to load updated role permissions", zap.Error(err))
				return
			}

		case RoleDeletedSubject:
			event := &RoleIDEvent{}
			if err := json.Unmarshal(msg.Data(), event); err != nil {
				p.logger.Error("failed to unmarshal message event data", zap.Error(err))
				return
			}

			// Remove role from local state
			p.deleteRole(event.RoleID, event.Job, event.Grade)

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

	if _, err := p.js.Publish(ctx, string(BaseSubject)+"."+string(subj), out); err != nil {
		return err
	}

	return nil
}
