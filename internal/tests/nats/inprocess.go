package nats

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/nats-io/nats-server/v2/server"
	nats "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// Based on a-h "Adrian Hesketh" code from <https://github.com/nats-io/nats.go/issues/467#issuecomment-1771424369>
func NewInProcessNATSServer() (*nats.Conn, *events.JSWrapper, func(), error) {
	tmp, err := os.MkdirTemp("", "nats_test")
	if err != nil {
		err = fmt.Errorf("failed to create temp directory for NATS storage. %w", err)
		return nil, nil, nil, err
	}
	server, err := server.NewServer(&server.Options{
		DontListen: true, // Don't make a TCP socket.
		JetStream:  true,
		StoreDir:   tmp,
	})
	if err != nil {
		err = fmt.Errorf("failed to create NATS server. %w", err)
		return nil, nil, nil, err
	}
	// Add logs to stdout.
	// server.ConfigureLogger()
	server.Start()
	cleanup := func() {
		server.Shutdown()
		os.RemoveAll(tmp)
	}

	if !server.ReadyForConnections(time.Second * 5) {
		err = errors.New("failed to start server after 5 seconds")
		return nil, nil, nil, err
	}

	// Create a connection.
	conn, err := nats.Connect("", nats.InProcessServer(server))
	if err != nil {
		err = fmt.Errorf("failed to connect to server. %w", err)
		return nil, nil, nil, err
	}

	// Create a JetStream client.
	rawJs, err := jetstream.New(conn)
	if err != nil {
		err = fmt.Errorf("failed to create jetstream. %w", err)
		return nil, nil, nil, err
	}

	js := events.NewJSWrapper(rawJs, config.NATS{
		Replicas: 1,
	}, nil)

	return conn, js, cleanup, err
}
