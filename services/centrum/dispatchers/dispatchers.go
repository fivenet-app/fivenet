package dispatchers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2025/pkg/nats/store"
	"github.com/fivenet-app/fivenet/v2025/pkg/tracker"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorscentrum "github.com/fivenet-app/fivenet/v2025/services/centrum/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type DispatchersDB struct {
	logger *zap.Logger

	db       *sql.DB
	js       *events.JSWrapper
	enricher *mstlystcdata.Enricher
	tracker  tracker.ITracker

	store *store.Store[centrum.Dispatchers, *centrum.Dispatchers]
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger   *zap.Logger
	DB       *sql.DB
	JS       *events.JSWrapper
	Enricher *mstlystcdata.Enricher
	Tracker  tracker.ITracker
}

func New(p Params) *DispatchersDB {
	ctxCancel, cancel := context.WithCancel(context.Background())

	logger := p.Logger.Named("centrum.dispatchers")
	d := &DispatchersDB{
		logger:   logger,
		db:       p.DB,
		js:       p.JS,
		enricher: p.Enricher,
		tracker:  p.Tracker,
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		st, err := store.New[centrum.Dispatchers, *centrum.Dispatchers](ctxCancel, logger, p.JS, "centrum_dispatchers")
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
	tColleagueProps := table.FivenetJobColleagueProps.AS("colleague_props")
	tUserProps := table.FivenetUserProps
	tCentrumDispatchers := table.FivenetCentrumDispatchers
	tColleague := tables.User().AS("colleague")
	tAvatar := table.FivenetFiles.AS("avatar")

	stmt := tCentrumDispatchers.
		SELECT(
			tCentrumDispatchers.Job,
			tCentrumDispatchers.UserID,
			tColleague.ID,
			tColleague.Firstname,
			tColleague.Lastname,
			tColleague.Job,
			tColleague.Dateofbirth,
			tColleague.PhoneNumber,
			tColleagueProps.UserID,
			tColleagueProps.Job,
			tColleagueProps.NamePrefix,
			tColleagueProps.NameSuffix,
			tUserProps.AvatarFileID.AS("colleague.avatar_file_id"),
			tAvatar.FilePath.AS("colleague.avatar"),
		).
		FROM(
			tCentrumDispatchers.
				INNER_JOIN(tColleague,
					tColleague.ID.EQ(tCentrumDispatchers.UserID),
				).
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tCentrumDispatchers.UserID),
				).
				LEFT_JOIN(tColleagueProps,
					tColleagueProps.UserID.EQ(tColleague.ID).
						AND(tColleagueProps.Job.EQ(tColleague.Job)),
				).
				LEFT_JOIN(tAvatar,
					tAvatar.ID.EQ(tUserProps.AvatarFileID),
				),
		)

	if job != "" {
		stmt = stmt.WHERE(
			tCentrumDispatchers.Job.EQ(jet.String(job)),
		)
	}

	var dest []*jobs.Colleague
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return fmt.Errorf("failed to query centrum dispatchers. %w", err)
		}
	}

	perJob := map[string][]*jobs.Colleague{}
	for _, user := range dest {
		if _, ok := perJob[user.Job]; !ok {
			perJob[user.Job] = []*jobs.Colleague{}
		}

		s.enricher.EnrichJobName(user)

		perJob[user.Job] = append(perJob[user.Job], user)
	}

	if job != "" {
		if err := s.updateDispatchersInKV(ctx, job, perJob[job]); err != nil {
			return fmt.Errorf("failed to update dispatchers for specific job. %w", err)
		}
	} else {
		for job, dispatchers := range perJob {
			if err := s.updateDispatchersInKV(ctx, job, dispatchers); err != nil {
				return fmt.Errorf("failed to update dispatchers for all jobs. %w", err)
			}
		}
	}

	return nil
}

func (s *DispatchersDB) SetUserState(ctx context.Context, job string, userId int32, signon bool) error {
	tCentrumDispatchers := table.FivenetCentrumDispatchers

	if signon {
		if um, ok := s.tracker.GetUserMarkerById(userId); !ok || um.Hidden {
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
			WHERE(jet.AND(
				tCentrumDispatchers.Job.EQ(jet.String(job)),
				tCentrumDispatchers.UserID.EQ(jet.Int32(userId)),
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
