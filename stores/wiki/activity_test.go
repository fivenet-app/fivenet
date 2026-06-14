package wiki

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	wikiactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/wiki/activity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreCountPageActivity(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)

	expectedQuery := regexp.QuoteMeta(`FROM fivenet_wiki_pages_activity AS page_activity`) +
		`(?s).*` + regexp.QuoteMeta(`page_activity.page_id = ?`)

	mock.ExpectQuery(expectedQuery).
		WithArgs(int64(42)).
		WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(7)))

	total, err := store.CountPageActivity(t.Context(), 42)
	require.NoError(t, err)
	assert.Equal(t, int64(7), total)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreListPageActivity(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)

	expectedQuery := regexp.QuoteMeta(`FROM fivenet_wiki_pages_activity AS page_activity`) +
		`(?s).*` + regexp.QuoteMeta(`LEFT JOIN fivenet_user AS creator ON`) +
		`(?s).*` + regexp.QuoteMeta(`page_activity.page_id = ?`) +
		`(?s).*` + regexp.QuoteMeta(`LIMIT ? OFFSET ?;`)

	now := time.Unix(0, 0).UTC()
	mock.ExpectQuery(expectedQuery).
		WithArgs(int64(42), int64(25), int64(5)).
		WillReturnRows(sqlmock.NewRows([]string{
			"page_activity.id",
			"page_activity.created_at",
			"page_activity.page_id",
			"page_activity.activity_type",
			"page_activity.creator_id",
			"page_activity.creator_job",
			"page_activity.reason",
			"page_activity.data",
			"creator.id",
			"creator.job",
			"creator.job_grade",
			"creator.firstname",
			"creator.lastname",
		}).AddRow(
			int64(9),
			now,
			int64(42),
			int32(wikiactivity.PageActivityType_PAGE_ACTIVITY_TYPE_UPDATED),
			int32(7),
			"police",
			"reason",
			nil,
			int32(7),
			"police",
			int32(5),
			"Jane",
			"Doe",
		))

	activity, err := store.ListPageActivity(t.Context(), 42, 5, 25)
	require.NoError(t, err)
	require.Len(t, activity, 1)
	require.NotNil(t, activity[0].GetCreator())
	assert.Equal(t, "police", activity[0].GetCreator().GetJob())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreAddPageActivity(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)

	expectedQuery := regexp.QuoteMeta(`INSERT INTO fivenet_wiki_pages_activity`) +
		`(?s).*` + regexp.QuoteMeta(`VALUES (?, ?, ?, ?, ?, ?)`)

	mock.ExpectExec(expectedQuery).
		WithArgs(
			int64(42),
			int32(wikiactivity.PageActivityType_PAGE_ACTIVITY_TYPE_CREATED),
			int32(7),
			"police",
			nil,
			nil,
		).
		WillReturnResult(sqlmock.NewResult(9, 1))

	id, err := store.AddPageActivity(t.Context(), db, &wikiactivity.PageActivity{
		PageId:       42,
		ActivityType: wikiactivity.PageActivityType_PAGE_ACTIVITY_TYPE_CREATED,
		CreatorId:    func() *int32 { v := int32(7); return &v }(),
		CreatorJob:   "police",
	})
	require.NoError(t, err)
	assert.Equal(t, int64(9), id)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreCountPageChildren(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)

	expectedQuery := regexp.QuoteMeta(`FROM fivenet_wiki_pages`) +
		`(?s).*` + regexp.QuoteMeta(`parent_id = ?`) +
		`(?s).*` + regexp.QuoteMeta(`deleted_at IS NULL`) +
		`(?s).*` + regexp.QuoteMeta(`LIMIT ?;`)

	mock.ExpectQuery(expectedQuery).
		WithArgs(int64(42), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(2)))

	total, err := store.CountPageChildren(t.Context(), 42)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	require.NoError(t, mock.ExpectationsWereMet())
}
