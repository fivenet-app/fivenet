package access

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func requireInt32(t *testing.T, value int) int32 {
	t.Helper()

	converted, ok := utils.ToInt32Checked(value)
	require.Truef(t, ok, "value %d must fit in int32", value)
	return converted
}

func TestSubjectCleanupOrphanSubjectsStatementShape(t *testing.T) {
	t.Parallel()

	stmt := NewSubjectResolver(nil).cleanupOrphanSubjectsStmt()
	sql, args := stmt.Sql()

	assert.Contains(t, sql, "LIMIT ?")
	assert.Contains(t, args, subjectCleanupDeleteLimit)
}

func TestSubjectCleanupStaleJobGradeSubjectsStatementShape(t *testing.T) {
	t.Parallel()

	stmt := NewSubjectResolver(nil).cleanupStaleJobGradeSubjectsStmt()
	sql, args := stmt.Sql()

	assert.Contains(t, sql, "LIMIT ?")
	assert.Contains(t, args, subjectCleanupDeleteLimit)
}

func TestResolveActorSubjectsUsesOnlyEffectiveWorkContext(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	resolver := NewSubjectResolver(db)
	mock.ExpectQuery(`(?s)SELECT.*fivenet_acl_subject_job_grade_scopes.*`).
		WithArgs(
			SubjectSpecificityUser, requireInt32(t, -1), requireInt32(t, 7),
			SubjectSpecificityQualification, requireInt32(t, -1), requireInt32(t, 7),
			SubjectSpecificityJobGrade, "mechanic", requireInt32(t, 3),
		).
		WillReturnRows(sqlmock.NewRows([]string{
			"subject_id",
			"specificity",
			"grade_specificity",
		}))

	subjects, err := resolver.ResolveActorSubjects(t.Context(), db, &userinfo.UserInfo{
		UserId:   7,
		Job:      "mechanic",
		JobGrade: 3,
	})
	require.NoError(t, err)
	assert.Empty(t, subjects)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestResolveActorSubjectsCanIncludeAllUserJobs(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	resolver := NewSubjectResolver(db)
	mock.ExpectQuery(`(?s)SELECT.*fivenet_user_jobs AS uj_resolve.*`).
		WithArgs(
			SubjectSpecificityUser, requireInt32(t, -1), requireInt32(t, 7),
			SubjectSpecificityQualification, requireInt32(t, -1), requireInt32(t, 7),
			SubjectSpecificityJobGrade, "mechanic", requireInt32(t, 3),
			SubjectSpecificityJobGrade, requireInt32(t, 7),
		).
		WillReturnRows(sqlmock.NewRows([]string{
			"subject_id",
			"specificity",
			"grade_specificity",
		}))

	subjects, err := resolver.ResolveActorSubjectsWithOptions(t.Context(), db, &userinfo.UserInfo{
		UserId:   7,
		Job:      "mechanic",
		JobGrade: 3,
	}, ActorSubjectOptions{IncludeAllUserJobs: true})
	require.NoError(t, err)
	assert.Empty(t, subjects)
	require.NoError(t, mock.ExpectationsWereMet())
}
