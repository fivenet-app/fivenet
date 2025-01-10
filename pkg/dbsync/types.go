package dbsync

import (
	"database/sql"

	pbsync "github.com/fivenet-app/fivenet/gen/go/proto/services/sync"
	"github.com/fivenet-app/fivenet/pkg/config"
	"go.uber.org/zap"
)

type syncer struct {
	logger *zap.Logger
	db     *sql.DB

	cfg *config.DBSync

	cli pbsync.SyncServiceClient
}
