package syncstore

import (
	"context"
	"testing"

	livemapmarkers "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/livemap/markers"
	properties "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	livemapstore "github.com/fivenet-app/fivenet/v2026/stores/livemap"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

type roundTripMarkerStore struct {
	marker      *livemapmarkers.MarkerMarker
	createCalls int
}

func TestSplitPropertyMutationsUsesManagedRegistry(t *testing.T) {
	t.Parallel()

	userValue := &properties.PropertyValue{}
	userValue.SetBoolValue(true)
	userMutation := &properties.PropertyMutation{Key: "user.wanted"}
	userMutation.SetSetValue(userValue)
	userProps, userChanges, err := splitUserPropertyMutations(&pbsync.SetUserPropertiesRequest{
		UserId:            42,
		PropertyMutations: []*properties.PropertyMutation{userMutation},
	})
	require.NoError(t, err)
	require.Empty(t, userChanges)
	require.NotNil(t, userProps)
	require.True(t, userProps.GetWanted())

	vehicleValue := &properties.PropertyValue{}
	vehicleValue.SetBoolValue(true)
	vehicleMutation := &properties.PropertyMutation{Key: "vehicle.wanted"}
	vehicleMutation.SetSetValue(vehicleValue)
	vehicleProps, vehicleChanges, hasManaged, err := splitVehiclePropertyMutations(
		&pbsync.SetVehiclePropertiesRequest{
			Plate:             "ABC123",
			PropertyMutations: []*properties.PropertyMutation{vehicleMutation},
		},
	)
	require.NoError(t, err)
	require.Empty(t, vehicleChanges)
	require.True(t, hasManaged)
	require.True(t, vehicleProps.GetWanted())
}

func TestWantedStateClearReason(t *testing.T) {
	t.Parallel()

	mutation := &properties.PropertyMutation{
		Key:         "vehicle.wanted",
		ClearReason: new(true),
	}
	require.True(t, wantedStateClearReason([]*properties.PropertyMutation{mutation}))
	require.False(t, wantedStateClearReason([]*properties.PropertyMutation{{
		Key:         "user.wanted",
		ClearReason: new(true),
	}}))
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
