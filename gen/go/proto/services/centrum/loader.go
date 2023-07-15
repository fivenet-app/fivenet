package centrum

import (
	"context"

	dispatch "github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"github.com/galexrt/fivenet/pkg/utils/syncx"
	jet "github.com/go-jet/jet/v2/mysql"
)

/* func (s *Server) loadDispatches(ctx context.Context, id uint64) error {
	condition := tDispatchStatus.ID.IS_NULL().OR(
		tDispatchStatus.ID.EQ(
			jet.RawInt(`SELECT MAX(dispatchstatus.id) FROM fivenet_centrum_dispatches_status AS dispatchstatus WHERE dispatchstatus.dispatch_id = dispatch.id`),
		),
	)

	if id > 0 {
		condition = condition.AND(
			tDispatch.ID.EQ(jet.Uint64(id)),
		)
	}

	stmt := tDispatch.
		SELECT(
			tDispatch.ID,
			tDispatch.CreatedAt,
			tDispatch.UpdatedAt,
			tDispatch.Job,
			tDispatch.Message,
			tDispatch.Description,
			tDispatch.Attributes,
			tDispatch.X,
			tDispatch.Y,
			tDispatch.Anon,
			tDispatch.UserID,
			tDispatchStatus.ID,
			tDispatchStatus.CreatedAt,
			tDispatchStatus.UnitID,
			tDispatchStatus.Status,
			tDispatchStatus.Reason,
			tDispatchStatus.Code,
			tDispatchStatus.UserID,
			tDispatchUnit.UnitID,
			tDispatchUnit.DispatchID,
			tDispatchUnit.CreatedAt,
			tDispatchUnit.ExpiresAt,
		).
		FROM(
			tDispatch.
				LEFT_JOIN(tDispatchStatus,
					tDispatchStatus.DispatchID.EQ(tDispatch.ID),
				).
				LEFT_JOIN(tDispatchUnit,
					tDispatchUnit.DispatchID.EQ(tDispatch.ID),
				),
		).
		WHERE(condition).
		ORDER_BY(
			tDispatch.ID.ASC(),
		).LIMIT(120)

	dispatches := []*dispatch.Dispatch{}
	if err := stmt.QueryContext(ctx, s.db, &dispatches); err != nil {
		return err
	}

	for i := 0; i < len(dispatches); i++ {
		// Add units to the dispatch based on the unit assignments
		for k := 0; k < len(dispatches[i].Units); k++ {
			unit, ok := s.getUnit(ctx, &userinfo.UserInfo{
				Job: dispatches[i].Job,
			}, dispatches[i].Units[k].UnitId)
			if !ok {
				return ErrFailedQuery
			}

			dispatches[i].Units[k].Unit = unit
		}
	}

	return nil
} */

func (s *Server) loadUnits(ctx context.Context, id uint64) error {
	condition := tUnitStatus.ID.IS_NULL().OR(
		tUnitStatus.ID.EQ(
			jet.RawInt(`SELECT MAX(unitstatus.id) FROM fivenet_centrum_units_status AS unitstatus WHERE unitstatus.unit_id = unit.id`),
		),
	)

	if id > 0 {
		condition = condition.AND(
			tUnits.ID.EQ(jet.Uint64(id)),
		)
	}

	stmt := tUnits.
		SELECT(
			tUnits.ID,
			tUnits.Job,
			tUnits.Name,
			tUnits.Initials,
			tUnits.Color,
			tUnits.Description,
			tUnitStatus.ID,
			tUnitStatus.CreatedAt,
			tUnitStatus.UnitID,
			tUnitStatus.Status,
			tUnitStatus.Reason,
			tUnitStatus.Code,
			tUnitStatus.UserID,
			tUnitStatus.X,
			tUnitStatus.Y,
			tUnitUser.UnitID,
			tUnitUser.UserID,
			tUnitUser.Identifier,
		).
		FROM(
			tUnits.
				LEFT_JOIN(tUnitStatus,
					tUnitStatus.UnitID.EQ(tUnits.ID),
				).
				LEFT_JOIN(tUnitUser,
					tUnitUser.UnitID.EQ(tUnits.ID),
				),
		).
		WHERE(condition).
		ORDER_BY(
			tUnits.Job.ASC(),
			tUnits.Name.ASC(),
			tUnitStatus.Status.ASC(),
		)

	units := []*dispatch.Unit{}
	if err := stmt.QueryContext(ctx, s.db, &units); err != nil {
		return err
	}

	for i := 0; i < len(units); i++ {
		var err error
		units[i].Users, err = s.resolveUsersForUnit(ctx, units[i].Users)
		if err != nil {
			return err
		}

		jobUnits, _ := s.units.LoadOrStore(units[i].Job, &syncx.Map[uint64, *dispatch.Unit]{})
		jobUnits.Store(units[i].Id, units[i])
	}

	return nil
}
