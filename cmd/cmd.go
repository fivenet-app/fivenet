package cmd

import (
	"time"

	"github.com/alecthomas/kong"
	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	// GRPC Services

	pbauth "github.com/fivenet-app/fivenet/v2025/services/auth"
	pbcalendar "github.com/fivenet-app/fivenet/v2025/services/calendar"
	pbcentrum "github.com/fivenet-app/fivenet/v2025/services/centrum"
	pbcitizenstore "github.com/fivenet-app/fivenet/v2025/services/citizenstore"
	pbcompletor "github.com/fivenet-app/fivenet/v2025/services/completor"
	pbdmv "github.com/fivenet-app/fivenet/v2025/services/dmv"
	pbdocstore "github.com/fivenet-app/fivenet/v2025/services/docstore"
	pbinternet "github.com/fivenet-app/fivenet/v2025/services/internet"
	pbjobs "github.com/fivenet-app/fivenet/v2025/services/jobs"
	pblivemapper "github.com/fivenet-app/fivenet/v2025/services/livemapper"
	pbmailer "github.com/fivenet-app/fivenet/v2025/services/mailer"
	pbnotificator "github.com/fivenet-app/fivenet/v2025/services/notificator"
	pbqualifications "github.com/fivenet-app/fivenet/v2025/services/qualifications"
	pbrector "github.com/fivenet-app/fivenet/v2025/services/rector"
	pbstats "github.com/fivenet-app/fivenet/v2025/services/stats"
	pbsync "github.com/fivenet-app/fivenet/v2025/services/sync"
	pbwiki "github.com/fivenet-app/fivenet/v2025/services/wiki"

	// Modules
	"github.com/fivenet-app/fivenet/v2025/internal/modules"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2025/pkg/coords/postals"
	"github.com/fivenet-app/fivenet/v2025/pkg/croner"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbsync"
	"github.com/fivenet-app/fivenet/v2025/pkg/discord"
	"github.com/fivenet-app/fivenet/v2025/pkg/discord/commands"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/v2025/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/v2025/pkg/html/htmldiffer"
	"github.com/fivenet-app/fivenet/v2025/pkg/html/htmlsanitizer"
	"github.com/fivenet-app/fivenet/v2025/pkg/lang"
	"github.com/fivenet-app/fivenet/v2025/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2025/pkg/notifi"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/server"
	"github.com/fivenet-app/fivenet/v2025/pkg/server/admin"
	"github.com/fivenet-app/fivenet/v2025/pkg/server/api"
	"github.com/fivenet-app/fivenet/v2025/pkg/server/audit"
	"github.com/fivenet-app/fivenet/v2025/pkg/server/filestore"
	"github.com/fivenet-app/fivenet/v2025/pkg/server/images"
	"github.com/fivenet-app/fivenet/v2025/pkg/server/oauth2"
	"github.com/fivenet-app/fivenet/v2025/pkg/server/wk"
	"github.com/fivenet-app/fivenet/v2025/pkg/storage"
	"github.com/fivenet-app/fivenet/v2025/pkg/tracker"
	"github.com/fivenet-app/fivenet/v2025/query"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/centrumbot"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/centrumbrokers"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/centrummanager"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/centrumstate"
)

type Context struct{}

var Cli struct {
	Version kong.VersionFlag `help:"Print version information and quit"`

	Config         string        `help:"Config file path" env:"FIVENET_CONFIG_FILE"`
	StartTimeout   time.Duration `help:"App start timeout duration" default:"180s" env:"FIVENET_START_TIMEOUT"`
	SkipMigrations *bool         `help:"Disable the automatic DB migrations on startup." env:"FIVENET_SKIP_DB_MIGRATIONS"`

	Server  ServerCmd  `cmd:"" help:"Run FiveNet server."`
	Worker  WorkerCmd  `cmd:"" help:"Run FiveNet worker."`
	Discord DiscordCmd `cmd:"" help:"Run FiveNet Discord bot."`
	DBSync  DBSyncCmd  `cmd:"" name:"dbsync" help:"Run FiveNet database sync."`
	Tools   ToolsCmd   `cmd:"" help:"Run FiveNet tools/helpers."`
}

func getFxBaseOpts(startTimeout time.Duration, withServer bool) []fx.Option {
	opts := []fx.Option{
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			l := &fxevent.ZapLogger{Logger: log}
			// Show fx logs only when debug logs are enabled
			l.UseLogLevel(zap.DebugLevel)
			return l
		}),
		fx.StartTimeout(startTimeout),

		admin.Module,
		appconfig.Module,
		audit.Module,
		auth.AuthModule,
		auth.PermsModule,
		auth.TokenMgrModule,
		centrumbot.Module,
		config.Module,
		croner.AgentModule,
		croner.HandlerModule,
		croner.Module,
		croner.SchedulerModule,
		croner.StateModule,
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
		fx.Provide(centrumbrokers.New),
		storage.Module,
		housekeeper.Module,
		dbsync.Module,
		// Discord Bot
		discord.StateModule,
		discord.BotModule,
		commands.Module,
		fx.Provide(
			commands.AsCommand(commands.NewAbsentCommand),
			commands.AsCommand(commands.NewFivenetCommand),
			commands.AsCommand(commands.NewHelpCommand),
			commands.AsCommand(commands.NewSyncCommand),
		),

		fx.Provide(
			mstlystcdata.NewDocumentCategories,
			mstlystcdata.NewJobs,
			mstlystcdata.NewJobsSearch,
			mstlystcdata.NewLaws,
			mstlystcdata.NewEnricher,
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
			grpc.AsService(pbsync.NewServer),
		),

		fx.Invoke(func(*bluemonday.Policy) {}),
	}

	if withServer {
		opts = append(opts,
			fx.Invoke(func(admin.AdminServer) {}),
			fx.Invoke(func(croner.ICron) {}),
		)
	}

	return opts
}
