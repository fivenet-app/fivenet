package servers

import (
	"testing"
	"time"

	"github.com/fivenet-app/fivenet/pkg/events"
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

	// Auto stop server when test context is done
	go func() {
		<-m.t.Context().Done()
		m.Stop()
	}()
	m.server = ns
}

func (m *natsServer) GetURL() string {
	return m.server.ClientURL()
}

func (m *natsServer) GetClient() (*nats.Conn, error) {
	return nats.Connect(m.GetURL())
}

func (m *natsServer) GetJS() jetstream.JetStream {
	cli, err := m.GetClient()
	if err != nil {
		m.t.Fatalf("failed to get NATS client. %v", err)
		return nil
	}

	js, err := jetstream.New(cli)
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
	return fx.Provide(func() *events.JSWrapper {
		return &events.JSWrapper{
			JetStream: m.GetJS(),
		}
	})
}
