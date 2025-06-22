package cmd

import (
	"github.com/fivenet-app/fivenet/v2025/pkg/croner"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbsync"
	"github.com/fivenet-app/fivenet/v2025/pkg/demo"
	"github.com/fivenet-app/fivenet/v2025/pkg/discord"
	"github.com/fivenet-app/fivenet/v2025/pkg/discord/commands"
	"github.com/fivenet-app/fivenet/v2025/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/v2025/pkg/server"
	"github.com/fivenet-app/fivenet/v2025/pkg/tracker/manager"
	"github.com/fivenet-app/fivenet/v2025/pkg/userinfo"
	centrumbot "github.com/fivenet-app/fivenet/v2025/services/centrum/bot"
	centrumhousekeeper "github.com/fivenet-app/fivenet/v2025/services/centrum/housekeeper"
	pbdocuments "github.com/fivenet-app/fivenet/v2025/services/documents"
	pbjobs "github.com/fivenet-app/fivenet/v2025/services/jobs"
	"go.uber.org/fx"
)

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

func FxJobsHousekeeperOpts() []fx.Option {
	return []fx.Option{
		fx.Invoke(func(*pbjobs.Housekeeper) {}),
	}
}

func FxDocumentsWorkflowOpts() []fx.Option {
	return []fx.Option{
		fx.Invoke(func(*pbdocuments.Workflow) {}),
	}
}

func FxDiscordOpts() []fx.Option {
	return []fx.Option{
		fx.Invoke(func(*discord.Bot) {}),
		fx.Invoke(func(*commands.Cmds) {}),
	}
}

func FxDBSyncOpts() []fx.Option {
	return []fx.Option{
		fx.Invoke(func(*dbsync.Sync) {}),
	}
}
