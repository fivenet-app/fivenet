package documentsstore

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	resourcesdatabase "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	documentsapproval "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/approval"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbdocuments "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/documents"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreApprovalPolicyHelpers(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_documents_approval_policies AS approval_policy`)+`(?s).*`+regexp.QuoteMeta(`approval_policy.document_id = ?`)).
		WithArgs(int64(42), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{}))

	pol, err := store.GetApprovalPolicy(
		t.Context(),
		db,
		table.FivenetDocumentsApprovalPolicies.AS("approval_policy").DocumentID.EQ(mysql.Int64(42)),
	)
	require.NoError(t, err)
	require.Nil(t, pol)

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO fivenet_documents_approval_policies`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	require.NoError(
		t,
		store.CreateApprovalPolicy(t.Context(), db, 42, &documentsapproval.ApprovalPolicy{}),
	)

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_documents_approval_tasks AS approval_task`)+`(?s).*`+regexp.QuoteMeta(`approval_task.id = ?`)).
		WithArgs(int64(7), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{}))

	task, err := store.GetApprovalTask(t.Context(), db, 7)
	require.NoError(t, err)
	require.Nil(t, task)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreListApprovals(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_documents_approvals AS approval`) + `(?s).*` + regexp.QuoteMeta(`approval.document_id = ?`)).
		WithArgs(int64(42)).
		WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(0)))

	count, approvals, err := store.ListApprovals(
		t.Context(),
		ListApprovalsQuery{DocumentID: 42, Pagination: &resourcesdatabase.PaginationRequest{}},
	)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count.Total)
	assert.Empty(t, approvals)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreApprovalTaskWrites(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)
	snapDate := timestamp.New(time.Date(2026, 6, 14, 12, 0, 0, 0, time.UTC))
	userInfo := &userinfo.UserInfo{UserId: 3, Job: "doj"}
	label := "Reviewer"
	comment := "seed comment"
	seed := &pbdocuments.ApprovalTaskSeed{
		UserId:            9,
		Label:             &label,
		SignatureRequired: true,
		Comment:           &comment,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(fivenet_documents_approval_tasks.id) AS "C" FROM fivenet_documents_approval_tasks`)+`(?s).*`+regexp.QuoteMeta(`fivenet_documents_approval_tasks.document_id = ?`)+`(?s).*`+regexp.QuoteMeta(`fivenet_documents_approval_tasks.user_id = ?`)).
		WithArgs(int64(42), sqlmock.AnyArg(), int32(1), int32(9), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"C"}).AddRow(int32(0)))

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO fivenet_documents_approval_tasks`)).
		WithArgs(int64(42), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "Reviewer", true, sqlmock.AnyArg(), sqlmock.AnyArg(), "seed comment", int32(3), "doj").
		WillReturnResult(sqlmock.NewResult(7, 1))

	created, ensured, err := store.CreateApprovalTasks(
		t.Context(),
		db,
		userInfo,
		42,
		snapDate,
		[]*pbdocuments.ApprovalTaskSeed{seed},
	)
	require.NoError(t, err)
	assert.Equal(t, int32(1), created)
	assert.Equal(t, int32(0), ensured)

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM fivenet_documents_approval_tasks`)).
		WithArgs(int64(42), sqlmock.AnyArg(), int64(7), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	require.NoError(
		t,
		store.DeleteApprovalTasks(t.Context(), db, 42, snapDate, false, []int64{7}, 0),
	)

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE fivenet_documents_approval_tasks SET`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), int64(250)).
		WillReturnResult(sqlmock.NewResult(0, 2))
	affected, err := store.ExpireApprovalTasks(t.Context(), db)
	require.NoError(t, err)
	assert.Equal(t, int64(2), affected)

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE fivenet_documents_approvals SET`)).
		WithArgs(sqlmock.AnyArg(), int64(42)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE fivenet_documents_approval_tasks SET`)).
		WithArgs(sqlmock.AnyArg(), int64(42)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	require.NoError(t, store.ResetApprovalProgress(t.Context(), db, 42))

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreRecomputeApprovalPolicyTx(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)
	snapDate := timestamp.New(time.Date(2026, 6, 14, 12, 0, 0, 0, time.UTC))

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_documents_approval_policies AS approval_policy`)).
		WillReturnRows(sqlmock.NewRows([]string{}))
	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_documents`)).
		WillReturnRows(sqlmock.NewRows([]string{"creator_id"}).AddRow(int32(0)))
	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_documents_approvals`)).
		WillReturnRows(sqlmock.NewRows([]string{"approved", "declined"}).AddRow(int32(2), int32(1)))
	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_documents_approval_tasks`)).
		WillReturnRows(sqlmock.NewRows([]string{"total", "pending"}).AddRow(int32(3), int32(1)))
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO fivenet_documents_meta`)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE fivenet_documents_approval_policies SET assigned_count = ?, approved_count = ?, declined_count = ?, pending_count = ?, any_declined = ? WHERE fivenet_documents_approval_policies.document_id = ? LIMIT ?;`)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(t, store.RecomputeApprovalPolicyTx(t.Context(), db, 42, snapDate))
	require.NoError(t, mock.ExpectationsWereMet())
}
