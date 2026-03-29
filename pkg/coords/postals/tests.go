package postals

import (
	"github.com/fivenet-app/fivenet/v2026/pkg/coords"
)

func NewForTests() (Postals, error) {
	cs, err := coords.NewReadOnly([]*Postal{})
	if err != nil {
		return nil, err
	}

	return &postalStore{
		CoordsRO: cs,
		byCode:   map[string]*Postal{},
	}, nil
}
