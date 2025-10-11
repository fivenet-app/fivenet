package cmd

import (
	"github.com/fivenet-app/fivenet/v2025/pkg/storage"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/instance"
	"go.uber.org/fx"
)

type WorkerCmd struct {
	ModuleCentrum        bool `default:"true" help:"Start Centrum bot and housekeeper module"`
	ModuleUserTracker    bool `default:"true" help:"Start User tracker module"`
	ModuleJobsTimeclock  bool `default:"true" help:"Start Jobs timeclock housekeeper module"`
	ModuleDocsWorkflow   bool `default:"true" help:"Start Docstore Workflow module"`
	ModuleHousekeeper    bool `default:"true" help:"Start Housekeepr module"`
	ModuleUserInfoPoller bool `default:"true" help:"Start UserInfo poller module"`
}

func (c *WorkerCmd) Run(_ *Context) error {
	instance.SetComponent("worker")

	fxOpts := getFxBaseOpts(Cli.StartTimeout, true, true)
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
