package centrum

import (
	"context"

	dispatch "github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/pkg/utils/dbutils"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/proto"
)

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
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	return nil
}

func (s *Server) updateDispatchStatus(ctx context.Context, userInfo *userinfo.UserInfo, dsp *dispatch.Dispatch, in *dispatch.DispatchStatus) error {
	tDispatchStatus := table.FivenetCentrumDispatchesStatus
	stmt := tDispatchStatus.
		INSERT(
			tDispatchStatus.DispatchID,
			tDispatchStatus.UnitID,
			tDispatchStatus.Status,
			tDispatchStatus.Reason,
			tDispatchStatus.Code,
			tDispatchStatus.UserID,
		).
		VALUES(
			in.DispatchId,
			in.UnitId,
			in.Status,
			in.Reason,
			in.Code,
			in.UserId,
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

	data, err := proto.Marshal(status)
	if err != nil {
		return err
	}

	if len(dsp.Units) == 0 {
		s.events.JS.Publish(s.buildSubject(TopicDispatch, TypeDispatchStatus, userInfo.Job, 0), data)
	} else {
		for _, u := range dsp.Units {
			s.events.JS.Publish(s.buildSubject(TopicDispatch, TypeDispatchStatus, userInfo.Job, u.UnitId), data)
		}
	}

	return nil
}

func (s *Server) updateDispatchAssignments(ctx context.Context, userInfo *userinfo.UserInfo, dsp *dispatch.Dispatch, toAdd []uint64, toRemove []uint64) error {
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
					dsp.Units = utils.RemoveFromSlice(dsp.Units, i)

					if err := s.updateDispatchStatus(ctx, userInfo, dsp, &dispatch.DispatchStatus{
						DispatchId: dsp.Id,
						UnitId:     &toRemove[k],
						UserId:     &userInfo.UserId,
						Status:     dispatch.DISPATCH_STATUS_UNIT_UNASSIGNED,
					}); err != nil {
						return ErrFailedQuery
					}

					continue
				}
			}
		}
	}

	if len(toAdd) > 0 {
		addIds := make([]jet.Expression, len(toAdd))
		for i := 0; i < len(toAdd); i++ {
			addIds[i] = jet.Uint64(toAdd[i])
		}

		for _, id := range addIds {
			stmt := tDispatchUnit.
				INSERT(
					tDispatchUnit.DispatchID,
					tDispatchUnit.UnitID,
					tDispatchUnit.ExpiresAt,
				).
				VALUES(
					dsp.Id,
					id,
					jet.CURRENT_TIMESTAMP().ADD(jet.INTERVAL(16, jet.SECOND)),
				).
				ON_DUPLICATE_KEY_UPDATE(
					tDispatchUnit.ExpiresAt.SET(jet.RawTimestamp("VALUES(`expires_at`)")),
				)

			if _, err := stmt.ExecContext(ctx, tx); err != nil {
				if !dbutils.IsDuplicateError(err) {
					return err
				}
			}
		}

		needsUpdate := []uint64{}
		for k := 0; k < len(toAdd); k++ {
			found := false
			for i := 0; i < len(dsp.Units); i++ {
				if dsp.Units[i].UnitId == toAdd[k] {
					found = true
					break
				}
			}

			unit, ok := s.getUnit(userInfo.Job, toAdd[k])
			if !ok {
				return ErrFailedQuery
			}

			dsp.Units = append(dsp.Units, &dispatch.DispatchAssignment{
				UnitId:     unit.Id,
				DispatchId: dsp.Id,
				Unit:       unit,
			})

			if !found {
				needsUpdate = append(needsUpdate, unit.Id)
			}
		}

		for _, unitId := range needsUpdate {
			if err := s.updateDispatchStatus(ctx, userInfo, dsp, &dispatch.DispatchStatus{
				DispatchId: dsp.Id,
				UnitId:     &unitId,
				UserId:     &userInfo.UserId,
				Status:     dispatch.DISPATCH_STATUS_UNIT_ASSIGNED,
			}); err != nil {
				return ErrFailedQuery
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
		s.events.JS.Publish(s.buildSubject(TopicDispatch, TypeDispatchUpdated, userInfo.Job, toRemove[i]), data)
	}
	for i := 0; i < len(toAdd); i++ {
		s.events.JS.Publish(s.buildSubject(TopicDispatch, TypeDispatchUpdated, userInfo.Job, toAdd[i]), data)
	}

	// Unit is empty, set unit status to be unavailable automatically
	if len(dsp.Units) <= 0 {
		if err := s.updateDispatchStatus(ctx, userInfo, dsp, &dispatch.DispatchStatus{
			DispatchId: dsp.Id,
			Status:     dispatch.DISPATCH_STATUS_UNASSIGNED,
			UserId:     &userInfo.UserId,
		}); err != nil {
			return err
		}
	}

	return nil
}
