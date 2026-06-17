package mailerstore

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreListUserEmailsVisible(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)
	now := time.Unix(0, 0).UTC()

	expectedQuery := `(?s).*WITH user_subjects AS.*visible_sources AS.*winning_visibility AS.*email.deactivated IS FALSE.*` +
		`ORDER BY email.job ASC, email.label ASC`
	mock.ExpectQuery(expectedQuery).
		WillReturnRows(sqlmock.NewRows([]string{
			"email.id",
			"email.created_at",
			"email.updated_at",
			"email.deleted_at",
			"email.deactivated",
			"email.job",
			"email.user_id",
			"email.email",
			"email.email_changed",
			"email.label",
		}).AddRow(
			int64(22),
			now,
			now,
			nil,
			false,
			"police",
			int32(7),
			"visible@example.com",
			nil,
			"Visible",
		))

	pageSize := int64(10)
	emails, err := store.ListUserEmails(
		t.Context(),
		db,
		&userinfo.UserInfo{UserId: 7, Job: "police"},
		&database.PaginationRequest{PageSize: &pageSize},
		false,
		false,
	)
	require.NoError(t, err)
	require.Len(t, emails, 1)
	assert.Equal(t, int64(22), emails[0].GetId())
	assert.Equal(t, "visible@example.com", emails[0].GetEmail())
	require.NoError(t, mock.ExpectationsWereMet())
}
