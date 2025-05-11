package perms

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	cache "github.com/Code-Hex/go-generics-cache"
	"github.com/Code-Hex/go-generics-cache/policy/lru"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/permissions"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms/collections"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/model"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/puzpuzpuz/xsync/v4"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Permissions interface {
	GetAllPermissions(ctx context.Context) ([]*permissions.Permission, error)
	GetPermissionsByIDs(ctx context.Context, ids ...uint64) ([]*permissions.Permission, error)
	GetPermission(ctx context.Context, category Category, name Name) (*permissions.Permission, error)
	CreatePermission(ctx context.Context, category Category, name Name) (uint64, error)
	GetPermissionsOfUser(userInfo *userinfo.UserInfo) (collections.Permissions, error)

	GetRoles(ctx context.Context, excludeSystem bool) (collections.Roles, error)
	GetJobRoles(ctx context.Context, job string) (collections.Roles, error)
	GetJobRolesUpTo(ctx context.Context, job string, grade int32) (collections.Roles, error)
	GetClosestJobRole(ctx context.Context, job string, grade int32) (*model.FivenetRoles, error)
	CountRolesForJob(ctx context.Context, prefix string) (int64, error)

	GetRole(ctx context.Context, id uint64) (*model.FivenetRoles, error)
	GetRoleByJobAndGrade(ctx context.Context, job string, grade int32) (*model.FivenetRoles, error)
	GetRolePermissions(ctx context.Context, id uint64) ([]*permissions.Permission, error)

	CreateRole(ctx context.Context, job string, grade int32) (*model.FivenetRoles, error)
	DeleteRole(ctx context.Context, id uint64) error
	UpdateRolePermissions(ctx context.Context, id uint64, perms ...AddPerm) error
	RemovePermissionsFromRole(ctx context.Context, id uint64, perms ...uint64) error
	GetJobPermissions(ctx context.Context, job string) ([]*permissions.Permission, error)
	UpdateJobPermissions(ctx context.Context, job string, id uint64, val bool) error
	ApplyJobPermissions(ctx context.Context, job string) error
	ClearJobPermissions(ctx context.Context, job string) error

	Can(userInfo *userinfo.UserInfo, category Category, name Name) bool

	LookupAttributeByID(id uint64) (*cacheAttr, bool)
	GetAttribute(category Category, name Name, key Key) (*permissions.RoleAttribute, error)
	GetAttributeByIDs(ctx context.Context, ids ...uint64) ([]*permissions.RoleAttribute, error)
	CreateAttribute(ctx context.Context, permId uint64, key Key, aType permissions.AttributeTypes, validValues *permissions.AttributeValues) (uint64, error)
	UpdateAttribute(ctx context.Context, attributeId uint64, permId uint64, key Key, aType permissions.AttributeTypes, validValues *permissions.AttributeValues) error
	GetRoleAttributes(job string, grade int32) ([]*permissions.RoleAttribute, error)
	GetRoleAttributeByID(roleId uint64, attrId uint64) (*permissions.RoleAttribute, bool)
	FlattenRoleAttributes(job string, grade int32) ([]string, error)
	GetAllAttributes(ctx context.Context, job string, grade int32) ([]*permissions.RoleAttribute, error)
	AddOrUpdateAttributesToRole(ctx context.Context, job string, roleId uint64, attrs ...*permissions.RoleAttribute) error
	RemoveAttributesFromRole(ctx context.Context, roleId uint64, attrs ...*permissions.RoleAttribute) error

	GetJobAttrMaxVals(job string, attrId uint64) (*permissions.AttributeValues, bool)
	UpdateJobAttributeMaxValues(ctx context.Context, job string, attrId uint64, maxValues *permissions.AttributeValues) error
	ClearJobAttributes(ctx context.Context, job string) error

	Attr(userInfo *userinfo.UserInfo, category Category, name Name, key Key) (*permissions.AttributeValues, error)
	AttrStringList(userInfo *userinfo.UserInfo, category Category, name Name, key Key) (*permissions.StringList, error)
	AttrJobList(userInfo *userinfo.UserInfo, category Category, name Name, key Key) (*permissions.StringList, error)
	AttrJobGradeList(userInfo *userinfo.UserInfo, category Category, name Name, key Key) (*permissions.JobGradeList, error)

	SetDefaultRolePerms(ctx context.Context, defaultPerms []string) error
}

type userCacheKey struct {
	userId int32
	permId uint64
}

type Perms struct {
	logger *zap.Logger
	db     *sql.DB
	wg     sync.WaitGroup

	tracer trace.Tracer

	js     *events.JSWrapper
	jsCons jetstream.ConsumeContext

	// Enable fast "init" mode, skipping clean ups. Should only be used for dev/test environments
	devMode       bool
	startJobGrade int32

	permsMap *xsync.Map[uint64, *cachePerm]
	// Guard name to permission ID
	permsGuardToIDMap *xsync.Map[string, uint64]
	// Job name to map of grade numbers to role ID
	permsJobsRoleMap *xsync.Map[string, *xsync.Map[int32, uint64]]
	// Role ID to map of permissions ID and result
	permsRoleMap *xsync.Map[uint64, *xsync.Map[uint64, bool]]
	// Role ID to Job map
	roleIDToJobMap *xsync.Map[uint64, string]

	// Attribute map (key: ID of attribute)
	attrsMap *xsync.Map[uint64, *cacheAttr]
	// Role ID to map of role attributes
	attrsRoleMap *xsync.Map[uint64, *xsync.Map[uint64, *cacheRoleAttr]]
	// Perm ID to map `Key` -> cached attribute
	attrsPermsMap *xsync.Map[uint64, *xsync.Map[string, uint64]]
	// Job to map attr ID to job max value attribute
	attrsJobMaxValuesMap *xsync.Map[string, *xsync.Map[uint64, *permissions.AttributeValues]]

	userCanCacheTTL time.Duration
	userCanCache    *cache.Cache[userCacheKey, bool]
}

type JobPermission struct {
	PermissionID uint64
	Val          bool
}

type Params struct {
	fx.In

	LC        fx.Lifecycle
	Logger    *zap.Logger
	DB        *sql.DB
	TP        *tracesdk.TracerProvider
	JS        *events.JSWrapper
	Cfg       *config.Config
	AppConfig appconfig.IConfig
}

func New(p Params) (Permissions, error) {
	ctxCancel, cancel := context.WithCancel(context.Background())

	userCanCache := cache.NewContext(
		ctxCancel,
		cache.AsLRU[userCacheKey, bool](lru.WithCapacity(1024)),
		cache.WithJanitorInterval[userCacheKey, bool](15*time.Second),
	)

	ps := &Perms{
		logger: p.Logger,
		db:     p.DB,
		wg:     sync.WaitGroup{},

		tracer: p.TP.Tracer("perms"),

		js: p.JS,

		devMode:       false,
		startJobGrade: p.Cfg.Game.StartJobGrade,

		permsMap:          xsync.NewMap[uint64, *cachePerm](),
		permsGuardToIDMap: xsync.NewMap[string, uint64](),
		permsJobsRoleMap:  xsync.NewMap[string, *xsync.Map[int32, uint64]](),
		permsRoleMap:      xsync.NewMap[uint64, *xsync.Map[uint64, bool]](),
		roleIDToJobMap:    xsync.NewMap[uint64, string](),

		attrsMap:             xsync.NewMap[uint64, *cacheAttr](),
		attrsRoleMap:         xsync.NewMap[uint64, *xsync.Map[uint64, *cacheRoleAttr]](),
		attrsPermsMap:        xsync.NewMap[uint64, *xsync.Map[string, uint64]](),
		attrsJobMaxValuesMap: xsync.NewMap[string, *xsync.Map[uint64, *permissions.AttributeValues]](),

		userCanCacheTTL: 30 * time.Second,
		userCanCache:    userCanCache,
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		cfgDefaultPerms := p.AppConfig.Get().Perms.Default
		defaultPerms := make([]string, len(cfgDefaultPerms))
		for i := range cfgDefaultPerms {
			defaultPerms[i] = BuildGuard(Category(cfgDefaultPerms[i].Category), Name(cfgDefaultPerms[i].Name))
		}

		if err := ps.load(ctxStartup); err != nil {
			return err
		}
		ps.logger.Debug("permissions loaded")

		if err := ps.registerSubscriptions(ctxStartup, ctxCancel); err != nil {
			return fmt.Errorf("failed to register events subscriptions. %w", err)
		}
		ps.logger.Debug("registered events subscription")

		if err := ps.register(ctxStartup, defaultPerms); err != nil {
			return fmt.Errorf("failed to register permissions. %w", err)
		}

		// Skip apply job perms when in dev mode
		if !ps.devMode {
			ps.wg.Add(1)
			go func() {
				defer ps.wg.Done()

				if err := ps.ApplyJobPermissions(ctxCancel, ""); err != nil {
					ps.logger.Error("failed to apply job permissions", zap.Error(err))
					return
				}
				ps.logger.Debug("successfully applied job permissions")
			}()
		}

		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		ps.wg.Wait()

		if ps.jsCons != nil {
			ps.jsCons.Stop()
			ps.jsCons = nil
		}

		return nil
	}))

	return ps, nil
}

type cachePerm struct {
	ID        uint64
	Category  Category
	Name      Name
	GuardName string
}

type cacheAttr struct {
	ID           uint64
	PermissionID uint64
	Category     Category
	Name         Name
	Key          Key
	Type         permissions.AttributeTypes
	ValidValues  *permissions.AttributeValues
}

type cacheRoleAttr struct {
	Job          string
	AttrID       uint64
	PermissionID uint64
	Key          Key
	Type         permissions.AttributeTypes
	Value        *permissions.AttributeValues
}

func (p *Perms) load(ctx context.Context) error {
	ctx, span := p.tracer.Start(ctx, "perms-load")
	defer span.End()

	if err := p.loadPermissions(ctx); err != nil {
		return fmt.Errorf("failed to load permissions. %w", err)
	}

	if err := p.loadAttributes(ctx); err != nil {
		return fmt.Errorf("failed to load attributes. %w", err)
	}

	if err := p.loadRoles(ctx, 0); err != nil {
		return fmt.Errorf("failed to load roles. %w", err)
	}

	if err := p.loadRolePermissions(ctx, 0); err != nil {
		return fmt.Errorf("failed to load role permissions. %w", err)
	}

	if err := p.loadRoleAttributes(ctx, 0); err != nil {
		return fmt.Errorf("failed to load role attributes. %w", err)
	}

	if err := p.loadJobAttrs(ctx, ""); err != nil {
		return fmt.Errorf("failed to load job attributes. %w", err)
	}

	return nil
}

func (p *Perms) loadPermissions(ctx context.Context) error {
	tPerms := tPerms.AS("cachePerm")
	stmt := tPerms.
		SELECT(
			tPerms.ID,
			tPerms.Category,
			tPerms.Name,
			tPerms.GuardName,
		).
		FROM(tPerms)

	var dest []*cachePerm
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return fmt.Errorf("failed to query permissions. %w", err)
		}
	}

	for _, perm := range dest {
		p.permsMap.Store(perm.ID, &cachePerm{
			ID:        perm.ID,
			Category:  perm.Category,
			Name:      perm.Name,
			GuardName: BuildGuard(perm.Category, perm.Name),
		})
		p.permsGuardToIDMap.Store(BuildGuard(perm.Category, perm.Name), perm.ID)
	}

	return nil
}

func (p *Perms) loadAttributes(ctx context.Context) error {
	stmt := tAttrs.
		SELECT(
			tAttrs.ID.AS("id"),
			tAttrs.PermissionID.AS("permission_id"),
			tAttrs.Key.AS("key"),
			tAttrs.Type.AS("type"),
			tAttrs.ValidValues.AS("valid_values"),
		).
		FROM(tAttrs)

	var dest []struct {
		ID           uint64
		PermissionID uint64
		Key          Key
		Type         permissions.AttributeTypes
		ValidValues  *permissions.AttributeValues
	}
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return fmt.Errorf("failed to query attributes. %w", err)
		}
	}

	for _, attr := range dest {
		if err := p.addOrUpdateAttributeInMap(attr.PermissionID, attr.ID, attr.Key, attr.Type, attr.ValidValues); err != nil {
			return fmt.Errorf("failed to add/update attribute in map. %w", err)
		}
	}

	return nil
}

func (p *Perms) loadRoles(ctx context.Context, id uint64) error {
	stmt := tRoles.
		SELECT(
			tRoles.ID.AS("id"),
			tRoles.Job.AS("job"),
			tRoles.Grade.AS("grade"),
		).
		FROM(tRoles)

	if id != 0 {
		stmt = stmt.
			WHERE(tRoles.ID.EQ(jet.Uint64(id)))
	}

	var dest []struct {
		ID    uint64
		Job   string
		Grade int32
	}
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return fmt.Errorf("failed to query roles. %w", err)
		}
	}

	for _, role := range dest {
		grades, _ := p.permsJobsRoleMap.LoadOrCompute(role.Job, func() (*xsync.Map[int32, uint64], bool) {
			return xsync.NewMap[int32, uint64](), false
		})
		grades.Store(role.Grade, role.ID)

		p.roleIDToJobMap.Store(role.ID, role.Job)
	}

	return nil
}

func (p *Perms) loadRolePermissions(ctx context.Context, roleId uint64) error {
	stmt := tRolePerms.
		SELECT(
			tRolePerms.RoleID.AS("role_id"),
			tRolePerms.PermissionID.AS("id"),
			tRolePerms.Val.AS("val"),
		).
		FROM(tRolePerms.
			INNER_JOIN(tRoles,
				tRoles.ID.EQ(tRolePerms.RoleID),
			),
		).
		ORDER_BY(
			tRoles.Job.ASC(),
			tRoles.Grade.DESC(),
		)

	if roleId != 0 {
		stmt = stmt.WHERE(
			tRoles.ID.EQ(jet.Uint64(roleId)),
		)
	}

	var dest []struct {
		RoleID uint64
		ID     uint64
		Val    bool
	}
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return fmt.Errorf("failed to query role permissions. %w", err)
		}
	}

	for _, rolePerms := range dest {
		perms, _ := p.permsRoleMap.LoadOrCompute(rolePerms.RoleID, func() (*xsync.Map[uint64, bool], bool) {
			return xsync.NewMap[uint64, bool](), false
		})
		perms.Store(rolePerms.ID, rolePerms.Val)
	}

	return nil
}

func (p *Perms) loadJobAttrs(ctx context.Context, job string) error {
	stmt := tJobAttrs.
		SELECT(
			tJobAttrs.Job.AS("job"),
			tJobAttrs.AttrID.AS("attr_id"),
			tJobAttrs.MaxValues.AS("max_values"),
		).
		FROM(tJobAttrs).
		ORDER_BY(
			tJobAttrs.Job.ASC(),
		)

	if job != "" {
		stmt = stmt.WHERE(
			tJobAttrs.Job.EQ(jet.String(job)),
		)
	}

	var dest []struct {
		Job       string
		AttrID    uint64
		MaxValues *permissions.AttributeValues
	}
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return fmt.Errorf("failed to query job attributes. %w", err)
		}
	}

	// No attributes? Delete cached data
	if len(dest) == 0 {
		p.attrsJobMaxValuesMap.Delete(job)
	} else {
		for _, jobAttrs := range dest {
			attrs, _ := p.attrsJobMaxValuesMap.LoadOrCompute(jobAttrs.Job, func() (*xsync.Map[uint64, *permissions.AttributeValues], bool) {
				return xsync.NewMap[uint64, *permissions.AttributeValues](), false
			})
			attrs.Store(jobAttrs.AttrID, jobAttrs.MaxValues)
		}
	}

	return nil
}

func (p *Perms) loadRoleAttributes(ctx context.Context, roleId uint64) error {
	stmt := tRoleAttrs.
		SELECT(
			tRoleAttrs.AttrID.AS("attr_id"),
			tAttrs.PermissionID.AS("permission_id"),
			tRoleAttrs.RoleID.AS("role_id"),
			tAttrs.Key.AS("key"),
			tAttrs.Type.AS("type"),
			tRoleAttrs.Value.AS("value"),
		).
		FROM(
			tRoleAttrs.
				INNER_JOIN(tAttrs,
					tAttrs.ID.EQ(tRoleAttrs.AttrID),
				),
		)

	if roleId != 0 {
		stmt = stmt.WHERE(
			tRoleAttrs.RoleID.EQ(jet.Uint64(roleId)),
		)
	}

	var dest []struct {
		AttrID       uint64
		PermissionID uint64
		RoleID       uint64
		Key          Key
		Type         permissions.AttributeTypes
		Value        *permissions.AttributeValues
	}

	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return fmt.Errorf("failed to query role attributes. %w", err)
		}
	}

	for _, ra := range dest {
		p.updateRoleAttributeInMap(ra.RoleID, ra.PermissionID, ra.AttrID, ra.Key, ra.Type, ra.Value)
	}

	return nil
}

func (p *Perms) updateRoleAttributeInMap(roleId uint64, permId uint64, attrId uint64, key Key, aType permissions.AttributeTypes, value *permissions.AttributeValues) {
	job, ok := p.lookupJobForRoleID(roleId)
	if !ok {
		p.logger.Error("unable to resolve role id to job", zap.Uint64("role_id", roleId))
		return
	}

	attrRoleMap, _ := p.attrsRoleMap.LoadOrCompute(roleId, func() (*xsync.Map[uint64, *cacheRoleAttr], bool) {
		return xsync.NewMap[uint64, *cacheRoleAttr](), false
	})
	cached, ok := attrRoleMap.Load(attrId)
	if !ok {
		attrRoleMap.Store(attrId, &cacheRoleAttr{
			Job:          job,
			AttrID:       attrId,
			PermissionID: permId,
			Key:          key,
			Type:         aType,
			Value:        value,
		})
	} else {
		cached.Job = job
		cached.PermissionID = permId
		cached.Key = key
		cached.Type = aType
		cached.Value = value
	}
}

func (p *Perms) removeRoleAttributeFromMap(roleId uint64, attrId uint64) {
	attrMap, ok := p.attrsRoleMap.Load(roleId)
	if !ok {
		return
	}

	attrMap.Delete(attrId)
}
