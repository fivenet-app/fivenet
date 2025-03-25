package modules

import (
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/pkg/coords/postals"
	"github.com/fivenet-app/fivenet/pkg/croner"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/pkg/html/htmlsanitizer"
	"github.com/fivenet-app/fivenet/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/pkg/perms"
	"github.com/fivenet-app/fivenet/pkg/server/audit"
	"github.com/fivenet-app/fivenet/pkg/storage"
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
		perms.TestModule,
		fx.Provide(TestTokenMgr),
		fx.Provide(TestAudit),
		fx.Provide(postals.NewForTests),
		auth.AuthModule,
		auth.PermsModule,
		croner.HandlerModule,
		fx.Provide(croner.NewNoopCron),
		fx.Provide(storage.NewNoop),

		fx.Provide(
			mstlystcdata.NewDocumentCategories,
			mstlystcdata.NewJobs,
			mstlystcdata.NewLaws,
			mstlystcdata.NewEnricher,
			mstlystcdata.NewUserAwareEnricher,
		),

		fx.Invoke(func(*bluemonday.Policy) {}),
	}

	to = append(to, opts...)

	return to
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
