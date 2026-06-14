package citizens

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	usersprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/props"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestStoreGetUserPropsLoadsPropsAndLabels(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db, config.CustomDB{})

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_user_props AS user_props`)+`(?s).*`+regexp.QuoteMeta(`LEFT JOIN fivenet_files AS mugshot ON`)).
		WithArgs(int32(42), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"user_props.user_id",
			"user_props.updated_at",
			"user_props.wanted",
			"user_props.job",
			"user_props.job_grade",
			"user_props.traffic_infraction_points",
			"user_props.traffic_infraction_points_updated_at",
			"user_props.open_fines",
			"user_props.avatar_file_id",
			"user_props.mugshot_file_id",
			"mugshot.mugshot_file_id",
			"file_path",
		}).AddRow(
			int32(42),
			nil,
			true,
			"police",
			int32(2),
			uint32(7),
			nil,
			int64(12),
			nil,
			int64(100),
			int64(100),
			"/files/mugshot.jpg",
		))

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_user_labels`)+`(?s).*`+regexp.QuoteMeta(`INNER JOIN fivenet_user_labels_job AS citizen_label ON`)).
		WithArgs(int32(42), int64(25)).
		WillReturnRows(sqlmock.NewRows([]string{"citizen_label.id", "citizen_label.job", "citizen_label.name", "citizen_label.color"}))

	props, err := store.GetUserProps(t.Context(), db, 42)
	require.NoError(t, err)
	assert.Equal(t, int32(42), props.GetUserId())
	require.NotNil(t, props.GetLabels())
	require.Len(t, props.GetLabels().GetList(), 0)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreHandleUserPropsChangesUpdatesWanted(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db, config.CustomDB{})
	zeroInt64 := int64(0)
	zeroUint32 := uint32(0)
	wanted := true
	wantedFalse := false
	x := &usersprops.UserProps{
		UserId:                  42,
		Wanted:                  &wanted,
		TrafficInfractionPoints: &zeroUint32,
		OpenFines:               &zeroInt64,
		MugshotFileId:           &zeroInt64,
	}
	in := proto.Clone(x).(*usersprops.UserProps)
	in.Wanted = &wantedFalse

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO fivenet_user_props`) + `(?s).*` + regexp.QuoteMeta(`ON DUPLICATE KEY UPDATE`)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	activities, err := store.HandleUserPropsChanges(t.Context(), db, x, in, nil, "manual")
	require.NoError(t, err)
	require.Len(t, activities, 1)
	require.NoError(t, mock.ExpectationsWereMet())
}
