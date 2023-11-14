package manager

import (
	"context"
	"time"

	"github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"github.com/galexrt/fivenet/gen/go/proto/resources/timestamp"
	errorscentrum "github.com/galexrt/fivenet/gen/go/proto/services/centrum/errors"
	eventscentrum "github.com/galexrt/fivenet/gen/go/proto/services/centrum/events"
	centrumutils "github.com/galexrt/fivenet/gen/go/proto/services/centrum/utils"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/pkg/utils/dbutils"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/paulmach/orb"
	"google.golang.org/protobuf/proto"
)

const DispatchExpirationTime = 31 * time.Second

func (s *Manager) UpdateDispatchStatus(ctx context.Context, job string, dsp *dispatch.Dispatch, in *dispatch.DispatchStatus) error {
	// If the dispatch status is the same and is a status that shouldn't be duplicated, don't update the status again
	if dsp.Status != nil &&
		dsp.Status.Status == in.Status &&
		(in.Status == dispatch.StatusDispatch_STATUS_DISPATCH_NEW ||
			in.Status == dispatch.StatusDispatch_STATUS_DISPATCH_UNASSIGNED ||
			centrumutils.IsStatusDispatchComplete(in.Status)) {
		return nil
	}

	// If the dispatch is complete, we ignore any unit unassignments/accepts/declines
	if dsp.Status != nil && centrumutils.IsStatusDispatchComplete(dsp.Status.Status) &&
		(in.Status == dispatch.StatusDispatch_STATUS_DISPATCH_UNASSIGNED ||
			in.Status == dispatch.StatusDispatch_STATUS_DISPATCH_UNIT_UNASSIGNED ||
			in.Status == dispatch.StatusDispatch_STATUS_DISPATCH_UNIT_ACCEPTED ||
			in.Status == dispatch.StatusDispatch_STATUS_DISPATCH_UNIT_DECLINED) {
		return nil
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
		return errorscentrum.ErrFailedQuery
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return errorscentrum.ErrFailedQuery
	}

	status, err := s.getDispatchStatusFromDB(ctx, job, uint64(lastId))
	if err != nil {
		return errorscentrum.ErrFailedQuery
	}
	dsp.Status = status

	data, err := proto.Marshal(dsp)
	if err != nil {
		return errorscentrum.ErrFailedQuery
	}
	s.events.JS.PublishAsync(eventscentrum.BuildSubject(eventscentrum.TopicDispatch, eventscentrum.TypeDispatchStatus, job, 0), data)

	return nil
}

func (s *Manager) DispatchAssignmentExpirationTime() time.Time {
	return time.Now().Add(DispatchExpirationTime)
}

func (s *Manager) UpdateDispatchAssignments(ctx context.Context, job string, userId *int32, dsp *dispatch.Dispatch, toAdd []uint64, toRemove []uint64, expiresAt time.Time) error {
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
				tDispatchUnit.DispatchID.EQ(jet.Uint64(dsp.Id)),
				tDispatchUnit.UnitID.IN(removeIds...),
			))

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return err
		}

		toAnnounce := []uint64{}
		for i := len(dsp.Units) - 1; i >= 0; i-- {
			if i > len(dsp.Units)-1 {
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
		for _, unit := range toAnnounce {
			if err := s.UpdateDispatchStatus(ctx, job, dsp, &dispatch.DispatchStatus{
				DispatchId: dsp.Id,
				UnitId:     &unit,
				Status:     dispatch.StatusDispatch_STATUS_DISPATCH_UNIT_UNASSIGNED,
				UserId:     userId,
				X:          x,
				Y:          y,
				Postal:     postal,
			}); err != nil {
				return err
			}
		}
	}

	if len(toAdd) > 0 {
		// If expires at time is not zero
		var expiresAtTS *timestamp.Timestamp
		expiresAtVal := jet.NULL
		if !expiresAt.IsZero() {
			expiresAtTS = timestamp.New(expiresAt)
			expiresAtVal = jet.TimeT(expiresAt)
		}

		units := []uint64{}
		for i := 0; i < len(toAdd); i++ {
			// Skip already added units
			if utils.InSliceFunc(dsp.Units, func(in *dispatch.DispatchAssignment) bool {
				return in.UnitId == toAdd[i]
			}) {
				continue
			}

			unit, ok := s.GetUnit(job, toAdd[i])
			if !ok {
				continue
			}

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
						dsp.Id,
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

		for _, unitId := range units {
			unit, ok := s.GetUnit(job, unitId)
			if !ok {
				continue
			}

			dsp.Units = append(dsp.Units, &dispatch.DispatchAssignment{
				DispatchId: dsp.Id,
				UnitId:     unit.Id,
				Unit:       unit,
				ExpiresAt:  expiresAtTS,
			})
		}

		for _, unitId := range units {
			if err := s.UpdateDispatchStatus(ctx, job, dsp, &dispatch.DispatchStatus{
				DispatchId: dsp.Id,
				UnitId:     &unitId,
				UserId:     userId,
				Status:     dispatch.StatusDispatch_STATUS_DISPATCH_UNIT_ASSIGNED,
				X:          x,
				Y:          y,
				Postal:     postal,
			}); err != nil {
				return err
			}
		}
	}

	// Dispatch has not units assigned anymore
	if len(dsp.Units) == 0 {
		// Check dispatch status to not be completed/archived, etc.
		if dsp.Status != nil && !centrumutils.IsStatusDispatchComplete(dsp.Status.Status) {
			if err := s.UpdateDispatchStatus(ctx, job, dsp, &dispatch.DispatchStatus{
				DispatchId: dsp.Id,
				Status:     dispatch.StatusDispatch_STATUS_DISPATCH_UNASSIGNED,
				UserId:     userId,
				X:          x,
				Y:          y,
				Postal:     postal,
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Manager) DeleteDispatch(ctx context.Context, job string, id uint64) error {
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

	dsp, ok := s.GetDispatch(job, id)
	if !ok {
		return nil
	}

	data, err := proto.Marshal(dsp)
	if err != nil {
		return errorscentrum.ErrFailedQuery
	}
	s.events.JS.PublishAsync(eventscentrum.BuildSubject(eventscentrum.TopicDispatch, eventscentrum.TypeDispatchDeleted, job, 0), data)

	s.State.DeleteDispatch(job, id)

	return nil
}

func (s *Manager) getDispatchStatusFromDB(ctx context.Context, job string, id uint64) (*dispatch.DispatchStatus, error) {
	stmt := tDispatchStatus.
		SELECT(
			tDispatchStatus.ID,
			tDispatchStatus.CreatedAt,
			tDispatchStatus.DispatchID,
			tDispatchStatus.UnitID,
			tDispatchStatus.Status,
			tDispatchStatus.Reason,
			tDispatchStatus.Code,
			tDispatchStatus.UserID,
			tUsers.ID,
			tUsers.Identifier,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Dateofbirth,
			tUsers.PhoneNumber,
		).
		FROM(
			tDispatchStatus.
				LEFT_JOIN(
					tUsers,
					tUsers.ID.EQ(tDispatchStatus.UserID),
				),
		).
		WHERE(
			tDispatchStatus.ID.EQ(jet.Uint64(id)),
		).
		LIMIT(1)

	var dest dispatch.DispatchStatus
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, err
	}

	if dest.UnitId != nil {
		dest.Unit, _ = s.GetUnit(job, *dest.UnitId)
	}

	return &dest, nil
}

func (s *Manager) CreateDispatch(ctx context.Context, d *dispatch.Dispatch) (*dispatch.Dispatch, error) {
	postal := s.postals.Closest(d.X, d.Y)
	if postal != nil {
		d.Postal = postal.Code
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
			d.Job,
			d.Message,
			d.Description,
			d.Attributes,
			d.X,
			d.Y,
			d.Postal,
			d.Anon,
			d.CreatorId,
		)

	result, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return nil, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	if err := s.AddDispatchStatus(ctx, tx, &dispatch.DispatchStatus{
		DispatchId: uint64(lastId),
		UserId:     d.CreatorId,
		Status:     dispatch.StatusDispatch_STATUS_DISPATCH_NEW,
		X:          &d.X,
		Y:          &d.Y,
		Postal:     d.Postal,
	}); err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Load dispatch into cache
	if err := s.LoadDispatches(ctx, uint64(lastId)); err != nil {
		return nil, err
	}

	dsp, ok := s.GetDispatch(d.Job, uint64(lastId))
	if !ok {
		return nil, err
	}
	// Hide user info when dispatch is anonymous
	if dsp.Anon {
		dsp.Creator = nil
	}
	metricsDispatchLastID.WithLabelValues(d.Job).Set(float64(lastId))

	data, err := proto.Marshal(dsp)
	if err != nil {
		return nil, err
	}
	s.events.JS.PublishAsync(eventscentrum.BuildSubject(eventscentrum.TopicDispatch, eventscentrum.TypeDispatchCreated, d.Job, 0), data)

	s.State.DispatchLocations[dsp.Job].Add(dsp)

	return dsp, nil
}

func (s *Manager) UpdateDispatch(ctx context.Context, userJob string, userId *int32, dsp *dispatch.Dispatch, publish bool) error {
	stmt := tDispatch.
		UPDATE(
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
		return err
	}

	if !s.State.DispatchLocations[dsp.Job].Has(dsp, func(p orb.Pointer) bool {
		return p.(*dispatch.Dispatch).Id == dsp.Id
	}) {
		s.State.DispatchLocations[dsp.Job].Add(dsp)
	}

	// Load dispatch into cache
	if err := s.LoadDispatches(ctx, dsp.Id); err != nil {
		return err
	}

	dsp, ok := s.GetDispatch(userJob, dsp.Id)
	if !ok {
		return errorscentrum.ErrFailedQuery
	}

	if publish {
		data, err := proto.Marshal(dsp)
		if err != nil {
			return err
		}

		s.events.JS.PublishAsync(eventscentrum.BuildSubject(eventscentrum.TopicDispatch, eventscentrum.TypeDispatchUpdated, userJob, 0), data)
	}

	return nil
}

func (s *Manager) AddDispatchStatus(ctx context.Context, tx qrm.DB, status *dispatch.DispatchStatus) error {
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

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return errorscentrum.ErrFailedQuery
	}

	return nil
}

func (s *Manager) TakeDispatch(ctx context.Context, job string, userId int32, unitId uint64, resp dispatch.TakeDispatchResp, dispatchIds []uint64) error {
	unit, ok := s.GetUnit(job, unitId)
	if !ok {
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

	for _, dispatchId := range dispatchIds {
		dsp, ok := s.GetDispatch(job, dispatchId)
		if !ok {
			return errorscentrum.ErrFailedQuery
		}

		// If the dispatch center is in central command mode, units can't self assign dispatches
		if settings.Mode == dispatch.CentrumMode_CENTRUM_MODE_CENTRAL_COMMAND {
			if !utils.InSliceFunc(dsp.Units, func(in *dispatch.DispatchAssignment) bool {
				return in.UnitId == unitId
			}) {
				return errorscentrum.ErrModeForbidsAction
			}
		}

		// If dispatch is completed, disallow to accept the dispatch
		if dsp.Status != nil && centrumutils.IsStatusDispatchComplete(dsp.Status.Status) {
			return errorscentrum.ErrDispatchAlreadyCompleted
		}

		status := dispatch.StatusDispatch_STATUS_DISPATCH_UNSPECIFIED

		tDispatchUnit := table.FivenetCentrumDispatchesAsgmts
		// Dispatch accepted
		if resp == dispatch.TakeDispatchResp_TAKE_DISPATCH_RESP_ACCEPTED {
			status = dispatch.StatusDispatch_STATUS_DISPATCH_UNIT_ACCEPTED

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
					return errorscentrum.ErrFailedQuery
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
				dsp.Units = append(dsp.Units, &dispatch.DispatchAssignment{
					DispatchId: dsp.Id,
					UnitId:     unit.Id,
					Unit:       unit,
					CreatedAt:  timestamp.Now(),
				})
			}

			if accepted {
				// Set unit to busy when unit accepts a dispatch
				if unit.Status == nil || unit.Status.Status != dispatch.StatusUnit_STATUS_UNIT_BUSY {
					if err := s.UpdateUnitStatus(ctx, job, unit, &dispatch.UnitStatus{
						UnitId:    unit.Id,
						Status:    dispatch.StatusUnit_STATUS_UNIT_BUSY,
						UserId:    &userId,
						CreatorId: &userId,
						X:         x,
						Y:         y,
						Postal:    postal,
					}); err != nil {
						return errorscentrum.ErrFailedQuery
					}
				}
			}
		} else {
			// Dispatch declined
			status = dispatch.StatusDispatch_STATUS_DISPATCH_UNIT_DECLINED

			stmt := tDispatchUnit.
				DELETE().
				WHERE(jet.AND(
					tDispatchUnit.DispatchID.EQ(jet.Uint64(dsp.Id)),
					tDispatchUnit.UnitID.EQ(jet.Uint64(unit.Id)),
				)).
				LIMIT(1)

			if _, err := stmt.ExecContext(ctx, s.db); err != nil {
				if !dbutils.IsDuplicateError(err) {
					return errorscentrum.ErrFailedQuery
				}
			}

			// Remove the unit's assignment
			for k, u := range dsp.Units {
				if u.UnitId == unit.Id {
					dsp.Units = utils.RemoveFromSlice(dsp.Units, k)
					break
				}
			}
		}

		if err := s.UpdateDispatchStatus(ctx, job, dsp, &dispatch.DispatchStatus{
			DispatchId: dispatchId,
			Status:     status,
			UnitId:     &unitId,
			UserId:     &userId,
			X:          x,
			Y:          y,
			Postal:     postal,
		}); err != nil {
			return errorscentrum.ErrFailedQuery
		}
	}

	return nil
}
