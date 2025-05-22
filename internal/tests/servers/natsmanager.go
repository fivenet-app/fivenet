package servers

import (
	"testing"
	"time"

	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
)

type natsServer struct {
	t *testing.T

	server *server.Server
}

func NewNATSServer(t *testing.T, setup bool) *natsServer {
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
	if err != nil {
		m.t.Fatalf("failed to setup mysql test server. %v", err)
	}

	// Start the server via goroutine
	go ns.Start()
	// Wait for server to be ready for connections
	if !ns.ReadyForConnections(8 * time.Second) {
		m.t.Fatal("nats: not ready connection after 8 seconds")
	}
	m.server = ns

	// Auto stop server when test is done
	m.t.Cleanup(m.Stop)
}

func (m *natsServer) GetURL() string {
	return m.server.ClientURL()
}

func (m *natsServer) GetConn() *nats.Conn {
	conn, err := nats.Connect(m.GetURL())
	if err != nil {
		m.t.Fatalf("failed to get NATS client. %v", err)
		return nil
	}

	return conn
}

func (m *natsServer) GetJS() jetstream.JetStream {
	conn := m.GetConn()

	js, err := jetstream.New(conn)
	if err != nil {
		m.t.Fatalf("failed to create JetStream client. %v", err)
		return nil
	}

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
