package documentsstore

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	documentsactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/activity"
	resourcesdatabase "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	usershort "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/short"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreListAppliesFiltersAndSortFallback(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	store := New(db)

	closed := true
	onlyDrafts := false
	from := timestamp.New(time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC))
	to := timestamp.New(time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC))
	query := ListQuery{
		Search:      "fire",
		CategoryIDs: []int64{7, 8},
		CreatorIDs:  []int32{3, 4},
		From:        from,
		To:          to,
		Closed:      &closed,
		DocumentIDs: []int64{11, 12},
		OnlyDrafts:  &onlyDrafts,
		Sort: &resourcesdatabase.Sort{
			Columns: []*resourcesdatabase.SortByColumn{{Id: "unknown", Desc: true}},
		},
		Offset:             0,
		Limit:              20,
		IncludePhoneNumber: true,
		UserInfo: &userinfo.UserInfo{
			UserId:         3,
			Job:            "doj",
			JobGrade:       16,
			Superuser:      true,
			Enabled:        true,
			AccountId:      3,
			License:        "license",
			CanBeSuperuser: true,
		},
	}

	expectedQuery := `(?s).*`

	mock.ExpectQuery(expectedQuery).
		WithArgs(
			true,
			int32(3),
			"doj",
			true,
			"fire",
			int64(7), int64(8),
			int32(3), int32(4),
			from.AsTime(),
			to.AsTime(),
			true,
			int64(11), int64(12),
			false,
			int64(20),
			int64(0),
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	docs, err := store.List(t.Context(), query)
	require.NoError(t, err)
	assert.Empty(t, docs)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreGetIncludesContentAndPhoneNumber(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	store := New(db)
	query := GetQuery{
		DocumentID:         42,
		WithContent:        true,
		IncludePhoneNumber: true,
		UserInfo: &userinfo.UserInfo{
			UserId:    3,
			Job:       "doj",
			JobGrade:  16,
			Superuser: true,
		},
	}

	expectedQuery := regexp.QuoteMeta(`SELECT document.id`) +
		`(?s).*` + regexp.QuoteMeta(`document.data`) +
		`(?s).*` + regexp.QuoteMeta(`document.content_json`) +
		`(?s).*` + regexp.QuoteMeta(`creator.phone_number`) +
		`(?s).*` + regexp.QuoteMeta(`ORDER BY document.created_at DESC, document.updated_at DESC LIMIT ?`)

	mock.ExpectQuery(expectedQuery).
		WithArgs(int32(3), int64(42), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	doc, err := store.Get(t.Context(), query)
	require.NoError(t, err)
	assert.Nil(t, doc)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreGetDocumentMetaAndUpdateOwner(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	store := New(db)

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_documents_meta AS document_meta`) + `(?s).*` + regexp.QuoteMeta(`document_meta.document_id = ?`) + `(?s).*` + regexp.QuoteMeta(`LIMIT ?`)).
		WithArgs(int64(42), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{}))

	meta, err := store.GetDocumentMeta(t.Context(), db, 42)
	require.NoError(t, err)
	require.NotNil(t, meta)
	assert.Equal(t, int64(42), meta.GetDocumentId())

	userInfo := &userinfo.UserInfo{UserId: 3, Job: "doj"}
	newOwner := &usershort.UserShort{UserId: 9, Job: "new"}

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE fivenet_documents AS document SET`)).
		WithArgs(int32(9), int64(42), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO fivenet_documents_activity`)).
		WithArgs(int64(42), documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_OWNER_CHANGED, int32(3), "doj", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))
	require.NoError(t, store.UpdateDocumentOwner(t.Context(), db, 42, userInfo, newOwner))

	require.NoError(t, mock.ExpectationsWereMet())
}
