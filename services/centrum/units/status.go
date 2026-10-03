package units

import (
	"context"
	"errors"
	"time"

	centrumunits "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/units"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	eventscentrum "github.com/fivenet-app/fivenet/v2026/services/centrum/events"
	colleagueshydrator "github.com/fivenet-app/fivenet/v2026/stores/jobs/colleagues/hydrator"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (s *UnitDB) GetStatusByID(
	ctx context.Context,
	tx qrm.DB,
	id int64,
) (*centrumunits.UnitStatus, error) {
	tUnitStatus := table.FivenetCentrumUnitsStatus.AS("unit_status")
	tColleagueProps := table.FivenetJobColleagueProps.AS("colleague_props")
	tColleague := table.FivenetUser.AS("colleague")
	tUserProps := table.FivenetUserProps.AS("user_props")
	tAvatar := table.FivenetFiles.AS("profile_picture")

	stmt := tUnitStatus.
		SELECT(
			tUnitStatus.ID,
			tUnitStatus.CreatedAt,
			tUnitStatus.UnitID,
			tUnitStatus.Status,
			tUnitStatus.Reason,
			tUnitStatus.Code,
			tUnitStatus.UserID,
			tUnitStatus.CreatorID,
			tUnitStatus.X,
			tUnitStatus.Y,
			tUnitStatus.Postal,
			tUnitStatus.CreatorJob,
			tColleague.ID,
			tColleague.Firstname,
			tColleague.Lastname,
			tColleague.Job,
			tColleague.JobGrade,
			tColleague.Sex,
			tColleague.Dateofbirth,
			tColleague.PhoneNumber,
			tColleagueProps.UserID,
			tColleagueProps.Job,
			tColleagueProps.NamePrefix,
			tColleagueProps.NameSuffix,
			tUserProps.AvatarFileID.AS("colleague.profile_picture_file_id"),
			tAvatar.FilePath.AS("colleague.profile_picture"),
		).
		FROM(
			tUnitStatus.
				LEFT_JOIN(tColleague,
					tColleague.ID.EQ(tUnitStatus.UserID),
				).
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tUnitStatus.UserID),
				).
				LEFT_JOIN(tColleagueProps,
					mysql.AND(
						tColleagueProps.UserID.EQ(tColleague.ID),
						tColleagueProps.Job.EQ(tColleague.Job),
					),
				).
				LEFT_JOIN(tAvatar,
					tAvatar.ID.EQ(tUserProps.AvatarFileID),
				),
		).
		WHERE(
			tUnitStatus.ID.EQ(mysql.Int64(id)),
		).
		ORDER_BY(tUnitStatus.ID.DESC()).
		LIMIT(1)

	var dest centrumunits.UnitStatus
	if err := stmt.QueryContext(ctx, tx, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		} else {
			return nil, nil
		}
	}

	// We can't use the units store to get the unit as we might be in a "locked" update unit call

	return &dest, nil
}

func (s *UnitDB) GetLastStatus(
	ctx context.Context,
	tx qrm.DB,
	unitId int64,
) (*centrumunits.UnitStatus, error) {
	tUnitStatus := table.FivenetCentrumUnitsStatus.AS("unit_status")
	tColleagueProps := table.FivenetJobColleagueProps.AS("colleague_props")
	tColleague := table.FivenetUser.AS("colleague")
	tUserProps := table.FivenetUserProps.AS("user_props")
	tAvatar := table.FivenetFiles.AS("profile_picture")

	stmt := tUnitStatus.
		SELECT(
			tUnitStatus.ID,
			tUnitStatus.CreatedAt,
			tUnitStatus.UnitID,
			tUnitStatus.Status,
			tUnitStatus.Reason,
			tUnitStatus.Code,
			tUnitStatus.UserID,
			tUnitStatus.CreatorID,
			tUnitStatus.X,
			tUnitStatus.Y,
			tUnitStatus.Postal,
			tUnitStatus.CreatorJob,
			tColleague.ID,
			tColleague.Firstname,
			tColleague.Lastname,
			tColleague.Job,
			tColleague.JobGrade,
			tColleague.Sex,
			tColleague.Dateofbirth,
			tColleague.PhoneNumber,
			tColleagueProps.UserID,
			tColleagueProps.Job,
			tColleagueProps.NamePrefix,
			tColleagueProps.NameSuffix,
			tUserProps.AvatarFileID.AS("colleague.profile_picture_file_id"),
			tAvatar.FilePath.AS("colleague.profile_picture"),
		).
		FROM(
			tUnitStatus.
				LEFT_JOIN(tColleague,
					tColleague.ID.EQ(tUnitStatus.UserID),
				).
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tUnitStatus.UserID),
				).
				LEFT_JOIN(tColleagueProps,
					mysql.AND(
						tColleagueProps.UserID.EQ(tColleague.ID),
						tColleagueProps.Job.EQ(tColleague.Job),
					),
				).
				LEFT_JOIN(tAvatar,
					tAvatar.ID.EQ(tUserProps.AvatarFileID),
				),
		).
		WHERE(mysql.AND(
			tUnitStatus.UnitID.EQ(mysql.Int64(unitId)),
			tUnitStatus.Status.NOT_IN(
				mysql.Int32(int32(centrumunits.StatusUnit_STATUS_UNIT_USER_ADDED)),
				mysql.Int32(int32(centrumunits.StatusUnit_STATUS_UNIT_USER_REMOVED)),
			),
		)).
		ORDER_BY(tUnitStatus.ID.DESC()).
		LIMIT(1)

	var dest centrumunits.UnitStatus
	if err := stmt.QueryContext(ctx, tx, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		} else {
			return nil, nil
		}
	}

	return &dest, nil
}

func (s *UnitDB) AddStatus(
	ctx context.Context,
	tx qrm.DB,
	status *centrumunits.UnitStatus,
) (*centrumunits.UnitStatus, error) {
	// AddStatus only persists the status row. Callers that publish status events
	// must do so after their surrounding DB transaction has committed.
	tUnitStatus := table.FivenetCentrumUnitsStatus
	stmt := tUnitStatus.
		INSERT(
			tUnitStatus.CreatedAt,
			tUnitStatus.UnitID,
			tUnitStatus.Status,
			tUnitStatus.Reason,
			tUnitStatus.Code,
			tUnitStatus.UserID,
			tUnitStatus.X,
			tUnitStatus.Y,
			tUnitStatus.Postal,
			tUnitStatus.CreatorID,
			tUnitStatus.CreatorJob,
		).
		VALUES(
			mysql.CURRENT_TIMESTAMP(),
			status.GetUnitId(),
			status.GetStatus(),
			status.Reason,
			status.Code,
			status.UserId,
			status.X,
			status.Y,
			status.Postal,
			status.CreatorId,
			status.CreatorJob,
		)

	res, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return nil, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	newStatus, err := s.GetStatusByID(ctx, tx, lastId)
	if err != nil {
		return nil, err
	}

	return newStatus, nil
}

func (s *UnitDB) publishStatus(
	ctx context.Context,
	status *centrumunits.UnitStatus,
	job string,
) error {
	data, err := proto.Marshal(status)
	if err != nil {
		return err
	}

	if _, err := s.js.Publish(
		ctx,
		eventscentrum.BuildSubject(eventscentrum.TopicUnit, eventscentrum.TypeUnitStatus, job),
		data,
	); err != nil {
		return err
	}

	return nil
}

func (s *UnitDB) UpdateStatus(
	ctx context.Context,
	unitId int64,
	in *centrumunits.UnitStatus,
) (*centrumunits.UnitStatus, bool, error) {
	unit, err := s.Get(ctx, unitId)
	if err != nil {
		return nil, false, err
	}

	// Suppress only a true duplicate. A reason/code or actor change is a new,
	// client-visible status update even when the status enum is unchanged.
	if unit.GetStatus() != nil &&
		sameUnitStatusContent(unit.GetStatus(), in) &&
		(in.GetStatus() == centrumunits.StatusUnit_STATUS_UNIT_ON_BREAK ||
			in.GetStatus() == centrumunits.StatusUnit_STATUS_UNIT_BUSY ||
			in.GetStatus() == centrumunits.StatusUnit_STATUS_UNIT_UNAVAILABLE ||
			in.GetStatus() == centrumunits.StatusUnit_STATUS_UNIT_AVAILABLE) &&
		// Additionally if the status is under 2 minutes disallow the same status update
		(unit.GetStatus().GetCreatedAt() == nil || time.Since(unit.GetStatus().GetCreatedAt().AsTime()) < 2*time.Minute) {
		s.logger.Debug(
			"skipping duplicate unit status update",
			zap.Int64("unit_id", unitId),
			zap.String("status", in.GetStatus().String()),
		)
		return unit.GetStatus(), false, nil
	}

	if unit.GetAttributes() != nil &&
		unit.GetAttributes().Has(centrumunits.UnitAttribute_UNIT_ATTRIBUTE_STATIC) {
		// Only allow a static unit to be set busy, on break or unavailable
		if in.GetStatus() != centrumunits.StatusUnit_STATUS_UNIT_BUSY &&
			in.GetStatus() != centrumunits.StatusUnit_STATUS_UNIT_ON_BREAK &&
			in.GetStatus() != centrumunits.StatusUnit_STATUS_UNIT_UNAVAILABLE {
			return unit.GetStatus(), false, nil
		}
	}

	s.logger.Debug(
		"updating unit status",
		zap.Int64("unit_id", unitId),
		zap.String("status", in.GetStatus().String()),
	)

	if in.UserId != nil {
		var err error
		in.User, err = s.colleagueHydrator.GetBasicByUserID(
			ctx,
			s.db,
			nil,
			in.GetUserId(),
			colleagueshydrator.ResolveOpts{
				Scope: colleagueshydrator.JobScope{
					Mode: colleagueshydrator.JobScopeExplicit,
					Job:  unit.GetJob(),
				},
			},
		)
		if err != nil {
			return nil, false, err
		}

		if um, ok := s.tracker.GetUserMarkerById(in.GetUserId()); ok {
			in.X = &um.X
			in.Y = &um.Y
			in.Postal = um.Postal
		}
	}
	if in.CreatorId != nil {
		// If the creator of the status is the same as the user, no need to query the db
		if in.UserId != nil && in.GetCreatorId() == in.GetUserId() {
			in.SetCreator(in.GetUser())
		} else {
			var err error
			in.Creator, err = s.colleagueHydrator.GetBasicByUserID(
				ctx,
				s.db,
				nil,
				in.GetCreatorId(),
				colleagueshydrator.ResolveOpts{
					Scope: colleagueshydrator.JobScope{
						Mode: colleagueshydrator.JobScopeExplicit,
						Job:  unit.GetJob(),
					},
				},
			)
			if err != nil {
				return nil, false, err
			}
		}
	}

	if in.GetCreatedAt() == nil {
		in.CreatedAt = timestamp.Now()
	}

	tUnitStatus := table.FivenetCentrumUnitsStatus
	stmt := tUnitStatus.
		INSERT(
			tUnitStatus.CreatedAt,
			tUnitStatus.UnitID,
			tUnitStatus.Status,
			tUnitStatus.Reason,
			tUnitStatus.Code,
			tUnitStatus.UserID,
			tUnitStatus.X,
			tUnitStatus.Y,
			tUnitStatus.Postal,
			tUnitStatus.CreatorID,
			tUnitStatus.CreatorJob,
		).
		VALUES(
			in.GetCreatedAt(),
			in.GetUnitId(),
			in.GetStatus(),
			in.Reason,
			in.Code,
			in.UserId,
			in.X,
			in.Y,
			in.Postal,
			in.CreatorId,
			in.CreatorJob,
		)

	res, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		return nil, false, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return nil, false, err
	}
	in.SetId(lastId)

	if err := s.updateStatusInKV(ctx, in.GetUnitId(), in); err != nil {
		return nil, false, err
	}

	if err := s.publishStatus(ctx, in, unit.GetJob()); err != nil {
		return nil, false, err
	}

	return in, true, nil
}
