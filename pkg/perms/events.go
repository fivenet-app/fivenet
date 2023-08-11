package perms

import (
	"fmt"
	"strings"

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
	var err error
	p.eventSub, err = p.events.NC.Subscribe(fmt.Sprintf("%s.>", BaseSubject), p.handleMessage)
	if err != nil {
		return err
	}

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
		}

		if err := p.loadRolePermissions(p.ctx, event.RoleID); err != nil {
			p.logger.Error("failed to update role permissions", zap.Error(err))
		}

	case RoleAttrUpdateSubject:
		event := &RoleAttrUpdateEvent{}
		if err := json.Unmarshal(msg.Data, event); err != nil {
			p.logger.Error("failed to unmarshal message event data", zap.Error(err))
		}

		if err := p.loadRoleAttributes(p.ctx, event.RoleID); err != nil {
			p.logger.Error("failed to update role permissions", zap.Error(err))
		}

	default:
		p.logger.Error("unknown perms message received")
	}
}

func (p *Perms) publishMessage(subj events.Type, data any) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if err := p.events.NC.Publish(fmt.Sprintf("%s.%s", BaseSubject, subj), out); err != nil {
		return err
	}

	return nil
}
