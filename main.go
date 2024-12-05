package main

import (
	"os"
	"runtime"
	"time"

	"github.com/alecthomas/kong"
	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	// GRPC Services
	pbauth "github.com/fivenet-app/fivenet/gen/go/proto/services/auth"
	pbcalendar "github.com/fivenet-app/fivenet/gen/go/proto/services/calendar"
	pbcentrum "github.com/fivenet-app/fivenet/gen/go/proto/services/centrum"
	pbcitizenstore "github.com/fivenet-app/fivenet/gen/go/proto/services/citizenstore"
	pbcompletor "github.com/fivenet-app/fivenet/gen/go/proto/services/completor"
	pbdmv "github.com/fivenet-app/fivenet/gen/go/proto/services/dmv"
	pbdocstore "github.com/fivenet-app/fivenet/gen/go/proto/services/docstore"
	pbinternet "github.com/fivenet-app/fivenet/gen/go/proto/services/internet"
	pbjobs "github.com/fivenet-app/fivenet/gen/go/proto/services/jobs"
	pblivemapper "github.com/fivenet-app/fivenet/gen/go/proto/services/livemapper"
	pbmailer "github.com/fivenet-app/fivenet/gen/go/proto/services/mailer"
	pbnotificator "github.com/fivenet-app/fivenet/gen/go/proto/services/notificator"
	pbqualifications "github.com/fivenet-app/fivenet/gen/go/proto/services/qualifications"
	pbrector "github.com/fivenet-app/fivenet/gen/go/proto/services/rector"
	pbstats "github.com/fivenet-app/fivenet/gen/go/proto/services/stats"
	pbwiki "github.com/fivenet-app/fivenet/gen/go/proto/services/wiki"

	// Modules
	"github.com/fivenet-app/fivenet/gen/go/proto/services/centrum/centrumbot"
	"github.com/fivenet-app/fivenet/gen/go/proto/services/centrum/centrummanager"
	"github.com/fivenet-app/fivenet/gen/go/proto/services/centrum/centrumstate"
	"github.com/fivenet-app/fivenet/internal/modules"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/pkg/coords/postals"
	"github.com/fivenet-app/fivenet/pkg/croner"
	"github.com/fivenet-app/fivenet/pkg/discord"
	"github.com/fivenet-app/fivenet/pkg/events"
	"github.com/fivenet-app/fivenet/pkg/grpc"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/pkg/html/htmldiffer"
	"github.com/fivenet-app/fivenet/pkg/html/htmlsanitizer"
	"github.com/fivenet-app/fivenet/pkg/lang"
	"github.com/fivenet-app/fivenet/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/pkg/notifi"
	"github.com/fivenet-app/fivenet/pkg/perms"
	"github.com/fivenet-app/fivenet/pkg/server"
	"github.com/fivenet-app/fivenet/pkg/server/admin"
	"github.com/fivenet-app/fivenet/pkg/server/api"
	"github.com/fivenet-app/fivenet/pkg/server/audit"
	"github.com/fivenet-app/fivenet/pkg/server/filestore"
	"github.com/fivenet-app/fivenet/pkg/server/images"
	"github.com/fivenet-app/fivenet/pkg/server/oauth2"
	"github.com/fivenet-app/fivenet/pkg/server/wk"
	"github.com/fivenet-app/fivenet/pkg/storage"
	"github.com/fivenet-app/fivenet/pkg/tracker"
	"github.com/fivenet-app/fivenet/query"
)

type Context struct{}

type ServerCmd struct {
	ModuleCronAgent bool `help:"Run the cron agent, should only be used for single container/binary deployments." default:"false"`
}

func (c *ServerCmd) Run(ctx *Context) error {
	fxOpts := getFxBaseOpts(cli.StartTimeout)
	fxOpts = append(fxOpts, fx.Invoke(func(server.HTTPServer) {}))

	if c.ModuleCronAgent {
		fxOpts = append(fxOpts, fx.Invoke(func(*croner.Agent) {}))
		fxOpts = append(fxOpts, fx.Invoke(func(*pbjobs.Housekeeper) {}))
	}

	app := fx.New(fxOpts...)
	app.Run()

	return nil
}

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
	fxOpts := getFxBaseOpts(cli.StartTimeout)

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

type DiscordCmd struct {
	ModuleCronAgent bool `help:"Run the cron agent." default:"false"`
}

func (c *DiscordCmd) Run(ctx *Context) error {
	fxOpts := getFxBaseOpts(cli.StartTimeout)
	fxOpts = append(fxOpts, fx.Invoke(func(*discord.Bot) {}))

	if c.ModuleCronAgent {
		fxOpts = append(fxOpts, fx.Invoke(func(*croner.Agent) {}))
	}

	app := fx.New(fxOpts...)
	app.Run()

	return nil
}

var cli struct {
	Config       string        `help:"Alternative config file (env var: FIVENET_CONFIG_FILE)"`
	StartTimeout time.Duration `help:"App start timeout duration" default:"180s"`

	Server  ServerCmd  `cmd:"" help:"Run FiveNet server."`
	Worker  WorkerCmd  `cmd:"" help:"Run FiveNet worker."`
	Discord DiscordCmd `cmd:"" help:"Run FiveNet Discord bot."`
}

func getFxBaseOpts(startTimeout time.Duration) []fx.Option {
	return []fx.Option{
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.StartTimeout(startTimeout),

		admin.Module,
		appconfig.Module,
		audit.Module,
		audit.RetentionModule,
		auth.AuthModule,
		auth.PermsModule,
		auth.TokenMgrModule,
		centrumbot.Module,
		config.Module,
		croner.AgentModule,
		croner.HandlerModule,
		croner.Module,
		croner.SchedulerModule,
		discord.BotModule,
		events.Module,
		grpc.ServerModule,
		htmlsanitizer.Module,
		htmldiffer.Module,
		lang.Module,
		centrummanager.HousekeeperModule,
		centrummanager.Module,
		modules.LoggerModule,
		modules.TracerProviderModule,
		perms.Module,
		query.Module,
		server.HTTPEngineModule,
		server.HTTPServerModule,
		centrumstate.StateModule,
		storage.Module,
		housekeeper.Module,

		fx.Provide(
			mstlystcdata.NewCache,
			mstlystcdata.NewEnricher,
			mstlystcdata.NewSearcher,
			mstlystcdata.NewUserAwareEnricher,
			notifi.New,
			postals.New,
			tracker.New,
			tracker.NewManager,
			userinfo.NewUIRetriever,

			// GRPC Service Helpers, Housekeepers and Co.
			pbjobs.NewHousekeeper,
			pbdocstore.NewWorkflow,

			// HTTP Services
			server.AsService(api.New),
			server.AsService(filestore.New),
			server.AsService(images.New),
			server.AsService(oauth2.New),
			server.AsService(wk.New),
		),

		// GRPC Services
		fx.Provide(
			grpc.AsService(pbauth.NewServer),
			grpc.AsService(pbcalendar.NewServer),
			grpc.AsService(pbcentrum.NewServer),
			grpc.AsService(pbcitizenstore.NewServer),
			grpc.AsService(pbcompletor.NewServer),
			grpc.AsService(pbdmv.NewServer),
			grpc.AsService(pbdocstore.NewServer),
			grpc.AsService(pbjobs.NewServer),
			grpc.AsService(pblivemapper.NewServer),
			grpc.AsService(pbmailer.NewServer),
			grpc.AsService(pbnotificator.NewServer),
			grpc.AsService(pbqualifications.NewServer),
			grpc.AsService(pbrector.NewServer),
			grpc.AsService(pbstats.NewServer),
			grpc.AsService(pbwiki.NewServer),
			grpc.AsService(pbinternet.NewServer),
		),

		fx.Invoke(func(*bluemonday.Policy) {}),
		fx.Invoke(func(admin.AdminServer) {}),
		fx.Invoke(func(croner.ICron) {}),
	}
}

func main() {
	// https://github.com/DataDog/go-profiler-notes/blob/main/block.md#overhead
	// Thanks, to the authors of this document!
	runtime.SetBlockProfileRate(20000)
	runtime.SetMutexProfileFraction(100)

	ctx := kong.Parse(&cli)

	// Cli flag overrides env var
	if cli.Config != "" {
		if err := os.Setenv("FIVENET_CONFIG_FILE", cli.Config); err != nil {
			panic(err)
		}
	}
	if cli.StartTimeout <= 0 {
		cli.StartTimeout = 180 * time.Second
	}

	err := ctx.Run(&Context{})
	ctx.FatalIfErrorf(err)
}
