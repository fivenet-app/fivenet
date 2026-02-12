package dbsync

import (
	"context"
	"errors"
	"fmt"

	syncactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/sync/activity"
	syncdata "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/sync/data"
	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	dbsyncconfig "github.com/fivenet-app/fivenet/v2026/pkg/dbsync/config"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/zap"
)

type accountsSync struct {
	*syncer

	state *dbsyncconfig.TableSyncState
}

func newAccountsSync(s *syncer, state *dbsyncconfig.TableSyncState) *accountsSync {
	return &accountsSync{
		syncer: s,
		state:  state,
	}
}

func (s *accountsSync) Sync(ctx context.Context) error {
	accounts, err := s.fetchAccounts(ctx)
	if err != nil {
		return err
	}

	s.logger.Debug("accountsSync", zap.Int("len", len(accounts)))

	if len(accounts) == 0 {
		return nil
	}

	// Log a warning when no jobs are left after filtering
	if len(accounts) == 0 {
		s.logger.Warn("no jobs left after filtering")
		return nil
	}

	// Sync jobs to FiveNet server
	if s.cli != nil {
		if err := s.sendData(ctx, &pbsync.SendDataRequest{
			Data: &pbsync.SendDataRequest_Accounts{
				Accounts: &syncdata.DataAccounts{
					AccountUpdates: accounts,
				},
			},
		}); err != nil {
			return err
		}
	}

	s.state.Set(0, nil)

	return nil
}

func (s *accountsSync) fetchAccounts(ctx context.Context) ([]*syncactivity.AccountUpdate, error) {
	limit := int64(200)
	q := s.cfg.Tables.Accounts.GetQuery(s.state, 0, limit)
	s.logger.Debug("accounts sync query", zap.String("query", q))

	accountsResults := []*syncactivity.AccountUpdate{}
	if _, err := qrm.Query(ctx, s.db, q, []any{}, &accountsResults); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("failed to query accounts. %w", err)
		}
	}

	return accountsResults, nil
}
