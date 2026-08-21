package citizenshydrator

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"

	citizenslicenses "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/citizens/licenses"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	users "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users"
	usersprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/props"
	usershort "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/short"
	permscitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2026/pkg/perms"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	citizensstore "github.com/fivenet-app/fivenet/v2026/stores/citizens"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"stores.citizens.hydrator",
	fx.Provide(New),
)

type ResolveOpts struct {
	InfoOnly bool

	IncludePhoneNumber bool
	IncludeProps       bool
	IncludeLabels      bool
	IncludeLicenses    bool
}

type IHydrator interface {
	ListByUserID(
		ctx context.Context,
		db qrm.DB,
		userInfo *userinfo.UserInfo,
		userIDs []int32,
		opts ResolveOpts,
	) ([]*users.User, error)
	HydrateByUserID(
		ctx context.Context,
		db qrm.DB,
		userInfo *userinfo.UserInfo,
		userIDs []int32,
		opts ResolveOpts,
	) (map[int32]*users.User, error)
	HydrateTargets(
		ctx context.Context,
		db qrm.DB,
		userInfo *userinfo.UserInfo,
		targets []Target,
		opts ResolveOpts,
	) error

	ListShortByUserID(
		ctx context.Context,
		db qrm.DB,
		userInfo *userinfo.UserInfo,
		userIDs []int32,
		opts ResolveOpts,
	) ([]*usershort.UserShort, error)
	HydrateShortByUserID(
		ctx context.Context,
		db qrm.DB,
		userInfo *userinfo.UserInfo,
		userIDs []int32,
		opts ResolveOpts,
	) (map[int32]*usershort.UserShort, error)
	HydrateShortTargets(
		ctx context.Context,
		db qrm.DB,
		userInfo *userinfo.UserInfo,
		targets []ShortTarget,
		opts ResolveOpts,
	) error
	HydrateShortTargetsSafeFunc(
		userInfo *userinfo.UserInfo,
	) func(
		ctx context.Context,
		db qrm.DB,
		targets []ShortTarget,
	) error

	GetShortByUserID(
		ctx context.Context,
		db qrm.DB,
		userID int32,
	) (*usershort.UserShort, error)
	GetShortByUserIDSafeFunc(
		userInfo *userinfo.UserInfo,
	) func(
		ctx context.Context,
		db qrm.DB,
		userID int32,
	) (*usershort.UserShort, error)
}

type Target struct {
	UserID int32
	Set    func(*users.User)
}

type ShortTarget struct {
	UserID int32
	Set    func(*usershort.UserShort)
}

type Hydrator struct {
	db       *sql.DB
	customDB *config.CustomDB
	appCfg   appconfig.IConfig
	perms    perms.Permissions
	enricher mstlystcdata.IUserAwareEnricher
	store    citizensstore.IStore
}

type Params struct {
	fx.In

	DB                *sql.DB
	CustomDB          *config.CustomDB
	AppConfig         appconfig.IConfig
	Perms             perms.Permissions
	UserAwareEnricher mstlystcdata.IUserAwareEnricher
	Store             citizensstore.IStore
}

func New(p Params) IHydrator {
	return &Hydrator{
		db:       p.DB,
		customDB: p.CustomDB,
		appCfg:   p.AppConfig,
		perms:    p.Perms,
		enricher: p.UserAwareEnricher,
		store:    p.Store,
	}
}

func resolveShortOpts(
	userInfo *userinfo.UserInfo,
	perms perms.Permissions,
) ResolveOpts {
	opts := ResolveOpts{}
	if userInfo == nil {
		return opts
	}

	userFields, err := permscitizens.CitizensService.ListCitizens.FieldsTyped.Get(perms, userInfo)
	if err != nil {
		return opts
	}

	opts.IncludePhoneNumber = userFields.Contains(
		permscitizens.CitizensServiceListCitizensFieldsPermValuePhoneNumber,
	)
	return opts
}

func (h *Hydrator) ListByUserID(
	ctx context.Context,
	db qrm.DB,
	userInfo *userinfo.UserInfo,
	userIDs []int32,
	opts ResolveOpts,
) ([]*users.User, error) {
	if len(userIDs) == 0 {
		return []*users.User{}, nil
	}
	if db == nil {
		db = h.db
	}

	userIDs = uniqueUserIDs(userIDs)

	usersByID, orderedIDs, err := h.loadUsers(ctx, db, userIDs, opts)
	if err != nil {
		return nil, err
	}

	if opts.IncludeProps || opts.IncludeLabels || opts.IncludeLicenses {
		if err := h.loadProps(ctx, db, userInfo, usersByID, orderedIDs, opts); err != nil {
			return nil, err
		}
	}

	out := make([]*users.User, 0, len(orderedIDs))
	for _, userID := range orderedIDs {
		if user, ok := usersByID[userID]; ok {
			out = append(out, user)
		}
	}

	return out, nil
}

func (h *Hydrator) HydrateByUserID(
	ctx context.Context,
	db qrm.DB,
	userInfo *userinfo.UserInfo,
	userIDs []int32,
	opts ResolveOpts,
) (map[int32]*users.User, error) {
	usersList, err := h.ListByUserID(ctx, db, userInfo, userIDs, opts)
	if err != nil {
		return nil, err
	}

	byUserID := make(map[int32]*users.User, len(usersList))
	for _, user := range usersList {
		byUserID[user.GetUserId()] = user
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

	usersByID, err := h.HydrateByUserID(ctx, db, userInfo, userIDs, opts)
	if err != nil {
		return err
	}

	for _, target := range targets {
		if target.Set == nil {
			continue
		}

		user, ok := usersByID[target.UserID]
		if !ok {
			continue
		}

		target.Set(user)
	}

	return nil
}

func (h *Hydrator) ListShortByUserID(
	ctx context.Context,
	db qrm.DB,
	userInfo *userinfo.UserInfo,
	userIDs []int32,
	opts ResolveOpts,
) ([]*usershort.UserShort, error) {
	usersList, err := h.ListByUserID(ctx, db, userInfo, userIDs, opts)
	if err != nil {
		return nil, err
	}

	out := make([]*usershort.UserShort, 0, len(usersList))
	for _, user := range usersList {
		out = append(out, user.UserShort())
	}

	return out, nil
}

func (h *Hydrator) HydrateShortByUserID(
	ctx context.Context,
	db qrm.DB,
	userInfo *userinfo.UserInfo,
	userIDs []int32,
	opts ResolveOpts,
) (map[int32]*usershort.UserShort, error) {
	usersList, err := h.ListShortByUserID(ctx, db, userInfo, userIDs, opts)
	if err != nil {
		return nil, err
	}

	byUserID := make(map[int32]*usershort.UserShort, len(usersList))
	for _, user := range usersList {
		byUserID[user.GetUserId()] = user
	}

	return byUserID, nil
}

func (h *Hydrator) HydrateShortTargets(
	ctx context.Context,
	db qrm.DB,
	userInfo *userinfo.UserInfo,
	targets []ShortTarget,
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

	usersByID, err := h.HydrateShortByUserID(ctx, db, userInfo, userIDs, opts)
	if err != nil {
		return err
	}

	for _, target := range targets {
		if target.Set == nil {
			continue
		}

		user, ok := usersByID[target.UserID]
		if !ok {
			continue
		}

		target.Set(user)
	}

	return nil
}

func (h *Hydrator) HydrateShortTargetsSafeFunc(
	userInfo *userinfo.UserInfo,
) func(
	ctx context.Context,
	db qrm.DB,
	targets []ShortTarget,
) error {
	opts := resolveShortOpts(userInfo, h.perms)

	return func(ctx context.Context, db qrm.DB, targets []ShortTarget) error {
		return h.HydrateShortTargets(ctx, db, userInfo, targets, opts)
	}
}

func (h *Hydrator) GetShortByUserID(
	ctx context.Context,
	db qrm.DB,
	userID int32,
) (*usershort.UserShort, error) {
	if userID <= 0 {
		return nil, nil
	}
	if db == nil {
		db = h.db
	}

	usersList, err := h.ListShortByUserID(
		ctx,
		db,
		nil,
		[]int32{userID},
		ResolveOpts{IncludePhoneNumber: true},
	)
	if err != nil {
		return nil, err
	}
	if len(usersList) == 0 {
		return nil, nil
	}

	return usersList[0], nil
}

func (h *Hydrator) GetShortByUserIDSafeFunc(
	userInfo *userinfo.UserInfo,
) func(
	ctx context.Context,
	db qrm.DB,
	userID int32,
) (*usershort.UserShort, error) {
	opts := resolveShortOpts(userInfo, h.perms)

	return func(ctx context.Context, db qrm.DB, userID int32) (*usershort.UserShort, error) {
		if userID <= 0 {
			return nil, nil
		}
		if db == nil {
			db = h.db
		}

		usersList, err := h.ListShortByUserID(ctx, db, userInfo, []int32{userID}, opts)
		if err != nil {
			return nil, err
		}
		if len(usersList) == 0 {
			return nil, nil
		}

		return usersList[0], nil
	}
}

func (h *Hydrator) loadUsers(
	ctx context.Context,
	db qrm.DB,
	userIDs []int32,
	opts ResolveOpts,
) (map[int32]*users.User, []int32, error) {
	tUser := table.FivenetUser.AS("user")

	selectors := []mysql.Projection{
		tUser.ID,
		tUser.Firstname,
		tUser.Lastname,
		tUser.Job,
		tUser.JobGrade,
		tUser.Dateofbirth,
		tUser.Sex,
		tUser.Height,
		h.customDB.Columns.User.GetVisum(tUser.Alias()),
	}
	if opts.IncludePhoneNumber {
		selectors = append(selectors, tUser.PhoneNumber)
	}

	userIDExprs := make([]mysql.Expression, 0, len(userIDs))
	for _, userID := range userIDs {
		userIDExprs = append(userIDExprs, mysql.Int32(userID))
	}

	stmt := tUser.
		SELECT(selectors[0], selectors[1:]...).
		FROM(tUser).
		WHERE(tUser.ID.IN(userIDExprs...)).
		LIMIT(int64(len(userIDs)))

	dest := []*users.User{}
	if err := stmt.QueryContext(ctx, db, &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return map[int32]*users.User{}, userIDs, nil
		}
		return nil, nil, fmt.Errorf("failed to load citizens by id: %w", err)
	}

	usersByID := make(map[int32]*users.User, len(dest))
	for _, user := range dest {
		if user == nil || user.GetUserId() <= 0 {
			continue
		}
		usersByID[user.GetUserId()] = user
	}

	return usersByID, userIDs, nil
}

func (h *Hydrator) loadProps(
	ctx context.Context,
	db qrm.DB,
	userInfo *userinfo.UserInfo,
	usersByID map[int32]*users.User,
	orderedIDs []int32,
	opts ResolveOpts,
) error {
	userIDs := make([]int32, 0, len(orderedIDs))
	for _, userID := range orderedIDs {
		if _, ok := usersByID[userID]; ok {
			userIDs = append(userIDs, userID)
		}
	}
	if len(userIDs) == 0 {
		return nil
	}

	propsByUserID, err := h.loadUserProps(ctx, db, userIDs)
	if err != nil {
		return err
	}

	publicJobs := h.appCfg.Get().GetJobInfo().GetPublicJobs()

	for _, userID := range orderedIDs {
		user := usersByID[userID]
		if user == nil {
			continue
		}

		props := propsByUserID[userID]
		if props == nil {
			props = &usersprops.UserProps{UserId: userID}
		}

		props.Default()
		if props.JobName != nil {
			grade := props.GetJobGradeNumber()
			props.Job, props.JobGrade = h.enricher.GetJobGrade(props.GetJobName(), grade)
		}

		if opts.IncludeProps || opts.IncludeLabels || opts.IncludeLicenses {
			user.Props = props
		}

		if props.JobName != nil && !slices.Contains(publicJobs, user.GetJob()) {
			user.Job = props.GetJobName()
			if props.JobGradeNumber != nil {
				user.JobGrade = props.GetJobGradeNumber()
			} else {
				user.JobGrade = 0
			}

			h.enricher.EnrichJobInfo(user)
		} else {
			h.enricher.EnrichJobInfoSafe(userInfo, user)
		}

		if opts.IncludeLabels {
			labels, err := h.store.GetUserLabelsForUser(ctx, userInfo, userID)
			if err != nil {
				return err
			}
			user.Props.Labels = labels
		}

		if opts.IncludeLicenses && !opts.InfoOnly {
			licenses, err := h.loadLicenses(ctx, db, userID)
			if err != nil {
				return err
			}
			user.Licenses = licenses
		}
	}

	return nil
}

func (h *Hydrator) loadUserProps(
	ctx context.Context,
	db qrm.DB,
	userIDs []int32,
) (map[int32]*usersprops.UserProps, error) {
	if len(userIDs) == 0 {
		return map[int32]*usersprops.UserProps{}, nil
	}

	tUserProps := table.FivenetUserProps.AS("user_props")
	tFiles := table.FivenetFiles.AS("mugshot")
	userIDExprs := make([]mysql.Expression, 0, len(userIDs))
	for _, userID := range userIDs {
		userIDExprs = append(userIDExprs, mysql.Int32(userID))
	}

	stmt := tUserProps.
		SELECT(
			tUserProps.UserID,
			tUserProps.UpdatedAt,
			tUserProps.Wanted,
			tUserProps.Job,
			tUserProps.JobGrade,
			tUserProps.TrafficInfractionPoints,
			tUserProps.TrafficInfractionPointsUpdatedAt,
			tUserProps.OpenFines,
			tUserProps.BloodType,
			tUserProps.MugshotFileID,
			tUserProps.Email,
			tFiles.ID.AS("mugshot.mugshot_file_id"),
			tFiles.FilePath,
		).
		FROM(
			tUserProps.
				LEFT_JOIN(tFiles,
					tFiles.ID.EQ(tUserProps.MugshotFileID),
				),
		).
		WHERE(tUserProps.UserID.IN(userIDExprs...))

	dest := []*usersprops.UserProps{}
	if err := stmt.QueryContext(ctx, db, &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return map[int32]*usersprops.UserProps{}, nil
		}
		return nil, fmt.Errorf("failed to load citizen props: %w", err)
	}

	propsByUserID := make(map[int32]*usersprops.UserProps, len(dest))
	for _, props := range dest {
		if props == nil || props.GetUserId() <= 0 {
			continue
		}
		propsByUserID[props.GetUserId()] = props
	}

	return propsByUserID, nil
}

func (h *Hydrator) loadLicenses(
	ctx context.Context,
	db qrm.Queryable,
	userID int32,
) ([]*citizenslicenses.License, error) {
	tCitizenLicenses := table.FivenetUserLicenses
	tLicenses := table.FivenetLicenses

	stmt := tCitizenLicenses.
		SELECT(
			tLicenses.Type.AS("type"),
			tLicenses.Label.AS("label"),
		).
		FROM(
			tCitizenLicenses.
				LEFT_JOIN(tLicenses,
					tCitizenLicenses.Type.EQ(tLicenses.Type)),
		).
		WHERE(tCitizenLicenses.UserID.EQ(mysql.Int32(userID))).
		LIMIT(15)

	var dest []*citizenslicenses.License
	if err := stmt.QueryContext(ctx, db, &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return []*citizenslicenses.License{}, nil
		}
		return nil, fmt.Errorf("failed to load citizen licenses: %w", err)
	}

	return dest, nil
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
