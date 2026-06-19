package documentsstore

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	documentscategory "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/category"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreCategoriesReadWrite(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	store := New(db)
	userInfo := &userinfo.UserInfo{Job: "doj", Superuser: true}

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_documents_categories AS category`) + `(?s).*` + regexp.QuoteMeta(`ORDER BY category.sort_key ASC`)).
		WillReturnRows(sqlmock.NewRows([]string{"category.id", "category.name", "category.job"}).AddRow(int64(1), "Alpha", "doj"))

	categories, err := store.ListCategories(t.Context(), userInfo)
	require.NoError(t, err)
	assert.Len(t, categories, 1)
	assert.Equal(t, int64(1), categories[0].GetId())

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_documents_categories AS category`)+`(?s).*`+regexp.QuoteMeta(`category.id = ?`)+`(?s).*`+regexp.QuoteMeta(`LIMIT ?`)).
		WithArgs("doj", int64(1), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"category.id", "category.name", "category.job"}).AddRow(int64(1), "Alpha", "doj"))

	category, err := store.GetCategory(t.Context(), 1, userInfo)
	require.NoError(t, err)
	require.NotNil(t, category)
	assert.Equal(t, int64(1), category.GetId())

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO fivenet_documents_categories`)).
		WithArgs("Beta", nil, "doj", nil, nil).
		WillReturnResult(sqlmock.NewResult(2, 1))
	lastID, err := store.CreateCategory(
		t.Context(),
		db,
		&documentscategory.Category{Name: "Beta"},
		userInfo,
	)
	require.NoError(t, err)
	assert.Equal(t, int64(2), lastID)

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE fivenet_documents_categories SET`)).
		WithArgs("Beta", nil, "doj", nil, nil, int64(2), "doj", int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	require.NoError(
		t,
		store.UpdateCategory(
			t.Context(),
			db,
			&documentscategory.Category{Id: 2, Name: "Beta"},
			userInfo,
		),
	)

	deletedAt := timestamp.New(time.Date(2026, 6, 14, 12, 0, 0, 0, time.UTC))
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE fivenet_documents_categories SET`)).
		WithArgs(sqlmock.AnyArg(), "doj", int64(2), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	require.NoError(t, store.DeleteCategory(t.Context(), db, 2, userInfo, deletedAt))

	require.NoError(t, mock.ExpectationsWereMet())
}
