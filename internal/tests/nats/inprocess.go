package nats

import (
	"testing"
	"time"

	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/events"
	"github.com/nats-io/nats-server/v2/server"
	natsgo "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
)

// ServerOptions configures a test NATS server.
type ServerOptions struct {
	// InProcess connects the client directly to the server, without opening a TCP listener.
	InProcess bool
	// Config configures the JetStream wrapper returned by GetJS.
	Config config.NATS
}

// Server is a NATS server, client connection, and JetStream client owned by a test.
type Server struct {
	server *server.Server
	conn   *natsgo.Conn
	rawJS  jetstream.JetStream
	cfg    config.NATS
}

// NewServer starts a JetStream-enabled NATS server and registers cleanup with t.
func NewServer(t *testing.T, opts ServerOptions) *Server {
	t.Helper()
	started := time.Now()
	if opts.Config.Replicas == 0 {
		opts.Config.Replicas = 1
	}

	tmp := t.TempDir()
	ns, err := server.NewServer(&server.Options{
		JetStream:  true,
		StoreDir:   tmp,
		Port:       -1,
		DontListen: opts.InProcess,
	})
	if err != nil {
		t.Fatalf("create NATS server: %v", err)
	}

	s := &Server{
		server: ns,
		cfg:    opts.Config,
	}
	ns.Start()
	if !ns.ReadyForConnections(8 * time.Second) {
		_ = s.Close()
		t.Fatal("NATS server was not ready for connections after 8 seconds")
	}

	if opts.InProcess {
		s.conn, err = natsgo.Connect("", natsgo.InProcessServer(ns))
	} else {
		s.conn, err = natsgo.Connect(ns.ClientURL())
	}
	if err != nil {
		_ = s.Close()
		t.Fatalf("connect to NATS server: %v", err)
	}

	s.rawJS, err = jetstream.New(s.conn)
	if err != nil {
		_ = s.Close()
		t.Fatalf("create JetStream client: %v", err)
	}

	t.Cleanup(func() {
		if err := s.Close(); err != nil {
			t.Errorf("clean up NATS server: %v", err)
		}
	})
	t.Logf(
		"test NATS server ready: in_process=%t duration=%s",
		opts.InProcess,
		time.Since(started).Round(time.Millisecond),
	)

	return s
}

func (s *Server) GetConn() *natsgo.Conn {
	return s.conn
}

// GetJetStream returns the unwrapped JetStream client.
func (s *Server) GetJetStream() jetstream.JetStream {
	return s.rawJS
}

// GetJS returns a JetStream client that applies the fixture's NATS configuration.
func (s *Server) GetJS() *events.JSWrapper {
	return s.NewJSWrapper(s.cfg, nil)
}

// NewJSWrapper returns a JetStream wrapper using cfg and shutdowner.
func (s *Server) NewJSWrapper(cfg config.NATS, shutdowner fx.Shutdowner) *events.JSWrapper {
	return events.NewJSWrapper(s.rawJS, cfg, shutdowner)
}

// Close stops the client and server, and removes the server's storage directory.
func (s *Server) Close() error {
	if s == nil {
		return nil
	}
	if s.conn != nil {
		s.conn.Close()
		s.conn = nil
	}
	if s.server != nil {
		s.server.Shutdown()
		s.server = nil
	}

	return nil
}
