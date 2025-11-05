package dbsync

import (
	"context"
	"database/sql"
	"fmt"

	pbsync "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/sync"
	"go.uber.org/zap"
)

type syncer struct {
	logger *zap.Logger
	db     *sql.DB

	cfg *DBSyncConfig

	cli pbsync.SyncServiceClient
}

func (s *syncer) sendData(ctx context.Context, data *pbsync.SendDataRequest) error {
	if s.cli != nil {
		if _, err := s.cli.SendData(ctx, data); err != nil {
			return fmt.Errorf("failed to send data to server. %w", err)
		}
	}
	return nil
}
