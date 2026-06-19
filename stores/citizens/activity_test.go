package citizensstore

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreCountUserActivityAppliesTargetFilter(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db, &config.CustomDB{})

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_user_activity AS user_activity`) + `(?s).*` + regexp.QuoteMeta(`user_activity.target_user_id = ?`)).
		WithArgs(int32(42)).
		WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(7)))

	total, err := store.CountUserActivity(t.Context(), CountUserActivityOptions{
		UserActivityOptions: UserActivityOptions{UserID: 42},
	})
	require.NoError(t, err)
	assert.Equal(t, int64(7), total)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreListUserActivityAppliesSortAndJoin(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db, &config.CustomDB{})

	expectedQuery := regexp.QuoteMeta(`FROM fivenet_user_activity AS user_activity`) +
		`(?s).*` + regexp.QuoteMeta(`INNER JOIN fivenet_user AS target_user ON`) +
		`(?s).*` + regexp.QuoteMeta(`LEFT JOIN fivenet_user AS source_user ON`) +
		`(?s).*` + regexp.QuoteMeta(`user_activity.target_user_id = ?`) +
		`(?s).*` + regexp.QuoteMeta(`ORDER BY user_activity.created_at DESC, user_activity.id DESC LIMIT ? OFFSET ?`)

	mock.ExpectQuery(expectedQuery).
		WithArgs(int32(42), int64(20), int64(0)).
		WillReturnRows(sqlmock.NewRows([]string{
			"user_activity.id",
			"user_activity.created_at",
			"user_activity.source_user_id",
			"user_activity.target_user_id",
			"user_activity.type",
			"user_activity.reason",
			"user_activity.data",
			"target_user.id",
			"target_user.job",
			"target_user.job_grade",
			"target_user.firstname",
			"target_user.lastname",
			"source_user.id",
			"source_user.job",
			"source_user.job_grade",
			"source_user.firstname",
			"source_user.lastname",
		}).AddRow(
			int64(1),
			time.Now(),
			nil,
			int32(42),
			int32(1),
			"updated",
			[]byte(`{}`),
			int32(42),
			"police",
			int32(2),
			"Jane",
			"Doe",
			nil,
			nil,
			nil,
			nil,
			nil,
		))

	activities, err := store.ListUserActivity(t.Context(), ListUserActivityOptions{
		UserActivityOptions: UserActivityOptions{UserID: 42},
		Offset:              0,
		Limit:               20,
	})
	require.NoError(t, err)
	require.Len(t, activities, 1)
	assert.Equal(t, int32(42), activities[0].GetTargetUser().GetUserId())
	require.NoError(t, mock.ExpectationsWereMet())
}
