package postals

import (
	"github.com/fivenet-app/fivenet/v2026/pkg/coords"
)

func NewForTests() (Postals, error) {
	return coords.NewReadOnly([]*Postal{})
}
