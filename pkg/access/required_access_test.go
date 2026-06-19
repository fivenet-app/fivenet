package access

import (
	"testing"

	resourcesaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/access"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeAccessSeedsFallbackWhenEmpty(t *testing.T) {
	t.Parallel()

	current := &resourcesaccess.Access{}
	fallback := &resourcesaccess.Access{
		Jobs: []*resourcesaccess.JobAccess{
			{
				Job:          "police",
				MinimumGrade: 3,
				Access:       7,
			},
		},
	}

	out, err := NormalizeAccess(current, nil, fallback, 10)
	require.NoError(t, err)
	require.NotNil(t, out)
	require.Len(t, out.GetJobs(), 1)
	assert.Equal(t, "police", out.GetJobs()[0].GetJob())
	assert.Equal(t, int32(3), out.GetJobs()[0].GetMinimumGrade())
	assert.Equal(t, int32(7), out.GetJobs()[0].GetAccess())
}

func TestApplyRequiredAccessOverlayUpgradesExistingEntry(t *testing.T) {
	t.Parallel()

	current := &resourcesaccess.Access{
		Jobs: []*resourcesaccess.JobAccess{
			{
				Job:            "police",
				MinimumGrade:   3,
				Access:         2,
				Required:       new(true),
				RequiredAccess: func() *int32 { v := int32(2); return &v }(),
			},
		},
	}
	required := &resourcesaccess.Access{
		Jobs: []*resourcesaccess.JobAccess{
			{
				Job:            "police",
				MinimumGrade:   3,
				Access:         7,
				Required:       new(true),
				RequiredAccess: func() *int32 { v := int32(2); return &v }(),
			},
		},
	}

	out, err := ApplyRequiredAccessOverlay(current, required, 10)
	require.NoError(t, err)
	require.Len(t, out.GetJobs(), 1)
	assert.True(t, out.GetJobs()[0].GetRequired())
	assert.Equal(t, int32(7), out.GetJobs()[0].GetAccess())
	assert.Equal(t, int32(2), out.GetJobs()[0].GetRequiredAccess())
}

func TestApplyRequiredAccessOverlayInsertsMissingEntry(t *testing.T) {
	t.Parallel()

	current := &resourcesaccess.Access{}
	required := &resourcesaccess.Access{
		Users: []*resourcesaccess.UserAccess{
			{
				UserId:         42,
				Access:         8,
				Required:       new(true),
				RequiredAccess: func() *int32 { v := int32(8); return &v }(),
			},
		},
	}

	out, err := ApplyRequiredAccessOverlay(current, required, 10)
	require.NoError(t, err)
	require.Len(t, out.GetUsers(), 1)
	assert.Equal(t, int32(42), out.GetUsers()[0].GetUserId())
	assert.True(t, out.GetUsers()[0].GetRequired())
	assert.Equal(t, int32(8), out.GetUsers()[0].GetAccess())
	assert.Equal(t, int32(8), out.GetUsers()[0].GetRequiredAccess())
}

func TestNormalizeRequiredAccessFloorsBackfillsMissingFloor(t *testing.T) {
	t.Parallel()

	input := &resourcesaccess.Access{
		Jobs: []*resourcesaccess.JobAccess{
			{
				Job:          "police",
				MinimumGrade: 3,
				Access:       2,
				Required:     new(true),
			},
		},
		Users: []*resourcesaccess.UserAccess{
			{
				UserId:   42,
				Access:   5,
				Required: new(true),
			},
		},
		Qualifications: []*resourcesaccess.QualificationAccess{
			{
				QualificationId: 7,
				Access:          4,
				Required:        new(true),
			},
		},
	}

	out := NormalizeRequiredAccessFloors(input)
	require.NotNil(t, out)
	require.Len(t, out.GetJobs(), 1)
	require.Len(t, out.GetUsers(), 1)
	require.Len(t, out.GetQualifications(), 1)

	assert.Equal(t, int32(2), out.GetJobs()[0].GetRequiredAccess())
	assert.Equal(t, int32(5), out.GetUsers()[0].GetRequiredAccess())
	assert.Equal(t, int32(4), out.GetQualifications()[0].GetRequiredAccess())
}

func TestNormalizeRequiredAccessFloorsClampsBelowFloor(t *testing.T) {
	t.Parallel()

	input := &resourcesaccess.Access{
		Users: []*resourcesaccess.UserAccess{
			{
				UserId:         42,
				Access:         1,
				Required:       new(true),
				RequiredAccess: func() *int32 { v := int32(3); return &v }(),
			},
		},
	}

	out := NormalizeRequiredAccessFloors(input)
	require.NotNil(t, out)
	require.Len(t, out.GetUsers(), 1)
	assert.Equal(t, int32(3), out.GetUsers()[0].GetAccess())
	assert.Equal(t, int32(3), out.GetUsers()[0].GetRequiredAccess())
}

func TestApplyRequiredAccessOverlayClampsAccessToFloor(t *testing.T) {
	t.Parallel()

	floor := int32(2)
	current := &resourcesaccess.Access{
		Qualifications: []*resourcesaccess.QualificationAccess{
			{
				QualificationId: 7,
				Access:          1,
				Required:        new(true),
				RequiredAccess:  &floor,
			},
		},
	}
	required := &resourcesaccess.Access{
		Qualifications: []*resourcesaccess.QualificationAccess{
			{
				QualificationId: 7,
				Access:          2,
				Required:        new(true),
				RequiredAccess:  &floor,
			},
		},
	}

	out, err := ApplyRequiredAccessOverlay(current, required, 10)
	require.NoError(t, err)
	require.Len(t, out.GetQualifications(), 1)
	assert.Equal(t, int32(2), out.GetQualifications()[0].GetAccess())
	assert.Equal(t, int32(2), out.GetQualifications()[0].GetRequiredAccess())
}

func TestNormalizeAccessReturnsLimitErrorWhenRequiredWouldExceedLimit(t *testing.T) {
	t.Parallel()

	current := &resourcesaccess.Access{
		Users: []*resourcesaccess.UserAccess{
			{
				UserId: 1,
				Access: 2,
			},
		},
	}
	required := &resourcesaccess.Access{
		Users: []*resourcesaccess.UserAccess{
			{
				UserId:   2,
				Access:   7,
				Required: new(true),
			},
		},
	}

	out, err := NormalizeAccess(current, required, nil, 1)
	require.Error(t, err)
	require.Nil(t, out)

	var limitErr *AccessEntryLimitError
	require.ErrorAs(t, err, &limitErr)
	assert.Equal(t, "users", limitErr.Kind)
	assert.Equal(t, 1, limitErr.Max)
	assert.Equal(t, 2, limitErr.Actual)
}

func TestSanitizeJobAccessEntriesDropsMissingJobsWithoutMutatingInput(t *testing.T) {
	t.Parallel()

	input := []*resourcesaccess.JobAccess{
		{
			Job:          "police",
			MinimumGrade: 3,
			Access:       2,
		},
		{
			Job:          "ghost",
			MinimumGrade: 1,
			Access:       9,
		},
	}

	jobs := mstlystcdata.NewDummyJobs(map[string]*jobs.Job{
		"police": {
			Name: "police",
			Grades: []*jobs.JobGrade{
				{Grade: 0},
				{Grade: 3},
				{Grade: 5},
			},
		},
	})

	out, err := SanitizeJobAccessEntries(jobs, input)
	require.NoError(t, err)
	require.Len(t, out, 1)
	assert.Equal(t, "police", out[0].GetJob())
	assert.Equal(t, int32(3), out[0].GetMinimumGrade())
	assert.Equal(t, int32(2), out[0].GetAccess())

	require.Len(t, input, 2)
	assert.Equal(t, "ghost", input[1].GetJob())
}

func TestSanitizeJobAccessEntriesRepairsInvalidGrade(t *testing.T) {
	t.Parallel()

	input := []*resourcesaccess.JobAccess{
		{
			Job:          "police",
			MinimumGrade: 99,
			Access:       2,
		},
	}

	jobs := mstlystcdata.NewDummyJobs(map[string]*jobs.Job{
		"police": {
			Name: "police",
			Grades: []*jobs.JobGrade{
				{Grade: 0},
				{Grade: 3},
			},
		},
	})

	out, err := SanitizeJobAccessEntries(jobs, input)
	require.NoError(t, err)
	require.Len(t, out, 1)
	assert.Equal(t, int32(3), out[0].GetMinimumGrade())
}

func TestSanitizeAccessJobsDropsStaleJobEntries(t *testing.T) {
	t.Parallel()

	input := &resourcesaccess.Access{
		Jobs: []*resourcesaccess.JobAccess{
			{
				Job:          "ghost",
				MinimumGrade: 1,
				Access:       9,
			},
		},
		Users: []*resourcesaccess.UserAccess{
			{
				UserId:   42,
				Access:   3,
				Required: new(true),
			},
		},
	}

	out, err := SanitizeAccessJobs(mstlystcdata.NewDummyJobs(nil), input)
	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Empty(t, out.GetJobs())
	require.Len(t, out.GetUsers(), 1)
	assert.Equal(t, int32(42), out.GetUsers()[0].GetUserId())
	assert.True(t, out.GetUsers()[0].GetRequired())
}
