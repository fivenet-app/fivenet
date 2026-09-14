package appconfig

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/settings"
	serverconfig "github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestUpdateConfigInDBUpsertsConfigAndMarksSetupComplete(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	mock.ExpectExec(
		regexp.QuoteMeta(
			`INSERT INTO fivenet_config`,
		) + `(?s).*` + regexp.QuoteMeta(
			`setup_complete`,
		) + `(?s).*` + regexp.QuoteMeta(
			`ON DUPLICATE KEY UPDATE`,
		),
	).WillReturnResult(sqlmock.NewResult(0, 1))

	c := &Config{db: db}
	require.NoError(t, c.updateConfigInDB(t.Context(), &settings.AppConfig{}))
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestApplyInitialConfigAndStartupOverride(t *testing.T) {
	t.Parallel()

	initialPrivacyPolicy := "https://example.com/privacy"
	overrideImprint := "https://override.example.com/imprint"
	c := &Config{
		logger: zap.NewNop(),
		serverConfig: &serverconfig.Config{
			AppConfig: serverconfig.AppConfig{
				Initial: serverconfig.AppConfigInitial{
					Website: serverconfig.AppConfigWebsite{
						Links: &serverconfig.AppConfigLinks{PrivacyPolicy: &initialPrivacyPolicy},
					},
				},
				Override: serverconfig.AppConfigOverride{
					Enabled: true,
					Website: serverconfig.AppConfigWebsite{
						Links: &serverconfig.AppConfigLinks{Imprint: &overrideImprint},
					},
				},
			},
		},
	}
	cfg := &settings.AppConfig{}
	cfg.Default()

	c.applyInitialConfig(cfg)
	require.True(t, cfg.GetWebsite().GetLinks().HasPrivacyPolicy())
	assert.Equal(t, initialPrivacyPolicy, cfg.GetWebsite().GetLinks().GetPrivacyPolicy())
	assert.False(t, cfg.GetWebsite().GetLinks().HasImprint())

	c.applyStartupOverride(cfg)
	assert.Equal(t, initialPrivacyPolicy, cfg.GetWebsite().GetLinks().GetPrivacyPolicy())
	require.True(t, cfg.GetWebsite().GetLinks().HasImprint())
	assert.Equal(t, overrideImprint, cfg.GetWebsite().GetLinks().GetImprint())
}
