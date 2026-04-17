package syncers

import (
	"context"
	"errors"
	"fmt"

	syncactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/sync/activity"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	dbsyncconfig "github.com/fivenet-app/fivenet/v2026/pkg/dbsync/config"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/zap"
)

const maxAccountsPerSendRequeust = 100

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
	limit := s.cfg.Limits.Accounts
	var total int64
	prevCursorID := ""

	for batches := 0; ; batches++ {
		accounts, cursorID, err := s.fetchAccounts(ctx)
		if err != nil {
			return total, err
		}

		fetched := int64(len(accounts))
		s.logger.Debug("accountsSync", zap.Int64("len", fetched))
		if fetched == 0 {
			break
		}

		// Sync accounts to FiveNet server (in batches if higher than API limit).
		for start := 0; start < len(accounts); start += sync.MaxAccountsPerRequest {
			end := min(start+maxUsersPerSendRequest, len(accounts))
			req := &pbsync.SendAccountsRequest{
				AccountUpdates: accounts[start:end],
			}
			if err := s.send(
				ctx,
				req,
				func(ctx context.Context, cli pbsync.SyncServiceClient) error {
					_, err := cli.SendAccounts(ctx, req)
					return err
				},
			); err != nil {
				return total, err
			}
		}
		total += fetched

		// Nothing left for this cycle.
		if fetched < limit {
			break
		}

		// Guard against starvation when data changes continuously under high write load.
		if batches+1 >= maxDrainBatchesPerSync {
			s.logger.Info(
				"accounts sync hit drain batch cap; remaining updates continue next interval",
				zap.Int64("fetched", fetched),
				zap.Stringp("cursor_id", cursorID),
			)
			break
		}

		// Guard against non-advancing cursor loops.
		if cursorID == nil || *cursorID == "" {
			s.logger.Info(
				"accounts sync cursor missing, stopping drain loop",
				zap.Int64("fetched", fetched),
			)
			break
		}
		if *cursorID == prevCursorID {
			s.logger.Info(
				"accounts sync cursor did not advance, stopping drain loop",
				zap.String("cursor_id", *cursorID),
				zap.Int64("fetched", fetched),
			)
			break
		}
		prevCursorID = *cursorID
		s.state.SetCursor(nil, cursorID)
	}

	// Accounts sync intentionally processes full snapshots each interval.
	s.state.ResetCursor()

	return total, nil
}

func (s *AccountsSync) fetchAccounts(
	ctx context.Context,
) ([]*syncactivity.AccountUpdate, *string, error) {
	limit := s.cfg.Limits.Accounts
	q := s.cfg.Tables.Accounts.GetQuery(s.state, 0, limit)
	s.logger.Debug("accounts sync query", zap.String("query", q))

	accountsResults := []*syncactivity.AccountUpdate{}
	if _, err := qrm.Query(ctx, s.db, q, []any{}, &accountsResults); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, nil, fmt.Errorf("failed to query accounts. %w", err)
		}
	}

	for _, acc := range accountsResults {
		if acc.GetGroup() != "" && acc.GetGroups() != nil {
			// If both Group and Groups are set, we prioritize Groups and add the single Group to it
			acc.GetGroups().AddGroup(acc.GetGroup())
			acc.Group = nil
		}
	}

	var cursorID *string
	if len(accountsResults) > 0 {
		lastID := accountsResults[len(accountsResults)-1].GetLicense()
		if lastID != "" {
			cursorID = &lastID
		}
	}

	return accountsResults, cursorID, nil
}
