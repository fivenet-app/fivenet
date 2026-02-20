package syncers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	dbsyncconfig "github.com/fivenet-app/fivenet/v2026/pkg/dbsync/config"
	"go.uber.org/zap"
)

type Syncer struct {
	logger *zap.Logger
	db     *sql.DB

	cfg *dbsyncconfig.DBSyncConfig

	cli pbsync.SyncServiceClient
}

func New(
	logger *zap.Logger,
	db *sql.DB,
	cfg *dbsyncconfig.DBSyncConfig,
	cli pbsync.SyncServiceClient,
) *Syncer {
	return &Syncer{
		logger: logger,
		db:     db,
		cfg:    cfg,
		cli:    cli,
	}
}

func (s *Syncer) sendData(ctx context.Context, data *pbsync.SendDataRequest) error {
	if s.cfg.Destination.DryRun {
		s.logger.Info("dry run enabled, not sending data to server")
		out, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			s.logger.Error("failed to marshal data for dry run output", zap.Error(err))
		} else {
			fmt.Println("Data that would have been sent to the server:", string(out))
		}
		return nil
	}

	if s.cli != nil {
		if _, err := s.cli.SendData(ctx, data); err != nil {
			return fmt.Errorf("failed to send data to server. %w", err)
		}
	}

	return nil
}
