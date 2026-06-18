package access

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	resourcesaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/access"
	documentsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/access"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateDeniedAccessUsesSingleRow(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	access := NewDocumentsSubjectObjectAccess(db)

	mock.ExpectQuery("(?s)SELECT.*fivenet_acl_subjects.*").
		WithArgs(sqlmock.AnyArg(), int64(42), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(int32(1)))
	mock.ExpectExec("(?s)INSERT INTO .*fivenet_documents_access.*").
		WithArgs(int64(7), int64(42), int32(documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT), false).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = access.createDeniedAccess(
		t.Context(),
		db,
		7,
		42,
		SubjectAccessOptions{
			BlockedAccess: int32(documentsaccess.AccessLevel_ACCESS_LEVEL_BLOCKED),
			DeniedAccessLevels: []int32{
				int32(documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW),
				int32(documentsaccess.AccessLevel_ACCESS_LEVEL_COMMENT),
				int32(documentsaccess.AccessLevel_ACCESS_LEVEL_STATUS),
				int32(documentsaccess.AccessLevel_ACCESS_LEVEL_ACCESS),
				int32(documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT),
			},
		},
	)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSubjectAccessRowsToProto(t *testing.T) {
	t.Parallel()

	rows := []subjectAccessRow{
		{
			ID:              1,
			TargetID:        10,
			Access:          int32(documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW),
			Effect:          true,
			SubjectType:     int16(SubjectTypeJobGrade),
			ACLJob:          new("5net"),
			ACLMinimumGrade: new(int32(5)),
		},
		{
			ID:              2,
			TargetID:        10,
			Access:          int32(documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT),
			Effect:          true,
			SubjectType:     int16(SubjectTypeJobGrade),
			ACLJob:          new("5net"),
			ACLMinimumGrade: new(int32(6)),
		},
		{
			ID:              3,
			TargetID:        10,
			Access:          int32(documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT),
			Effect:          false,
			SubjectType:     int16(SubjectTypeJobGrade),
			ACLJob:          new("5net"),
			ACLMinimumGrade: new(int32(7)),
		},
		{
			ID:              4,
			TargetID:        10,
			Access:          int32(documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT),
			Effect:          true,
			SubjectType:     int16(SubjectTypeUser),
			SubjectUserID:   new(int32(42)),
			UserJob:         new("5net"),
			UserJobGrade:    new(int32(8)),
			UserFirstname:   new("Ada"),
			UserLastname:    new("Lovelace"),
			UserDateofbirth: new("1815-12-10"),
		},
	}

	out := subjectAccessRowsToProto(rows, SubjectAccessOptions{
		BlockedAccess: int32(documentsaccess.AccessLevel_ACCESS_LEVEL_BLOCKED),
		DeniedAccessLevels: []int32{
			int32(documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW),
			int32(documentsaccess.AccessLevel_ACCESS_LEVEL_COMMENT),
			int32(documentsaccess.AccessLevel_ACCESS_LEVEL_STATUS),
			int32(documentsaccess.AccessLevel_ACCESS_LEVEL_ACCESS),
			int32(documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT),
		},
	})

	require.Len(t, out.GetJobs(), 3)
	assert.Equal(t, int64(3), out.GetJobs()[0].GetId())
	assert.Equal(t, int32(7), out.GetJobs()[0].GetMinimumGrade())
	assert.Equal(
		t,
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_BLOCKED),
		out.GetJobs()[0].GetAccess(),
	)
	assert.Equal(
		t,
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW),
		out.GetJobs()[0].GetRequiredAccess(),
	)
	assert.Equal(t, int64(1), out.GetJobs()[1].GetId())
	assert.Equal(t, int32(5), out.GetJobs()[1].GetMinimumGrade())
	assert.Equal(
		t,
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW),
		out.GetJobs()[1].GetAccess(),
	)
	assert.Equal(
		t,
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW),
		out.GetJobs()[1].GetRequiredAccess(),
	)
	assert.Equal(t, int64(2), out.GetJobs()[2].GetId())
	assert.Equal(t, int32(6), out.GetJobs()[2].GetMinimumGrade())
	assert.Equal(
		t,
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT),
		out.GetJobs()[2].GetAccess(),
	)
	assert.Equal(
		t,
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT),
		out.GetJobs()[2].GetRequiredAccess(),
	)

	require.Len(t, out.GetUsers(), 1)
	assert.Equal(t, int32(42), out.GetUsers()[0].GetUserId())
	assert.Equal(
		t,
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT),
		out.GetUsers()[0].GetAccess(),
	)
	assert.Equal(
		t,
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT),
		out.GetUsers()[0].GetRequiredAccess(),
	)
	require.NotNil(t, out.GetUsers()[0].GetUser())
	assert.Equal(t, "Ada", out.GetUsers()[0].GetUser().GetFirstname())
	assert.Equal(t, "Lovelace", out.GetUsers()[0].GetUser().GetLastname())
}

func TestSubjectAccessRowsToProtoPreservesAllowAndDenyForSameSubject(t *testing.T) {
	t.Parallel()

	rows := []subjectAccessRow{
		{
			ID:              1,
			TargetID:        10,
			Access:          int32(documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW),
			Effect:          true,
			SubjectType:     int16(SubjectTypeJobGrade),
			ACLJob:          new("5net"),
			ACLMinimumGrade: new(int32(5)),
		},
		{
			ID:              2,
			TargetID:        10,
			Access:          int32(documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT),
			Effect:          false,
			SubjectType:     int16(SubjectTypeJobGrade),
			ACLJob:          new("5net"),
			ACLMinimumGrade: new(int32(5)),
		},
	}

	out := subjectAccessRowsToProto(rows, SubjectAccessOptions{
		BlockedAccess: int32(documentsaccess.AccessLevel_ACCESS_LEVEL_BLOCKED),
		DeniedAccessLevels: []int32{
			int32(documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW),
			int32(documentsaccess.AccessLevel_ACCESS_LEVEL_COMMENT),
			int32(documentsaccess.AccessLevel_ACCESS_LEVEL_STATUS),
			int32(documentsaccess.AccessLevel_ACCESS_LEVEL_ACCESS),
			int32(documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT),
		},
	})

	require.Len(t, out.GetJobs(), 2)
	assert.Equal(t, int64(2), out.GetJobs()[0].GetId())
	assert.Equal(
		t,
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_BLOCKED),
		out.GetJobs()[0].GetAccess(),
	)
	assert.Equal(
		t,
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW),
		out.GetJobs()[0].GetRequiredAccess(),
	)
	assert.Equal(t, int64(1), out.GetJobs()[1].GetId())
	assert.Equal(
		t,
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW),
		out.GetJobs()[1].GetAccess(),
	)
	assert.Equal(
		t,
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW),
		out.GetJobs()[1].GetRequiredAccess(),
	)
}

func TestCompareSubjectAccess(t *testing.T) {
	t.Parallel()

	current := subjectAccessRowsToProto([]subjectAccessRow{
		{
			ID:              1,
			TargetID:        10,
			Access:          int32(documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW),
			Effect:          true,
			SubjectType:     int16(SubjectTypeJobGrade),
			ACLJob:          new("5net"),
			ACLMinimumGrade: new(int32(5)),
		},
		{
			ID:            2,
			TargetID:      10,
			Access:        int32(documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT),
			Effect:        true,
			SubjectType:   int16(SubjectTypeUser),
			SubjectUserID: new(int32(42)),
		},
	}, SubjectAccessOptions{})
	in := subjectAccessRowsToProto([]subjectAccessRow{
		{
			ID:              1,
			TargetID:        10,
			Access:          int32(documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT),
			Effect:          true,
			SubjectType:     int16(SubjectTypeJobGrade),
			ACLJob:          new("5net"),
			ACLMinimumGrade: new(int32(5)),
		},
		{
			ID:            3,
			TargetID:      10,
			Access:        int32(documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW),
			Effect:        true,
			SubjectType:   int16(SubjectTypeUser),
			SubjectUserID: new(int32(99)),
		},
	}, SubjectAccessOptions{})

	changes := compareSubjectAccess(current, in)

	require.Len(t, changes.Jobs.ToUpdate, 1)
	assert.Equal(t, int64(1), changes.Jobs.ToUpdate[0].GetId())
	assert.Equal(
		t,
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT),
		changes.Jobs.ToUpdate[0].GetAccess(),
	)
	require.Empty(t, changes.Jobs.ToCreate)
	require.Empty(t, changes.Jobs.ToDelete)
	require.Len(t, changes.Users.ToCreate, 1)
	assert.Equal(t, int32(99), changes.Users.ToCreate[0].GetUserId())
	require.Len(t, changes.Users.ToDelete, 1)
	assert.Equal(t, int32(42), changes.Users.ToDelete[0].GetUserId())
	require.Empty(t, changes.Users.ToUpdate)
}

func TestCompareSubjectAccessKeepsAllowAndDenySeparate(t *testing.T) {
	t.Parallel()

	current := &resourcesaccess.Access{
		Jobs: []*resourcesaccess.JobAccess{{
			Id:           1,
			TargetId:     10,
			Job:          "5net",
			MinimumGrade: 5,
			Access:       int32(documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW),
			RequiredAccess: func() *int32 {
				v := int32(documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW)
				return &v
			}(),
		}},
	}
	in := &resourcesaccess.Access{
		Jobs: []*resourcesaccess.JobAccess{{
			Id:           2,
			TargetId:     10,
			Job:          "5net",
			MinimumGrade: 5,
			Access:       int32(documentsaccess.AccessLevel_ACCESS_LEVEL_BLOCKED),
			RequiredAccess: func() *int32 {
				v := int32(documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW)
				return &v
			}(),
		}},
	}

	changes := compareSubjectAccess(current, in)

	require.Len(t, changes.Jobs.ToCreate, 1)
	assert.Equal(t, int64(2), changes.Jobs.ToCreate[0].GetId())
	require.Len(t, changes.Jobs.ToDelete, 1)
	assert.Equal(t, int64(1), changes.Jobs.ToDelete[0].GetId())
	require.Empty(t, changes.Jobs.ToUpdate)
}
