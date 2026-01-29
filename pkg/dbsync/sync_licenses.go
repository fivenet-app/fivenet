package dbsync

import (
	"context"
	"errors"

	syncdata "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/sync/data"
	userslicenses "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/licenses"
	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/zap"
)

type licensesSync struct {
	*syncer

	state *TableSyncState
}

func newLicensesSync(s *syncer, state *TableSyncState) *licensesSync {
	return &licensesSync{
		syncer: s,
		state:  state,
	}
}

func (s *licensesSync) Sync(ctx context.Context) error {
	if !s.cfg.Tables.Licenses.Enabled {
		return nil
	}

	limit := int64(200)

	q := s.cfg.Tables.Licenses.GetQuery(s.state, 0, limit)
	s.logger.Debug("licenses sync query", zap.String("query", q))

	licenses := []*userslicenses.License{}
	if _, err := qrm.Query(ctx, s.db, q, []any{}, &licenses); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	s.logger.Debug("licensesSync", zap.Int("len", len(licenses)))

	if len(licenses) == 0 {
		return nil
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
			return err
		}
	}

	s.state.Set(0, nil)

	return nil
}
