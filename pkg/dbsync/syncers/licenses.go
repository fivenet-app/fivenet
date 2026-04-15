package syncers

import (
	"context"
	"errors"

	userslicenses "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/citizens/licenses"
	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	dbsyncconfig "github.com/fivenet-app/fivenet/v2026/pkg/dbsync/config"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/zap"
)

type LicensesSync struct {
	*Syncer
}

func NewLicensesSync(s *Syncer, state *dbsyncconfig.TableSyncState) *LicensesSync {
	return &LicensesSync{
		Syncer: s,
	}
}

func (s *LicensesSync) Sync(ctx context.Context) (int64, error) {
	limit := s.cfg.Limits.Licenses

	q := s.cfg.Tables.Licenses.GetQuery(limit)
	s.logger.Debug("licenses sync query", zap.String("query", q))

	licenses := []*userslicenses.License{}
	if _, err := qrm.Query(ctx, s.db, q, []any{}, &licenses); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, err
		}
	}

	count := int64(len(licenses))
	s.logger.Debug("licensesSync", zap.Int64("len", count))
	if count == 0 {
		return 0, nil
	}

	// Sync licenses to FiveNet server
	req := &pbsync.SendLicensesRequest{
		Licenses: licenses,
	}
	if err := s.send(ctx, req, func(ctx context.Context, cli pbsync.SyncServiceClient) error {
		_, err := cli.SendLicenses(ctx, req)
		return err
	}); err != nil {
		return 0, err
	}

	return count, nil
}
