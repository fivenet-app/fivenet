package colleagueshydrator

import (
	"context"
	"database/sql"
	"errors"
	"maps"
	"slices"

	jobscolleagues "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/colleagues"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	permsjobs "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/jobs/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2026/pkg/perms"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	jobsstore "github.com/fivenet-app/fivenet/v2026/stores/jobs"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"stores.jobs.colleagues.hydrator",
	fx.Provide(New),
)

type JobScopeMode int

const (
	JobScopeCaller JobScopeMode = iota
	JobScopePrimary
	JobScopeExplicit
)

type JobScope struct {
	Mode JobScopeMode
	Job  string
}

// ResolveOpts controls job resolution and optional, permission-gated details.
type ResolveOpts struct {
	Scope JobScope

	// IncludeNote requests the colleague note. It is returned only when the
	// requester has permission or is a job admin.
	IncludeNote bool
	// IncludeLabels requests the additional labels lookup. It is returned only
	// when the requester has permission or is a job admin. Keep this opt-in
	// unless labels become part of the standard colleague hydration contract.
	IncludeLabels bool
}

type IHydrator interface {
	GetBasicByUserID(
		ctx context.Context,
		db qrm.DB,
		userInfo *userinfo.UserInfo,
		userID int32,
		opts ResolveOpts,
	) (*jobscolleagues.Colleague, error)
	ListBasicByUserID(
		ctx context.Context,
		db qrm.DB,
		userInfo *userinfo.UserInfo,
		userIDs []int32,
		opts ResolveOpts,
	) ([]*jobscolleagues.Colleague, error)
	HydrateBasicByUserID(
		ctx context.Context,
		db qrm.DB,
		userInfo *userinfo.UserInfo,
		userIDs []int32,
		opts ResolveOpts,
	) (map[int32]*jobscolleagues.Colleague, error)
	HydrateBasicTargets(
		ctx context.Context,
		db qrm.DB,
		userInfo *userinfo.UserInfo,
		targets []BasicTarget,
		opts ResolveOpts,
	) error

	GetByUserID(
		ctx context.Context,
		db qrm.DB,
		userInfo *userinfo.UserInfo,
		userID int32,
		opts ResolveOpts,
	) (*jobscolleagues.Colleague, error)
	ListByUserID(
		ctx context.Context,
		db qrm.DB,
		userInfo *userinfo.UserInfo,
		userIDs []int32,
		opts ResolveOpts,
	) ([]*jobscolleagues.Colleague, error)
	HydrateByUserID(
		ctx context.Context,
		db qrm.DB,
		userInfo *userinfo.UserInfo,
		userIDs []int32,
		opts ResolveOpts,
	) (map[int32]*jobscolleagues.Colleague, error)
	HydrateTargets(
		ctx context.Context,
		db qrm.DB,
		userInfo *userinfo.UserInfo,
		targets []Target,
		opts ResolveOpts,
	) error
}

type Target struct {
	UserID int32
	Set    func(*jobscolleagues.Colleague)
}

type BasicTarget struct {
	UserID int32
	Set    func(*jobscolleagues.Colleague)
}

type Hydrator struct {
	db       *sql.DB
	perms    perms.Permissions
	enricher mstlystcdata.IUserAwareEnricher
	store    jobsstore.IStore
}

type Params struct {
	fx.In

	DB                *sql.DB
	Perms             perms.Permissions
	UserAwareEnricher mstlystcdata.IUserAwareEnricher
	Store             jobsstore.IStore
}

func New(p Params) IHydrator {
	return &Hydrator{
		db:       p.DB,
		perms:    p.Perms,
		enricher: p.UserAwareEnricher,
		store:    p.Store,
	}
}

type resolvedFields struct {
	note   bool
	labels bool
}

func (h *Hydrator) resolveFields(
	userInfo *userinfo.UserInfo,
	opts ResolveOpts,
) (resolvedFields, bool, error) {
	if userInfo == nil {
		return resolvedFields{}, false, nil
	}
	jobAdmin := userInfo.GetJobAdmin()
	if !opts.IncludeNote && !opts.IncludeLabels {
		return resolvedFields{}, jobAdmin, nil
	}

	fields, err := permsjobs.ColleaguesService.GetColleague.TypesTyped.Get(h.perms, userInfo)
	if err != nil {
		return resolvedFields{}, false, err
	}
	return resolvedFields{
		note: opts.IncludeNote &&
			(jobAdmin || fields.Contains(permsjobs.ColleaguesServiceGetColleagueTypesPermValueNote)),
		labels: opts.IncludeLabels &&
			(jobAdmin || fields.Contains(permsjobs.ColleaguesServiceGetColleagueTypesPermValueLabels)),
	}, jobAdmin, nil
}

func (h *Hydrator) ListByUserID(
	ctx context.Context,
	db qrm.DB,
	userInfo *userinfo.UserInfo,
	userIDs []int32,
	opts ResolveOpts,
) ([]*jobscolleagues.Colleague, error) {
	if len(userIDs) == 0 {
		return []*jobscolleagues.Colleague{}, nil
	}
	if db == nil {
		db = h.db
	}

	fields, jobAdmin, err := h.resolveFields(userInfo, opts)
	if err != nil {
		return nil, err
	}

	resolvedUserIDs := uniqueUserIDs(userIDs)

	userJobs, err := h.resolveJob(ctx, db, userInfo, resolvedUserIDs, opts)
	if err != nil {
		return nil, err
	}

	colleaguesByUserID := make(map[int32]*jobscolleagues.Colleague, len(userIDs))
	fallbackByUserID, err := h.loadFallbackByUserID(ctx, db, resolvedUserIDs)
	if err != nil {
		return nil, err
	}
	maps.Copy(colleaguesByUserID, fallbackByUserID)

	enrichJobInfo := h.enricher.EnrichJobInfoSafeFunc(userInfo)
	for _, jobGroup := range userJobs {
		colleagues, err := h.store.ListColleaguesByUserIDs(
			ctx,
			db,
			jobsstore.ListColleaguesByUserIDsQuery{
				Job:         jobGroup.job,
				UserIDs:     jobGroup.userIDs,
				WithColumns: columnsForFields(fields),
			},
		)
		if err != nil {
			return nil, err
		}

		if fields.labels {
			if err := h.attachLabels(ctx, db, jobGroup.job, colleagues, jobAdmin); err != nil {
				return nil, err
			}
		}

		for _, colleague := range colleagues {
			colleaguesByUserID[colleague.GetUserId()] = colleague
		}
	}

	for _, colleague := range colleaguesByUserID {
		enrichJobInfo(colleague)
	}

	out := make([]*jobscolleagues.Colleague, 0, len(userIDs))
	for _, userID := range userIDs {
		if colleague, ok := colleaguesByUserID[userID]; ok {
			out = append(out, colleague)
		}
	}

	return out, nil
}

func (h *Hydrator) loadFallbackByUserID(
	ctx context.Context,
	db qrm.DB,
	userIDs []int32,
) (map[int32]*jobscolleagues.Colleague, error) {
	if len(userIDs) == 0 {
		return map[int32]*jobscolleagues.Colleague{}, nil
	}
	if db == nil {
		db = h.db
	}

	tUser := table.FivenetUser.AS("colleague")
	tUserProps := table.FivenetUserProps.AS("colleague_user_props")
	tAvatar := table.FivenetFiles.AS("colleague_profile_picture")

	userIdExprs := make([]mysql.Expression, 0, len(userIDs))
	for _, userID := range userIDs {
		userIdExprs = append(userIdExprs, mysql.Int32(userID))
	}

	stmt := tUser.
		SELECT(
			tUser.ID,
			tUser.Job.AS("colleague.job"),
			tUser.JobGrade.AS("colleague.job_grade"),
			tUser.Firstname,
			tUser.Lastname,
			tUser.Dateofbirth,
			tUser.PhoneNumber,
			tUserProps.AvatarFileID.AS("colleague.profile_picture_file_id"),
			tAvatar.FilePath.AS("colleague.profile_picture"),
			tUserProps.Email.AS("colleague.email"),
		).
		FROM(
			tUser.
				LEFT_JOIN(tUserProps, tUserProps.UserID.EQ(tUser.ID)).
				LEFT_JOIN(tAvatar, tAvatar.ID.EQ(tUserProps.AvatarFileID)),
		).
		WHERE(mysql.AND(
			tUser.ID.IN(userIdExprs...),
			tUser.DeletedAt.IS_NULL(),
		)).
		LIMIT(int64(len(userIDs)))

	colleagues := []*jobscolleagues.Colleague{}
	if err := stmt.QueryContext(ctx, db, &colleagues); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	out := make(map[int32]*jobscolleagues.Colleague, len(colleagues))
	for _, colleague := range colleagues {
		out[colleague.GetUserId()] = colleague
	}

	return out, nil
}

func (h *Hydrator) GetBasicByUserID(
	ctx context.Context,
	db qrm.DB,
	userInfo *userinfo.UserInfo,
	userID int32,
	opts ResolveOpts,
) (*jobscolleagues.Colleague, error) {
	colleagues, err := h.ListBasicByUserID(ctx, db, userInfo, []int32{userID}, opts)
	if err != nil || len(colleagues) == 0 {
		return nil, err
	}
	return colleagues[0], nil
}

func (h *Hydrator) ListBasicByUserID(
	ctx context.Context,
	db qrm.DB,
	userInfo *userinfo.UserInfo,
	userIDs []int32,
	opts ResolveOpts,
) ([]*jobscolleagues.Colleague, error) {
	opts.IncludeNote = false
	opts.IncludeLabels = false
	return h.ListByUserID(ctx, db, userInfo, userIDs, opts)
}

func (h *Hydrator) HydrateBasicByUserID(
	ctx context.Context,
	db qrm.DB,
	userInfo *userinfo.UserInfo,
	userIDs []int32,
	opts ResolveOpts,
) (map[int32]*jobscolleagues.Colleague, error) {
	colleagues, err := h.ListBasicByUserID(ctx, db, userInfo, userIDs, opts)
	if err != nil {
		return nil, err
	}

	byUserID := make(map[int32]*jobscolleagues.Colleague, len(colleagues))
	for _, colleague := range colleagues {
		byUserID[colleague.GetUserId()] = colleague
	}
	return byUserID, nil
}

func (h *Hydrator) HydrateBasicTargets(
	ctx context.Context,
	db qrm.DB,
	userInfo *userinfo.UserInfo,
	targets []BasicTarget,
	opts ResolveOpts,
) error {
	if len(targets) == 0 {
		return nil
	}

	userIDs := make([]int32, 0, len(targets))
	for _, target := range targets {
		if target.UserID > 0 && target.Set != nil {
			userIDs = append(userIDs, target.UserID)
		}
	}
	colleaguesByUserID, err := h.HydrateBasicByUserID(ctx, db, userInfo, userIDs, opts)
	if err != nil {
		return err
	}
	for _, target := range targets {
		if colleague, ok := colleaguesByUserID[target.UserID]; target.Set != nil && ok {
			target.Set(colleague)
		}
	}
	return nil
}

func (h *Hydrator) GetByUserID(
	ctx context.Context,
	db qrm.DB,
	userInfo *userinfo.UserInfo,
	userID int32,
	opts ResolveOpts,
) (*jobscolleagues.Colleague, error) {
	colleagues, err := h.ListByUserID(ctx, db, userInfo, []int32{userID}, opts)
	if err != nil {
		return nil, err
	}
	if len(colleagues) == 0 {
		return nil, nil
	}
	return colleagues[0], nil
}

func (h *Hydrator) HydrateByUserID(
	ctx context.Context,
	db qrm.DB,
	userInfo *userinfo.UserInfo,
	userIDs []int32,
	opts ResolveOpts,
) (map[int32]*jobscolleagues.Colleague, error) {
	colleagues, err := h.ListByUserID(ctx, db, userInfo, userIDs, opts)
	if err != nil {
		return nil, err
	}

	byUserID := make(map[int32]*jobscolleagues.Colleague, len(colleagues))
	for _, colleague := range colleagues {
		byUserID[colleague.GetUserId()] = colleague
	}

	return byUserID, nil
}

func (h *Hydrator) HydrateTargets(
	ctx context.Context,
	db qrm.DB,
	userInfo *userinfo.UserInfo,
	targets []Target,
	opts ResolveOpts,
) error {
	if len(targets) == 0 {
		return nil
	}

	userIDs := make([]int32, 0, len(targets))
	for _, target := range targets {
		if target.UserID <= 0 || target.Set == nil {
			continue
		}

		userIDs = append(userIDs, target.UserID)
	}
	if len(userIDs) == 0 {
		return nil
	}

	colleaguesByUserID, err := h.HydrateByUserID(ctx, db, userInfo, userIDs, opts)
	if err != nil {
		return err
	}

	for _, target := range targets {
		if target.Set == nil {
			continue
		}

		colleague, ok := colleaguesByUserID[target.UserID]
		if !ok {
			continue
		}

		target.Set(colleague)
	}

	return nil
}

type jobGroup struct {
	job     string
	userIDs []int32
}

func (h *Hydrator) resolveJob(
	ctx context.Context,
	db qrm.DB,
	userInfo *userinfo.UserInfo,
	userIDs []int32,
	opts ResolveOpts,
) ([]jobGroup, error) {
	switch opts.Scope.Mode {
	case JobScopePrimary:
		return h.resolvePrimaryJob(ctx, db, userIDs, opts.Scope.Job)

	case JobScopeExplicit:
		if opts.Scope.Job == "" {
			return nil, nil
		}
		return []jobGroup{{job: opts.Scope.Job, userIDs: userIDs}}, nil

	default:
		job := opts.Scope.Job
		if job == "" && userInfo != nil {
			job = userInfo.GetJob()
		}
		if job == "" {
			return nil, errors.New("colleague hydrator requires a job for caller-scoped hydration")
		}
		return []jobGroup{{job: job, userIDs: userIDs}}, nil
	}
}

func (h *Hydrator) resolvePrimaryJob(
	ctx context.Context,
	db qrm.DB,
	userIDs []int32,
	fallbackJob string,
) ([]jobGroup, error) {
	tUser := table.FivenetUser.AS("user")

	userIdExprs := make([]mysql.Expression, 0, len(userIDs))
	for _, value := range userIDs {
		userIdExprs = append(userIdExprs, mysql.Int32(value))
	}

	stmt := tUser.
		SELECT(
			tUser.ID.AS("user_id"),
			tUser.Job.AS("job"),
		).
		FROM(tUser).
		WHERE(mysql.AND(
			tUser.ID.IN(userIdExprs...),
			tUser.DeletedAt.IS_NULL(),
		)).
		LIMIT(int64(len(userIDs)))

	var rows []struct {
		UserID int32  `alias:"user_id"`
		Job    string `alias:"job"`
	}
	if err := stmt.QueryContext(ctx, db, &rows); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	grouped := make(map[string][]int32, len(rows))
	for _, row := range rows {
		if row.Job == "" {
			continue
		}
		grouped[row.Job] = append(grouped[row.Job], row.UserID)
	}

	if len(grouped) == 0 {
		return []jobGroup{
			{
				job:     fallbackJob,
				userIDs: userIDs,
			},
		}, nil
	}

	out := make([]jobGroup, 0, len(grouped))
	for job, ids := range grouped {
		out = append(out, jobGroup{
			job:     job,
			userIDs: slices.Compact(ids),
		})
	}

	return out, nil
}

func (h *Hydrator) attachLabels(
	ctx context.Context,
	db qrm.DB,
	job string,
	colleagues []*jobscolleagues.Colleague,
	includeDeleted bool,
) error {
	userIDs := make([]int32, 0, len(colleagues))
	for _, colleague := range colleagues {
		userIDs = append(userIDs, colleague.GetUserId())
	}

	labels, err := h.store.GetUsersLabels(ctx, db, job, userIDs, includeDeleted)
	if err != nil {
		return err
	}

	labelsByUser := make(map[int32]*jobsstore.UserLabels, len(labels))
	for _, userLabels := range labels {
		labelsByUser[userLabels.UserId] = userLabels
	}

	for _, colleague := range colleagues {
		userLabels, ok := labelsByUser[colleague.GetUserId()]
		if !ok {
			continue
		}
		if colleague.GetProps() == nil {
			colleague.Props = &jobscolleagues.ColleagueProps{
				UserId: colleague.GetUserId(),
				Job:    job,
			}
		}
		colleague.Props.Labels = userLabels.Labels
	}

	return nil
}

func columnsForFields(
	fields resolvedFields,
) mysql.ProjectionList {
	columns := mysql.ProjectionList{}
	if fields.note {
		columns = append(columns, table.FivenetJobColleagueProps.AS("colleague_props").Note)
	}

	return columns
}

func uniqueUserIDs(userIDs []int32) []int32 {
	seen := make(map[int32]struct{}, len(userIDs))
	out := make([]int32, 0, len(userIDs))
	for _, userID := range userIDs {
		if userID <= 0 {
			continue
		}
		if _, ok := seen[userID]; ok {
			continue
		}
		seen[userID] = struct{}{}
		out = append(out, userID)
	}

	return out
}
