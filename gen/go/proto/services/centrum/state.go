package centrum

import (
	"context"
	"fmt"
	"time"

	dispatch "github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	users "github.com/galexrt/fivenet/gen/go/proto/resources/users"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/nats-io/nats.go"
	"github.com/puzpuzpuz/xsync/v2"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (s *Server) watchForEvents() error {
	msgCh := make(chan *nats.Msg, 256)

	sub, err := s.events.JS.ChanSubscribe(fmt.Sprintf("%s.>", BaseSubject), msgCh, nats.DeliverLastPerSubject())
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
						var dest dispatch.Unit
						if err := proto.Unmarshal(msg.Data, &dest); err != nil {
							return err
						}

						if dest.Status.Status == dispatch.UNIT_STATUS_USER_ADDED {
							s.userIDToUnitID.Store(*dest.Status.UserId, dest.Status.UnitId)
						} else if dest.Status.Status == dispatch.UNIT_STATUS_USER_REMOVED {
							s.userIDToUnitID.Delete(*dest.Status.UserId)
						}

						unit, ok := s.getUnitsMap(job).Load(dest.Status.UnitId)
						if ok {
							unit.Status = dest.Status
							//unit.Statuses = append(unit.Statuses, &dest)

						} else {
							// "Cache/State miss" load from database
							if err := s.loadUnits(ctx, dest.Status.UnitId); err != nil {
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
					user, err := s.resolveUserById(ctx, userId)
					if err != nil {
						s.logger.Error("failed to get user info from db", zap.Error(err))
						continue
					}

					s.handleRemovedUserFromDisponents(ctx, user)

					if s.handleRemovedUserFromUnit(ctx, user) {
						continue
					}
				}
			}()
		}
	}
}

func (s *Server) handleRemovedUserFromDisponents(ctx context.Context, user *users.UserShort) {
	if s.checkIfUserIsDisponent(user.Job, user.UserId) {
		if err := s.dispatchCenterSignOn(ctx, user.Job, user.UserId, false); err != nil {
			s.logger.Error("failed to remove user from disponents", zap.Error(err))
			return
		}
	}
}

func (s *Server) handleRemovedUserFromUnit(ctx context.Context, user *users.UserShort) bool {
	unitId, ok := s.userIDToUnitID.Load(user.UserId)
	if !ok {
		// Nothing to do
		return false
	}

	unit, ok := s.getUnit(user.Job, unitId)
	if !ok {
		return false
	}

	if err := s.updateUnitAssignments(ctx, &userinfo.UserInfo{
		UserId: user.UserId,
		Job:    user.Job,
	}, unit, nil, []int32{user.UserId}); err != nil {
		s.logger.Error("failed to remove user from unit", zap.Error(err))
		return false
	}

	s.userIDToUnitID.Delete(user.UserId)

	return true
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

// Handle expired dispatch unit assignments
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
	s.dispatches.Range(func(job string, value *xsync.MapOf[uint64, *dispatch.Dispatch]) bool {
		value.Range(func(id uint64, dsp *dispatch.Dispatch) bool {
			for i := len(dsp.Units) - 1; i >= 0; i-- {
				unit, _ := s.getUnit(job, dsp.Units[i].UnitId)

				// If unit isn't empty, continue with the loop
				if len(unit.Users) > 0 {
					continue
				}

				if err := s.updateDispatchAssignments(ctx, job, nil, dsp, nil, []uint64{unit.Id}); err != nil {
					s.logger.Error("failed to remove empty unit from dispatch", zap.Error(err))
					continue
				}
			}

			return true
		})

		return true
	})

	return nil
}
