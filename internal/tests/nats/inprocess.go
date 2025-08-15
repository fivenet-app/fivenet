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

// NewInProcessNATSServer based on a-h "Adrian Hesketh" code from <https://github.com/nats-io/nats.go/issues/467#issuecomment-1771424369>
// Used to create an in-process NATS server for testing purposes.
// It returns a connection to the server, a JetStream client, a cleanup function to call after tests, and an error if any.
func NewInProcessNATSServer() (*nats.Conn, *events.JSWrapper, func() error, error) {
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
	cleanup := func() error {
		server.Shutdown()
		if err := os.RemoveAll(tmp); err != nil {
			return fmt.Errorf("failed to remove temp directory for NATS storage. %w", err)
		}
		return nil
	}

	if !server.ReadyForConnections(5 * time.Second) {
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
