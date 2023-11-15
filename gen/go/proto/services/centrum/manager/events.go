package manager

import (
	"fmt"

	"github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	eventscentrum "github.com/galexrt/fivenet/gen/go/proto/services/centrum/events"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (s *Manager) registerSubscriptions() error {
	if _, err := s.events.JS.Subscribe(fmt.Sprintf("%s.*.%s.>", eventscentrum.BaseSubject, eventscentrum.TopicGeneral), s.watchTopicGeneral, nats.DeliverLastPerSubject()); err != nil {
		s.logger.Error("failed to subscribe to centrum general topic", zap.Error(err))
		return err
	}

	if _, err := s.events.JS.Subscribe(fmt.Sprintf("%s.*.%s.>", eventscentrum.BaseSubject, eventscentrum.TopicUnit), s.watchTopicUnits, nats.DeliverLastPerSubject()); err != nil {
		s.logger.Error("failed to subscribe to centrum units topic", zap.Error(err))
		return err
	}

	if _, err := s.events.JS.Subscribe(fmt.Sprintf("%s.*.%s.>", eventscentrum.BaseSubject, eventscentrum.TopicDispatch), s.watchTopicDispatches, nats.DeliverLastPerSubject()); err != nil {
		s.logger.Error("failed to subscribe to centrum dispatch topic", zap.Error(err))
		return err
	}

	return nil
}

func (s *Manager) watchTopicGeneral(msg *nats.Msg) {
	msg.Ack()

	job, _, tType := eventscentrum.SplitSubject(msg.Subject)

	switch tType {
	case eventscentrum.TypeGeneralDisponents:
		var dest dispatch.DisponentsChange
		if err := proto.Unmarshal(msg.Data, &dest); err != nil {
			s.logger.Error("failed to unmarshal disponents message", zap.Error(err))
			return
		}

		s.UpdateDisponents(job, dest.Disponents)

	case eventscentrum.TypeGeneralSettings:
		var dest dispatch.Settings
		if err := proto.Unmarshal(msg.Data, &dest); err != nil {
			s.logger.Error("failed to unmarshal settings message", zap.Error(err))
			return
		}

		if err := s.State.UpdateSettings(job, &dest); err != nil {
			s.logger.Error("failed to update settings", zap.Error(err))
			return
		}
	}
}

func (s *Manager) watchTopicUnits(msg *nats.Msg) {
	msg.Ack()

	job, _, tType := eventscentrum.SplitSubject(msg.Subject)

	switch tType {
	case eventscentrum.TypeUnitCreated:
		fallthrough
	case eventscentrum.TypeUnitStatus:
		fallthrough
	case eventscentrum.TypeUnitUpdated:
		dest := &dispatch.Unit{}
		if err := proto.Unmarshal(msg.Data, dest); err != nil {
			s.logger.Error("failed to unmarshal unit message", zap.Error(err))
			return
		}

		if err := s.State.UpdateUnit(job, dest.Id, dest); err != nil {
			s.logger.Error("failed to update unit", zap.Error(err))
			return
		}

		if tType == eventscentrum.TypeUnitStatus {
			if dest.Status.Status == dispatch.StatusUnit_STATUS_UNIT_USER_ADDED {
				s.SetUnitForUser(*dest.Status.UserId, dest.Status.UnitId)
			} else if dest.Status.Status == dispatch.StatusUnit_STATUS_UNIT_USER_REMOVED {
				s.UnsetUnitIDForUser(*dest.Status.UserId)
			}
		}

	case eventscentrum.TypeUnitDeleted:
		dest := &dispatch.Unit{}
		if err := proto.Unmarshal(msg.Data, dest); err != nil {
			s.logger.Error("failed to unmarshal unit message", zap.Error(err))
			return
		}

		s.State.DeleteUnit(job, dest.Id)
	}
}

func (s *Manager) watchTopicDispatches(msg *nats.Msg) {
	msg.Ack()

	job, _, tType := eventscentrum.SplitSubject(msg.Subject)

	switch tType {
	case eventscentrum.TypeDispatchCreated:
		fallthrough
	case eventscentrum.TypeDispatchStatus:
		fallthrough
	case eventscentrum.TypeDispatchUpdated:
		dest := &dispatch.Dispatch{}
		if err := proto.Unmarshal(msg.Data, dest); err != nil {
			s.logger.Error("failed to unmarshal dispatch message", zap.Error(err))
			return
		}

		if err := s.State.UpdateDispatch(job, dest.Id, dest); err != nil {
			s.logger.Error("failed to update dispatch", zap.Error(err))
			return
		}

	case eventscentrum.TypeDispatchDeleted:
		dest := &dispatch.Dispatch{}
		if err := proto.Unmarshal(msg.Data, dest); err != nil {
			s.logger.Error("failed to unmarshal dispatch message", zap.Error(err))
			return
		}

		s.State.DeleteDispatch(job, dest.Id)
	}
}
