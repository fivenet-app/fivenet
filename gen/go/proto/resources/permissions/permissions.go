package permissions

import (
	"slices"

	timestamp "github.com/fivenet-app/fivenet/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
)

func (x *Role) GetJobGrade() int32 {
	return x.GetGrade()
}

func (x *Role) SetJob(job string) {
	x.Job = job
}

func (x *Role) SetJobLabel(label string) {
	x.JobLabel = &label
}

func (x *Role) SetJobGrade(grade int32) {
	x.Grade = grade
}

func (x *Role) SetJobGradeLabel(label string) {
	x.JobGradeLabel = &label
}

func ConvertFromRole(p *model.FivenetRoles) *Role {
	var createdAt *timestamp.Timestamp
	if p.CreatedAt != nil {
		createdAt = timestamp.New(*p.CreatedAt)
	}

	return &Role{
		Id:        p.ID,
		CreatedAt: createdAt,
		Job:       p.Job,
		Grade:     p.Grade,
	}
}

type AttributeTypes string

const (
	StringListAttributeType   AttributeTypes = "StringList"
	JobListAttributeType      AttributeTypes = "JobList"
	JobGradeListAttributeType AttributeTypes = "JobGradeList"
)

func (x *AttributeValues) Default(aType AttributeTypes) {
	switch AttributeTypes(aType) {
	case StringListAttributeType:
		if x.GetStringList() == nil || x.GetStringList().Strings == nil {
			x.ValidValues = &AttributeValues_StringList{
				StringList: &StringList{
					Strings: nil,
				},
			}
		}

	case JobListAttributeType:
		if x.GetJobList() == nil || x.GetJobList().Strings == nil {
			x.ValidValues = &AttributeValues_JobList{
				JobList: &StringList{
					Strings: nil,
				},
			}
		}

	case JobGradeListAttributeType:
		if x.GetJobGradeList() == nil || x.GetJobGradeList().Jobs == nil {
			x.ValidValues = &AttributeValues_JobGradeList{
				JobGradeList: &JobGradeList{
					Jobs: nil,
				},
			}
		}
	}
}

func (x *AttributeValues) Check(aType AttributeTypes, validVals *AttributeValues, maxVals *AttributeValues) (bool, bool) {
	if validVals == nil && maxVals == nil {
		return true, false
	}

	switch AttributeTypes(aType) {
	case StringListAttributeType:
		var valid []string
		if validVals != nil && validVals.GetStringList() != nil && validVals.GetStringList().Strings != nil {
			valid = validVals.GetStringList().Strings
		}
		var max []string
		if maxVals != nil && maxVals.GetStringList() != nil && maxVals.GetStringList().Strings != nil {
			max = maxVals.GetStringList().Strings
		}

		return ValidateStringList(x.GetStringList(), valid, max)

	case JobListAttributeType:
		var valid []string
		if validVals != nil && validVals.GetJobList() != nil && validVals.GetJobList().Strings != nil {
			valid = validVals.GetJobList().Strings
		}
		var max []string
		if maxVals != nil && maxVals.GetJobList() != nil && maxVals.GetJobList().Strings != nil {
			max = maxVals.GetJobList().Strings
		}

		return ValidateJobList(x.GetJobList(), valid, max)

	case JobGradeListAttributeType:
		var valid map[string]int32

		if validVals != nil && validVals.GetJobGradeList() != nil && validVals.GetJobGradeList().Jobs != nil {
			valid = validVals.GetJobGradeList().Jobs
		}
		var max map[string]int32
		if maxVals != nil && maxVals.GetJobGradeList() != nil && maxVals.GetJobGradeList().Jobs != nil {
			max = maxVals.GetJobGradeList().Jobs
		}

		return ValidateJobGradeList(x.GetJobGradeList(), valid, max)

	default:
		return false, false
	}
}

func ValidateStringList(in *StringList, validVals []string, maxVals []string) (bool, bool) {
	// If more values than valid/max values in the list, it can't be valid
	if len(in.Strings) > len(maxVals) || (len(validVals) > 0 && len(in.Strings) > len(validVals)) {
		in.Strings = []string{}
		return true, true
	}

	changed := false
	for i := range in.Strings {
		if !slices.Contains(maxVals, in.Strings[i]) {
			in.Strings = slices.Delete(in.Strings, i, i+1)
			changed = true
			continue
		}

		if len(validVals) > 0 && !slices.Contains(validVals, in.Strings[i]) {
			in.Strings = slices.Delete(in.Strings, i, i+1)
			changed = true
			continue
		}
	}

	return true, changed
}

func ValidateJobList(in *StringList, validVals []string, maxVals []string) (bool, bool) {
	// If more values than valid/max values in the list, it can't be valid
	if len(in.Strings) > len(maxVals) || (len(validVals) > 0 && len(in.Strings) > len(validVals)) {
		in.Strings = []string{}
		return true, true
	}

	changed := false
	for i := range in.Strings {
		if !slices.Contains(maxVals, in.Strings[i]) {
			in.Strings = slices.Delete(in.Strings, i, i+1)
			changed = true
			continue
		}

		if len(validVals) > 0 && !slices.Contains(validVals, in.Strings[i]) {
			// Remove invalid jobs from list
			in.Strings = slices.Delete(in.Strings, i, i+1)
			changed = true
			continue
		}
	}

	return true, changed
}

func ValidateJobGradeList(in *JobGradeList, validVals map[string]int32, maxVals map[string]int32) (bool, bool) {
	changed := false
	for job, grade := range in.Jobs {
		if vg, ok := maxVals[job]; ok {
			if grade > vg {
				in.Jobs[job] = vg
				changed = true
			}
		} else {
			delete(in.Jobs, job)
			changed = true
		}

		// If valid vals are empty/ nil, don't check them
		if len(validVals) > 0 {
			if vg, ok := validVals[job]; ok {
				if grade > vg {
					in.Jobs[job] = vg
					changed = true
				}
			} else {
				delete(in.Jobs, job)
				changed = true
			}
		}
	}

	return true, changed
}
