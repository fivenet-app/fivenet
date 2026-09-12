package centrum

import (
	"context"
	"testing"

	centrumsettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/settings"
	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbcentrum "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/centrum"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

type testCentrumStreamServer struct {
	//nolint:containedctx // Used in test env/setup only.
	ctx context.Context
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
	assert.Equal(t, "police", user.GetJob())
	assert.Equal(t, int32(3), user.GetJobGrade())
}

func TestStreamRequestsSnapshotAfterOwnSettingsChange(t *testing.T) {
	t.Parallel()

	feed := make(chan *feedEvent, 1)
	feed <- &feedEvent{
		Sequence: 1,
		Job:      "ambulance",
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
