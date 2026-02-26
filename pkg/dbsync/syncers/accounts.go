package syncers

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

type AccountsSync struct {
	*Syncer

	state *dbsyncconfig.TableSyncState
}

func NewAccountsSync(
	s *Syncer,
	state *dbsyncconfig.TableSyncState,
) *AccountsSync {
	return &AccountsSync{
		Syncer: s,
		state:  state,
	}
}

func (s *AccountsSync) Sync(ctx context.Context) (int64, error) {
	accounts, err := s.fetchAccounts(ctx)
	if err != nil {
		return 0, err
	}

	s.logger.Debug("accountsSync", zap.Int("len", len(accounts)))

	if len(accounts) == 0 {
		return 0, nil
	}

	count := int64(len(accounts))
	if count == 0 {
		return 0, nil
	}

	// Sync jobs to FiveNet server
	if err := s.sendData(ctx, &pbsync.SendDataRequest{
		Data: &pbsync.SendDataRequest_Accounts{
			Accounts: &syncdata.DataAccounts{
				AccountUpdates: accounts,
			},
		},
	}); err != nil {
		return 0, err
	}

	s.state.Set(0, nil)

	return count, nil
}

func (s *AccountsSync) fetchAccounts(ctx context.Context) ([]*syncactivity.AccountUpdate, error) {
	limit := s.cfg.Limits.Accounts
	q := s.cfg.Tables.Accounts.GetQuery(s.state, 0, limit)
	s.logger.Debug("accounts sync query", zap.String("query", q))

	accountsResults := []*syncactivity.AccountUpdate{}
	if _, err := qrm.Query(ctx, s.db, q, []any{}, &accountsResults); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("failed to query accounts. %w", err)
		}
	}

	for _, acc := range accountsResults {
		if acc.GetGroup() != "" && acc.GetGroups() != nil {
			// If both Group and Groups are set, we prioritize Groups and add the single Group to it
			acc.GetGroups().AddGroup(acc.GetGroup())
			acc.Group = nil
		}
	}

	return accountsResults, nil
}
