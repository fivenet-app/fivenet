package auth

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreListCharacters(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	store := New(db, &config.CustomDB{})

	expectedQuery := regexp.QuoteMeta(`FROM fivenet_user AS user`) +
		`(?s).*` + regexp.QuoteMeta(`LEFT JOIN fivenet_user_props AS user_props ON`) +
		`(?s).*` + regexp.QuoteMeta(`LEFT JOIN fivenet_files AS profile_picture ON`) +
		`(?s).*` + regexp.QuoteMeta(`user.account_id = ?`) +
		`(?s).*` + regexp.QuoteMeta(`user.license = ?`) +
		`(?s).*` + regexp.QuoteMeta(`ORDER BY user.id LIMIT ?;`)

	mock.ExpectQuery(expectedQuery).
		WillReturnRows(sqlmock.NewRows([]string{
			"user.id",
			"user.user_id",
			"user.account_id",
			"user.identifier",
			"user.job",
			"user.job_grade",
			"user.firstname",
			"user.lastname",
			"user.dateofbirth",
			"user.sex",
			"user.height",
			"user.phone_number",
			"user.profile_picture_file_id",
			"user.profile_picture",
			"character.group",
			"user.visum",
			"user.playtime",
		}).AddRow(
			int64(11),
			int64(11),
			int64(3),
			"license-11",
			"police",
			int32(2),
			"John",
			"Doe",
			"1990-01-01",
			"m",
			float64(180),
			"555-0100",
			int64(22),
			"/avatar.png",
			"leo",
			int32(12),
			int32(45),
		))

	chars, err := store.ListCharacters(t.Context(), 3, "license-11")
	require.NoError(t, err)
	require.Len(t, chars, 1)
	assert.Equal(t, "leo", chars[0].GetGroup())
	require.NotNil(t, chars[0].GetChar())
	assert.Equal(t, int32(11), chars[0].GetChar().GetUserId())
	assert.Equal(t, "John", chars[0].GetChar().GetFirstname())
	require.NoError(t, mock.ExpectationsWereMet())
}
