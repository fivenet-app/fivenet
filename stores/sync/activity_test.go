package syncstore

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	citizenslabels "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/citizens/labels"
	livemapmarkers "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/livemap/markers"
	activity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/sync/activity"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	usersprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/props"
	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	citizensstore "github.com/fivenet-app/fivenet/v2026/stores/citizens"
	livemapstore "github.com/fivenet-app/fivenet/v2026/stores/livemap"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

type roundTripMarkerStore struct {
	marker      *livemapmarkers.MarkerMarker
	createCalls int
}

func (s *roundTripMarkerStore) CreateMarker(
	_ context.Context,
	marker *livemapmarkers.MarkerMarker,
	_ *int32,
	_ string,
) (int64, error) {
	s.createCalls++
	s.marker = cloneRoundTripMarker(marker)
	if s.marker.GetId() == 0 {
		s.marker.SetId(77)
	}
	return s.marker.GetId(), nil
}

func (s *roundTripMarkerStore) UpdateMarker(
	_ context.Context,
	marker *livemapmarkers.MarkerMarker,
	_ string,
) error {
	s.marker = cloneRoundTripMarker(marker)
	return nil
}

func (s *roundTripMarkerStore) DeleteMarker(
	_ context.Context,
	_ int64,
	_ *timestamp.Timestamp,
) error {
	return nil
}

func (s *roundTripMarkerStore) GetMarker(
	_ context.Context,
	_ int64,
) (*livemapmarkers.MarkerMarker, error) {
	if s.marker == nil {
		return nil, qrm.ErrNoRows
	}
	return cloneRoundTripMarker(s.marker), nil
}

func (s *roundTripMarkerStore) ListActiveMarkers(
	_ context.Context,
) ([]*livemapmarkers.MarkerMarker, error) {
	if s.marker == nil {
		return []*livemapmarkers.MarkerMarker{}, nil
	}
	return []*livemapmarkers.MarkerMarker{cloneRoundTripMarker(s.marker)}, nil
}

func (s *roundTripMarkerStore) ListDeletedMarkers(
	_ context.Context,
) ([]*livemapmarkers.MarkerMarker, error) {
	return []*livemapmarkers.MarkerMarker{}, nil
}

var _ livemapstore.IStore = (*roundTripMarkerStore)(nil)

func cloneRoundTripMarker(marker *livemapmarkers.MarkerMarker) *livemapmarkers.MarkerMarker {
	return proto.Clone(marker).(*livemapmarkers.MarkerMarker)
}

func TestAddMarkerPreservesPublicFlag(t *testing.T) {
	t.Parallel()

	store := &roundTripMarkerStore{}
	srv := &Store{livemapStore: store}

	public := true
	marker := &livemapmarkers.MarkerMarker{}
	marker.SetName("Sperrzone")
	marker.SetJob("police")
	marker.SetJobLabel("Police")
	marker.SetX(10)
	marker.SetY(20)
	marker.SetPublic(public)

	resp, err := srv.AddMarker(t.Context(), &pbsync.AddMarkerRequest{
		Marker: marker,
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, 1, store.createCalls)
	require.NotNil(t, store.marker)
	require.True(t, store.marker.GetPublic())
	require.NotNil(t, store.marker.GetExpiresAt())
	require.True(t, store.marker.GetPublic())
}

func TestAddMarkerAllowsNilCreatorID(t *testing.T) {
	t.Parallel()

	store := &roundTripMarkerStore{}
	srv := &Store{livemapStore: store}

	marker := &livemapmarkers.MarkerMarker{}
	marker.SetName("Sperrzone")
	marker.SetJob("police")
	marker.SetJobLabel("Police")
	marker.SetX(10)
	marker.SetY(20)

	resp, err := srv.AddMarker(t.Context(), &pbsync.AddMarkerRequest{
		Marker: marker,
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, 1, store.createCalls)
	require.NotNil(t, store.marker)
	require.Nil(t, store.marker.CreatorId)
}

func TestAddUserPropsPreservesExistingLabelsWhenSyncPayloadHasLabels(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)
	store.citizensStore = citizensstore.New(citizensstore.Params{
		DB:           store.db,
		CustomDB:     &config.CustomDB{},
		LabelsAccess: access.NewCitizenLabelsSubjectObjectAccess(store.db),
	})

	mock.ExpectBegin()
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
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		))

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_user_labels INNER JOIN fivenet_user_labels_job AS label ON`)+`(?s).*`+regexp.QuoteMeta(`WHERE (fivenet_user_labels.user_id = ?)`)+`(?s).*`+regexp.QuoteMeta(`LIMIT ?;`)).
		WithArgs(int32(42), int64(25)).
		WillReturnRows(sqlmock.NewRows([]string{
			"label.id",
			"label.job",
			"label.name",
			"label.color",
			"label.icon",
			"label.settings",
			"label.expiresAt",
		}).AddRow(
			int64(1),
			"police",
			"Hidden",
			"#ffffff",
			nil,
			nil,
			nil,
		).AddRow(
			int64(2),
			"police",
			"Visible",
			"#ffffff",
			nil,
			nil,
			nil,
		))

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM fivenet_user_labels`)+`(?s).*`+regexp.QuoteMeta(`label_id IN (?)`)+`(?s).*`+regexp.QuoteMeta(`LIMIT ?;`)).
		WithArgs(int32(42), int64(1), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	incoming := &activity.UserProps{
		Props: &usersprops.UserProps{
			UserId: 42,
			Labels: &citizenslabels.Labels{
				List: []*citizenslabels.Label{{Id: 2}},
			},
		},
	}

	resp, err := store.AddUserProps(t.Context(), &pbsync.AddUserPropsRequest{UserProps: incoming})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NoError(t, mock.ExpectationsWereMet())
}
