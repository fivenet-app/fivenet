package main

import (
	"fmt"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/coords/postals"
	"github.com/galexrt/fivenet/pkg/discord"
	"github.com/galexrt/fivenet/pkg/events"
	"github.com/galexrt/fivenet/pkg/grpc"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/notifi"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/pkg/sentry"
	"github.com/galexrt/fivenet/pkg/server"
	"github.com/galexrt/fivenet/pkg/server/admin"
	"github.com/galexrt/fivenet/pkg/server/audit"
	"github.com/galexrt/fivenet/pkg/storage"
	"github.com/galexrt/fivenet/pkg/tracker"
	"github.com/galexrt/fivenet/query"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	grpcserver "google.golang.org/grpc"
	// GRPC Services
	pbauth "github.com/galexrt/fivenet/gen/go/proto/services/auth"
	pbcentrum "github.com/galexrt/fivenet/gen/go/proto/services/centrum"
	"github.com/galexrt/fivenet/gen/go/proto/services/centrum/bot"
	"github.com/galexrt/fivenet/gen/go/proto/services/centrum/manager"
	"github.com/galexrt/fivenet/gen/go/proto/services/centrum/state"
	pbcitizenstore "github.com/galexrt/fivenet/gen/go/proto/services/citizenstore"
	pbcompletor "github.com/galexrt/fivenet/gen/go/proto/services/completor"
	pbdmv "github.com/galexrt/fivenet/gen/go/proto/services/dmv"
	pbdocstore "github.com/galexrt/fivenet/gen/go/proto/services/docstore"
	pbfilestore "github.com/galexrt/fivenet/gen/go/proto/services/filestore"
	pbjobs "github.com/galexrt/fivenet/gen/go/proto/services/jobs"
	pblivemapper "github.com/galexrt/fivenet/gen/go/proto/services/livemapper"
	pbnotificator "github.com/galexrt/fivenet/gen/go/proto/services/notificator"
	pbrector "github.com/galexrt/fivenet/gen/go/proto/services/rector"
)

type Context struct{}

type ServerCmd struct{}

func (c *ServerCmd) Run(ctx *Context) error {
	fxOpts := getFxBaseOpts()
	fxOpts = append(fxOpts,
		fx.Invoke(func(*grpcserver.Server) {}),
		fx.Invoke(func(server.HTTPServer) {}),
	)

	fx.New(fxOpts...).Run()

	return nil
}

type WorkerCmd struct {
	ModuleAuditRetention     bool `help:"Start Audit log retention module" default:"true"`
	ModuleDiscordBot         bool `help:"Start Discord bot module" default:"true"`
	ModuleCentrumBot         bool `help:"Start Centrum Bot module" default:"true"`
	ModuleCentrumHousekeeper bool `help:"Start Centrum Housekeeper module" default:"true"`
}

func (c *WorkerCmd) Run(ctx *Context) error {
	fxOpts := getFxBaseOpts()

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

	fx.New(fxOpts...).Run()

	return nil
}

var cli struct {
	Config string `help:"Alternative config file (env var: FIVENET_CONFIG_FILE)"`

	Server ServerCmd `cmd:"" help:"Run FiveNet server."`
	Worker WorkerCmd `cmd:"" help:"Run FiveNet worker."`
}

func getFxBaseOpts() []fx.Option {
	return []fx.Option{
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.StartTimeout(180 * time.Second),

		LoggerModule,
		config.Module,
		admin.Module,
		server.HTTPServerModule,
		grpc.ServerModule,
		server.TracerProviderModule,
		auth.AuthModule,
		auth.TokenMgrModule,
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
		sentry.Module,
		storage.Module,

		fx.Provide(
			mstlystcdata.NewCache,
			mstlystcdata.NewEnricher,
			mstlystcdata.NewSearcher,
			notifi.New,
			tracker.New,
			userinfo.NewUIRetriever,
			postals.New,
		),

		// GRPC Services
		fx.Provide(
			grpc.AsService(pbauth.NewServer),
			grpc.AsService(pbcentrum.NewServer),
			grpc.AsService(pbcitizenstore.NewServer),
			grpc.AsService(pbcompletor.NewServer),
			grpc.AsService(pbdmv.NewServer),
			grpc.AsService(pbdocstore.NewServer),
			grpc.AsService(pbjobs.NewServer),
			grpc.AsService(pblivemapper.NewServer),
			grpc.AsService(pbnotificator.NewServer),
			grpc.AsService(pbrector.NewServer),
			grpc.AsService(pbfilestore.NewServer),
		),

		fx.Invoke(func(admin.AdminServer) {}),
	}
}

func main() {
	ctx := kong.Parse(&cli)

	// Cli flag overrides env var
	if cli.Config != "" {
		if err := os.Setenv("FIVENET_CONFIG_FILE", cli.Config); err != nil {
			panic(err)
		}
	}

	err := ctx.Run(&Context{})
	ctx.FatalIfErrorf(err)
}

var LoggerModule = fx.Module("logger",
	fx.Provide(
		NewLogger,
	),
)

func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	// Logger Setup
	loggerConfig := zap.NewProductionConfig()
	level, err := zapcore.ParseLevel(cfg.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to parse log level from config. %w", err)
	}
	loggerConfig.Level.SetLevel(level)

	logger, err := loggerConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to configure logger. %w", err)
	}

	return logger, nil
}
