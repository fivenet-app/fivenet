package cmd

import (
	"github.com/fivenet-app/fivenet/v2025/pkg/croner"
	"github.com/fivenet-app/fivenet/v2025/pkg/demo"
	"github.com/fivenet-app/fivenet/v2025/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/v2025/pkg/tracker"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/centrumbot"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/centrummanager"
	pbdocuments "github.com/fivenet-app/fivenet/v2025/services/documents"
	pbjobs "github.com/fivenet-app/fivenet/v2025/services/jobs"
	"go.uber.org/fx"
)

type WorkerCmd struct {
	ModuleCentrumBot         bool `help:"Start Centrum bot module" default:"true"`
	ModuleCentrumHousekeeper bool `help:"Start Centrum Housekeeper module" default:"true"`
	ModuleUserTracker        bool `help:"Start User tracker module" default:"true"`
	ModuleJobsTimeclock      bool `help:"Start Jobs timeclock housekeeper module" default:"true"`
	ModuleDocsWorkflow       bool `help:"Start Docstore Workflow module" default:"true"`
	ModuleHousekeeper        bool `help:"Start Housekeepr module" default:"true"`
}

func (c *WorkerCmd) Run(ctx *Context) error {
	fxOpts := getFxBaseOpts(Cli.StartTimeout, true)
	fxOpts = append(fxOpts, fx.Invoke(func(*demo.Demo) {}))

	if c.ModuleCentrumBot {
		fxOpts = append(fxOpts, fx.Invoke(func(*centrumbot.Manager) {}))
	}
	if c.ModuleCentrumHousekeeper {
		fxOpts = append(fxOpts, fx.Invoke(func(*centrummanager.Housekeeper) {}))
	}
	if c.ModuleUserTracker {
		fxOpts = append(fxOpts, fx.Invoke(func(*tracker.Manager) {}))
	}
	if c.ModuleJobsTimeclock {
		fxOpts = append(fxOpts, fx.Invoke(func(*pbjobs.Housekeeper) {}))
	}
	if c.ModuleDocsWorkflow {
		fxOpts = append(fxOpts, fx.Invoke(func(*pbdocuments.Workflow) {}))
	}
	if c.ModuleHousekeeper {
		fxOpts = append(fxOpts, fx.Invoke(func(*housekeeper.Housekeeper) {}))
	}

	// Always run cron agent in worker
	fxOpts = append(fxOpts, fx.Invoke(func(*croner.Executor) {}))

	app := fx.New(fxOpts...)
	app.Run()

	return nil
}
