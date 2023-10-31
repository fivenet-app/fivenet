package manager

import (
	"context"
	"fmt"

	"github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	eventscentrum "github.com/galexrt/fivenet/gen/go/proto/services/centrum/events"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
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
					s.handleRemovedUserFromDisponents(ctx, userInfo.Job, userInfo.UserID)
					s.handleRemovedUserFromUnit(ctx, userInfo.Job, userInfo.UserID)
				}
			}()
		}
	}
}

func (s *Manager) handleRemovedUserFromDisponents(ctx context.Context, job string, userId int32) {
	if s.CheckIfUserIsDisponent(job, userId) {
		if err := s.DisponentSignOn(ctx, job, userId, false); err != nil {
			s.logger.Error("failed to remove user from disponents", zap.Error(err))
			return
		}
	}
}

func (s *Manager) handleRemovedUserFromUnit(ctx context.Context, job string, userId int32) bool {
	unitId, ok := s.UserIDToUnitID.Load(userId)
	if !ok {
		// Nothing to do
		return false
	}

	s.UserIDToUnitID.Delete(userId)

	unit, ok := s.GetUnit(job, unitId)
	if !ok {
		return false
	}

	if err := s.UpdateUnitAssignments(ctx, &userinfo.UserInfo{
		UserId: userId,
		Job:    job,
	}, unit, nil, []int32{userId}); err != nil {
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
		func() error {
			ctx, span := s.tracer.Start(s.ctx, "centrum-state-events")
			defer span.End()

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
							return err
						}

						s.Disponents.Store(job, dest.Disponents)

					case eventscentrum.TypeGeneralSettings:
						var dest dispatch.Settings
						if err := proto.Unmarshal(msg.Data, &dest); err != nil {
							return err
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
						dest := &dispatch.Dispatch{}
						if err := proto.Unmarshal(msg.Data, dest); err != nil {
							return err
						}

						s.GetDispatchesMap(job).Store(dest.Id, dest)

					case eventscentrum.TypeDispatchUpdated:
						dest := &dispatch.Dispatch{}
						if err := proto.Unmarshal(msg.Data, dest); err != nil {
							return err
						}

						s.GetDispatchesMap(job).Store(dest.Id, dest)

					case eventscentrum.TypeDispatchDeleted:
						dest := &dispatch.Dispatch{}
						if err := proto.Unmarshal(msg.Data, dest); err != nil {
							return err
						}

						s.GetDispatchesMap(job).Delete(dest.Id)

					}

				case eventscentrum.TopicUnit:
					switch tType {
					case eventscentrum.TypeUnitDeleted:
						dest := &dispatch.Unit{}
						if err := proto.Unmarshal(msg.Data, dest); err != nil {
							return err
						}

						units, ok := s.Units.Load(dest.Job)
						if ok {
							units.Delete(dest.Id)
						}

					case eventscentrum.TypeUnitUpdated:
						dest := &dispatch.Unit{}
						if err := proto.Unmarshal(msg.Data, dest); err != nil {
							return err
						}

						s.GetUnitsMap(job).Store(dest.Id, dest)

					case eventscentrum.TypeUnitStatus:
						dest := &dispatch.Unit{}
						if err := proto.Unmarshal(msg.Data, dest); err != nil {
							return err
						}

						if dest.Status.Status == dispatch.StatusUnit_STATUS_UNIT_USER_ADDED {
							s.UserIDToUnitID.Store(*dest.Status.UserId, dest.Status.UnitId)
						} else if dest.Status.Status == dispatch.StatusUnit_STATUS_UNIT_USER_REMOVED {
							s.UserIDToUnitID.Delete(*dest.Status.UserId)
						}

						unit, ok := s.GetUnitsMap(job).Load(dest.Id)
						if ok {
							if dest.Status.Status != dispatch.StatusUnit_STATUS_UNIT_USER_ADDED && dest.Status.Status != dispatch.StatusUnit_STATUS_UNIT_USER_REMOVED {
								unit.Status = dest.Status
							}
						} else {
							// "Cache/State miss" load from database
							if err := s.LoadUnits(ctx, dest.Id); err != nil {
								s.logger.Error("failed to load unit", zap.Error(err))
							}
						}
					}
				}
			}

			return nil
		}()
	}
}
