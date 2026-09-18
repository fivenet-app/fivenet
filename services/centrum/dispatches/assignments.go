package dispatches

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"time"

	centrumdispatches "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatches"
	centrumsettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/settings"
	centrumunits "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/units"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorscentrum "github.com/fivenet-app/fivenet/v2026/services/centrum/errors"
	centrumutils "github.com/fivenet-app/fivenet/v2026/services/centrum/utils"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (s *DispatchDB) LoadDispatchAssignments(
	ctx context.Context,
	dispatchID int64,
) ([]*centrumdispatches.DispatchAssignment, error) {
	tDispatchAssignment := table.FivenetCentrumDispatchesAsgmts.AS("dispatch_assignment")
	stmt := tDispatchAssignment.
		SELECT(
			tDispatchAssignment.DispatchID,
			tDispatchAssignment.UnitID,
			tDispatchAssignment.CreatedAt,
			tDispatchAssignment.ExpiresAt,
		).
		FROM(tDispatchAssignment).
		ORDER_BY(tDispatchAssignment.CreatedAt.ASC()).
		WHERE(tDispatchAssignment.DispatchID.EQ(mysql.Int64(dispatchID)))

	dest := []*centrumdispatches.DispatchAssignment{}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil && !errors.Is(err, qrm.ErrNoRows) {
		return nil, err
	}
	for i := range dest {
		unit, err := s.units.Get(ctx, dest[i].GetUnitId())
		if unit == nil || err != nil {
			return nil, fmt.Errorf(
				"no unit found for dispatch id %d with id %d",
				dispatchID,
				dest[i].GetUnitId(),
			)
		}
		dest[i].Unit = unit
	}
	return dest, nil
}

func (s *DispatchDB) UpdateAssignments(
	ctx context.Context,
	creatorJob *string,
	creatorId *int32,
	dspId int64,
	toAdd []int64,
	toRemove []int64,
	expiresAt time.Time,
) error {
	_, err := s.updateAssignments(
		ctx,
		creatorJob,
		creatorId,
		dspId,
		toAdd,
		toRemove,
		expiresAt,
		false,
	)
	return err
}

// UpdateExpiredAssignments removes only assignments that are still expired at
// the time of the mutation, then updates the live projection and statuses.
func (s *DispatchDB) UpdateExpiredAssignments(
	ctx context.Context,
	creatorJob *string,
	dspId int64,
	toRemove []int64,
) (int, error) {
	return s.updateAssignments(ctx, creatorJob, nil, dspId, nil, toRemove, time.Time{}, true)
}

func (s *DispatchDB) updateAssignments(
	ctx context.Context,
	creatorJob *string,
	creatorId *int32,
	dspId int64,
	toAdd []int64,
	toRemove []int64,
	expiresAt time.Time,
	expiredOnly bool,
) (int, error) {
	s.logger.Debug(
		"updating dispatch assignments",
		zap.Int32p("user_id", creatorId),
		zap.Int64("dispatch_id", dspId),
		zap.Int64s("toAdd", toAdd),
		zap.Int64s("toRemove", toRemove),
	)

	if len(toAdd) == 0 && len(toRemove) == 0 {
		return 0, nil
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

	tDispatchUnit := table.FivenetCentrumDispatchesAsgmts

	// If expires at time is not zero
	expiresAtVal := mysql.NULL
	if !expiresAt.IsZero() {
		expiresAtVal = mysql.TimeT(expiresAt)
	}

	resolvedUnits := map[int64]*centrumunits.Unit{}
	if len(toAdd) > 0 {
		for i := range toAdd {
			unit, err := s.units.Get(ctx, toAdd[i])
			if err != nil {
				continue
			}

			// Skip empty units
			if len(unit.GetUsers()) == 0 {
				continue
			}

			// Only add unit to dispatch if not already assigned/in list
			resolvedUnits[toAdd[i]] = unit
		}
	}

	type pendingDispatchStatus struct {
		status *centrumdispatches.DispatchStatus
		jobs   []string
	}
	pendingStatuses := []pendingDispatchStatus{}
	dsp, err := s.Get(ctx, dspId)
	if err != nil {
		return 0, err
	}
	if dsp == nil {
		return 0, fmt.Errorf("dispatch %d not found", dspId)
	}
	jobs := slices.Clone(dsp.GetJobs().GetJobStrings())

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	// The dispatch projection can briefly outlive its database row (for example
	// while a delete event is propagating). Lock the source row before writing
	// assignments so this race becomes a normal not-found result instead of a
	// foreign-key error.
	dispatchExists := struct {
		ID int64 `alias:"id"`
	}{}
	dispatchStmt := table.FivenetCentrumDispatches.
		SELECT(table.FivenetCentrumDispatches.ID).
		FROM(table.FivenetCentrumDispatches).
		WHERE(table.FivenetCentrumDispatches.ID.EQ(mysql.Int64(dspId))).
		FOR(mysql.UPDATE())
	if err := dispatchStmt.QueryContext(ctx, tx, &dispatchExists); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return 0, errorscentrum.ErrDispatchNotFound
		}
		return 0, err
	}

	existingRows := []struct {
		UnitID int64 `alias:"unit_id"`
	}{}
	stmt := tDispatchUnit.
		SELECT(tDispatchUnit.UnitID.AS("unit_id")).
		FROM(tDispatchUnit).
		WHERE(tDispatchUnit.DispatchID.EQ(mysql.Int64(dspId)))
	if expiredOnly {
		stmt = stmt.FOR(mysql.UPDATE())
	}
	if err := stmt.QueryContext(ctx, tx, &existingRows); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, err
		}
	}

	existingUnits := map[int64]struct{}{}
	existingAssignmentIDs := make([]int64, 0, len(existingRows))
	for _, row := range existingRows {
		existingUnits[row.UnitID] = struct{}{}
		existingAssignmentIDs = append(existingAssignmentIDs, row.UnitID)
	}
	expiredUnits := map[int64]struct{}{}
	if expiredOnly {
		expiredIDs, err := selectExpiredAssignmentIDs(ctx, tx, dspId, toRemove)
		if err != nil {
			return 0, err
		}
		for _, unitID := range expiredIDs {
			expiredUnits[unitID] = struct{}{}
		}
	}
	slices.Sort(existingAssignmentIDs)
	s.logger.Debug(
		"loaded persisted dispatch assignments before assignment update",
		zap.Int64("dispatch_id", dspId),
		zap.Int64s("persisted_assignment_unit_ids", existingAssignmentIDs),
	)

	actualRemove := make([]int64, 0, len(toRemove))
	for _, unitId := range toRemove {
		if expiredOnly {
			if _, ok := expiredUnits[unitId]; !ok {
				continue
			}
		}
		if _, ok := existingUnits[unitId]; ok {
			actualRemove = append(actualRemove, unitId)
			delete(existingUnits, unitId)
		}
	}
	slices.Sort(actualRemove)
	actualRemove = slices.Compact(actualRemove)

	// A requested removal must reconcile an assignment that remains only in the
	// KV projection and emit the corresponding status transition.
	effectiveRemove := slices.Clone(actualRemove)
	removeCandidates := toRemove
	if expiredOnly {
		removeCandidates = actualRemove
	}
	for _, unitId := range removeCandidates {
		if slices.ContainsFunc(
			dsp.GetUnits(),
			func(assignment *centrumdispatches.DispatchAssignment) bool {
				return assignment.GetUnitId() == unitId
			},
		) {
			effectiveRemove = append(effectiveRemove, unitId)
		}
	}
	slices.Sort(effectiveRemove)
	effectiveRemove = slices.Compact(effectiveRemove)

	actualAdd := make([]int64, 0, len(resolvedUnits))
	for unitId := range resolvedUnits {
		if _, ok := existingUnits[unitId]; ok {
			continue
		}

		actualAdd = append(actualAdd, unitId)
		existingUnits[unitId] = struct{}{}
	}
	slices.Sort(actualAdd)
	actualAdd = slices.Compact(actualAdd)

	s.logger.Debug(
		"resolved dispatch assignment changes",
		zap.Int64("dispatch_id", dspId),
		zap.Int64s("requested_add", toAdd),
		zap.Int64s("requested_remove", toRemove),
		zap.Int64s("actual_add", actualAdd),
		zap.Int64s("actual_remove", actualRemove),
		zap.Int64s("effective_remove", effectiveRemove),
		zap.Int("remaining_assignments", len(existingUnits)),
		zap.String("current_status", dsp.GetStatus().GetStatus().String()),
	)

	if len(actualRemove) > 0 {
		removeIds := make([]mysql.Expression, len(actualRemove))
		for i := range actualRemove {
			removeIds[i] = mysql.Int64(actualRemove[i])
		}

		removeWhere := mysql.AND(
			tDispatchUnit.DispatchID.EQ(mysql.Int64(dspId)),
			tDispatchUnit.UnitID.IN(removeIds...),
		)
		if expiredOnly {
			removeWhere = mysql.AND(
				removeWhere,
				tDispatchUnit.ExpiresAt.IS_NOT_NULL(),
				tDispatchUnit.ExpiresAt.LT_EQ(
					mysql.CURRENT_TIMESTAMP().SUB(mysql.INTERVAL(2, mysql.SECOND)),
				),
			)
		}

		stmt := tDispatchUnit.
			DELETE().
			WHERE(removeWhere).
			LIMIT(int64(len(actualRemove)))

		result, err := stmt.ExecContext(ctx, tx)
		if err != nil {
			return 0, err
		}
		if expiredOnly {
			deleted, err := result.RowsAffected()
			if err != nil {
				return 0, err
			}
			if deleted != int64(len(actualRemove)) {
				return 0, fmt.Errorf(
					"expired assignment delete race: expected %d rows, deleted %d",
					len(actualRemove),
					deleted,
				)
			}
		}
	}

	if len(resolvedUnits) > 0 {
		unitIDs := make([]int64, 0, len(resolvedUnits))
		for unitId := range resolvedUnits {
			unitIDs = append(unitIDs, unitId)
		}
		if err := s.insertDispatchAssignments(
			ctx,
			tx,
			dspId,
			unitIDs,
			expiresAtVal,
			true,
		); err != nil {
			return 0, err
		}
	}

	for _, unitId := range effectiveRemove {
		pendingStatuses = append(pendingStatuses, pendingDispatchStatus{
			status: &centrumdispatches.DispatchStatus{
				CreatedAt:  timestamp.Now(),
				DispatchId: dspId,
				UnitId:     &unitId,
				Status:     centrumdispatches.StatusDispatch_STATUS_DISPATCH_UNIT_UNASSIGNED,
				UserId:     creatorId,
				X:          x,
				Y:          y,
				Postal:     postal,
				CreatorJob: creatorJob,
			},
			jobs: jobs,
		})
	}

	for _, unitId := range actualAdd {
		pendingStatuses = append(pendingStatuses, pendingDispatchStatus{
			status: &centrumdispatches.DispatchStatus{
				CreatedAt:  timestamp.Now(),
				DispatchId: dspId,
				UnitId:     &unitId,
				UserId:     creatorId,
				Status:     centrumdispatches.StatusDispatch_STATUS_DISPATCH_UNIT_ASSIGNED,
				X:          x,
				Y:          y,
				Postal:     postal,
				CreatorJob: creatorJob,
			},
			jobs: jobs,
		})
	}

	var currentStatusID int64
	if len(existingUnits) == 0 &&
		(len(effectiveRemove) > 0 || len(actualAdd) > 0) &&
		dsp.GetStatus() != nil &&
		!centrumutils.IsStatusDispatchComplete(dsp.GetStatus().GetStatus()) {
		pendingStatuses = append(pendingStatuses, pendingDispatchStatus{
			status: &centrumdispatches.DispatchStatus{
				CreatedAt:  timestamp.Now(),
				DispatchId: dspId,
				Status:     centrumdispatches.StatusDispatch_STATUS_DISPATCH_UNASSIGNED,
				UserId:     creatorId,
				X:          x,
				Y:          y,
				Postal:     postal,
				CreatorJob: creatorJob,
			},
			jobs: jobs,
		})
	}

	s.logger.Debug(
		"prepared dispatch assignment statuses",
		zap.Int64("dispatch_id", dspId),
		zap.Int("status_count", len(pendingStatuses)),
	)

	persistedStatuses := make([]pendingDispatchStatus, 0, len(pendingStatuses))
	for i := range pendingStatuses {
		status, err := s.AddDispatchStatus(ctx, tx, pendingStatuses[i].status)
		if err != nil {
			return 0, err
		}

		if pendingStatuses[i].status.GetStatus() == centrumdispatches.StatusDispatch_STATUS_DISPATCH_UNASSIGNED {
			currentStatusID = status.GetId()
		}

		persistedStatuses = append(persistedStatuses, pendingDispatchStatus{
			status: status,
			jobs:   pendingStatuses[i].jobs,
		})
	}

	// Commit assignment and status rows before updating dispatch KV or publishing events.
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	finalAssignments, err := s.LoadDispatchAssignments(ctx, dspId)
	if err != nil {
		return 0, err
	}
	finalAssignmentIDs := make([]int64, 0, len(finalAssignments))
	for _, assignment := range finalAssignments {
		finalAssignmentIDs = append(finalAssignmentIDs, assignment.GetUnitId())
	}
	s.logger.Debug(
		"loaded assignments after assignment update",
		zap.Int64("dispatch_id", dspId),
		zap.Int64s("persisted_assignment_unit_ids", finalAssignmentIDs),
	)

	key := centrumutils.IdKey(dspId)
	if err := s.store.ComputeUpdate(
		ctx,
		key,
		func(key string, dsp *centrumdispatches.Dispatch) (*centrumdispatches.Dispatch, bool, error) {
			if dsp == nil {
				s.logger.Error(
					"nil dispatch in dispatch assignment update",
					zap.String("key", key),
					zap.Int64("dispatch_id", dspId),
				)
				return dsp, false, nil
			}

			changed := len(dsp.GetUnits()) != len(finalAssignments)
			if !changed {
				for i := range finalAssignments {
					if proto.Equal(dsp.GetUnits()[i], finalAssignments[i]) {
						continue
					}
					changed = true
					break
				}
			}
			dsp.Units = finalAssignments
			if currentStatusID > 0 {
				for _, pending := range persistedStatuses {
					if pending.status.GetId() != currentStatusID {
						continue
					}
					dsp.Status = pending.status
					changed = true
					break
				}
			}

			return dsp, changed, nil
		},
	); err != nil {
		return 0, err
	}

	// Timers are an optimization for prompt expiry notification. The committed
	// assignment rows remain authoritative and the housekeeper cron recovers if
	// a timer operation fails.
	s.updateAssignmentExpirationTimers(ctx, dspId, actualRemove, resolvedUnits, expiresAt)

	for i := range persistedStatuses {
		s.logger.Debug(
			"publishing dispatch assignment status",
			zap.Int64("dispatch_id", dspId),
			zap.Int64("status_id", persistedStatuses[i].status.GetId()),
			zap.String("status", persistedStatuses[i].status.GetStatus().String()),
			zap.Int64p("unit_id", persistedStatuses[i].status.UnitId),
		)
		if err := s.publishDispatchStatus(
			ctx,
			persistedStatuses[i].status,
			persistedStatuses[i].jobs,
		); err != nil {
			return 0, err
		}
	}

	return len(actualRemove), nil
}

func selectExpiredAssignmentIDs(
	ctx context.Context,
	tx *sql.Tx,
	dspID int64,
	unitIDs []int64,
) ([]int64, error) {
	if len(unitIDs) == 0 {
		return nil, nil
	}

	removeIDs := make([]mysql.Expression, len(unitIDs))
	for i := range unitIDs {
		removeIDs[i] = mysql.Int64(unitIDs[i])
	}

	tDispatchAssignment := table.FivenetCentrumDispatchesAsgmts
	stmt := tDispatchAssignment.
		SELECT(tDispatchAssignment.UnitID.AS("unit_id")).
		FROM(tDispatchAssignment).
		WHERE(mysql.AND(
			tDispatchAssignment.DispatchID.EQ(mysql.Int64(dspID)),
			tDispatchAssignment.UnitID.IN(removeIDs...),
			tDispatchAssignment.ExpiresAt.IS_NOT_NULL(),
			tDispatchAssignment.ExpiresAt.LT_EQ(
				mysql.CURRENT_TIMESTAMP().SUB(mysql.INTERVAL(2, mysql.SECOND)),
			),
		)).
		FOR(mysql.UPDATE())

	var rows []struct {
		UnitID int64 `alias:"unit_id"`
	}
	if err := stmt.QueryContext(ctx, tx, &rows); err != nil {
		return nil, err
	}

	result := make([]int64, 0, len(rows))
	for _, row := range rows {
		result = append(result, row.UnitID)
	}
	return result, nil
}

// DeleteExpiredAssignments removes only assignments that have actually
// expired. It is used for dispatches whose live KV projection was already
// removed for archival, so no dispatch projection or status event can be
// updated.
func (s *DispatchDB) DeleteExpiredAssignments(
	ctx context.Context,
	dspID int64,
	unitIDs []int64,
) (int64, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	expiredIDs, err := selectExpiredAssignmentIDs(ctx, tx, dspID, unitIDs)
	if err != nil {
		return 0, err
	}
	if len(expiredIDs) == 0 {
		return 0, tx.Commit()
	}

	removeIDs := make([]mysql.Expression, len(expiredIDs))
	for i := range expiredIDs {
		removeIDs[i] = mysql.Int64(expiredIDs[i])
	}

	tDispatchAssignment := table.FivenetCentrumDispatchesAsgmts
	stmt := tDispatchAssignment.
		DELETE().
		WHERE(mysql.AND(
			tDispatchAssignment.DispatchID.EQ(mysql.Int64(dspID)),
			tDispatchAssignment.UnitID.IN(removeIDs...),
		)).
		LIMIT(int64(len(expiredIDs)))

	result, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return 0, err
	}

	deleted, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return deleted, nil
}

// insertDispatchAssignments persists one or more assignments in the caller's
// transaction. Upsert is used by explicit assignment updates; self-taking
// passes false so a missing unit-specific assignment is handled separately.
func (s *DispatchDB) insertDispatchAssignments(
	ctx context.Context,
	tx *sql.Tx,
	dspId int64,
	unitIDs []int64,
	expiresAt interface{},
	upsert bool,
) error {
	if len(unitIDs) == 0 {
		return nil
	}

	tDispatchUnit := table.FivenetCentrumDispatchesAsgmts
	stmt := tDispatchUnit.
		INSERT(
			tDispatchUnit.DispatchID,
			tDispatchUnit.UnitID,
			tDispatchUnit.ExpiresAt,
		)
	for _, unitId := range unitIDs {
		stmt = stmt.VALUES(dspId, unitId, expiresAt)
	}

	if upsert {
		stmt = stmt.ON_DUPLICATE_KEY_UPDATE(
			tDispatchUnit.ExpiresAt.SET(mysql.RawTimestamp("VALUES(`expires_at`)")),
		)
	}

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		if !upsert || !dbutils.IsDuplicateError(err) {
			return err
		}
	}

	return nil
}

func (s *DispatchDB) updateAssignmentExpirationTimers(
	ctx context.Context,
	dspID int64,
	actualRemove []int64,
	resolvedUnits map[int64]*centrumunits.Unit,
	expiresAt time.Time,
) {
	for _, unitID := range actualRemove {
		if err := s.CancelAssignmentExpiration(ctx, dspID, unitID); err != nil {
			s.logger.Warn(
				"failed to cancel dispatch assignment expiration timer",
				zap.Int64("dispatch_id", dspID),
				zap.Int64("unit_id", unitID),
				zap.Error(err),
			)
		}
	}
	for unitID := range resolvedUnits {
		if expiresAt.IsZero() {
			if err := s.CancelAssignmentExpiration(ctx, dspID, unitID); err != nil {
				s.logger.Warn(
					"failed to cancel forced dispatch assignment expiration timer",
					zap.Int64("dispatch_id", dspID),
					zap.Int64("unit_id", unitID),
					zap.Error(err),
				)
			}
		} else {
			if err := s.ScheduleAssignmentExpiration(ctx, dspID, unitID, expiresAt); err != nil {
				s.logger.Warn(
					"failed to schedule dispatch assignment expiration timer",
					zap.Int64("dispatch_id", dspID),
					zap.Int64("unit_id", unitID),
					zap.Error(err),
				)
			}
		}
	}
}

func (s *DispatchDB) acceptDispatchAssignment(
	ctx context.Context,
	dspId int64,
	unitId int64,
	allowSelfTake bool,
) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}
	defer tx.Rollback()

	// Lock the dispatch row while taking a snapshot of its latest status.
	dispatchTable := table.FivenetCentrumDispatches.AS("dispatch_for_update")
	statusRow := struct {
		Status sql.NullInt32 `alias:"status"`
	}{}
	statusStmt := dispatchTable.
		SELECT(mysql.RawInt(
			"(SELECT `status` FROM `fivenet_centrum_dispatches_status` WHERE `dispatch_id` = `dispatch_for_update`.`id` ORDER BY `id` DESC LIMIT 1)",
		).AS("status")).
		FROM(dispatchTable).
		WHERE(dispatchTable.ID.EQ(mysql.Int64(dspId))).
		FOR(mysql.UPDATE())
	if err := statusStmt.QueryContext(ctx, tx, &statusRow); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return errorscentrum.ErrNotPartOfDispatch
		}
		return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	if !statusRow.Status.Valid {
		return errorscentrum.ErrNotPartOfDispatch
	}
	currentStatus := centrumdispatches.StatusDispatch(statusRow.Status.Int32)
	if centrumutils.IsStatusDispatchComplete(currentStatus) {
		return errorscentrum.ErrDispatchAlreadyCompleted
	}

	tDispatchUnit := table.FivenetCentrumDispatchesAsgmts.AS("dispatch_assignment")
	loadAssignment := func() (*timestamp.Timestamp, error) {
		assignmentStmt := tDispatchUnit.
			SELECT(
				tDispatchUnit.ExpiresAt.AS("expires_at"),
			).
			FROM(tDispatchUnit).
			WHERE(mysql.AND(
				tDispatchUnit.DispatchID.EQ(mysql.Int64(dspId)),
				tDispatchUnit.UnitID.EQ(mysql.Int64(unitId)),
			)).
			FOR(mysql.UPDATE())

		var assignmentRow struct {
			ExpiresAt *timestamp.Timestamp `alias:"expires_at"`
		}
		if err := assignmentStmt.QueryContext(ctx, tx, &assignmentRow); err != nil {
			return assignmentRow.ExpiresAt, err
		}
		return assignmentRow.ExpiresAt, nil
	}

	// This lookup is scoped to the dispatch and unit. ErrNoRows therefore
	// means this unit has no persisted assignment, which is the valid
	// self-take case for an unassigned dispatch.
	expiresAt, assignmentErr := loadAssignment()
	switch {
	case assignmentErr == nil:
		// The assignment row exists. A nil expiry denotes an already-accepted
		// assignment; only a pending assignment needs expiry validation.
		if expiresAt != nil {
			// An offer may be accepted for two seconds after its deadline,
			// but never later. Keep the cutoff in the conditional update so
			// the row lock and mutation use one database-clock decision.
			clearStmt := tDispatchUnit.
				UPDATE().
				SET(tDispatchUnit.ExpiresAt.SET(mysql.TimestampExp(mysql.NULL))).
				WHERE(mysql.AND(
					tDispatchUnit.DispatchID.EQ(mysql.Int64(dspId)),
					tDispatchUnit.UnitID.EQ(mysql.Int64(unitId)),
					tDispatchUnit.ExpiresAt.GT_EQ(
						mysql.CURRENT_TIMESTAMP().SUB(mysql.INTERVAL(2, mysql.SECOND)),
					),
				))
			result, err := clearStmt.ExecContext(ctx, tx)
			if err != nil {
				return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
			}
			if rows, err := result.RowsAffected(); err != nil || rows == 0 {
				return errorscentrum.ErrNotPartOfDispatch
			}
		}

	case errors.Is(assignmentErr, qrm.ErrNoRows):
		// A missing row is valid only for self-taking a dispatch that is
		// currently unassigned. NEW and UNIT_DECLINED both represent an
		// available dispatch to self-assign.
		if !allowSelfTake ||
			(!centrumutils.IsStatusDispatchUnassigned(currentStatus) &&
				currentStatus != centrumdispatches.StatusDispatch_STATUS_DISPATCH_UNIT_DECLINED) {
			return errorscentrum.ErrNotPartOfDispatch
		}

		if err := s.insertDispatchAssignments(
			ctx,
			tx,
			dspId,
			[]int64{unitId},
			mysql.NULL,
			false,
		); err != nil {
			if !dbutils.IsDuplicateError(err) {
				return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
			}

			// Another self-take won the insert race. Re-read the row under
			// the lock and apply the same expiry validation to it.
			expiresAt, err = loadAssignment()
			if err != nil {
				return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
			}
			// The concurrent taker must have accepted the assignment, which
			// clears its expiry. A pending expiry cannot be self-taken here.
			if expiresAt != nil {
				return errorscentrum.ErrNotPartOfDispatch
			}
		}

	default:
		return errswrap.NewError(assignmentErr, errorscentrum.ErrFailedQuery)
	}

	if err := tx.Commit(); err != nil {
		return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	return nil
}

func (s *DispatchDB) TakeDispatch(
	ctx context.Context,
	userJob string,
	userId int32,
	unitId int64,
	resp centrumdispatches.TakeDispatchResp,
	dispatchIds []int64,
) error {
	settings, err := s.settings.Get(ctx, userJob)
	if err != nil {
		return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	// If the dispatch center is in central command mode, units can't self assign dispatches
	if settings.GetMode() == centrumsettings.CentrumMode_CENTRUM_MODE_CENTRAL_COMMAND {
		return errorscentrum.ErrModeForbidsAction
	}

	unit, err := s.units.Get(ctx, unitId)
	if err != nil {
		return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	var x, y *float64
	var postal *string
	if um, ok := s.tracker.GetUserMarkerById(userId); ok {
		x = &um.X
		y = &um.Y
		postal = um.Postal
	}

	tDispatchUnit := table.FivenetCentrumDispatchesAsgmts

	for _, dspId := range dispatchIds {
		var statusToPublish *centrumdispatches.DispatchStatus
		var publishJobs []string
		action := "declined"
		var result sql.Result

		if resp == centrumdispatches.TakeDispatchResp_TAKE_DISPATCH_RESP_ACCEPTED {
			action = "accepted"
			if err := s.acceptDispatchAssignment(
				ctx,
				dspId,
				unit.GetId(),
				settings.GetMode() != centrumsettings.CentrumMode_CENTRUM_MODE_CENTRAL_COMMAND,
			); err != nil {
				return err
			}
		} else {
			stmt := tDispatchUnit.
				DELETE().
				WHERE(mysql.AND(
					tDispatchUnit.DispatchID.EQ(mysql.Int64(dspId)),
					tDispatchUnit.UnitID.EQ(mysql.Int64(unit.GetId())),
				)).
				LIMIT(1)

			result, err = stmt.ExecContext(ctx, s.db)
			if err != nil {
				if !dbutils.IsDuplicateError(err) {
					return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
				}
			}
		}

		rowsAffected := int64(-1)
		if resp != centrumdispatches.TakeDispatchResp_TAKE_DISPATCH_RESP_ACCEPTED && result != nil {
			if rows, err := result.RowsAffected(); err == nil {
				rowsAffected = rows
			}
		}
		assignments, assignmentErr := s.LoadDispatchAssignments(ctx, dspId)
		if assignmentErr != nil {
			return errswrap.NewError(assignmentErr, errorscentrum.ErrFailedQuery)
		} else {
			assignmentIDs := make([]int64, 0, len(assignments))
			assignmentExpired := false
			for _, assignment := range assignments {
				assignmentIDs = append(assignmentIDs, assignment.GetUnitId())
				if resp == centrumdispatches.TakeDispatchResp_TAKE_DISPATCH_RESP_ACCEPTED &&
					assignment.GetUnitId() == unit.GetId() &&
					assignment.GetExpiresAt() != nil &&
					assignment.GetExpiresAt().AsTime().Before(time.Now().Add(-2*time.Second)) {
					assignmentExpired = true
				}
			}
			s.logger.Debug(
				"persisted take dispatch assignment mutation",
				zap.Int64("dispatch_id", dspId),
				zap.Int64("unit_id", unit.GetId()),
				zap.String("action", action),
				zap.Int64("rows_affected", rowsAffected),
				zap.Int64s("persisted_assignment_unit_ids", assignmentIDs),
			)
			if assignmentExpired {
				return errorscentrum.ErrNotPartOfDispatch
			}
		}

		key := centrumutils.IdKey(dspId)
		if err := s.store.ComputeUpdate(
			ctx,
			key,
			func(key string, dsp *centrumdispatches.Dispatch) (*centrumdispatches.Dispatch, bool, error) {
				// If dispatch is nil or completed, disallow to accept the dispatch
				if dsp == nil ||
					(dsp.GetStatus() != nil && centrumutils.IsStatusDispatchComplete(dsp.GetStatus().GetStatus())) {
					return nil, false, errorscentrum.ErrDispatchAlreadyCompleted
				}

				var status centrumdispatches.StatusDispatch

				// Dispatch accepted
				if resp == centrumdispatches.TakeDispatchResp_TAKE_DISPATCH_RESP_ACCEPTED {
					status = centrumdispatches.StatusDispatch_STATUS_DISPATCH_UNIT_ACCEPTED

					found := false
					// Set unit expires at to nil
					for _, ua := range dsp.GetUnits() {
						if ua.GetUnitId() == unit.GetId() {
							found = true
							// A direct assignment has already been accepted. Do not emit
							// another acceptance status or rewrite the projection.
							if ua.GetExpiresAt() == nil {
								return dsp, false, nil
							}
							ua.ExpiresAt = nil
							break
						}
					}

					if !found {
						dsp.Units = append(dsp.Units, &centrumdispatches.DispatchAssignment{
							DispatchId: dsp.GetId(),
							UnitId:     unit.GetId(),
							Unit:       unit,
							CreatedAt:  timestamp.Now(),
						})
					}

					// Set unit to busy when unit accepts a dispatch.
					if unit.GetStatus() == nil ||
						unit.GetStatus().GetStatus() != centrumunits.StatusUnit_STATUS_UNIT_BUSY {
						if _, _, err := s.units.UpdateStatus(
							ctx,
							unit.GetId(),
							&centrumunits.UnitStatus{
								CreatedAt:  timestamp.Now(),
								UnitId:     unit.GetId(),
								Status:     centrumunits.StatusUnit_STATUS_UNIT_BUSY,
								UserId:     &userId,
								CreatorId:  &userId,
								X:          x,
								Y:          y,
								Postal:     postal,
								CreatorJob: &userJob,
							},
						); err != nil {
							return nil, false, err
						}
					}
				} else {
					// Dispatch declined
					status = centrumdispatches.StatusDispatch_STATUS_DISPATCH_UNIT_DECLINED

					// Remove the unit's assignment
					dsp.Units = slices.DeleteFunc(
						dsp.GetUnits(),
						func(in *centrumdispatches.DispatchAssignment) bool {
							return in.GetUnitId() == unit.GetId()
						},
					)
				}

				persistedStatus, err := s.AddDispatchStatus(
					ctx,
					s.db,
					&centrumdispatches.DispatchStatus{
						CreatedAt:  timestamp.Now(),
						DispatchId: dspId,
						Status:     status,
						UnitId:     &unitId,
						UserId:     &userId,
						X:          x,
						Y:          y,
						Postal:     postal,
						CreatorJob: &userJob,
					},
				)
				if err != nil {
					return nil, false, err
				}
				dsp.SetStatus(persistedStatus)
				statusToPublish = persistedStatus
				publishJobs = slices.Clone(dsp.GetJobs().GetJobStrings())

				return dsp, true, nil
			},
		); err != nil {
			// Ignore errors that are "okay" to encounter
			if !errors.Is(err, errorscentrum.ErrDispatchAlreadyCompleted) {
				return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
			}
		}

		if statusToPublish != nil {
			if err := s.publishDispatchStatus(ctx, statusToPublish, publishJobs); err != nil {
				return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
			}
		}

		if err := s.CancelAssignmentExpiration(ctx, dspId, unit.GetId()); err != nil {
			s.logger.Warn(
				"failed to cancel dispatch assignment expiration timer after response",
				zap.Int64("dispatch_id", dspId),
				zap.Int64("unit_id", unit.GetId()),
				zap.String("action", action),
				zap.Error(err),
			)
		}
	}

	return nil
}
