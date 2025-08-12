package permissions

import (
	"iter"
	"slices"
)

type AttributeTypes string

const (
	StringListAttributeType   AttributeTypes = "StringList"
	JobListAttributeType      AttributeTypes = "JobList"
	JobGradeListAttributeType AttributeTypes = "JobGradeList"
)

func (x *AttributeValues) Default(aType AttributeTypes) {
	switch aType {
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
		if x.GetJobGradeList() == nil {
			x.ValidValues = &AttributeValues_JobGradeList{
				JobGradeList: &JobGradeList{
					Jobs:        map[string]int32{},
					FineGrained: false,
					Grades:      map[string]*JobGrades{},
				},
			}
		}
		if x.GetJobGradeList().Jobs == nil {
			x.GetJobGradeList().Jobs = map[string]int32{}
		}
		if x.GetJobGradeList().Grades == nil {
			x.GetJobGradeList().Grades = map[string]*JobGrades{}
		}
	}
}

func (x *StringList) Contains(items ...string) bool {
	if x == nil || x.Strings == nil {
		return false
	}

	// Check if all items are in the list
	for _, item := range items {
		if !slices.Contains(x.GetStrings(), item) {
			return false
		}
	}

	return true
}

func (x *StringList) Len() int {
	if x == nil || x.Strings == nil {
		return 0
	}

	return len(x.GetStrings())
}

func (x *JobGradeList) HasJobGrade(job string, grade int32) bool {
	if x == nil {
		return false
	}

	if x.GetFineGrained() {
		if x.Grades == nil {
			return false
		}

		// Check if the job exists in the list and the grade is allowed in the fine grained list
		grades, ok := x.GetGrades()[job]
		if !ok {
			return false
		}
		if grades == nil || len(grades.GetGrades()) == 0 {
			return false
		}

		return slices.Contains(grades.GetGrades(), grade)
	} else {
		if x.Jobs == nil {
			return false
		}

		// Check if the job exists in the list and the grade is "in range"
		if g, ok := x.GetJobs()[job]; ok {
			if g >= grade {
				return true
			}
		}
	}

	return false
}

func (x *JobGradeList) Len() int {
	if x == nil {
		return 0
	}

	if x.GetFineGrained() {
		return len(x.GetGrades())
	}

	return len(x.GetJobs())
}

func (x *JobGradeList) Iter() iter.Seq2[string, []int32] {
	return func(yield func(string, []int32) bool) {
		if x.GetFineGrained() {
			for job, v := range x.GetGrades() {
				if !yield(job, v.GetGrades()) {
					return
				}
			}
		} else {
			for job, v := range x.GetJobs() {
				if !yield(job, []int32{v}) {
					return
				}
			}
		}
	}
}

func (x *AttributeValues) Check(
	aType AttributeTypes,
	validVals *AttributeValues,
	maxVals *AttributeValues,
) (bool, bool) {
	if validVals == nil && maxVals == nil {
		return true, false
	}

	switch aType {
	case StringListAttributeType:
		var valid []string
		if validVals != nil && validVals.GetStringList() != nil &&
			validVals.GetStringList().Strings != nil {
			valid = validVals.GetStringList().GetStrings()
		}
		var maxV []string
		if maxVals != nil && maxVals.GetStringList() != nil &&
			maxVals.GetStringList().Strings != nil {
			maxV = maxVals.GetStringList().GetStrings()
		}

		return ValidateStringList(x.GetStringList(), valid, maxV)

	case JobListAttributeType:
		var valid []string
		if validVals != nil && validVals.GetJobList() != nil &&
			validVals.GetJobList().Strings != nil {
			valid = validVals.GetJobList().GetStrings()
		}
		var maxV []string
		if maxVals != nil && maxVals.GetJobList() != nil && maxVals.GetJobList().Strings != nil {
			maxV = maxVals.GetJobList().GetStrings()
		}

		return ValidateJobList(x.GetJobList(), valid, maxV)

	case JobGradeListAttributeType:
		var valid map[string]int32

		if validVals != nil && validVals.GetJobGradeList() != nil &&
			validVals.GetJobGradeList().Jobs != nil {
			valid = validVals.GetJobGradeList().GetJobs()
		}
		var maxV map[string]int32
		if maxVals != nil && maxVals.GetJobGradeList() != nil &&
			maxVals.GetJobGradeList().Jobs != nil {
			maxV = maxVals.GetJobGradeList().GetJobs()
		}

		return ValidateJobGradeList(x.GetJobGradeList(), valid, maxV)

	default:
		return false, false
	}
}

func ValidateStringList(in *StringList, validVals []string, maxVals []string) (bool, bool) {
	// If more values than valid/max values in the list, it can't be valid
	if len(in.GetStrings()) > len(maxVals) ||
		(len(validVals) > 0 && len(in.GetStrings()) > len(validVals)) {
		in.Strings = []string{}
		return true, true
	}

	changed := false
	for i := len(in.GetStrings()) - 1; i >= 0; i-- {
		if !slices.Contains(maxVals, in.GetStrings()[i]) {
			in.Strings = slices.Delete(in.GetStrings(), i, i+1)
			changed = true
			continue
		}

		if len(validVals) > 0 && !slices.Contains(validVals, in.GetStrings()[i]) {
			in.Strings = slices.Delete(in.GetStrings(), i, i+1)
			changed = true
			continue
		}
	}

	return true, changed
}

func ValidateJobList(in *StringList, validVals []string, maxVals []string) (bool, bool) {
	return ValidateStringList(in, validVals, maxVals)
}

func ValidateJobGradeList(
	in *JobGradeList,
	validVals map[string]int32,
	maxVals map[string]int32,
) (bool, bool) {
	changed := false
	if !in.GetFineGrained() {
		if len(in.GetGrades()) > 0 {
			in.Grades = map[string]*JobGrades{}
			changed = true
		}

		for job, grade := range in.GetJobs() {
			if vg, ok := maxVals[job]; ok {
				if vg > 0 {
					if grade > vg {
						in.Jobs[job] = vg
						changed = true
					}
				} else {
					// Valid grade for job is less than 0, remove job (invalid input case)
					delete(in.GetJobs(), job)
					changed = true
				}
			} else {
				delete(in.GetJobs(), job)
				changed = true
			}

			// If valid vals are empty/ nil, don't check them
			if len(validVals) > 0 {
				if vg, ok := validVals[job]; ok {
					if vg > 0 {
						if grade > vg {
							in.Jobs[job] = vg
							changed = true
						}
					} else {
						// Valid grade for job is less than 0, remove job (invalid input case)
						delete(in.GetJobs(), job)
						changed = true
					}
				} else {
					delete(in.GetJobs(), job)
					changed = true
				}
			}
		}
	} else {
		if len(in.GetJobs()) > 0 {
			in.Jobs = map[string]int32{}
		}

		for job, grades := range in.GetGrades() {
			if grades == nil || len(grades.GetGrades()) == 0 {
				delete(in.GetGrades(), job)
				changed = true
				continue
			}

			for _, grade := range grades.GetGrades() {
				if vg, ok := maxVals[job]; ok {
					currentLen := len(in.GetGrades()[job].GetGrades())
					// Remove all grades that are greater than the max grade
					in.Grades[job].Grades = slices.DeleteFunc(in.GetGrades()[job].GetGrades(), func(ig int32) bool {
						return grade > vg
					})

					if currentLen != len(in.GetGrades()[job].GetGrades()) {
						changed = true
					}
				} else {
					delete(in.GetGrades(), job)
					changed = true
				}

				// If valid vals are empty/ nil, don't check them
				if len(validVals) > 0 {
					if vg, ok := validVals[job]; ok {
						currentLen := len(in.GetGrades()[job].GetGrades())
						// Remove all grades that are greater than the max grade
						in.Grades[job].Grades = slices.DeleteFunc(in.GetGrades()[job].GetGrades(), func(ig int32) bool {
							return grade > vg
						})

						if currentLen != len(in.GetGrades()[job].GetGrades()) {
							changed = true
						}
					} else {
						delete(in.GetGrades(), job)
						changed = true
					}
				}
			}
		}
	}

	return true, changed
}
