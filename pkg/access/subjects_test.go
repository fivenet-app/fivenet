package access

import (
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
		false,
		10,
		11,
	)

	sql, args := stmt.Sql()

	require.Contains(t, sql, "WITH actor_subjects AS")
	assert.Contains(t, sql, "fivenet_acl_subject_users")
	assert.Contains(t, sql, "fivenet_acl_subject_qualifications")
	assert.Contains(t, sql, "fivenet_acl_subject_job_grade_scopes")
	assert.Contains(t, sql, "fivenet_user_jobs")
	assert.Contains(t, sql, "fivenet_documents_visibility_public")
	assert.Contains(t, sql, "fivenet_documents_visibility_creator")
	assert.Contains(t, sql, "fivenet_documents_visibility_subject")
	assert.Contains(t, sql, "UNION")
	assert.Contains(t, sql, "SELECT DISTINCT doc_ids.id AS \"id\"")
	assert.NotContains(t, sql, "ROW_NUMBER() OVER")
	assert.NotContains(t, sql, "matching_acl")
	assert.NotContains(t, sql, "visible_objects")
	assert.NotContains(t, sql, "fivenet_documents_access")
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
		false,
		table.FivenetDocuments.ID.GT(mysql.Int(0)),
	)

	sql, _ := stmt.Sql()

	require.Contains(t, sql, "WITH actor_subjects AS")
	assert.Contains(t, sql, "COUNT(doc_ids.id) AS \"exact_total\"")
	assert.Contains(t, sql, "fivenet_documents_visibility_public")
	assert.Contains(t, sql, "fivenet_documents_visibility_creator")
	assert.Contains(t, sql, "fivenet_documents_visibility_subject")
	assert.NotContains(t, sql, "ROW_NUMBER() OVER")
	assert.Contains(t, sql, "fivenet_documents.deleted_at IS NULL")
}

func TestSubjectObjectAccessCountStatementIncludesDeletedForSuperuser(t *testing.T) {
	t.Parallel()

	access := NewDocumentsSubjectObjectAccess(nil)

	stmt := access.CountVisibleByConditionStatement(
		&userinfo.UserInfo{
			UserId:    1,
			Job:       "admin",
			JobGrade:  10,
			Superuser: true,
		},
		2,
		true,
		table.FivenetDocuments.ID.GT(mysql.Int(0)),
	)

	sql, _ := stmt.Sql()

	assert.NotContains(t, sql, "fivenet_documents.deleted_at IS NULL")
	assert.Contains(t, sql, "COUNT(visible_ids.id) AS \"exact_total\"")
}

func TestWikiPageSubjectObjectAccessVisibleIDsStatementShape(t *testing.T) {
	t.Parallel()

	access := NewWikiPageSubjectObjectAccess(nil)

	stmt := access.VisibleIDsStatement(
		&userinfo.UserInfo{
			UserId:   7,
			Job:      "police",
			JobGrade: 6,
		},
		2,
		false,
		10,
		11,
	)

	sql, args := stmt.Sql()

	assert.Contains(t, sql, "fivenet_wiki_pages_visibility_public")
	assert.Contains(t, sql, "fivenet_wiki_pages_visibility_creator")
	assert.Contains(t, sql, "fivenet_wiki_pages_visibility_subject")
	assert.NotEmpty(t, args)
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
