package servers

import (
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

var TestNATSServer *natsServer

type natsServer struct {
	server *server.Server
}

func init() {
	TestNATSServer = &natsServer{}
}

func (m *natsServer) Setup() {
	opts := &server.Options{
		JetStream: true,
		Username:  "fivenet",
		Password:  "changeme",
		Port:      -1,
	}

	// Initialize new server with options
	ns, err := server.NewServer(opts)
	if err != nil {
		panic(err)
	}

	// Start the server via goroutine
	go ns.Start()
	// Wait for server to be ready for connections
	if !ns.ReadyForConnections(4 * time.Second) {
		panic("not ready for connection")
	}

	nc, err := nats.Connect(ns.ClientURL())
	if err != nil {
		panic(err)
	}

	js, err := nc.JetStream()
	if err != nil {
		panic(err)
	}

	buckets := []string{"stat", "basic", "list", "listnr"}
	for _, bucket := range buckets {
		_, err = js.CreateKeyValue(&nats.KeyValueConfig{
			Bucket:  bucket,
			History: 5,
			Storage: nats.MemoryStorage,
		})
		if err != nil {
			panic(err)
		}
	}

	m.server = ns
}

func (m *natsServer) GetURL() string {
	return m.server.ClientURL()
}

func (m *natsServer) GetClient() (*nats.Conn, error) {
	return nats.Connect(m.GetURL())
}

func (m *natsServer) GetJS() (nats.JetStreamContext, error) {
	cli, err := m.GetClient()
	if err != nil {
		return nil, err
	}

	return cli.JetStream()
}

func (m *natsServer) Stop() {
	if m == nil {
		return
	}

	m.Stop()
	m = nil
}
