package qualifications

import (
	"testing"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs"
	qualificationsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications/access"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/stretchr/testify/require"
)

func TestQualificationCreationJobAccessIncludesCreatorAndHighestGrade(t *testing.T) {
	t.Parallel()

	entries := qualificationCreationJobAccess(
		&userinfo.UserInfo{Job: "ambulance", JobGrade: 4},
		&jobs.Job{
			Name:   "ambulance",
			Grades: []*jobs.JobGrade{{Grade: 4}, {Grade: 16}},
		},
		42,
	)

	require.Len(t, entries, 2)
	require.Equal(t, int32(4), entries[0].GetMinimumGrade())
	require.Equal(t, int32(16), entries[1].GetMinimumGrade())
	require.True(t, entries[1].GetRequired())
	require.Equal(t, int32(qualificationsaccess.AccessLevel_ACCESS_LEVEL_EDIT), entries[0].GetAccess())
	require.Equal(t, int32(qualificationsaccess.AccessLevel_ACCESS_LEVEL_EDIT), entries[1].GetAccess())
}

func TestQualificationCreationJobAccessDoesNotDuplicateHighestGrade(t *testing.T) {
	t.Parallel()

	entries := qualificationCreationJobAccess(
		&userinfo.UserInfo{Job: "ambulance", JobGrade: 16},
		&jobs.Job{
			Name:   "ambulance",
			Grades: []*jobs.JobGrade{{Grade: 4}, {Grade: 16}},
		},
		42,
	)

	require.Len(t, entries, 1)
	require.Equal(t, int32(16), entries[0].GetMinimumGrade())
	require.True(t, entries[0].GetRequired())
}
