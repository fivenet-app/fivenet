package sync

import (
	"context"
	"fmt"

	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"go.uber.org/zap"
)

func (s *Server) handleUserLocations(
	ctx context.Context,
	data *pbsync.SendDataRequest_UserLocations,
) (int64, error) {
	tLocations := table.FivenetCentrumUserLocations

	// Handle clear all
	if data.UserLocations.ClearAll != nil && data.UserLocations.GetClearAll() {
		stmt := tLocations.
			DELETE().
			WHERE(tLocations.UserID.IS_NOT_NULL().OR(tLocations.UserID.IS_NULL()))

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return 0, fmt.Errorf("failed to execute user locations clear all statement. %w", err)
		}
	}

	stmt := tLocations.
		INSERT(
			tLocations.UserID,
			tLocations.Job,
			tLocations.JobGrade,
			tLocations.X,
			tLocations.Y,
			tLocations.Hidden,
		)

	atLeastOne := false
	toDelete := []int32{}
	for _, location := range data.UserLocations.GetUsers() {
		// Collect user locations are marked for removal
		if location.GetRemove() {
			toDelete = append(toDelete, location.GetUserId())
			continue
		}

		jg := mysql.NULL
		if location.JobGrade != nil {
			jg = mysql.Int32(location.GetJobGrade())
		}

		stmt = stmt.
			VALUES(
				location.GetUserId(),
				location.GetJob(),
				jg,
				location.GetCoords().GetX(),
				location.GetCoords().GetY(),
				location.GetHidden(),
			)
		atLeastOne = true
	}

	stmt = stmt.
		ON_DUPLICATE_KEY_UPDATE(
			tLocations.Job.SET(mysql.RawString("VALUES(`job`)")),
			tLocations.JobGrade.SET(mysql.RawInt("VALUES(`job_grade`)")),
			tLocations.X.SET(mysql.RawFloat("VALUES(`x`)")),
			tLocations.Y.SET(mysql.RawFloat("VALUES(`y`)")),
			tLocations.Hidden.SET(mysql.RawBool("VALUES(`hidden`)")),
		)

	rowsAffected := int64(0)
	if atLeastOne {
		res, err := stmt.ExecContext(ctx, s.db)
		if err != nil {
			s.logger.Debug(
				"failed to execute user locations insert statement",
				zap.Any("data", data),
				zap.Error(err),
			)
			return 0, fmt.Errorf("failed to execute user locations insert statement. %w", err)
		}

		rowsAffected, err = res.RowsAffected()
		if err != nil {
			return 0, fmt.Errorf(
				"failed to retrieve rows affected for user locations insert. %w",
				err,
			)
		}
	}

	// Delete any user locations that have been marked for removal
	if len(toDelete) > 0 {
		userIds := []mysql.Expression{}
		for _, userId := range toDelete {
			userIds = append(userIds, mysql.Int32(userId))
		}

		delStmt := tLocations.
			DELETE().
			WHERE(tLocations.UserID.IN(userIds...)).
			LIMIT(int64(len(toDelete)))

		res, err := delStmt.ExecContext(ctx, s.db)
		if err != nil {
			return 0, fmt.Errorf("failed to execute user locations delete statement. %w", err)
		}
		rows, err := res.RowsAffected()
		if err != nil {
			return 0, fmt.Errorf(
				"failed to retrieve rows affected for user locations delete. %w",
				err,
			)
		}
		rowsAffected += rows
	}

	return rowsAffected, nil
}
