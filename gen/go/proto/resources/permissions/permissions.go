package permissions

import (
	timestamp "github.com/galexrt/fivenet/gen/go/proto/resources/timestamp"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/query/fivenet/model"
)

func (x *Role) GetJobGrade() int32 {
	return x.GetGrade()
}

func (x *Role) SetJobLabel(label string) {
	x.JobLabel = &label
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

func (x *AttributeValues) Check(aType AttributeTypes, validVals *AttributeValues, maxVals *AttributeValues) bool {
	if validVals == nil && maxVals == nil {
		return true
	}

	switch AttributeTypes(aType) {
	case StringListAttributeType:
		if maxVals != nil {
			if !ValidateStringList(x.GetStringList().Strings, validVals.GetStringList().Strings, maxVals.GetStringList().Strings) {
				return false
			}
		} else {
			if !ValidateStringList(x.GetStringList().Strings, validVals.GetStringList().Strings, nil) {
				return false
			}
		}
	case JobListAttributeType:
		if maxVals != nil {
			if !ValidateJobList(x.GetJobList().Strings, validVals.GetJobList().Strings, maxVals.GetJobList().Strings) {
				return false
			}
		} else {
			if !ValidateJobList(x.GetJobList().Strings, validVals.GetJobList().Strings, nil) {
				return false
			}
		}
	case JobGradeListAttributeType:
		if maxVals != nil {
			if !ValidateJobGradeList(x.GetJobGradeList().Jobs, validVals.GetJobGradeList().Jobs, maxVals.GetJobGradeList().Jobs) {
				return false
			}
		} else {
			if !ValidateJobGradeList(x.GetJobGradeList().Jobs, validVals.GetJobGradeList().Jobs, nil) {
				return false
			}
		}
	}

	return true
}

func ValidateStringList(in []string, validVals []string, maxVals []string) bool {
	// If more values than valid/max values in the list, it can't be valid
	if len(in) > len(validVals) || len(in) > len(maxVals) {
		return false
	}

	for i := 0; i < len(in); i++ {
		if !utils.InStringSlice(maxVals, in[i]) {
			return false
		}

		if validVals != nil && !utils.InStringSlice(validVals, in[i]) {
			return false
		}
	}

	return true
}

func ValidateJobList(in []string, validVals []string, maxVals []string) bool {
	for i := 0; i < len(in); i++ {
		if !utils.InStringSlice(maxVals, in[i]) {
			return false
		}

		if validVals != nil && !utils.InStringSlice(validVals, in[i]) {
			// Remove invalid jobs from list
			utils.RemoveFromStringSlice(in, i)
		}
	}

	return true
}

func ValidateJobGradeList(in map[string]int32, validVals map[string]int32, maxVals map[string]int32) bool {
	// TODO validate job grade list, valid vals will contain one rank and that is the "highest" it can be
	for job, grade := range in {
		if vg, ok := maxVals[job]; ok {
			if grade > vg {
				return false
			}
		} else {
			return false
		}

		if validVals != nil {
			if vg, ok := validVals[job]; ok {
				if grade > vg {
					return false
				}
			} else {
				return false
			}
		}
	}

	return true
}
