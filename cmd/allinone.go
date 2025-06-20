package cmd

import (
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/instance"
	"go.uber.org/fx"
)

type AllInOneCmd struct{}

func (c *AllInOneCmd) Run(ctx *Context) error {
	instance.SetComponent("allinone")

	fxOpts := getFxBaseOpts(Cli.StartTimeout, true)
	fxOpts = append(fxOpts, FxServerOpts()...)
	fxOpts = append(fxOpts, FxDemoOpts()...)
	fxOpts = append(fxOpts, FxUserInfoPollerOpts()...)
	fxOpts = append(fxOpts, FxCronerOpts()...)
	fxOpts = append(fxOpts, FxJobsHousekeeperOpts()...)
	fxOpts = append(fxOpts, FxHousekeeperOpts()...)
	fxOpts = append(fxOpts, FxCentrumOpts()...)
	fxOpts = append(fxOpts, FxTrackerOpts()...)
	fxOpts = append(fxOpts, FxDocumentsWorkflowOpts()...)

	app := fx.New(fxOpts...)
	app.Run()

	return nil
}
