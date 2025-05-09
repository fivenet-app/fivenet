package postals

import (
	"github.com/fivenet-app/fivenet/v2025/pkg/coords"
)

func NewForTests() (Postals, error) {
	return coords.NewReadOnly([]*Postal{})
}
