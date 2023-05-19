package perms

import "github.com/galexrt/fivenet/pkg/utils"

type AttributeTypes string

const (
	StringListAttributeType   AttributeTypes = "StringList"
	JobListAttributeType      AttributeTypes = "JobList"
	JobGradeListAttributeType AttributeTypes = "JobGradeList"
)

type Key string

type StringList []string

type JobList []string
type JobGradeList map[string]int32

func ValidateStringList(in []string, validVals []string) bool {
	// If more values than valid values in the list, it can't be valid
	if len(in) > len(validVals) {
		return false
	}

	for i := 0; i < len(in); i++ {
		if !utils.InStringSlice(validVals, in[i]) {
			return false
		}
	}

	return true
}

func ValidateJobList(in []string, jobs []string) bool {
	for k, v := range in {
		if !utils.InStringSlice(jobs, v) {
			// Remove invalid jobs from list
			utils.RemoveFromStringSlice(in, k)
		}
	}

	return true
}

func ValidateJobGradeList(in map[string]int32) bool {

	// TODO validate job grade list, valid vals will contain one rank and that is the "highest" it can have

	return true
}
