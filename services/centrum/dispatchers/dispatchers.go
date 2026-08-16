package dispatchers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	centrumdispatchers "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatchers"
	jobscolleagues "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/colleagues"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/pkg/events"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2026/pkg/nats/store"
	"github.com/fivenet-app/fivenet/v2026/pkg/tracker"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorscentrum "github.com/fivenet-app/fivenet/v2026/services/centrum/errors"
	colleagueshydrator "github.com/fivenet-app/fivenet/v2026/stores/jobs/colleagues/hydrator"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type DispatchersDB struct {
	logger *zap.Logger

	db                *sql.DB
	js                *events.JSWrapper
	enricher          mstlystcdata.IEnricher
	tracker           tracker.ITracker
	colleagueHydrator colleagueshydrator.IHydrator

	store *store.Store[centrumdispatchers.Dispatchers, *centrumdispatchers.Dispatchers]
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger            *zap.Logger
	DB                *sql.DB
	JS                *events.JSWrapper
	Enricher          mstlystcdata.IEnricher
	Tracker           tracker.ITracker
	ColleagueHydrator colleagueshydrator.IHydrator
}

func New(p Params) *DispatchersDB {
	ctxCancel, cancel := context.WithCancel(context.Background())

	logger := p.Logger.Named("centrum.dispatchers")
	d := &DispatchersDB{
		logger:            logger,
		db:                p.DB,
		js:                p.JS,
		enricher:          p.Enricher,
		tracker:           p.Tracker,
		colleagueHydrator: p.ColleagueHydrator,
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		st, err := store.New[centrumdispatchers.Dispatchers, *centrumdispatchers.Dispatchers](
			ctxCancel,
			logger,
			p.JS,
			"centrum_dispatchers",
		)
		if err != nil {
			return err
		}

		if err := st.Start(ctxCancel, false); err != nil {
			return err
		}
		d.store = st

		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		return nil
	}))

	return d
}

func (s *DispatchersDB) LoadFromDB(ctx context.Context, job string) error {
	tCentrumDispatchers := table.FivenetCentrumDispatchers

	type dispatchersRow struct {
		Job    string `alias:"job"`
		UserID int32  `alias:"user_id"`
	}

	stmt := tCentrumDispatchers.
		SELECT(
			tCentrumDispatchers.Job.AS("job"),
			tCentrumDispatchers.UserID.AS("user_id"),
		).
		FROM(tCentrumDispatchers)

	if job != "" {
		stmt = stmt.WHERE(tCentrumDispatchers.Job.EQ(mysql.String(job)))
	}

	if job != "" {
		stmt = stmt.
			WHERE(
				tCentrumDispatchers.Job.EQ(mysql.String(job)),
			)
	}

	var rows []*dispatchersRow
	if err := stmt.QueryContext(ctx, s.db, &rows); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return fmt.Errorf("failed to query centrum dispatchers. %w", err)
		}
	}

	if len(rows) == 0 {
		return nil
	}

	userIDs := make([]int32, 0, len(rows))
	for _, row := range rows {
		userIDs = append(userIDs, row.UserID)
	}

	dispatchersByUserID := map[int32]*jobscolleagues.Colleague{}
	if len(userIDs) > 0 {
		byUserID, err := s.colleagueHydrator.HydrateByUserID(
			ctx,
			s.db,
			userIDs,
			colleagueshydrator.ResolveOpts{
				PropsJobMode: colleagueshydrator.PropsJobModePrimary,
			},
		)
		if err != nil {
			return fmt.Errorf("failed to hydrate centrum dispatchers. %w", err)
		}
		dispatchersByUserID = byUserID
	}

	perJob := map[string][]*jobscolleagues.Colleague{}
	for _, row := range rows {
		dispatcher := dispatchersByUserID[row.UserID]
		if dispatcher == nil {
			continue
		}

		s.enricher.EnrichJobName(dispatcher)

		perJob[dispatcher.GetJob()] = append(perJob[dispatcher.GetJob()], dispatcher)
	}

	if job != "" {
		if err := s.updateDispatchersInKV(ctx, job, perJob[job]); err != nil {
			return fmt.Errorf("failed to update dispatchers for specific job. %w", err)
		}
		return nil
	}

	for job, dispatchers := range perJob {
		if err := s.updateDispatchersInKV(ctx, job, dispatchers); err != nil {
			return fmt.Errorf("failed to update dispatchers for all jobs. %w", err)
		}
	}

	return nil
}

func (s *DispatchersDB) SetUserState(
	ctx context.Context,
	job string,
	userId int32,
	signon bool,
) error {
	tCentrumDispatchers := table.FivenetCentrumDispatchers

	if signon {
		if um, ok := s.tracker.GetUserMarkerById(userId); !ok || um.GetHidden() {
			return errorscentrum.ErrNotOnDuty
		}

		stmt := tCentrumDispatchers.
			INSERT(
				tCentrumDispatchers.Job,
				tCentrumDispatchers.UserID,
			).
			VALUES(
				job,
				userId,
			)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			if !dbutils.IsDuplicateError(err) {
				return err
			}
		}
	} else {
		stmt := tCentrumDispatchers.
			DELETE().
			WHERE(mysql.AND(
				tCentrumDispatchers.Job.EQ(mysql.String(job)),
				tCentrumDispatchers.UserID.EQ(mysql.Int32(userId)),
			)).
			LIMIT(1)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}
	}

	// Load updated dispatchers into state
	if err := s.LoadFromDB(ctx, job); err != nil {
		return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	return nil
}
