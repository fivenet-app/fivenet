package cmd

import (
	"github.com/fivenet-app/fivenet/pkg/croner"
	"github.com/fivenet-app/fivenet/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/pkg/server/audit"
	"github.com/fivenet-app/fivenet/pkg/tracker"
	"github.com/fivenet-app/fivenet/services/centrum/centrumbot"
	"github.com/fivenet-app/fivenet/services/centrum/centrummanager"
	pbdocstore "github.com/fivenet-app/fivenet/services/docstore"
	pbjobs "github.com/fivenet-app/fivenet/services/jobs"
	"go.uber.org/fx"
)

type WorkerCmd struct {
	ModuleAuditRetention     bool `help:"Start Audit log retention module" default:"true"`
	ModuleCentrumBot         bool `help:"Start Centrum bot module" default:"true"`
	ModuleCentrumHousekeeper bool `help:"Start Centrum Housekeeper module" default:"true"`
	ModuleUserTracker        bool `help:"Start User tracker module" default:"true"`
	ModuleJobsTimeclock      bool `help:"Start Jobs timeclock housekeeper module" default:"true"`
	ModuleDocsWorkflow       bool `help:"Start Docstore Workflow module" default:"true"`
	ModuleHousekeeper        bool `help:"Start Housekeepr module" default:"true"`
}

func (c *WorkerCmd) Run(ctx *Context) error {
	fxOpts := getFxBaseOpts(Cli.StartTimeout, true)

	if c.ModuleAuditRetention {
		fxOpts = append(fxOpts, fx.Invoke(func(*audit.Retention) {}))
	}
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
		fxOpts = append(fxOpts, fx.Invoke(func(*pbdocstore.Workflow) {}))
	}
	if c.ModuleHousekeeper {
		fxOpts = append(fxOpts, fx.Invoke(func(*housekeeper.Housekeeper) {}))
	}

	// Only run cron agent in worker
	fxOpts = append(fxOpts, fx.Invoke(func(*croner.Agent) {}))

	app := fx.New(fxOpts...)
	app.Run()

	return nil
}
