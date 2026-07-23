package utils

import (
	"testing"

	syncdata "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/sync/data"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeUserJobsUsesJobsWhenScalarJobEmpty(t *testing.T) {
	t.Parallel()

	user := &syncdata.DataUser{
		UserId:   11,
		Job:      "",
		JobGrade: 0,
		Jobs: []*users.UserJob{
			{Job: "", Grade: 0, IsPrimary: false},
			{Job: "ems", Grade: 1, IsPrimary: true},
			{Job: "police", Grade: 3, IsPrimary: false},
		},
	}

	NormalizeUserJobs(user, "", 0)

	assert.Equal(t, "ems", user.GetJob())
	assert.Equal(t, int32(1), user.GetJobGrade())
	require.Len(t, user.GetJobs(), 2)
	assert.Equal(t, "ems", user.GetJobs()[0].GetJob())
	assert.True(t, user.GetJobs()[0].GetIsPrimary())
	assert.Equal(t, "police", user.GetJobs()[1].GetJob())
	assert.False(t, user.GetJobs()[1].GetIsPrimary())
}
