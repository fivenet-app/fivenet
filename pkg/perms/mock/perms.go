package mock

import (
	"context"

	"github.com/galexrt/fivenet/gen/go/proto/resources/permissions"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/pkg/perms/collections"
	"github.com/galexrt/fivenet/pkg/perms/helpers"
	"github.com/galexrt/fivenet/query/fivenet/model"
)

type PermCounter struct {
	UserPerms map[int32]map[string]int
}

func NewPermCounter() *PermCounter {
	return &PermCounter{
		UserPerms: map[int32]map[string]int{},
	}
}

func (p *PermCounter) GetUser(userId int32) map[string]int {
	if _, ok := p.UserPerms[userId]; !ok {
		p.UserPerms[userId] = map[string]int{}
	}

	return p.UserPerms[userId]
}

func (p *PermCounter) IncUser(userId int32, name string) {
	if _, ok := p.UserPerms[userId]; !ok {
		p.UserPerms[userId] = map[string]int{}
	}

	if _, ok := p.UserPerms[userId][name]; !ok {
		p.UserPerms[userId][name] = 1
	} else {
		p.UserPerms[userId][name]++
	}
}

type PermsMock struct {
	perms.Permissions

	Counter *PermCounter

	UserPerms map[int32]map[string]interface{}
	UserRoles map[int32]map[string]interface{}
}

func NewMock() *PermsMock {
	return &PermsMock{
		Counter:   NewPermCounter(),
		UserPerms: map[int32]map[string]interface{}{},
		UserRoles: map[int32]map[string]interface{}{},
	}
}

func (p *PermsMock) AddUserPerm(userId int32, perm string) {
	if _, ok := p.UserPerms[userId]; !ok {
		p.UserPerms[userId] = map[string]interface{}{}
	}

	p.UserPerms[userId][perm] = nil
}

func (p *PermsMock) RemoveUserPerm(userId int32, perm string) {
	if _, ok := p.UserPerms[userId]; !ok {
		return
	}

	delete(p.UserPerms[userId], perm)
}

// Implementation of perms.Permissions

func (p *PermsMock) GetAllPermissions(ctx context.Context) ([]*permissions.Permission, error) {
	ps := []*permissions.Permission{}

	track := map[string]interface{}{}
	i := 0
	// Add all permissions once collected by iterating over `p.UserPerms`
	for _, v := range p.UserPerms {
		for k := range v {
			if _, ok := track[k]; !ok {
				ps = append(ps, &permissions.Permission{
					Id:        uint64(i),
					Name:      k,
					GuardName: k,
				})
				i++
				track[k] = nil
			}
		}
	}

	return ps, nil
}

func (p *PermsMock) GetPermissionsByIDs(ctx context.Context, ids ...uint64) ([]*permissions.Permission, error) {
	return nil, nil
}

func (p *PermsMock) CreatePermission(ctx context.Context, category perms.Category, name perms.Name) (uint64, error) {
	return 0, nil
}

func (p *PermsMock) GetPermissionsOfUser(userInfo *userinfo.UserInfo) (collections.Permissions, error) {
	ps := collections.Permissions{}

	if _, ok := p.UserPerms[userInfo.UserId]; !ok {
		return nil, nil
	}

	track := map[string]interface{}{}
	i := 0
	for k := range p.UserPerms[userInfo.UserId] {
		if _, ok := track[k]; !ok {
			ps = append(ps, &model.FivenetPermissions{
				ID:        uint64(i),
				Name:      k,
				GuardName: k,
			})
			i++
			track[k] = nil
		}
	}

	return ps, nil
}

func (p *PermsMock) GetRoles(ctx context.Context, prefix string) (collections.Roles, error) {
	r := collections.Roles{}

	track := map[string]interface{}{}
	i := 0
	for _, v := range p.UserRoles {
		for k := range v {
			if _, ok := track[k]; !ok {
				r = append(r, &model.FivenetRoles{
					ID:    uint64(i),
					Job:   "ambulance",
					Grade: 0,
				})
				i++
				track[k] = nil
			}
		}
	}

	return r, nil
}

func (p *PermsMock) CountRolesForJob(ctx context.Context, prefix string) (int64, error) {
	return 0, nil
}

func (p *PermsMock) GetRole(ctx context.Context, id uint64) (*model.FivenetRoles, error) {
	return nil, nil
}

func (p *PermsMock) GetRoleByJobAndGrade(ctx context.Context, job string, grade int32) (*model.FivenetRoles, error) {
	return nil, nil
}

func (p *PermsMock) GetRolePermissions(ctx context.Context, id uint64) ([]*permissions.Permission, error) {
	return nil, nil
}

func (p *PermsMock) CreateRole(ctx context.Context, job string, grade int32) (*model.FivenetRoles, error) {
	return nil, nil
}

func (p *PermsMock) DeleteRole(ctx context.Context, id uint64) error {
	return nil
}

func (p *PermsMock) UpdateRolePermissions(ctx context.Context, id uint64, perms ...perms.AddPerm) error {
	return nil
}

func (p *PermsMock) RemovePermissionsFromRole(ctx context.Context, id uint64, perms ...uint64) error {
	return nil
}

func (p *PermsMock) GetUserRoles(userId int32) (collections.Roles, error) {
	r := collections.Roles{}

	uRoles, ok := p.UserRoles[userId]
	if !ok {
		return r, nil
	}

	i := 0
	for range uRoles {
		r = append(r, &model.FivenetRoles{
			ID:    uint64(i),
			Job:   "ambulance",
			Grade: 0,
		})
		i++
	}

	return nil, nil
}

func (p *PermsMock) AddUserToRoles(ctx context.Context, userId int32, roles ...string) error {
	if _, ok := p.UserRoles[userId]; !ok {
		p.UserRoles[userId] = map[string]interface{}{}
	}

	for _, role := range roles {
		p.UserRoles[userId][role] = role
	}

	return nil
}

func (p *PermsMock) RemoveUserFromRoles(ctx context.Context, userId int32, roles ...string) error {
	if _, ok := p.UserRoles[userId]; !ok {
		return nil
	}

	for _, role := range roles {
		delete(p.UserRoles[userId], role)
	}

	return nil
}

func (p *PermsMock) Can(userInfo *userinfo.UserInfo, category perms.Category, name perms.Name) bool {
	guard := helpers.Guard(perms.BuildGuard(category, name))
	p.Counter.IncUser(userInfo.UserId, guard)

	if _, ok := p.UserPerms[userInfo.UserId]; !ok {
		return false
	}

	_, ok := p.UserPerms[userInfo.UserId][guard]
	return ok
}
