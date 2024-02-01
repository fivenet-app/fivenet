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
	"github.com/galexrt/fivenet/gen/go/proto/resources/permissions"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/galexrt/fivenet/pkg/perms/collections"
	"github.com/galexrt/fivenet/query/fivenet/model"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/nats-io/nats.go"
	"github.com/puzpuzpuz/xsync/v3"
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

	Can(userInfo *userinfo.UserInfo, category Category, name Name) bool

	LookupAttributeByID(id uint64) (*cacheAttr, bool)
	GetAttribute(category Category, name Name, key Key) (*permissions.RoleAttribute, error)
	GetAttributeByIDs(ctx context.Context, ids ...uint64) ([]*permissions.RoleAttribute, error)
	CreateAttribute(ctx context.Context, permId uint64, key Key, aType permissions.AttributeTypes, validValues any, defaultValues any) (uint64, error)
	UpdateAttribute(ctx context.Context, attributeId uint64, permId uint64, key Key, aType permissions.AttributeTypes, validValues any, defaultValues any) error
	GetRoleAttributes(job string, grade int32) ([]*permissions.RoleAttribute, error)
	GetRoleAttributeByID(roleId uint64, attrId uint64) (*permissions.RoleAttribute, bool)
	FlattenRoleAttributes(job string, grade int32) ([]string, error)
	GetAllAttributes(ctx context.Context, job string, grade int32) ([]*permissions.RoleAttribute, error)
	AddOrUpdateAttributesToRole(ctx context.Context, job string, grade int32, roleId uint64, attrs ...*permissions.RoleAttribute) error
	RemoveAttributesFromRole(ctx context.Context, roleId uint64, attrs ...*permissions.RoleAttribute) error
	ClearAttributeFromRole(ctx context.Context, roleId uint64, attrs ...*permissions.RoleAttribute) error
	UpdateRoleAttributeMaxValues(ctx context.Context, roleId uint64, attrId uint64, maxValues *permissions.AttributeValues) error
	GetClosestRoleAttrMaxVals(job string, grade int32, permId uint64, key Key) (*permissions.AttributeValues, uint64)

	Attr(userInfo *userinfo.UserInfo, category Category, name Name, key Key) (any, error)
}

type userCacheKey struct {
	userId int32
	permId uint64
}

type Perms struct {
	logger *zap.Logger
	db     *sql.DB
	cfg    *config.Config
	wg     sync.WaitGroup

	tracer trace.Tracer
	ctx    context.Context

	js    nats.JetStreamContext
	jsSub *nats.Subscription

	permsMap *xsync.MapOf[uint64, *cachePerm]
	// Guard name to permission ID
	permsGuardToIDMap *xsync.MapOf[string, uint64]
	// Job name to map of grade numbers to role ID
	permsJobsRoleMap *xsync.MapOf[string, *xsync.MapOf[int32, uint64]]
	// Role ID to map of permissions ID and result
	permsRoleMap *xsync.MapOf[uint64, *xsync.MapOf[uint64, bool]]

	// Attribute map (key: ID of attribute)
	attrsMap *xsync.MapOf[uint64, *cacheAttr]
	// Role ID to map of role attributes
	attrsRoleMap *xsync.MapOf[uint64, *xsync.MapOf[uint64, *cacheRoleAttr]]
	// Perm ID to map Key -> cached attribute
	attrsPermsMap *xsync.MapOf[uint64, *xsync.MapOf[string, uint64]]

	userCanCacheTTL time.Duration
	userCanCache    *cache.Cache[userCacheKey, bool]
}

type JobPermission struct {
	PermissionID uint64
	Val          bool
}

type Params struct {
	fx.In

	LC     fx.Lifecycle
	Logger *zap.Logger
	DB     *sql.DB
	TP     *tracesdk.TracerProvider
	JS     nats.JetStreamContext
	Config *config.Config
}

func New(p Params) (Permissions, error) {
	ctx, cancel := context.WithCancel(context.Background())

	userCanCache := cache.NewContext(
		ctx,
		cache.AsLRU[userCacheKey, bool](lru.WithCapacity(1024)),
		cache.WithJanitorInterval[userCacheKey, bool](15*time.Second),
	)

	ps := &Perms{
		logger: p.Logger,
		db:     p.DB,
		cfg:    p.Config,
		wg:     sync.WaitGroup{},

		tracer: p.TP.Tracer("perms"),
		ctx:    ctx,

		js: p.JS,

		permsMap:          xsync.NewMapOf[uint64, *cachePerm](),
		permsGuardToIDMap: xsync.NewMapOf[string, uint64](),
		permsJobsRoleMap:  xsync.NewMapOf[string, *xsync.MapOf[int32, uint64]](),
		permsRoleMap:      xsync.NewMapOf[uint64, *xsync.MapOf[uint64, bool]](),

		attrsMap:      xsync.NewMapOf[uint64, *cacheAttr](),
		attrsRoleMap:  xsync.NewMapOf[uint64, *xsync.MapOf[uint64, *cacheRoleAttr]](),
		attrsPermsMap: xsync.NewMapOf[uint64, *xsync.MapOf[string, uint64]](),

		userCanCacheTTL: 30 * time.Second,
		userCanCache:    userCanCache,
	}

	p.LC.Append(fx.StartHook(func(ctx context.Context) error {
		if err := ps.init(ctx); err != nil {
			return err
		}

		ps.wg.Add(1)
		go func() {
			defer ps.wg.Done()
			if err := ps.ApplyJobPermissions(ctx, ""); err != nil {
				ps.logger.Error("failed to apply job permissions", zap.Error(err))
				return
			}
		}()

		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		ps.wg.Wait()

		cancel()
		return ps.stop()
	}))

	return ps, nil
}

func (p *Perms) init(ctx context.Context) error {
	ctx, span := p.tracer.Start(ctx, "perms-init")
	defer span.End()

	if err := p.load(); err != nil {
		return err
	}

	if err := p.registerEvents(ctx); err != nil {
		return err
	}

	cfgDefaultPerms := p.cfg.Game.Auth.DefaultPermissions
	defaultPerms := make([]string, len(p.cfg.Game.Auth.DefaultPermissions))
	for i := 0; i < len(p.cfg.Game.Auth.DefaultPermissions); i++ {
		defaultPerms[i] = BuildGuard(Category(cfgDefaultPerms[i].Category), Name(cfgDefaultPerms[i].Name))
	}

	if err := p.register(ctx, defaultPerms); err != nil {
		return fmt.Errorf("failed to register permissions. %w", err)
	}

	return nil
}

type cachePerm struct {
	ID        uint64
	Category  Category
	Name      Name
	GuardName string
}

type cacheAttr struct {
	ID            uint64
	PermissionID  uint64
	Category      Category
	Name          Name
	Key           Key
	Type          permissions.AttributeTypes
	ValidValues   *permissions.AttributeValues
	DefaultValues *permissions.AttributeValues
}

type cacheRoleAttr struct {
	AttrID       uint64
	PermissionID uint64
	Key          Key
	Type         permissions.AttributeTypes
	Value        *permissions.AttributeValues
	Max          *permissions.AttributeValues
}

func (p *Perms) load() error {
	ctx, span := p.tracer.Start(p.ctx, "perms-load")
	defer span.End()

	if err := p.loadPermissions(ctx); err != nil {
		return err
	}

	if err := p.loadAttributes(ctx); err != nil {
		return err
	}

	if err := p.loadRoleIDs(ctx); err != nil {
		return err
	}

	if err := p.loadRolePermissions(ctx, 0); err != nil {
		return err
	}

	if err := p.loadRoleAttributes(ctx, 0); err != nil {
		return err
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
		return err
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
			tAttrs.DefaultValues.AS("default_values"),
		).
		FROM(tAttrs)

	var dest []struct {
		ID            uint64
		PermissionID  uint64
		Key           Key
		Type          permissions.AttributeTypes
		ValidValues   string
		DefaultValues string
	}
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		return err
	}

	for _, attr := range dest {
		if err := p.addOrUpdateAttributeInMap(attr.PermissionID, attr.ID, attr.Key, attr.Type, attr.ValidValues, attr.DefaultValues); err != nil {
			return err
		}
	}

	return nil
}

func (p *Perms) loadRoleIDs(ctx context.Context) error {
	stmt := tRoles.
		SELECT(
			tRoles.ID.AS("id"),
			tRoles.Job.AS("job"),
			tRoles.Grade.AS("grade"),
		).
		FROM(tRoles)

	var dest []struct {
		ID    uint64
		Job   string
		Grade int32
	}
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		return err
	}

	for _, role := range dest {
		grades, _ := p.permsJobsRoleMap.LoadOrCompute(role.Job, xsync.NewMapOf[int32, uint64])
		grades.Store(role.Grade, role.ID)
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
				tRoles.ID.EQ(tRolePerms.RoleID)),
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
		return err
	}

	for _, rolePerms := range dest {
		perms, _ := p.permsRoleMap.LoadOrCompute(rolePerms.RoleID, xsync.NewMapOf[uint64, bool])
		perms.Store(rolePerms.ID, rolePerms.Val)
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
			tRoleAttrs.MaxValues.AS("max_values"),
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
		Value        string
		MaxValues    string
	}

	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	for _, ra := range dest {
		a, ok := p.LookupAttributeByID(ra.AttrID)
		if !ok {
			return fmt.Errorf("unable to find attribute ID %d for role %d", ra.AttrID, ra.RoleID)
		}

		if err := p.addOrUpdateRoleAttributeInMap(ra.RoleID, ra.PermissionID, ra.AttrID, ra.Key, ra.Type, ra.Value, ra.MaxValues); err != nil {
			// Reset the attribute value to null/ empty because we encountered an error
			if err := p.addOrUpdateAttributesToRole(p.ctx, ra.RoleID, &permissions.RoleAttribute{
				RoleId: ra.RoleID,
				AttrId: ra.AttrID,
				Value:  a.DefaultValues,
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *Perms) addOrUpdateRoleAttributeInMap(roleId uint64, permId uint64, attrId uint64, key Key, aType permissions.AttributeTypes, rawValue string, rawMax string) error {
	val := &permissions.AttributeValues{}
	if err := p.convertRawValue(val, rawValue, aType); err != nil {
		return err
	}

	var max *permissions.AttributeValues
	if rawMax != "" {
		max = &permissions.AttributeValues{}
		if err := p.convertRawValue(max, rawMax, aType); err != nil {
			return err
		}
	}

	p.updateRoleAttributeInMap(roleId, permId, attrId, key, aType, val, max)

	return nil
}

func (p *Perms) updateRoleAttributeInMap(roleId uint64, permId uint64, attrId uint64, key Key, aType permissions.AttributeTypes, value *permissions.AttributeValues, max *permissions.AttributeValues) {
	attrRoleMap, _ := p.attrsRoleMap.LoadOrCompute(roleId, xsync.NewMapOf[uint64, *cacheRoleAttr])
	cached, ok := attrRoleMap.Load(attrId)
	if !ok {
		attrRoleMap.Store(attrId, &cacheRoleAttr{
			AttrID:       attrId,
			PermissionID: permId,
			Key:          key,
			Type:         aType,
			Value:        value,
			Max:          max,
		})
	} else {
		cached.PermissionID = permId
		cached.Key = key
		cached.Type = aType
		cached.Value = value
		cached.Max = max
	}
}

func (p *Perms) removeRoleAttributeFromMap(roleId uint64, attrId uint64) {
	attrMap, ok := p.attrsRoleMap.Load(roleId)
	if !ok {
		return
	}

	attrMap.Delete(attrId)
}

func (p *Perms) stop() error {
	return p.jsSub.Unsubscribe()
}
