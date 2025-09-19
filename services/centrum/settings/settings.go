package settings

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"sync"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2025/pkg/nats/store"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type SettingsDB struct {
	logger *zap.Logger

	db       *sql.DB
	js       *events.JSWrapper
	enricher *mstlystcdata.Enricher

	store *store.Store[centrum.Settings, *centrum.Settings]

	publicJobsMu sync.RWMutex
	publicJobs   []*jobs.Job
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger   *zap.Logger
	JS       *events.JSWrapper
	DB       *sql.DB
	Cfg      *config.Config
	Enricher *mstlystcdata.Enricher
}

func New(p Params) *SettingsDB {
	ctxCancel, cancel := context.WithCancel(context.Background())

	logger := p.Logger.Named("centrum.settings")
	d := &SettingsDB{
		logger:   logger,
		db:       p.DB,
		js:       p.JS,
		enricher: p.Enricher,

		publicJobsMu: sync.RWMutex{},
		publicJobs:   []*jobs.Job{},
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		st, err := store.New[centrum.Settings, *centrum.Settings](
			ctxCancel,
			logger,
			p.JS,
			"centrum_settings",
			store.WithOnRemoteUpdatedFn(
				func(ctx context.Context, oldValue *centrum.Settings, newValue *centrum.Settings) (*centrum.Settings, error) {
					if newValue == nil {
						return newValue, nil
					}

					d.publicJobsMu.Lock()
					defer d.publicJobsMu.Unlock()

					if newValue.GetPublic() {
						// If the job is public, add it to the public jobs list
						if !slices.ContainsFunc(d.publicJobs, func(j *jobs.Job) bool {
							return j.GetName() == newValue.GetJob()
						}) {
							if j := d.enricher.GetJobByName(newValue.GetJob()); j != nil {
								d.publicJobs = append(d.publicJobs, j)
							}
						}
					} else if i := slices.IndexFunc(d.publicJobs, func(j *jobs.Job) bool {
						return j.GetName() == newValue.GetJob()
					}); i != -1 {
						// If the job is not public, remove it from the public jobs list
						d.publicJobs = slices.Delete(d.publicJobs, i, 1)
					}

					return newValue, nil
				},
			),
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

func (s *SettingsDB) LoadFromDB(ctx context.Context, job string) error {
	tCentrumSettings := table.FivenetCentrumSettings.AS("settings")

	stmt := tCentrumSettings.
		SELECT(
			tCentrumSettings.Job,
			tCentrumSettings.Enabled,
			tCentrumSettings.Type,
			tCentrumSettings.Public,
			tCentrumSettings.Mode,
			tCentrumSettings.FallbackMode,
			tCentrumSettings.PredefinedStatus,
			tCentrumSettings.Timings,
			tCentrumSettings.Configuration,
		).
		FROM(tCentrumSettings)

	if job != "" {
		stmt = stmt.
			WHERE(
				tCentrumSettings.Job.EQ(mysql.String(job)),
			).
			LIMIT(1)
	}

	var dest []*centrum.Settings
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	for _, settings := range dest {
		j := settings.GetJob()
		settings.Default(j)

		var err error
		settings.Access, settings.OfferedAccess, err = s.loadAccess(ctx, j)
		if err != nil {
			return fmt.Errorf("failed to load access for job %s. %w", j, err)
		}

		settings.EffectiveAccess = s.calculateEffectiveAccess(
			settings.GetJob(),
			settings.GetOfferedAccess(),
		)

		if err := s.updateInKV(ctx, j, settings); err != nil {
			return err
		}
	}

	return nil
}

func (s *SettingsDB) updateDB(ctx context.Context, job string, settings *centrum.Settings) error {
	tCentrumSettings := table.FivenetCentrumSettings

	stmt := tCentrumSettings.
		INSERT(
			tCentrumSettings.Job,
			tCentrumSettings.Enabled,
			tCentrumSettings.Type,
			tCentrumSettings.Public,
			tCentrumSettings.Mode,
			tCentrumSettings.FallbackMode,
			tCentrumSettings.PredefinedStatus,
			tCentrumSettings.Timings,
			tCentrumSettings.Configuration,
		).
		VALUES(
			job,
			settings.GetEnabled(),
			settings.GetType(),
			settings.GetPublic(),
			settings.GetMode(),
			settings.GetFallbackMode(),
			settings.GetPredefinedStatus(),
			settings.GetTimings(),
			settings.GetConfiguration(),
		).
		ON_DUPLICATE_KEY_UPDATE(
			tCentrumSettings.Enabled.SET(mysql.Bool(settings.GetEnabled())),
			tCentrumSettings.Type.SET(mysql.Int32(int32(settings.GetType()))),
			tCentrumSettings.Public.SET(mysql.Bool(settings.GetPublic())),
			tCentrumSettings.Mode.SET(mysql.Int32(int32(settings.GetMode()))),
			tCentrumSettings.FallbackMode.SET(mysql.Int32(int32(settings.GetFallbackMode()))),
			tCentrumSettings.PredefinedStatus.SET(mysql.StringExp(mysql.Raw("VALUES(`predefined_status`)"))),
			tCentrumSettings.Timings.SET(mysql.StringExp(mysql.Raw("VALUES(`timings`)"))),
			tCentrumSettings.Configuration.SET(mysql.StringExp(mysql.Raw("VALUES(`configuration`)"))),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return err
	}

	current, err := s.Get(ctx, job)
	if err != nil {
		return fmt.Errorf("failed to get current settings for job %s. %w", job, err)
	}

	if err := s.handleAccessChanges(ctx, job, current.Access, settings.Access); err != nil {
		return fmt.Errorf("failed to handle access changes for job %s. %w", job, err)
	}

	if err := s.handleAccessChanges(ctx, job, current.OfferedAccess, settings.OfferedAccess); err != nil {
		return fmt.Errorf("failed to handle offered access changes for job %s. %w", job, err)
	}

	// Load settings from database so they are updated in the "cache"
	if err := s.LoadFromDB(ctx, job); err != nil {
		return err
	}

	return nil
}

// loadAccess loads the given access for a job from the database.
func (s *SettingsDB) loadAccess(
	ctx context.Context,
	job string,
) (*centrum.CentrumAccess, *centrum.CentrumAccess, error) {
	tCentrumJobAccess := table.FivenetCentrumJobAccess.AS("centrum_job_access")

	stmt := tCentrumJobAccess.
		SELECT(
			tCentrumJobAccess.ID,
			tCentrumJobAccess.SourceJob,
			tCentrumJobAccess.Job,
			tCentrumJobAccess.MinimumGrade,
			tCentrumJobAccess.Access,
			tCentrumJobAccess.AcceptedAt.AS("centrum_job_access.accepted_at"),
		).
		FROM(tCentrumJobAccess).
		WHERE(mysql.OR(
			tCentrumJobAccess.SourceJob.EQ(mysql.String(job)),
			tCentrumJobAccess.Job.EQ(mysql.String(job)),
		)).
		LIMIT(25)

	var accesses []*centrum.CentrumJobAccess
	if err := stmt.QueryContext(ctx, s.db, &accesses); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, nil, fmt.Errorf("failed to load offered access for job %s. %w", job, err)
		}
	}

	access := &centrum.CentrumAccess{}
	offered := &centrum.CentrumAccess{}
	for _, v := range accesses {
		if v.GetSourceJob() == job {
			access.Jobs = append(access.Jobs, v)
		} else {
			offered.Jobs = append(offered.Jobs, v)
		}
	}

	return access, offered, nil
}

func (s *SettingsDB) handleAccessChanges(
	ctx context.Context,
	job string,
	current *centrum.CentrumAccess,
	incoming *centrum.CentrumAccess,
) error {
	removed, updated, created, err := s.compareAccess(job, current, incoming)
	if err != nil {
		return fmt.Errorf("failed to compare access for job %s. %w", job, err)
	}

	if len(removed) == 0 && len(updated) == 0 && len(created) == 0 {
		s.logger.Debug("No changes in access settings",
			zap.String("job", job),
		)
		return nil
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction for job %s. %w", job, err)
	}
	defer tx.Rollback()

	// Collect jobs that need a permissions recalculation
	jobs := map[string]any{}

	if len(removed) > 0 {
		var ids []mysql.Expression
		for _, access := range removed {
			ids = append(ids, mysql.Int64(access.GetId()))
			jobs[access.GetJob()] = nil
		}

		stmt := table.FivenetCentrumJobAccess.
			DELETE().
			WHERE(
				table.FivenetCentrumJobAccess.ID.IN(ids...),
			).
			LIMIT(int64(len(removed)))

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return fmt.Errorf(
				"failed to delete removed access for job %s. %w",
				job,
				err,
			)
		}
	}

	for _, access := range updated {
		jobs[access.GetJob()] = nil

		stmt := table.FivenetCentrumJobAccess.
			UPDATE(
				table.FivenetCentrumJobAccess.Job,
				table.FivenetCentrumJobAccess.MinimumGrade,
				table.FivenetCentrumJobAccess.Access,
				table.FivenetCentrumJobAccess.AcceptedAt,
			).
			SET(
				access.GetJob(),
				access.GetMinimumGrade(),
				access.GetAccess(),
				access.GetAcceptedAt(),
			).
			WHERE(
				table.FivenetCentrumJobAccess.ID.EQ(mysql.Int64(access.GetId())),
			).
			LIMIT(1)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return fmt.Errorf(
				"failed to update access for job %s. %w",
				job,
				err,
			)
		}
	}

	for _, access := range created {
		jobs[access.GetJob()] = nil

		stmt := table.FivenetCentrumJobAccess.
			INSERT(
				table.FivenetCentrumJobAccess.SourceJob,
				table.FivenetCentrumJobAccess.Job,
				table.FivenetCentrumJobAccess.MinimumGrade,
				table.FivenetCentrumJobAccess.Access,
			).
			VALUES(
				job,
				access.GetJob(),
				access.GetMinimumGrade(),
				access.GetAccess(),
			)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			if !dbutils.IsDuplicateError(err) {
				return fmt.Errorf(
					"failed to create access for job %s. %w",
					job,
					err,
				)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf(
			"failed to commit transaction for job %s access changes. %w",
			job,
			err,
		)
	}

	// Load settings for jobs that had access removed to trigger perms recalculation
	var errs error
	for j := range jobs {
		if err := s.LoadFromDB(ctx, j); err != nil {
			errs = multierr.Append(
				errs,
				fmt.Errorf(
					"failed to reload settings for job %s after access changes. %w",
					j,
					err,
				),
			)
		}
	}

	if errs != nil {
		return fmt.Errorf("failed to reload settings after access changes. %w", errs)
	}

	return nil
}

func (s *SettingsDB) compareAccess(
	j string,
	currentSettings *centrum.CentrumAccess,
	incoming *centrum.CentrumAccess,
) ([]*centrum.CentrumJobAccess, []*centrum.CentrumJobAccess, []*centrum.CentrumJobAccess, error) {
	currentAccess := make(map[string]*centrum.CentrumJobAccess)
	if currentSettings != nil {
		for _, access := range currentSettings.GetJobs() {
			currentAccess[access.GetJob()] = access
		}
	}

	removed := []*centrum.CentrumJobAccess{}
	updated := []*centrum.CentrumJobAccess{}
	created := []*centrum.CentrumJobAccess{}

	if incoming != nil {
		for _, offered := range incoming.GetJobs() {
			job := offered.GetJob()
			if current, exists := currentAccess[job]; exists {
				// Compare the access details to check for updates
				if current.GetAccess() != offered.GetAccess() ||
					current.GetMinimumGrade() != offered.GetMinimumGrade() ||
					(current.GetAcceptedAt() == nil && offered.GetAcceptedAt() != nil) ||
					(current.GetAcceptedAt() != nil && offered.GetAcceptedAt() == nil) ||
					(current.GetAcceptedAt() != nil && offered.GetAcceptedAt() != nil && current.GetAcceptedAt().AsTime() != offered.GetAcceptedAt().AsTime()) {
					// If the job is different, make sure to "reset" the accepted at dates
					if current.GetJob() != offered.GetJob() {
						if current.AcceptedAt != nil {
							offered.AcceptedAt = timestamp.New(current.AcceptedAt.AsTime())
						} else {
							offered.AcceptedAt = nil // Reset if no accepted at date is currently set
						}
					}
					// Prevent access level being changed if the source job is different
					if current.GetAccess() != offered.GetAccess() &&
						j != offered.GetSourceJob() {
						offered.Access = current.GetAccess()
					}
					updated = append(updated, offered)
				}
				delete(currentAccess, job) // Remove from currentAccess as it's still valid
			} else {
				// If the job is not in currentAccess, it means it's a new entry (it can't be accepted yet)
				offered.AcceptedAt = nil
				created = append(created, offered)
			}
		}
	}

	// Any remaining jobs in currentAccess are removed
	for _, access := range currentAccess {
		removed = append(removed, access)
	}

	return removed, updated, created, nil
}

func (s *SettingsDB) calculateEffectiveAccess(
	job string,
	offeredAccess *centrum.CentrumAccess,
) *centrum.EffectiveAccess {
	effectiveAccess := &centrum.EffectiveAccess{
		Dispatches: &centrum.EffectiveDispatchAccess{},
	}

	// Calculate effective access: intersection of offered and accepted accesses
	// settings.Access: jobs this job offers access to (outbound)
	// settings.AcceptedAccess: jobs this job has accepted access from (inbound)
	// For each job in settings.Access, check if that job offers access to this job
	if offeredAccess != nil {
		for _, offered := range offeredAccess.GetJobs() {
			if offered.GetSourceJob() == job {
				continue // Skip self-references
			}

			if offered.AcceptedAt == nil {
				continue // Skip unaccepted accesses
			}

			effectiveAccess.Dispatches.Jobs = append(
				effectiveAccess.Dispatches.Jobs,
				&centrum.JobAccessEntry{
					Job:    offered.GetSourceJob(),
					Access: offered.GetAccess(),
				},
			)
		}
	}

	return effectiveAccess
}

func (s *SettingsDB) Update(
	ctx context.Context,
	job string,
	in *centrum.Settings,
) (*centrum.Settings, error) {
	current, err := s.Get(ctx, job)
	if err != nil {
		return nil, err
	}

	// Ensure job is set in the settings
	if in.Job == "" {
		in.Job = job
	}

	current.Merge(in)

	current.EffectiveAccess = s.calculateEffectiveAccess(
		current.GetJob(),
		current.GetOfferedAccess(),
	)

	if err := s.updateDB(ctx, job, current); err != nil {
		return nil, err
	}

	return current, nil
}

func (s *SettingsDB) GetPublicJobs() []*jobs.Job {
	s.publicJobsMu.RLock()
	defer s.publicJobsMu.RUnlock()

	publicJobs := make([]*jobs.Job, 0, len(s.publicJobs))
	for _, job := range s.publicJobs {
		publicJobs = append(publicJobs, proto.Clone(job).(*jobs.Job))
	}

	return publicJobs
}
