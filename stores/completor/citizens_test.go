package completorstore

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreCompleteCitizensAppliesSearchAndCustomFilter(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(
		db,
		&config.CustomDB{
			Conditions: dbutils.CustomConditions{
				User: dbutils.UserConditions{FilterEmptyName: true},
			},
		},
	)

	expectedQuery := regexp.QuoteMeta(`FROM fivenet_user AS user_short`) +
		`(?s).*` + regexp.QuoteMeta(`user_short.firstname != ?`) +
		`(?s).*` + regexp.QuoteMeta(`user_short.lastname != ?`) +
		`(?s).*` + regexp.QuoteMeta(`CONCAT(user_short.firstname, ?, user_short.lastname) LIKE ?`) +
		`(?s).*` + regexp.QuoteMeta(`ORDER BY user_short.lastname ASC LIMIT ?;`)

	mock.ExpectQuery(expectedQuery).
		WithArgs(true, "", "", " ", "%John%", int64(20)).
		WillReturnRows(sqlmock.NewRows([]string{
			"user_short.id",
			"user_short.firstname",
			"user_short.lastname",
			"user_short.dateofbirth",
		}).AddRow(int32(3), "John", "Doe", "1990-01-01"))

	users, err := store.CompleteCitizens(t.Context(), CitizensQuery{Search: "John"})
	require.NoError(t, err)
	require.Len(t, users, 1)
	assert.Equal(t, int32(3), users[0].GetUserId())
	assert.Equal(t, "John", users[0].GetFirstname())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreCompleteCitizensAddsCurrentJobColumnsAndOrdering(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(
		db,
		&config.CustomDB{
			Conditions: dbutils.CustomConditions{
				User: dbutils.UserConditions{FilterEmptyName: true},
			},
		},
	)

	expectedQuery := regexp.QuoteMeta(`FROM fivenet_user AS user_short`) +
		`(?s).*` + regexp.QuoteMeta(`user_short.firstname != ?`) +
		`(?s).*` + regexp.QuoteMeta(`user_short.lastname != ?`) +
		`(?s).*` + regexp.QuoteMeta(`user_short.job = ?`) +
		`(?s).*` + regexp.QuoteMeta(`OR (user_short.id IN (?, ?))`) +
		`(?s).*` + regexp.QuoteMeta(`ORDER BY user_short.id IN (?, ?) DESC, user_short.lastname ASC LIMIT ?;`)

	mock.ExpectQuery(expectedQuery).
		WithArgs(true, "", "", "police", int32(7), int32(9), int32(7), int32(9), int64(20)).
		WillReturnRows(sqlmock.NewRows([]string{
			"user_short.id",
			"user_short.firstname",
			"user_short.lastname",
			"user_short.dateofbirth",
			"user_short.job",
			"user_short.job_grade",
		}).AddRow(int32(7), "Jane", "Doe", "1992-02-02", "police", int32(12)))

	users, err := store.CompleteCitizens(t.Context(), CitizensQuery{
		CurrentJob:  true,
		UserJob:     "police",
		UserIDs:     []int32{7, 9},
		UserIDsOnly: true,
	})
	require.NoError(t, err)
	require.Len(t, users, 1)
	assert.Equal(t, "police", users[0].GetJob())
	assert.Equal(t, int32(12), users[0].GetJobGrade())
	require.NoError(t, mock.ExpectationsWereMet())
}
