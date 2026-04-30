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
	testCategory Category = "test.Service"
	testView     Name     = "View"
	testEdit     Name     = "Edit"
	testFields   Key      = "Fields"
)

func TestCanResolvesDefaultFallbackAndGradeOverrides(t *testing.T) {
	t.Parallel()

	ps := newTestPerms()
	addTestPermission(ps, 1, testCategory, testView)
	addTestPermission(ps, 2, testCategory, testEdit)
	addTestPermission(ps, 3, testCategory, "Admin")
	addTestRole(ps, 100, DefaultRoleJob, 0, map[int64]bool{1: true})
	addTestRole(ps, 200, "police", 0, map[int64]bool{2: true})
	addTestRole(ps, 201, "police", 1, map[int64]bool{1: false})
	addTestRole(ps, 202, "police", 2, map[int64]bool{1: true})

	user := &userinfo.UserInfo{
		UserId:   10,
		Job:      "police",
		JobGrade: 0,
	}

	assert.True(t, ps.Can(user, testCategory, testView), "default role grants missing job permissions")
	assert.True(t, ps.Can(user, testCategory, testEdit), "lower job grade grants are inherited")

	user.JobGrade = 1
	assert.False(t, ps.Can(user, testCategory, testView), "grade deny overrides default grant")
	assert.True(t, ps.Can(user, testCategory, testEdit), "lower grade grant remains effective")

	user.JobGrade = 2
	assert.True(t, ps.Can(user, testCategory, testView), "higher grade allow overrides lower deny")
	assert.False(t, ps.Can(user, testCategory, "Missing"), "unknown permission is denied")

	user.Superuser = true
	assert.True(t, ps.Can(user, testCategory, "Admin"), "superusers bypass role checks for known permissions")
}

func TestCanCacheKeyIncludesJobAndGrade(t *testing.T) {
	t.Parallel()

	ps := newTestPerms()
	addTestPermission(ps, 1, testCategory, testView)
	addTestRole(ps, 100, DefaultRoleJob, 0, nil)
	addTestRole(ps, 200, "police", 0, map[int64]bool{1: false})
	addTestRole(ps, 201, "police", 1, map[int64]bool{1: true})
	addTestRole(ps, 300, "ambulance", 0, map[int64]bool{1: false})

	user := &userinfo.UserInfo{
		UserId:   10,
		Job:      "police",
		JobGrade: 0,
	}

	require.False(t, ps.Can(user, testCategory, testView))

	user.JobGrade = 1
	assert.True(t, ps.Can(user, testCategory, testView), "same user at a different grade must not reuse stale cache")

	user.Job = "ambulance"
	user.JobGrade = 0
	assert.False(t, ps.Can(user, testCategory, testView), "same user at a different job must not reuse stale cache")
}

func TestClearUserCanCacheAllowsRolePermissionChangesToTakeEffect(t *testing.T) {
	t.Parallel()

	ps := newTestPerms()
	addTestPermission(ps, 1, testCategory, testView)
	addTestRole(ps, 100, DefaultRoleJob, 0, nil)
	addTestRole(ps, 200, "police", 0, map[int64]bool{1: false})

	user := &userinfo.UserInfo{
		UserId:   10,
		Job:      "police",
		JobGrade: 0,
	}

	require.False(t, ps.Can(user, testCategory, testView))

	rolePerms, ok := ps.permsRoleMap.Load(200)
	require.True(t, ok)
	rolePerms.Store(1, true)
	ps.clearUserCanCache()

	assert.True(t, ps.Can(user, testCategory, testView), "permission updates must clear cached decisions")
}

func TestGetPermissionsOfUserFiltersDeniedPermissions(t *testing.T) {
	t.Parallel()

	ps := newTestPerms()
	addTestPermission(ps, 1, testCategory, testView)
	addTestPermission(ps, 2, testCategory, testEdit)
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

	assert.ElementsMatch(t, []string{BuildGuard(testCategory, testEdit)}, perms.GuardNames())
}

func TestAttrUsesClosestRoleGradeAndReturnsClone(t *testing.T) {
	t.Parallel()

	ps := newTestPerms()
	addTestPermission(ps, 1, testCategory, testView)
	addTestRole(ps, 100, DefaultRoleJob, 0, nil)
	addTestRole(ps, 200, "police", 0, nil)
	addTestRole(ps, 201, "police", 2, nil)
	addTestAttribute(ps, 10, 1, testCategory, testView, testFields, stringListValues("name", "phone", "address"))
	addTestRoleAttribute(ps, 200, 1, 10, testFields, stringListValues("name"))
	addTestRoleAttribute(ps, 201, 1, 10, testFields, stringListValues("phone"))

	user := &userinfo.UserInfo{
		UserId:   10,
		Job:      "police",
		JobGrade: 1,
	}

	fields, err := ps.AttrStringList(user, testCategory, testView, testFields)
	require.NoError(t, err)
	assert.Equal(t, []string{"name"}, fields.GetStrings())

	user.JobGrade = 2
	fields, err = ps.AttrStringList(user, testCategory, testView, testFields)
	require.NoError(t, err)
	assert.Equal(t, []string{"phone"}, fields.GetStrings())

	fields.Strings = append(fields.GetStrings(), "mutated")
	fields, err = ps.AttrStringList(user, testCategory, testView, testFields)
	require.NoError(t, err)
	assert.Equal(t, []string{"phone"}, fields.GetStrings(), "attribute values returned to callers must be cloned")

	superuserFields, err := ps.AttrStringList(&userinfo.UserInfo{
		UserId:    11,
		Job:       "unemployed",
		JobGrade:  0,
		Superuser: true,
	}, testCategory, testView, testFields)
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

func addTestPermission(p *Perms, id int64, category Category, name Name) {
	p.permsMap.Store(id, &cachePerm{
		ID:        id,
		Category:  category,
		Name:      name,
		GuardName: BuildGuard(category, name),
	})
	p.permsGuardToIDMap.Store(BuildGuard(category, name), id)
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
	category Category,
	name Name,
	key Key,
	validValues *permissionsattributes.AttributeValues,
) {
	p.attrsMap.Store(attrId, &cacheAttr{
		ID:           attrId,
		PermissionID: permId,
		Category:     category,
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
	roleAttrs, _ := p.attrsRoleMap.LoadOrCompute(roleId, func() (*xsync.Map[int64, *cacheRoleAttr], bool) {
		return xsync.NewMap[int64, *cacheRoleAttr](), false
	})
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
