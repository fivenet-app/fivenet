package statsstore

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreLoadPublicStats(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)

	expectations := []struct {
		query string
		value int64
	}{
		{query: `FROM fivenet_accounts`, value: 10},
		{query: `FROM fivenet_documents`, value: 20},
		{query: `FROM fivenet_centrum_dispatches`, value: 30},
		{query: `FROM fivenet_user_activity`, value: 40},
		{query: `FROM fivenet_job_timeclock`, value: 50},
		{query: `FROM fivenet_user`, value: 60},
	}

	for _, expectation := range expectations {
		mock.ExpectQuery(regexp.QuoteMeta(expectation.query)).
			WillReturnRows(sqlmock.NewRows([]string{"value"}).AddRow(expectation.value))
	}

	got, err := store.LoadPublicStats(t.Context())
	require.NoError(t, err)
	require.Len(t, got, len(expectations))
	assert.Equal(t, int32(10), got["users_registered"].GetValue())
	assert.Equal(t, int32(20), got["documents_created"].GetValue())
	assert.Equal(t, int32(30), got["dispatches_created"].GetValue())
	assert.Equal(t, int32(40), got["citizen_activity"].GetValue())
	assert.Equal(t, int32(50), got["timeclock_tracked"].GetValue())
	assert.Equal(t, int32(60), got["citizens_total"].GetValue())
	require.NoError(t, mock.ExpectationsWereMet())
}
