package events

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/reqs"
	"github.com/fivenet-app/fivenet/v2025/pkg/server/admin"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Default `defaultAsyncPubAckInflight` is `4000` (`nats.go`).
const DefaultDefaultAsyncPubAckInflight = 256

var Module = fx.Module("events",
	fx.Provide(
		New,
	),
)

var metricNATSAsyncPending = promauto.NewGauge(prometheus.GaugeOpts{
	Namespace: admin.MetricsNamespace,
	Subsystem: "nats",
	Name:      "jetstream_async_pending_count",
	Help:      "NATS JetStreamn async pending count.",
})

type Subject string

type Topic string

type Type string

type Params struct {
	fx.In

	LC         fx.Lifecycle
	Shutdowner fx.Shutdowner

	Logger *zap.Logger
	Config *config.Config
}

type Result struct {
	fx.Out

	NC *nats.Conn
	JS *JSWrapper

	Req *reqs.NatsReqs
}

func New(p Params) (res Result, err error) {
	logger := p.Logger.Named("events")

	connOpts := []nats.Option{
		nats.Name("FiveNet"),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			if !nc.IsClosed() {
				logger.Error("nats: disconnected", zap.Error(err))
			}
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			logger.Info("nats: reconnected")
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			if err := nc.LastError(); err != nil {
				logger.Error("nats: connection closed", zap.Error(err))

				if err := p.Shutdowner.Shutdown(fx.ExitCode(1)); err != nil {
					logger.Fatal(
						"failed to shutdown app after nats connection close",
						zap.Error(err),
					)
				}
			}
		}),
	}

	if p.Config.NATS.NKey != nil {
		nKeyOpt, err := nats.NkeyOptionFromSeed(*p.Config.NATS.NKey)
		if err != nil {
			return res, fmt.Errorf("failed to read nats nkey. %w", err)
		}

		connOpts = append(connOpts, nKeyOpt)
	} else if p.Config.NATS.Creds != nil {
		// Use NATS credentials file if provided
		connOpts = append(connOpts, nats.UserCredentials(*p.Config.NATS.Creds))
	}

	// Connect to NATS
	nc, err := nats.Connect(p.Config.NATS.URL, connOpts...)
	if err != nil {
		return res, err
	}
	res.NC = nc

	// Create JetStream context
	js, err := jetstream.New(
		nc,
		jetstream.WithPublishAsyncMaxPending(DefaultDefaultAsyncPubAckInflight),
	)
	if err != nil {
		return res, err
	}

	res.JS = NewJSWrapper(js, p.Config.NATS, p.Shutdowner)

	res.Req = reqs.NewNatsReqs(nc)
	if err := res.Req.ValidateAll(); err != nil {
		if !p.Config.IgnoreRequirements {
			return res, fmt.Errorf("failed to validate nats requirements. %w", err)
		}
		p.Logger.Warn("ignoring failed nats requirements", zap.Error(err))
	}

	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	// Run migrations and collect basic metric
	p.LC.Append(fx.StartHook(func(_ context.Context) error {
		if err := runMigrations(ctx, p.Logger, res.JS); err != nil {
			return fmt.Errorf("failed to run nats migrations. %w", err)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return

				case <-time.After(5 * time.Second):
					metricNATSAsyncPending.Set(float64(res.JS.PublishAsyncPending()))
				}
			}
		}()

		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		wg.Wait()

		return nc.Drain()
	}))

	return res, nil
}
