package perms

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	permissionspermissions "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/permissions"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/cache"
	"github.com/nats-io/nats.go"
	"github.com/puzpuzpuz/xsync/v4"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Permissions defines the interface for permissions management.
//
// It provides methods for managing permissions, attributes, and roles,
// as well as checking permissions and attributes for users.
//
//nolint:interfacebloat // This interface is designed to be implemented by a concrete type.
type Permissions interface {
	// Permissions
	SetDefaultRolePerms(ctx context.Context, defaultPerms []string) error
	GetAllPermissions(ctx context.Context) ([]*permissionspermissions.Permission, error)
	GetPermissionsByIDs(
		ctx context.Context,
		ids ...int64,
	) ([]*permissionspermissions.Permission, error)
	GetPermission(
		ctx context.Context,
		namespace Namespace,
		service Service,
		name Name,
	) (*permissionspermissions.Permission, error)
	GetPermissionsOfUser(userInfo *userinfo.UserInfo) ([]*permissionspermissions.Permission, error)

	// Attributes
	GetAllAttributes(ctx context.Context) ([]*permissionsattributes.RoleAttribute, error)

	// Roles management
	GetRoles(ctx context.Context, excludeSystem bool) ([]*permissionspermissions.Role, error)
	GetRole(ctx context.Context, id int64) (*permissionspermissions.Role, error)
	GetRoleByJobAndGrade(
		ctx context.Context,
		job string,
		grade int32,
	) (*permissionspermissions.Role, error)
	GetJobRoles(ctx context.Context, job string) ([]*permissionspermissions.Role, error)
	GetJobRolesUpTo(
		ctx context.Context,
		job string,
		grade int32,
	) ([]*permissionspermissions.Role, error)
	GetClosestJobRole(
		ctx context.Context,
		job string,
		grade int32,
	) (*permissionspermissions.Role, error)
	CountRolesForJob(ctx context.Context, prefix string) (int64, error)
	CreateRole(ctx context.Context, job string, grade int32) (*permissionspermissions.Role, error)
	DeleteRole(ctx context.Context, roleId int64) error
	GetRolePermissions(
		ctx context.Context,
		roleId int64,
	) ([]*permissionspermissions.Permission, error)
	GetEffectiveRolePermissions(
		ctx context.Context,
		id int64,
	) ([]*permissionspermissions.Permission, error)
	UpdateRolePermissions(ctx context.Context, roleId int64, perms ...AddPerm) error
	RemovePermissionsFromRole(ctx context.Context, roleId int64, perms ...int64) error

	// Role Attributes management
	GetRoleAttributes(
		ctx context.Context,
		job string,
		grade int32,
	) ([]*permissionsattributes.RoleAttribute, error)
	GetEffectiveRoleAttributes(
		ctx context.Context,
		job string,
		grade int32,
	) ([]*permissionsattributes.RoleAttribute, error)
	UpdateRoleAttributes(
		ctx context.Context,
		job string,
		roleId int64,
		attrs ...*permissionsattributes.RoleAttribute,
	) error
	RemoveAttributesFromRole(
		ctx context.Context,
		roleId int64,
		attrs ...*permissionsattributes.RoleAttribute,
	) error
	RemoveAttributesFromRoleByPermission(
		ctx context.Context,
		roleId int64,
		permissionId int64,
	) error

	// Limits - Job permissions
	GetJobPermissions(ctx context.Context, job string) ([]*permissionspermissions.Permission, error)
	UpdateJobPermissions(
		ctx context.Context,
		job string,
		perms ...*permissionspermissions.PermItem,
	) error
	ApplyJobPermissions(ctx context.Context, job string) error
	ClearJobPermissions(ctx context.Context, job string) error

	// Limits - Job attributes (max values)
	GetJobAttributes(
		ctx context.Context,
		job string,
	) ([]*permissionsattributes.RoleAttribute, error)
	UpdateJobAttributes(
		ctx context.Context,
		job string,
		attrs ...*permissionsattributes.RoleAttribute,
	) error
	ClearJobAttributes(ctx context.Context, job string) error

	// Perms Check
	Can(userInfo *userinfo.UserInfo, perm PermissionRef) bool
	CanRaw(userInfo *userinfo.UserInfo, namespace string, service string, name string) bool
	CanServiceMethod(userInfo *userinfo.UserInfo, serviceMethod string) bool
	CanProto(userInfo *userinfo.UserInfo, perm *permissionspermissions.Permission) bool

	// Attribute retrieval/"check"
	Attr(
		userInfo *userinfo.UserInfo,
		namespace Namespace,
		service Service,
		name Name,
		key Key,
	) (*permissionsattributes.AttributeValues, error)
	AttrStringList(
		userInfo *userinfo.UserInfo,
		attr AttrRef[StringListAttr],
	) (*permissionsattributes.StringList, error)
	AttrJobList(
		userInfo *userinfo.UserInfo,
		attr AttrRef[JobListAttr],
	) (*permissionsattributes.StringList, error)
	AttrJobGradeList(
		userInfo *userinfo.UserInfo,
		attr AttrRef[JobGradeListAttr],
	) (*permissionsattributes.JobGradeList, error)
}

type userCacheKey struct {
	userId int32
	job    string
	grade  int32
	permId int64
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

	permsMap *xsync.Map[int64, *cachePerm]
	// Guard name to permission ID
	permsGuardToIDMap *xsync.Map[string, int64]
	// Job name to map of grade numbers to role ID
	permsJobsRoleMap *xsync.Map[string, *xsync.Map[int32, int64]]
	// Role ID to map of permissions ID and result
	permsRoleMap *xsync.Map[int64, *xsync.Map[int64, bool]]
	// Role ID to Job map
	roleIDToJobMap *xsync.Map[int64, string]

	// Attribute map (key: ID of attribute)
	attrsMap *xsync.Map[int64, *cacheAttr]
	// Role ID to map of role attributes
	attrsRoleMap *xsync.Map[int64, *xsync.Map[int64, *cacheRoleAttr]]
	// Perm ID to map `Key` -> cached attribute (key is name of attribute)
	attrsPermsMap *xsync.Map[int64, *xsync.Map[string, int64]]

	userCanCacheTTL time.Duration
	userCanCache    *cache.LRUCache[userCacheKey, bool]
}

type JobPermission struct {
	PermissionID int64
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

	userCanCache := cache.NewLRUCache[userCacheKey, bool](p.Cfg.Auth.PermsCacheSize)

	logger := p.Logger.WithOptions(
		zap.IncreaseLevel(
			p.Cfg.Log.LevelOverrides.Get(config.LoggingComponentPerms, p.Cfg.LogLevel),
		),
	)

	ps := &Perms{
		logger: logger,
		db:     p.DB,
		wg:     sync.WaitGroup{},

		tracer: p.TP.Tracer("perms"),

		nc: p.NC,

		cleanupRolesForMissingJobs: p.Cfg.Game.CleanupRolesForMissingJobs,
		startJobGrade:              p.Cfg.Game.StartJobGrade,

		permsMap:          xsync.NewMap[int64, *cachePerm](),
		permsGuardToIDMap: xsync.NewMap[string, int64](),
		permsJobsRoleMap:  xsync.NewMap[string, *xsync.Map[int32, int64]](),
		permsRoleMap:      xsync.NewMap[int64, *xsync.Map[int64, bool]](),
		roleIDToJobMap:    xsync.NewMap[int64, string](),

		attrsMap:      xsync.NewMap[int64, *cacheAttr](),
		attrsRoleMap:  xsync.NewMap[int64, *xsync.Map[int64, *cacheRoleAttr]](),
		attrsPermsMap: xsync.NewMap[int64, *xsync.Map[string, int64]](),

		userCanCacheTTL: p.Cfg.Auth.PermsCacheTTL,
		userCanCache:    userCanCache,
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		go userCanCache.StartJanitor(ctxCancel, p.Cfg.Auth.PermsCacheTTL+1*time.Second)

		return ps.init(ctxCancel, ctxStartup, p)
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		ps.wg.Wait()

		if ps.ncSub != nil {
			if err := ps.ncSub.Unsubscribe(); err != nil {
				ps.logger.Error("failed to unsubscribe from previous perms subject", zap.Error(err))
			}
			ps.ncSub = nil
		}

		return nil
	}))

	return ps, nil
}

type cachePerm struct {
	ID        int64
	Namespace Namespace
	Service   Service
	Name      Name
	GuardName string
	Order     *int32
	Icon      *string
}

type cacheAttr struct {
	ID           int64
	PermissionID int64
	Namespace    Namespace
	Service      Service
	Name         Name
	Key          Key
	Type         permissionsattributes.AttributeTypes
	ValidValues  *permissionsattributes.AttributeValues
}

type cacheRoleAttr struct {
	Job          string
	AttrID       int64
	PermissionID int64
	Key          Key
	Type         permissionsattributes.AttributeTypes
	Value        *permissionsattributes.AttributeValues
}

func (ps *Perms) init(ctxCancel context.Context, ctxStartup context.Context, params Params) error {
	cfgDefaultPerms := params.AppConfig.Get().GetPerms().GetDefault()
	defaultPerms := make([]string, len(cfgDefaultPerms))
	for i := range cfgDefaultPerms {
		split := strings.Split(cfgDefaultPerms[i].GetCategory(), ".")
		namespace := strings.Join(split[:len(split)-1], ".")
		svc := split[len(split)-1]

		defaultPerms[i] = BuildGuard(
			Namespace(namespace),
			Service(svc),
			Name(cfgDefaultPerms[i].GetName()),
		)
	}

	if err := ps.loadData(ctxStartup); err != nil {
		return err
	}
	ps.logger.Debug("permissions loaded")

	if err := ps.registerSubscriptions(ctxCancel); err != nil {
		return fmt.Errorf("failed to register events subscriptions. %w", err)
	}
	ps.logger.Debug("registered events subscription")

	if err := ps.register(ctxStartup, defaultPerms); err != nil {
		return fmt.Errorf("failed to register permissions. %w", err)
	}

	ps.wg.Go(func() {
		if err := ps.ApplyJobPermissions(ctxCancel, ""); err != nil {
			ps.logger.Error("failed to apply job permissions", zap.Error(err))
			return
		}
		ps.logger.Debug("successfully applied job permissions")
	})

	return nil
}
