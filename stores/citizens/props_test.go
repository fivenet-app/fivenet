package citizensstore

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	citizenslabels "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/citizens/labels"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	usersprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/props"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestStoreGetUserPropsLoadsScalarPropsOnly(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(testParams(db))

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_user_property_values`)).
		WithArgs(int32(42)).
		WillReturnRows(sqlmock.NewRows([]string{
			"fivenet_user_property_values.user_id",
			"fivenet_user_property_values.property_key",
			"fivenet_user_property_values.value_type",
			"fivenet_user_property_values.value_string",
			"fivenet_user_property_values.value_integer",
			"fivenet_user_property_values.value_number",
			"fivenet_user_property_values.value_bool",
			"fivenet_user_property_values.value_timestamp",
			"fivenet_user_property_values.value_date",
			"fivenet_user_property_values.updated_at",
		}).AddRow(
			int32(42),
			"user.wanted", userPropertyValueTypeBool, nil, nil, nil, true, nil, nil, time.Now(),
		).AddRow(
			int32(
				42,
			), "user.job", userPropertyValueTypeString, "police", nil, nil, nil, nil, nil, time.Now(),
		).
			AddRow(
				int32(
					42,
				), "user.job_grade", userPropertyValueTypeInteger, nil, int64(2), nil, nil, nil, nil, time.Now(),
			).
			AddRow(
				int32(
					42,
				), "user.traffic_infraction_points", userPropertyValueTypeInteger, nil, int64(7), nil, nil, nil, nil, time.Now(),
			).
			AddRow(
				int32(
					42,
				), "user.open_fines", userPropertyValueTypeInteger, nil, int64(12), nil, nil, nil, nil, time.Now(),
			),
		)
	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_user_property_state`)).
		WithArgs(int32(42)).
		WillReturnRows(sqlmock.NewRows([]string{
			"fivenet_user_property_state.user_id",
			"fivenet_user_property_state.property_key",
			"fivenet_user_property_state.reason",
			"fivenet_user_property_state.active_at",
			"fivenet_user_property_state.expires_at",
			"fivenet_user_property_state.updated_at",
		}).AddRow(
			int32(42), "user.wanted", nil,
			time.Date(2026, 7, 22, 10, 0, 0, 0, time.UTC),
			time.Date(2026, 7, 23, 10, 0, 0, 0, time.UTC),
			time.Now(),
		))
	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_user_files AS mugshot_file`)+`(?s).*`+regexp.QuoteMeta(`LEFT JOIN fivenet_files AS mugshot ON`)).
		WithArgs(int32(42), "mugshot", int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"mugshot_file_id", "mugshot.mugshot_file_id", "file_path"}).AddRow(int64(100), int64(100), "/files/mugshot.jpg"))

	props, err := store.GetUserProps(t.Context(), db, 42)
	require.NoError(t, err)
	assert.Equal(t, int32(42), props.GetUserId())
	assert.Equal(
		t,
		timestamp.New(time.Date(2026, 7, 22, 10, 0, 0, 0, time.UTC)),
		props.GetWantedAt(),
	)
	assert.Equal(
		t,
		timestamp.New(time.Date(2026, 7, 23, 10, 0, 0, 0, time.UTC)),
		props.GetWantedTill(),
	)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestNormalizeWantedChangeSetsStartAndClearsTimestamps(t *testing.T) {
	t.Parallel()

	wanted := true
	current := &usersprops.UserProps{Wanted: new(bool)}
	in := &usersprops.UserProps{Wanted: &wanted}

	before := time.Now().UTC()
	normalizeWantedChange(current, in)
	after := time.Now().UTC()

	require.NotNil(t, in.GetWantedAt())
	assert.False(t, in.GetWantedAt().AsTime().Before(before))
	assert.False(t, in.GetWantedAt().AsTime().After(after))

	wantedTill := timestamp.New(after.Add(time.Hour))
	current = &usersprops.UserProps{
		Wanted:     &wanted,
		WantedAt:   in.WantedAt,
		WantedTill: wantedTill,
	}
	in = &usersprops.UserProps{Wanted: &wanted}
	normalizeWantedChange(current, in)
	assert.Equal(t, current.GetWantedAt(), in.GetWantedAt())
	assert.Equal(t, wantedTill, in.GetWantedTill())

	wantedFalse := false
	in.Wanted = &wantedFalse
	normalizeWantedChange(current, in)
	assert.Nil(t, in.WantedAt)
	assert.Nil(t, in.WantedTill)
}

func TestMergeUserPropsLabelsPreservesHiddenLabels(t *testing.T) {
	t.Parallel()

	current := []*citizenslabels.Label{
		{Id: 1, Name: "hidden"},
		{Id: 2, Name: "visible"},
	}
	visible := []*citizenslabels.Label{
		{Id: 2, Name: "visible"},
	}
	requested := []*citizenslabels.Label{}

	merged := mergeUserPropsLabels(current, visible, requested)

	require.Len(t, merged, 1)
	assert.Equal(t, int64(1), merged[0].GetId())
}

func TestStoreHandleUserPropsChangesUpdatesWanted(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(testParams(db))
	zeroInt64 := int64(0)
	zeroUint32 := uint32(0)
	wanted := false
	wantedTrue := true
	x := &usersprops.UserProps{
		UserId:                  42,
		Wanted:                  &wanted,
		TrafficInfractionPoints: &zeroUint32,
		OpenFines:               &zeroInt64,
		MugshotFileId:           &zeroInt64,
	}
	in := proto.Clone(x).(*usersprops.UserProps)
	in.Wanted = &wantedTrue
	wantedTill := timestamp.New(time.Date(2026, 7, 24, 10, 0, 0, 0, time.UTC))
	in.SetWantedTill(wantedTill)

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO fivenet_user_props`) + `(?s).*` + regexp.QuoteMeta(`ON DUPLICATE KEY UPDATE`)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM fivenet_user_files`)).
		WithArgs(int32(42), "mugshot").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`(?s)INSERT INTO fivenet_user_property_values`).
		WithArgs(int32(42), "user.wanted", uint8(userPropertyValueTypeBool), nil, nil, nil, true, nil, nil, uint8(userPropertyValueTypeBool)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`(?s)INSERT INTO fivenet_user_property_state`).
		WillReturnResult(sqlmock.NewResult(0, 1))

	activities, err := store.HandleUserPropsChanges(t.Context(), db, x, in, nil, "manual")
	require.NoError(t, err)
	require.Len(t, activities, 1)
	assert.Equal(t, wantedTill, activities[0].GetData().GetWantedChange().GetWantedTill())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSyncUserMugshotFileRespectsTotalUserFileLimit(t *testing.T) {
	t.Parallel()

	t.Run("inserts when a slot remains", func(t *testing.T) {
		t.Parallel()
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		t.Cleanup(func() { _ = db.Close() })

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM fivenet_user_files`)).
			WithArgs(int32(42), "mugshot").
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectQuery(`SELECT COUNT\(fivenet_user_files\.file_id\).*FROM fivenet_user_files.*WHERE fivenet_user_files\.user_id = \?;`).
			WithArgs(int32(42)).
			WillReturnRows(sqlmock.NewRows([]string{"total"}).AddRow(MaxUserFiles - 1))
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO fivenet_user_files`)).
			WithArgs(int32(42), int64(99), "mugshot").
			WillReturnResult(sqlmock.NewResult(0, 1))

		require.NoError(t, syncUserMugshotFile(t.Context(), db, 42, 99))
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rejects when all slots are occupied", func(t *testing.T) {
		t.Parallel()
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		t.Cleanup(func() { _ = db.Close() })

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM fivenet_user_files`)).
			WithArgs(int32(42), "mugshot").
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectQuery(`SELECT COUNT\(fivenet_user_files\.file_id\).*FROM fivenet_user_files.*WHERE fivenet_user_files\.user_id = \?;`).
			WithArgs(int32(42)).
			WillReturnRows(sqlmock.NewRows([]string{"total"}).AddRow(MaxUserFiles))

		require.EqualError(
			t,
			syncUserMugshotFile(t.Context(), db, 42, 99),
			"user file limit reached",
		)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
