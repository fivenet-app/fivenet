package collections

import (
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/model"
)

// Attributes provides methods for you to manage array data more easily.
type Attributes []*model.FivenetAttrs

// Len returns the number of elements of the array.
// @return int64
func (u Attributes) Len() (length int64) {
	return int64(len(u))
}

// IDs returns an array of the attribute array's ids.
// @return []uint
func (u Attributes) IDs() (IDs []uint64) {
	for _, attribute := range u {
		IDs = append(IDs, attribute.ID)
	}
	return IDs
}

// Names returns an array of the attribute array's key names.
// @return []string
func (u Attributes) Names() (names []string) {
	for _, attribute := range u {
		names = append(names, attribute.Key)
	}
	return names
}
