package servers

import (
	"testing"

	testnats "github.com/fivenet-app/fivenet/v2026/internal/tests/nats"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/events"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
)

type natsServer struct {
	server *testnats.Server
}

func NewNATSServer(t *testing.T) *natsServer {
	t.Helper()
	return &natsServer{server: testnats.NewServer(t, testnats.ServerOptions{})}
}

func (m *natsServer) GetURL() string {
	return m.server.GetConn().ConnectedUrl()
}

func (m *natsServer) GetConn() *nats.Conn {
	return m.server.GetConn()
}

func (m *natsServer) GetJS() jetstream.JetStream {
	return m.server.GetJetStream()
}

func (m *natsServer) FxProvide() fx.Option {
	return fx.Provide(
		func(cfg *config.Config, shutdowner fx.Shutdowner) events.Result {
			return events.Result{
				NC: m.GetConn(),
				JS: m.server.NewJSWrapper(cfg.NATS, shutdowner),
			}
		},
	)
}
