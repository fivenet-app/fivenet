package access

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
