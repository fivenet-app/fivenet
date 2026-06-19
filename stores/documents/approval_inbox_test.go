package documentsstore

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	resourcesdatabase "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	documentsapproval "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/approval"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/stretchr/testify/require"
)

func TestStoreListApprovalTasksInboxUsesVisibilityForNonSuperuser(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	store := New(db)
	onlyDrafts := true
	query := ListApprovalTasksInboxQuery{
		Pagination: &resourcesdatabase.PaginationRequest{},
		UserInfo: &userinfo.UserInfo{
			UserId:   3,
			Job:      "doj",
			JobGrade: 16,
		},
		Statuses: []documentsapproval.ApprovalTaskStatus{
			documentsapproval.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING,
			documentsapproval.ApprovalTaskStatus_APPROVAL_TASK_STATUS_DECLINED,
		},
		NotAlreadyActed: true,
		OnlyDrafts:      &onlyDrafts,
	}

	countQuery := `(?s).*WITH user_subjects AS.*visible_sources AS.*winning_visibility AS.*` +
		regexp.QuoteMeta(`SELECT COUNT(approval_task.id) AS "data_count.total"`) +
		`.*` + regexp.QuoteMeta(`document_short.draft = ?`) +
		`.*` + regexp.QuoteMeta(`approval_task.assignee_kind = ?`) +
		`.*` + regexp.QuoteMeta(`approval_task.user_id = ?`) +
		`.*` + regexp.QuoteMeta(`approval_task.job = ?`) +
		`.*` + regexp.QuoteMeta(`approval_task.minimum_grade <= ?`) +
		`.*` + regexp.QuoteMeta(`NOT (EXISTS (`) +
		`.*` + regexp.QuoteMeta(`approval_task.status IN (?, ?)`)
	mock.ExpectQuery(countQuery).
		WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(1)))

	selectQuery := `(?s).*WITH user_subjects AS.*visible_sources AS.*winning_visibility AS.*` +
		regexp.QuoteMeta(`SELECT approval_task.id`) +
		`.*` + regexp.QuoteMeta(`document_short.draft = ?`) +
		`.*` + regexp.QuoteMeta(`approval_task.assignee_kind = ?`) +
		`.*` + regexp.QuoteMeta(`approval_task.user_id = ?`) +
		`.*` + regexp.QuoteMeta(`approval_task.job = ?`) +
		`.*` + regexp.QuoteMeta(`approval_task.minimum_grade <= ?`) +
		`.*` + regexp.QuoteMeta(`NOT (EXISTS (`) +
		`.*` + regexp.QuoteMeta(`approval_task.status IN (?, ?)`) +
		`.*` + regexp.QuoteMeta(`ORDER BY approval_task.created_at ASC LIMIT ? OFFSET ?;`)
	mock.ExpectQuery(selectQuery).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	count, tasks, err := store.ListApprovalTasksInbox(t.Context(), query)
	require.NoError(t, err)
	require.Equal(t, int64(1), count.Total)
	require.Empty(t, tasks)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreListApprovalTasksInboxSuperuserBypassesVisibilityCte(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	store := New(db)
	onlyDrafts := false
	query := ListApprovalTasksInboxQuery{
		Pagination: &resourcesdatabase.PaginationRequest{},
		UserInfo: &userinfo.UserInfo{
			UserId:    3,
			Job:       "doj",
			JobGrade:  16,
			Superuser: true,
		},
		OnlyDrafts: &onlyDrafts,
	}

	countQuery := regexp.QuoteMeta(`SELECT COUNT(approval_task.id) AS "data_count.total"`) +
		`(?s).*` + regexp.QuoteMeta(`FROM fivenet_documents_approval_tasks AS approval_task INNER JOIN fivenet_documents AS document_short ON (document_short.id = approval_task.document_id)`) +
		`(?s).*` + regexp.QuoteMeta(`approval_task.status IN (?)`) +
		`(?s).*` + regexp.QuoteMeta(`document_short.draft = ?`)
	mock.ExpectQuery(countQuery).
		WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(1)))

	selectQuery := regexp.QuoteMeta(`SELECT approval_task.id`) +
		`(?s).*` + regexp.QuoteMeta(`FROM fivenet_documents_approval_tasks AS approval_task INNER JOIN fivenet_documents AS document_short ON (document_short.id = approval_task.document_id)`) +
		`(?s).*` + regexp.QuoteMeta(`approval_task.status IN (?)`) +
		`(?s).*` + regexp.QuoteMeta(`document_short.draft = ?`) +
		`(?s).*` + regexp.QuoteMeta(`ORDER BY approval_task.created_at ASC LIMIT ? OFFSET ?;`)
	mock.ExpectQuery(selectQuery).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	count, tasks, err := store.ListApprovalTasksInbox(t.Context(), query)
	require.NoError(t, err)
	require.Equal(t, int64(1), count.Total)
	require.Empty(t, tasks)
	require.NoError(t, mock.ExpectationsWereMet())
}
