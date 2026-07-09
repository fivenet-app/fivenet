package livemap

import (
	"context"
	"testing"
	"time"

	livemapmarkers "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/livemap/markers"
	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	timestamp "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	usershort "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/short"
	pblivemap "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/livemap"
	"github.com/fivenet-app/fivenet/v2026/internal/tests/permsstub"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2026/pkg/perms"
	errorslivemap "github.com/fivenet-app/fivenet/v2026/services/livemap/errors"
	livemapstore "github.com/fivenet-app/fivenet/v2026/stores/livemap"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/puzpuzpuz/xsync/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type markerTestStore struct {
	markers      map[int64]*livemapmarkers.MarkerMarker
	nextID       int64
	createCalls  int
	updateCalls  int
	deleteCalls  int
	lastUpdate   *livemapmarkers.MarkerMarker
	lastDeleteAt *timestamp.Timestamp
}

func newMarkerTestStore(markers ...*livemapmarkers.MarkerMarker) *markerTestStore {
	store := &markerTestStore{
		markers: make(map[int64]*livemapmarkers.MarkerMarker, len(markers)),
		nextID:  100,
	}

	for _, marker := range markers {
		store.markers[marker.GetId()] = cloneMarker(marker)
	}

	return store
}

func (s *markerTestStore) CreateMarker(
	_ context.Context,
	marker *livemapmarkers.MarkerMarker,
	creatorID int32,
	job string,
) (int64, error) {
	s.createCalls++

	if marker.GetId() <= 0 {
		s.nextID++
		marker.SetId(s.nextID)
	}
	marker.SetJob(job)
	marker.SetCreatorId(creatorID)
	marker.SetCreator(newCreator(creatorID, job, 0))
	s.markers[marker.GetId()] = cloneMarker(marker)

	return marker.GetId(), nil
}

func (s *markerTestStore) UpdateMarker(
	_ context.Context,
	marker *livemapmarkers.MarkerMarker,
	job string,
) error {
	s.updateCalls++
	if existing, ok := s.markers[marker.GetId()]; ok {
		if marker.GetCreatorId() == 0 {
			marker.SetCreatorId(existing.GetCreatorId())
		}
		if marker.GetCreator() == nil && existing.GetCreator() != nil {
			marker.SetCreator(proto.Clone(existing.GetCreator()).(*usershort.UserShort))
		}
	}
	marker.SetJob(job)
	s.lastUpdate = cloneMarker(marker)
	s.markers[marker.GetId()] = cloneMarker(marker)

	return nil
}

func (s *markerTestStore) DeleteMarker(
	_ context.Context,
	id int64,
	deletedAt *timestamp.Timestamp,
) error {
	s.deleteCalls++
	s.lastDeleteAt = deletedAt
	if marker, ok := s.markers[id]; ok {
		marker.DeletedAt = deletedAt
	}

	return nil
}

func (s *markerTestStore) GetMarker(
	_ context.Context,
	id int64,
) (*livemapmarkers.MarkerMarker, error) {
	marker, ok := s.markers[id]
	if !ok {
		return nil, qrm.ErrNoRows
	}

	return cloneMarker(marker), nil
}

func (s *markerTestStore) ListActiveMarkers(
	_ context.Context,
) ([]*livemapmarkers.MarkerMarker, error) {
	out := make([]*livemapmarkers.MarkerMarker, 0, len(s.markers))
	for _, marker := range s.markers {
		if marker.GetDeletedAt() != nil {
			continue
		}
		out = append(out, cloneMarker(marker))
	}

	return out, nil
}

func (s *markerTestStore) ListDeletedMarkers(
	_ context.Context,
) ([]*livemapmarkers.MarkerMarker, error) {
	out := make([]*livemapmarkers.MarkerMarker, 0, len(s.markers))
	for _, marker := range s.markers {
		if marker.GetDeletedAt() == nil {
			continue
		}
		out = append(out, cloneMarker(marker))
	}

	return out, nil
}

var _ livemapstore.IStore = (*markerTestStore)(nil)

type testLivemapPerms struct {
	permsstub.Permissions
}

func (p *testLivemapPerms) AttrStringList(
	_ *pbuserinfo.UserInfo,
	_ perms.AttrRef[perms.StringListAttr],
) (*permissionsattributes.StringList, error) {
	return &permissionsattributes.StringList{
		Strings: []string{"Any"},
	}, nil
}

func newMarkerServer(store livemapstore.IStore, ps perms.Permissions) *Server {
	return &Server{
		logger:              zap.NewNop(),
		ps:                  ps,
		enricher:            mstlystcdata.NewDummyEnricher(),
		store:               store,
		markersCache:        xsync.NewMap[string, []*livemapmarkers.MarkerMarker](),
		markersDeletedCache: xsync.NewMap[string, []int64](),
		markersPublicCache:  newMarkerPublicCache(),
	}
}

func newUserInfo(userID int32, job string, grade int32, admin bool) *pbuserinfo.UserInfo {
	userInfo := &pbuserinfo.UserInfo{}
	userInfo.SetUserId(userID)
	userInfo.SetJob(job)
	userInfo.SetJobGrade(grade)
	userInfo.SetJobAdmin(admin)
	return userInfo
}

func newCreator(userID int32, job string, grade int32) *usershort.UserShort {
	creator := &usershort.UserShort{}
	creator.SetUserId(userID)
	creator.SetJob(job)
	creator.SetJobGrade(grade)
	return creator
}

func newMarkerRequest(id int64, public *bool) *livemapmarkers.MarkerMarker {
	marker := &livemapmarkers.MarkerMarker{}
	marker.SetId(id)
	marker.SetName("Marker")
	marker.SetPostal("12345")
	marker.SetX(1.5)
	marker.SetY(2.5)
	marker.SetJob("police")
	if public != nil {
		marker.SetPublic(*public)
	}
	return marker
}

func cloneMarker(marker *livemapmarkers.MarkerMarker) *livemapmarkers.MarkerMarker {
	return proto.Clone(marker).(*livemapmarkers.MarkerMarker)
}

func TestCreateOrUpdateMarkerPublicAccess(t *testing.T) {
	t.Parallel()

	t.Run("non-admin create cannot set public", func(t *testing.T) {
		t.Parallel()

		store := newMarkerTestStore()
		srv := newMarkerServer(store, &testLivemapPerms{})
		ctx := auth.ContextWithUserInfo(t.Context(), newUserInfo(10, "police", 3, false))

		public := true
		_, err := srv.CreateOrUpdateMarker(ctx, &pblivemap.CreateOrUpdateMarkerRequest{
			Marker: newMarkerRequest(0, &public),
		})

		require.ErrorIs(t, err, errorslivemap.ErrMarkerDenied)
		require.Zero(t, store.createCalls)
	})

	t.Run("admin create can set public", func(t *testing.T) {
		t.Parallel()

		store := newMarkerTestStore()
		srv := newMarkerServer(store, &testLivemapPerms{})
		ctx := auth.ContextWithUserInfo(t.Context(), newUserInfo(10, "police", 3, true))

		public := true
		resp, err := srv.CreateOrUpdateMarker(ctx, &pblivemap.CreateOrUpdateMarkerRequest{
			Marker: newMarkerRequest(0, &public),
		})

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, resp.Marker)
		require.True(t, resp.Marker.GetPublic())
		require.Equal(t, 1, store.createCalls)
	})

	t.Run("omitted public is preserved on edit", func(t *testing.T) {
		t.Parallel()

		existing := newMarkerRequest(42, new(true))
		existing.SetCreatorId(10)
		existing.SetCreator(newCreator(10, "police", 3))
		existing.SetJob("police")
		existing.SetJobLabel("Police")

		store := newMarkerTestStore(existing)
		srv := newMarkerServer(store, &testLivemapPerms{})
		ctx := auth.ContextWithUserInfo(t.Context(), newUserInfo(10, "police", 3, false))

		req := newMarkerRequest(42, nil)
		req.SetName("Updated")

		resp, err := srv.CreateOrUpdateMarker(ctx, &pblivemap.CreateOrUpdateMarkerRequest{
			Marker: req,
		})

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, resp.Marker)
		require.NotNil(t, store.lastUpdate)
		require.True(t, store.lastUpdate.GetPublic())
		require.True(t, resp.Marker.GetPublic())
	})

	t.Run("non-admin cannot clear public", func(t *testing.T) {
		t.Parallel()

		existing := newMarkerRequest(42, new(true))
		existing.SetCreatorId(10)
		existing.SetCreator(newCreator(10, "police", 3))
		existing.SetJob("police")
		existing.SetJobLabel("Police")

		store := newMarkerTestStore(existing)
		srv := newMarkerServer(store, &testLivemapPerms{})
		ctx := auth.ContextWithUserInfo(t.Context(), newUserInfo(10, "police", 3, false))

		clearPublic := false
		_, err := srv.CreateOrUpdateMarker(ctx, &pblivemap.CreateOrUpdateMarkerRequest{
			Marker: newMarkerRequest(42, &clearPublic),
		})

		require.ErrorIs(t, err, errorslivemap.ErrMarkerDenied)
		require.Zero(t, store.updateCalls)
	})

	t.Run("admin can clear public", func(t *testing.T) {
		t.Parallel()

		existing := newMarkerRequest(42, new(true))
		existing.SetCreatorId(10)
		existing.SetCreator(newCreator(10, "police", 3))
		existing.SetJob("police")
		existing.SetJobLabel("Police")

		store := newMarkerTestStore(existing)
		srv := newMarkerServer(store, &testLivemapPerms{})
		ctx := auth.ContextWithUserInfo(t.Context(), newUserInfo(10, "police", 3, true))

		clearPublic := false
		resp, err := srv.CreateOrUpdateMarker(ctx, &pblivemap.CreateOrUpdateMarkerRequest{
			Marker: newMarkerRequest(42, &clearPublic),
		})

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, resp.Marker)
		require.NotNil(t, store.lastUpdate)
		require.False(t, store.lastUpdate.GetPublic())
		require.False(t, resp.Marker.GetPublic())
	})
}

func TestMarkerCacheMutationKeepsPublicBuckets(t *testing.T) {
	store := newMarkerTestStore()
	srv := newMarkerServer(store, &testLivemapPerms{})
	ctx := auth.ContextWithUserInfo(t.Context(), newUserInfo(10, "police", 3, true))

	public := true
	resp, err := srv.CreateOrUpdateMarker(ctx, &pblivemap.CreateOrUpdateMarkerRequest{
		Marker: newMarkerRequest(0, &public),
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Marker)

	active, deleted := srv.markersPublicCache.Snapshot()
	require.Len(t, active, 1)
	require.Empty(t, deleted)
	require.Equal(t, resp.Marker.GetId(), active[0].GetId())
	require.Equal(t, "Marker", active[0].GetName())
	_, ok := srv.markersCache.Load("police")
	require.False(t, ok)
	_, ok = srv.markersDeletedCache.Load("police")
	require.False(t, ok)

	updatedMarker := newMarkerRequest(resp.Marker.GetId(), nil)
	updatedMarker.SetName("Updated Marker")
	updatedResp, err := srv.CreateOrUpdateMarker(ctx, &pblivemap.CreateOrUpdateMarkerRequest{
		Marker: updatedMarker,
	})
	require.NoError(t, err)
	require.NotNil(t, updatedResp)
	require.NotNil(t, updatedResp.Marker)

	active, deleted = srv.markersPublicCache.Snapshot()
	require.Len(t, active, 1)
	require.Empty(t, deleted)
	require.Equal(t, resp.Marker.GetId(), active[0].GetId())
	require.Equal(t, "Updated Marker", active[0].GetName())

	_, err = srv.DeleteMarker(ctx, &pblivemap.DeleteMarkerRequest{Id: resp.Marker.GetId()})
	require.NoError(t, err)

	active, deleted = srv.markersPublicCache.Snapshot()
	require.Empty(t, active)
	require.Len(t, deleted, 1)
	require.Equal(t, resp.Marker.GetId(), deleted[0])
	_, ok = srv.markersCache.Load("police")
	require.False(t, ok)
	_, ok = srv.markersDeletedCache.Load("police")
	require.False(t, ok)

	_, err = srv.DeleteMarker(ctx, &pblivemap.DeleteMarkerRequest{Id: resp.Marker.GetId()})
	require.NoError(t, err)

	active, deleted = srv.markersPublicCache.Snapshot()
	require.Len(t, active, 1)
	require.Empty(t, deleted)
	require.Equal(t, resp.Marker.GetId(), active[0].GetId())
	require.Equal(t, "Updated Marker", active[0].GetName())
}

func TestGetMarkerMarkersIncludesPublicMarkers(t *testing.T) {
	t.Parallel()

	publicMarker := newMarkerRequest(100, new(true))
	publicMarker.SetExpiresAt(timestamp.New(time.Now().Add(time.Hour)))

	privateMarker := newMarkerRequest(101, new(false))
	privateMarker.SetJob("police")
	privateMarker.SetJobLabel("Police")

	expiredPublicMarker := newMarkerRequest(102, new(true))
	expiredPublicMarker.SetExpiresAt(timestamp.New(time.Now().Add(-time.Hour)))

	srv := &Server{
		markersCache:        xsync.NewMap[string, []*livemapmarkers.MarkerMarker](),
		markersDeletedCache: xsync.NewMap[string, []int64](),
		markersPublicCache:  newMarkerPublicCache(),
	}
	srv.markersPublicCache.Replace(
		[]*livemapmarkers.MarkerMarker{publicMarker, expiredPublicMarker},
		[]int64{200},
	)
	srv.markersCache.Store("police", []*livemapmarkers.MarkerMarker{privateMarker})
	srv.markersDeletedCache.Store("police", []int64{300})

	updated, deleted := srv.getMarkerMarkers(
		&permissionsattributes.StringList{Strings: []string{"police"}},
	)

	require.Len(t, updated, 2)
	require.Contains(t, updated, privateMarker)
	require.Contains(t, updated, publicMarker)
	require.Contains(t, deleted, int64(300))
	require.Contains(t, deleted, int64(102))
	require.Contains(t, deleted, int64(200))
}

func TestRefreshMarkersRebuildsPublicAndJobCaches(t *testing.T) {
	store := newMarkerTestStore(
		func() *livemapmarkers.MarkerMarker {
			marker := newMarkerRequest(100, new(false))
			marker.SetJob("police")
			marker.SetJobLabel("Police")
			return marker
		}(),
		func() *livemapmarkers.MarkerMarker {
			marker := newMarkerRequest(101, new(true))
			marker.SetJob("ems")
			marker.SetJobLabel("EMS")
			return marker
		}(),
		func() *livemapmarkers.MarkerMarker {
			marker := newMarkerRequest(102, new(false))
			marker.SetJob("police")
			marker.SetJobLabel("Police")
			marker.SetDeletedAt(timestamp.New(time.Now()))
			return marker
		}(),
		func() *livemapmarkers.MarkerMarker {
			marker := newMarkerRequest(103, new(true))
			marker.SetJob("fire")
			marker.SetJobLabel("Fire")
			marker.SetDeletedAt(timestamp.New(time.Now()))
			return marker
		}(),
	)

	srv := newMarkerServer(store, &testLivemapPerms{})

	require.NoError(t, srv.refreshMarkers(t.Context()))

	privateActive, ok := srv.markersCache.Load("police")
	require.True(t, ok)
	require.Len(t, privateActive, 1)
	require.Equal(t, int64(100), privateActive[0].GetId())

	privateDeleted, ok := srv.markersDeletedCache.Load("police")
	require.True(t, ok)
	require.Equal(t, []int64{102}, privateDeleted)

	active, deleted := srv.markersPublicCache.Snapshot()
	require.Len(t, active, 1)
	require.Equal(t, int64(101), active[0].GetId())
	require.Len(t, deleted, 1)
	require.Equal(t, int64(103), deleted[0])
}
