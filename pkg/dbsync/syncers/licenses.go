package syncers

import (
	"context"
	"errors"

	syncdata "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/sync/data"
	userslicenses "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/licenses"
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
	limit := int64(200)

	q := s.cfg.Tables.Licenses.GetQuery(0, limit)
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
	if s.cli != nil {
		if err := s.sendData(ctx, &pbsync.SendDataRequest{
			Data: &pbsync.SendDataRequest_Licenses{
				Licenses: &syncdata.DataLicenses{
					Licenses: licenses,
				},
			},
		}); err != nil {
			return 0, err
		}
	}

	return count, nil
}
