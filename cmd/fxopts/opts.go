package fxopts

import (
	"github.com/fivenet-app/fivenet/v2026/pkg/croner"
	"github.com/fivenet-app/fivenet/v2026/pkg/demo"
	"github.com/fivenet-app/fivenet/v2026/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/v2026/pkg/server"
	"github.com/fivenet-app/fivenet/v2026/pkg/tracker/manager"
	"github.com/fivenet-app/fivenet/v2026/pkg/userinfo"
	pbcalendar "github.com/fivenet-app/fivenet/v2026/services/calendar"
	centrumbot "github.com/fivenet-app/fivenet/v2026/services/centrum/bot"
	centrumhousekeeper "github.com/fivenet-app/fivenet/v2026/services/centrum/housekeeper"
	pbcitizens "github.com/fivenet-app/fivenet/v2026/services/citizens"
	pbdocuments "github.com/fivenet-app/fivenet/v2026/services/documents"
	pbjobs "github.com/fivenet-app/fivenet/v2026/services/jobs"
	pbvehicles "github.com/fivenet-app/fivenet/v2026/services/vehicles"
	"go.uber.org/fx"
)

const DefaultStopTimeout = 180 // seconds

// Option groups for fx modules, to be reused across commands.

func FxServerOpts() []fx.Option {
	return []fx.Option{
		fx.Invoke(func(server.HTTPServer) {}),
	}
}

func FxDemoOpts() []fx.Option {
	return []fx.Option{
		fx.Invoke(func(*demo.Demo) {}),
	}
}

func FxUserInfoPollerOpts() []fx.Option {
	return []fx.Option{
		fx.Invoke(func(*userinfo.Poller) {}),
	}
}

func FxCronerOpts() []fx.Option {
	return []fx.Option{
		fx.Invoke(func(*croner.Scheduler) {}),
		fx.Invoke(func(*croner.Executor) {}),
	}
}

func FxHousekeeperOpts() []fx.Option {
	return []fx.Option{
		fx.Invoke(func(*housekeeper.Housekeeper) {}),
	}
}

func FxCentrumOpts() []fx.Option {
	return []fx.Option{
		fx.Invoke(func(*centrumbot.Manager) {}),
		fx.Invoke(func(*centrumhousekeeper.Housekeeper) {}),
	}
}

func FxTrackerOpts() []fx.Option {
	return []fx.Option{
		fx.Invoke(func(*manager.Manager) {}),
	}
}

func FxServiceHousekeeperOpts() []fx.Option {
	return []fx.Option{
		fx.Invoke(func(*pbcitizens.Housekeeper) {}),
		fx.Invoke(func(*pbjobs.Housekeeper) {}),
		fx.Invoke(func(*pbvehicles.Housekeeper) {}),
		fx.Invoke(func(*pbdocuments.Workflow) {}),
		fx.Invoke(func(*pbcalendar.BirthdaySyncer) {}),
	}
}
