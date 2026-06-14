package wiki

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreGetPageOrderInfo(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)

	expectedQuery := regexp.QuoteMeta(`FROM fivenet_wiki_pages AS page_order_info`) +
		`(?s).*` + regexp.QuoteMeta(`page_order_info.id = ?`) +
		`(?s).*` + regexp.QuoteMeta(`LIMIT ?;`)

	mock.ExpectQuery(expectedQuery).
		WithArgs(int64(42), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"page_order_info.id",
			"page_order_info.job",
			"page_order_info.parent_id",
			"page_order_info.startpage",
			"page_order_info.sort_rank",
		}).AddRow(int64(42), "police", nil, false, "1000"))

	info, err := store.GetPageOrderInfo(t.Context(), db, 42)
	require.NoError(t, err)
	require.NotNil(t, info)
	assert.Equal(t, int64(42), info.ID)
	assert.Equal(t, "police", info.Job)
	assert.Nil(t, info.ParentID)
	assert.False(t, info.Startpage)
	assert.Equal(t, "1000", info.SortRank)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreInsertPageGroupRankUsesGapWhenAvailable(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)

	expectedQuery := regexp.QuoteMeta(`FROM fivenet_wiki_pages AS page_rank_row`) +
		`(?s).*` + regexp.QuoteMeta(`page_rank_row.job = ?`) +
		`(?s).*` + regexp.QuoteMeta(`page_rank_row.parent_id IS NULL`) +
		`(?s).*` + regexp.QuoteMeta(`page_rank_row.startpage = ?`) +
		`(?s).*` + regexp.QuoteMeta(`ORDER BY page_rank_row.sort_rank ASC, page_rank_row.id ASC`) +
		`(?s).*` + regexp.QuoteMeta(`FOR UPDATE`)

	mock.ExpectQuery(expectedQuery).
		WithArgs("police", false).
		WillReturnRows(sqlmock.NewRows([]string{"page_rank_row.id", "page_rank_row.sort_rank"}).
			AddRow(int64(1), "1000").
			AddRow(int64(2), "3000"))

	afterID := int64(1)
	rank, err := store.InsertPageGroupRank(t.Context(), db, "police", nil, false, 0, nil, &afterID)
	require.NoError(t, err)
	assert.Equal(t, utils.FormatRank(2000), rank)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreInsertPageGroupRankRebalancesWhenNoGap(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)

	rankQuery := regexp.QuoteMeta(`FROM fivenet_wiki_pages AS page_rank_row`) +
		`(?s).*` + regexp.QuoteMeta(`page_rank_row.job = ?`) +
		`(?s).*` + regexp.QuoteMeta(`page_rank_row.parent_id IS NULL`) +
		`(?s).*` + regexp.QuoteMeta(`page_rank_row.startpage = ?`) +
		`(?s).*` + regexp.QuoteMeta(`ORDER BY page_rank_row.sort_rank ASC, page_rank_row.id ASC`) +
		`(?s).*` + regexp.QuoteMeta(`FOR UPDATE`)

	mock.ExpectQuery(rankQuery).
		WithArgs("police", false).
		WillReturnRows(sqlmock.NewRows([]string{"page_rank_row.id", "page_rank_row.sort_rank"}).
			AddRow(int64(1), "1000").
			AddRow(int64(2), "1001"))
	mock.ExpectQuery(rankQuery).
		WithArgs("police", false).
		WillReturnRows(sqlmock.NewRows([]string{"page_rank_row.id", "page_rank_row.sort_rank"}).
			AddRow(int64(1), "1000").
			AddRow(int64(2), "1001"))

	updateQuery := regexp.QuoteMeta(`UPDATE fivenet_wiki_pages SET`) +
		`(?s).*` + regexp.QuoteMeta(`sort_rank = ?`) +
		`(?s).*` + regexp.QuoteMeta(`WHERE`) +
		`(?s).*` + regexp.QuoteMeta(`id = ?`) +
		`(?s).*` + regexp.QuoteMeta(`job = ?`) +
		`(?s).*` + regexp.QuoteMeta(`deleted_at IS NULL`) +
		`(?s).*` + regexp.QuoteMeta(`LIMIT ?;`)

	mock.ExpectExec(updateQuery).
		WithArgs(utils.FormatRank(1000), int64(1), "police", int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(updateQuery).
		WithArgs(utils.FormatRank(2000), int64(2), "police", int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectQuery(rankQuery).
		WithArgs("police", false).
		WillReturnRows(sqlmock.NewRows([]string{"page_rank_row.id", "page_rank_row.sort_rank"}).
			AddRow(int64(1), "1000").
			AddRow(int64(2), "2000"))

	afterID := int64(1)
	rank, err := store.InsertPageGroupRank(t.Context(), db, "police", nil, false, 0, nil, &afterID)
	require.NoError(t, err)
	assert.Equal(t, utils.FormatRank(1500), rank)
	require.NoError(t, mock.ExpectationsWereMet())
}
