package permissionsattributes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateStringList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		in          *StringList
		validVals   []string
		maxVals     []string
		wantStrings []string
		wantChanged bool
	}{
		{
			name:        "clears_when_input_exceeds_max_len",
			in:          &StringList{Strings: []string{"a", "b", "c"}},
			maxVals:     []string{"a", "b"},
			wantStrings: []string{},
			wantChanged: true,
		},
		{
			name:        "removes_values_not_in_max",
			in:          &StringList{Strings: []string{"a", "x", "b"}},
			maxVals:     []string{"a", "b", "c"},
			wantStrings: []string{"a", "b"},
			wantChanged: true,
		},
		{
			name:        "removes_values_not_in_valid_when_provided",
			in:          &StringList{Strings: []string{"a", "b", "c"}},
			validVals:   []string{"a", "c"},
			maxVals:     []string{"a", "b", "c"},
			wantStrings: []string{},
			wantChanged: true,
		},
		{
			name:        "unchanged_when_all_values_allowed",
			in:          &StringList{Strings: []string{"a", "b"}},
			validVals:   []string{"a", "b", "c"},
			maxVals:     []string{"a", "b", "c"},
			wantStrings: []string{"a", "b"},
			wantChanged: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ok, changed := ValidateStringList(tt.in, tt.validVals, tt.maxVals)

			assert.True(t, ok)
			assert.Equal(t, tt.wantChanged, changed)
			assert.Equal(t, tt.wantStrings, tt.in.GetStrings())
		})
	}
}

func TestValidateJobList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		in          *StringList
		validVals   []string
		maxVals     []string
		wantStrings []string
		wantChanged bool
	}{
		{
			name:        "filters_like_string_list",
			in:          &StringList{Strings: []string{"police", "invalid"}},
			validVals:   []string{"police", "ems"},
			maxVals:     []string{"police", "ems", "doj"},
			wantStrings: []string{"police"},
			wantChanged: true,
		},
		{
			name:        "unchanged_when_valid",
			in:          &StringList{Strings: []string{"police"}},
			validVals:   []string{"police"},
			maxVals:     []string{"police", "ems"},
			wantStrings: []string{"police"},
			wantChanged: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ok, changed := ValidateJobList(tt.in, tt.validVals, tt.maxVals)

			assert.True(t, ok)
			assert.Equal(t, tt.wantChanged, changed)
			assert.Equal(t, tt.wantStrings, tt.in.GetStrings())
		})
	}
}

func TestValidateJobGradeListNonFineGrained(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		in            *JobGradeList
		validVals     map[string]int32
		maxVals       map[string]int32
		wantJobs      map[string]int32
		wantGrades    map[string][]int32
		wantChanged   bool
		wantFineGrain bool
	}{
		{
			name: "unchanged_when_jobs_are_in_range",
			in: &JobGradeList{
				FineGrained: false,
				Jobs:        map[string]int32{"police": 2},
				Grades:      map[string]*JobGrades{},
			},
			validVals:     map[string]int32{"police": 3},
			maxVals:       map[string]int32{"police": 4},
			wantJobs:      map[string]int32{"police": 2},
			wantGrades:    map[string][]int32{},
			wantChanged:   false,
			wantFineGrain: false,
		},
		{
			name: "caps_and_removes_jobs_and_clears_fine_map",
			in: &JobGradeList{
				FineGrained: false,
				Jobs:        map[string]int32{"police": 6, "ems": 2, "invalid": 1},
				Grades: map[string]*JobGrades{
					"police": {Grades: []int32{1, 2}},
				},
			},
			validVals:     map[string]int32{"police": 4, "ems": 2},
			maxVals:       map[string]int32{"police": 5, "ems": 3},
			wantJobs:      map[string]int32{"police": 4, "ems": 2},
			wantGrades:    map[string][]int32{},
			wantChanged:   true,
			wantFineGrain: false,
		},
		{
			name: "removes_job_with_non_positive_max",
			in: &JobGradeList{
				FineGrained: false,
				Jobs:        map[string]int32{"police": 1},
				Grades:      map[string]*JobGrades{},
			},
			maxVals:       map[string]int32{"police": 0},
			wantJobs:      map[string]int32{},
			wantGrades:    map[string][]int32{},
			wantChanged:   true,
			wantFineGrain: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ok, changed := ValidateJobGradeList(tt.in, tt.validVals, tt.maxVals)

			assert.True(t, ok)
			assert.Equal(t, tt.wantChanged, changed)
			assert.Equal(t, tt.wantFineGrain, tt.in.GetFineGrained())
			assert.Equal(t, tt.wantJobs, tt.in.GetJobs())
			assert.Equal(t, tt.wantGrades, flattenGrades(tt.in.GetGrades()))
		})
	}
}

func TestValidateJobGradeListFineGrained(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		in          *JobGradeList
		validVals   map[string]int32
		maxVals     map[string]int32
		wantJobs    map[string]int32
		wantGrades  map[string][]int32
		wantChanged bool
	}{
		{
			name: "unchanged_when_all_grades_are_valid",
			in: &JobGradeList{
				FineGrained: true,
				Jobs:        map[string]int32{},
				Grades: map[string]*JobGrades{
					"police": {Grades: []int32{1, 2}},
				},
			},
			validVals:   map[string]int32{"police": 3},
			maxVals:     map[string]int32{"police": 3},
			wantJobs:    map[string]int32{},
			wantGrades:  map[string][]int32{"police": {1, 2}},
			wantChanged: false,
		},
		{
			name: "filters_grades_and_removes_jobs_without_valid_entries",
			in: &JobGradeList{
				FineGrained: true,
				Jobs:        map[string]int32{"to-clear": 1},
				Grades: map[string]*JobGrades{
					"police": {Grades: []int32{1, 4, 5}},
					"ems":    {Grades: []int32{}},
					"doj":    {Grades: []int32{1}},
				},
			},
			validVals:   map[string]int32{"police": 3},
			maxVals:     map[string]int32{"police": 4, "doj": 2},
			wantJobs:    map[string]int32{},
			wantGrades:  map[string][]int32{"police": {1}},
			wantChanged: true,
		},
		{
			name: "removes_job_missing_in_max_values",
			in: &JobGradeList{
				FineGrained: true,
				Jobs:        map[string]int32{},
				Grades: map[string]*JobGrades{
					"invalid": {Grades: []int32{1}},
				},
			},
			maxVals:     map[string]int32{"police": 3},
			wantJobs:    map[string]int32{},
			wantGrades:  map[string][]int32{},
			wantChanged: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ok, changed := ValidateJobGradeList(tt.in, tt.validVals, tt.maxVals)

			assert.True(t, ok)
			assert.Equal(t, tt.wantChanged, changed)
			assert.Equal(t, tt.wantJobs, tt.in.GetJobs())
			assert.Equal(t, tt.wantGrades, flattenGrades(tt.in.GetGrades()))
		})
	}
}

func TestAttributeValuesCheckPerType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		aType          AttributeTypes
		in             *AttributeValues
		validVals      *AttributeValues
		maxVals        *AttributeValues
		wantOk         bool
		wantChanged    bool
		wantStringList []string
		wantJobList    []string
		wantJobs       map[string]int32
		wantGrades     map[string][]int32
	}{
		{
			name:  "string_list_type",
			aType: StringListAttributeType,
			in: &AttributeValues{
				ValidValues: &AttributeValues_StringList{
					StringList: &StringList{Strings: []string{"a", "x"}},
				},
			},
			validVals: &AttributeValues{
				ValidValues: &AttributeValues_StringList{
					StringList: &StringList{Strings: []string{"a", "b"}},
				},
			},
			maxVals: &AttributeValues{
				ValidValues: &AttributeValues_StringList{
					StringList: &StringList{Strings: []string{"a", "b", "c"}},
				},
			},
			wantOk:         true,
			wantChanged:    true,
			wantStringList: []string{"a"},
		},
		{
			name:  "job_list_type",
			aType: JobListAttributeType,
			in: &AttributeValues{
				ValidValues: &AttributeValues_JobList{
					JobList: &StringList{Strings: []string{"police", "x"}},
				},
			},
			validVals: &AttributeValues{
				ValidValues: &AttributeValues_JobList{
					JobList: &StringList{Strings: []string{"police"}},
				},
			},
			maxVals: &AttributeValues{
				ValidValues: &AttributeValues_JobList{
					JobList: &StringList{Strings: []string{"police", "ems"}},
				},
			},
			wantOk:      true,
			wantChanged: true,
			wantJobList: []string{},
		},
		{
			name:  "job_grade_list_type",
			aType: JobGradeListAttributeType,
			in: &AttributeValues{
				ValidValues: &AttributeValues_JobGradeList{JobGradeList: &JobGradeList{
					FineGrained: false,
					Jobs:        map[string]int32{"police": 6},
					Grades:      map[string]*JobGrades{},
				}},
			},
			validVals: &AttributeValues{
				ValidValues: &AttributeValues_JobGradeList{JobGradeList: &JobGradeList{
					Jobs: map[string]int32{"police": 4},
				}},
			},
			maxVals: &AttributeValues{
				ValidValues: &AttributeValues_JobGradeList{JobGradeList: &JobGradeList{
					Jobs: map[string]int32{"police": 5},
				}},
			},
			wantOk:      true,
			wantChanged: true,
			wantJobs:    map[string]int32{"police": 4},
			wantGrades:  map[string][]int32{},
		},
		{
			name:        "unknown_type",
			aType:       AttributeTypes("Unknown"),
			in:          &AttributeValues{},
			validVals:   &AttributeValues{},
			wantOk:      false,
			wantChanged: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ok, changed := tt.in.Check(tt.aType, tt.validVals, tt.maxVals)

			assert.Equal(t, tt.wantOk, ok)
			assert.Equal(t, tt.wantChanged, changed)

			if tt.wantStringList != nil {
				assert.Equal(t, tt.wantStringList, tt.in.GetStringList().GetStrings())
			}
			if tt.wantJobList != nil {
				assert.Equal(t, tt.wantJobList, tt.in.GetJobList().GetStrings())
			}
			if tt.wantJobs != nil {
				assert.Equal(t, tt.wantJobs, tt.in.GetJobGradeList().GetJobs())
			}
			if tt.wantGrades != nil {
				assert.Equal(t, tt.wantGrades, flattenGrades(tt.in.GetJobGradeList().GetGrades()))
			}
		})
	}
}

func flattenGrades(in map[string]*JobGrades) map[string][]int32 {
	out := make(map[string][]int32, len(in))
	for job, grades := range in {
		if grades == nil {
			out[job] = nil
			continue
		}

		out[job] = append([]int32(nil), grades.GetGrades()...)
	}

	return out
}
