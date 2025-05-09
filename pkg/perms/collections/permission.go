package collections

import (
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/model"
)

// Permissions provides methods for you to manage array data more easily.
type Permissions []*model.FivenetPermissions

// Len returns the number of elements of the array.
// @return int64
func (u Permissions) Len() (length int64) {
	return int64(len(u))
}

// IDs returns an array of the permission array's ids.
// @return []uint
func (u Permissions) IDs() (IDs []uint64) {
	for _, permission := range u {
		IDs = append(IDs, permission.ID)
	}
	return IDs
}

// Names returns an array of the permission array's names.
// @return []string
func (u Permissions) Names() (names []string) {
	for _, permission := range u {
		names = append(names, permission.Name)
	}
	return names
}

// GuardNames returns an array of the permission array's guard names.
// @return []string
func (u Permissions) GuardNames() (guards []string) {
	for _, permission := range u {
		guards = append(guards, permission.GuardName)
	}
	return guards
}
