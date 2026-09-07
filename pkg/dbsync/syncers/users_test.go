package syncers

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	syncdata "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/sync/data"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users"
	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	dbsyncconfig "github.com/fivenet-app/fivenet/v2026/pkg/dbsync/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type deleteUsersClient struct {
	pbsync.SyncServiceClient

	err error
	req *pbsync.DeleteUsersRequest
}

func (c *deleteUsersClient) DeleteUsers(
	_ context.Context,
	req *pbsync.DeleteUsersRequest,
	_ ...grpc.CallOption,
) (*pbsync.DeleteDataResponse, error) {
	c.req = req
	return &pbsync.DeleteDataResponse{}, c.err
}

func newSingleUserSync(
	t *testing.T,
	client pbsync.SyncServiceClient,
) (*UsersSync, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	query := "SELECT id, identifier FROM users WHERE id = ? LIMIT 1"
	syncer := &UsersSync{
		Syncer: &Syncer{
			logger: zap.NewNop(),
			db:     db,
			cli:    client,
			cfg: &dbsyncconfig.DBSyncConfig{
				Tables: dbsyncconfig.DBSyncSourceTables{
					Users: dbsyncconfig.UsersTable{
						DBSyncTable: dbsyncconfig.DBSyncTable{Query: &query},
					},
				},
			},
		},
		logger: zap.NewNop(),
	}
	return syncer, mock
}

func TestSyncUserDeletesDestinationWhenSourceUserIsMissing(t *testing.T) {
	t.Parallel()

	client := &deleteUsersClient{}
	syncer, mock := newSingleUserSync(t, client)
	mock.ExpectQuery("SELECT id, identifier FROM users WHERE id = \\? LIMIT 1").
		WithArgs(int32(42)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "identifier"}))

	require.NoError(t, syncer.SyncUser(t.Context(), 42))
	require.NotNil(t, client.req)
	require.Equal(t, []int32{42}, client.req.GetUserIds())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSyncUserRejectsInvalidUserID(t *testing.T) {
	t.Parallel()

	syncer, mock := newSingleUserSync(t, &deleteUsersClient{})
	for _, userID := range []int32{0, -1} {
		require.Error(t, syncer.SyncUser(t.Context(), userID))
	}
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSyncUserReturnsDeleteError(t *testing.T) {
	t.Parallel()

	deleteErr := errors.New("delete unavailable")
	client := &deleteUsersClient{err: deleteErr}
	syncer, mock := newSingleUserSync(t, client)
	mock.ExpectQuery("SELECT id, identifier FROM users WHERE id = \\? LIMIT 1").
		WithArgs(int32(42)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "identifier"}))

	err := syncer.SyncUser(t.Context(), 42)
	require.Error(t, err)
	require.ErrorIs(t, err, deleteErr)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSyncUserDryRunDoesNotCallDeleteRPC(t *testing.T) {
	t.Parallel()

	client := &deleteUsersClient{}
	syncer, mock := newSingleUserSync(t, client)
	syncer.cfg.Destination.DryRun = true
	mock.ExpectQuery("SELECT id, identifier FROM users WHERE id = \\? LIMIT 1").
		WithArgs(int32(42)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "identifier"}))

	require.NoError(t, syncer.SyncUser(t.Context(), 42))
	require.Nil(t, client.req)
	require.NoError(t, mock.ExpectationsWereMet())
}

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
