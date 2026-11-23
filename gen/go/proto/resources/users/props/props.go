package usersprops

import (
	citizenslabels "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/citizens/labels"
	properties "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
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

// TemplateDefault fills the optional scalar values that legacy document
// templates access as fields. Go's template package treats an unset optional
// protobuf field as a nil-pointer dereference, rather than rendering it as an
// empty value. Keep this separate from Default because the latter is also
// used by API hydration, where field presence conveys property visibility.
func (x *UserProps) TemplateDefault() {
	if x == nil {
		return
	}

	x.Default()

	if x.JobName == nil {
		x.JobName = new(string)
	}
	if x.JobGradeNumber == nil {
		x.JobGradeNumber = new(int32)
	}
	if x.BloodType == nil {
		x.BloodType = new(string)
	}
	if x.MugshotFileId == nil {
		x.MugshotFileId = new(int64)
	}
	if x.Email == nil {
		x.Email = new(string)
	}
	if x.Properties == nil {
		x.Properties = &properties.Properties{}
	}
}
