package documentsstore

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	content "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/content"
	documentsactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/activity"
	documentscomment "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/comment"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreCommentCounters(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	store := New(db)

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_documents_comments AS comments`) + `(?s).*` + regexp.QuoteMeta(`comments.document_id = ?`) + `(?s).*` + regexp.QuoteMeta(`comments.deleted_at IS NULL`)).
		WithArgs(int64(42)).
		WillReturnRows(sqlmock.NewRows([]string{"comment_count"}).AddRow(int32(3)))

	count, err := store.CountComments(t.Context(), db, 42, false)
	require.NoError(t, err)
	assert.Equal(t, int32(3), count)

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_documents_comments AS comments`) + `(?s).*` + regexp.QuoteMeta(`comments.document_id = ?`) + `(?s).*` + regexp.QuoteMeta(`comments.deleted_at IS NULL`)).
		WithArgs(int64(42)).
		WillReturnRows(sqlmock.NewRows([]string{"comment_count"}).AddRow(int32(3)))

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO fivenet_documents_meta`)).
		WithArgs(int64(42), int32(3), int32(3)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	require.NoError(t, store.UpdateCommentsCount(t.Context(), db, 42))
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreCommentReads(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	store := New(db)

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_documents_comments AS comment`)+`(?s).*`+regexp.QuoteMeta(`comment.document_id = ?`)+`(?s).*`+regexp.QuoteMeta(`comment.deleted_at IS NULL`)+`(?s).*`+regexp.QuoteMeta(`ORDER BY comment.created_at DESC`)+`(?s).*`+regexp.QuoteMeta(`LIMIT ?`)).
		WithArgs(int64(42), int64(8), int64(0)).
		WillReturnRows(sqlmock.NewRows([]string{}))

	comments, err := store.ListComments(t.Context(), 42, &userinfo.UserInfo{Superuser: false}, 0, 8)
	require.NoError(t, err)
	assert.Empty(t, comments)

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_documents_comments AS comment`)+`(?s).*`+regexp.QuoteMeta(`comment.id = ?`)+`(?s).*`+regexp.QuoteMeta(`LIMIT ?`)).
		WithArgs(int64(99), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"comment.id", "comment.document_id", "comment.creator_id"}).AddRow(int64(99), int64(42), int32(7)))

	comment, err := store.GetComment(t.Context(), 99, &userinfo.UserInfo{Superuser: true})
	require.NoError(t, err)
	require.NotNil(t, comment)
	assert.Equal(t, int64(99), comment.GetId())
	assert.Equal(t, int64(42), comment.GetDocumentId())
	assert.Equal(t, int32(7), comment.GetCreatorId())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreCommentWrites(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	store := New(db)
	comment := &documentscomment.Comment{
		Id:         7,
		DocumentId: 42,
		Content:    content.Empty(),
	}
	userInfo := &userinfo.UserInfo{UserId: 3, Job: "doj"}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO fivenet_documents_comments`)).
		WithArgs(int64(42), sqlmock.AnyArg(), int32(3), "doj").
		WillReturnResult(sqlmock.NewResult(7, 1))
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO fivenet_documents_activity`)).
		WithArgs(int64(42), documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_COMMENT_ADDED, sqlmock.AnyArg(), "doj", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(8, 1))
	lastID, err := store.CreateComment(t.Context(), db, comment, userInfo)
	require.NoError(t, err)
	assert.Equal(t, int64(7), lastID)

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE fivenet_documents_comments SET`)).
		WithArgs(sqlmock.AnyArg(), int64(7), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO fivenet_documents_activity`)).
		WithArgs(int64(42), documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_COMMENT_UPDATED, sqlmock.AnyArg(), "doj", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(9, 1))
	require.NoError(t, store.UpdateComment(t.Context(), db, comment, userInfo))

	deletedAt := timestamp.New(time.Date(2026, 6, 14, 12, 0, 0, 0, time.UTC))
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE fivenet_documents_comments SET`)).
		WithArgs(sqlmock.AnyArg(), int64(7), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO fivenet_documents_activity`)).
		WithArgs(int64(42), documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_COMMENT_DELETED, sqlmock.AnyArg(), "doj", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(10, 1))
	require.NoError(
		t,
		store.DeleteComment(
			t.Context(),
			db,
			comment,
			userInfo,
			deletedAt,
			documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_COMMENT_DELETED,
		),
	)

	require.NoError(t, mock.ExpectationsWereMet())
}
