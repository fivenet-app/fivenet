package userinfo

import (
	"context"
	"errors"
	"testing"
	"time"

	accounts "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts"
	notificationsevents "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications/events"
	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/notifi"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
)

type retrieverTestMsg struct {
	subject   string
	data      []byte
	ackCalls  int
	termCalls int
}

func (m *retrieverTestMsg) Metadata() (*jetstream.MsgMetadata, error) { return nil, nil }
func (m *retrieverTestMsg) Data() []byte                              { return m.data }
func (m *retrieverTestMsg) Headers() nats.Header                      { return nil }
func (m *retrieverTestMsg) Subject() string                           { return m.subject }
func (m *retrieverTestMsg) Reply() string                             { return "" }
func (m *retrieverTestMsg) Ack() error {
	m.ackCalls++
	return nil
}
func (m *retrieverTestMsg) DoubleAck(context.Context) error  { return nil }
func (m *retrieverTestMsg) Nak() error                       { return nil }
func (m *retrieverTestMsg) NakWithDelay(time.Duration) error { return nil }
func (m *retrieverTestMsg) InProgress() error                { return nil }
func (m *retrieverTestMsg) Term() error {
	m.termCalls++
	return nil
}
func (m *retrieverTestMsg) TermWithReason(string) error { return nil }

var _ jetstream.Msg = (*retrieverTestMsg)(nil)

type retrieverTestNotifi struct {
	notifi.INotifi

	err   error
	calls int
}

func (n *retrieverTestNotifi) SendUserEvent(
	context.Context,
	int32,
	*notificationsevents.UserEvent,
) error {
	n.calls++
	return n.err
}

func userInfoChangedJSON(t *testing.T, userID int32, accountID int64) []byte {
	t.Helper()
	data, err := protojson.Marshal(&pbuserinfo.UserInfoChanged{
		UserId:    userID,
		AccountId: accountID,
	})
	require.NoError(t, err)
	return data
}

func TestRetrieverHandleMsgMalformedPayloadTerminatesMessage(t *testing.T) {
	t.Parallel()

	msg := &retrieverTestMsg{subject: "userinfo.changed", data: []byte("{")}
	retriever := &Retriever{logger: zap.NewNop(), notifi: &retrieverTestNotifi{}}

	retriever.handleMsg(msg)

	assert.Equal(t, 1, msg.termCalls)
	assert.Zero(t, msg.ackCalls)
}

func TestRetrieverHandleMsgDeliveryFailureLeavesMessageUnacknowledged(t *testing.T) {
	t.Parallel()

	msg := &retrieverTestMsg{
		subject: "userinfo.changed",
		data:    userInfoChangedJSON(t, 42, 7),
	}
	notification := &retrieverTestNotifi{err: errors.New("temporary failure")}
	retriever := &Retriever{logger: zap.NewNop(), notifi: notification}

	retriever.handleMsg(msg)

	assert.Equal(t, 1, notification.calls)
	assert.Zero(t, msg.ackCalls)
	assert.Zero(t, msg.termCalls)
}

func TestRetrieverHandleMsgSuccessfulDeliveryAcknowledgesMessage(t *testing.T) {
	t.Parallel()

	msg := &retrieverTestMsg{
		subject: "userinfo.changed",
		data:    userInfoChangedJSON(t, 42, 7),
	}
	notification := &retrieverTestNotifi{}
	retriever := &Retriever{logger: zap.NewNop(), notifi: notification}

	retriever.handleMsg(msg)

	assert.Equal(t, 1, notification.calls)
	assert.Equal(t, 1, msg.ackCalls)
	assert.Zero(t, msg.termCalls)
}

func TestRetrieverHandleMsgInvalidIDsTerminateMessage(t *testing.T) {
	t.Parallel()

	msg := &retrieverTestMsg{
		subject: "userinfo.changed",
		data:    userInfoChangedJSON(t, 0, 7),
	}
	notification := &retrieverTestNotifi{}
	retriever := &Retriever{logger: zap.NewNop(), notifi: notification}

	retriever.handleMsg(msg)

	assert.Equal(t, 0, notification.calls)
	assert.Zero(t, msg.ackCalls)
	assert.Equal(t, 1, msg.termCalls)
}

func TestCheckAndSetSuperuserAcceptsConfigAdminMembership(t *testing.T) {
	t.Parallel()

	retriever := &Retriever{
		configAdminGroups: []string{"config-admin"},
	}

	userInfo := &pbuserinfo.UserInfo{
		Groups: &accounts.AccountGroups{
			Groups: []string{"config-admin"},
		},
		License: "license-42",
	}

	retriever.checkAndSetSuperuser(userInfo)

	assert.True(t, userInfo.GetCanBeSuperuser())
	assert.True(t, userInfo.GetCanBeConfigAdmin())
	assert.False(t, userInfo.GetSuperuser())
}
