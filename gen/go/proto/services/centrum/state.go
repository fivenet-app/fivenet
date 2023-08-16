package centrum

import (
	"fmt"
	"time"

	dispatch "github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (s *Server) watchForEvents() error {
	msgCh := make(chan *nats.Msg, 256)

	sub, err := s.events.JS.ChanSubscribe(fmt.Sprintf("%s.>", BaseSubject), msgCh)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	// Watch for events from message queue
	for {
		func() error {
			ctx, span := s.tracer.Start(s.ctx, "centrum-initial-cache")
			defer span.End()

			select {
			case <-s.ctx.Done():
				return nil
			case msg := <-msgCh:
				msg.Ack()
				job, topic, tType := s.splitSubject(msg.Subject)

				switch topic {
				case TopicGeneral:
					switch tType {
					case TypeGeneralDisponents:
						var dest DisponentsChange
						if err := proto.Unmarshal(msg.Data, &dest); err != nil {
							return err
						}

						s.disponents.Store(job, dest.Disponents)

					case TypeGeneralSettings:
						var dest dispatch.Settings
						if err := proto.Unmarshal(msg.Data, &dest); err != nil {
							return err
						}

						s.settings.Store(job, &dest)
					}

				case TopicDispatch:
					switch tType {
					case TypeDispatchCreated:
						var dest dispatch.Dispatch
						if err := proto.Unmarshal(msg.Data, &dest); err != nil {
							return err
						}

						s.getDispatchesMap(job).Store(dest.Id, &dest)

					case TypeDispatchDeleted:
						var dest dispatch.Dispatch
						if err := proto.Unmarshal(msg.Data, &dest); err != nil {
							return err
						}

						s.getDispatchesMap(job).Delete(dest.Id)

					case TypeDispatchUpdated:
						var dest dispatch.Dispatch
						if err := proto.Unmarshal(msg.Data, &dest); err != nil {
							return err
						}

						s.getDispatchesMap(job).Store(dest.Id, &dest)

					case TypeDispatchStatus:
						var dest dispatch.DispatchStatus
						if err := proto.Unmarshal(msg.Data, &dest); err != nil {
							return err
						}

						dsp, ok := s.getDispatchesMap(job).Load(dest.Id)
						if ok {
							dsp.Status = &dest
						} else {
							// "Cache/State miss" load from database
							s.loadDispatches(ctx, dest.DispatchId)
						}

					}

				case TopicUnit:
					switch tType {
					case TypeUnitUserAssigned:
						var dest dispatch.Unit
						if err := proto.Unmarshal(msg.Data, &dest); err != nil {
							return err
						}

						unit, ok := s.getUnitsMap(job).Load(dest.Id)
						if ok {
							unit.Users = dest.Users
						} else {
							s.loadUnits(ctx, dest.Id)
						}

					case TypeUnitCreated:
						var dest dispatch.Unit
						if err := proto.Unmarshal(msg.Data, &dest); err != nil {
							return err
						}

						s.getUnitsMap(job).Store(dest.Id, &dest)

					case TypeUnitDeleted:
						var dest dispatch.Unit
						if err := proto.Unmarshal(msg.Data, &dest); err != nil {
							return err
						}

						units, ok := s.units.Load(dest.Job)
						if ok {
							units.Delete(dest.Id)
						}

					case TypeUnitUpdated:
						var dest dispatch.Unit
						if err := proto.Unmarshal(msg.Data, &dest); err != nil {
							return err
						}

						s.getUnitsMap(job).Store(dest.Id, &dest)

					case TypeUnitStatus:
						var dest dispatch.UnitStatus
						if err := proto.Unmarshal(msg.Data, &dest); err != nil {
							return err
						}

						unit, ok := s.getUnitsMap(job).Load(dest.Id)
						if ok {
							unit.Status = &dest
						} else {
							// "Cache/State miss" load from database
							s.loadUnits(ctx, dest.UnitId)
						}
					}
				}
			}

			return nil
		}()
	}
}

func (s *Server) watchForUserChanges() {
	userCh := s.tracker.Subscribe()

	for {
		select {
		case <-s.ctx.Done():
			return
		case event := <-userCh:
			func() {
				ctx, span := s.tracer.Start(s.ctx, "centrum-watch-users")
				defer span.End()

				for _, userId := range event.Added {
					unitId, err := s.loadUnitIDForUserID(ctx, userId)
					if err != nil {
						s.logger.Error("failed to load user unit id", zap.Error(err))
						continue
					}

					if unitId == 0 {
						continue
					}

					s.userIDToUnitID.Store(userId, unitId)
				}

				for _, userId := range event.Removed {
					unitId, ok := s.userIDToUnitID.Load(userId)
					if !ok {
						// Nothing to do
						continue
					}

					if err := s.updateUnitStatus(ctx, &userinfo.UserInfo{}, &dispatch.UnitStatus{
						UnitId:    unitId,
						Status:    dispatch.UNIT_STATUS_USER_REMOVED,
						UserId:    &userId,
						CreatorId: &userId,
					}); err != nil {
						s.logger.Error("failed to update user's unit status", zap.Error(err))
						continue
					}

					s.userIDToUnitID.Delete(userId)
				}
			}()
		}
	}
}

func (s *Server) housekeeper() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-time.After(1 * time.Second):
			// TODO take care of housekeeper tasks such as:
			// * Remove empty units from dispatches (if no other unit is assigned to dispatch update status to UNASSIGNED)
			// * Set dispatches to status ARCHIVED when older than 60 minutes
		}
	}
}
