package syncers

import (
	"context"
	"database/sql"
	"fmt"

	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	dbsyncconfig "github.com/fivenet-app/fivenet/v2026/pkg/dbsync/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/protoutils"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
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

func (s *Syncer) send(
	ctx context.Context,
	req proto.Message,
	call func(context.Context, pbsync.SyncServiceClient) error,
) error {
	if s.cfg.Destination.DryRun {
		out, err := protoutils.MarshalToPrettyJSON(req)
		if err != nil {
			s.logger.Error("failed to marshal data for dry run output", zap.Error(err))
		} else {
			//nolint:forbidigo // This is a debug log, so it's fine to print to stdout.
			fmt.Println("Data that would have been sent to the server:", string(out))
		}
		return nil
	}

	if s.cli != nil {
		if err := call(ctx, s.cli); err != nil {
			return fmt.Errorf("failed to send data to server. %w", err)
		}
	}

	return nil
}
