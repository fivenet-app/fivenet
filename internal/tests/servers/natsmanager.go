package servers

import (
	"fmt"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

var TestNATSServer = &natsServer{}

type natsServer struct {
	server *server.Server
}

func (m *natsServer) Setup() error {
	opts := &server.Options{
		JetStream: true,
		Port:      -1,
	}

	// Initialize new server with options
	ns, err := server.NewServer(opts)
	if err != nil {
		return err
	}

	// Start the server via goroutine
	go ns.Start()
	// Wait for server to be ready for connections
	if !ns.ReadyForConnections(4 * time.Second) {
		return fmt.Errorf("not ready for connection")
	}

	m.server = ns

	return nil
}

func (m *natsServer) GetURL() string {
	return m.server.ClientURL()
}

func (m *natsServer) GetClient() (*nats.Conn, error) {
	return nats.Connect(m.GetURL())
}

func (m *natsServer) GetJS() (jetstream.JetStream, error) {
	cli, err := m.GetClient()
	if err != nil {
		return nil, err
	}

	return jetstream.New(cli)
}

func (m *natsServer) Stop() {
	if m == nil {
		return
	}

	m.server.Shutdown()
	m.server = nil
}

func (m *natsServer) Reset() {
	if m == nil {
		return
	}

	m.Stop()
	m.Setup()
}
