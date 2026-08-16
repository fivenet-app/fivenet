package colleagueshydrator

import (
	"context"
	"database/sql"

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

var tColleagueProps = table.FivenetJobColleagueProps.AS("colleague_props")

var Module = fx.Module(
	"stores.jobs.colleagues.hydrator",
	fx.Provide(New),
)

type IHydrator interface {
	ListByUserID(
		ctx context.Context,
		db qrm.DB,
		userInfo *userinfo.UserInfo,
		job string,
		userIDs []int32,
		infoOnly bool,
	) ([]*jobscolleagues.Colleague, error)
	HydrateByUserID(
		ctx context.Context,
		db qrm.DB,
		userInfo *userinfo.UserInfo,
		job string,
		userIDs []int32,
		infoOnly bool,
	) (map[int32]*jobscolleagues.Colleague, error)
	HydrateTargets(
		ctx context.Context,
		db qrm.DB,
		userInfo *userinfo.UserInfo,
		job string,
		targets []Target,
		infoOnly bool,
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
) (*perms.TypedStringList[permsjobs.ColleaguesServiceGetColleagueTypesPermValue], error) {
	fields := perms.NewTypedStringList[permsjobs.ColleaguesServiceGetColleagueTypesPermValue]()
	if !infoOnly {
		var err error
		fields, err = permsjobs.ColleaguesService.GetColleague.TypesTyped.Get(h.perms, userInfo)
		if err != nil {
			return nil, err
		}
	}

	if userInfo.GetJobAdmin() {
		fields.Set(permsjobs.ColleaguesServiceGetColleagueTypesPermValueNote)
	}

	return fields, nil
}

func (h *Hydrator) ListByUserID(
	ctx context.Context,
	db qrm.DB,
	userInfo *userinfo.UserInfo,
	job string,
	userIDs []int32,
	infoOnly bool,
) ([]*jobscolleagues.Colleague, error) {
	if len(userIDs) == 0 {
		return []*jobscolleagues.Colleague{}, nil
	}
	if db == nil {
		db = h.db
	}

	fields, err := h.getFields(userInfo, infoOnly)
	if err != nil {
		return nil, err
	}

	colleagues, err := h.store.ListColleaguesByUserIDs(
		ctx,
		db,
		jobsstore.ListColleaguesByUserIDsQuery{
			Job:         job,
			UserIDs:     uniqueUserIDs(userIDs),
			WithColumns: columnsForFields(fields),
		},
	)
	if err != nil {
		return nil, err
	}

	if fields.Contains(permsjobs.ColleaguesServiceGetColleagueTypesPermValueLabels) ||
		userInfo.GetJobAdmin() {
		if err := h.attachLabels(ctx, db, job, colleagues, userInfo.GetJobAdmin()); err != nil {
			return nil, err
		}
	}

	enrichJobInfo := h.enricher.EnrichJobInfoSafeFunc(userInfo)
	for _, colleague := range colleagues {
		enrichJobInfo(colleague)
	}

	return colleagues, nil
}

func (h *Hydrator) HydrateByUserID(
	ctx context.Context,
	db qrm.DB,
	userInfo *userinfo.UserInfo,
	job string,
	userIDs []int32,
	infoOnly bool,
) (map[int32]*jobscolleagues.Colleague, error) {
	colleagues, err := h.ListByUserID(ctx, db, userInfo, job, userIDs, infoOnly)
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
	job string,
	targets []Target,
	infoOnly bool,
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

	colleaguesByUserID, err := h.HydrateByUserID(ctx, db, userInfo, job, userIDs, infoOnly)
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
			columns = append(columns, tColleagueProps.Note)
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
