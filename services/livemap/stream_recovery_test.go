package livemap

import (
	"context"
	"io"
	"sync"
	"testing"
	"time"

	livemapmarkers "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/livemap/markers"
	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pblivemap "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/livemap"
	"github.com/fivenet-app/fivenet/v2026/internal/tests/nats"
	"github.com/fivenet-app/fivenet/v2026/internal/tests/permsstub"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2026/pkg/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/tracker"
	pkguserinfo "github.com/fivenet-app/fivenet/v2026/pkg/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/broker"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/puzpuzpuz/xsync/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

type streamRecoveryPerms struct{ permsstub.Permissions }

func (p *streamRecoveryPerms) AttrJobList(
	*pbuserinfo.UserInfo,
	perms.AttrRef[perms.JobListAttr],
) (*permissionsattributes.StringList, error) {
	return &permissionsattributes.StringList{}, nil
}

func (p *streamRecoveryPerms) AttrJobGradeList(
	*pbuserinfo.UserInfo,
	perms.AttrRef[perms.JobGradeListAttr],
) (*permissionsattributes.JobGradeList, error) {
	return &permissionsattributes.JobGradeList{
		Jobs: map[string]int32{"ambulance": 10},
	}, nil
}

type streamRecoveryTracker struct {
	*streamTestTracker

	marker *livemapmarkers.UserMarker
}

func (t *streamRecoveryTracker) GetUserMarkerById(id int32) (*livemapmarkers.UserMarker, bool) {
	if t.marker != nil && t.marker.GetUserId() == id {
		return t.marker, true
	}
	return nil, false
}

func (t *streamRecoveryTracker) ListTrackedJobs() []string {
	return []string{"ambulance"}
}

type streamRecoveryUserInfo struct {
	pkguserinfo.UserInfoRetriever

	userInfo *pbuserinfo.UserInfo
}

func (r *streamRecoveryUserInfo) GetUserInfo(context.Context, int32) (*pbuserinfo.UserInfo, error) {
	return proto.Clone(r.userInfo).(*pbuserinfo.UserInfo), nil
}

type streamRecoveryChanges struct {
	ch chan *pbuserinfo.UserInfoChanged
}

func (c *streamRecoveryChanges) SubscribeUserInfoChanges() chan *pbuserinfo.UserInfoChanged {
	return c.ch
}

func (*streamRecoveryChanges) UnsubscribeUserInfoChanges(chan *pbuserinfo.UserInfoChanged) {}

func (*streamRecoveryChanges) SubscribeAccountGroupsChanges() chan *pbuserinfo.AccountGroupsChanged {
	return make(chan *pbuserinfo.AccountGroupsChanged)
}

func (*streamRecoveryChanges) UnsubscribeAccountGroupsChanges(chan *pbuserinfo.AccountGroupsChanged) {
}

type streamRecoveryServer struct {
	pblivemap.LivemapService_StreamServer

	ctx context.Context
	mu  sync.Mutex
	ch  chan *pblivemap.StreamResponse
}

func (s *streamRecoveryServer) SetHeader(metadata.MD) error  { return nil }
func (s *streamRecoveryServer) SendHeader(metadata.MD) error { return nil }
func (s *streamRecoveryServer) SetTrailer(metadata.MD)       {}
func (s *streamRecoveryServer) Context() context.Context     { return s.ctx }
func (s *streamRecoveryServer) SendMsg(any) error            { return nil }
func (s *streamRecoveryServer) RecvMsg(any) error            { return io.EOF }

func (s *streamRecoveryServer) Send(resp *pblivemap.StreamResponse) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	select {
	case s.ch <- resp:
		return nil
	case <-s.ctx.Done():
		return s.ctx.Err()
	}
}

var (
	_ pblivemap.LivemapService_StreamServer = (*streamRecoveryServer)(nil)
	_ grpc.ServerStream                     = (*streamRecoveryServer)(nil)
)

func TestStreamRecoversAfterUserConsumerDeletion(t *testing.T) {
	t.Parallel()

	natsServer := nats.NewServer(t, nats.ServerOptions{InProcess: true})
	js := natsServer.GetJS()

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	grade := int32(3)
	user := newUserInfo(10, "police", grade, false)
	marker := &livemapmarkers.UserMarker{UserId: 42, Job: "ambulance", JobGrade: &grade, X: 1, Y: 2}
	trackerStub := &streamRecoveryTracker{
		streamTestTracker: &streamTestTracker{markers: []*livemapmarkers.UserMarker{marker}},
		marker:            &livemapmarkers.UserMarker{UserId: 10, Job: "police", JobGrade: &grade},
	}

	kv, err := js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket:  tracker.BucketUserLoc,
		History: 1,
		Storage: jetstream.MemoryStorage,
	})
	require.NoError(t, err)
	data, err := proto.Marshal(marker)
	require.NoError(t, err)
	_, err = kv.Put(ctx, "ambulance.3.42", data)
	require.NoError(t, err)

	streamCtx := auth.ContextWithUserInfo(ctx, user)
	responses := make(chan *pblivemap.StreamResponse, 32)
	streamSrv := &streamRecoveryServer{ctx: streamCtx, ch: responses}
	server := &Server{
		logger:              zap.NewNop(),
		js:                  js,
		perms:               &streamRecoveryPerms{},
		enricher:            mstlystcdata.NewDummyEnricher(),
		tracker:             trackerStub,
		userinfo:            &streamRecoveryUserInfo{userInfo: user},
		userinfoChanges:     &streamRecoveryChanges{ch: make(chan *pbuserinfo.UserInfoChanged)},
		markersCache:        xsync.NewMap[string, []*livemapmarkers.MarkerMarker](),
		markersDeletedCache: xsync.NewMap[string, []int64](),
		markersPublicCache:  newMarkerPublicCache(),
		broker:              broker.NewWithResyncOnSlowSubscriber[*brokerEvent](32),
	}
	go server.broker.Start(ctx)

	done := make(chan error, 1)
	go func() { done <- server.Stream(&pblivemap.StreamRequest{}, streamSrv) }()

	// The initial ACL, marker snapshot, and user snapshot establish the stream.
	for range 3 {
		select {
		case <-responses:
		case <-time.After(time.Second):
			t.Fatal("timed out waiting for initial stream response")
		}
	}

	var kvStream jetstream.Stream
	require.Eventually(t, func() bool {
		kvStream, err = js.Stream(ctx, "KV_"+tracker.BucketUserLoc)
		if err != nil {
			return false
		}
		for name := range kvStream.ConsumerNames(ctx).Name() {
			if err := kvStream.DeleteConsumer(ctx, name); err == nil {
				return true
			}
		}
		return false
	}, time.Second, 10*time.Millisecond)

	// Recreating the stream produces the same ACL and snapshots without ending the RPC.
	for range 3 {
		select {
		case <-responses:
		case <-time.After(time.Second):
			t.Fatal("timed out waiting for stream resync response")
		}
	}

	marker.SetX(9)
	data, err = proto.Marshal(marker)
	require.NoError(t, err)
	_, err = kv.Put(ctx, "ambulance.3.42", data)
	require.NoError(t, err)

	select {
	case resp := <-responses:
		require.InEpsilon(t, float64(9), resp.GetUserUpdates().GetUpdates()[0].GetX(), 0.0001)

	case <-time.After(time.Second):
		t.Fatal("timed out waiting for user marker update after consumer recovery")
	}

	cancel()
	require.NoError(t, <-done)
}
