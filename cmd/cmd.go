package cmd

import (
	"time"

	"github.com/alecthomas/kong"
	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	// GRPC Services

	pbauth "github.com/fivenet-app/fivenet/services/auth"
	pbcalendar "github.com/fivenet-app/fivenet/services/calendar"
	pbcentrum "github.com/fivenet-app/fivenet/services/centrum"
	pbcitizenstore "github.com/fivenet-app/fivenet/services/citizenstore"
	pbcompletor "github.com/fivenet-app/fivenet/services/completor"
	pbdmv "github.com/fivenet-app/fivenet/services/dmv"
	pbdocstore "github.com/fivenet-app/fivenet/services/docstore"
	pbinternet "github.com/fivenet-app/fivenet/services/internet"
	pbjobs "github.com/fivenet-app/fivenet/services/jobs"
	pblivemapper "github.com/fivenet-app/fivenet/services/livemapper"
	pbmailer "github.com/fivenet-app/fivenet/services/mailer"
	pbnotificator "github.com/fivenet-app/fivenet/services/notificator"
	pbqualifications "github.com/fivenet-app/fivenet/services/qualifications"
	pbrector "github.com/fivenet-app/fivenet/services/rector"
	pbstats "github.com/fivenet-app/fivenet/services/stats"
	pbsync "github.com/fivenet-app/fivenet/services/sync"
	pbwiki "github.com/fivenet-app/fivenet/services/wiki"

	// Modules
	"github.com/fivenet-app/fivenet/internal/modules"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/pkg/coords/postals"
	"github.com/fivenet-app/fivenet/pkg/croner"
	"github.com/fivenet-app/fivenet/pkg/dbsync"
	"github.com/fivenet-app/fivenet/pkg/discord"
	"github.com/fivenet-app/fivenet/pkg/discord/commands"
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
	"github.com/fivenet-app/fivenet/services/centrum/centrumbot"
	"github.com/fivenet-app/fivenet/services/centrum/centrummanager"
	"github.com/fivenet-app/fivenet/services/centrum/centrumstate"
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
	DBSync  DBSyncCmd  `cmd:"" help:"Run FiveNet database sync."`
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
		dbsync.Module,
		// Discord Bot
		discord.StateModule,
		discord.BotModule,
		commands.Module,

		fx.Provide(
			mstlystcdata.NewDocumentCategories,
			mstlystcdata.NewJobs,
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
