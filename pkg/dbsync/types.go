package dbsync

import (
	"database/sql"

	pbsync "github.com/fivenet-app/fivenet/gen/go/proto/services/sync"
	"go.uber.org/zap"
)

type syncer struct {
	logger *zap.Logger
	db     *sql.DB

	cfg *DBSync

	cli pbsync.SyncServiceClient
}
