package settingsstore

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/settings"
	"github.com/stretchr/testify/require"
)

func TestStoreUpdateAppConfigUpsertsConfig(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO fivenet_config`) + `(?s).*` + regexp.QuoteMeta(`ON DUPLICATE KEY UPDATE`)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(t, store.UpdateAppConfig(t.Context(), &settings.AppConfig{}))
	require.NoError(t, mock.ExpectationsWereMet())
}
