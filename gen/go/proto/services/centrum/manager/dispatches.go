package manager

import (
	"context"
	"time"

	"github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"github.com/galexrt/fivenet/gen/go/proto/resources/timestamp"
	errorscentrum "github.com/galexrt/fivenet/gen/go/proto/services/centrum/errors"
	eventscentrum "github.com/galexrt/fivenet/gen/go/proto/services/centrum/events"
	centrumutils "github.com/galexrt/fivenet/gen/go/proto/services/centrum/utils"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/pkg/utils/dbutils"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
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
		marker, ok := s.tracker.GetUserById(*userId)
		if ok {
			x = &marker.Info.X
			y = &marker.Info.Y
			postal = marker.Info.Postal
		}
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return errorscentrum.ErrFailedQuery
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

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

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return errorscentrum.ErrFailedQuery
		}

		for i := len(dsp.Units) - 1; i >= 0; i-- {
			if i > len(dsp.Units)-1 {
				break
			}

			for k := 0; k < len(toRemove); k++ {
				if dsp.Units[i].UnitId == toRemove[k] {
					dsp.Units = utils.RemoveFromSlice(dsp.Units, i)

					if err := s.UpdateDispatchStatus(ctx, job, dsp, &dispatch.DispatchStatus{
						DispatchId: dsp.Id,
						UnitId:     &toRemove[k],
						UserId:     userId,
						Status:     dispatch.StatusDispatch_STATUS_DISPATCH_UNIT_UNASSIGNED,
						X:          x,
						Y:          y,
						Postal:     postal,
					}); err != nil {
						return errorscentrum.ErrFailedQuery
					}
					break
				}
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

		stmt := tDispatchUnit.
			INSERT(
				tDispatchUnit.DispatchID,
				tDispatchUnit.UnitID,
				tDispatchUnit.ExpiresAt,
			)
		needInsert := false
		for k := 0; k < len(toAdd); k++ {
			found := false
			for i := 0; i < len(dsp.Units); i++ {
				if dsp.Units[i].UnitId == toAdd[k] {
					found = true
					break
				}
			}

			unit, ok := s.GetUnit(job, toAdd[k])
			if !ok {
				return errorscentrum.ErrFailedQuery
			}

			if len(unit.Users) <= 0 {
				continue
			}

			// Only add unit to dispatch if not already assigned
			if !found {
				stmt = stmt.
					VALUES(
						dsp.Id,
						unit.Id,
						expiresAtVal,
					)
				needInsert = true

				dsp.Units = append(dsp.Units, &dispatch.DispatchAssignment{
					UnitId:     unit.Id,
					DispatchId: dsp.Id,
					Unit:       unit,
					ExpiresAt:  expiresAtTS,
				})

				if err := s.UpdateDispatchStatus(ctx, job, dsp, &dispatch.DispatchStatus{
					DispatchId: dsp.Id,
					UnitId:     &unit.Id,
					UserId:     userId,
					Status:     dispatch.StatusDispatch_STATUS_DISPATCH_UNIT_ASSIGNED,
					X:          x,
					Y:          y,
					Postal:     postal,
				}); err != nil {
					return errorscentrum.ErrFailedQuery
				}
			}
		}

		if needInsert {
			stmt = stmt.ON_DUPLICATE_KEY_UPDATE(
				tDispatchUnit.ExpiresAt.SET(jet.RawTimestamp("VALUES(`expires_at`)")),
			)

			if _, err := stmt.ExecContext(ctx, tx); err != nil {
				if !dbutils.IsDuplicateError(err) {
					return errorscentrum.ErrFailedQuery
				}
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return errorscentrum.ErrFailedQuery
	}

	// Dispatch has not units assigned anymore
	if len(dsp.Units) <= 0 {
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
				return errorscentrum.ErrFailedQuery
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

	s.GetDispatchesMap(job).Delete(id)
	s.State.DispatchLocations[dsp.Job].Remove(dsp, nil)

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
		dsp.CreatorId = nil
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

func (s *Manager) UpdateDispatch(ctx context.Context, userInfo *userinfo.UserInfo, dsp *dispatch.Dispatch) error {
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
			userInfo.Job,
			dsp.Message,
			dsp.Description,
			dsp.Attributes,
			dsp.X,
			dsp.Y,
			dsp.Postal,
			dsp.Anon,
			userInfo.UserId,
		).
		WHERE(jet.AND(
			tDispatch.Job.EQ(jet.String(userInfo.Job)),
			tDispatch.ID.EQ(jet.Uint64(dsp.Id)),
		))

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return errorscentrum.ErrFailedQuery
	}

	s.State.DispatchLocations[dsp.Job].Add(dsp)

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
