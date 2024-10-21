package modules

import (
	"database/sql"
	"time"

	"github.com/fivenet-app/fivenet/internal/tests/servers"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/pkg/coords/postals"
	"github.com/fivenet-app/fivenet/pkg/croner"
	"github.com/fivenet-app/fivenet/pkg/events"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/pkg/htmlsanitizer"
	"github.com/fivenet-app/fivenet/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/pkg/perms"
	"github.com/fivenet-app/fivenet/pkg/server/audit"
	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/fx"
)

func GetFxTestOpts(opts ...fx.Option) []fx.Option {
	to := []fx.Option{
		LoggerModule,
		config.TestModule,
		fx.Provide(appconfig.NewTest),
		htmlsanitizer.Module,
		TracerProviderModule,
		perms.Module,
		fx.Provide(TestDBClient),
		fx.Provide(TestJetStreamClient),
		fx.Provide(TestUserInfoRetriever),
		fx.Provide(TestTokenMgr),
		fx.Provide(TestAudit),
		fx.Provide(postals.NewForTests),
		auth.AuthModule,
		auth.PermsModule,
		croner.HandlerModule,
		fx.Provide(croner.NewNoopCron),

		fx.Provide(
			mstlystcdata.NewCache,
			mstlystcdata.NewEnricher,
			mstlystcdata.NewUserAwareEnricher,
			mstlystcdata.NewSearcher,
		),

		fx.Invoke(func(*bluemonday.Policy) {}),
	}

	to = append(to, opts...)

	return to
}

func TestConfig() (*config.Config, error) {
	cfg, err := config.LoadTestConfig()

	if cfg != nil {
		cfg.NATS.URL = servers.TestNATSServer.GetURL()
		cfg.Cache.RefreshTime = 1 * time.Hour
	}

	return cfg, err
}

func TestDBClient() (*sql.DB, error) {
	db, err := servers.TestDBServer.DB()

	return db, err
}

func TestJetStreamClient() (*events.JSWrapper, error) {
	js, err := servers.TestNATSServer.GetJS()
	if err != nil {
		return nil, err
	}

	return &events.JSWrapper{
		JetStream: js,
	}, nil
}

func TestUserInfoRetriever() userinfo.UserInfoRetriever {
	return userinfo.NewMockUserInfoRetriever(map[int32]*userinfo.UserInfo{})
}

func TestTokenMgr() *auth.TokenMgr {
	return auth.NewTokenMgr("")
}

func TestAudit() audit.IAuditer {
	return &audit.Noop{}
}
