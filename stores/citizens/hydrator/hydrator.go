package citizenshydrator

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"

	citizenslicenses "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/citizens/licenses"
	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	users "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users"
	usersprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/props"
	usershort "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/short"
	permscitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2026/pkg/perms"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/model"
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

func resolvedFieldsForProperties(
	propertyList *permissionsattributes.PropertyList,
	enabled map[string]bool,
) resolvedFields {
	canUse := func(key string) bool {
		return propertyList.Contains(key) && enabled[key]
	}
	return resolvedFields{
		phoneNumber:             canUse("user.phone_number"),
		licenses:                canUse("user.licenses"),
		wanted:                  canUse("user.wanted"),
		job:                     canUse("user.job"),
		trafficInfractionPoints: canUse("user.traffic_infraction_points"),
		openFines:               canUse("user.open_fines"),
		bloodType:               canUse("user.blood_type"),
		mugshot:                 canUse("user.mugshot"),
		labels:                  canUse("user.labels"),
		email:                   canUse("user.email"),
	}
}

func (h *Hydrator) resolveFields(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
) (resolvedFields, error) {
	if userInfo == nil {
		return resolvedFields{}, nil
	}

	propertyList, err := h.perms.AttrPropertyListContext(ctx,
		userInfo,
		permscitizens.CitizensService.ListCitizens.Properties,
	)
	if err != nil {
		return resolvedFields{}, err
	}

	keys := propertyList.GetKeys()
	enabled := make(map[string]bool, len(keys))
	if h.store == nil {
		for _, key := range keys {
			enabled[key] = true
		}
	} else {
		definitions, err := h.store.ResolvePropertyDefinitions(ctx, keys)
		if err != nil {
			return resolvedFields{}, err
		}
		for key, definition := range definitions {
			enabled[key] = definition.GetEnabled()
		}
	}
	return resolvedFieldsForProperties(propertyList, enabled), nil
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

	fields, err := h.resolveFields(ctx, userInfo)
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

	fields, err := h.resolveFields(ctx, userInfo)
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

	userIDExprs := make([]mysql.Expression, 0, len(userIDs))
	for _, userID := range userIDs {
		userIDExprs = append(userIDExprs, mysql.Int32(userID))
	}

	keys := make([]mysql.Expression, 0, 7)
	if fields.wanted {
		keys = append(keys, mysql.String("user.wanted"))
	}
	if fields.job {
		keys = append(keys, mysql.String("user.job"), mysql.String("user.job_grade"))
	}
	if fields.trafficInfractionPoints {
		keys = append(keys, mysql.String("user.traffic_infraction_points"))
	}
	if fields.openFines {
		keys = append(keys, mysql.String("user.open_fines"))
	}
	if fields.bloodType {
		keys = append(keys, mysql.String("user.blood_type"))
	}
	if fields.email {
		keys = append(keys, mysql.String("user.email"))
	}

	propsByUserID := map[int32]*usersprops.UserProps{}
	ensure := func(userID int32) *usersprops.UserProps {
		if propsByUserID[userID] == nil {
			propsByUserID[userID] = &usersprops.UserProps{UserId: userID}
		}
		return propsByUserID[userID]
	}
	if len(keys) > 0 {
		tValues := table.FivenetUserPropertyValues
		rows := []model.FivenetUserPropertyValues{}
		stmt := tValues.SELECT(tValues.AllColumns).
			FROM(tValues).
			WHERE(mysql.AND(tValues.UserID.IN(userIDExprs...), tValues.PropertyKey.IN(keys...)))
		if err := stmt.QueryContext(ctx, db, &rows); err != nil && !errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("failed to load citizen property values: %w", err)
		}
		for _, row := range rows {
			props := ensure(row.UserID)
			switch row.PropertyKey {
			case "user.wanted":
				props.Wanted = row.ValueBool
			case "user.job":
				props.JobName = row.ValueString
			case "user.job_grade":
				if row.ValueInteger != nil {
					grade := int32(*row.ValueInteger)
					props.JobGradeNumber = &grade
				}
			case "user.traffic_infraction_points":
				if row.ValueInteger != nil {
					points := uint32(*row.ValueInteger)
					props.TrafficInfractionPoints = &points
				}
			case "user.open_fines":
				props.OpenFines = row.ValueInteger
			case "user.blood_type":
				props.BloodType = row.ValueString
			case "user.email":
				props.Email = row.ValueString
			}
		}
	}
	if fields.wanted || fields.trafficInfractionPoints {
		tState := table.FivenetUserPropertyState
		states := []model.FivenetUserPropertyState{}
		stmt := tState.SELECT(tState.AllColumns).
			FROM(tState).
			WHERE(tState.UserID.IN(userIDExprs...))
		if err := stmt.QueryContext(
			ctx,
			db,
			&states,
		); err != nil &&
			!errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("failed to load citizen property state: %w", err)
		}
		for _, state := range states {
			props := ensure(state.UserID)
			if state.PropertyKey == "user.wanted" && fields.wanted {
				if state.ActiveAt != nil {
					props.WantedAt = timestamp.New(*state.ActiveAt)
				}
				if state.ExpiresAt != nil {
					props.WantedTill = timestamp.New(*state.ExpiresAt)
				}
			}
			if state.PropertyKey == "user.traffic_infraction_points" &&
				fields.trafficInfractionPoints {
				props.TrafficInfractionPointsUpdatedAt = timestamp.New(state.UpdatedAt)
			}
		}
	}
	if fields.mugshot {
		tUserFiles := table.FivenetUserFiles.AS("mugshot_file")
		tFiles := table.FivenetFiles.AS("mugshot")
		mugshots := []*usersprops.UserProps{}
		stmt := tUserFiles.SELECT(
			tUserFiles.UserID,
			tUserFiles.FileID.AS("mugshot_file_id"),
			tFiles.ID.AS("mugshot.mugshot_file_id"),
			tFiles.FilePath,
		).FROM(tUserFiles.LEFT_JOIN(tFiles, tFiles.ID.EQ(tUserFiles.FileID))).WHERE(mysql.AND(
			tUserFiles.UserID.IN(userIDExprs...),
			tUserFiles.Kind.EQ(mysql.String("mugshot")),
		))
		if err := stmt.QueryContext(
			ctx,
			db,
			&mugshots,
		); err != nil &&
			!errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("failed to load citizen mugshots: %w", err)
		}
		for _, mugshot := range mugshots {
			if mugshot == nil || mugshot.GetUserId() <= 0 {
				continue
			}
			props := ensure(mugshot.GetUserId())
			props.MugshotFileId = mugshot.MugshotFileId
			props.Mugshot = mugshot.Mugshot
		}
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
