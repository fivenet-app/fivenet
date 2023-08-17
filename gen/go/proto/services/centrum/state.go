package centrum

import (
	"context"
	"fmt"
	"time"

	dispatch "github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	jet "github.com/go-jet/jet/v2/mysql"
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
						fallthrough
					case TypeDispatchStatus:
						fallthrough
					case TypeDispatchUpdated:
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

					}

				case TopicUnit:
					switch tType {
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

						if dest.Status == dispatch.UNIT_STATUS_USER_ADDED {
							s.userIDToUnitID.Store(*dest.UserId, dest.UnitId)
						} else if dest.Status == dispatch.UNIT_STATUS_USER_REMOVED {
							s.userIDToUnitID.Delete(*dest.UserId)
						}

						unit, ok := s.getUnitsMap(job).Load(dest.UnitId)
						if ok {
							unit.Status = &dest

						} else {
							// "Cache/State miss" load from database
							if err := s.loadUnits(ctx, dest.UnitId); err != nil {
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

					user, err := s.resolveUserById(ctx, userId)
					if err != nil {
						s.logger.Error("failed to get user info from db", zap.Error(err))
						continue
					}

					unit, ok := s.getUnit(user.Job, unitId)
					if !ok {
						continue
					}

					if err := s.updateUnitAssignments(ctx, &userinfo.UserInfo{
						UserId: userId,
						Job:    user.Job,
					}, unit, nil, []int32{userId}); err != nil {
						s.logger.Error("failed to remove user from unit", zap.Error(err))
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
			func() {
				ctx, span := s.tracer.Start(s.ctx, "centrum-housekeeper")
				defer span.End()

				if err := s.handleDispatchAssignmentExpiration(ctx); err != nil {
					s.logger.Error("failed to handle expired dispatch assignments", zap.Error(err))
				}

				if err := s.archiveDispatches(ctx); err != nil {
					s.logger.Error("failed to archive dispatches", zap.Error(err))
				}

				if err := s.removeDispatchesFromEmptyUnits(ctx); err != nil {
					s.logger.Error("failed to clean empty units from dispatches", zap.Error(err))
				}
			}()
		}
	}
}

// TODO handle expired `createdAt` dispatch unit assignments
func (s *Server) handleDispatchAssignmentExpiration(ctx context.Context) error {
	stmt := tDispatchUnit.
		SELECT(
			tDispatchUnit.DispatchID.AS("dispatch_id"),
			tDispatchUnit.UnitID.AS("unit_id"),
			tUnits.Job.AS("job"),
		).
		FROM(
			tDispatchUnit.
				INNER_JOIN(tUnits,
					tUnits.ID.EQ(tDispatchUnit.UnitID),
				),
		).
		WHERE(jet.AND(
			tDispatchUnit.ExpiresAt.LT(jet.NOW()),
		))

	var dest []*struct {
		DispatchID uint64
		UnitID     uint64
		Job        string
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return err
	}

	for _, ua := range dest {
		dsp, ok := s.getDispatch(ua.Job, ua.DispatchID)
		if !ok {
			continue
		}

		if err := s.updateDispatchAssignments(ctx, ua.Job, nil, dsp, nil, []uint64{ua.UnitID}); err != nil {
			return err
		}
	}

	return nil
}

// Set `COMPLETED`/`CANCELLED` dispatches to status `ARCHIVED` when the status is older than 20 minutes
func (s *Server) archiveDispatches(ctx context.Context) error {
	stmt := tDispatchStatus.
		SELECT(
			tDispatchStatus.DispatchID.AS("dispatch_id"),
			tDispatch.Job.AS("job"),
		).
		FROM(
			tDispatchStatus.
				INNER_JOIN(tDispatch,
					tDispatch.ID.EQ(tDispatchStatus.DispatchID),
				),
		).
		WHERE(jet.AND(
			tDispatchStatus.CreatedAt.LT_EQ(
				jet.CURRENT_TIMESTAMP().SUB(jet.INTERVAL(20, jet.MINUTE)),
			),
			tDispatchStatus.Status.IN(
				jet.Int16(int16(dispatch.DISPATCH_STATUS_COMPLETED)),
				jet.Int16(int16(dispatch.DISPATCH_STATUS_CANCELLED)),
			),
			tDispatchStatus.ID.IS_NULL().OR(
				tDispatchStatus.ID.EQ(
					jet.RawInt("SELECT MAX(`dispatchstatus`.`id`) FROM `fivenet_centrum_dispatches_status` AS `dispatchstatus` WHERE `dispatchstatus`.`dispatch_id` = `dispatch`.`id`"),
				),
			),
		))

	var dest []*struct {
		DispatchID uint64
		Job        string
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return err
	}

	for _, ds := range dest {
		dsp, ok := s.getDispatch(ds.Job, ds.DispatchID)
		if !ok {
			continue
		}

		if err := s.updateDispatchStatus(ctx, ds.Job, dsp, &dispatch.DispatchStatus{
			DispatchId: dsp.Id,
			Status:     dispatch.DISPATCH_STATUS_ARCHIVED,
			UserId:     dsp.Status.UserId,
		}); err != nil {
			return err
		}

		s.getDispatchesMap(ds.Job).Delete(ds.DispatchID)
	}

	return nil
}

// Remove empty units from dispatches (if no other unit is assigned to dispatch update status to UNASSIGNED) by
// iterating over the dispatches and making sure the assigned units aren't empty
func (s *Server) removeDispatchesFromEmptyUnits(ctx context.Context) error {

	// TODO

	return nil
}
