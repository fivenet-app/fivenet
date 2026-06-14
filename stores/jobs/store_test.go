package jobs

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/stretchr/testify/require"
)

func newTestStore(t *testing.T) (*Store, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	return &Store{db: db, customDB: &config.CustomDB{}}, mock
}
