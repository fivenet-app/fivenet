package perms

import (
	"fmt"
	"strings"
	"time"

	"github.com/galexrt/fivenet/pkg/events"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

const (
	BaseSubject events.Subject = "perms"

	RolePermUpdateSubject events.Type = "roleperm.update"
	RoleAttrUpdateSubject events.Type = "roleattr.update"
)

type RolePermUpdateEvent struct {
	RoleID uint64
}

type RoleAttrUpdateEvent struct {
	RoleID uint64
}

func (p *Perms) registerEvents() error {
	cfg := &nats.StreamConfig{
		Name:      "PERMS",
		Retention: nats.InterestPolicy,
		Subjects:  []string{fmt.Sprintf("%s.>", BaseSubject)},
		Discard:   nats.DiscardOld,
		MaxAge:    10 * time.Second,
	}

	if _, err := p.events.JS.CreateOrUpdateStream(cfg); err != nil {
		return err
	}

	sub, err := p.events.JS.Subscribe(fmt.Sprintf("%s.>", BaseSubject), p.handleMessage, nats.DeliverNew())
	if err != nil {
		return err
	}
	p.eventSub = sub

	return nil
}

func (p *Perms) handleMessage(msg *nats.Msg) {
	msg.Ack()
	p.logger.Debug("received message", zap.String("subject", msg.Subject))

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

	default:
		p.logger.Error("unknown type of perms events message received")
	}
}

func (p *Perms) publishMessage(subj events.Type, data any) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if _, err := p.events.JS.Publish(fmt.Sprintf("%s.%s", BaseSubject, subj), out); err != nil {
		return err
	}

	return nil
}
