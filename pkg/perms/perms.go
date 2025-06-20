package perms

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/permissions"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms/collections"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/cache"
	"github.com/nats-io/nats.go"
	"github.com/puzpuzpuz/xsync/v4"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Permissions interface {
	// Permissions management
	SetDefaultRolePerms(ctx context.Context, defaultPerms []string) error
	GetAllPermissions(ctx context.Context) ([]*permissions.Permission, error)
	GetPermissionsByIDs(ctx context.Context, ids ...uint64) ([]*permissions.Permission, error)
	GetPermission(ctx context.Context, category Category, name Name) (*permissions.Permission, error)
	CreatePermission(ctx context.Context, category Category, name Name) (uint64, error)
	GetPermissionsOfUser(userInfo *userinfo.UserInfo) (collections.Permissions, error)

	// Attributes management
	GetAllAttributes(ctx context.Context) ([]*permissions.RoleAttribute, error)
	CreateAttribute(ctx context.Context, permId uint64, key Key, aType permissions.AttributeTypes, validValues *permissions.AttributeValues) (uint64, error)
	UpdateAttribute(ctx context.Context, attributeId uint64, permId uint64, key Key, aType permissions.AttributeTypes, validValues *permissions.AttributeValues) error

	// Roles management
	GetRoles(ctx context.Context, excludeSystem bool) (collections.Roles, error)
	GetRole(ctx context.Context, id uint64) (*permissions.Role, error)
	GetRoleByJobAndGrade(ctx context.Context, job string, grade int32) (*permissions.Role, error)
	GetJobRoles(ctx context.Context, job string) (collections.Roles, error)
	GetJobRolesUpTo(ctx context.Context, job string, grade int32) (collections.Roles, error)
	GetClosestJobRole(ctx context.Context, job string, grade int32) (*permissions.Role, error)
	CountRolesForJob(ctx context.Context, prefix string) (int64, error)
	CreateRole(ctx context.Context, job string, grade int32) (*permissions.Role, error)
	DeleteRole(ctx context.Context, id uint64) error
	GetRolePermissions(ctx context.Context, id uint64) ([]*permissions.Permission, error)
	GetEffectiveRolePermissions(ctx context.Context, id uint64) ([]*permissions.Permission, error)
	UpdateRolePermissions(ctx context.Context, id uint64, perms ...AddPerm) error
	RemovePermissionsFromRole(ctx context.Context, id uint64, perms ...uint64) error

	// Role Attributes management
	GetRoleAttributes(ctx context.Context, job string, grade int32) ([]*permissions.RoleAttribute, error)
	FlattenRoleAttributes(job string, grade int32) ([]string, error)
	GetEffectiveRoleAttributes(ctx context.Context, job string, grade int32) ([]*permissions.RoleAttribute, error)
	UpdateRoleAttributes(ctx context.Context, job string, roleId uint64, attrs ...*permissions.RoleAttribute) error
	RemoveAttributesFromRole(ctx context.Context, roleId uint64, attrs ...*permissions.RoleAttribute) error
	RemoveAttributesFromRoleByPermission(ctx context.Context, roleId uint64, permissionId uint64) error

	// Limit - Job permissions
	GetJobPermissions(ctx context.Context, job string) ([]*permissions.Permission, error)
	UpdateJobPermissions(ctx context.Context, job string, perms ...*permissions.PermItem) error
	ApplyJobPermissions(ctx context.Context, job string) error
	ClearJobPermissions(ctx context.Context, job string) error

	// Limit - Job attributes (max values)
	GetJobAttributes(ctx context.Context, job string) ([]*permissions.RoleAttribute, error)
	UpdateJobAttributes(ctx context.Context, job string, attrs ...*permissions.RoleAttribute) error
	ClearJobAttributes(ctx context.Context, job string) error

	// Perms Check
	Can(userInfo *userinfo.UserInfo, category Category, name Name) bool
	// Attribute retrieval/"check"
	Attr(userInfo *userinfo.UserInfo, category Category, name Name, key Key) (*permissions.AttributeValues, error)
	AttrStringList(userInfo *userinfo.UserInfo, category Category, name Name, key Key) (*permissions.StringList, error)
	AttrJobList(userInfo *userinfo.UserInfo, category Category, name Name, key Key) (*permissions.StringList, error)
	AttrJobGradeList(userInfo *userinfo.UserInfo, category Category, name Name, key Key) (*permissions.JobGradeList, error)
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

	nc    *nats.Conn
	ncSub *nats.Subscription

	cleanupRolesForMissingJobs bool
	startJobGrade              int32

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
	// Perm ID to map `Key` -> cached attribute (key is name of attribute)
	attrsPermsMap *xsync.Map[uint64, *xsync.Map[string, uint64]]

	userCanCacheTTL time.Duration
	userCanCache    *cache.LRUCache[userCacheKey, bool]
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
	NC        *nats.Conn
	Cfg       *config.Config
	AppConfig appconfig.IConfig
}

func New(p Params) (Permissions, error) {
	ctxCancel, cancel := context.WithCancel(context.Background())

	userCanCache := cache.NewLRUCache[userCacheKey, bool](1024)

	ps := &Perms{
		logger: p.Logger,
		db:     p.DB,
		wg:     sync.WaitGroup{},

		tracer: p.TP.Tracer("perms"),

		nc: p.NC,

		cleanupRolesForMissingJobs: p.Cfg.Game.CleanupRolesForMissingJobs,
		startJobGrade:              p.Cfg.Game.StartJobGrade,

		permsMap:          xsync.NewMap[uint64, *cachePerm](),
		permsGuardToIDMap: xsync.NewMap[string, uint64](),
		permsJobsRoleMap:  xsync.NewMap[string, *xsync.Map[int32, uint64]](),
		permsRoleMap:      xsync.NewMap[uint64, *xsync.Map[uint64, bool]](),
		roleIDToJobMap:    xsync.NewMap[uint64, string](),

		attrsMap:      xsync.NewMap[uint64, *cacheAttr](),
		attrsRoleMap:  xsync.NewMap[uint64, *xsync.Map[uint64, *cacheRoleAttr]](),
		attrsPermsMap: xsync.NewMap[uint64, *xsync.Map[string, uint64]](),

		userCanCacheTTL: 30 * time.Second,
		userCanCache:    userCanCache,
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		go userCanCache.StartJanitor(ctxCancel, 31)

		return ps.init(ctxCancel, ctxStartup, p)
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		ps.wg.Wait()

		if ps.ncSub != nil {
			ps.ncSub.Unsubscribe()
			ps.ncSub = nil
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
	Order     *int32
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

func (p *Perms) init(ctxCancel context.Context, ctxStartup context.Context, params Params) error {
	cfgDefaultPerms := params.AppConfig.Get().Perms.Default
	defaultPerms := make([]string, len(cfgDefaultPerms))
	for i := range cfgDefaultPerms {
		defaultPerms[i] = BuildGuard(Category(cfgDefaultPerms[i].Category), Name(cfgDefaultPerms[i].Name))
	}

	if err := p.loadData(ctxStartup); err != nil {
		return err
	}
	p.logger.Debug("permissions loaded")

	if err := p.registerSubscriptions(ctxCancel); err != nil {
		return fmt.Errorf("failed to register events subscriptions. %w", err)
	}
	p.logger.Debug("registered events subscription")

	if err := p.register(ctxStartup, defaultPerms); err != nil {
		return fmt.Errorf("failed to register permissions. %w", err)
	}

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()

		if err := p.ApplyJobPermissions(ctxCancel, ""); err != nil {
			p.logger.Error("failed to apply job permissions", zap.Error(err))
			return
		}
		p.logger.Debug("successfully applied job permissions")
	}()

	return nil
}
