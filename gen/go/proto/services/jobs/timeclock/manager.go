package timeclock

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"time"

	"github.com/galexrt/fivenet/gen/go/proto/resources/jobs"
	"github.com/galexrt/fivenet/gen/go/proto/resources/livemap"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/pkg/tracker"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var (
	tTimeClock = table.FivenetJobsTimeclock.AS("timeclock_entry")
	tUser      = table.Users.AS("user")
)

type Manager struct {
	logger *zap.Logger
	wg     sync.WaitGroup

	tracer  trace.Tracer
	db      *sql.DB
	tracker tracker.ITracker
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger  *zap.Logger
	TP      *tracesdk.TracerProvider
	DB      *sql.DB
	Perms   perms.Permissions
	Tracker tracker.ITracker
}

func New(p Params) *Manager {
	ctx, cancel := context.WithCancel(context.Background())

	m := &Manager{
		logger: p.Logger.Named("jobs.timeclock"),
		wg:     sync.WaitGroup{},

		tracer:  p.TP.Tracer("jobs"),
		db:      p.DB,
		tracker: p.Tracker,
	}

	p.LC.Append(fx.StartHook(func(_ context.Context) error {
		m.wg.Add(1)
		go func() {
			defer m.wg.Done()
			m.runTimeclock(ctx)
		}()

		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		m.wg.Wait()

		return nil
	}))

	return m
}

func (s *Manager) runTimeclock(ctx context.Context) {
	userCh := s.tracker.Subscribe()
	defer s.tracker.Unsubscribe(userCh)

	ticker := time.NewTicker(45 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			func() {
				ctx, span := s.tracer.Start(ctx, "jobs-timeclock")
				defer span.End()

				users, err := s.tracker.GetAllActiveUsers()
				if err != nil {
					s.logger.Error("failed to load active users for timeclock entries", zap.Error(err))
					return
				}

				if err := s.addTimeclockEntries(ctx, users); err != nil {
					s.logger.Error("failed to add timeclock entries", zap.Error(err))
				}
			}()

		case event := <-userCh:
			s.logger.Debug("received user changes", zap.Int("added", len(event.Added)), zap.Int("removed", len(event.Removed)))
			func() {
				ctx, span := s.tracer.Start(ctx, "jobs-timeclock")
				defer span.End()

				if err := s.addTimeclockEntries(ctx, event.Added); err != nil {
					s.logger.Error("failed to add timeclock entries", zap.Error(err))
				}

				for _, userInfo := range event.Removed {
					if err := s.endTimeclockEntry(ctx, userInfo.UserId); err != nil {
						s.logger.Error("failed to end timeclock entry", zap.Error(err))
						continue
					}
				}
			}()
		}
	}
}

func (s *Manager) addTimeclockEntries(ctx context.Context, users []*livemap.UserMarker) error {
	for _, userMarker := range users {
		if err := s.addTimeclockEntry(ctx, userMarker.UserId); err != nil {
			s.logger.Error("failed to add timeclock entry", zap.Error(err))
			continue
		}
	}

	return nil
}

func (s *Manager) addTimeclockEntry(ctx context.Context, userId int32) error {
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

func (s *Manager) endTimeclockEntry(ctx context.Context, userId int32) error {
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
