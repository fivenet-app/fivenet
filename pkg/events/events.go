package events

import (
	"context"
	"sync"
	"time"

	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/server/admin"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Default `defaultAsyncPubAckInflight` is `4000` (`nats.go`)
const DefaultDefaultAsyncPubAckInflight = 256

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

	JS *JSWrapper
}

func New(p Params) (res Result, err error) {
	// Connect to NATS
	nc, err := nats.Connect(p.Config.NATS.URL,
		nats.Name("FiveNet"),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			if !nc.IsClosed() {
				p.Logger.Error("nats: disconnected", zap.Error(err))
			}
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			p.Logger.Info("nats: reconnected")
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			if err := nc.LastError(); err != nil {
				p.Logger.Error("nats: connection closed", zap.Error(err))

				if err := p.Shutdowner.Shutdown(fx.ExitCode(1)); err != nil {
					p.Logger.Fatal("failed to shutdown app after nats connection close", zap.Error(err))
				}
			}
		}),
	)
	if err != nil {
		return res, err
	}

	js, err := jetstream.New(nc, jetstream.WithPublishAsyncMaxPending(DefaultDefaultAsyncPubAckInflight))
	if err != nil {
		return res, err
	}

	res.JS = NewJSWrapper(js, p.Config.NATS, p.Shutdowner)

	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	// Collect basic metric
	p.LC.Append(fx.StartHook(func(_ context.Context) error {
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
