package documentsstore

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreGetTemplateOrderInfo(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(testParams(db))

	expectedQuery := regexp.QuoteMeta(`FROM fivenet_documents_templates AS template_order_info`) +
		`(?s).*` + regexp.QuoteMeta(`template_order_info.id = ?`) +
		`(?s).*` + regexp.QuoteMeta(`LIMIT ?;`)

	mock.ExpectQuery(expectedQuery).
		WithArgs(int64(42), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"template_order_info.id",
			"template_order_info.creator_job",
			"template_order_info.sort_rank",
		}).AddRow(int64(42), "doj", "000000001000"))

	info, err := store.GetTemplateOrderInfo(t.Context(), db, 42)
	require.NoError(t, err)
	require.NotNil(t, info)
	assert.Equal(t, int64(42), info.ID)
	assert.Equal(t, "doj", info.CreatorJob)
	assert.Equal(t, "000000001000", info.SortRank)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreInsertTemplateGroupRankUsesGapWhenAvailable(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(testParams(db))

	expectedQuery := regexp.QuoteMeta("FROM fivenet_documents_templates AS `row`") +
		`(?s).*` + regexp.QuoteMeta("`row`.creator_job = ?") +
		`(?s).*` + regexp.QuoteMeta("ORDER BY `row`.sort_rank ASC, `row`.id ASC") +
		`(?s).*` + regexp.QuoteMeta(`FOR UPDATE`)

	mock.ExpectQuery(expectedQuery).
		WithArgs("doj").
		WillReturnRows(sqlmock.NewRows([]string{"row.id", "row.sort_rank"}).
			AddRow(int64(1), "000000001000").
			AddRow(int64(2), "000000003000"))

	afterID := int64(1)
	rank, err := store.InsertTemplateGroupRank(t.Context(), db, "doj", 0, nil, &afterID)
	require.NoError(t, err)
	assert.Equal(t, "000000002000", rank)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreInsertTemplateGroupRankRebalancesWhenNoGap(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(testParams(db))

	rankQuery := regexp.QuoteMeta("FROM fivenet_documents_templates AS `row`") +
		`(?s).*` + regexp.QuoteMeta("`row`.creator_job = ?") +
		`(?s).*` + regexp.QuoteMeta("ORDER BY `row`.sort_rank ASC, `row`.id ASC") +
		`(?s).*` + regexp.QuoteMeta(`FOR UPDATE`)

	mock.ExpectQuery(rankQuery).
		WithArgs("doj").
		WillReturnRows(sqlmock.NewRows([]string{"row.id", "row.sort_rank"}).
			AddRow(int64(1), "000000001000").
			AddRow(int64(2), "000000001001"))
	mock.ExpectQuery(rankQuery).
		WithArgs("doj").
		WillReturnRows(sqlmock.NewRows([]string{"row.id", "row.sort_rank"}).
			AddRow(int64(1), "000000001000").
			AddRow(int64(2), "000000001001"))

	updateQuery := regexp.QuoteMeta(`UPDATE fivenet_documents_templates SET`) +
		`(?s).*` + regexp.QuoteMeta(`sort_rank = ?`) +
		`(?s).*` + regexp.QuoteMeta(`WHERE`) +
		`(?s).*` + regexp.QuoteMeta(`id = ?`) +
		`(?s).*` + regexp.QuoteMeta(`creator_job = ?`) +
		`(?s).*` + regexp.QuoteMeta(`deleted_at IS NULL`) +
		`(?s).*` + regexp.QuoteMeta(`LIMIT ?;`)

	mock.ExpectExec(updateQuery).
		WithArgs("000000001000", int64(1), "doj", int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(updateQuery).
		WithArgs("000000002000", int64(2), "doj", int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectQuery(rankQuery).
		WithArgs("doj").
		WillReturnRows(sqlmock.NewRows([]string{"row.id", "row.sort_rank"}).
			AddRow(int64(1), "000000001000").
			AddRow(int64(2), "000000002000"))

	afterID := int64(1)
	rank, err := store.InsertTemplateGroupRank(t.Context(), db, "doj", 0, nil, &afterID)
	require.NoError(t, err)
	assert.Equal(t, "000000001500", rank)
	require.NoError(t, mock.ExpectationsWereMet())
}
