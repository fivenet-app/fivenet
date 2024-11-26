package centrummanager

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/timestamp"
	users "github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/zap"
)

func (s *Manager) loadData(ctx context.Context) error {
	ctx, span := s.tracer.Start(ctx, "centrum-loaddata")
	defer span.End()

	s.logger.Debug("loading settings")
	if err := s.LoadSettingsFromDB(ctx, ""); err != nil {
		return fmt.Errorf("failed to load centrum settings: %w", err)
	}

	s.logger.Debug("loading disponents")
	if err := s.LoadDisponentsFromDB(ctx, ""); err != nil {
		return fmt.Errorf("failed to load centrum disponents: %w", err)
	}

	s.logger.Debug("loading units")
	if err := s.LoadUnitsFromDB(ctx, 0); err != nil {
		return fmt.Errorf("failed to load centrum units: %w", err)
	}

	s.logger.Debug("loading dispatches")
	if err := s.LoadDispatchesFromDB(ctx, nil); err != nil {
		return fmt.Errorf("failed to load centrum dispatches: %w", err)
	}

	s.logger.Debug("loaded all centrum data")
	return nil
}

func (s *Manager) LoadSettingsFromDB(ctx context.Context, job string) error {
	tCentrumSettings := tCentrumSettings.AS("settings")
	stmt := tCentrumSettings.
		SELECT(
			tCentrumSettings.Job,
			tCentrumSettings.Enabled,
			tCentrumSettings.Mode,
			tCentrumSettings.FallbackMode,
			tCentrumSettings.PredefinedStatus,
			tCentrumSettings.Timings,
		).
		FROM(tCentrumSettings)

	if job != "" {
		stmt = stmt.WHERE(
			tCentrumSettings.Job.EQ(jet.String(job)),
		)
	}

	var dest []*centrum.Settings
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	for _, settings := range dest {
		settings.Default(settings.Job)

		if err := s.State.UpdateSettings(ctx, settings.Job, settings); err != nil {
			return err
		}
	}

	return nil
}

func (s *Manager) LoadDisponentsFromDB(ctx context.Context, job string) error {
	stmt := tCentrumUsers.
		SELECT(
			tCentrumUsers.Job,
			tCentrumUsers.UserID,
			tUsers.ID,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Job,
			tUsers.Dateofbirth,
			tUsers.PhoneNumber,
			tUserProps.Avatar.AS("usershort.avatar"),
		).
		FROM(
			tCentrumUsers.
				INNER_JOIN(tUsers,
					tUsers.ID.EQ(tCentrumUsers.UserID),
				).
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tCentrumUsers.UserID).
						AND(tUsers.Job.EQ(jet.String(job))),
				),
		)

	if job != "" {
		stmt = stmt.WHERE(
			tCentrumUsers.Job.EQ(jet.String(job)),
		)
	}

	var dest []*users.UserShort
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	perJob := map[string][]*users.UserShort{}
	for _, j := range s.appCfg.Get().UserTracker.LivemapJobs {
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

	if job != "" {
		if err := s.UpdateDisponents(ctx, job, perJob[job]); err != nil {
			return err
		}
	} else {
		for job, us := range perJob {
			if err := s.UpdateDisponents(ctx, job, us); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Manager) LoadUnitsFromDB(ctx context.Context, id uint64) error {
	condition := jet.Bool(true)

	if id > 0 {
		condition = condition.AND(
			tUnits.ID.EQ(jet.Uint64(id)),
		)
	}

	stmt := tUnits.
		SELECT(
			tUnits.ID,
			tUnits.CreatedAt,
			tUnits.UpdatedAt,
			tUnits.Job,
			tUnits.Name,
			tUnits.Initials,
			tUnits.Color,
			tUnits.Description,
			tUnits.Attributes,
			tUnits.HomePostal,
			tUnitUser.UnitID,
			tUnitUser.UserID,
		).
		FROM(
			tUnits.
				LEFT_JOIN(tUnitUser,
					tUnitUser.UnitID.EQ(tUnits.ID),
				),
		).
		WHERE(condition).
		ORDER_BY(
			tUnits.Job.ASC(),
			tUnits.Name.ASC(),
		)

	units := []*centrum.Unit{}
	if err := stmt.QueryContext(ctx, s.db, &units); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	for i := 0; i < len(units); i++ {
		access, err := s.ListUnitAccess(ctx, units[i].Id)
		if err != nil {
			return err
		}
		units[i].Access = access

		status, err := s.GetLastUnitStatus(ctx, s.db, units[i].Job, units[i].Id)
		if err != nil {
			return err
		}
		units[i].Status = status

		if err := s.resolveUsersForUnit(ctx, &units[i].Users); err != nil {
			return err
		}

		if err := s.State.UpdateUnit(ctx, units[i].Job, units[i].Id, units[i]); err != nil {
			return err
		}

		for _, user := range units[i].Users {
			if err := s.SetUnitForUser(ctx, user.User.Job, user.UserId, units[i].Id); err != nil {
				s.logger.Error("failed to set user's unit id", zap.Error(err))
			}
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
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, err
		}

		return 0, nil
	}

	return dest.UnitID, nil
}

func (s *Manager) LoadDispatchesFromDB(ctx context.Context, cond jet.BoolExpression) error {
	condition := tDispatchStatus.ID.IS_NULL().OR(
		jet.AND(
			tDispatchStatus.ID.EQ(
				jet.RawInt("SELECT MAX(`dispatchstatus`.`id`) FROM `fivenet_centrum_dispatches_status` AS `dispatchstatus` WHERE `dispatchstatus`.`dispatch_id` = `dispatch`.`id`"),
			).
				// Don't load archived dispatches into cache
				AND(tDispatchStatus.Status.NOT_IN(
					jet.Int16(int16(centrum.StatusDispatch_STATUS_DISPATCH_ARCHIVED)),
					jet.Int16(int16(centrum.StatusDispatch_STATUS_DISPATCH_CANCELLED)),
					jet.Int16(int16(centrum.StatusDispatch_STATUS_DISPATCH_COMPLETED)),
				)),
		),
	)

	if cond != nil {
		condition = condition.AND(cond)
	}

	tUsers := tUsers.AS("user")
	stmt := tDispatch.
		SELECT(
			tDispatch.ID,
			tDispatch.CreatedAt,
			tDispatch.UpdatedAt,
			tDispatch.Job,
			tDispatch.Message,
			tDispatch.Description,
			tDispatch.Attributes,
			tDispatch.References,
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

	dsps := []*centrum.Dispatch{}
	if err := stmt.QueryContext(ctx, s.db, &dsps); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	if len(dsps) == 0 {
		return nil
	}

	publicJobs := s.appCfg.Get().JobInfo.PublicJobs
	for i := 0; i < len(dsps); i++ {
		var err error
		dsps[i].Units, err = s.LoadDispatchAssignments(ctx, dsps[i].Job, dsps[i].Id)
		if err != nil {
			return err
		}

		if dsps[i].CreatorId != nil {
			dsps[i].Creator, err = s.ResolveUserById(ctx, *dsps[i].CreatorId)
			if err != nil {
				return err
			}

			if dsps[i].Creator != nil {
				// Clear dispatch creator's job info if not a visible job
				if !slices.Contains(publicJobs, dsps[i].Creator.Job) {
					dsps[i].Creator.Job = ""
				}
				dsps[i].Creator.JobGrade = 0
			}
		}

		if dsps[i].Postal == nil {
			postal := s.postals.Closest(dsps[i].X, dsps[i].Y)
			dsps[i].Postal = postal.Code

			if _, err := s.UpdateDispatch(ctx, dsps[i].Job, dsps[i].CreatorId, dsps[i], false); err != nil {
				return err
			}
		}

		// Ensure dispatch has a status
		if dsps[i].Status == nil {
			dsps[i].Status, err = s.AddDispatchStatus(ctx, s.db, dsps[i].Job, &centrum.DispatchStatus{
				CreatedAt:  timestamp.Now(),
				DispatchId: dsps[i].Id,
				Status:     centrum.StatusDispatch_STATUS_DISPATCH_NEW,
				Postal:     dsps[i].Postal,
				X:          &dsps[i].X,
				Y:          &dsps[i].Y,
			}, false)
			if err != nil {
				return err
			}
		}

		if _, err := s.UpdateDispatch(ctx, dsps[i].Job, dsps[i].CreatorId, dsps[i], true); err != nil {
			return err
		}
	}

	return nil
}

func (s *Manager) LoadDispatchAssignments(ctx context.Context, job string, dispatchId uint64) ([]*centrum.DispatchAssignment, error) {
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

	dest := []*centrum.DispatchAssignment{}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	// Resolve units based on the dispatch unit assignments
	for i := 0; i < len(dest); i++ {
		unit, err := s.GetUnit(ctx, job, dest[i].UnitId)
		if unit == nil || err != nil {
			return nil, fmt.Errorf("no unit found with id")
		}

		dest[i].Unit = unit
	}

	return dest, nil
}
