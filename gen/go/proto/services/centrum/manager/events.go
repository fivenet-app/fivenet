package manager

import (
	"context"
	"fmt"

	"github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	eventscentrum "github.com/galexrt/fivenet/gen/go/proto/services/centrum/events"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (s *Manager) watchUserChanges() {
	userCh := s.tracker.Subscribe()

	for {
		select {
		case <-s.ctx.Done():
			return

		case event := <-userCh:
			func() {
				ctx, span := s.tracer.Start(s.ctx, "centrum-watch-users")
				defer span.End()

				for _, userInfo := range event.Added {
					if _, ok := s.UserIDToUnitID.Load(userInfo.UserID); !ok {
						unitId, err := s.LoadUnitIDForUserID(ctx, userInfo.UserID)
						if err != nil {
							s.logger.Error("failed to load user unit id", zap.Error(err))
							continue
						}
						if unitId == 0 {
							continue
						}

						s.UserIDToUnitID.Store(userInfo.UserID, unitId)
					}
				}

				for _, userInfo := range event.Removed {
					s.handleRemoveUserFromDisponents(ctx, userInfo.Job, userInfo.UserID)
					s.handleRemoveUserFromUnit(ctx, userInfo.Job, userInfo.UserID)
				}
			}()
		}
	}
}

func (s *Manager) handleRemoveUserFromDisponents(ctx context.Context, job string, userId int32) {
	if s.CheckIfUserIsDisponent(job, userId) {
		if err := s.DisponentSignOn(ctx, job, userId, false); err != nil {
			s.logger.Error("failed to remove user from disponents", zap.Error(err))
			return
		}
	}
}

func (s *Manager) handleRemoveUserFromUnit(ctx context.Context, job string, userId int32) bool {
	unitId, ok := s.UserIDToUnitID.Load(userId)
	if !ok {
		// Nothing to do
		return false
	}

	unit, ok := s.GetUnit(job, unitId)
	if !ok {
		s.UserIDToUnitID.Delete(userId)
		return false
	}

	if err := s.UpdateUnitAssignments(ctx, job, &userId, unit, nil, []int32{userId}); err != nil {
		s.logger.Error("failed to remove user from unit", zap.Error(err))
		return false
	}

	return true
}

func (s *Manager) watchEvents() error {
	msgCh := make(chan *nats.Msg, 256)

	sub, err := s.events.JS.ChanSubscribe(fmt.Sprintf("%s.>", eventscentrum.BaseSubject), msgCh, nats.DeliverLastPerSubject())
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	// Watch for events from message queue
	for {
		select {
		case <-s.ctx.Done():
			return nil

		case msg := <-msgCh:
			msg.Ack()

			job, topic, tType := eventscentrum.SplitSubject(msg.Subject)

			switch topic {
			case eventscentrum.TopicGeneral:
				switch tType {
				case eventscentrum.TypeGeneralDisponents:
					var dest dispatch.DisponentsChange
					if err := proto.Unmarshal(msg.Data, &dest); err != nil {
						s.logger.Error("failed to unmarshal disponents message", zap.Error(err))
						continue
					}

					s.Disponents.Store(job, dest.Disponents)

				case eventscentrum.TypeGeneralSettings:
					var dest dispatch.Settings
					if err := proto.Unmarshal(msg.Data, &dest); err != nil {
						s.logger.Error("failed to unmarshal settings message", zap.Error(err))
						continue
					}

					settings, ok := s.Settings.LoadOrStore(job, &dest)
					if ok {
						settings.Enabled = dest.Enabled
						settings.Mode = dest.Mode
						settings.FallbackMode = dest.FallbackMode
					}
				}

			case eventscentrum.TopicDispatch:
				switch tType {
				case eventscentrum.TypeDispatchCreated:
					fallthrough
				case eventscentrum.TypeDispatchStatus:
					fallthrough
				case eventscentrum.TypeDispatchUpdated:
					dest := &dispatch.Dispatch{}
					if err := proto.Unmarshal(msg.Data, dest); err != nil {
						s.logger.Error("failed to unmarshal dispatch message", zap.Error(err))
						continue
					}

					if dispatch, ok := s.GetDispatchesMap(job).LoadOrStore(dest.Id, dest); ok {
						dispatch.Merge(dest)
					}

				case eventscentrum.TypeDispatchDeleted:
					dest := &dispatch.Dispatch{}
					if err := proto.Unmarshal(msg.Data, dest); err != nil {
						s.logger.Error("failed to unmarshal dispatch message", zap.Error(err))
						continue
					}

					s.GetDispatchesMap(job).Delete(dest.Id)

				}

			case eventscentrum.TopicUnit:
				switch tType {
				case eventscentrum.TypeUnitUpdated:
					dest := &dispatch.Unit{}
					if err := proto.Unmarshal(msg.Data, dest); err != nil {
						s.logger.Error("failed to unmarshal unit message", zap.Error(err))
						continue
					}

					if unit, ok := s.GetUnitsMap(job).LoadOrStore(dest.Id, dest); ok {
						unit.Merge(dest)
					}

				case eventscentrum.TypeUnitStatus:
					dest := &dispatch.Unit{}
					if err := proto.Unmarshal(msg.Data, dest); err != nil {
						s.logger.Error("failed to unmarshal unit message", zap.Error(err))
						continue
					}

					if unit, ok := s.GetUnitsMap(job).LoadOrStore(dest.Id, dest); ok {
						unit.Merge(dest)
					}

					if dest.Status.Status == dispatch.StatusUnit_STATUS_UNIT_USER_ADDED {
						s.UserIDToUnitID.Store(*dest.Status.UserId, dest.Status.UnitId)
					} else if dest.Status.Status == dispatch.StatusUnit_STATUS_UNIT_USER_REMOVED {
						s.UserIDToUnitID.Delete(*dest.Status.UserId)
					}

				case eventscentrum.TypeUnitDeleted:
					dest := &dispatch.Unit{}
					if err := proto.Unmarshal(msg.Data, dest); err != nil {
						s.logger.Error("failed to unmarshal unit message", zap.Error(err))
						continue
					}

					s.GetUnitsMap(dest.Job).Delete(dest.Id)

				}
			}
		}
	}
}
