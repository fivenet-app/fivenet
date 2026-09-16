package units

import (
	"context"
	"errors"
	"fmt"
	"slices"

	centrumunits "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/units"
	jobscolleagues "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/colleagues"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	colleagueshydrator "github.com/fivenet-app/fivenet/v2026/stores/jobs/colleagues/hydrator"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/zap"
)

func (s *UnitDB) UpdateUnitAssignments(
	ctx context.Context,
	creatorJob string,
	creatorId *int32,
	unitId int64,
	toAdd []int32,
	toRemove []int32,
) error {
	syncUserIds, err := s.applyUnitAssignmentChanges(
		ctx,
		creatorJob,
		creatorId,
		unitId,
		toAdd,
		toRemove,
	)
	if err != nil {
		if len(syncUserIds) == 0 {
			return err
		}
	}

	var sideEffectErr error
	if err != nil {
		sideEffectErr = errors.Join(sideEffectErr, err)
	}
	for _, userId := range syncUserIds {
		if err := s.SyncUserUnitMapping(ctx, userId); err != nil {
			sideEffectErr = errors.Join(sideEffectErr, err)
		}
	}

	if len(syncUserIds) == 0 {
		if err := s.SyncUnitMembership(ctx, unitId); err != nil {
			sideEffectErr = errors.Join(sideEffectErr, err)
		}
	}

	return sideEffectErr
}

// RemoveUnitAssignments removes user membership rows and updates unit cache/status only.
// It does not create, clear, or delete tracker mappings; callers must use tracker
// UnsetUnitIDForUser or DeleteUserMapping explicitly for the intended tracker lifecycle.
func (s *UnitDB) RemoveUnitAssignments(
	ctx context.Context,
	creatorJob string,
	creatorId *int32,
	unitId int64,
	userIds []int32,
) error {
	_, err := s.applyUnitAssignmentChanges(
		ctx,
		creatorJob,
		creatorId,
		unitId,
		nil,
		userIds,
	)
	return err
}

func (s *UnitDB) applyUnitAssignmentChanges(
	ctx context.Context,
	_ string,
	creatorId *int32,
	unitId int64,
	toAdd []int32,
	toRemove []int32,
) ([]int32, error) {
	s.logger.Debug(
		"updating unit assignments",
		zap.Int64("unit_id", unitId),
		zap.Int32s("toAdd", toAdd),
		zap.Int32s("toRemove", toRemove),
	)

	if len(toAdd) == 0 && len(toRemove) == 0 {
		return nil, nil
	}

	type pendingStatus struct {
		status *centrumunits.UnitStatus
		job    string
	}

	var x, y *float64
	var postal *string
	if creatorId != nil {
		if um, ok := s.tracker.GetUserMarkerById(*creatorId); ok {
			x = &um.X
			y = &um.Y
			postal = um.Postal
		}
	}

	tUnitUser := table.FivenetCentrumUnitsUsers

	addIds := []int32{}
	pendingStatuses := []pendingStatus{}
	unit, err := s.Get(ctx, unitId)
	if err != nil {
		return nil, err
	}
	if unit == nil {
		return nil, fmt.Errorf("unit %d not found", unitId)
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	if len(toAdd) > 0 {
		for i := range toAdd {
			um, ok := s.tracker.GetUserMarkerById(toAdd[i])
			if !ok || um.GetHidden() {
				continue
			}

			addIds = append(addIds, toAdd[i])
		}
	}

	existingRows := []struct {
		UserID int32 `alias:"user_id"`
	}{}
	stmt := tUnitUser.
		SELECT(tUnitUser.UserID).
		FROM(tUnitUser).
		WHERE(tUnitUser.UnitID.EQ(mysql.Int64(unitId)))
	if err := stmt.QueryContext(ctx, tx, &existingRows); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	existingUsers := map[int32]struct{}{}
	for _, row := range existingRows {
		existingUsers[row.UserID] = struct{}{}
	}

	actualRemove := make([]int32, 0, len(toRemove))
	for _, userId := range toRemove {
		if _, ok := existingUsers[userId]; ok {
			actualRemove = append(actualRemove, userId)
			delete(existingUsers, userId)
		}
	}
	slices.Sort(actualRemove)
	actualRemove = slices.Compact(actualRemove)

	actualAdd := make([]int32, 0, len(addIds))
	for _, userId := range addIds {
		if _, ok := existingUsers[userId]; ok {
			continue
		}

		actualAdd = append(actualAdd, userId)
	}
	slices.Sort(actualAdd)
	actualAdd = slices.Compact(actualAdd)

	eligibleAdd := make([]int32, 0, len(actualAdd))
	for _, userId := range actualAdd {
		ok, err := s.isEligibleJobMember(ctx, tx, unit.GetJob(), userId)
		if err != nil {
			return nil, err
		}
		if !ok {
			s.logger.Debug(
				"dropping unit assignment for user not in job",
				zap.Int64("unit_id", unitId),
				zap.String("job", unit.GetJob()),
				zap.Int32("user_id", userId),
			)
			continue
		}

		eligibleAdd = append(eligibleAdd, userId)
	}

	for _, userId := range eligibleAdd {
		existingUsers[userId] = struct{}{}
	}

	toAddUsers := []*jobscolleagues.Colleague{}
	if len(eligibleAdd) > 0 {
		var err error
		byUserID, err := s.colleagueHydrator.HydrateByUserID(
			ctx,
			s.db,
			nil,
			eligibleAdd,
			colleagueshydrator.ResolveOpts{
				Scope: colleagueshydrator.JobScope{
					Mode: colleagueshydrator.JobScopeExplicit,
					Job:  unit.GetJob(),
				},
			},
		)
		if err != nil {
			return nil, err
		}
		for _, userId := range eligibleAdd {
			if user, ok := byUserID[userId]; ok {
				toAddUsers = append(toAddUsers, user)
			}
		}
	}

	if len(toRemove) > 0 {
		removeIds := make([]mysql.Expression, len(toRemove))
		for i := range toRemove {
			removeIds[i] = mysql.Int32(toRemove[i])
		}

		stmt := tUnitUser.
			DELETE().
			WHERE(mysql.AND(
				tUnitUser.UnitID.EQ(mysql.Int64(unitId)),
				tUnitUser.UserID.IN(removeIds...),
			)).
			LIMIT(int64(len(removeIds)))

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return nil, err
		}
	}

	if len(eligibleAdd) > 0 {
		stmt := tUnitUser.
			INSERT(
				tUnitUser.UnitID,
				tUnitUser.UserID,
			)

		for _, id := range eligibleAdd {
			stmt = stmt.
				VALUES(
					unitId,
					id,
				)
		}

		stmt = stmt.
			ON_DUPLICATE_KEY_UPDATE(
				tUnitUser.UnitID.SET(mysql.RawInt("VALUES(`unit_id`)")),
			)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			if !dbutils.IsDuplicateError(err) {
				return nil, err
			}
		}
	}

	for _, userId := range actualRemove {
		pendingStatuses = append(pendingStatuses, pendingStatus{
			status: &centrumunits.UnitStatus{
				CreatedAt:  timestamp.Now(),
				UnitId:     unit.GetId(),
				Status:     centrumunits.StatusUnit_STATUS_UNIT_USER_REMOVED,
				UserId:     &userId,
				CreatorId:  creatorId,
				X:          x,
				Y:          y,
				Postal:     postal,
				CreatorJob: new(unit.GetJob()),
			},
			job: unit.GetJob(),
		})
	}

	for _, user := range toAddUsers {
		pendingStatuses = append(pendingStatuses, pendingStatus{
			status: &centrumunits.UnitStatus{
				CreatedAt:  timestamp.Now(),
				UnitId:     unit.GetId(),
				Status:     centrumunits.StatusUnit_STATUS_UNIT_USER_ADDED,
				UserId:     &user.UserId,
				CreatorId:  creatorId,
				X:          x,
				Y:          y,
				Postal:     postal,
				CreatorJob: new(user.GetJob()),
			},
			job: unit.GetJob(),
		})
	}

	if len(existingUsers) == 0 && (len(actualRemove) > 0 || len(eligibleAdd) > 0) {
		pendingStatuses = append(pendingStatuses, pendingStatus{
			status: &centrumunits.UnitStatus{
				CreatedAt:  timestamp.Now(),
				UnitId:     unit.GetId(),
				Status:     centrumunits.StatusUnit_STATUS_UNIT_UNAVAILABLE,
				UserId:     creatorId,
				X:          x,
				Y:          y,
				Postal:     postal,
				CreatorId:  creatorId,
				CreatorJob: new(unit.GetJob()),
			},
			job: unit.GetJob(),
		})
	}

	persistedStatuses := make([]pendingStatus, 0, len(pendingStatuses))
	for i := range pendingStatuses {
		status, err := s.AddStatus(ctx, tx, pendingStatuses[i].status)
		if err != nil {
			return nil, err
		}

		persistedStatuses = append(persistedStatuses, pendingStatus{
			status: status,
			job:    pendingStatuses[i].job,
		})
	}

	// The DB commit is the source-of-truth boundary. Unit KV and tracker
	// mappings are reconciled only after assignment and status rows are durable.
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	var sideEffectErr error
	if err := s.refreshUnitCacheFromDB(ctx, unitId); err != nil {
		sideEffectErr = errors.Join(sideEffectErr, err)
	}

	for i := range persistedStatuses {
		if err := s.publishStatus(
			ctx,
			persistedStatuses[i].status,
			persistedStatuses[i].job,
		); err != nil {
			sideEffectErr = errors.Join(sideEffectErr, err)
		}
	}

	syncUserIds := make([]int32, 0, len(eligibleAdd)+len(toRemove))
	syncUserIds = append(syncUserIds, eligibleAdd...)
	syncUserIds = append(syncUserIds, toRemove...)
	slices.Sort(syncUserIds)
	syncUserIds = slices.Compact(syncUserIds)

	if sideEffectErr != nil {
		return syncUserIds, sideEffectErr
	}

	return syncUserIds, nil
}
