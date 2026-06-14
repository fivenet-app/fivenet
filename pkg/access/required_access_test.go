package access

import (
	"testing"

	resourcesaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/access"
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
				Job:          "police",
				MinimumGrade: 3,
				Access:       2,
			},
		},
	}
	required := &resourcesaccess.Access{
		Jobs: []*resourcesaccess.JobAccess{
			{
				Job:          "police",
				MinimumGrade: 3,
				Access:       7,
				Required:     boolPtr(true),
			},
		},
	}

	out, err := ApplyRequiredAccessOverlay(current, required, 10)
	require.NoError(t, err)
	require.Len(t, out.GetJobs(), 1)
	assert.True(t, out.GetJobs()[0].GetRequired())
	assert.Equal(t, int32(7), out.GetJobs()[0].GetAccess())
}

func TestApplyRequiredAccessOverlayInsertsMissingEntry(t *testing.T) {
	t.Parallel()

	current := &resourcesaccess.Access{}
	required := &resourcesaccess.Access{
		Users: []*resourcesaccess.UserAccess{
			{
				UserId:   42,
				Access:   8,
				Required: boolPtr(true),
			},
		},
	}

	out, err := ApplyRequiredAccessOverlay(current, required, 10)
	require.NoError(t, err)
	require.Len(t, out.GetUsers(), 1)
	assert.Equal(t, int32(42), out.GetUsers()[0].GetUserId())
	assert.True(t, out.GetUsers()[0].GetRequired())
	assert.Equal(t, int32(8), out.GetUsers()[0].GetAccess())
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
				Required: boolPtr(true),
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

func boolPtr(v bool) *bool {
	return &v
}
