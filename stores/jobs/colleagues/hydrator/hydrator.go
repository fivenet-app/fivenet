package colleagueshydrator

import (
	"context"
	"database/sql"
	"errors"
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

type PropsJobMode int

const (
	PropsJobModeCaller PropsJobMode = iota
	PropsJobModePrimary
	PropsJobModeExplicit
)

// ResolveOpts options to adjust the resolve behavior.
type ResolveOpts struct {
	UserInfo *userinfo.UserInfo
	InfoOnly bool

	PropsJobMode PropsJobMode
	PropsJob     string
}

type IHydrator interface {
	ListByUserID(
		ctx context.Context,
		db qrm.DB,
		userIDs []int32,
		args ResolveOpts,
	) ([]*jobscolleagues.Colleague, error)
	GetShortByUserID(
		ctx context.Context,
		db qrm.DB,
		userID int32,
		args ResolveOpts,
	) (*jobscolleagues.Colleague, error)

	HydrateByUserID(
		ctx context.Context,
		db qrm.DB,
		userIDs []int32,
		args ResolveOpts,
	) (map[int32]*jobscolleagues.Colleague, error)
	HydrateTargets(
		ctx context.Context,
		db qrm.DB,
		targets []Target,
		args ResolveOpts,
	) error
}

type Target struct {
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

func (h *Hydrator) getFields(
	userInfo *userinfo.UserInfo,
	infoOnly bool,
) (*perms.TypedStringList[permsjobs.ColleaguesServiceGetColleagueTypesPermValue], bool, error) {
	fields := perms.NewTypedStringList[permsjobs.ColleaguesServiceGetColleagueTypesPermValue]()
	if userInfo == nil {
		return fields, false, nil
	}
	jobAdmin := userInfo.GetJobAdmin()
	if !infoOnly {
		var err error
		fields, err = permsjobs.ColleaguesService.GetColleague.TypesTyped.Get(h.perms, userInfo)
		if err != nil {
			return nil, false, err
		}
	}

	if jobAdmin {
		fields.Set(permsjobs.ColleaguesServiceGetColleagueTypesPermValueNote)
	}

	return fields, jobAdmin, nil
}

func (h *Hydrator) ListByUserID(
	ctx context.Context,
	db qrm.DB,
	userIDs []int32,
	args ResolveOpts,
) ([]*jobscolleagues.Colleague, error) {
	if len(userIDs) == 0 {
		return []*jobscolleagues.Colleague{}, nil
	}
	if db == nil {
		db = h.db
	}

	fields, jobAdmin, err := h.getFields(args.UserInfo, args.InfoOnly)
	if err != nil {
		return nil, err
	}

	jobGroups, err := h.resolveJobGroups(ctx, db, uniqueUserIDs(userIDs), args)
	if err != nil {
		return nil, err
	}

	colleaguesByUserID := make(map[int32]*jobscolleagues.Colleague, len(userIDs))
	for _, jobGroup := range jobGroups {
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

		if fields.Contains(permsjobs.ColleaguesServiceGetColleagueTypesPermValueLabels) ||
			jobAdmin {
			if err := h.attachLabels(ctx, db, jobGroup.job, colleagues, jobAdmin); err != nil {
				return nil, err
			}
		}

		enrichJobInfo := h.enricher.EnrichJobInfoSafeFunc(args.UserInfo)
		for _, colleague := range colleagues {
			enrichJobInfo(colleague)
			colleaguesByUserID[colleague.GetUserId()] = colleague
		}
	}

	out := make([]*jobscolleagues.Colleague, 0, len(userIDs))
	for _, userID := range userIDs {
		if colleague, ok := colleaguesByUserID[userID]; ok {
			out = append(out, colleague)
		}
	}

	return out, nil
}

func (h *Hydrator) GetShortByUserID(
	ctx context.Context,
	db qrm.DB,
	userID int32,
	args ResolveOpts,
) (*jobscolleagues.Colleague, error) {
	args.InfoOnly = true

	colleagues, err := h.ListByUserID(ctx, db, []int32{userID}, args)
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
	userIDs []int32,
	args ResolveOpts,
) (map[int32]*jobscolleagues.Colleague, error) {
	colleagues, err := h.ListByUserID(ctx, db, userIDs, args)
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
	targets []Target,
	args ResolveOpts,
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

	colleaguesByUserID, err := h.HydrateByUserID(ctx, db, userIDs, args)
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

func (h *Hydrator) resolveJobGroups(
	ctx context.Context,
	db qrm.DB,
	userIDs []int32,
	args ResolveOpts,
) ([]jobGroup, error) {
	switch args.PropsJobMode {
	case PropsJobModePrimary:
		return h.resolvePrimaryJobGroups(ctx, db, userIDs)

	case PropsJobModeExplicit:
		if args.PropsJob == "" {
			return nil, nil
		}
		return []jobGroup{{job: args.PropsJob, userIDs: userIDs}}, nil

	default:
		job := args.PropsJob
		if job == "" && args.UserInfo != nil {
			job = args.UserInfo.GetJob()
		}
		if job == "" {
			return nil, errors.New("colleague hydrator requires a job for caller-scoped hydration")
		}
		return []jobGroup{{job: job, userIDs: userIDs}}, nil
	}
}

func (h *Hydrator) resolvePrimaryJobGroups(
	ctx context.Context,
	db qrm.DB,
	userIDs []int32,
) ([]jobGroup, error) {
	tUser := table.FivenetUser.AS("user")

	userIdExprs := make([]mysql.Expression, 0, len(userIDs))
	for _, value := range userIDs {
		userIdExprs = append(userIdExprs, mysql.Int32(value))
	}

	stmt := tUser.
		SELECT(
			tUser.ID.AS("user_id"),
			tUser.Job,
		).
		FROM(tUser).
		WHERE(tUser.ID.IN(userIdExprs...)).
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
		return nil, nil
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
	fields *perms.TypedStringList[permsjobs.ColleaguesServiceGetColleagueTypesPermValue],
) mysql.ProjectionList {
	columns := mysql.ProjectionList{}
	for _, field := range fields.Values() {
		switch field {
		case permsjobs.ColleaguesServiceGetColleagueTypesPermValueNote:
			columns = append(columns, table.FivenetJobColleagueProps.AS("colleague_props").Note)
		}
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
