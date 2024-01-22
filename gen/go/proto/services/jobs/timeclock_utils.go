package jobs

import (
	"context"
	"errors"

	"github.com/galexrt/fivenet/gen/go/proto/resources/jobs"
	"github.com/galexrt/fivenet/pkg/tracker"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/zap"
)

func (s *Server) runTimeclock() {
	userCh := s.tracker.Subscribe()

	for {
		select {
		case <-s.ctx.Done():
			return

		case event := <-userCh:
			func() {
				ctx, span := s.tracer.Start(s.ctx, "jobs-timeclock")
				defer span.End()

				if len(event.Added) > 0 {
					if err := s.addTimeclockEntries(ctx, event.Added); err != nil {
						s.logger.Error("failed to add timeclock entries", zap.Error(err))
					}
				} else {
					if err := s.addTimeclockEntries(ctx, event.Current); err != nil {
						s.logger.Error("failed to add current timeclock entries", zap.Error(err))
					}
				}

				for _, userInfo := range event.Removed {
					if err := s.endTimeclockEntry(ctx, userInfo.UserID); err != nil {
						s.logger.Error("failed to end timeclock entry", zap.Error(err))
						continue
					}
				}
			}()
		}
	}
}

func (s *Server) addTimeclockEntries(ctx context.Context, users []*tracker.UserInfo) error {
	for _, userInfo := range users {
		if err := s.addTimeclockEntry(ctx, userInfo.UserID); err != nil {
			s.logger.Error("failed to add timeclock entry", zap.Error(err))
			continue
		}
	}

	return nil
}

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

	// If start time is not null, the entry is (already) active, keep using it
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

func (s *Server) getTimeclockStats(ctx context.Context, condition jet.BoolExpression) (*jobs.TimeclockStats, error) {
	stmt := tTimeClock.
		SELECT(
			tTimeClock.Job.AS("timeclock_stats.job"),
			jet.SUM(tTimeClock.SpentTime).AS("timeclock_stats.spent_time_sum"),
			jet.AVG(tTimeClock.SpentTime).AS("timeclock_stats.spent_time_avg"),
			jet.MAX(tTimeClock.SpentTime).AS("timeclock_stats.spent_time_max"),
		).
		FROM(tTimeClock).
		WHERE(jet.AND(
			condition,
			tTimeClock.Date.BETWEEN(jet.CURRENT_DATE().SUB(jet.INTERVAL(7, jet.DAY)), jet.CURRENT_TIMESTAMP()),
		))

	var dest jobs.TimeclockStats
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, err
	}

	return &dest, nil
}

func (s *Server) getTimeclockWeeklyStats(ctx context.Context, condition jet.BoolExpression) ([]*jobs.TimeclockWeeklyStats, error) {
	stmt := tTimeClock.
		SELECT(
			jet.RawString("CONCAT(YEAR(timeclock_entry.`date`), ' - ', WEEK(timeclock_entry.`date`)) AS `timeclock_weekly_stats.date`"),
			jet.SUM(tTimeClock.SpentTime).AS("timeclock_weekly_stats.sum"),
			jet.AVG(tTimeClock.SpentTime).AS("timeclock_weekly_stats.avg"),
			jet.MAX(tTimeClock.SpentTime).AS("timeclock_weekly_stats.max"),
		).
		FROM(tTimeClock).
		WHERE(jet.AND(
			condition,
		)).
		GROUP_BY(
			jet.RawString("WEEK(timeclock_entry.`date`)"),
		).
		ORDER_BY(
			jet.RawString("`timeclock_weekly_stats.date` DESC"),
		).
		LIMIT(12)

	var dest []*jobs.TimeclockWeeklyStats
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest, nil
}
