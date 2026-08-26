package usersprops

import (
	citizenslabels "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/citizens/labels"
)

func (x *UserProps) Default() {
	if x.Wanted == nil {
		v := false
		x.Wanted = &v
	}

	if x.TrafficInfractionPoints == nil {
		v := uint32(0)
		x.TrafficInfractionPoints = &v
	}

	if x.OpenFines == nil {
		v := int64(0)
		x.OpenFines = &v
	}

	if x.GetLabels() == nil {
		x.Labels = &citizenslabels.Labels{}
	}
}
