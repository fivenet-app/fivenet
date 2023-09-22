package jobs

import (
	"context"
	"errors"

	"github.com/galexrt/fivenet/gen/go/proto/resources/jobs"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) addTimeclockEntry(ctx context.Context, userId int32) error {
	stmt := tTimeClock.
		SELECT(
			tTimeClock.UserID,
			tTimeClock.StartTime,
		).
		FROM(tTimeClock).
		WHERE(jet.AND(
			tTimeClock.UserID.EQ(jet.Int32(userId)),
		)).
		ORDER_BY(tTimeClock.Date.DESC()).
		LIMIT(1)

	var dest jobs.TimeclockEntry
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	// If start time is not null, the entry is active, keep using it
	if dest.StartTime != nil {
		return nil
	}

	tTimeClock := table.FivenetJobsTimeclock
	insert := tTimeClock.
		INSERT(
			tTimeClock.Job,
			tTimeClock.UserID,
			tTimeClock.Date,
		).
		VALUES(
			tUser.SELECT(tUser.Job).FROM(tUser).WHERE(tUser.ID.EQ(jet.Int32(userId))),
			userId,
			jet.CURRENT_DATE(),
		).
		ON_DUPLICATE_KEY_UPDATE(
			tTimeClock.StartTime.SET(jet.CURRENT_TIMESTAMP()),
		)

	if _, err := insert.ExecContext(ctx, s.db); err != nil {
		return err
	}

	return nil
}

func (s *Server) endTimeclockEntry(ctx context.Context, userId int32) error {
	stmt := tTimeClock.
		UPDATE(
			tTimeClock.EndTime,
		).
		SET(
			tTimeClock.EndTime.SET(jet.CURRENT_TIMESTAMP()),
		).
		WHERE(jet.AND(
			tTimeClock.UserID.EQ(jet.Int32(userId)),
			tTimeClock.StartTime.IS_NOT_NULL(),
			tTimeClock.EndTime.IS_NULL(),
		))

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return err
	}

	return nil
}
