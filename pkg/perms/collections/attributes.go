package collections

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/permissions"
)

// Attributes provides methods for you to manage array data more easily.
type Attributes []*permissions.RoleAttribute

// Len returns the number of elements of the array.
//
//	@return	int64
func (u Attributes) Len() int64 {
	return int64(len(u))
}

// IDs returns an array of the attribute array's ids.
//
//	@return	[]uint64
func (u Attributes) IDs() []uint64 {
	ids := make([]uint64, 0, len(u))
	for _, attribute := range u {
		ids = append(ids, attribute.GetAttrId())
	}
	return ids
}

// Names returns an array of the attribute array's key names.
//
//	@return	[]string
func (u Attributes) Names() []string {
	names := make([]string, 0, len(u))
	for _, attribute := range u {
		names = append(names, attribute.GetKey())
	}
	return names
}
