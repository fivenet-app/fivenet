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

func (s *Server) watchStateEvents() error {
	msgCh := make(chan *nats.Msg, 256)

	sub, err := s.events.JS.ChanSubscribe(fmt.Sprintf("%s.>", BaseSubject), msgCh, nats.DeliverLastPerSubject())
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	// Watch for events from message queue
	for {
		func() error {
			ctx, span := s.tracer.Start(s.ctx, "centrum-watch-events")
			defer span.End()

			select {
			case <-s.ctx.Done():
				return nil

			case msg := <-msgCh:
				msg.Ack()

				job, topic, tType := splitSubject(msg.Subject)

				switch topic {
				case TopicGeneral:
					switch tType {
					case TypeGeneralDisponents:
						var dest DisponentsChange
						if err := proto.Unmarshal(msg.Data, &dest); err != nil {
							return err
						}

						s.state.Disponents.Store(job, dest.Disponents)

					case TypeGeneralSettings:
						var dest dispatch.Settings
						if err := proto.Unmarshal(msg.Data, &dest); err != nil {
							return err
						}

						s.state.Settings.Store(job, &dest)
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

						units, ok := s.state.Units.Load(dest.Job)
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

						if dest.Status.Status == dispatch.StatusUnit_STATUS_UNIT_USER_ADDED {
							s.state.UserIDToUnitID.Store(*dest.Status.UserId, dest.Status.UnitId)
						} else if dest.Status.Status == dispatch.StatusUnit_STATUS_UNIT_USER_REMOVED {
							s.state.UserIDToUnitID.Delete(*dest.Status.UserId)
						}

						unit, ok := s.getUnitsMap(job).Load(dest.Id)
						if ok {
							unit.Status = dest.Status
						} else {
							// "Cache/State miss" load from database
							if err := s.loadUnits(ctx, dest.Id); err != nil {
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

func (s *Server) watchUserChanges() {
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
					unitId, err := s.loadUnitIDForUserID(ctx, userInfo.UserID)
					if err != nil {
						s.logger.Error("failed to load user unit id", zap.Error(err))
						continue
					}
					if unitId == 0 {
						continue
					}

					s.state.UserIDToUnitID.Store(userInfo.UserID, unitId)
				}

				for _, userInfo := range event.Removed {
					user, err := s.resolveUserShortById(ctx, userInfo.UserID)
					if err != nil {
						s.logger.Error("failed to get user info from db", zap.Error(err))
						continue
					}

					s.handleRemovedUserFromDisponents(ctx, user)
					s.handleRemovedUserFromUnit(ctx, user)
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
	unitId, ok := s.state.UserIDToUnitID.Load(user.UserId)
	if !ok {
		// Nothing to do
		return false
	}

	s.state.UserIDToUnitID.Delete(user.UserId)

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

	return true
}

func (s *Server) housekeeper() {
	go s.runHandleDispatchAssignmentExpiration()

	go s.runArchiveDispatches()

	go s.runRemoveDispatchesFromEmptyUnits()

	for {
		select {
		case <-s.ctx.Done():
			return

		case <-time.After(1 * time.Second):
			func() {
				ctx, span := s.tracer.Start(s.ctx, "centrum-housekeeper")
				defer span.End()

				if err := s.checkIfBotsAreNeeded(ctx); err != nil {
					s.logger.Error("failed to clean empty units from dispatches", zap.Error(err))
				}
			}()
		}
	}
}

func (s *Server) runHandleDispatchAssignmentExpiration() {
	for {
		select {
		case <-s.ctx.Done():
			return

		case <-time.After(1 * time.Second):
			func() {
				ctx, span := s.tracer.Start(s.ctx, "centrum-dispatch-assignment-expiration")
				defer span.End()

				if err := s.handleDispatchAssignmentExpiration(ctx); err != nil {
					s.logger.Error("failed to handle expired dispatch assignments", zap.Error(err))
				}
			}()
		}
	}
}

func (s *Server) runArchiveDispatches() {
	for {
		select {
		case <-s.ctx.Done():
			return

		case <-time.After(5 * time.Second):
			func() {
				ctx, span := s.tracer.Start(s.ctx, "centrum-dispatch-archival")
				defer span.End()

				if err := s.archiveDispatches(ctx); err != nil {
					s.logger.Error("failed to archive dispatches", zap.Error(err))
				}

				if err := s.cleanupDispatches(ctx); err != nil {
					s.logger.Error("failed to cleanup dispatches", zap.Error(err))
				}
			}()
		}
	}
}

func (s *Server) runRemoveDispatchesFromEmptyUnits() {
	for {
		select {
		case <-s.ctx.Done():
			return

		case <-time.After(3 * time.Second):
			func() {
				ctx, span := s.tracer.Start(s.ctx, "centrum-units-empty")
				defer span.End()

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
			tDispatchUnit.ExpiresAt.IS_NOT_NULL(),
			tDispatchUnit.ExpiresAt.LT_EQ(jet.NOW()),
		))

	var dest []*struct {
		DispatchID uint64
		UnitID     uint64
		Job        string
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return err
	}

	assignments := map[string]map[uint64][]uint64{}
	for _, ua := range dest {
		if _, ok := assignments[ua.Job]; !ok {
			assignments[ua.Job] = map[uint64][]uint64{}
		}
		if _, ok := assignments[ua.Job][ua.DispatchID]; !ok {
			assignments[ua.Job][ua.DispatchID] = []uint64{}
		}

		assignments[ua.Job][ua.DispatchID] = append(assignments[ua.Job][ua.DispatchID], ua.UnitID)
	}

	for job, dsps := range assignments {
		for dispatchId, units := range dsps {
			dsp, ok := s.getDispatch(job, dispatchId)
			if !ok {
				continue
			}

			if err := s.updateDispatchAssignments(ctx, job, nil, dsp, nil, units); err != nil {
				return err
			}
		}
	}

	return nil
}

// Set `COMPLETED`/`CANCELLED` dispatches to status `ARCHIVED` when the status is older than 5 minutes
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
		// Dispatches that are at 5 minutes or older, have completed/cancelled or no status set
		WHERE(jet.AND(
			tDispatchStatus.CreatedAt.LT_EQ(
				jet.CURRENT_TIMESTAMP().SUB(jet.INTERVAL(5, jet.MINUTE)),
			),
			tDispatchStatus.ID.IS_NULL().OR(
				jet.AND(
					tDispatchStatus.ID.EQ(
						jet.RawInt("SELECT MAX(`dispatchstatus`.`id`) FROM `fivenet_centrum_dispatches_status` AS `dispatchstatus` WHERE `dispatchstatus`.`dispatch_id` = `dispatch`.`id`"),
					),
					tDispatchStatus.Status.IN(
						jet.Int16(int16(dispatch.StatusDispatch_STATUS_DISPATCH_COMPLETED)),
						jet.Int16(int16(dispatch.StatusDispatch_STATUS_DISPATCH_CANCELLED)),
					),
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

		// Ignore already archived dispatches
		if dsp.Status != nil && dsp.Status.Status == dispatch.StatusDispatch_STATUS_DISPATCH_ARCHIVED {
			continue
		}

		if err := s.updateDispatchStatus(ctx, ds.Job, dsp, &dispatch.DispatchStatus{
			DispatchId: dsp.Id,
			Status:     dispatch.StatusDispatch_STATUS_DISPATCH_ARCHIVED,
			UserId:     dsp.Status.UserId,
		}); err != nil {
			return err
		}

		s.getDispatchesMap(ds.Job).Delete(ds.DispatchID)
	}

	return nil
}

func (s *Server) cleanupDispatches(ctx context.Context) error {
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
				jet.CURRENT_TIMESTAMP().SUB(jet.INTERVAL(2, jet.HOUR)),
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

		s.getDispatchesMap(ds.Job).Delete(ds.DispatchID)

		if err := s.deleteDispatch(ctx, dsp.Job, dsp.Id); err != nil {
			return err
		}
	}

	return nil
}

// Remove empty units from dispatches (if no other unit is assigned to dispatch update status to UNASSIGNED) by
// iterating over the dispatches and making sure the assigned units aren't empty
func (s *Server) removeDispatchesFromEmptyUnits(ctx context.Context) error {
	s.state.Dispatches.Range(func(job string, value *xsync.MapOf[uint64, *dispatch.Dispatch]) bool {
		value.Range(func(id uint64, dsp *dispatch.Dispatch) bool {
			for i := len(dsp.Units) - 1; i >= 0; i-- {
				unit, _ := s.getUnit(job, dsp.Units[i].UnitId)
				// If unit isn't empty, continue with the loop
				if unit != nil && len(unit.Users) > 0 {
					continue
				}

				if err := s.updateDispatchAssignments(ctx, job, nil, dsp, nil, []uint64{dsp.Units[i].UnitId}); err != nil {
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

func (s *Server) checkIfBotsAreNeeded(ctx context.Context) error {
	/*s.state.Settings.Range(func(job string, value *dispatch.Settings) bool {
		if s.checkIfBotNeeded(job) {
			if err := s.botManager.Start(job); err != nil {
				s.logger.Error("failed to start dispatch center bot for job", zap.String("job", job))
			}
		} else {
			if err := s.botManager.Stop(job); err != nil {
				s.logger.Error("failed to stop dispatch center bot for job", zap.String("job", job))
			}
		}

		return true
	})*/

	return nil
}
