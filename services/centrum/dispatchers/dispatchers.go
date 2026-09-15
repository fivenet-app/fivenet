package dispatchers

import (
	"context"
	"database/sql"
	"fmt"
	"slices"

	centrumdispatchers "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatchers"
	jobscolleagues "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/colleagues"
	"github.com/fivenet-app/fivenet/v2026/pkg/events"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2026/pkg/nats/store"
	"github.com/fivenet-app/fivenet/v2026/pkg/tracker"
	errorscentrum "github.com/fivenet-app/fivenet/v2026/services/centrum/errors"
	colleagueshydrator "github.com/fivenet-app/fivenet/v2026/stores/jobs/colleagues/hydrator"
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

func (s *DispatchersDB) SetUserState(
	ctx context.Context,
	job string,
	userId int32,
	signon bool,
) error {
	if job == "" || userId <= 0 {
		return fmt.Errorf("invalid dispatcher identity: job=%q user_id=%d", job, userId)
	}

	var dispatcher *jobscolleagues.Colleague
	if signon {
		um, ok := s.tracker.GetUserMarkerById(userId)
		if !ok || um.GetHidden() || um.GetJob() != job {
			return errorscentrum.ErrNotOnDuty
		}

		byUserID, err := s.colleagueHydrator.HydrateByUserID(
			ctx,
			s.db,
			nil,
			[]int32{userId},
			colleagueshydrator.ResolveOpts{
				Scope: colleagueshydrator.JobScope{
					Mode: colleagueshydrator.JobScopeExplicit,
					Job:  job,
				},
			},
		)
		if err != nil {
			return fmt.Errorf("failed to hydrate dispatcher %d for %s. %w", userId, job, err)
		}

		dispatcher = byUserID[userId]
		if dispatcher == nil {
			dispatcher = &jobscolleagues.Colleague{UserId: userId, Job: job}
		} else if dispatcher.GetJob() == "" {
			dispatcher.Job = job
		}
		s.enricher.EnrichJobName(dispatcher)
	}

	return s.store.ComputeUpdate(
		ctx,
		job,
		func(_ string, existing *centrumdispatchers.Dispatchers) (*centrumdispatchers.Dispatchers, bool, error) {
			if existing == nil {
				existing = &centrumdispatchers.Dispatchers{Job: job}
			}

			idx := slices.IndexFunc(
				existing.GetDispatchers(),
				func(candidate *jobscolleagues.Colleague) bool {
					return candidate.GetUserId() == userId
				},
			)
			if signon {
				if idx >= 0 {
					return existing, false, nil
				}

				existing.Dispatchers = append(existing.Dispatchers, dispatcher)
				return existing, true, nil
			}

			if idx < 0 {
				return existing, false, nil
			}

			existing.Dispatchers = slices.Delete(existing.Dispatchers, idx, idx+1)
			return existing, true, nil
		},
	)
}
