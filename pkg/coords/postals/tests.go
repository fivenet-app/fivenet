package postals

import (
	"github.com/galexrt/fivenet/pkg/coords"
)

func NewForTests() Postals {
	return coords.New[*Postal]()
}
