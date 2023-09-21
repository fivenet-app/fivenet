package centrum

import (
	"context"
	"time"

	dispatch "github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"github.com/galexrt/fivenet/gen/go/proto/resources/timestamp"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/pkg/utils/dbutils"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/proto"
)

const DispatchExpirationTime = 31 * time.Second

func (s *Server) getDispatch(job string, id uint64) (*dispatch.Dispatch, bool) {
	dispatches, ok := s.dispatches.Load(job)
	if !ok {
		return nil, false
	}

	return dispatches.Load(id)
}

func (s *Server) listDispatches(job string) ([]*dispatch.Dispatch, error) {
	ds := []*dispatch.Dispatch{}

	dispatches, ok := s.dispatches.Load(job)
	if !ok {
		return nil, nil
	}

	dispatches.Range(func(id uint64, dispatch *dispatch.Dispatch) bool {
		ds = append(ds, dispatch)
		return true
	})

	slices.SortFunc(ds, func(a, b *dispatch.Dispatch) int {
		return int(b.Id - a.Id)
	})

	return ds, nil
}

func (s *Server) getDispatchStatusFromDB(ctx context.Context, id uint64) (*dispatch.DispatchStatus, error) {
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

	return &dest, nil
}

func (s *Server) addDispatchStatus(ctx context.Context, tx qrm.DB, status *dispatch.DispatchStatus) error {
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
		return err
	}

	return nil
}

func (s *Server) updateDispatchStatus(ctx context.Context, job string, dsp *dispatch.Dispatch, in *dispatch.DispatchStatus) error {
	// If the dispatch status is the same and is a status that shouldn't be duplicated, don't update the status again
	if dsp.Status != nil &&
		dsp.Status.Status == in.Status &&
		(in.Status == dispatch.StatusDispatch_STATUS_DISPATCH_NEW ||
			in.Status == dispatch.StatusDispatch_STATUS_DISPATCH_UNASSIGNED ||
			in.Status == dispatch.StatusDispatch_STATUS_DISPATCH_COMPLETED ||
			in.Status == dispatch.StatusDispatch_STATUS_DISPATCH_CANCELLED ||
			in.Status == dispatch.StatusDispatch_STATUS_DISPATCH_ARCHIVED) {
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
		return err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return err
	}

	status, err := s.getDispatchStatusFromDB(ctx, uint64(lastId))
	if err != nil {
		return err
	}
	dsp.Status = status

	data, err := proto.Marshal(dsp)
	if err != nil {
		return err
	}

	if len(dsp.Units) == 0 {
		s.events.JS.PublishAsync(s.buildSubject(TopicDispatch, TypeDispatchStatus, job, 0), data)
	} else {
		for _, u := range dsp.Units {
			s.events.JS.PublishAsync(s.buildSubject(TopicDispatch, TypeDispatchStatus, job, u.UnitId), data)
		}
	}

	return nil
}

func (s *Server) updateDispatchAssignments(ctx context.Context, job string, userId *int32, dsp *dispatch.Dispatch, toAdd []uint64, toRemove []uint64) error {
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
		return ErrFailedQuery
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
			return err
		}

		for i := 0; i < len(dsp.Units); i++ {
			for k := 0; k < len(toRemove); k++ {
				if dsp.Units[i].UnitId == toRemove[k] {
					if err := s.updateDispatchStatus(ctx, job, dsp, &dispatch.DispatchStatus{
						DispatchId: dsp.Id,
						UnitId:     &toRemove[k],
						UserId:     userId,
						Status:     dispatch.StatusDispatch_STATUS_DISPATCH_UNIT_UNASSIGNED,
						X:          x,
						Y:          y,
						Postal:     postal,
					}); err != nil {
						return ErrFailedQuery
					}

					dsp.Units = utils.RemoveFromSlice(dsp.Units, i)

					break
				}
			}
		}
	}

	if len(toAdd) > 0 {
		expiresAt := time.Now().Add(DispatchExpirationTime)
		expiresAtTS := timestamp.New(expiresAt)
		for k := 0; k < len(toAdd); k++ {
			found := false
			for i := 0; i < len(dsp.Units); i++ {
				if dsp.Units[i].UnitId == toAdd[k] {
					found = true
					break
				}
			}

			unit, ok := s.getUnit(job, toAdd[k])
			if !ok {
				return ErrFailedQuery
			}

			if len(unit.Users) <= 0 {
				continue
			}

			expiresAt := time.Now().Add(DispatchExpirationTime)
			stmt := tDispatchUnit.
				INSERT(
					tDispatchUnit.DispatchID,
					tDispatchUnit.UnitID,
					tDispatchUnit.ExpiresAt,
				).
				VALUES(
					dsp.Id,
					unit.Id,
					expiresAt,
				).
				ON_DUPLICATE_KEY_UPDATE(
					tDispatchUnit.ExpiresAt.SET(jet.RawTimestamp("VALUES(`expires_at`)")),
				)

			if _, err := stmt.ExecContext(ctx, tx); err != nil {
				if !dbutils.IsDuplicateError(err) {
					return err
				}
			}

			// Only add unit to dispatch if not already assigned
			if !found {
				dsp.Units = append(dsp.Units, &dispatch.DispatchAssignment{
					UnitId:     unit.Id,
					DispatchId: dsp.Id,
					Unit:       unit,
					ExpiresAt:  expiresAtTS,
				})

				if err := s.updateDispatchStatus(ctx, job, dsp, &dispatch.DispatchStatus{
					DispatchId: dsp.Id,
					UnitId:     &unit.Id,
					UserId:     userId,
					Status:     dispatch.StatusDispatch_STATUS_DISPATCH_UNIT_ASSIGNED,
					X:          x,
					Y:          y,
					Postal:     postal,
				}); err != nil {
					return ErrFailedQuery
				}
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return ErrFailedQuery
	}

	data, err := proto.Marshal(dsp)
	if err != nil {
		return err
	}

	for i := 0; i < len(toRemove); i++ {
		s.events.JS.PublishAsync(s.buildSubject(TopicDispatch, TypeDispatchUpdated, job, toRemove[i]), data)
	}
	for i := 0; i < len(toAdd); i++ {
		s.events.JS.PublishAsync(s.buildSubject(TopicDispatch, TypeDispatchUpdated, job, toAdd[i]), data)
	}

	// Dispatch has not units assigned anymore
	if len(dsp.Units) <= 0 {
		// Check dispatch status to not be completed/archived, etc.
		if dsp.Status != nil && (dsp.Status.Status != dispatch.StatusDispatch_STATUS_DISPATCH_ARCHIVED &&
			dsp.Status.Status != dispatch.StatusDispatch_STATUS_DISPATCH_CANCELLED &&
			dsp.Status.Status != dispatch.StatusDispatch_STATUS_DISPATCH_COMPLETED) {
			if err := s.updateDispatchStatus(ctx, job, dsp, &dispatch.DispatchStatus{
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
