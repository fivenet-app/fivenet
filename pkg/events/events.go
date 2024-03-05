package events

import (
	"context"
	"sync"
	"time"

	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/server/admin"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/fx"
)

var (
	metricsNats = promauto.NewGauge(prometheus.GaugeOpts{
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
	Config *config.BaseConfig
}

type Result struct {
	fx.Out

	JS nats.JetStreamContext
}

func New(p Params) (res Result, err error) {
	// Connect to NATS
	nc, err := nats.Connect(p.Config.NATS.URL, nats.Name("FiveNet"))
	if err != nil {
		return res, err
	}

	// Default `defaultAsyncPubAckInflight` is `4000` (`nats.go`)
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return res, err
	}
	res.JS = js

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

				case <-time.After(10 * time.Second):
					metricsNats.Set(float64(js.PublishAsyncPending()))
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
