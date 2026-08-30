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

type IHydrator interface {
	GetBasicByUserID(
		ctx context.Context,
		db qrm.DB,
		userInfo *userinfo.UserInfo,
		userID int32,
	) (*usershort.UserShort, error)
	ListBasicByUserID(
		ctx context.Context,
		db qrm.DB,
		userInfo *userinfo.UserInfo,
		userIDs []int32,
	) ([]*usershort.UserShort, error)
	HydrateBasicByUserID(
		ctx context.Context,
		db qrm.DB,
		userInfo *userinfo.UserInfo,
		userIDs []int32,
	) (map[int32]*usershort.UserShort, error)
	HydrateBasicTargets(
		ctx context.Context,
		db qrm.DB,
		userInfo *userinfo.UserInfo,
		targets []BasicTarget,
	) error
	HydrateBasicTargetsSafeFunc(
		userInfo *userinfo.UserInfo,
	) func(ctx context.Context, db qrm.DB, targets []BasicTarget) error

	GetByUserID(
		ctx context.Context,
		db qrm.DB,
		userInfo *userinfo.UserInfo,
		userID int32,
	) (*users.User, error)
	GetBasicByUserIDSafeFunc(
		userInfo *userinfo.UserInfo,
	) func(ctx context.Context, db qrm.DB, userID int32) (*usershort.UserShort, error)
	ListByUserID(
		ctx context.Context,
		db qrm.DB,
		userInfo *userinfo.UserInfo,
		userIDs []int32,
	) ([]*users.User, error)
	HydrateByUserID(
		ctx context.Context,
		db qrm.DB,
		userInfo *userinfo.UserInfo,
		userIDs []int32,
	) (map[int32]*users.User, error)
	HydrateTargets(
		ctx context.Context,
		db qrm.DB,
		userInfo *userinfo.UserInfo,
		targets []Target,
	) error
}

type Target struct {
	UserID int32
	Set    func(*users.User)
}

type BasicTarget struct {
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

type resolvedFields struct {
	phoneNumber             bool
	licenses                bool
	wanted                  bool
	job                     bool
	trafficInfractionPoints bool
	openFines               bool
	bloodType               bool
	mugshot                 bool
	labels                  bool
	email                   bool
}

func (f resolvedFields) hasProps() bool {
	return f.wanted || f.job || f.trafficInfractionPoints || f.openFines ||
		f.bloodType || f.mugshot || f.labels || f.email
}

func (h *Hydrator) resolveFields(userInfo *userinfo.UserInfo) (resolvedFields, error) {
	if userInfo == nil {
		return resolvedFields{}, nil
	}

	fields, err := permscitizens.CitizensService.ListCitizens.FieldsTyped.Get(h.perms, userInfo)
	if err != nil {
		return resolvedFields{}, err
	}

	return resolvedFields{
		phoneNumber: fields.Contains(
			permscitizens.CitizensServiceListCitizensFieldsPermValuePhoneNumber,
		),
		licenses: fields.Contains(
			permscitizens.CitizensServiceListCitizensFieldsPermValueLicenses,
		),
		wanted: fields.Contains(
			permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsWanted,
		),
		job: fields.Contains(
			permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsJob,
		),
		trafficInfractionPoints: fields.Contains(
			permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsTrafficInfractionPoints,
		),
		openFines: fields.Contains(
			permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsOpenFines,
		),
		bloodType: fields.Contains(
			permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsBloodType,
		),
		mugshot: fields.Contains(
			permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsMugshot,
		),
		labels: fields.Contains(
			permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsLabels,
		),
		email: fields.Contains(
			permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsEmail,
		),
	}, nil
}

func (h *Hydrator) GetBasicByUserID(
	ctx context.Context,
	db qrm.DB,
	userInfo *userinfo.UserInfo,
	userID int32,
) (*usershort.UserShort, error) {
	users, err := h.ListBasicByUserID(ctx, db, userInfo, []int32{userID})
	if err != nil || len(users) == 0 {
		return nil, err
	}
	return users[0], nil
}

func (h *Hydrator) ListBasicByUserID(
	ctx context.Context,
	db qrm.DB,
	userInfo *userinfo.UserInfo,
	userIDs []int32,
) ([]*usershort.UserShort, error) {
	if len(userIDs) == 0 {
		return []*usershort.UserShort{}, nil
	}
	if db == nil {
		db = h.db
	}

	fields, err := h.resolveFields(userInfo)
	if err != nil {
		return nil, err
	}
	userIDs = uniqueUserIDs(userIDs)

	usersByID, orderedIDs, err := h.loadBasicUsers(ctx, db, userIDs, fields)
	if err != nil {
		return nil, err
	}

	out := make([]*usershort.UserShort, 0, len(orderedIDs))
	for _, userID := range orderedIDs {
		if user, ok := usersByID[userID]; ok {
			out = append(out, user)
		}
	}
	h.enrichShortsSafe(userInfo, out...)

	return out, nil
}

func (h *Hydrator) HydrateBasicByUserID(
	ctx context.Context,
	db qrm.DB,
	userInfo *userinfo.UserInfo,
	userIDs []int32,
) (map[int32]*usershort.UserShort, error) {
	usersList, err := h.ListBasicByUserID(ctx, db, userInfo, userIDs)
	if err != nil {
		return nil, err
	}

	byUserID := make(map[int32]*usershort.UserShort, len(usersList))
	for _, user := range usersList {
		byUserID[user.GetUserId()] = user
	}

	return byUserID, nil
}

func (h *Hydrator) HydrateBasicTargets(
	ctx context.Context,
	db qrm.DB,
	userInfo *userinfo.UserInfo,
	targets []BasicTarget,
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

	usersByID, err := h.HydrateBasicByUserID(ctx, db, userInfo, userIDs)
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

func (h *Hydrator) HydrateBasicTargetsSafeFunc(
	userInfo *userinfo.UserInfo,
) func(context.Context, qrm.DB, []BasicTarget) error {
	return func(ctx context.Context, db qrm.DB, targets []BasicTarget) error {
		return h.HydrateBasicTargets(ctx, db, userInfo, targets)
	}
}

func (h *Hydrator) GetBasicByUserIDSafeFunc(
	userInfo *userinfo.UserInfo,
) func(context.Context, qrm.DB, int32) (*usershort.UserShort, error) {
	return func(ctx context.Context, db qrm.DB, userID int32) (*usershort.UserShort, error) {
		return h.GetBasicByUserID(ctx, db, userInfo, userID)
	}
}

func (h *Hydrator) GetByUserID(
	ctx context.Context,
	db qrm.DB,
	userInfo *userinfo.UserInfo,
	userID int32,
) (*users.User, error) {
	users, err := h.ListByUserID(ctx, db, userInfo, []int32{userID})
	if err != nil || len(users) == 0 {
		return nil, err
	}
	return users[0], nil
}

func (h *Hydrator) ListByUserID(
	ctx context.Context,
	db qrm.DB,
	userInfo *userinfo.UserInfo,
	userIDs []int32,
) ([]*users.User, error) {
	if len(userIDs) == 0 {
		return []*users.User{}, nil
	}
	if db == nil {
		db = h.db
	}

	fields, err := h.resolveFields(userInfo)
	if err != nil {
		return nil, err
	}
	userIDs = uniqueUserIDs(userIDs)

	usersByID, orderedIDs, err := h.loadUsers(ctx, db, userIDs, fields)
	if err != nil {
		return nil, err
	}
	if fields.hasProps() || fields.licenses {
		if err := h.loadProps(ctx, db, userInfo, usersByID, orderedIDs, fields); err != nil {
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
) (map[int32]*users.User, error) {
	usersList, err := h.ListByUserID(ctx, db, userInfo, userIDs)
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
	usersByID, err := h.HydrateByUserID(ctx, db, userInfo, userIDs)
	if err != nil {
		return err
	}
	for _, target := range targets {
		if user, ok := usersByID[target.UserID]; target.Set != nil && ok {
			target.Set(user)
		}
	}
	return nil
}

func (h *Hydrator) enrichShortsSafe(
	userInfo *userinfo.UserInfo,
	users ...*usershort.UserShort,
) {
	enrichJobInfo := h.enricher.EnrichJobInfoSafeFunc(userInfo)
	for _, user := range users {
		if user != nil {
			enrichJobInfo(user)
		}
	}
}

func (h *Hydrator) loadUsers(
	ctx context.Context,
	db qrm.DB,
	userIDs []int32,
	fields resolvedFields,
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
	if fields.phoneNumber {
		selectors = append(selectors, tUser.PhoneNumber)
	}

	userIDExprs := make([]mysql.Expression, 0, len(userIDs))
	for _, userID := range userIDs {
		userIDExprs = append(userIDExprs, mysql.Int32(userID))
	}

	stmt := tUser.
		SELECT(selectors[0], selectors[1:]...).
		FROM(tUser).
		WHERE(mysql.AND(
			tUser.ID.IN(userIDExprs...),
			tUser.DeletedAt.IS_NULL(),
		)).
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

func (h *Hydrator) loadBasicUsers(
	ctx context.Context,
	db qrm.DB,
	userIDs []int32,
	fields resolvedFields,
) (map[int32]*usershort.UserShort, []int32, error) {
	tUser := table.FivenetUser.AS("user_short")
	selectors := []mysql.Projection{
		tUser.ID,
		tUser.Firstname,
		tUser.Lastname,
		tUser.Job,
		tUser.JobGrade,
	}
	if fields.phoneNumber {
		selectors = append(selectors, tUser.PhoneNumber)
	}

	userIDExprs := make([]mysql.Expression, 0, len(userIDs))
	for _, userID := range userIDs {
		userIDExprs = append(userIDExprs, mysql.Int32(userID))
	}

	stmt := tUser.
		SELECT(selectors[0], selectors[1:]...).
		FROM(tUser).
		WHERE(mysql.AND(
			tUser.ID.IN(userIDExprs...),
			tUser.DeletedAt.IS_NULL(),
		)).
		LIMIT(int64(len(userIDs)))

	dest := []*usershort.UserShort{}
	if err := stmt.QueryContext(ctx, db, &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return map[int32]*usershort.UserShort{}, userIDs, nil
		}
		return nil, nil, fmt.Errorf("failed to load basic citizens by id: %w", err)
	}

	usersByID := make(map[int32]*usershort.UserShort, len(dest))
	for _, user := range dest {
		if user != nil && user.GetUserId() > 0 {
			usersByID[user.GetUserId()] = user
		}
	}
	return usersByID, userIDs, nil
}

func (h *Hydrator) loadProps(
	ctx context.Context,
	db qrm.DB,
	userInfo *userinfo.UserInfo,
	usersByID map[int32]*users.User,
	orderedIDs []int32,
	fields resolvedFields,
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

	propsByUserID, err := h.loadUserProps(ctx, db, userIDs, fields)
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

		if fields.job {
			props.Default()
		}
		if fields.job && props.JobName != nil {
			grade := props.GetJobGradeNumber()
			props.Job, props.JobGrade = h.enricher.GetJobGrade(props.GetJobName(), grade)
		}

		if fields.hasProps() {
			user.Props = props
		}

		if fields.job && props.JobName != nil && !slices.Contains(publicJobs, user.GetJob()) {
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

		if fields.labels {
			labels, err := h.store.GetUserLabelsForUser(ctx, db, userInfo, userID)
			if err != nil {
				return err
			}
			user.Props.Labels = labels
		}

		if fields.licenses {
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
	fields resolvedFields,
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

	selectors := []mysql.Projection{tUserProps.UserID}
	if fields.wanted {
		selectors = append(selectors, tUserProps.Wanted)
	}
	if fields.job {
		selectors = append(selectors, tUserProps.Job, tUserProps.JobGrade)
	}
	if fields.trafficInfractionPoints {
		selectors = append(
			selectors,
			tUserProps.TrafficInfractionPoints,
			tUserProps.TrafficInfractionPointsUpdatedAt,
		)
	}
	if fields.openFines {
		selectors = append(selectors, tUserProps.OpenFines)
	}
	if fields.bloodType {
		selectors = append(selectors, tUserProps.BloodType)
	}
	if fields.mugshot {
		selectors = append(
			selectors,
			tUserProps.MugshotFileID,
			tFiles.ID.AS("mugshot.mugshot_file_id"),
			tFiles.FilePath,
		)
	}
	if fields.email {
		selectors = append(selectors, tUserProps.Email)
	}

	stmt := tUserProps.
		SELECT(selectors[0], selectors[1:]...).
		FROM(
			tUserProps.LEFT_JOIN(tFiles, tFiles.ID.EQ(tUserProps.MugshotFileID)),
		).
		WHERE(tUserProps.UserID.IN(userIDExprs...)).
		LIMIT(int64(len(userIDs)))

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
