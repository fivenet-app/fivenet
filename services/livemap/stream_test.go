package livemap

import (
	"context"
	"testing"
	"time"

	livemapmarkers "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/livemap/markers"
	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	resourcestracker "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/tracker"
	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pblivemap "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/livemap"
	"github.com/fivenet-app/fivenet/v2026/pkg/nats/store"
	pkgtracker "github.com/fivenet-app/fivenet/v2026/pkg/tracker"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

type streamTestTracker struct {
	markers []*livemapmarkers.UserMarker
}

func (t *streamTestTracker) ListTrackedJobs() []string {
	return nil
}

func (t *streamTestTracker) GetUserMarkerById(_ int32) (*livemapmarkers.UserMarker, bool) {
	return nil, false
}

func (t *streamTestTracker) IsUserOnDuty(_ int32) bool {
	return false
}

func (t *streamTestTracker) Subscribe(
	_ context.Context,
) (store.IKVWatcher[livemapmarkers.UserMarker, *livemapmarkers.UserMarker], error) {
	return nil, nil
}

func (t *streamTestTracker) GetFilteredUserMarkers(
	_ *permissionsattributes.JobGradeList,
	_ *pbuserinfo.UserInfo,
) []*livemapmarkers.UserMarker {
	return t.markers
}

func (t *streamTestTracker) GetUserMapping(_ int32) (*resourcestracker.UserMapping, error) {
	return nil, nil
}

func (t *streamTestTracker) SetUserMapping(
	_ context.Context,
	_ *resourcestracker.UserMapping,
) error {
	return nil
}

func (t *streamTestTracker) SetUserMappingForUser(
	_ context.Context,
	_ int32,
	_ *int64,
) error {
	return nil
}

func (t *streamTestTracker) UnsetUnitIDForUser(_ context.Context, _ int32) error {
	return nil
}

func (t *streamTestTracker) ListUserMappings(
	_ context.Context,
) (map[int32]*resourcestracker.UserMapping, error) {
	return nil, nil
}

var _ pkgtracker.ITracker = (*streamTestTracker)(nil)

type streamTestMsg struct {
	subject string
	headers nats.Header
	data    []byte
}

func (m *streamTestMsg) Metadata() (*jetstream.MsgMetadata, error) {
	return nil, nil
}

func (m *streamTestMsg) Data() []byte {
	return m.data
}

func (m *streamTestMsg) Headers() nats.Header {
	return m.headers
}

func (m *streamTestMsg) Subject() string {
	return m.subject
}

func (m *streamTestMsg) Reply() string {
	return ""
}

func (m *streamTestMsg) Ack() error {
	return nil
}

func (m *streamTestMsg) DoubleAck(context.Context) error {
	return nil
}

func (m *streamTestMsg) Nak() error {
	return nil
}

func (m *streamTestMsg) NakWithDelay(time.Duration) error {
	return nil
}

func (m *streamTestMsg) InProgress() error {
	return nil
}

func (m *streamTestMsg) Term() error {
	return nil
}

func (m *streamTestMsg) TermWithReason(string) error {
	return nil
}

var _ jetstream.Msg = (*streamTestMsg)(nil)

func TestProcessMessageOnDutyTransitionEnqueuesSnapshot(t *testing.T) {
	t.Parallel()

	grade := int32(3)
	userMarker := &livemapmarkers.UserMarker{
		UserId:   10,
		Job:      "police",
		JobGrade: &grade,
		X:        1,
		Y:        2,
	}
	data, err := proto.Marshal(userMarker)
	require.NoError(t, err)

	srv := &Server{
		tracker: &streamTestTracker{
			markers: []*livemapmarkers.UserMarker{userMarker},
		},
	}
	userInfo := newUserInfo(10, "police", 3, false)
	userOnDuty := false
	usersJobs := &permissionsattributes.JobGradeList{
		Jobs: map[string]int32{"police": 10},
	}

	var sent []*pblivemap.StreamResponse
	err = srv.processMessage(
		&streamTestMsg{
			subject: "$KV.userloc.police.3.10",
			headers: nats.Header{},
			data:    data,
		},
		userInfo,
		&userOnDuty,
		usersJobs,
		func(resp *pblivemap.StreamResponse) error {
			sent = append(sent, resp)
			return nil
		},
	)

	require.NoError(t, err)
	require.True(t, userOnDuty)
	require.Len(t, sent, 2)
	require.Equal(t, true, sent[0].GetUserOnDuty())
	_, ok := sent[0].Data.(*pblivemap.StreamResponse_Snapshot)
	require.True(t, ok)
	require.Len(t, sent[0].GetSnapshot().GetMarkers(), 1)
	require.Len(t, sent[1].GetUserUpdates().GetUpdates(), 1)
}
