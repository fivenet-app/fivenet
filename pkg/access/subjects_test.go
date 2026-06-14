package access

import (
	"strings"
	"testing"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubjectObjectAccessVisibleIDsStatementShape(t *testing.T) {
	t.Parallel()

	access := NewDocumentsSubjectObjectAccess(nil)

	stmt := access.VisibleIDsStatement(
		&userinfo.UserInfo{
			UserId:   7,
			Job:      "police",
			JobGrade: 6,
		},
		2,
		10,
		11,
	)

	sql, args := stmt.Sql()
	compactSQL := strings.Join(strings.Fields(sql), " ")

	require.Contains(t, sql, "WITH actor_subjects AS")
	assert.Contains(t, sql, "fivenet_acl_subject_users")
	assert.Contains(t, sql, "fivenet_acl_subject_qualifications")
	assert.Contains(t, sql, "fivenet_acl_subject_job_grade_scopes")
	assert.Contains(t, sql, "fivenet_user_jobs")
	assert.Contains(t, sql, "fivenet_documents_access")
	assert.Contains(t, sql, "ROW_NUMBER() OVER")
	assert.Contains(t, sql, "PARTITION BY fivenet_documents_access.target_id")
	assert.Contains(
		t,
		compactSQL,
		"fivenet_documents_access.effect IS TRUE) AND (fivenet_documents_access.access >= ?",
	)
	assert.Contains(
		t,
		compactSQL,
		"fivenet_documents_access.effect IS FALSE) AND (fivenet_documents_access.access = ?",
	)
	assert.Contains(t, sql, "HAVING")
	assert.Contains(t, sql, "fivenet_documents.public IS TRUE")
	require.NotEmpty(t, args)
}

func TestSubjectObjectAccessCountStatementShape(t *testing.T) {
	t.Parallel()

	access := NewDocumentsSubjectObjectAccess(nil)

	stmt := access.CountVisibleByConditionStatement(
		&userinfo.UserInfo{
			UserId:   7,
			Job:      "police",
			JobGrade: 6,
		},
		2,
		table.FivenetDocuments.ID.GT(mysql.Int(0)),
	)

	sql, _ := stmt.Sql()

	require.Contains(t, sql, "WITH actor_subjects AS")
	assert.Contains(t, sql, "COUNT(visible_objects.id) AS \"exact_total\"")
	assert.NotContains(t, sql, "ORDER BY visible_objects.id")
	assert.Contains(t, sql, "fivenet_documents.deleted_at IS NULL")
}

func TestSubjectConstants(t *testing.T) {
	t.Parallel()

	assert.Equal(t, SubjectTypeUser, SubjectType(1))
	assert.Equal(t, SubjectTypeQualification, SubjectType(2))
	assert.Equal(t, SubjectTypeJobGrade, SubjectType(3))
	assert.Equal(t, AccessEffectDeny, AccessEffect(0))
	assert.Equal(t, AccessEffectAllow, AccessEffect(1))
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
