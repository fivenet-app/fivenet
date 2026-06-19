package settingsstore

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/laws"
	pbsettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/settings"
	"github.com/stretchr/testify/require"
)

func TestStoreCreateOrUpdateLawBookCreatesLawBook(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT COALESCE(MAX(lawbook.sort_order), ?) AS "sort_order" FROM fivenet_lawbooks AS lawbook;`)).
		WillReturnRows(sqlmock.NewRows([]string{"sort_order"}).AddRow(int32(-1)))
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO fivenet_lawbooks`) + `(?s).*` + regexp.QuoteMeta(`ON DUPLICATE KEY UPDATE`)).
		WillReturnResult(sqlmock.NewResult(11, 1))
	mock.ExpectCommit()
	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_lawbooks AS law_book WHERE`) + `(?s).*` + regexp.QuoteMeta(`LIMIT ?;`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "sort_order", "deleted_at", "name", "description"}).AddRow(
			int64(11),
			time.Now(),
			time.Now(),
			int32(0),
			nil,
			"Traffic",
			nil,
		))

	lawBook, err := store.CreateOrUpdateLawBook(
		t.Context(),
		&pbsettings.CreateOrUpdateLawBookRequest{LawBook: &laws.LawBook{Name: "Traffic"}},
		false,
	)
	require.NoError(t, err)
	require.NotNil(t, lawBook)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreCreateOrUpdateLawCreatesLaw(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT lawbook.id AS "id" FROM fivenet_lawbooks AS lawbook WHERE`) + `(?s).*` + regexp.QuoteMeta(`LIMIT ?;`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(3)))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT COALESCE(MAX(law.sort_order), ?) AS "sort_order" FROM fivenet_lawbooks_laws AS law WHERE`) + `(?s).*` + regexp.QuoteMeta(`law.lawbook_id = ?;`)).
		WillReturnRows(sqlmock.NewRows([]string{"sort_order"}).AddRow(int32(-1)))
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO fivenet_lawbooks_laws`) + `(?s).*` + regexp.QuoteMeta(`ON DUPLICATE KEY UPDATE`)).
		WillReturnResult(sqlmock.NewResult(21, 1))
	mock.ExpectCommit()
	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_lawbooks_laws AS law WHERE`) + `(?s).*` + regexp.QuoteMeta(`LIMIT ?;`)).
		WillReturnRows(sqlmock.NewRows([]string{"law.id", "law.created_at", "law.updated_at", "law.deleted_at", "law.lawbook_id", "law.sort_order", "law.name", "law.description", "law.hint", "law.fine", "law.detention_time", "law.stvo_points"}).AddRow(
			int64(21),
			time.Now(),
			time.Now(),
			nil,
			int64(3),
			int32(0),
			"Speeding",
			nil,
			nil,
			nil,
			nil,
			nil,
		))

	law, refreshIDs, err := store.CreateOrUpdateLaw(
		t.Context(),
		&pbsettings.CreateOrUpdateLawRequest{Law: &laws.Law{LawbookId: 3, Name: "Speeding"}},
		false,
	)
	require.NoError(t, err)
	require.NotNil(t, law)
	require.Len(t, refreshIDs, 1)
	require.NoError(t, mock.ExpectationsWereMet())
}
