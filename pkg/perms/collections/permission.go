package collections

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/permissions"
)

// Permissions provides methods for you to manage array data more easily.
type Permissions []*permissions.Permission

// Len returns the number of elements of the array.
//
//	@return	int64
func (u Permissions) Len() int64 {
	return int64(len(u))
}

// IDs returns an array of the permission array's ids.
//
//	@return	[]uint
func (u Permissions) IDs() []int64 {
	ids := make([]int64, 0, len(u))
	for _, permission := range u {
		ids = append(ids, permission.GetId())
	}
	return ids
}

// Names returns an array of the permission array's names.
//
//	@return	[]string
func (u Permissions) Names() []string {
	names := make([]string, 0, len(u))
	for _, permission := range u {
		names = append(names, permission.GetName())
	}
	return names
}

// GuardNames returns an array of the permission array's guard names.
//
//	@return	[]string
func (u Permissions) GuardNames() []string {
	guards := make([]string, 0, len(u))
	for _, permission := range u {
		guards = append(guards, permission.GetGuardName())
	}
	return guards
}
