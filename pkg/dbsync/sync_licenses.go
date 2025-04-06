package dbsync

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/sync"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	pbsync "github.com/fivenet-app/fivenet/gen/go/proto/services/sync"
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

	limit := 200

	sQuery := s.cfg.Tables.Licenses
	query := prepareStringQuery(sQuery, s.state, 0, limit)

	licenses := []*users.License{}
	if _, err := qrm.Query(ctx, s.db, query, []any{}, &licenses); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	s.logger.Debug("licensesSync", zap.Any("licenses", licenses))

	if len(licenses) == 0 {
		return nil
	}

	// Sync licenses to FiveNet server
	if s.cli != nil {
		if _, err := s.cli.SendData(ctx, &pbsync.SendDataRequest{
			Data: &pbsync.SendDataRequest_Licenses{
				Licenses: &sync.DataLicenses{
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
