package documents

import "slices"

func (x *SignatureTask) GetJobGrade() int32 {
	return x.GetMinimumGrade()
}

func (x *SignatureTask) SetJobGrade(grade int32) {
	x.MinimumGrade = &grade
}

func (x *SignatureTask) SetJobGradeLabel(label string) {
	x.JobGradeLabel = &label
}

func (x *SignatureTask) SetJob(job string) {
	x.Job = &job
}

func (x *SignatureTask) SetJobLabel(label string) {
	x.JobLabel = &label
}

// HasType In case no types are specified? All types are allowed.
func (x *SignatureTypes) HasType(t SignatureType) bool {
	if x == nil {
		return true
	}

	return slices.Contains(x.Types, t)
}
