package syncers

import (
	"testing"

	syncdata "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/sync/data"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users"
	dbsyncconfig "github.com/fivenet-app/fivenet/v2026/pkg/dbsync/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestCleanupUserJobUsesJobsWhenScalarJobEmpty(t *testing.T) {
	t.Parallel()

	syncer := &UsersSync{}
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

	syncer.cleanupUserJob(user)

	assert.Equal(t, "ems", user.GetJob())
	assert.Equal(t, int32(1), user.GetJobGrade())
	require.Len(t, user.GetJobs(), 2)
	assert.Equal(t, "ems", user.GetJobs()[0].GetJob())
	assert.True(t, user.GetJobs()[0].GetIsPrimary())
	assert.Equal(t, "police", user.GetJobs()[1].GetJob())
	assert.False(t, user.GetJobs()[1].GetIsPrimary())
}

func TestApplyFiltersAndTransformationsSkipsInvalidIdentity(t *testing.T) {
	t.Parallel()

	syncer := &UsersSync{
		Syncer: &Syncer{
			logger: zap.NewNop(),
			cfg:    &dbsyncconfig.DBSyncConfig{},
		},
		logger: zap.NewNop(),
	}
	users := []*syncdata.DataUser{
		nil,
		{Identifier: "license:zero"},
		{UserId: -1, Identifier: "license:negative"},
		{UserId: 43},
		{UserId: 44, Identifier: " \t"},
		{UserId: 42, Identifier: "license:valid"},
	}

	filtered := syncer.applyFiltersAndTransformations(users, dbsyncconfig.UsersTable{})
	require.Len(t, filtered, 1)
	require.Equal(t, int32(42), filtered[0].GetUserId())
	require.Equal(t, "license:valid", filtered[0].GetIdentifier())
}

func TestCleanupUserPhoneNumbersSetsUserIDOnFallback(t *testing.T) {
	t.Parallel()

	syncer := &UsersSync{}
	user := &syncdata.DataUser{
		UserId:      11,
		PhoneNumber: new("555-0100"),
	}

	syncer.cleanupUserPhoneNumbers(user)

	require.Len(t, user.GetPhoneNumbers(), 1)
	assert.Equal(t, int32(11), user.GetPhoneNumbers()[0].GetUserId())
}

func TestCleanupUserPhoneNumbersSetsUserIDOnAllNumbers(t *testing.T) {
	t.Parallel()

	syncer := &UsersSync{}
	user := &syncdata.DataUser{
		UserId: 11,
		PhoneNumbers: []*users.PhoneNumber{
			{UserId: 99, Number: "555-0100", IsPrimary: true},
			{Number: "555-0101"},
		},
	}

	syncer.cleanupUserPhoneNumbers(user)

	require.Len(t, user.GetPhoneNumbers(), 2)
	for _, phoneNumber := range user.GetPhoneNumbers() {
		assert.Equal(t, int32(11), phoneNumber.GetUserId())
	}
}
