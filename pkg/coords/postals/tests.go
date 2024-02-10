package postals

import (
	"github.com/galexrt/fivenet/pkg/coords"
)

func NewForTests() (Postals, error) {
	return coords.NewReadOnly[*Postal]([]*Postal{})
}
