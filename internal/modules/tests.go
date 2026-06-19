package modules

import (
	"time"

	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/audit"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2026/pkg/coords/postals"
	"github.com/fivenet-app/fivenet/v2026/pkg/croner"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2026/pkg/perms"
	htmlsanitizer "github.com/fivenet-app/fivenet/v2026/pkg/sanitizer/html"
	"github.com/fivenet-app/fivenet/v2026/pkg/storage"
	"github.com/fivenet-app/fivenet/v2026/pkg/userinfo"
	authstore "github.com/fivenet-app/fivenet/v2026/stores/auth"
	calendarstore "github.com/fivenet-app/fivenet/v2026/stores/calendar"
	citizensstore "github.com/fivenet-app/fivenet/v2026/stores/citizens"
	completorstore "github.com/fivenet-app/fivenet/v2026/stores/completor"
	documentsstore "github.com/fivenet-app/fivenet/v2026/stores/documents"
	jobsstore "github.com/fivenet-app/fivenet/v2026/stores/jobs"
	livemapstore "github.com/fivenet-app/fivenet/v2026/stores/livemap"
	mailerstore "github.com/fivenet-app/fivenet/v2026/stores/mailer"
	notificationsstore "github.com/fivenet-app/fivenet/v2026/stores/notifications"
	qualificationsstore "github.com/fivenet-app/fivenet/v2026/stores/qualifications"
	settingsstore "github.com/fivenet-app/fivenet/v2026/stores/settings"
	statsstore "github.com/fivenet-app/fivenet/v2026/stores/stats"
	usersstore "github.com/fivenet-app/fivenet/v2026/stores/users"
	vehiclesstore "github.com/fivenet-app/fivenet/v2026/stores/vehicles"
	wikistore "github.com/fivenet-app/fivenet/v2026/stores/wiki"
	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/fx"
)

func GetFxTestOpts(opts ...fx.Option) []fx.Option {
	to := []fx.Option{
		fx.StartTimeout(240 * time.Second),

		LoggerModule,
		config.TestModule,
		fx.Provide(appconfig.NewTest),
		htmlsanitizer.Module,
		TracerProviderModule,
		perms.Module,
		fx.Provide(TestTokenMgr),
		fx.Provide(TestAudit),
		fx.Provide(postals.NewForTests),
		auth.AuthModule,
		auth.PermsModule,
		croner.HandlersModule,
		fx.Provide(croner.NewNoopRegistry),
		fx.Provide(storage.NewNoop),

		fx.Provide(
			mstlystcdata.NewDummyEnricher,
			mstlystcdata.NewDummyUserAwareEnricher,
			mstlystcdata.NewDocumentCategories,
			mstlystcdata.NewJobs,
			mstlystcdata.NewLaws,
		),

		fx.Provide(
			authstore.New,
			calendarstore.New,
			citizensstore.New,
			completorstore.New,
			jobsstore.New,
			livemapstore.New,
			mailerstore.New,
			notificationsstore.New,
			qualificationsstore.New,
			settingsstore.New,
			statsstore.New,
			documentsstore.New,
			usersstore.New,
			vehiclesstore.New,
			wikistore.New,
		),

		fx.Invoke(func(*bluemonday.Policy) {}),
	}

	to = append(to, opts...)

	return to
}

func TestUserInfoRetriever() userinfo.UserInfoRetriever {
	return userinfo.NewMockUserInfoRetriever(map[int32]*pbuserinfo.UserInfo{})
}

func TestTokenMgr() *auth.TokenMgr {
	return auth.NewTokenMgr("")
}

func TestAudit() audit.IAuditer {
	return &audit.Noop{}
}
