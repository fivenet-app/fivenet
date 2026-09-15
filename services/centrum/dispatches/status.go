package dispatches

import (
	"context"
	"errors"
	"fmt"

	centrumdispatches "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatches"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	eventscentrum "github.com/fivenet-app/fivenet/v2026/services/centrum/events"
	centrumutils "github.com/fivenet-app/fivenet/v2026/services/centrum/utils"
	colleagueshydrator "github.com/fivenet-app/fivenet/v2026/stores/jobs/colleagues/hydrator"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (s *DispatchDB) GetStatusByID(
	ctx context.Context,
	tx qrm.DB,
	id int64,
) (*centrumdispatches.DispatchStatus, error) {
	tDispatchStatus := table.FivenetCentrumDispatchesStatus.AS("dispatch_status")
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
			tDispatchStatus.X,
			tDispatchStatus.Y,
			tDispatchStatus.Postal,
			tDispatchStatus.CreatorJob,
		).
		FROM(tDispatchStatus).
		WHERE(tDispatchStatus.ID.EQ(mysql.Int64(id))).ORDER_BY(tDispatchStatus.ID.DESC()).
		LIMIT(1)

	var status centrumdispatches.DispatchStatus
	if err := stmt.QueryContext(ctx, tx, &status); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if status.GetUserId() > 0 {
		colleague, err := s.colleagueshydrator.GetBasicByUserID(
			ctx,
			tx,
			nil,
			status.GetUserId(),
			colleagueshydrator.ResolveOpts{
				Scope: colleagueshydrator.JobScope{Mode: colleagueshydrator.JobScopePrimary},
			},
		)
		if err != nil {
			return nil, err
		}
		status.User = colleague
	}
	if status.UnitId != nil && status.GetUnitId() > 0 && status.GetUser() != nil {
		unit, err := s.units.Get(ctx, status.GetUnitId())
		if err != nil {
			return nil, err
		}
		status.Unit = unit
	}
	return &status, nil
}

// AddDispatchStatus persists a status row and returns the hydrated row.
func (s *DispatchDB) AddDispatchStatus(
	ctx context.Context,
	tx qrm.DB,
	status *centrumdispatches.DispatchStatus,
) (*centrumdispatches.DispatchStatus, error) {
	statusTable := table.FivenetCentrumDispatchesStatus
	stmt := statusTable.
		INSERT(
			statusTable.CreatedAt, statusTable.DispatchID, statusTable.Status, statusTable.Reason,
			statusTable.Code, statusTable.UnitID, statusTable.UserID, statusTable.X, statusTable.Y,
			statusTable.Postal, statusTable.CreatorJob,
		).
		VALUES(
			status.GetCreatedAt(),
			status.GetDispatchId(), status.GetStatus(), status.Reason, status.Code,
			status.UnitId, status.UserId, status.X, status.Y, status.Postal, status.CreatorJob,
		)

	result, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return s.GetStatusByID(ctx, tx, id)
}

func (s *DispatchDB) UpdateStatus(
	ctx context.Context,
	dspId int64,
	in *centrumdispatches.DispatchStatus,
) (*centrumdispatches.DispatchStatus, error) {
	dsp, err := s.Get(ctx, dspId)
	if err != nil {
		if errors.Is(err, jetstream.ErrKeyNotFound) {
			return nil, fmt.Errorf(
				"dispatch %d projection not found: %w",
				dspId,
				err,
			)
		}

		return nil, err
	}

	if dsp == nil {
		return nil, fmt.Errorf("dispatch %d not found", dspId)
	}

	if dsp.GetStatus() != nil {
		// If the dispatch status is the same and is a status that shouldn't be duplicated, don't update the status again
		if dsp.GetStatus().GetStatus() == in.GetStatus() &&
			(in.GetStatus() == centrumdispatches.StatusDispatch_STATUS_DISPATCH_NEW ||
				in.GetStatus() == centrumdispatches.StatusDispatch_STATUS_DISPATCH_UNASSIGNED) {
			s.logger.Debug(
				"skipping dispatch status update due to being new or same status",
				zap.Int64("dispatch_id", dsp.GetId()),
				zap.String("status", in.GetStatus().String()),
			)
			return dsp.GetStatus(), nil
		}

		// If the dispatch is complete, we ignore any unit unassignments/accepts/declines
		if centrumutils.IsStatusDispatchComplete(dsp.GetStatus().GetStatus()) &&
			(in.GetStatus() == centrumdispatches.StatusDispatch_STATUS_DISPATCH_UNASSIGNED ||
				in.GetStatus() == centrumdispatches.StatusDispatch_STATUS_DISPATCH_UNIT_UNASSIGNED ||
				in.GetStatus() == centrumdispatches.StatusDispatch_STATUS_DISPATCH_UNIT_ACCEPTED ||
				in.GetStatus() == centrumdispatches.StatusDispatch_STATUS_DISPATCH_UNIT_DECLINED) {
			return dsp.GetStatus(), nil
		}
	}

	s.logger.Debug(
		"updating dispatch status",
		zap.Int64("dispatch_id", dspId),
		zap.String("status", in.GetStatus().String()),
	)

	if in.GetUserId() > 0 {
		var err error
		in.User, err = s.colleagueshydrator.GetBasicByUserID(
			ctx,
			s.db,
			nil,
			in.GetUserId(),
			colleagueshydrator.ResolveOpts{
				Scope: colleagueshydrator.JobScope{
					Mode: colleagueshydrator.JobScopeExplicit,
					Job:  dsp.GetFirstJob(),
				},
			},
		)
		if err != nil {
			return nil, err
		}

		if um, ok := s.tracker.GetUserMarkerById(in.GetUserId()); ok {
			in.X = &um.X
			in.Y = &um.Y
			in.Postal = um.Postal
		}
	}

	// Set postal code using coordinates if empty
	if !in.HasPostal() && in.HasX() && in.HasY() {
		if postal, exists := s.postals.Closest(in.GetX(), in.GetY()); exists {
			in.SetPostal(*postal.Code)
		}
	}

	if in.GetCreatedAt() == nil {
		in.CreatedAt = timestamp.Now()
	}

	in, err = s.AddDispatchStatus(ctx, s.db, in)
	if err != nil {
		return nil, err
	}

	if err := s.updateStatusInKV(ctx, in.GetDispatchId(), in); err != nil {
		return nil, err
	}

	if centrumutils.IsStatusDispatchComplete(in.GetStatus()) {
		if err := s.ScheduleCleanup(ctx, in.GetDispatchId()); err != nil {
			s.logger.Warn(
				"failed to schedule completed dispatch cleanup; recovery audit will retry",
				zap.Int64("dispatch_id", in.GetDispatchId()),
				zap.Error(err),
			)
		}
	} else if err := s.CancelScheduledCleanup(ctx, in.GetDispatchId()); err != nil {
		s.logger.Warn(
			"failed to cancel completed dispatch cleanup",
			zap.Int64("dispatch_id", in.GetDispatchId()),
			zap.Error(err),
		)
	}

	if err := s.publishDispatchStatus(ctx, in, dsp.GetJobs().GetJobStrings()); err != nil {
		return nil, fmt.Errorf(
			"failed to publish dispatch status event (message: '%+v'). %w",
			in,
			err,
		)
	}

	return in, nil
}

func (s *DispatchDB) publishDispatchStatus(
	ctx context.Context,
	status *centrumdispatches.DispatchStatus,
	jobs []string,
) error {
	data, err := proto.Marshal(status)
	if err != nil {
		return err
	}
	for _, job := range jobs {
		if _, err := s.js.Publish(ctx, eventscentrum.BuildSubject(
			eventscentrum.TopicDispatch, eventscentrum.TypeDispatchStatus, job,
		), data); err != nil {
			return err
		}
	}
	return nil
}
