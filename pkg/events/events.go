package events

import (
	"context"
	"sync"
	"time"

	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/server/admin"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/fx"
)

// Default `defaultAsyncPubAckInflight` is `4000` (`nats.go`)
const DefaultDefaultAsyncPubAckInflight = 256

var (
	metricNATSAsyncPending = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: admin.MetricsNamespace,
		Subsystem: "nats",
		Name:      "jetstream_async_pending_count",
		Help:      "NATS JetStreamn async pending count.",
	})
)

type Subject string

type Topic string

type Type string

type Params struct {
	fx.In

	LC     fx.Lifecycle
	Config *config.Config
}

type Result struct {
	fx.Out

	JS jetstream.JetStream
}

func New(p Params) (res Result, err error) {
	// Connect to NATS
	nc, err := nats.Connect(p.Config.NATS.URL, nats.Name("FiveNet"))
	if err != nil {
		return res, err
	}

	res.JS, err = jetstream.New(nc, jetstream.WithPublishAsyncMaxPending(DefaultDefaultAsyncPubAckInflight))
	if err != nil {
		return res, err
	}

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
