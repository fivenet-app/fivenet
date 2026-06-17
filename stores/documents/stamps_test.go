package documentsstore

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreListUsableStampsUsesAclBranchesForNonSuperuser(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	store := New(db)

	mock.ExpectQuery(`(?s).*FROM fivenet_documents_stamps AS stamp.*stamp\.name = \?.*fivenet_documents_access.*subject_acl_user_exists.*subject_acl_qualification_exists.*subject_acl_job_grade_exists.*`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	pag, stamps, err := store.ListUsableStamps(t.Context(), ListUsableStampsQuery{
		UserInfo: &userinfo.UserInfo{UserId: 3, Job: "doj", JobGrade: 16},
	})
	require.NoError(t, err)
	assert.NotNil(t, pag)
	assert.Empty(t, stamps)
	require.NoError(t, mock.ExpectationsWereMet())
}
