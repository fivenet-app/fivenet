package cmd

import (
	"github.com/fivenet-app/fivenet/v2025/pkg/storage"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/instance"
	"go.uber.org/fx"
)

type WorkerCmd struct {
	ModuleCentrum        bool `help:"Start Centrum bot and housekeeper module" default:"true"`
	ModuleUserTracker    bool `help:"Start User tracker module" default:"true"`
	ModuleJobsTimeclock  bool `help:"Start Jobs timeclock housekeeper module" default:"true"`
	ModuleDocsWorkflow   bool `help:"Start Docstore Workflow module" default:"true"`
	ModuleHousekeeper    bool `help:"Start Housekeepr module" default:"true"`
	ModuleUserInfoPoller bool `help:"Start UserInfo poller module" default:"true"`
}

func (c *WorkerCmd) Run(ctx *Context) error {
	instance.SetComponent("worker")

	fxOpts := getFxBaseOpts(Cli.StartTimeout, true)
	fxOpts = append(fxOpts, FxDemoOpts()...)
	fxOpts = append(fxOpts, FxCronerOpts()...)

	if c.ModuleCentrum {
		fxOpts = append(fxOpts, FxCentrumOpts()...)
	}
	if c.ModuleUserTracker {
		fxOpts = append(fxOpts, FxTrackerOpts()...)
	}
	if c.ModuleJobsTimeclock {
		fxOpts = append(fxOpts, FxJobsHousekeeperOpts()...)
	}
	if c.ModuleDocsWorkflow {
		fxOpts = append(fxOpts, FxDocumentsWorkflowOpts()...)
	}
	if c.ModuleHousekeeper {
		fxOpts = append(fxOpts, FxHousekeeperOpts()...)
		fxOpts = append(fxOpts, fx.Invoke(func(*storage.MetricsCollector) {}))
	}
	if c.ModuleUserInfoPoller {
		fxOpts = append(fxOpts, FxUserInfoPollerOpts()...)
	}

	app := fx.New(fxOpts...)
	app.Run()

	return nil
}
