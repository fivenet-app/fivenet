package settingsstore

import (
	"context"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/settings"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
)

var tConfig = table.FivenetConfig

func (s *Store) UpdateAppConfig(ctx context.Context, cfg *settings.AppConfig) error {
	stmt := tConfig.
		INSERT(
			tConfig.Key,
			tConfig.AppConfig,
		).
		VALUES(
			1,
			cfg,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tConfig.AppConfig.SET(mysql.RawString("VALUES(`app_config`)")),
		)

	_, err := stmt.ExecContext(ctx, s.db)
	return err
}
