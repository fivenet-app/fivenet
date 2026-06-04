package cmd

import (
	"github.com/fivenet-app/fivenet/v2026/cmd/fxopts"
	"github.com/fivenet-app/fivenet/v2026/pkg/storage"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/instance"
	"go.uber.org/fx"
)

type WorkerCmd struct {
	ModuleCentrum        bool `default:"true" help:"Start Centrum bot and housekeeper module"`
	ModuleUserTracker    bool `default:"true" help:"Start User tracker module"`
	ModuleHousekeeper    bool `default:"true" help:"Start Housekeepr modules"`
	ModuleUserInfoPoller bool `default:"true" help:"Start UserInfo poller module"`
}

func (c *WorkerCmd) Run(cli *CLI) error {
	instance.SetComponent("worker")

	fxOpts := fxopts.GetFxBaseOpts(cli.StartTimeout, true, true)
	fxOpts = append(fxOpts, fxopts.FxDemoOpts()...)
	fxOpts = append(fxOpts, fxopts.FxCronerOpts()...)

	if c.ModuleCentrum {
		fxOpts = append(fxOpts, fxopts.FxCentrumOpts()...)
	}
	if c.ModuleUserTracker {
		fxOpts = append(fxOpts, fxopts.FxTrackerOpts()...)
	}
	if c.ModuleHousekeeper {
		fxOpts = append(fxOpts, fxopts.FxServiceHousekeeperOpts()...)
		fxOpts = append(fxOpts, fxopts.FxHousekeeperOpts()...)
		fxOpts = append(fxOpts, fx.Invoke(func(*storage.MetricsCollector) {}))
	}
	if c.ModuleUserInfoPoller {
		fxOpts = append(fxOpts, fxopts.FxUserInfoPollerOpts()...)
	}

	app := fx.New(fxOpts...)
	app.Run()

	return nil
}
