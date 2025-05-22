package perms

import (
	"context"
	"fmt"
	"strings"

	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/nats-io/nats.go"
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

func (p *Perms) registerSubscriptions(ctxCancel context.Context) error {
	if p.ncSub != nil {
		p.ncSub.Unsubscribe()
		p.ncSub = nil
	}

	ncSub, err := p.nc.Subscribe(fmt.Sprintf("%s.>", BaseSubject), p.handleMessageFunc(ctxCancel))
	if err != nil {
		return fmt.Errorf("failed to subscribe to events. %w", err)
	}
	p.ncSub = ncSub

	return nil
}

func (p *Perms) handleMessageFunc(ctx context.Context) nats.MsgHandler {
	return func(msg *nats.Msg) {
		p.logger.Debug("received event message", zap.String("subject", msg.Subject))

		switch events.Type(strings.TrimPrefix(msg.Subject, string(BaseSubject)+".")) {
		case RoleCreatedSubject:
			fallthrough
		case RolePermUpdateSubject:
			event := &RoleIDEvent{}
			if err := json.Unmarshal(msg.Data, event); err != nil {
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
			if err := json.Unmarshal(msg.Data, event); err != nil {
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
			if err := json.Unmarshal(msg.Data, event); err != nil {
				p.logger.Error("failed to unmarshal message event data", zap.Error(err))
				return
			}

			// Remove role from local state
			p.deleteRole(event.RoleID, event.Job, event.Grade)

		case JobAttrUpdateSubject:
			event := &JobAttrUpdateEvent{}
			if err := json.Unmarshal(msg.Data, event); err != nil {
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
		return fmt.Errorf("failed to marshal data. %w", err)
	}

	if err := p.nc.Publish(string(BaseSubject)+"."+string(subj), out); err != nil {
		return fmt.Errorf("failed to publish message to subject %s. %w", string(BaseSubject)+"."+string(subj), err)
	}

	return nil
}
