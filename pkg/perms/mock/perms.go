package mock

import (
	"strings"

	"github.com/galexrt/arpanet/pkg/perms"
	"github.com/galexrt/arpanet/pkg/perms/collections"
	"github.com/galexrt/arpanet/pkg/perms/helpers"
	"github.com/galexrt/arpanet/query/arpanet/model"
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

func (p *PermsMock) GetAllPermissions() (collections.Permissions, error) {
	ps := collections.Permissions{}

	track := map[string]interface{}{}
	i := 0
	// Add all permissions once collected by iterating over `p.UserPerms`
	for _, v := range p.UserPerms {
		for k := range v {
			if _, ok := track[k]; !ok {
				ps = append(ps, &model.ArpanetPermissions{
					ID:        uint64(i),
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

func (p *PermsMock) CreatePermission(name string, description string) error {
	return nil
}

func (p *PermsMock) GetAllPermissionsOfUser(userId int32) (collections.Permissions, error) {
	ps := collections.Permissions{}

	if _, ok := p.UserPerms[userId]; !ok {
		return nil, nil
	}

	track := map[string]interface{}{}
	i := 0
	for k := range p.UserPerms[userId] {
		if _, ok := track[k]; !ok {
			ps = append(ps, &model.ArpanetPermissions{
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

func (p *PermsMock) GetAllPermissionsByPrefixOfUser(userId int32, prefix string) (collections.Permissions, error) {
	ps, nil := p.GetAllPermissionsOfUser(userId)

	name := helpers.Guard(prefix) + "-"

	return ps.HasPrefix(name), nil
}

func (p *PermsMock) GetSuffixOfPermissionsByPrefixOfUser(userId int32, prefix string) ([]string, error) {
	ps, _ := p.GetAllPermissionsByPrefixOfUser(userId, prefix)

	prefix = helpers.Guard(prefix) + "-"

	list := make([]string, len(ps))
	for i := 0; i < len(ps); i++ {
		list[i] = strings.TrimPrefix(ps[i].GuardName, prefix)
	}

	return list, nil
}

func (p *PermsMock) GetRoles(prefix string) (collections.Roles, error) {
	r := collections.Roles{}

	track := map[string]interface{}{}
	i := 0
	for _, v := range p.UserRoles {
		for k := range v {
			if _, ok := track[k]; !ok {
				r = append(r, &model.ArpanetRoles{
					ID:        uint64(i),
					Name:      k,
					GuardName: k,
				})
				i++
				track[k] = nil
			}
		}
	}

	return r, nil
}

func (p *PermsMock) GetUserRoles(userId int32) (collections.Roles, error) {
	r := collections.Roles{}

	uRoles, ok := p.UserRoles[userId]
	if !ok {
		return r, nil
	}

	i := 0
	for k := range uRoles {
		r = append(r, &model.ArpanetRoles{
			ID:        uint64(i),
			Name:      k,
			GuardName: k,
		})
		i++
	}

	return nil, nil
}

func (p *PermsMock) AddUserRoles(userId int32, roles ...string) error {
	if _, ok := p.UserRoles[userId]; !ok {
		p.UserRoles[userId] = map[string]interface{}{}
	}

	for _, role := range roles {
		p.UserRoles[userId][role] = role
	}

	return nil
}

func (p *PermsMock) RemoveUserRoles(userId int32, roles ...string) error {
	if _, ok := p.UserRoles[userId]; !ok {
		return nil
	}

	for _, role := range roles {
		delete(p.UserRoles[userId], role)
	}

	return nil
}

func (p *PermsMock) Can(userId int32, perm ...string) bool {
	guard := helpers.Guard(strings.Join(perm, "."))
	p.Counter.IncUser(userId, guard)

	if _, ok := p.UserPerms[userId]; !ok {
		return false
	}

	_, ok := p.UserPerms[userId][guard]
	return ok
}
