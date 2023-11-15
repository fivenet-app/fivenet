package manager

import (
	"context"
	"errors"
	"fmt"

	dispatch "github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	users "github.com/galexrt/fivenet/gen/go/proto/resources/users"
	"github.com/galexrt/fivenet/pkg/utils"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/paulmach/orb"
)

func (s *Manager) LoadSettings(ctx context.Context, job string) error {
	tCentrumSettings := tCentrumSettings.AS("settings")
	stmt := tCentrumSettings.
		SELECT(
			tCentrumSettings.Job,
			tCentrumSettings.Enabled,
			tCentrumSettings.Mode,
			tCentrumSettings.FallbackMode,
		).
		FROM(tCentrumSettings)

	if job != "" {
		stmt = stmt.WHERE(
			tCentrumSettings.Job.EQ(jet.String(job)),
		)
	}

	var dest []*dispatch.Settings
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return err
	}

	for _, settings := range dest {
		if err := s.State.UpdateSettings(settings.Job, settings); err != nil {
			return err
		}
	}

	return nil
}

func (s *Manager) LoadDisponents(ctx context.Context, job string) error {
	stmt := tCentrumUsers.
		SELECT(
			tCentrumUsers.Job,
			tCentrumUsers.UserID,
			tUsers.ID,
			tUsers.Identifier,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Job,
			tUsers.Dateofbirth,
			tUsers.PhoneNumber,
		).
		FROM(
			tCentrumUsers.
				INNER_JOIN(tUsers,
					tUsers.ID.EQ(tCentrumUsers.UserID),
				),
		)

	if job != "" {
		stmt = stmt.WHERE(
			tCentrumUsers.Job.EQ(jet.String(job)),
		)
	}

	var dest []*users.UserShort
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return err
	}

	perJob := map[string][]*users.UserShort{}
	for _, j := range s.trackedJobs {
		if _, ok := perJob[j]; !ok {
			perJob[j] = []*users.UserShort{}
		}
	}

	for _, user := range dest {
		if _, ok := perJob[user.Job]; !ok {
			perJob[user.Job] = []*users.UserShort{}
		}

		s.enricher.EnrichJobName(user)

		perJob[user.Job] = append(perJob[user.Job], user)
	}

	if len(perJob) == 0 {
		s.UpdateDisponents("", nil)
	} else {
		for job, us := range perJob {
			s.UpdateDisponents(job, us)
		}
	}

	return nil
}

func (s *Manager) LoadUnits(ctx context.Context, id uint64) error {
	condition := tUnitStatus.ID.EQ(
		jet.RawInt("SELECT MAX(`unitstatus`.`id`) FROM `fivenet_centrum_units_status` AS `unitstatus` WHERE `unitstatus`.`unit_id` = `unit`.`id`  AND `unitstatus`.`status` NOT IN (2, 3)"),
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
			tUnitStatus.Postal,
			tUnitUser.UnitID,
			tUnitUser.UserID,
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
		if err := s.resolveUsersForUnit(ctx, &units[i].Users); err != nil {
			return err
		}

		um := s.GetUnitsMap(units[i].Job)
		if unit, loaded := um.LoadOrStore(units[i].Id, units[i]); loaded {
			unit.Merge(units[i])
		}

		for _, user := range units[i].Users {
			s.SetUnitForUser(user.UserId, units[i].Id)
		}
	}

	return nil
}

func (s *Manager) LoadUnitIDForUserID(ctx context.Context, userId int32) (uint64, error) {
	stmt := tUnitUser.
		SELECT(
			tUnitUser.UnitID.AS("unit_id"),
		).
		FROM(tUnitUser).
		WHERE(
			tUnitUser.UserID.EQ(jet.Int32(userId)),
		).
		LIMIT(1)

	var dest struct {
		UnitID uint64
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return 0, err
		}

		return 0, nil
	}

	return dest.UnitID, nil
}

func (s *Manager) LoadDispatches(ctx context.Context, id uint64) error {
	condition := tDispatchStatus.ID.IS_NULL().OR(
		tDispatchStatus.ID.EQ(
			jet.RawInt("SELECT MAX(`dispatchstatus`.`id`) FROM `fivenet_centrum_dispatches_status` AS `dispatchstatus` WHERE `dispatchstatus`.`dispatch_id` = `dispatch`.`id`"),
		).
			// Don't load archived dispatches into cache
			AND(tDispatchStatus.Status.NOT_IN(
				jet.Int16(int16(dispatch.StatusDispatch_STATUS_DISPATCH_ARCHIVED)),
				jet.Int16(int16(dispatch.StatusDispatch_STATUS_DISPATCH_CANCELLED)),
				jet.Int16(int16(dispatch.StatusDispatch_STATUS_DISPATCH_COMPLETED)),
			)),
	)

	if id > 0 {
		condition = condition.AND(
			tDispatch.ID.EQ(jet.Uint64(id)),
		)
	}

	tUsers := tUsers.AS("user_short")
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
			tDispatch.Postal,
			tDispatch.Anon,
			tDispatch.CreatorID,
			tDispatchStatus.ID,
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
			tUsers.ID,
			tUsers.Identifier,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Sex,
			tUsers.Dateofbirth,
			tUsers.PhoneNumber,
		).
		FROM(
			tDispatch.
				LEFT_JOIN(tDispatchStatus,
					tDispatchStatus.DispatchID.EQ(tDispatch.ID),
				).
				LEFT_JOIN(tUsers,
					tUsers.ID.EQ(tDispatchStatus.UserID),
				),
		).
		WHERE(condition).
		ORDER_BY(
			tDispatch.ID.DESC(),
		).
		LIMIT(200)

	dispatches := []*dispatch.Dispatch{}
	if err := stmt.QueryContext(ctx, s.db, &dispatches); err != nil {
		return err
	}

	for i := 0; i < len(dispatches); i++ {
		var err error
		dispatches[i].Units, err = s.LoadDispatchAssignments(ctx, dispatches[i].Job, dispatches[i].Id)
		if err != nil {
			return err
		}

		if dispatches[i].CreatorId != nil {
			dispatches[i].Creator, err = s.ResolveUserById(ctx, *dispatches[i].CreatorId)
			if err != nil {
				return err
			}

			// Clear dispatch creator's job info if not a visible job
			if !utils.InSlice(s.publicJobs, dispatches[i].Creator.Job) {
				dispatches[i].Creator.Job = ""
			}
			dispatches[i].Creator.JobGrade = 0
		}

		if dispatches[i].Postal == nil {
			postal := s.postals.Closest(dispatches[i].X, dispatches[i].Y)
			dispatches[i].Postal = postal.Code

			if err := s.UpdateDispatch(ctx, dispatches[i].Job, dispatches[i].CreatorId, dispatches[i], false); err != nil {
				return err
			}
		}

		dm := s.GetDispatchesMap(dispatches[i].Job)
		if dispatch, loaded := dm.LoadOrStore(dispatches[i].Id, dispatches[i]); loaded {
			if dispatch.X != dispatches[i].X || dispatch.Y != dispatches[i].Y {
				s.State.GetDispatchLocations(dispatch.Job).Remove(dispatch, nil)
			}

			dispatch.Merge(dispatches[i])
		}

		// Ensure dispatch has status new if it currently doesn't have one
		if dispatches[i].Status == nil {
			if err := s.UpdateDispatchStatus(ctx, dispatches[i].Job, dispatches[i], &dispatch.DispatchStatus{
				DispatchId: dispatches[i].Id,
				Status:     dispatch.StatusDispatch_STATUS_DISPATCH_NEW,
			}); err != nil {
				return err
			}
		}

		if !s.State.GetDispatchLocations(dispatches[i].Job).Has(dispatches[i], func(p orb.Pointer) bool {
			return p.(*dispatch.Dispatch).Id == dispatches[i].Id
		}) {
			s.State.GetDispatchLocations(dispatches[i].Job).Add(dispatches[i])
		}
	}

	return nil
}

func (s *Manager) LoadDispatchAssignments(ctx context.Context, job string, dispatchId uint64) ([]*dispatch.DispatchAssignment, error) {
	stmt := tDispatch.
		SELECT(
			tDispatchUnit.DispatchID,
			tDispatchUnit.UnitID,
			tDispatchUnit.CreatedAt,
			tDispatchUnit.ExpiresAt,
		).
		FROM(tDispatchUnit).
		ORDER_BY(
			tDispatchUnit.CreatedAt.ASC(),
		).
		WHERE(
			tDispatchUnit.DispatchID.EQ(jet.Uint64(dispatchId)),
		)

	dest := []*dispatch.DispatchAssignment{}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, err
	}

	// Resolve units based on the dispatch unit assignments
	for i := 0; i < len(dest); i++ {
		unit, ok := s.GetUnit(job, dest[i].UnitId)
		if !ok {
			return nil, fmt.Errorf("no unit found with id")
		}

		dest[i].Unit = unit
	}

	return dest, nil
}
