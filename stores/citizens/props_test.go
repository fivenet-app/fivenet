package citizensstore

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	citizenslabels "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/citizens/labels"
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

	props, err := store.GetUserProps(t.Context(), db, 42)
	require.NoError(t, err)
	assert.Equal(t, int32(42), props.GetUserId())
	require.NoError(t, mock.ExpectationsWereMet())
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
