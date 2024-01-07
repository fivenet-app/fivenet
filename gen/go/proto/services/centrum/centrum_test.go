package centrum

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/galexrt/fivenet/gen/go/proto/services/centrum/manager"
	"github.com/galexrt/fivenet/gen/go/proto/services/centrum/state"
	"github.com/galexrt/fivenet/internal/tests/servers"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/coords/postals"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/pkg/server/audit"
	"github.com/galexrt/fivenet/pkg/tracker"
	"github.com/stretchr/testify/require"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

func TestMain(m *testing.M) {
	if err := servers.TestDBServer.Setup(); err != nil {
		fmt.Println("failed to setup mysql test server: %w", err)
		return
	}
	if err := servers.TestNATSServer.Setup(); err != nil {
		fmt.Println("failed to setup nats test server: %w", err)
		return
	}

	code := m.Run()

	defer servers.TestDBServer.Stop()
	defer servers.TestNATSServer.Stop()

	defer os.Exit(code)
}

func TestBasicCentrumFlow(t *testing.T) {
	defer servers.TestDBServer.Reset()

	db, err := servers.TestDBServer.DB()
	require.NoError(t, err)

	ctx := context.Background()
	logger := zap.NewNop()
	tp := tracesdk.NewTracerProvider()

	cfg, err := config.LoadTest()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	cfg.NATS.URL = servers.TestNATSServer.GetURL()
	cfg.Cache.RefreshTime = 1 * time.Hour

	js, err := servers.TestNATSServer.GetJS()
	require.NoError(t, err)

	fxLC := fxtest.NewLifecycle(t)

	p, err := perms.New(perms.Params{
		LC:     fxLC,
		Logger: logger,
		DB:     db,
		TP:     tp,
		JS:     js,
		Config: cfg,
	})
	require.NoError(t, err)

	aud := &audit.Noop{}

	c, err := mstlystcdata.NewCache(mstlystcdata.Params{
		LC:     fxLC,
		Logger: logger,
		TP:     tp,
		DB:     db,
		Config: cfg,
	})
	require.NoError(t, err)
	enricher := mstlystcdata.NewEnricher(c, p, cfg)

	state, err := state.New(state.Params{
		LC:     fxLC,
		Logger: logger,
		JS:     js,
		Config: cfg,
	})
	require.NoError(t, err)

	mgr := manager.New(manager.Params{
		LC:       fxLC,
		Logger:   logger,
		DB:       db,
		TP:       tp,
		JS:       js,
		Enricher: enricher,
		Config:   cfg,
		State:    state,
	})

	tracker := tracker.NewForTests(ctx)

	srv, err := NewServer(Params{
		LC:      fxLC,
		Logger:  logger,
		TP:      tp,
		DB:      db,
		Perms:   p,
		Audit:   aud,
		JS:      js,
		Config:  cfg,
		Manager: mgr,
		Tracker: tracker,
		Postals: postals.NewForTests(),
	})
	require.NoError(t, err)

	fxLC.RequireStart()
	defer fxLC.RequireStop()

	client, _, cancel := NewTestCentrumServiceClient(srv)
	defer cancel()
	_ = ctx
	_ = client

	// TODO add tests
}
