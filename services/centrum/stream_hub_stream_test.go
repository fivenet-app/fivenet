package centrum

import (
	"context"
	"testing"
	"time"

	centrumsettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/settings"
	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbcentrum "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/centrum"
	pkguserinfo "github.com/fivenet-app/fivenet/v2026/pkg/userinfo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

type testCentrumStreamServer struct {
	ctx context.Context
}

type authoritativeStreamUserInfo struct {
	pkguserinfo.UserInfoRetriever

	user *pbuserinfo.UserInfo
}

func (r *authoritativeStreamUserInfo) GetUserInfo(
	context.Context,
	int32,
) (*pbuserinfo.UserInfo, error) {
	return r.user, nil
}

func (s *testCentrumStreamServer) SetHeader(metadata.MD) error  { return nil }
func (s *testCentrumStreamServer) SendHeader(metadata.MD) error { return nil }
func (s *testCentrumStreamServer) SetTrailer(metadata.MD)       {}
func (s *testCentrumStreamServer) Context() context.Context     { return s.ctx }
func (s *testCentrumStreamServer) SendMsg(any) error            { return nil }
func (s *testCentrumStreamServer) RecvMsg(any) error            { return nil }
func (s *testCentrumStreamServer) Send(*pbcentrum.StreamResponse) error {
	return nil
}

func TestStreamRequestsSnapshotAfterPrimaryJobChange(t *testing.T) {
	t.Parallel()

	newJob := "police"
	newGrade := int32(3)
	changes := make(chan *pbuserinfo.UserInfoChanged, 1)
	changes <- &pbuserinfo.UserInfoChanged{
		UserId:      42,
		NewJob:      &newJob,
		NewJobGrade: &newGrade,
	}

	user := &pbuserinfo.UserInfo{UserId: 42, Job: "ambulance", JobGrade: 1}
	err := (&Server{}).stream(
		t.Context(),
		&testCentrumStreamServer{ctx: t.Context()},
		user,
		nil,
		make(chan *feedEvent),
		changes,
		0,
	)

	require.ErrorIs(t, err, errUserInfoChanged)
	assert.Equal(t, "ambulance", user.GetJob())
	assert.Equal(t, int32(1), user.GetJobGrade())
}

func TestRefreshStreamUserInfoUsesAuthoritativeState(t *testing.T) {
	t.Parallel()

	user := &pbuserinfo.UserInfo{UserId: 42, Job: "ambulance", JobGrade: 1}
	server := &Server{userinfo: &authoritativeStreamUserInfo{
		user: &pbuserinfo.UserInfo{UserId: 42, Job: "police", JobGrade: 3},
	}}

	require.NoError(t, server.refreshStreamUserInfo(t.Context(), user))
	assert.Equal(t, "police", user.GetJob())
	assert.Equal(t, int32(3), user.GetJobGrade())
}

func TestStreamRequestsSnapshotAfterOwnSettingsChange(t *testing.T) {
	t.Parallel()

	feed := make(chan *feedEvent, 1)
	feed <- &feedEvent{
		Sequence: 1,
		Jobs:     []string{"ambulance"},
		Response: &pbcentrum.StreamResponse{
			Change: &pbcentrum.StreamResponse_Settings{
				Settings: &centrumsettings.Settings{Job: "ambulance"},
			},
		},
	}

	err := (&Server{}).stream(
		t.Context(),
		&testCentrumStreamServer{ctx: t.Context()},
		&pbuserinfo.UserInfo{UserId: 42, Job: "ambulance", JobGrade: 1},
		nil,
		feed,
		make(chan *pbuserinfo.UserInfoChanged),
		0,
	)

	require.ErrorIs(t, err, errAccessChanged)
}

func TestStreamRequestsSnapshotAfterFeedWorkerRestart(t *testing.T) {
	t.Parallel()

	feed := make(chan *feedEvent, 1)
	feed <- &feedEvent{Sequence: 1, Resync: true}

	err := (&Server{}).stream(
		t.Context(),
		&testCentrumStreamServer{ctx: t.Context()},
		&pbuserinfo.UserInfo{UserId: 42, Job: "ambulance", JobGrade: 1},
		nil,
		feed,
		make(chan *pbuserinfo.UserInfoChanged),
		0,
	)

	require.ErrorIs(t, err, errFeedResync)
}

func TestStreamRequestsSnapshotAfterSequenceGap(t *testing.T) {
	t.Parallel()

	feed := make(chan *feedEvent, 2)
	feed <- &feedEvent{Sequence: 1}
	feed <- &feedEvent{Sequence: 3}

	err := (&Server{}).stream(
		t.Context(),
		&testCentrumStreamServer{ctx: t.Context()},
		&pbuserinfo.UserInfo{UserId: 42, Job: "ambulance", JobGrade: 1},
		nil,
		feed,
		make(chan *pbuserinfo.UserInfoChanged),
		0,
	)

	require.ErrorIs(t, err, errFeedResync)
}

func TestStreamIgnoresEventsCoveredBySnapshot(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(t.Context(), 20*time.Millisecond)
	defer cancel()
	feed := make(chan *feedEvent, 1)
	feed <- &feedEvent{Sequence: 1, Resync: true}

	err := (&Server{}).stream(
		ctx,
		&testCentrumStreamServer{ctx: ctx},
		&pbuserinfo.UserInfo{UserId: 42, Job: "ambulance", JobGrade: 1},
		nil,
		feed,
		make(chan *pbuserinfo.UserInfoChanged),
		1,
	)

	assert.NoError(t, err)
}

func TestStreamIgnoresEventsOutsideAuthorizedJobs(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(t.Context(), 20*time.Millisecond)
	defer cancel()
	feed := make(chan *feedEvent, 1)
	feed <- &feedEvent{
		Sequence: 1,
		Jobs:     []string{"police"},
		Response: &pbcentrum.StreamResponse{
			Change: &pbcentrum.StreamResponse_Settings{
				Settings: &centrumsettings.Settings{Job: "police"},
			},
		},
	}

	err := (&Server{}).stream(
		ctx,
		&testCentrumStreamServer{ctx: ctx},
		&pbuserinfo.UserInfo{UserId: 42, Job: "ambulance", JobGrade: 1},
		[]string{"ambulance"},
		feed,
		make(chan *pbuserinfo.UserInfoChanged),
		0,
	)

	assert.NoError(t, err)
}

func TestStreamRequestsSnapshotAfterOwnSettingsDeletion(t *testing.T) {
	t.Parallel()

	feed := make(chan *feedEvent, 1)
	feed <- &feedEvent{
		Sequence: 1,
		Response: &pbcentrum.StreamResponse{
			Change: &pbcentrum.StreamResponse_SettingsDeleted{
				SettingsDeleted: "ambulance",
			},
		},
	}

	err := (&Server{}).stream(
		t.Context(),
		&testCentrumStreamServer{ctx: t.Context()},
		&pbuserinfo.UserInfo{UserId: 42, Job: "ambulance", JobGrade: 1},
		nil,
		feed,
		make(chan *pbuserinfo.UserInfoChanged),
		0,
	)

	require.ErrorIs(t, err, errAccessChanged)
}
