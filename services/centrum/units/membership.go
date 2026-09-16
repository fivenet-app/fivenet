package units

import (
	"context"
	"errors"
	"fmt"

	centrumunits "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/units"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

// SyncUserUnitMapping reconciles membership and mappings from tracker duty
// state. Unlike ReconcileUserJobChange, it may repair a matching assignment.
func (s *UnitDB) SyncUserUnitMapping(ctx context.Context, userId int32) error {
	if userId <= 0 {
		return fmt.Errorf("invalid user ID: %d", userId)
	}

	currentMapping, hasMapping, err := s.tracker.GetUserMapping(userId)
	if err != nil {
		return err
	}

	var currentUnitId int64
	if currentMapping != nil && currentMapping.UnitId != nil {
		currentUnitId = currentMapping.GetUnitId()
	}

	unitId, err := s.LoadUnitIDForUserID(ctx, userId)
	if err != nil {
		return err
	}
	if unitId > 0 {
		unit, err := s.Get(ctx, unitId)
		if err != nil {
			return err
		}

		eligible, err := s.IsEligibleUnitMember(ctx, unit.GetJob(), userId)
		if err != nil {
			return err
		}

		// Unit membership is tied to the user's current duty context. Do not
		// retain an old-job or off-duty assignment until the periodic audit runs.
		if !eligible {
			if err := s.RemoveUnitAssignments(
				ctx,
				"",
				&userId,
				unitId,
				[]int32{userId},
			); err != nil {
				return err
			}
			unitId = 0
		}
	}

	var targetUnitId *int64
	if unitId > 0 && s.tracker.IsUserOnDuty(userId) {
		targetUnitId = &unitId
	}

	switch {
	case targetUnitId != nil:
		if err := s.tracker.SetUserMappingForUser(ctx, userId, targetUnitId); err != nil {
			return err
		}
	case s.tracker.IsUserOnDuty(userId):
		if err := s.tracker.UnsetUnitIDForUser(ctx, userId); err != nil {
			return err
		}
	case hasMapping:
		if err := s.tracker.DeleteUserMapping(ctx, userId); err != nil {
			return err
		}
	}

	var errs error
	if currentUnitId > 0 && currentUnitId != unitId {
		errs = errors.Join(errs, s.SyncUnitMembership(ctx, currentUnitId))
	}
	if unitId > 0 {
		errs = errors.Join(errs, s.SyncUnitMembership(ctx, unitId))
	}

	return errs
}

// ReconcileUserJobChange removes a durable unit assignment when it no longer
// matches the active marker's effective duty job. The user-info event is a
// trigger only: a location job override remains authoritative while it is
// active.
func (s *UnitDB) ReconcileUserJobChange(
	ctx context.Context,
	userId int32,
	_ string,
) (bool, error) {
	if userId <= 0 {
		return false, fmt.Errorf("invalid user ID: %d", userId)
	}

	unitId, err := s.LoadUnitIDForUserID(ctx, userId)
	if err != nil {
		return false, err
	}
	if unitId <= 0 {
		return false, nil
	}

	unit, err := s.Get(ctx, unitId)
	if err != nil {
		return false, err
	}
	if unit == nil {
		return false, nil
	}
	eligible, err := s.IsEligibleUnitMember(ctx, unit.GetJob(), userId)
	if err != nil {
		return false, err
	}
	if eligible {
		return false, nil
	}

	if err := s.UpdateUnitAssignments(ctx, "", nil, unitId, nil, []int32{userId}); err != nil {
		return false, err
	}

	return true, nil
}

func (s *UnitDB) SyncUnitMembership(ctx context.Context, unitId int64) error {
	if unitId <= 0 {
		return fmt.Errorf("invalid unit ID: %d", unitId)
	}

	units, err := s.loadUnitsFromDB(ctx, unitId)
	if err != nil {
		return err
	}

	if len(units) == 0 {
		return s.syncMissingUnitMembership(ctx, unitId)
	}

	return s.syncLoadedUnitMembership(ctx, units[0])
}

func (s *UnitDB) syncLoadedUnitMembership(ctx context.Context, unit *centrumunits.Unit) error {
	if unit == nil {
		return nil
	}

	unitId := unit.GetId()
	if err := s.updateInKV(ctx, unitId, unit); err != nil {
		return err
	}

	userIds := make(map[int32]struct{}, len(unit.GetUsers()))
	var errs error
	for _, user := range unit.GetUsers() {
		userId := user.GetUserId()
		if userId <= 0 {
			continue
		}

		eligible, err := s.IsEligibleUnitMember(ctx, unit.GetJob(), userId)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		if !eligible {
			errs = errors.Join(errs, s.clearTrackerMappingForUnit(ctx, userId, unitId))
			continue
		}

		userIds[userId] = struct{}{}
		if err := s.tracker.SetUserMappingForUser(ctx, userId, &unit.Id); err != nil {
			s.logger.Error(
				"failed to set user's unit id",
				zap.Int64("unit_id", unitId),
				zap.Int32("user_id", userId),
				zap.Error(err),
			)
			errs = errors.Join(errs, err)
		}
	}

	return errors.Join(
		errs,
		s.clearStaleTrackerMappingsForUnit(ctx, unitId, userIds),
	)
}

// IsEligibleUnitMember keeps tracker mappings tied to the active marker job.
// A matching marker may contain a location job/grade override.
func (s *UnitDB) IsEligibleUnitMember(ctx context.Context, job string, userId int32) (bool, error) {
	return s.isEligibleJobMember(ctx, s.db, job, userId)
}

func (s *UnitDB) isEligibleJobMember(ctx context.Context, db qrm.DB, job string, userID int32) (bool, error) {
	marker, found := s.tracker.GetUserMarkerById(userID)
	if !found || marker == nil || marker.GetHidden() || !s.tracker.IsUserOnDuty(userID) || marker.GetJob() != job {
		return false, nil
	}

	inJob, err := s.UserInJob(ctx, db, job, userID)
	if err != nil || inJob {
		return inJob, err
	}

	// A matching active marker is the character's trusted effective duty job,
	// including an override from fivenet_centrum_user_locations.
	return true, nil
}

func (s *UnitDB) syncMissingUnitMembership(ctx context.Context, unitId int64) error {
	var errs error
	if err := s.deleteInKV(ctx, unitId); err != nil && !errors.Is(err, jetstream.ErrKeyNotFound) {
		errs = errors.Join(errs, err)
	}

	return errors.Join(
		errs,
		s.clearStaleTrackerMappingsForUnit(ctx, unitId, nil),
	)
}

func (s *UnitDB) clearStaleTrackerMappingsForUnit(
	ctx context.Context,
	unitId int64,
	validUserIds map[int32]struct{},
) error {
	mappings, err := s.tracker.ListUserMappings(ctx)
	if err != nil {
		return err
	}

	var errs error
	for userId, mapping := range mappings {
		if userId <= 0 || mapping == nil || mapping.UnitId == nil || mapping.GetUnitId() != unitId {
			continue
		}

		if _, ok := validUserIds[userId]; ok {
			continue
		}

		errs = errors.Join(errs, s.tracker.DeleteUserMapping(ctx, userId))
	}

	return errs
}

func (s *UnitDB) clearTrackerMappingForUnit(ctx context.Context, userId int32, unitId int64) error {
	mapping, ok, err := s.tracker.GetUserMapping(userId)
	if err != nil {
		return err
	}
	if !ok || mapping == nil || mapping.UnitId == nil || mapping.GetUnitId() != unitId {
		return nil
	}

	return s.tracker.DeleteUserMapping(ctx, userId)
}

func (s *UnitDB) LoadUnitIDForUserID(ctx context.Context, userId int32) (int64, error) {
	tUnitUser := table.FivenetCentrumUnitsUsers.AS("unit_assignment")

	stmt := tUnitUser.
		SELECT(
			tUnitUser.UnitID.AS("unit_id"),
		).
		FROM(tUnitUser).
		WHERE(
			tUnitUser.UserID.EQ(mysql.Int32(userId)),
		).
		LIMIT(1)

	var dest struct {
		UnitID int64
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, err
		}

		return 0, nil
	}

	return dest.UnitID, nil
}

func (s *UnitDB) UserInJob(
	ctx context.Context,
	db qrm.DB,
	job string,
	userID int32,
) (bool, error) {
	if s.jobs == nil {
		return false, errors.New("jobs store unavailable")
	}

	return s.jobs.UserInJob(ctx, db, job, userID)
}
