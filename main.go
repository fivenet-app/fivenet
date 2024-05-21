package main

import (
	"os"
	"runtime"
	"time"

	"github.com/alecthomas/kong"
	"github.com/fivenet-app/fivenet/internal/modules"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/pkg/coords/postals"
	"github.com/fivenet-app/fivenet/pkg/discord"
	"github.com/fivenet-app/fivenet/pkg/events"
	"github.com/fivenet-app/fivenet/pkg/grpc"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/pkg/htmlsanitizer"
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
	"github.com/fivenet-app/fivenet/pkg/storage"
	"github.com/fivenet-app/fivenet/pkg/tracker"
	"github.com/fivenet-app/fivenet/query"
	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	grpcserver "google.golang.org/grpc"
	// GRPC Services
	pbauth "github.com/fivenet-app/fivenet/gen/go/proto/services/auth"
	pbcalendar "github.com/fivenet-app/fivenet/gen/go/proto/services/calendar"
	pbcentrum "github.com/fivenet-app/fivenet/gen/go/proto/services/centrum"
	pbcitizenstore "github.com/fivenet-app/fivenet/gen/go/proto/services/citizenstore"
	pbcompletor "github.com/fivenet-app/fivenet/gen/go/proto/services/completor"
	pbdmv "github.com/fivenet-app/fivenet/gen/go/proto/services/dmv"
	pbdocstore "github.com/fivenet-app/fivenet/gen/go/proto/services/docstore"
	pbjobs "github.com/fivenet-app/fivenet/gen/go/proto/services/jobs"
	pblivemapper "github.com/fivenet-app/fivenet/gen/go/proto/services/livemapper"
	pbmessenger "github.com/fivenet-app/fivenet/gen/go/proto/services/messenger"
	pbnotificator "github.com/fivenet-app/fivenet/gen/go/proto/services/notificator"
	pbqualifications "github.com/fivenet-app/fivenet/gen/go/proto/services/qualifications"
	pbrector "github.com/fivenet-app/fivenet/gen/go/proto/services/rector"

	// Modules
	"github.com/fivenet-app/fivenet/gen/go/proto/services/centrum/bot"
	"github.com/fivenet-app/fivenet/gen/go/proto/services/centrum/manager"
	"github.com/fivenet-app/fivenet/gen/go/proto/services/centrum/state"
)

type Context struct{}

type FrontendCmd struct{}

func (c *FrontendCmd) Run(ctx *Context) error {
	fxOpts := getFxBaseOpts(cli.StartTimeout)
	fxOpts = append(fxOpts,
		fx.Invoke(func(server.HTTPServer) {}),
	)

	app := fx.New(fxOpts...)
	app.Run()

	return nil
}

type ServerCmd struct {
	ServeFrontend bool `help:"Serve HTTP Frontend."`
}

func (c *ServerCmd) Run(ctx *Context) error {
	fxOpts := getFxBaseOpts(cli.StartTimeout)
	fxOpts = append(fxOpts, fx.Invoke(func(*grpcserver.Server) {}))

	if c.ServeFrontend {
		fxOpts = append(fxOpts, fx.Invoke(func(server.HTTPServer) {}))
	}

	app := fx.New(fxOpts...)
	app.Run()

	return nil
}

type WorkerCmd struct {
	ModuleAuditRetention     bool `help:"Start Audit log retention module" default:"true"`
	ModuleDiscordBot         bool `help:"Start Discord bot module" default:"true"`
	ModuleCentrumBot         bool `help:"Start Centrum bot module" default:"true"`
	ModuleCentrumHousekeeper bool `help:"Start Centrum Housekeeper module" default:"true"`
	ModuleUserTracker        bool `help:"Start User tracker module" default:"true"`
}

func (c *WorkerCmd) Run(ctx *Context) error {
	fxOpts := getFxBaseOpts(cli.StartTimeout)

	if c.ModuleAuditRetention {
		fxOpts = append(fxOpts, fx.Invoke(func(*audit.Retention) {}))
	}
	if c.ModuleCentrumBot {
		fxOpts = append(fxOpts, fx.Invoke(func(*bot.Manager) {}))
	}
	if c.ModuleCentrumHousekeeper {
		fxOpts = append(fxOpts, fx.Invoke(func(*manager.Housekeeper) {}))
	}
	if c.ModuleDiscordBot {
		fxOpts = append(fxOpts, fx.Invoke(func(*discord.Bot) {}))
	}
	if c.ModuleUserTracker {
		fxOpts = append(fxOpts, fx.Invoke(func(*tracker.Manager) {}))
	}

	app := fx.New(fxOpts...)
	app.Run()

	return nil
}

var cli struct {
	Config       string        `help:"Alternative config file (env var: FIVENET_CONFIG_FILE)"`
	StartTimeout time.Duration `help:"App start timeout duration"`

	Frontend FrontendCmd `cmd:"" help:"Run FiveNet frontend."`
	Server   ServerCmd   `cmd:"" help:"Run FiveNet server."`
	Worker   WorkerCmd   `cmd:"" help:"Run FiveNet worker."`
}

func getFxBaseOpts(startTimeout time.Duration) []fx.Option {
	return []fx.Option{
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.StartTimeout(startTimeout),

		modules.LoggerModule,
		htmlsanitizer.Module,
		config.Module,
		appconfig.Module,
		modules.TracerProviderModule,
		admin.Module,
		server.HTTPEngineModule,
		server.HTTPServerModule,
		grpc.ServerModule,
		auth.AuthModule,
		auth.TokenMgrModule,
		auth.PermsModule,
		query.Module,
		perms.Module,
		events.Module,
		audit.Module,
		audit.RetentionModule,
		state.StateModule,
		bot.Module,
		manager.Module,
		manager.HousekeeperModule,
		discord.BotModule,
		storage.Module,

		fx.Provide(
			mstlystcdata.NewCache,
			mstlystcdata.NewEnricher,
			mstlystcdata.NewUserAwareEnricher,
			mstlystcdata.NewSearcher,
			notifi.New,
			tracker.NewManager,
			tracker.New,
			userinfo.NewUIRetriever,
			postals.New,

			// HTTP Services
			server.AsService(api.New),
			server.AsService(oauth2.New),
			server.AsService(images.New),
			server.AsService(filestore.New),
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
			grpc.AsService(pbmessenger.NewServer),
			grpc.AsService(pbnotificator.NewServer),
			grpc.AsService(pbqualifications.NewServer),
			grpc.AsService(pbrector.NewServer),
		),

		fx.Invoke(func(*bluemonday.Policy) {}),
		fx.Invoke(func(admin.AdminServer) {}),
	}
}

func main() {
	// https://github.com/DataDog/go-profiler-notes/blob/main/block.md#overhead
	// Thanks, to the authors of this document!
	runtime.SetBlockProfileRate(20000)

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
