package main

import (
	"fmt"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/discord"
	"github.com/galexrt/fivenet/pkg/events"
	"github.com/galexrt/fivenet/pkg/grpc"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/notifi"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/pkg/server"
	"github.com/galexrt/fivenet/pkg/server/admin"
	"github.com/galexrt/fivenet/pkg/server/audit"
	"github.com/galexrt/fivenet/pkg/tracker"
	"github.com/galexrt/fivenet/pkg/tracker/postals"
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
	pbjobs "github.com/galexrt/fivenet/gen/go/proto/services/jobs"
	pblivemapper "github.com/galexrt/fivenet/gen/go/proto/services/livemapper"
	pbnotificator "github.com/galexrt/fivenet/gen/go/proto/services/notificator"
	pbrector "github.com/galexrt/fivenet/gen/go/proto/services/rector"
)

var CLI struct {
	Config string `help:"Alternative config file (env var: FIVENET_CONFIG_FILE)"`

	Server struct {
	} `cmd:"" help:"Run FiveNet server."`

	Worker struct {
	} `cmd:"" help:"Run FiveNet worker."`
}

func main() {
	ctx := kong.Parse(&CLI)

	// Cli flag always overrides env var
	if CLI.Config != "" {
		if err := os.Setenv("FIVENET_CONFIG_FILE", CLI.Config); err != nil {
			panic(err)
		}
	}

	fxOpts := []fx.Option{
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
		),

		fx.Invoke(func(admin.AdminServer) {}),
	}

	switch ctx.Command() {
	case "server":
		fxOpts = append(fxOpts,
			fx.Invoke(func(*grpcserver.Server) {}),
			fx.Invoke(func(server.HTTPServer) {}),
		)

	case "worker":
		fxOpts = append(fxOpts,
			fx.Invoke(func(*audit.Retention) {}),
			fx.Invoke(func(*discord.Bot) {}),
			fx.Invoke(func(*bot.Manager) {}),
			fx.Invoke(func(*manager.Housekeeper) {}),
		)

	default:
		panic(ctx.Error)
	}

	fx.New(fxOpts...).Run()
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
