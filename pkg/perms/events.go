package perms

import (
	"context"
	"fmt"
	"strings"

	permissionsevents "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/events"
	"github.com/fivenet-app/fivenet/v2026/pkg/events"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/protoutils"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

const (
	BaseSubject events.Subject = "perms"

	RoleCreatedSubject      events.Type = "roleperm.create"
	RolePermUpdateSubject   events.Type = "roleperm.update"
	RoleDeletedSubject      events.Type = "roleperm.delete"
	RoleAttrUpdateSubject   events.Type = "roleattr.update"
	JobLimitsUpdatedSubject events.Type = "joblimits.update"
)

func (ps *Perms) registerSubscriptions(ctxCancel context.Context) error {
	if ps.ncSub != nil {
		if err := ps.ncSub.Unsubscribe(); err != nil {
			ps.logger.Error("failed to unsubscribe from previous perms subject", zap.Error(err))
		}
		ps.ncSub = nil
	}

	ncSub, err := ps.nc.Subscribe(fmt.Sprintf("%s.>", BaseSubject), ps.handleMessageFunc(ctxCancel))
	if err != nil {
		return fmt.Errorf("failed to subscribe to events. %w", err)
	}
	ps.ncSub = ncSub

	return nil
}

func (ps *Perms) handleMessageFunc(ctx context.Context) nats.MsgHandler {
	return func(msg *nats.Msg) {
		ps.logger.Debug("received event message", zap.String("subject", msg.Subject))

		switch events.Type(strings.TrimPrefix(msg.Subject, string(BaseSubject)+".")) {
		case RoleCreatedSubject:
			fallthrough
		case RolePermUpdateSubject:
			event := &permissionsevents.RoleIDEvent{}
			if err := protojson.Unmarshal(msg.Data, event); err != nil {
				ps.logger.Error("failed to unmarshal message event data", zap.Error(err))
				return
			}

			if err := ps.loadRoles(ctx, event.GetRoleId()); err != nil {
				ps.logger.Error("failed to load role for role data load", zap.Error(err))
				return
			}

			if err := ps.loadRolePermissions(ctx, event.GetRoleId()); err != nil {
				ps.logger.Error("failed to load updated role permissions", zap.Error(err))
				return
			}

		case RoleAttrUpdateSubject:
			event := &permissionsevents.RoleIDEvent{}
			if err := protojson.Unmarshal(msg.Data, event); err != nil {
				ps.logger.Error("failed to unmarshal message event data", zap.Error(err))
				return
			}

			if err := ps.loadRoles(ctx, event.GetRoleId()); err != nil {
				ps.logger.Error("failed to load role for role data load", zap.Error(err))
				return
			}

			if err := ps.loadRoleAttributes(ctx, event.GetRoleId()); err != nil {
				ps.logger.Error("failed to load updated role permissions", zap.Error(err))
				return
			}

		case RoleDeletedSubject:
			event := &permissionsevents.RoleIDEvent{}
			if err := protojson.Unmarshal(msg.Data, event); err != nil {
				ps.logger.Error("failed to unmarshal message event data", zap.Error(err))
				return
			}

			// Remove role from local state
			ps.deleteRole(event.GetRoleId(), event.GetJob(), event.GetGrade())

		case JobLimitsUpdatedSubject:
			event := &permissionsevents.JobLimitsUpdatedEvent{}
			if err := protojson.Unmarshal(msg.Data, event); err != nil {
				ps.logger.Error("failed to unmarshal message event data", zap.Error(err))
				return
			}

			if err := ps.loadJobRoles(ctx, event.GetJob()); err != nil {
				ps.logger.Error(
					"failed to load job role permissions and attributes",
					zap.Error(err),
				)
				return
			}

		default:
			ps.logger.Error("unknown type of perms events message received")
		}
	}
}

func (ps *Perms) publishMessage(_ context.Context, subj events.Type, msg proto.Message) error {
	out, err := protoutils.MarshalToJSON(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal data. %w", err)
	}

	if err := ps.nc.Publish(string(BaseSubject)+"."+string(subj), out); err != nil {
		return fmt.Errorf(
			"failed to publish message to subject %s. %w",
			string(BaseSubject)+"."+string(subj),
			err,
		)
	}

	return nil
}
