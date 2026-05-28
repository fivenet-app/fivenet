package servers

import (
	"testing"
	"time"

	"github.com/fivenet-app/fivenet/v2026/pkg/events"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

type natsServer struct {
	t *testing.T

	server *server.Server
}

func NewNATSServer(t *testing.T, setup bool) *natsServer {
	t.Helper()

	s := &natsServer{
		t: t,
	}

	if setup {
		s.Setup()
	}

	return s
}

func (m *natsServer) Setup() {
	opts := &server.Options{
		JetStream: true,
		Port:      -1,
	}

	// Initialize new server with options
	ns, err := server.NewServer(opts)
	require.NoError(m.t, err, "failed to setup mysql test server")

	// Start the server via goroutine
	go ns.Start()
	// Wait for server to be ready for connections
	require.True(
		m.t,
		ns.ReadyForConnections(8*time.Second),
		"nats: not ready connection after 8 seconds",
	)
	m.server = ns

	// Auto stop server when test is done
	m.t.Cleanup(m.Stop)
}

func (m *natsServer) GetURL() string {
	return m.server.ClientURL()
}

func (m *natsServer) GetConn() *nats.Conn {
	conn, err := nats.Connect(m.GetURL())
	require.NoError(m.t, err, "failed to get NATS client")

	return conn
}

func (m *natsServer) GetJS() jetstream.JetStream {
	conn := m.GetConn()

	js, err := jetstream.New(conn)
	require.NoError(m.t, err, "failed to create JetStream client")

	return js
}

func (m *natsServer) Stop() {
	if m == nil {
		return
	}

	m.server.Shutdown()
	m.server = nil
}

func (m *natsServer) FxProvide() fx.Option {
	return fx.Provide(func() events.Result {
		return events.Result{
			NC: m.GetConn(),
			JS: &events.JSWrapper{
				JetStream: m.GetJS(),
			},
		}
	})
}
