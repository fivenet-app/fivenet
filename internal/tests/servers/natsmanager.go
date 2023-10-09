package servers

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var TestNATSServer *natsServer

type natsServer struct {
	pool     *dockertest.Pool
	resource *dockertest.Resource

	stopped bool
}

func init() {
	TestNATSServer = &natsServer{}
}

func (m *natsServer) Setup() {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	var err error
	m.pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %q", err)
	}
	m.pool.MaxWait = 3 * time.Minute

	// uses pool to try to connect to Docker
	err = m.pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %q", err)
	}

	// pulls an image, creates a container based on it and runs it
	m.resource, err = m.pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: "docker.io/library/nats",
			Tag:        "2.10.1-alpine3.18",
			Cmd: []string{
				"--user=fivenet",
				"--pass=changeme",
				"--jetstream",
			},
		},
		func(config *docker.HostConfig) {
			// set AutoRemove to true so that stopped container goes away by itself
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
		},
	)
	if err != nil {
		log.Fatalf("Could not start resource: %q", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := m.pool.Retry(func() error {
		client, err := nats.Connect(m.GetURL(), nats.Name("FiveNet natsmanager"))
		if err != nil {
			return err
		}
		defer client.Close()

		if !client.IsConnected() {
			return fmt.Errorf("NATS client not connected")
		}

		return nil
	}); err != nil {
		log.Fatalf("Could not connect to nats: %q", err)
	}
}

func (m *natsServer) GetURL() string {
	return fmt.Sprintf("nats://fivenet:changeme@localhost:%s", m.resource.GetPort("4222/tcp"))
}

func (m *natsServer) Stop() {
	if m.stopped {
		return
	}
	m.stopped = true

	// You can't defer this because os.Exit doesn't care for defer
	if err := m.pool.Purge(m.resource); err != nil {
		log.Fatalf("Could not purge container resource: %q", err)
	}
}
