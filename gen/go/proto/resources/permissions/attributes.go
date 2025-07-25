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
		if !slices.Contains(x.Strings, item) {
			return false
		}
	}

	return true
}

func (x *StringList) Len() int {
	if x == nil || x.Strings == nil {
		return 0
	}

	return len(x.Strings)
}

func (x *JobGradeList) HasJobGrade(job string, grade int32) bool {
	if x == nil {
		return false
	}

	if x.FineGrained {
		if x.Grades == nil {
			return false
		}

		// Check if the job exists in the list and the grade is allowed in the fine grained list
		grades, ok := x.Grades[job]
		if !ok {
			return false
		}
		if grades == nil || len(grades.Grades) == 0 {
			return false
		}

		return slices.Contains(grades.Grades, grade)
	} else {
		if x.Jobs == nil {
			return false
		}

		// Check if the job exists in the list and the grade is "in range"
		if g, ok := x.Jobs[job]; ok {
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

	if x.FineGrained {
		return len(x.Grades)
	}

	return len(x.Jobs)
}

func (x *JobGradeList) Iter() iter.Seq2[string, []int32] {
	return func(yield func(string, []int32) bool) {
		if x.FineGrained {
			for job, v := range x.Grades {
				if !yield(job, v.Grades) {
					return
				}
			}
		} else {
			for job, v := range x.Jobs {
				if !yield(job, []int32{v}) {
					return
				}
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
	for i := len(in.Strings) - 1; i >= 0; i-- {
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
	for i := len(in.Strings) - 1; i >= 0; i-- {
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
	if !in.FineGrained {
		if len(in.Grades) > 0 {
			in.Grades = map[string]*JobGrades{}
			changed = true
		}

		for job, grade := range in.Jobs {
			if vg, ok := maxVals[job]; ok {
				if vg > 0 {
					if grade > vg {
						in.Jobs[job] = vg
						changed = true
					}
				} else {
					// Valid grade for job is less than 0, remove job (invalid input case)
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
					if vg > 0 {
						if grade > vg {
							in.Jobs[job] = vg
							changed = true
						}
					} else {
						// Valid grade for job is less than 0, remove job (invalid input case)
						delete(in.Jobs, job)
						changed = true
					}
				} else {
					delete(in.Jobs, job)
					changed = true
				}
			}
		}
	} else {
		if len(in.Jobs) > 0 {
			in.Jobs = map[string]int32{}
		}

		for job, grades := range in.Grades {
			if grades == nil || len(grades.Grades) == 0 {
				delete(in.Grades, job)
				changed = true
				continue
			}

			for _, grade := range grades.Grades {
				if vg, ok := maxVals[job]; ok {
					currentLen := len(in.Grades[job].Grades)
					// Remove all grades that are greater than the max grade
					in.Grades[job].Grades = slices.DeleteFunc(in.Grades[job].Grades, func(ig int32) bool {
						return grade > vg
					})

					if currentLen != len(in.Grades[job].Grades) {
						changed = true
					}
				} else {
					delete(in.Grades, job)
					changed = true
				}

				// If valid vals are empty/ nil, don't check them
				if len(validVals) > 0 {
					if vg, ok := validVals[job]; ok {
						currentLen := len(in.Grades[job].Grades)
						// Remove all grades that are greater than the max grade
						in.Grades[job].Grades = slices.DeleteFunc(in.Grades[job].Grades, func(ig int32) bool {
							return grade > vg
						})

						if currentLen != len(in.Grades[job].Grades) {
							changed = true
						}
					} else {
						delete(in.Grades, job)
						changed = true
					}
				}
			}
		}
	}

	return true, changed
}
