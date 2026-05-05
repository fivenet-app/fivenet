package perms

import (
	"testing"
	"time"

	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/cache"
	"github.com/puzpuzpuz/xsync/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

const (
	testNamespace = "test"
	testService   = "test.Service"
	testNameView  = "View"
	testNameEdit  = "Edit"
	testKeyFields = "Fields"
)

func TestPermissionRefFromServiceMethodSplitsNamespaceAndService(t *testing.T) {
	t.Parallel()

	ref, ok := PermissionRefFromServiceMethod("calendar.CalendarService/CreateCalendar")
	require.True(t, ok)
	assert.Equal(t, Namespace("calendar"), ref.Namespace())
	assert.Equal(t, Service("CalendarService"), ref.Service())
	assert.Equal(t, Name("CreateCalendar"), ref.Name())

	_, ok = PermissionRefFromServiceMethod("CalendarService/CreateCalendar")
	assert.False(t, ok)
	_, ok = PermissionRefFromServiceMethod("calendar./CreateCalendar")
	assert.False(t, ok)
}

func TestCanServiceMethodUsesSplitService(t *testing.T) {
	t.Parallel()

	ps := newTestPerms()
	addTestPermission(ps, 1, "calendar", "CalendarService", "CreateCalendar")
	addTestRole(ps, 100, DefaultRoleJob, 0, map[int64]bool{1: true})

	user := &userinfo.UserInfo{
		UserId:   10,
		Job:      "unemployed",
		JobGrade: 0,
	}

	assert.True(t, ps.CanServiceMethod(user, "calendar.CalendarService/CreateCalendar"))
	assert.False(t, ps.CanServiceMethod(user, "calendar.CalendarService/DeleteCalendar"))
}

func TestCanResolvesDefaultFallbackAndGradeOverrides(t *testing.T) {
	t.Parallel()

	ps := newTestPerms()
	addTestPermission(ps, 1, testNamespace, testService, testNameView)
	addTestPermission(ps, 2, testNamespace, testService, testNameEdit)
	addTestPermission(ps, 3, testNamespace, testService, "Admin")
	addTestRole(ps, 100, DefaultRoleJob, 0, map[int64]bool{1: true})
	addTestRole(ps, 200, "police", 0, map[int64]bool{2: true})
	addTestRole(ps, 201, "police", 1, map[int64]bool{1: false})
	addTestRole(ps, 202, "police", 2, map[int64]bool{1: true})

	user := &userinfo.UserInfo{
		UserId:   10,
		Job:      "police",
		JobGrade: 0,
	}

	assert.True(
		t,
		ps.CanRaw(user, testNamespace, testService, testNameView),
		"default role grants missing job permissions",
	)
	assert.True(
		t,
		ps.CanRaw(user, testNamespace, testService, testNameEdit),
		"lower job grade grants are inherited",
	)

	user.JobGrade = 1
	assert.False(
		t,
		ps.CanRaw(user, testNamespace, testService, testNameView),
		"grade deny overrides default grant",
	)
	assert.True(
		t,
		ps.CanRaw(user, testNamespace, testService, testNameEdit),
		"lower grade grant remains effective",
	)

	user.JobGrade = 2
	assert.True(
		t,
		ps.CanRaw(user, testNamespace, testService, testNameView),
		"higher grade allow overrides lower deny",
	)
	assert.False(
		t,
		ps.CanRaw(user, testNamespace, testService, "Missing"),
		"unknown permission is denied",
	)

	user.Superuser = true
	assert.True(
		t,
		ps.CanRaw(user, testNamespace, testService, "Admin"),
		"superusers bypass role checks for known permissions",
	)
}

func TestCanCacheKeyIncludesJobAndGrade(t *testing.T) {
	t.Parallel()

	ps := newTestPerms()
	addTestPermission(ps, 1, testNamespace, testService, testNameView)
	addTestRole(ps, 100, DefaultRoleJob, 0, nil)
	addTestRole(ps, 200, "police", 0, map[int64]bool{1: false})
	addTestRole(ps, 201, "police", 1, map[int64]bool{1: true})
	addTestRole(ps, 300, "ambulance", 0, map[int64]bool{1: false})

	user := &userinfo.UserInfo{
		UserId:   10,
		Job:      "police",
		JobGrade: 0,
	}

	require.False(t, ps.CanRaw(user, testNamespace, testService, testNameView))

	user.JobGrade = 1
	assert.True(
		t,
		ps.CanRaw(user, testNamespace, testService, testNameView),
		"same user at a different grade must not reuse stale cache",
	)

	user.Job = "ambulance"
	user.JobGrade = 0
	assert.False(
		t,
		ps.CanRaw(user, testNamespace, testService, testNameView),
		"same user at a different job must not reuse stale cache",
	)
}

func TestClearUserCanCacheAllowsRolePermissionChangesToTakeEffect(t *testing.T) {
	t.Parallel()

	ps := newTestPerms()
	addTestPermission(ps, 1, testNamespace, testService, testNameView)
	addTestRole(ps, 100, DefaultRoleJob, 0, nil)
	addTestRole(ps, 200, "police", 0, map[int64]bool{1: false})

	user := &userinfo.UserInfo{
		UserId:   10,
		Job:      "police",
		JobGrade: 0,
	}

	require.False(t, ps.CanRaw(user, testNamespace, testService, testNameView))

	rolePerms, ok := ps.permsRoleMap.Load(200)
	require.True(t, ok)
	rolePerms.Store(1, true)
	ps.clearUserCanCache()

	assert.True(
		t,
		ps.CanRaw(user, testNamespace, testService, testNameView),
		"permission updates must clear cached decisions",
	)
}

func TestGetPermissionsOfUserFiltersDeniedPermissions(t *testing.T) {
	t.Parallel()

	ps := newTestPerms()
	addTestPermission(ps, 1, testNamespace, testService, testNameView)
	addTestPermission(ps, 2, testNamespace, testService, testNameEdit)
	addTestRole(ps, 100, DefaultRoleJob, 0, map[int64]bool{
		1: true,
		2: true,
	})
	addTestRole(ps, 200, "police", 0, map[int64]bool{
		1: false,
	})

	perms, err := ps.GetPermissionsOfUser(&userinfo.UserInfo{
		UserId:   10,
		Job:      "police",
		JobGrade: 0,
	})
	require.NoError(t, err)

	assert.ElementsMatch(
		t,
		[]string{BuildGuard(testNamespace, testService, testNameEdit)},
		[]string{perms[0].GetGuardName()},
	)
}

func TestAttrUsesClosestRoleGradeAndReturnsClone(t *testing.T) {
	t.Parallel()

	ps := newTestPerms()
	addTestPermission(ps, 1, testNamespace, testService, testNameView)
	addTestRole(ps, 100, DefaultRoleJob, 0, nil)
	addTestRole(ps, 200, "police", 0, nil)
	addTestRole(ps, 201, "police", 2, nil)
	addTestAttribute(
		ps,
		10,
		1,
		testNamespace,
		testNameView,
		testKeyFields,
		stringListValues("name", "phone", "address"),
	)
	addTestRoleAttribute(ps, 200, 1, 10, testKeyFields, stringListValues("name"))
	addTestRoleAttribute(ps, 201, 1, 10, testKeyFields, stringListValues("phone"))

	user := &userinfo.UserInfo{
		UserId:   10,
		Job:      "police",
		JobGrade: 1,
	}
	fieldsRef := NewStringListAttrRef(
		NewPermissionRef(testNamespace, testService, testNameView),
		testKeyFields,
	)

	fields, err := ps.AttrStringList(user, fieldsRef)
	require.NoError(t, err)
	assert.Equal(t, []string{"name"}, fields.GetStrings())

	user.JobGrade = 2
	fields, err = ps.AttrStringList(user, fieldsRef)
	require.NoError(t, err)
	assert.Equal(t, []string{"phone"}, fields.GetStrings())

	fields.Strings = append(fields.GetStrings(), "mutated")
	fields, err = ps.AttrStringList(user, fieldsRef)
	require.NoError(t, err)
	assert.Equal(
		t,
		[]string{"phone"},
		fields.GetStrings(),
		"attribute values returned to callers must be cloned",
	)

	superuserFields, err := ps.AttrStringList(&userinfo.UserInfo{
		UserId:    11,
		Job:       "unemployed",
		JobGrade:  0,
		Superuser: true,
	}, fieldsRef)
	require.NoError(t, err)
	assert.Equal(t, []string{"name", "phone", "address"}, superuserFields.GetStrings())
}

func newTestPerms() *Perms {
	return &Perms{
		logger:        zap.NewNop(),
		startJobGrade: 0,

		permsMap:          xsync.NewMap[int64, *cachePerm](),
		permsGuardToIDMap: xsync.NewMap[string, int64](),
		permsJobsRoleMap:  xsync.NewMap[string, *xsync.Map[int32, int64]](),
		permsRoleMap:      xsync.NewMap[int64, *xsync.Map[int64, bool]](),
		roleIDToJobMap:    xsync.NewMap[int64, string](),

		attrsMap:      xsync.NewMap[int64, *cacheAttr](),
		attrsRoleMap:  xsync.NewMap[int64, *xsync.Map[int64, *cacheRoleAttr]](),
		attrsPermsMap: xsync.NewMap[int64, *xsync.Map[string, int64]](),

		userCanCacheTTL: time.Hour,
		userCanCache:    cache.NewLRUCache[userCacheKey, bool](64),
	}
}

func addTestPermission(p *Perms, id int64, namespace Namespace, service Service, name Name) {
	p.permsMap.Store(id, &cachePerm{
		ID:        id,
		Namespace: namespace,
		Service:   service,
		Name:      name,
		GuardName: BuildGuard(namespace, service, name),
	})
	p.permsGuardToIDMap.Store(BuildGuard(namespace, service, name), id)
}

func addTestRole(p *Perms, id int64, job string, grade int32, rolePerms map[int64]bool) {
	grades, _ := p.permsJobsRoleMap.LoadOrCompute(job, func() (*xsync.Map[int32, int64], bool) {
		return xsync.NewMap[int32, int64](), false
	})
	grades.Store(grade, id)
	p.roleIDToJobMap.Store(id, job)

	permsMap := xsync.NewMap[int64, bool]()
	for permId, val := range rolePerms {
		permsMap.Store(permId, val)
	}
	p.permsRoleMap.Store(id, permsMap)
}

func addTestAttribute(
	p *Perms,
	attrId int64,
	permId int64,
	category Namespace,
	name Name,
	key Key,
	validValues *permissionsattributes.AttributeValues,
) {
	p.attrsMap.Store(attrId, &cacheAttr{
		ID:           attrId,
		PermissionID: permId,
		Namespace:    category,
		Name:         name,
		Key:          key,
		Type:         permissionsattributes.StringListAttributeType,
		ValidValues:  validValues,
	})

	attrs, _ := p.attrsPermsMap.LoadOrCompute(permId, func() (*xsync.Map[string, int64], bool) {
		return xsync.NewMap[string, int64](), false
	})
	attrs.Store(string(key), attrId)
}

func addTestRoleAttribute(
	p *Perms,
	roleId int64,
	permId int64,
	attrId int64,
	key Key,
	value *permissionsattributes.AttributeValues,
) {
	roleAttrs, _ := p.attrsRoleMap.LoadOrCompute(
		roleId,
		func() (*xsync.Map[int64, *cacheRoleAttr], bool) {
			return xsync.NewMap[int64, *cacheRoleAttr](), false
		},
	)
	roleAttrs.Store(attrId, &cacheRoleAttr{
		Job:          "police",
		AttrID:       attrId,
		PermissionID: permId,
		Key:          key,
		Type:         permissionsattributes.StringListAttributeType,
		Value:        value,
	})
}

func stringListValues(values ...string) *permissionsattributes.AttributeValues {
	return &permissionsattributes.AttributeValues{
		ValidValues: &permissionsattributes.AttributeValues_StringList{
			StringList: &permissionsattributes.StringList{
				Strings: values,
			},
		},
	}
}
