package manager

import (
	"context"
	"fmt"
	"time"

	"github.com/galexrt/fivenet/gen/go/proto/resources/centrum"
	"github.com/galexrt/fivenet/gen/go/proto/resources/timestamp"
	errorscentrum "github.com/galexrt/fivenet/gen/go/proto/services/centrum/errors"
	eventscentrum "github.com/galexrt/fivenet/gen/go/proto/services/centrum/events"
	"github.com/galexrt/fivenet/gen/go/proto/services/centrum/state"
	centrumutils "github.com/galexrt/fivenet/gen/go/proto/services/centrum/utils"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/pkg/utils/dbutils"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/paulmach/orb"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const DispatchExpirationTime = 31 * time.Second

func (s *Manager) UpdateDispatchStatus(ctx context.Context, job string, dsp *centrum.Dispatch, in *centrum.DispatchStatus) (*centrum.DispatchStatus, error) {
	// If the dispatch status is the same and is a status that shouldn't be duplicated, don't update the status again
	if dsp.Status != nil &&
		dsp.Status.Status == in.Status &&
		(in.Status == centrum.StatusDispatch_STATUS_DISPATCH_NEW ||
			in.Status == centrum.StatusDispatch_STATUS_DISPATCH_UNASSIGNED) {
		s.logger.Debug("skipping dispatch status update due to same status", zap.Uint64("dispatch_id", dsp.Id), zap.String("status", in.Status.String()))
		return nil, nil
	}

	// If the dispatch is complete, we ignore any unit unassignments/accepts/declines
	if dsp.Status != nil && centrumutils.IsStatusDispatchComplete(dsp.Status.Status) &&
		(in.Status == centrum.StatusDispatch_STATUS_DISPATCH_UNASSIGNED ||
			in.Status == centrum.StatusDispatch_STATUS_DISPATCH_UNIT_UNASSIGNED ||
			in.Status == centrum.StatusDispatch_STATUS_DISPATCH_UNIT_ACCEPTED ||
			in.Status == centrum.StatusDispatch_STATUS_DISPATCH_UNIT_DECLINED) {
		return nil, nil
	}

	s.logger.Debug("updating dispatch status", zap.Uint64("dispatch_id", dsp.Id), zap.String("status", in.Status.String()))

	if in.UserId != nil {
		var err error
		in.User, err = s.resolveUserShortById(ctx, *in.UserId)
		if err != nil {
			return nil, errorscentrum.ErrFailedQuery
		}
	}

	tDispatchStatus := table.FivenetCentrumDispatchesStatus
	stmt := tDispatchStatus.
		INSERT(
			tDispatchStatus.CreatedAt,
			tDispatchStatus.DispatchID,
			tDispatchStatus.UnitID,
			tDispatchStatus.Status,
			tDispatchStatus.Reason,
			tDispatchStatus.Code,
			tDispatchStatus.UserID,
			tDispatchStatus.X,
			tDispatchStatus.Y,
			tDispatchStatus.Postal,
		).
		VALUES(
			jet.CURRENT_TIMESTAMP(),
			in.DispatchId,
			in.UnitId,
			in.Status,
			in.Reason,
			in.Code,
			in.UserId,
			in.X,
			in.Y,
			in.Postal,
		)

	res, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		return nil, errorscentrum.ErrFailedQuery
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return nil, errorscentrum.ErrFailedQuery
	}
	in.Id = uint64(lastId)

	if err := s.State.UpdateDispatchStatus(ctx, job, in.DispatchId, in); err != nil {
		return nil, err
	}

	data, err := proto.Marshal(in)
	if err != nil {
		return nil, errorscentrum.ErrFailedQuery
	}

	if _, err := s.js.Publish(eventscentrum.BuildSubject(eventscentrum.TopicDispatch, eventscentrum.TypeDispatchStatus, job), data); err != nil {
		return nil, fmt.Errorf("failed to publish dispatch status event (size: %d, message: '%+v'): %w", len(data), dsp, err)
	}

	return in, nil
}

func (s *Manager) DispatchAssignmentExpirationTime() time.Time {
	return time.Now().Add(DispatchExpirationTime)
}

func (s *Manager) UpdateDispatchAssignments(ctx context.Context, job string, userId *int32, dspId uint64, toAdd []uint64, toRemove []uint64, expiresAt time.Time) error {
	s.logger.Debug("updating dispatch assignments", zap.String("job", job), zap.Uint64("dispatch_id", dspId), zap.Uint64s("toAdd", toAdd), zap.Uint64s("toRemove", toRemove))

	var x, y *float64
	var postal *string
	if userId != nil {
		if marker, ok := s.tracker.GetUserById(*userId); ok {
			x = &marker.Info.X
			y = &marker.Info.Y
			postal = marker.Info.Postal
		}
	}

	tDispatchUnit := table.FivenetCentrumDispatchesAsgmts

	if len(toRemove) > 0 {
		removeIds := make([]jet.Expression, len(toRemove))
		for i := 0; i < len(toRemove); i++ {
			removeIds[i] = jet.Uint64(toRemove[i])
		}

		stmt := tDispatchUnit.
			DELETE().
			WHERE(jet.AND(
				tDispatchUnit.DispatchID.EQ(jet.Uint64(dspId)),
				tDispatchUnit.UnitID.IN(removeIds...),
			))

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return err
		}
	}

	var expiresAtTS *timestamp.Timestamp
	// If expires at time is not zero
	expiresAtVal := jet.NULL
	if !expiresAt.IsZero() {
		expiresAtTS = timestamp.New(expiresAt)
		expiresAtVal = jet.TimeT(expiresAt)
	}

	if len(toAdd) > 0 {
		units := []uint64{}
		for i := 0; i < len(toAdd); i++ {
			dsp := s.GetDispatch(job, toAdd[i])
			if dsp == nil {
				continue
			}

			// Skip already added units
			if utils.InSliceFunc(dsp.Units, func(in *centrum.DispatchAssignment) bool {
				return in.UnitId == toAdd[i]
			}) {
				continue
			}

			unit := s.GetUnit(job, toAdd[i])
			if unit == nil {
				continue
			}

			// Skip empty units
			if len(unit.Users) == 0 {
				continue
			}

			// Only add unit to dispatch if not already assigned/in list
			units = append(units, toAdd[i])
		}

		if len(units) > 0 {
			stmt := tDispatchUnit.
				INSERT(
					tDispatchUnit.DispatchID,
					tDispatchUnit.UnitID,
					tDispatchUnit.ExpiresAt,
				)

			for _, unitId := range units {
				stmt = stmt.
					VALUES(
						dspId,
						unitId,
						expiresAtVal,
					)
			}

			stmt = stmt.ON_DUPLICATE_KEY_UPDATE(
				tDispatchUnit.ExpiresAt.SET(jet.RawTimestamp("VALUES(`expires_at`)")),
			)

			if _, err := stmt.ExecContext(ctx, s.db); err != nil {
				if !dbutils.IsDuplicateError(err) {
					return err
				}
			}
		}
	}

	store := s.State.DispatchesStore()

	key := state.JobIdKey(job, dspId)
	if err := store.ComputeUpdate(key, true, func(key string, dsp *centrum.Dispatch) (*centrum.Dispatch, error) {
		if len(toRemove) > 0 {
			toAnnounce := []uint64{}
			for i := len(dsp.Units) - 1; i >= 0; i-- {
				if i > (len(dsp.Units) - 1) {
					break
				}

				for k := 0; k < len(toRemove); k++ {
					if dsp.Units[i].UnitId != toRemove[k] {
						continue
					}

					dsp.Units = utils.RemoveFromSlice(dsp.Units, i)
					toAnnounce = append(toAnnounce, toRemove[k])
				}
			}

			// Send updates
			for _, unitId := range toAnnounce {
				if err := s.AddDispatchStatus(ctx, s.db, job, &centrum.DispatchStatus{
					DispatchId: dsp.Id,
					UnitId:     &unitId,
					Status:     centrum.StatusDispatch_STATUS_DISPATCH_UNIT_UNASSIGNED,
					UserId:     userId,
					X:          x,
					Y:          y,
					Postal:     postal,
				}); err != nil {
					return nil, err
				}
			}
		}

		if len(toAdd) > 0 {
			units := []uint64{}
			for i := 0; i < len(toAdd); i++ {
				// Skip already added units
				if utils.InSliceFunc(dsp.Units, func(in *centrum.DispatchAssignment) bool {
					return in.UnitId == toAdd[i]
				}) {
					continue
				}

				unit := s.GetUnit(job, toAdd[i])
				if unit == nil {
					continue
				}

				// Skip empty units
				if len(unit.Users) == 0 {
					continue
				}

				// Only add unit to dispatch if not already assigned/in list
				units = append(units, toAdd[i])
			}

			for _, unitId := range units {
				unit := s.GetUnit(job, unitId)
				if unit == nil {
					continue
				}

				dsp.Units = append(dsp.Units, &centrum.DispatchAssignment{
					DispatchId: dsp.Id,
					UnitId:     unit.Id,
					Unit:       unit,
					ExpiresAt:  expiresAtTS,
				})
			}

			for _, unitId := range units {
				if err := s.AddDispatchStatus(ctx, s.db, job, &centrum.DispatchStatus{
					DispatchId: dsp.Id,
					UnitId:     &unitId,
					UserId:     userId,
					Status:     centrum.StatusDispatch_STATUS_DISPATCH_UNIT_ASSIGNED,
					X:          x,
					Y:          y,
					Postal:     postal,
				}); err != nil {
					return nil, err
				}
			}
		}

		// Dispatch has not units assigned anymore
		if len(dsp.Units) == 0 {
			// Check dispatch status to not be completed/archived, etc.
			if dsp.Status != nil && !centrumutils.IsStatusDispatchComplete(dsp.Status.Status) {
				dsp.Status = &centrum.DispatchStatus{
					DispatchId: dsp.Id,
					Status:     centrum.StatusDispatch_STATUS_DISPATCH_UNASSIGNED,
					UserId:     userId,
					X:          x,
					Y:          y,
					Postal:     postal,
				}
				if err := s.AddDispatchStatus(ctx, s.db, job, dsp.Status); err != nil {
					return nil, err
				}
			}
		}

		return dsp, nil
	}); err != nil {
		return err
	}

	return nil
}

func (s *Manager) DeleteDispatch(ctx context.Context, job string, id uint64, allTheWay bool) error {
	if allTheWay {
		stmt := tDispatch.
			DELETE().
			WHERE(jet.AND(
				tDispatch.Job.EQ(jet.String(job)),
				tDispatch.ID.EQ(jet.Uint64(id)),
			)).
			LIMIT(1)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return errorscentrum.ErrFailedQuery
		}
	}

	dsp := s.GetDispatch(job, id)
	if dsp == nil {
		return nil
	}

	data, err := proto.Marshal(dsp)
	if err != nil {
		return errorscentrum.ErrFailedQuery
	}

	if _, err := s.js.Publish(eventscentrum.BuildSubject(eventscentrum.TopicDispatch, eventscentrum.TypeDispatchDeleted, job), data); err != nil {
		return err
	}

	s.State.DeleteDispatch(job, id)

	return nil
}

func (s *Manager) CreateDispatch(ctx context.Context, dsp *centrum.Dispatch) (*centrum.Dispatch, error) {
	postal := s.postals.Closest(dsp.X, dsp.Y)
	if postal != nil {
		dsp.Postal = postal.Code
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	tDispatch := table.FivenetCentrumDispatches
	stmt := tDispatch.
		INSERT(
			tDispatch.CreatedAt,
			tDispatch.Job,
			tDispatch.Message,
			tDispatch.Description,
			tDispatch.Attributes,
			tDispatch.X,
			tDispatch.Y,
			tDispatch.Postal,
			tDispatch.Anon,
			tDispatch.CreatorID,
		).
		VALUES(
			jet.CURRENT_TIMESTAMP(),
			dsp.Job,
			dsp.Message,
			dsp.Description,
			dsp.Attributes,
			dsp.X,
			dsp.Y,
			dsp.Postal,
			dsp.Anon,
			dsp.CreatorId,
		)

	result, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return nil, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	dsp.Id = uint64(lastId)

	var userId *int32
	if !dsp.Anon && dsp.CreatorId != nil {
		userId = dsp.CreatorId
	}

	dsp.Status = &centrum.DispatchStatus{
		DispatchId: dsp.Id,
		UserId:     userId,
		Status:     centrum.StatusDispatch_STATUS_DISPATCH_NEW,
		X:          &dsp.X,
		Y:          &dsp.Y,
		Postal:     dsp.Postal,
	}
	if err := s.AddDispatchStatus(ctx, tx, dsp.Job, dsp.Status); err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	if err := s.State.UpdateDispatch(ctx, dsp.Job, dsp.Id, dsp); err != nil {
		return nil, err
	}

	// Hide user info when dispatch is anonymous
	if dsp.Anon {
		dsp.Creator = nil
	}
	metricsDispatchLastID.WithLabelValues(dsp.Job).Set(float64(lastId))

	// Make sure dispatch is in the locations list
	s.State.GetDispatchLocations(dsp.Job).Add(dsp)

	data, err := proto.Marshal(dsp)
	if err != nil {
		return nil, err
	}

	if _, err := s.js.Publish(eventscentrum.BuildSubject(eventscentrum.TopicDispatch, eventscentrum.TypeDispatchCreated, dsp.Job), data); err != nil {
		return nil, err
	}

	return dsp, nil
}

func (s *Manager) UpdateDispatch(ctx context.Context, userJob string, userId *int32, dsp *centrum.Dispatch, publish bool) (*centrum.Dispatch, error) {
	dsp.UpdatedAt = timestamp.Now()
	stmt := tDispatch.
		UPDATE(
			tDispatch.UpdatedAt,
			tDispatch.Job,
			tDispatch.Message,
			tDispatch.Description,
			tDispatch.Attributes,
			tDispatch.X,
			tDispatch.Y,
			tDispatch.Postal,
			tDispatch.Anon,
			tDispatch.CreatorID,
		).
		SET(
			jet.CURRENT_TIMESTAMP(),
			dsp.Job,
			dsp.Message,
			dsp.Description,
			dsp.Attributes,
			dsp.X,
			dsp.Y,
			dsp.Postal,
			dsp.Anon,
			dsp.CreatorId,
		).
		WHERE(jet.AND(
			tDispatch.Job.EQ(jet.String(userJob)),
			tDispatch.ID.EQ(jet.Uint64(dsp.Id)),
		))

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	// Make sure dispatch is in the locations list
	if !s.State.GetDispatchLocations(dsp.Job).Has(dsp, func(p orb.Pointer) bool {
		return p.(*centrum.Dispatch).Id == dsp.Id
	}) {
		s.State.GetDispatchLocations(dsp.Job).Add(dsp)
	}

	if err := s.State.UpdateDispatch(ctx, dsp.Job, dsp.Id, dsp); err != nil {
		return nil, err
	}

	if publish {
		data, err := proto.Marshal(dsp)
		if err != nil {
			return nil, err
		}

		if _, err := s.js.Publish(eventscentrum.BuildSubject(eventscentrum.TopicDispatch, eventscentrum.TypeDispatchUpdated, userJob), data); err != nil {
			return nil, err
		}
	}

	return dsp, nil
}

func (s *Manager) AddDispatchStatus(ctx context.Context, tx qrm.DB, job string, status *centrum.DispatchStatus) error {
	tDispatchStatus := table.FivenetCentrumDispatchesStatus
	stmt := tDispatchStatus.
		INSERT(
			tDispatchStatus.DispatchID,
			tDispatchStatus.Status,
			tDispatchStatus.Reason,
			tDispatchStatus.Code,
			tDispatchStatus.UnitID,
			tDispatchStatus.UserID,
			tDispatchStatus.X,
			tDispatchStatus.Y,
			tDispatchStatus.Postal,
		).
		VALUES(
			status.DispatchId,
			status.Status,
			status.Reason,
			status.Code,
			status.UnitId,
			status.UserId,
			status.X,
			status.Y,
			status.Postal,
		)

	res, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return errorscentrum.ErrFailedQuery
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return errorscentrum.ErrFailedQuery
	}
	status.Id = uint64(lastId)

	data, err := proto.Marshal(status)
	if err != nil {
		return err
	}

	if _, err := s.js.Publish(eventscentrum.BuildSubject(eventscentrum.TopicDispatch, eventscentrum.TypeDispatchStatus, job), data); err != nil {
		return err
	}

	return nil
}

func (s *Manager) TakeDispatch(ctx context.Context, job string, userId int32, unitId uint64, resp centrum.TakeDispatchResp, dispatchIds []uint64) error {
	unit := s.GetUnit(job, unitId)
	if unit == nil {
		return errorscentrum.ErrFailedQuery
	}

	settings := s.GetSettings(job)
	var x, y *float64
	var postal *string
	if marker, ok := s.tracker.GetUserById(userId); ok {
		x = &marker.Info.X
		y = &marker.Info.Y
		postal = marker.Info.Postal
	}

	tDispatchUnit := table.FivenetCentrumDispatchesAsgmts

	store := s.State.DispatchesStore()

	for _, dispatchId := range dispatchIds {
		dsp := s.GetDispatch(job, dispatchId)
		if dsp == nil {
			return errorscentrum.ErrFailedQuery
		}

		key := state.JobIdKey(job, dispatchId)
		if err := store.ComputeUpdate(key, true, func(key string, dsp *centrum.Dispatch) (*centrum.Dispatch, error) {
			// If the dispatch center is in central command mode, units can't self assign dispatches
			if settings.Mode == centrum.CentrumMode_CENTRUM_MODE_CENTRAL_COMMAND {
				if !utils.InSliceFunc(dsp.Units, func(in *centrum.DispatchAssignment) bool {
					return in.UnitId == unitId
				}) {
					return nil, errorscentrum.ErrModeForbidsAction
				}
			}

			// If dispatch is completed, disallow to accept the dispatch
			if dsp.Status != nil && centrumutils.IsStatusDispatchComplete(dsp.Status.Status) {
				return nil, errorscentrum.ErrDispatchAlreadyCompleted
			}

			status := centrum.StatusDispatch_STATUS_DISPATCH_UNSPECIFIED

			// Dispatch accepted
			if resp == centrum.TakeDispatchResp_TAKE_DISPATCH_RESP_ACCEPTED {
				status = centrum.StatusDispatch_STATUS_DISPATCH_UNIT_ACCEPTED

				stmt := tDispatchUnit.
					INSERT(
						tDispatchUnit.DispatchID,
						tDispatchUnit.UnitID,
						tDispatchUnit.ExpiresAt,
					).
					VALUES(
						dsp.Id,
						unit.Id,
						jet.NULL,
					).
					ON_DUPLICATE_KEY_UPDATE(
						tDispatchUnit.ExpiresAt.SET(jet.TimestampExp(jet.NULL)),
					)

				if _, err := stmt.ExecContext(ctx, s.db); err != nil {
					if !dbutils.IsDuplicateError(err) {
						return nil, errorscentrum.ErrFailedQuery
					}
				}

				found := false
				accepted := true
				// Set unit expires at to nil
				for _, ua := range dsp.Units {
					if ua.UnitId == unit.Id {
						found = true
						// If there's no expiration time the unit has been directly assigned
						if ua.ExpiresAt == nil {
							accepted = false
						}
						ua.ExpiresAt = nil
						break
					}
				}

				if !found {
					dsp.Units = append(dsp.Units, &centrum.DispatchAssignment{
						DispatchId: dsp.Id,
						UnitId:     unit.Id,
						Unit:       unit,
						CreatedAt:  timestamp.Now(),
					})
				}

				if accepted {
					// Set unit to busy when unit accepts a dispatch
					if unit.Status == nil || unit.Status.Status != centrum.StatusUnit_STATUS_UNIT_BUSY {
						if _, err := s.UpdateUnitStatus(ctx, job, unit, &centrum.UnitStatus{
							UnitId:    unit.Id,
							Status:    centrum.StatusUnit_STATUS_UNIT_BUSY,
							UserId:    &userId,
							CreatorId: &userId,
							X:         x,
							Y:         y,
							Postal:    postal,
						}); err != nil {
							return nil, errorscentrum.ErrFailedQuery
						}
					}
				}
			} else {
				// Dispatch declined
				status = centrum.StatusDispatch_STATUS_DISPATCH_UNIT_DECLINED

				stmt := tDispatchUnit.
					DELETE().
					WHERE(jet.AND(
						tDispatchUnit.DispatchID.EQ(jet.Uint64(dsp.Id)),
						tDispatchUnit.UnitID.EQ(jet.Uint64(unit.Id)),
					)).
					LIMIT(1)

				if _, err := stmt.ExecContext(ctx, s.db); err != nil {
					if !dbutils.IsDuplicateError(err) {
						return nil, errorscentrum.ErrFailedQuery
					}
				}

				// Remove the unit's assignment
				for k, u := range dsp.Units {
					if u.UnitId == unit.Id {
						dsp.Units = utils.RemoveFromSlice(dsp.Units, k)
					}
				}
			}

			if _, err := s.UpdateDispatchStatus(ctx, job, dsp, &centrum.DispatchStatus{
				DispatchId: dispatchId,
				Status:     status,
				UnitId:     &unitId,
				UserId:     &userId,
				X:          x,
				Y:          y,
				Postal:     postal,
			}); err != nil {
				return nil, errorscentrum.ErrFailedQuery
			}

			return dsp, nil
		}); err != nil {
			return err
		}
	}

	return nil
}
