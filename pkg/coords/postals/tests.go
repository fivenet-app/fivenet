package postals

import (
	"github.com/fivenet-app/fivenet/pkg/coords"
)

func NewForTests() (Postals, error) {
	return coords.NewReadOnly([]*Postal{})
}
