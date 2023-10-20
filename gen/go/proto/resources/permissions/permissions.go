package permissions

import (
	timestamp "github.com/galexrt/fivenet/gen/go/proto/resources/timestamp"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/query/fivenet/model"
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
					Strings: []string{},
				},
			}
		}

	case JobListAttributeType:
		if x.GetJobList() == nil || x.GetJobList().Strings == nil {
			x.ValidValues = &AttributeValues_JobList{
				JobList: &StringList{
					Strings: []string{},
				},
			}
		}

	case JobGradeListAttributeType:
		if x.GetJobGradeList() == nil || x.GetJobGradeList().Jobs == nil {
			x.ValidValues = &AttributeValues_JobGradeList{
				JobGradeList: &JobGradeList{
					Jobs: map[string]int32{},
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
		if maxVals != nil {
			max = maxVals.GetStringList().Strings
		}

		return ValidateStringList(x.GetStringList(), valid, max)

	case JobListAttributeType:
		var valid []string
		if validVals != nil && validVals.GetJobList() != nil && validVals.GetJobList().Strings != nil {
			valid = validVals.GetJobList().Strings
		}
		var max []string
		if maxVals != nil {
			max = maxVals.GetJobList().Strings
		}

		return ValidateJobList(x.GetJobList(), valid, max)

	case JobGradeListAttributeType:
		var valid map[string]int32

		if validVals != nil && validVals.GetJobGradeList() != nil && validVals.GetJobGradeList().Jobs != nil {
			valid = validVals.GetJobGradeList().Jobs
		}
		var max map[string]int32
		if maxVals != nil {
			max = maxVals.GetJobGradeList().Jobs
		}

		return ValidateJobGradeList(x.GetJobGradeList(), valid, max)

	default:
		return false, false
	}
}

func ValidateStringList(in *StringList, validVals []string, maxVals []string) (bool, bool) {
	// If more values than valid/max values in the list, it can't be valid
	if (validVals != nil && len(in.Strings) > len(validVals)) || len(in.Strings) > len(maxVals) {
		in.Strings = []string{}
		return true, true
	}

	changed := false
	for i := 0; i < len(in.Strings); i++ {
		if !utils.InSlice(maxVals, in.Strings[i]) {
			in.Strings = utils.RemoveFromSlice(in.Strings, i)
			changed = true
			continue
		}

		if validVals != nil && !utils.InSlice(validVals, in.Strings[i]) {
			in.Strings = utils.RemoveFromSlice(in.Strings, i)
			changed = true
			continue
		}
	}

	return true, changed
}

func ValidateJobList(in *StringList, validVals []string, maxVals []string) (bool, bool) {
	// If more values than valid/max values in the list, it can't be valid
	if len(in.Strings) > len(maxVals) || (validVals != nil && len(in.Strings) > len(validVals)) {
		in.Strings = []string{}
		return true, true
	}

	changed := false
	for i := 0; i < len(in.Strings); i++ {
		if !utils.InSlice(maxVals, in.Strings[i]) {
			in.Strings = utils.RemoveFromSlice(in.Strings, i)
			changed = true
			continue
		}

		if validVals != nil && !utils.InSlice(validVals, in.Strings[i]) {
			// Remove invalid jobs from list
			in.Strings = utils.RemoveFromSlice(in.Strings, i)
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
				delete(in.Jobs, job)
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
					delete(in.Jobs, job)
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
