package perms

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

func (p *Perms) registerEvents(ctx context.Context) error {
	cfg := &nats.StreamConfig{
		Name:        "PERMS",
		Description: natsutils.Description,
		Retention:   nats.InterestPolicy,
		Subjects:    []string{fmt.Sprintf("%s.>", BaseSubject)},
		Discard:     nats.DiscardOld,
		MaxAge:      15 * time.Second,
		Storage:     nats.MemoryStorage,
	}

	if _, err := natsutils.CreateOrUpdateStream(ctx, p.js, cfg); err != nil {
		return err
	}

	sub, err := p.js.Subscribe(fmt.Sprintf("%s.>", BaseSubject), p.handleMessage, nats.DeliverNew())
	if err != nil {
		return err
	}
	p.jsSub = sub

	return nil
}

func (p *Perms) handleMessage(msg *nats.Msg) {
	p.logger.Debug("received event message", zap.String("subject", msg.Subject))

	switch events.Type(strings.TrimPrefix(msg.Subject, string(BaseSubject)+".")) {
	case RolePermUpdateSubject:
		event := &RolePermUpdateEvent{}
		if err := json.Unmarshal(msg.Data, event); err != nil {
			p.logger.Error("failed to unmarshal message event data", zap.Error(err))
			return
		}

		if err := p.loadRolePermissions(p.ctx, event.RoleID); err != nil {
			p.logger.Error("failed to update role permissions", zap.Error(err))
			return
		}

	case RoleAttrUpdateSubject:
		event := &RoleAttrUpdateEvent{}
		if err := json.Unmarshal(msg.Data, event); err != nil {
			p.logger.Error("failed to unmarshal message event data", zap.Error(err))
			return
		}

		if err := p.loadRoleAttributes(p.ctx, event.RoleID); err != nil {
			p.logger.Error("failed to update role permissions", zap.Error(err))
			return
		}

	case JobAttrUpdateSubject:
		event := &JobAttrUpdateEvent{}
		if err := json.Unmarshal(msg.Data, event); err != nil {
			p.logger.Error("failed to unmarshal message event data", zap.Error(err))
			return
		}

		if err := p.loadJobAttrs(p.ctx, event.Job); err != nil {
			p.logger.Error("failed to update job attributes", zap.Error(err))
			return
		}

	default:
		p.logger.Error("unknown type of perms events message received")
	}
}

func (p *Perms) publishMessage(subj events.Type, data any) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if _, err := p.js.Publish(fmt.Sprintf("%s.%s", BaseSubject, subj), out); err != nil {
		return err
	}

	return nil
}
