package reqs

import (
	"fmt"
	"sync"

	"github.com/nats-io/nats.go"
)

type NatsReqs struct {
	mu sync.Mutex

	nc      *nats.Conn
	version string
}

func NewNatsReqs(nc *nats.Conn) *NatsReqs {
	return &NatsReqs{
		nc: nc,
	}
}

func (r *NatsReqs) GetVersion() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.version
}

func (r *NatsReqs) ValidateAll() error {
	return r.ValidateVersion()
}

func (r *NatsReqs) ValidateVersion() error {
	r.version = r.nc.ConnectedServerVersion()

	// Skip if version is not set
	if r.version == "" {
		return nil
	}

	// Example: Check if the version is at least "2.11.4"
	if r.version < "2.11.4" {
		return fmt.Errorf("NATS version %s is not supported, requires at least NATS 2.11.4", r.version)
	}

	return nil
}
