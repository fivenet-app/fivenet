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

func (p *PermCounter) GetUser(userID int32) map[string]int {
	if _, ok := p.UserPerms[userID]; !ok {
		p.UserPerms[userID] = map[string]int{}
	}

	return p.UserPerms[userID]
}

func (p *PermCounter) IncUser(userID int32, name string) {
	if _, ok := p.UserPerms[userID]; !ok {
		p.UserPerms[userID] = map[string]int{}
	}

	if _, ok := p.UserPerms[userID][name]; !ok {
		p.UserPerms[userID][name] = 1
	} else {
		p.UserPerms[userID][name]++
	}
}

type PermsMock struct {
	perms.Permissions

	Counter *PermCounter

	UserPerms map[int32]map[string]interface{}
}

func NewMock() *PermsMock {
	return &PermsMock{
		Counter:   NewPermCounter(),
		UserPerms: map[int32]map[string]interface{}{},
	}
}

func (p *PermsMock) AddUserPerm(userID int32, perm string) {
	if _, ok := p.UserPerms[userID]; !ok {
		p.UserPerms[userID] = map[string]interface{}{}
	}

	p.UserPerms[userID][perm] = nil
}

func (p *PermsMock) RemoveUserPerm(userID int32, perm string) {
	if _, ok := p.UserPerms[userID]; !ok {
		return
	}

	delete(p.UserPerms[userID], perm)
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

func (p *PermsMock) GetAllPermissionsOfUser(userID int32) (collections.Permissions, error) {
	ps := collections.Permissions{}

	if _, ok := p.UserPerms[userID]; !ok {
		return nil, nil
	}

	track := map[string]interface{}{}
	i := 0
	for k := range p.UserPerms[userID] {
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

func (p *PermsMock) GetAllPermissionsByPrefixOfUser(userID int32, prefix string) (collections.Permissions, error) {
	ps, nil := p.GetAllPermissionsOfUser(userID)

	name := helpers.Guard(prefix) + "-"

	return ps.HasPrefix(name), nil
}

func (p *PermsMock) GetSuffixOfPermissionsByPrefixOfUser(userID int32, prefix string) ([]string, error) {
	ps, _ := p.GetAllPermissionsByPrefixOfUser(userID, prefix)

	prefix = helpers.Guard(prefix) + "-"

	list := make([]string, len(ps))
	for i := 0; i < len(ps); i++ {
		list[i] = strings.TrimPrefix(ps[i].GuardName, prefix)
	}

	return list, nil
}

func (p *PermsMock) Can(userID int32, perm ...string) bool {
	guard := helpers.Guard(strings.Join(perm, "."))
	p.Counter.IncUser(userID, guard)

	if _, ok := p.UserPerms[userID]; !ok {
		return false
	}

	_, ok := p.UserPerms[userID][guard]
	return ok
}
