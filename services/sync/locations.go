package sync

import (
	"context"
	"fmt"
	"time"

	syncdata "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/sync/data"
	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"go.uber.org/zap"
)

func (s *Server) SendUserLocations(
	ctx context.Context,
	req *pbsync.SendUserLocationsRequest,
) (*pbsync.SendDataResponse, error) {
	s.lastSyncedData.Store(time.Now().Unix())

	rowsAffected, err := s.handleUserLocations(ctx, req.GetUsers(), req.GetClearAll())
	if err != nil {
		return nil, fmt.Errorf("failed to handle user locations data. %w", err)
	}

	return &pbsync.SendDataResponse{
		RowsAffected: rowsAffected,
	}, nil
}

func (s *Server) handleUserLocations(
	ctx context.Context,
	users []*syncdata.CitizenLocations,
	clearAll bool,
) (int64, error) {
	tLocations := table.FivenetCentrumUserLocations

	// Handle clear all
	if clearAll {
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
	for _, location := range users {
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
				zap.Bool("clear_all", clearAll),
				zap.Any("users", users),
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
